module github.com/linkall-labs/connector/mongodb-sink

go 1.18

require (
	github.com/cloudevents/sdk-go/v2 v2.11.0
	github.com/golang/protobuf v1.5.2
	github.com/linkall-labs/connector/proto v0.0.0
	google.golang.org/protobuf v1.28.1
)

replace github.com/linkall-labs/connector/proto => ../../schemas/pkg

require (
	github.com/google/uuid v1.1.1 // indirect
	github.com/json-iterator/go v1.1.10 // indirect
	github.com/modern-go/concurrent v0.0.0-20180228061459-e0a39a4cb421 // indirect
	github.com/modern-go/reflect2 v0.0.0-20180701023420-4b7aa43c6742 // indirect
	go.uber.org/atomic v1.4.0 // indirect
	go.uber.org/multierr v1.1.0 // indirect
	go.uber.org/zap v1.10.0 // indirect
)
