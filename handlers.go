package main

import (
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel"
)

var tracer = otel.Tracer("repro-gcp-trace")

func HelloController(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := log.Ctx(ctx)

	// spanContext := trace.SpanFromContext(ctx).SpanContext()

	// println("spanContext.TraceID.String()")
	// println(spanContext.TraceID.String())

	log.Info().Msg("Reached HelloController")

	ctx, span := tracer.Start(ctx, "HelloController::wait100")
	time.Sleep(100 * time.Millisecond)
	span.End()

	_, span = tracer.Start(ctx, "HelloController::wait750")
	time.Sleep(750 * time.Millisecond)
	span.End()

	w.WriteHeader(http.StatusOK)
	log.Info().Msg("Leaving HelloController")
}
