/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	"fmt"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ConnectorSpec defines the desired state of Connector
type ConnectorSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	ConnectorType Type `json:"type"`
	// The image uri of your connector
	Image string `json:"image,omitempty"`
	// Specify the container in detail if needed
	Containers []v1.Container `json:"containers,omitempty"`

	Scaler *ScalerSpec `json:"scalerSpec,omitempty"`
}

type ScalerSpec struct {
	// +optional
	CheckInterval *int32 `json:"checkInterval,omitempty"`
	// +optional
	CooldownPeriod *int32 `json:"cooldownPeriod,omitempty"`
	// +optional
	MaxReplicaCount *int32 `json:"maxReplicaCount,omitempty"`
	// +optional
	MinReplicaCount *int32 `json:"minReplicaCount,omitempty"`

	Metadata     map[string]intstr.IntOrString `json:"metadata,omitempty"`
	ScalerSecret string                        `json:"scalerSecret,omitempty"`
}

// ConnectorStatus defines the observed state of Connector
type ConnectorStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Conditions []ConnectorCondition `json:"conditions,omitempty"`
}
type ConnectorCondition struct {
	Timestamp string                   `json:"timestamp" description:"Timestamp of this condition"`
	Type      ConnectorCreationStatus  `json:"type" description:"type of status condition"`
	Status    metav1.ConditionStatus   `json:"status" description:"status of the condition, one of True, False, Unknown"`
	Reason    ConnectorConditionReason `json:"reason,omitempty" description:"one-word CamelCase reason for the condition's last transition"`
	Message   string                   `json:"message,omitempty" description:"human-readable message indicating details about last transition"`
}
type ConnectorCreationStatus string
type ConnectorConditionReason string

const (
	ErrorCreatingAppScaledObject    ConnectorConditionReason = "ErrorCreatingAppScaledObject"
	AppScaledObjectCreated          ConnectorConditionReason = "AppScaledObjectCreated"
	TerminatingResources            ConnectorConditionReason = "TerminatingResources"
	AppScaledObjectTerminated       ConnectorConditionReason = "AppScaledObjectTerminated"
	AppScaledObjectTerminationError ConnectorConditionReason = "AppScaledObjectTerminationError"
	PendingCreation                 ConnectorConditionReason = "PendingCreation"
	HTTPScaledObjectIsReady         ConnectorConditionReason = "HTTPScaledObjectIsReady"
)
const (
	// Created indicates the resource has been created
	Created ConnectorCreationStatus = "Created"
	// Terminated indicates the resource has been terminated
	Terminated ConnectorCreationStatus = "Terminated"
	// Error indicates the resource had an error
	Error ConnectorCreationStatus = "Error"
	// Pending indicates the resource hasn't been created
	Pending ConnectorCreationStatus = "Pending"
	// Terminating indicates that the resource is marked for deletion but hasn't
	// been deleted yet
	Terminating ConnectorCreationStatus = "Terminating"
	// Unknown indicates the status is unavailable
	Unknown ConnectorCreationStatus = "Unknown"
	// Ready indicates the object is fully created
	Ready ConnectorCreationStatus = "Ready"
)

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Connector is the Schema for the connectors API
type Connector struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ConnectorSpec   `json:"spec"`
	Status ConnectorStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ConnectorList contains a list of Connector
type ConnectorList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Connector `json:"items"`
}

type ResourceType string

const (
	DeployResType ResourceType = "deployment"
	SvcResType    ResourceType = "service"
	HttpSOResType ResourceType = "httpScaledObject"
	SoResType     ResourceType = "scaledObject"
	TAResType     ResourceType = "triggerAuthentication"
)

// Type describes the pattern the image uses to obtain data and reflects to a corresponding scaler.
// +enum
type Type string

const (
	// Http type means that the image is a webserver waiting data to be pushed to it.
	// This type also reflects to a http scaler.
	Http Type = "http"
	// ActiveMQ type means that the image fetches data from an ActiveMQ queue
	// This type also reflects to a ActiveMQ scaler.
	ActiveMQ Type = "activemq"
	// ArtemisQue type means that the image fetches data from an ActiveMQ Artemis queue
	// This type also reflects to an artemis-queue scaler.
	ArtemisQue Type = "artemis-queue"
	// Kafka type means that the image fetches data from an Apache Kafka topic
	// This type also reflects to a Kafka scaler.
	Kafka Type = "kafka"
	// AWSCloudwatch type means that the image fetches data from AWS Cloudwatch
	// This type also reflects to a AWSCloudwatch scaler.
	AWSCloudwatch Type = "aws-cloudwatch"
	// AWSKinesisStream type means that the image fetches data from AWS Kinesis Stream
	// This type also reflects to a AWSKinesisStream scaler.
	AWSKinesisStream Type = "aws-kinesis-stream"
	// AWSSqsQueue type means that the image fetches data from AWS Sqs Queue
	// This type also reflects to a AWSSqsQueue scaler.
	AWSSqsQueue Type = "aws-sqs-queue"
	// AZUREAppInsights type means that the image fetches data from Azure Application Insights
	// This type also reflects to a AZUREAppInsights scaler.
	AZUREAppInsights Type = "azure-app-insights"
	// Rabbitmq type means that the image fetches data from Rabbitmq
	// This type also reflects to a Rabbitmq scaler.
	Rabbitmq Type = "rabbitmq"
)

func (in *Connector) String() string {

	return fmt.Sprintf("Image [%s], Mode [%s]",
		in.Spec.Image,
		in.Spec.ConnectorType)
}

func init() {
	SchemeBuilder.Register(&Connector{}, &ConnectorList{})
}
