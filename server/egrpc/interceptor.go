package egrpc

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"

	"github.com/eezz10001/ego/core/eerrors"
	"github.com/eezz10001/ego/core/elog"
	"github.com/eezz10001/ego/core/emetric"
	"github.com/eezz10001/ego/core/etrace"
	"github.com/eezz10001/ego/core/transport"
	"github.com/eezz10001/ego/core/util/xstring"
	"github.com/eezz10001/ego/internal/ecode"
	"github.com/eezz10001/ego/internal/egrpcinteceptor"
	"github.com/eezz10001/ego/internal/tools"
	"github.com/eezz10001/ego/internal/xcpu"
)

func prometheusUnaryServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	startTime := time.Now()
	serviceName, _ := egrpcinteceptor.SplitMethodName(info.FullMethod)
	emetric.ServerStartedCounter.Inc(emetric.TypeGRPCUnary, info.FullMethod, getPeerName(ctx), serviceName)
	resp, err := handler(ctx, req)
	statusInfo := ecode.Convert(err)

	emetric.ServerHandleHistogram.ObserveWithExemplar(time.Since(startTime).Seconds(), prometheus.Labels{
		"tid": etrace.ExtractTraceID(ctx),
	}, emetric.TypeGRPCUnary, info.FullMethod, getPeerName(ctx), serviceName)
	emetric.ServerHandleCounter.Inc(emetric.TypeGRPCUnary, info.FullMethod, getPeerName(ctx), statusInfo.Code().String(), http.StatusText(ecode.GrpcToHTTPStatusCode(statusInfo.Code())), serviceName)
	return resp, err
}

func prometheusStreamServerInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	startTime := time.Now()
	serviceName, _ := egrpcinteceptor.SplitMethodName(info.FullMethod)
	emetric.ServerStartedCounter.Inc(emetric.TypeGRPCStream, info.FullMethod, getPeerName(ss.Context()), serviceName)
	err := handler(srv, ss)
	statusInfo := ecode.Convert(err)
	emetric.ServerHandleHistogram.Observe(time.Since(startTime).Seconds(), emetric.TypeGRPCStream, info.FullMethod, getPeerName(ss.Context()), serviceName)
	emetric.ServerHandleCounter.Inc(emetric.TypeGRPCStream, info.FullMethod, getPeerName(ss.Context()), statusInfo.Message(), http.StatusText(ecode.GrpcToHTTPStatusCode(statusInfo.Code())), serviceName)
	return err
}

func traceUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	tracer := etrace.NewTracer(trace.SpanKindServer)
	attrs := []attribute.KeyValue{
		egrpcinteceptor.RPCSystemGRPC,
		egrpcinteceptor.GRPCKindUnary,
	}
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (reply interface{}, err error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			md = metadata.New(nil)
		}
		// Deprecated ???????????????v0.9.0??????
		etrace.CompatibleExtractGrpcTraceID(md)
		ctx, span := tracer.Start(ctx, info.FullMethod, transport.GrpcHeaderCarrier(md), trace.WithAttributes(attrs...))
		span.SetAttributes(
			semconv.RPCMethodKey.String(info.FullMethod),
			semconv.NetPeerNameKey.String(getPeerName(ctx)),
			semconv.NetPeerIPKey.String(getPeerIP(ctx)),
		)
		defer func() {
			if err != nil {
				span.RecordError(err)
				if e := eerrors.FromError(err); e != nil {
					span.SetAttributes(semconv.RPCGRPCStatusCodeKey.Int64(int64(e.Code)))
				}
				span.SetStatus(codes.Error, err.Error())
			} else {
				span.SetStatus(codes.Ok, "OK")
			}
			span.End()
		}()
		return handler(ctx, req)
	}
}

type contextedServerStream struct {
	grpc.ServerStream
	ctx context.Context

	receivedMessageID int
	sentMessageID     int
}

func (css *contextedServerStream) RecvMsg(m interface{}) error {
	err := css.ServerStream.RecvMsg(m)

	if err == nil {
		css.receivedMessageID++
		egrpcinteceptor.MessageReceived.Event(css.Context(), css.receivedMessageID, m)
	}

	return err
}

