package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	texporter "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
	"go.opentelemetry.io/otel"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	"github.com/yfuruyama/crzerolog"
)

func main() {
	initTracer()

	router := mux.NewRouter()

	rootLogger := zerolog.New(os.Stdout)
	loggingHandler := crzerolog.InjectLogger(&rootLogger)
	handler := loggingHandler(router)

	router.Use(otelmux.Middleware("repro-gcp-trace"))
	router.Use(traceRequest)

	// general routes
	router.HandleFunc("/hello", HelloController).Methods(http.MethodGet).Name("hello")

	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatal().Err(err).Msg("Startup failed")
	}

}

// Start a span named after the route name and request URI
func traceRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := log.Ctx(r.Context())

		log.Info().Msgf("Current route: %s", mux.CurrentRoute(r).GetName())
		ctx, requestSpan := tracer.Start(
			r.Context(),
			fmt.Sprintf("%s / %s", mux.CurrentRoute(r).GetName(), r.RequestURI),
		)
		defer requestSpan.End()

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func initTracer() {
	exporter, err := texporter.NewExporter()

	// exporter, err := stdout.NewExporter(
	// 	stdout.WithPrettyPrint(),
	// )
	if err != nil {
		log.Fatal().Err(err).Msgf("texporter.NewExporter failed.")
	}

	cfg := sdktrace.Config{
		DefaultSampler: sdktrace.AlwaysSample(),
	}
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithConfig(cfg),
		// sdktrace.WithSyncer(exporter),
		sdktrace.WithSyncer(exporter),
	)

	otel.SetTextMapPropagator(
		// propagation.TraceContext{},
		GCPPropagator{},
		// propagation.NewCompositeTextMapPropagator(
		// 	GCPPropagator{},
		// 	propagation.TraceContext{},
		// 	propagation.Baggage{},
		// ),
	)
	otel.SetTracerProvider(tp)
}
