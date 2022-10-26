module github.com/linkall-labs/sink-elasticsearch

go 1.17

require (
	github.com/cloudevents/sdk-go/v2 v2.12.0
	github.com/elastic/go-elasticsearch/v7 v7.17.1
	github.com/linkall-labs/cdk-go v0.0.0
	github.com/pkg/errors v0.9.1
	github.com/tidwall/gjson v1.14.0
)

replace github.com/linkall-labs/cdk-go v0.0.0 => ../../../cdk-go

require (
	github.com/google/go-cmp v0.5.7 // indirect
	github.com/google/uuid v1.1.2 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/sirupsen/logrus v1.9.0 // indirect
	github.com/tidwall/match v1.1.1 // indirect
	github.com/tidwall/pretty v1.2.0 // indirect
	go.uber.org/atomic v1.7.0 // indirect
	go.uber.org/goleak v1.1.12 // indirect
	go.uber.org/multierr v1.6.0 // indirect
	go.uber.org/zap v1.19.1 // indirect
	golang.org/x/sys v0.0.0-20220715151400-c0bba94af5f8 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)
