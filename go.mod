module github.com/nrempel/repro-gcp-trace

go 1.16

require (
	github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace v0.16.0
	github.com/golang/protobuf v1.4.3 // indirect
	github.com/gorilla/mux v1.8.0
	github.com/rs/zerolog v1.20.0
	github.com/yfuruyama/crzerolog v0.3.0
	go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux v0.17.0
	go.opentelemetry.io/otel v0.17.0
	go.opentelemetry.io/otel/exporters/stdout v0.17.0
	go.opentelemetry.io/otel/sdk v0.17.0
	go.opentelemetry.io/otel/trace v0.17.0
	golang.org/x/net v0.0.0-20210119194325-5f4716e94777 // indirect
	golang.org/x/sys v0.0.0-20210223095934-7937bea0104d // indirect
	golang.org/x/text v0.3.5 // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	google.golang.org/genproto v0.0.0-20210222152913-aa3ee6e6a81c // indirect
	google.golang.org/grpc v1.35.0 // indirect
)
