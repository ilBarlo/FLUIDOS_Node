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
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	nodecorev1alpha1 "github.com/fluidos-project/node/apis/nodecore/v1alpha1"
)

type Partition struct {
	Architecture     string            `json:"architecture"`
	Cpu              resource.Quantity `json:"cpu"`
	Memory           resource.Quantity `json:"memory"`
	Gpu              resource.Quantity `json:"gpu,omitempty"`
	EphemeralStorage resource.Quantity `json:"ephemeral-storage,omitempty"`
	Storage          resource.Quantity `json:"storage,omitempty"`
}

// ReservationSpec defines the desired state of Reservation
type ReservationSpec struct {

	// SolverID is the ID of the solver that asks for the reservation
	SolverID string `json:"solverID"`

	// This is the Node identity of the buyer FLUIDOS Node.
	Buyer nodecorev1alpha1.NodeIdentity `json:"buyer"`

	// BuyerClusterID is the Liqo ClusterID used by the seller to search a contract and the related resources during the peering phase.
	BuyerClusterID string `json:"buyerClusterID"`

	// This is the Node identity of the seller FLUIDOS Node.
	Seller nodecorev1alpha1.NodeIdentity `json:"seller"`

	// Parition is the partition of the flavour that is being reserved
	Partition *Partition `json:"partition,omitempty"`

	// Reserve indicates if the reservation is a reserve or not
	Reserve bool `json:"reserve,omitempty"`

	// Purchase indicates if the reservation is an purchase or not
	Purchase bool `json:"purchase,omitempty"`

	// PeeringCandidate is the reference to the PeeringCandidate of the Reservation
	PeeringCandidate nodecorev1alpha1.GenericRef `json:"peeringCandidate,omitempty"`
}

// ReservationStatus defines the observed state of Reservation
type ReservationStatus struct {
	// This is the current phase of the reservation
	Phase nodecorev1alpha1.PhaseStatus `json:"phase"`

	// ReservePhase is the current phase of the reservation
	ReservePhase nodecorev1alpha1.Phase `json:"reservePhase,omitempty"`

	// PurchasePhase is the current phase of the reservation
	PurchasePhase nodecorev1alpha1.Phase `json:"purchasePhase,omitempty"`

	// TransactionID is the ID of the transaction that this reservation is part of
	TransactionID string `json:"transactionID,omitempty"`

	// Contract is the reference to the Contract of the Reservation
	Contract nodecorev1alpha1.GenericRef `json:"contract,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// +kubebuilder:printcolumn:name="Solver ID",type=string,JSONPath=`.spec.solverID`
// +kubebuilder:printcolumn:name="Reserve",type=boolean,JSONPath=`.spec.reserve`
// +kubebuilder:printcolumn:name="Purchase",type=boolean,JSONPath=`.spec.purchase`
// +kubebuilder:printcolumn:name="Seller",type=string,JSONPath=`.spec.seller.name`
// +kubebuilder:printcolumn:name="Peering Candidate",type=string,priority=1,JSONPath=`.spec.peeringCandidate.name`
// +kubebuilder:printcolumn:name="Transaction ID",type=string,JSONPath=`.status.transactionID`
// +kubebuilder:printcolumn:name="Reserve Phase",type=string,priority=1,JSONPath=`.status.reservePhase.phase`
// +kubebuilder:printcolumn:name="Purchase Phase",type=string,priority=1,JSONPath=`.status.purchasePhase.phase`
// +kubebuilder:printcolumn:name="Contract Name",type=string,JSONPath=`.status.contract.name`
// +kubebuilder:printcolumn:name="Status",type=string,JSONPath=`.status.phase.phase`
// +kubebuilder:printcolumn:name="Message",type=string,priority=1,JSONPath=`.status.phase.message`
// Reservation is the Schema for the reservations API
type Reservation struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ReservationSpec   `json:"spec,omitempty"`
	Status ReservationStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ReservationList contains a list of Reservation
type ReservationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Reservation `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Reservation{}, &ReservationList{})
}
