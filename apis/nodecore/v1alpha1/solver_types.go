// Copyright 2022-2023 FLUIDOS Project
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package v1alpha1

import (
	resource "k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Phase string

// PhaseStatus represents the status of a phase of the solver. I.e. the status of the REAR phases.
type PhaseStatus struct {
	Phase          Phase  `json:"phase"`
	Message        string `json:"message,omitempty"`
	StartTime      string `json:"startTime,omitempty"`
	LastChangeTime string `json:"lastChangeTime,omitempty"`
	EndTime        string `json:"endTime,omitempty"`
}

// Selector represents the criteria for selecting Flavours.
/* type Selector struct {
	FlavourType      string `json:"type,omitempty"`
	Architecture     string `json:"architecture,omitempty"`
	Cpu              int    `json:"cpu,omitempty"`
	Memory           int    `json:"memory,omitempty"`
	EphemeralStorage int    `json:"ephemeral-storage,omitempty"`
	MoreThanCpu      int    `json:"moreThanCpu,omitempty"`
	MoreThanMemory   int    `json:"moreThanMemory,omitempty"`
	MoreThanEph      int    `json:"moreThanEph,omitempty"`
	LessThanCpu      int    `json:"lessThanCpu,omitempty"`
	LessThanMemory   int    `json:"lessThanMemory,omitempty"`
	LessThanEph      int    `json:"lessThanEph,omitempty"`
} */

type FlavourSelector struct {
	FlavourType   string         `json:"type"`
	Architecture  string         `json:"architecture"`
	RangeSelector *RangeSelector `json:"rangeSelector,omitempty"`
	MatchSelector *MatchSelector `json:"matchSelector,omitempty"`
}

// MatchSelector represents the criteria for selecting Flavours through a strict match.
type MatchSelector struct {
	Cpu              resource.Quantity `json:"cpu"`
	Memory           resource.Quantity `json:"memory"`
	Storage          resource.Quantity `json:"storage,omitempty"`
	EphemeralStorage resource.Quantity `json:"ephemeralStorage,omitempty"`
	Gpu              resource.Quantity `json:"gpu,omitempty"`
}

// RangeSelector represents the criteria for selecting Flavours through a range.
type RangeSelector struct {
	MinCpu     resource.Quantity `json:"minCpu,omitempty"`
	MinMemory  resource.Quantity `json:"minMemory,omitempty"`
	MinEph     resource.Quantity `json:"minEph,omitempty"`
	MinStorage resource.Quantity `json:"minStorage,omitempty"`
	MinGpu     resource.Quantity `json:"minGpu,omitempty"`
	MaxCpu     resource.Quantity `json:"MaxCpu,omitempty"`
	MaxMemory  resource.Quantity `json:"MaxMemory,omitempty"`
	MaxEph     resource.Quantity `json:"MaxEph,omitempty"`
	MaxStorage resource.Quantity `json:"MaxStorage,omitempty"`
	MaxGpu     resource.Quantity `json:"MaxGpu,omitempty"`
}

// SolverSpec defines the desired state of Solver
type SolverSpec struct {

	// Selector contains the flavour requirements for the solver.
	Selector *FlavourSelector `json:"selector,omitempty"`

	// IntentID is the ID of the intent that the Node Orchestrator is trying to solve.
	// It is used to link the solver with the intent.
	IntentID string `json:"intentID"`

	// FindCandidate is a flag that indicates if the solver should find a candidate to solve the intent.
	FindCandidate bool `json:"findCandidate,omitempty"`

	// ReserveAndBuy is a flag that indicates if the solver should reserve and buy the resources on the candidate.
	ReserveAndBuy bool `json:"reserveAndBuy,omitempty"`

	// EnstablishPeering is a flag that indicates if the solver should enstablish a peering with the candidate.
	EnstablishPeering bool `json:"enstablishPeering,omitempty"`
}

// SolverStatus defines the observed state of Solver
type SolverStatus struct {

	// FindCandidate describes the status of research of the candidate.
	// Rear Manager is looking for the best candidate Flavour to solve the Node Orchestrator request.
	FindCandidate Phase `json:"findCandidate,omitempty"`

	// ReserveAndBuy describes the status of the reservation and purchase of selected Flavour.
	// Rear Manager is trying to reserve and purchase the resources on the candidate FLUIDOS Node.
	ReserveAndBuy Phase `json:"reserveAndBuy,omitempty"`

	// Peering describes the status of the peering with the candidate.
	// Rear Manager is trying to enstablish a peering with the candidate FLUIDOS Node.
	Peering Phase `json:"peering,omitempty"`

	// DiscoveryPhase describes the status of the Discovery where the Discovery Manager
	// is looking for matching flavours outside the FLUIDOS Node
	DiscoveryPhase Phase `json:"discoveryPhase,omitempty"`

	// ReservationPhase describes the status of the Reservation where the Contract Manager
	// is reserving and purchasing the resources on the candidate node.
	ReservationPhase Phase `json:"reservationPhase,omitempty"`

	// ConsumePhase describes the status of the Consume phase where the VFM (Liqo) is enstablishing
	// a peering with the candidate node.
	ConsumePhase Phase `json:"consumePhase,omitempty"`

	// SolverPhase describes the status of the Solver generated by the Node Orchestrator.
	// It is usefull to understand if the solver is still running or if it has finished or failed.
	SolverPhase PhaseStatus `json:"solverPhase,omitempty"`

	// PeeringCandidate contains the candidate that the solver has eventually found to solve the intent.
	PeeringCandidate GenericRef `json:"peeringCandidate,omitempty"`

	// Allocation contains the allocation that the solver has eventually created for the intent.
	// It can correspond to a virtual node
	// The Node Orchestrator will use this allocation to fullfill the intent.
	Allocation GenericRef `json:"allocation,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Solver is the Schema for the solvers API
// +kubebuilder:printcolumn:name="Intent ID",type=string,JSONPath=`.spec.intentID`
// +kubebuilder:printcolumn:name="Find Candidate",type=boolean,JSONPath=`.spec.findCandidate`
// +kubebuilder:printcolumn:name="Reserve and Buy",type=boolean,JSONPath=`.spec.reserveAndBuy`
// +kubebuilder:printcolumn:name="Peering",type=boolean,JSONPath=`.spec.enstablishPeering`
// +kubebuilder:printcolumn:name="Candidate Phase",type=string,priority=1,JSONPath=`.status.findCandidate`
// +kubebuilder:printcolumn:name="Reserving Phase",type=string,priority=1,JSONPath=`.status.reserveAndBuy`
// +kubebuilder:printcolumn:name="Peering Phase",type=string,priority=1,JSONPath=`.status.peering`
// +kubebuilder:printcolumn:name="Status",type=string,JSONPath=`.status.solverPhase.phase`
// +kubebuilder:printcolumn:name="Message",type=string,JSONPath=`.status.solverPhase.message`
// +kubebuilder:printcolumn:name="Age",type=date,JSONPath=`.metadata.creationTimestamp`
// Solver is the Schema for the solvers API
type Solver struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SolverSpec   `json:"spec,omitempty"`
	Status SolverStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// SolverList contains a list of Solver
type SolverList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Solver `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Solver{}, &SolverList{})
}