func (css *contextedServerStream) SendMsg(m interface{}) error {
	err := css.ServerStream.SendMsg(m)

	css.sentMessageID++
	egrpcinteceptor.MessageSent.Event(css.Context(), css.sentMessageID, m)

	return err
}

// Context ...
func (css *contextedServerStream) Context() context.Context {
	return css.ctx
}

func traceStreamServerInterceptor() grpc.StreamServerInterceptor {
	tracer := etrace.NewTracer(trace.SpanKindServer)
	attrs := []attribute.KeyValue{
		semconv.RPCSystemKey.String("grpc"),
		egrpcinteceptor.GRPCKindStream,
	}
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		md, ok := metadata.FromIncomingContext(ss.Context())
		if !ok {
			md = metadata.New(nil)
		}
		// Deprecated ???????????????v0.9.0??????
		etrace.CompatibleExtractGrpcTraceID(md)
		ctx, span := tracer.Start(ss.Context(), info.FullMethod, transport.GrpcHeaderCarrier(md), trace.WithAttributes(attrs...))
		span.SetAttributes(
			semconv.RPCMethodKey.String(info.FullMethod),
			semconv.NetPeerNameKey.String(getPeerName(ctx)),
			semconv.NetPeerIPKey.String(getPeerIP(ctx)),
			etrace.CustomTag("rpc.grpc.kind", "stream"),
		)
		defer span.End()
		err := handler(srv, &contextedServerStream{
			ServerStream: ss,
			ctx:          ctx,
		})
		if err != nil {
			span.RecordError(err)
			if e := eerrors.FromError(err); e != nil {
				span.SetAttributes(semconv.RPCGRPCStatusCodeKey.Int64(int64(e.Code)))
			}
			span.SetStatus(codes.Error, err.Error())
		} else {
			span.SetStatus(codes.Ok, "OK")
		}
		return err
	}
}

func (c *Container) defaultStreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		var beg = time.Now()
		var fields = make([]elog.Field, 0, 20)
		var event = "normal"
		defer func() {
			cost := time.Since(beg)
			if rec := recover(); rec != nil {
				switch rec := rec.(type) {
				case error:
					err = rec
				default:
					err = fmt.Errorf("%v", rec)
				}
				stack := make([]byte, 4096)
				stack = stack[:runtime.Stack(stack, true)]
				fields = append(fields, elog.FieldStack(stack))
				event = "recover"
			}
			spbStatus := ecode.Convert(err)
			httpStatusCode := ecode.GrpcToHTTPStatusCode(spbStatus.Code())

			fields = append(fields,
				elog.FieldType("stream"),
				elog.FieldEvent(event),
				elog.FieldCode(int32(spbStatus.Code())),
				elog.FieldUniformCode(int32(httpStatusCode)),
				elog.FieldDescription(spbStatus.Message()),
				elog.FieldCost(time.Since(beg)),
				elog.FieldPeerName(getPeerName(stream.Context())),
				elog.FieldPeerIP(getPeerIP(stream.Context())),
			)

			if c.config.SlowLogThreshold > time.Duration(0) && c.config.SlowLogThreshold < cost {
				c.logger.Warn("slow", fields...)
			}

			if err != nil {
				fields = append(fields, elog.FieldErr(err))
				// ???????????????????????????
				if httpStatusCode >= http.StatusInternalServerError {
					// ???????????????????????????
					c.logger.Error("access", fields...)
				} else {
					// ?????????????????????warning
					c.logger.Warn("access", fields...)
				}
				return
			}
			c.logger.Info("access", fields...)
		}()
		return handler(srv, stream)
	}
}

