package main

import (
	"context"
	"encoding/binary"
	"fmt"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

const (
	httpHeader = "X-Cloud-Trace-Context"
)

// GCPPropagator propagator serializes SpanContext to/from X-Cloud-Trace-Context HTTP Headers.
type GCPPropagator struct{}

var _ propagation.TextMapPropagator = &GCPPropagator{}

// Inject injects a context into the carrier as HTTP headers.
func (gcpPropagator GCPPropagator) Inject(context context.Context, carrier propagation.TextMapCarrier) {
	log := log.Ctx(context)

	spanContext := trace.SpanFromContext(context).SpanContext()
	if !spanContext.TraceID.IsValid() || !spanContext.SpanID.IsValid() {
		return
	}

	spanID := binary.BigEndian.Uint64(spanContext.SpanID[:])
	header := fmt.Sprintf("%s/%d;o=%d", spanContext.TraceID.String(), spanID, spanContext.TraceFlags)
	log.Trace().Msgf("Inject: header set to %s", header)
	carrier.Set(httpHeader, header)
}

// Extract extracts a context from the carrier if it contains HTTP headers.
func (gcpPropagator GCPPropagator) Extract(context context.Context, carrier propagation.TextMapCarrier) context.Context {
	log := log.Ctx(context)

	if header := carrier.Get(httpHeader); header != "" {
		spanContext, err := extract(header)
		if err != nil {
			log.Error().Err(err).Msgf("Failed to extract %s", httpHeader)
		} else if err == nil && spanContext.IsValid() {
			return trace.ContextWithRemoteSpanContext(context, spanContext)
		}
	}

	return context
}

func extract(traceHeader string) (trace.SpanContext, error) {
	// Parsing X-Cloud-Trace-Context header from Google Cloud load balancer
	// More info: https://cloud.google.com/trace/docs/setup#force-trace
	// ;o=TRACE_TRUE may or may not exist
	spanContext := trace.SpanContext{}
	malformedHeaderError :=
		fmt.Errorf(
			"Malformed %s header. Header must be of the form "+
				"TRACE_ID/SPAN_ID;o=TRACE_TRUE "+
				"but %s was received.",
			httpHeader,
			traceHeader,
		)

	headerTokens := strings.Split(traceHeader, "/")
	if len(headerTokens) != 2 {
		return spanContext, malformedHeaderError
	}

	// Parse out TraceID
	traceIdHex := headerTokens[0]
	traceID, err := trace.TraceIDFromHex(traceIdHex)
	if err != nil {
		return spanContext, err
	}
	spanContext.TraceID = traceID

	spanIdAndOption := strings.Split(headerTokens[1], ";o=")
	if len(spanIdAndOption) < 1 {
		return spanContext, malformedHeaderError
	}

	// Parse out SpanID
	spanIdAsString := spanIdAndOption[0]
	spanId, err := strconv.ParseUint(spanIdAsString, 10, 64)
	if err != nil {
		return spanContext, err
	}
	binary.BigEndian.PutUint64(spanContext.SpanID[:], spanId)

	// Parse out option if it exists
	if len(spanIdAndOption) == 2 {
		switch spanIdAndOption[1] {
		case "1":
			spanContext.TraceFlags = trace.FlagsSampled
		}
	}

	return spanContext, nil
}

// Fields is propagation keys
func (gcpPropagator GCPPropagator) Fields() []string {
	return []string{httpHeader}
}
