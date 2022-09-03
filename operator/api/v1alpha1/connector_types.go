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
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ConnectorSpec defines the desired state of Connector
type ConnectorSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// The image uri of your connector
	Image string `json:"image,omitempty"`
	// Specify the container in detail if needed
	Containers []v1.Container `json:"containers,omitempty"`

	ExposePort *int32 `json:"exposePort"`

	ConfigRef string `json:"configRef,omitempty"`

	SecretRef string `json:"secretRef,omitempty"`

	ScalingRule *ScalingRule `json:"scalingRule,omitempty"`
}

type ScalingRule struct {
	// +optional
	MaxReplicaCount *int32 `json:"maxReplicaCount,omitempty"`
	// +optional
	MinReplicaCount *int32 `json:"minReplicaCount,omitempty"`
	// +optional
	CustomScaling *CustomScaling `json:"customScaling,omitempty"`
	// +optional
	HTTPScaling *HTTPScaling `json:"httpScaling,omitempty"`
}

type HTTPScaling struct {
	Host string `json:"host"`
	// +optional
	SvcType string `json:"svcType,omitempty"`
	// +optional
	PendingRequests int32 `json:"pendingRequests,omitempty" description:"The target metric value for the HPA (Default 100)"`
}

type CustomScaling struct {
	// +optional
	CheckInterval *int32 `json:"checkInterval,omitempty"`
	// +optional
	CooldownPeriod *int32 `json:"cooldownPeriod,omitempty"`

	Triggers []Trigger `json:"triggers"`
}

type Trigger struct {
	Type      string            `json:"type"`
	Metadata  map[string]string `json:"metadata"`
	SecretRef string            `json:"secretRef,omitempty"`
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

type TriggerType string

const (
	MySQL TriggerType = "mysql"
	SQS   TriggerType = "aws-sqs-queue"
)

type ResourceType string

const (
	DeployResType ResourceType = "deployment"
	SvcResType    ResourceType = "service"
	CLSvcResType  ResourceType = "cl-service"
	LBSvcResType  ResourceType = "lb-service"
	HttpSOResType ResourceType = "httpScaledObject"
	SoResType     ResourceType = "scaledObject"
	TAResType     ResourceType = "triggerAuthentication"
)

func (in *Connector) String() string {

	return fmt.Sprintf("Image [%s]",
		in.Spec.Image)
}

func init() {
	SchemeBuilder.Register(&Connector{}, &ConnectorList{})
}