func (c *Container) defaultUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (res interface{}, err error) {
		// ??????????????????????????????
		if c.config.EnableSkipHealthLog && info.FullMethod == "/grpc.health.v1.Health/Check" {
			return handler(ctx, req)
		}

		var beg = time.Now()
		// ????????????????????????????????????????????????????????????slice??????
		loggerKeys := transport.CustomContextKeys()
		var fields = make([]elog.Field, 0, 20+transport.CustomContextKeysLength())
		var event = "normal"

		// ?????????defer?????????????????????????????????ctx
		for _, key := range loggerKeys {
			if value := tools.GrpcHeaderValue(ctx, key); value != "" {
				ctx = transport.WithValue(ctx, key, value)
			}
		}

		// ??????????????????defer???recover handler?????????????????????panic
		defer func() {
			cost := time.Since(beg)
			if rec := recover(); rec != nil {
				switch recType := rec.(type) {
				case error:
					err = recType
				default:
					err = fmt.Errorf("%v", rec)
				}

				stack := make([]byte, 4096)
				stack = stack[:runtime.Stack(stack, true)]
				fields = append(fields, elog.FieldStack(stack))
				event = "recover"
			}

			isSlow := false
			if c.config.SlowLogThreshold > time.Duration(0) && c.config.SlowLogThreshold < cost {
				isSlow = true
			}

			// ?????????????????????????????????????????????????????????????????????????????????????????????????????????
			if err == nil && !c.config.EnableAccessInterceptor && !isSlow {
				return
			}

			spbStatus := ecode.Convert(err)
			httpStatusCode := ecode.GrpcToHTTPStatusCode(spbStatus.Code())

			fields = append(fields,
				elog.FieldType("unary"),
				elog.FieldCode(int32(spbStatus.Code())),
				elog.FieldUniformCode(int32(httpStatusCode)),
				elog.FieldDescription(spbStatus.Message()),
				elog.FieldEvent(event),
				elog.FieldMethod(info.FullMethod),
				elog.FieldCost(time.Since(beg)),
				elog.FieldPeerName(getPeerName(ctx)),
				elog.FieldPeerIP(getPeerIP(ctx)),
			)

			for _, key := range loggerKeys {
				if value := tools.ContextValue(ctx, key); value != "" {
					fields = append(fields, elog.FieldCustomKeyValue(key, value))
				}
			}

			if c.config.EnableTraceInterceptor && etrace.IsGlobalTracerRegistered() {
				fields = append(fields, elog.FieldTid(etrace.ExtractTraceID(ctx)))
			}

			if c.config.EnableAccessInterceptorReq {
				var reqMap = map[string]interface{}{
					"payload": xstring.JSON(req),
				}
				if md, ok := metadata.FromIncomingContext(ctx); ok {
					reqMap["metadata"] = md
				}
				fields = append(fields, elog.Any("req", reqMap))
			}
			if c.config.EnableAccessInterceptorRes {
				fields = append(fields, elog.Any("res", map[string]interface{}{
					"payload": xstring.JSON(res),
				}))
			}

			if isSlow {
				c.logger.Warn("slow", fields...)
			}

			if err != nil {
				fields = append(fields, elog.FieldErr(err))
				// ???????????????????????????
				if httpStatusCode >= http.StatusInternalServerError {
					// ???????????????????????????
					c.logger.Error("access", fields...)
				} else {
					// ?????????????????????warning
					c.logger.Warn("access", fields...)
				}
				return
			}

			if c.config.EnableAccessInterceptor {
				c.logger.Info("access", fields...)
			}
		}()

		if enableCPUUsage(ctx) {
			var stat = xcpu.Stat{}
			xcpu.ReadStat(&stat)
			if stat.Usage > 0 {
				// https://github.com/grpc/grpc-go/blob/master/Documentation/grpc-metadata.md
				header := metadata.Pairs("cpu-usage", strconv.Itoa(int(stat.Usage)))
				err = grpc.SetHeader(ctx, header)
				if err != nil {
					c.logger.Error("set header error", elog.FieldErr(err))
				}
			}
		}
		return handler(ctx, req)
	}
}

// enableCPUUsage ????????????cpu?????????
func enableCPUUsage(ctx context.Context) bool {
	return tools.GrpcHeaderValue(ctx, "enable-cpu-usage") == "true"
}

// getPeerName ????????????????????????
func getPeerName(ctx context.Context) string {
	return tools.GrpcHeaderValue(ctx, "app")
}

// getPeerIP ????????????ip
func getPeerIP(ctx context.Context) string {
	clientIP := tools.GrpcHeaderValue(ctx, "client-ip")
	if clientIP != "" {
		return clientIP
	}

	// ???grpc????????????ip
	pr, ok2 := peer.FromContext(ctx)
	if !ok2 {
		return ""
	}
	if pr.Addr == net.Addr(nil) {
		return ""
	}
	addSlice := strings.Split(pr.Addr.String(), ":")
	if len(addSlice) > 1 {
		return addSlice[0]
	}
	return ""
}
