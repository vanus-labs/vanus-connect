module github.com/linkall-labs/connector/sink/tencent-cloud/function

go 1.18

require (
	github.com/cloudevents/sdk-go/v2 v2.8.0
	github.com/google/uuid v1.1.1
	github.com/linkall-labs/cdk-go v0.0.0
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common v1.0.527
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/scf v1.0.527
)

require (
	github.com/google/go-cmp v0.5.5 // indirect
	github.com/json-iterator/go v1.1.10 // indirect
	github.com/modern-go/concurrent v0.0.0-20180228061459-e0a39a4cb421 // indirect
	github.com/modern-go/reflect2 v0.0.0-20180701023420-4b7aa43c6742 // indirect
	github.com/sirupsen/logrus v1.9.0 // indirect
	go.uber.org/atomic v1.4.0 // indirect
	go.uber.org/multierr v1.1.0 // indirect
	go.uber.org/zap v1.10.0 // indirect
	golang.org/x/sys v0.0.0-20220715151400-c0bba94af5f8 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/linkall-labs/cdk-go => ../../../cdk-go
