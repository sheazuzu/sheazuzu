/*
 *  trace.go
 *  Created on 22.02.2021
 *  Copyright (C) 2021 Volkswagen AG, All rights reserved.
 */

// Package tracing provides functions for a consistent and simple trace handling within the NGW application landscape.
package tracing

import (
	"bytes"
	"context"
	"contrib.go.opencensus.io/exporter/jaeger"
	"contrib.go.opencensus.io/exporter/jaeger/propagation"
	"go.opencensus.io/trace"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
)

var jaegerPropagation = &propagation.HTTPFormat{}

type TraceID struct{}

// RegisterTraceExporter creates and registers a new jaeger exporter with the provided params set as options.
func RegisterTraceExporter(agentURI, collectorURI, serviceName string) (*jaeger.Exporter, error) {
	// Port details: https://www.jaegertracing.io/docs/getting-started/
	je, err := jaeger.NewExporter(jaeger.Options{
		AgentEndpoint:     agentURI,
		CollectorEndpoint: collectorURI,
		Process: jaeger.Process{
			ServiceName: serviceName,
		},
	})
	if err != nil {
		return nil, err
	}

	// And now finally register it as a Trace Exporter
	trace.RegisterExporter(je)

	trace.ApplyConfig(trace.Config{
		DefaultSampler: trace.AlwaysSample(),
	})

	return je, nil
}

// StartSpan creates a span that is used to trace operations.
// It returns the new Context and the created span.
func StartSpan(ctx context.Context, name string) (context.Context, *trace.Span) {
	return trace.StartSpan(ctx, name)
}

// SpanFromContext returns the span stored in the given context.
func SpanFromContext(ctx context.Context) *trace.Span {
	return trace.FromContext(ctx)
}

// AddTraceIdToContext adds the id of the given span to the provided context und returns the context.
func AddTraceIdToContext(ctx context.Context, span *trace.Span) context.Context {
	spanContext := span.SpanContext()
	ctx = context.WithValue(ctx, TraceID{}, spanContext.TraceID.String())
	return ctx
}

// TraceHandler is a middleware that adds tracing to a http.Handler.
func TraceHandler(logger *zap.SugaredLogger, name string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			sw := statusWriter{ResponseWriter: w, status: 200}

			var ctx context.Context
			var span *trace.Span

			spanContext, ok := jaegerPropagation.SpanContextFromRequest(r)
			if ok {
				ctx, span = trace.StartSpanWithRemoteParent(r.Context(), name, spanContext, trace.WithSpanKind(trace.SpanKindServer))
			} else {
				ctx, span = trace.StartSpan(r.Context(), name, trace.WithSpanKind(trace.SpanKindServer))
			}
			ctx = AddTraceIdToContext(ctx, span)

			defer span.End()

			// SASVC-2017: Add the trace id to the response header
			w.Header().Add("X-Trace-ID", span.SpanContext().TraceID.String())

			span.AddAttributes(trace.StringAttribute("span.kind", "server"))
			span.AddAttributes(trace.StringAttribute("url", r.URL.String()))
			span.AddAttributes(trace.StringAttribute("method", r.Method))
			span.AddAttributes(trace.Int64Attribute("http.status.code", int64(sw.status)))

			tenant := r.URL.Query().Get("tenant")
			if len(tenant) > 0 {
				span.AddAttributes(trace.StringAttribute("tenant", tenant))
			}

			body, err := ioutil.ReadAll(r.Body)
			if err == nil && len(body) > 0 {
				logger.Debugw("using body from request",
					"body", string(body),
					"traceID", span.SpanContext().TraceID.String(),
					"spanID", span.SpanContext().SpanID.String())

				// set the body again, because close() has been called on the body, which would cause an error later on
				r.Body = ioutil.NopCloser(bytes.NewBuffer(body))
			}

			r = r.WithContext(ctx)

			next.ServeHTTP(&sw, r)

			switch sw.status {
			case 200:
				span.SetStatus(trace.Status{Code: int32(trace.StatusCodeOK), Message: "200 OK"})
				break
			case 204:
				span.SetStatus(trace.Status{Code: int32(trace.StatusCodeOK), Message: "204 No content"})
				break
			case 400:
				span.SetStatus(trace.Status{Code: int32(trace.StatusCodeInvalidArgument), Message: "400 Invalid Request"})
				break
			case 401:
				span.SetStatus(trace.Status{Code: int32(trace.StatusCodeUnauthenticated), Message: "401 Not authorized"})
				break
			case 404:
				span.SetStatus(trace.Status{Code: int32(trace.StatusCodeNotFound), Message: "404 Not found"})
				break
			case 429:
				span.SetStatus(trace.Status{Code: int32(trace.StatusCodeResourceExhausted), Message: "429 Too many requests"})
				break
			case 500:
				span.SetStatus(trace.Status{Code: int32(trace.StatusCodeInternal), Message: "500 Internal Server Error"})
				break
			default:
				span.SetStatus(trace.Status{Code: int32(trace.StatusCodeUnknown), Message: "No known status code"})
			}
		})
	}
}

// AddStatus takes an error and an error code and sets them as status in the provided span.
func AddStatus(span *trace.Span, err error, code int) {

	if span == nil || err == nil {
		return
	}

	span.SetStatus(trace.Status{
		Code:    int32(code),
		Message: err.Error(),
	})
}

// LoggerWithTraceID reads the traceId from the context and adds it to the logging context.
func LoggerWithTraceID(ctx context.Context, logger *zap.SugaredLogger) *zap.SugaredLogger {

	traceId := ctx.Value(TraceID{})
	if traceId != nil {
		logger = logger.With("traceId", traceId)
	}

	return logger
}

// TraceId returns the traceId set in the context. If no traceId was found, an empty string is returned.
func TraceId(ctx context.Context) string {

	traceId := ctx.Value(TraceID{})
	if traceId != nil {
		return traceId.(string)
	}

	return ""
}

type statusWriter struct {
	http.ResponseWriter
	status int
}

// WriteHeader writes the status code to the HTTP response header.
func (w *statusWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

// Write writes the provided data to the HTTP response.
// Additionally, it sets 200 OK as status in the header if no status was set before.
func (w *statusWriter) Write(b []byte) (int, error) {
	n, err := w.ResponseWriter.Write(b)
	return n, err
}
