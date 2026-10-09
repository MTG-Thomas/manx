// SPDX-License-Identifier: AGPL-3.0-or-later
// Generated from schema.json SHA-256 76de2248d578cf6ccf77247e591a62f05164558a8e26d9d95ed5c200cb275d7b; do not edit.
package contract

import (
	"bytes"
	"encoding/json"
	"fmt"
	"regexp"
	"time"
)

type Registration struct {
	ContractVersion string `json:"contract_version"`
	InstallationId  string `json:"installation_id"`
	ClientKind      string `json:"client_kind"`
	PublicKey       string `json:"public_key"`
	SerialNumber    string `json:"serial_number,omitempty"`
	SmbiosUuid      string `json:"smbios_uuid,omitempty"`
	Manufacturer    string `json:"manufacturer,omitempty"`
	Model           string `json:"model,omitempty"`
	Architecture    string `json:"architecture"`
	BootEnvironment string `json:"boot_environment"`
	ClientVersion   string `json:"client_version"`
	ObservedAt      string `json:"observed_at"`
}

func (v *Registration) UnmarshalJSON(data []byte) error {
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(data, &fields); err != nil {
		return err
	}
	if _, present := fields["contract_version"]; !present {
		return fmt.Errorf("contract_version required")
	}
	if _, present := fields["installation_id"]; !present {
		return fmt.Errorf("installation_id required")
	}
	if _, present := fields["client_kind"]; !present {
		return fmt.Errorf("client_kind required")
	}
	if _, present := fields["public_key"]; !present {
		return fmt.Errorf("public_key required")
	}
	if _, present := fields["architecture"]; !present {
		return fmt.Errorf("architecture required")
	}
	if _, present := fields["boot_environment"]; !present {
		return fmt.Errorf("boot_environment required")
	}
	if _, present := fields["client_version"]; !present {
		return fmt.Errorf("client_version required")
	}
	if _, present := fields["observed_at"]; !present {
		return fmt.Errorf("observed_at required")
	}
	type alias Registration
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	return decoder.Decode((*alias)(v))
}

func (v Registration) Validate() error {
	if v.ContractVersion != "bifrost.device-lifecycle/v1" {
		return fmt.Errorf("contract_version value")
	}
	if !regexp.MustCompile("^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$").MatchString(v.InstallationId) {
		return fmt.Errorf("installation_id pattern")
	}
	if v.ClientKind != "windows_bootstrap" && v.ClientKind != "manx_recovery" {
		return fmt.Errorf("client_kind value")
	}
	if len(v.PublicKey) < 100 {
		return fmt.Errorf("public_key length")
	}
	if len(v.PublicKey) > 160 {
		return fmt.Errorf("public_key length")
	}
	if !regexp.MustCompile("^[A-Za-z0-9+/]+={0,2}$").MatchString(v.PublicKey) {
		return fmt.Errorf("public_key pattern")
	}
	if v.SerialNumber != "" && len(v.SerialNumber) > 128 {
		return fmt.Errorf("serial_number length")
	}
	if v.SmbiosUuid != "" && len(v.SmbiosUuid) > 128 {
		return fmt.Errorf("smbios_uuid length")
	}
	if v.Manufacturer != "" && len(v.Manufacturer) > 128 {
		return fmt.Errorf("manufacturer length")
	}
	if v.Model != "" && len(v.Model) > 128 {
		return fmt.Errorf("model length")
	}
	if v.Architecture != "x64" && v.Architecture != "arm64" && v.Architecture != "x86" && v.Architecture != "unknown" {
		return fmt.Errorf("architecture value")
	}
	if v.BootEnvironment != "windows11" && v.BootEnvironment != "linux_rescue" {
		return fmt.Errorf("boot_environment value")
	}
	if len(v.ClientVersion) < 1 {
		return fmt.Errorf("client_version length")
	}
	if len(v.ClientVersion) > 40 {
		return fmt.Errorf("client_version length")
	}
	if _, err := time.Parse(time.RFC3339, v.ObservedAt); err != nil {
		return fmt.Errorf("observed_at timestamp")
	}
	return nil
}

type EventPayload struct {
	BootId     string `json:"boot_id"`
	Checkpoint string `json:"checkpoint"`
	Result     string `json:"result"`
	Diagnostic string `json:"diagnostic,omitempty"`
}

func (v *EventPayload) UnmarshalJSON(data []byte) error {
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(data, &fields); err != nil {
		return err
	}
	if _, present := fields["boot_id"]; !present {
		return fmt.Errorf("boot_id required")
	}
	if _, present := fields["checkpoint"]; !present {
		return fmt.Errorf("checkpoint required")
	}
	if _, present := fields["result"]; !present {
		return fmt.Errorf("result required")
	}
	type alias EventPayload
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	return decoder.Decode((*alias)(v))
}

func (v EventPayload) Validate() error {
	if !regexp.MustCompile("^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$").MatchString(v.BootId) {
		return fmt.Errorf("boot_id pattern")
	}
	if v.Checkpoint != "started" && v.Checkpoint != "registered" && v.Checkpoint != "approval_observed" && v.Checkpoint != "diagnostic_collected" && v.Checkpoint != "reconnected" && v.Checkpoint != "handoff_verified" {
		return fmt.Errorf("checkpoint value")
	}
	if v.Result != "ok" && v.Result != "unavailable" && v.Result != "intervention_required" {
		return fmt.Errorf("result value")
	}
	if v.Diagnostic != "" && (v.Diagnostic != "hardware_inventory" && v.Diagnostic != "network_connectivity" && v.Diagnostic != "session_recovery") {
		return fmt.Errorf("diagnostic value")
	}
	return nil
}

type ClientEvent struct {
	ContractVersion string       `json:"contract_version"`
	SessionId       string       `json:"session_id"`
	InstallationId  string       `json:"installation_id"`
	Sequence        int64        `json:"sequence"`
	EventId         string       `json:"event_id"`
	EventType       string       `json:"event_type"`
	Timestamp       string       `json:"timestamp"`
	Payload         EventPayload `json:"payload"`
}

func (v *ClientEvent) UnmarshalJSON(data []byte) error {
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(data, &fields); err != nil {
		return err
	}
	if _, present := fields["contract_version"]; !present {
		return fmt.Errorf("contract_version required")
	}
	if _, present := fields["session_id"]; !present {
		return fmt.Errorf("session_id required")
	}
	if _, present := fields["installation_id"]; !present {
		return fmt.Errorf("installation_id required")
	}
	if _, present := fields["sequence"]; !present {
		return fmt.Errorf("sequence required")
	}
	if _, present := fields["event_id"]; !present {
		return fmt.Errorf("event_id required")
	}
	if _, present := fields["event_type"]; !present {
		return fmt.Errorf("event_type required")
	}
	if _, present := fields["timestamp"]; !present {
		return fmt.Errorf("timestamp required")
	}
	if _, present := fields["payload"]; !present {
		return fmt.Errorf("payload required")
	}
	type alias ClientEvent
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	return decoder.Decode((*alias)(v))
}

func (v ClientEvent) Validate() error {
	if v.ContractVersion != "bifrost.device-lifecycle/v1" {
		return fmt.Errorf("contract_version value")
	}
	if !regexp.MustCompile("^[a-f0-9]{64}$").MatchString(v.SessionId) {
		return fmt.Errorf("session_id pattern")
	}
	if !regexp.MustCompile("^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$").MatchString(v.InstallationId) {
		return fmt.Errorf("installation_id pattern")
	}
	if v.Sequence < 1 {
		return fmt.Errorf("sequence range")
	}
	if v.Sequence > 4096 {
		return fmt.Errorf("sequence range")
	}
	if !regexp.MustCompile("^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$").MatchString(v.EventId) {
		return fmt.Errorf("event_id pattern")
	}
	if v.EventType != "checkpoint" && v.EventType != "diagnostic" && v.EventType != "intervention_requested" && v.EventType != "handoff_ready" && v.EventType != "handoff_verified" {
		return fmt.Errorf("event_type value")
	}
	if _, err := time.Parse(time.RFC3339, v.Timestamp); err != nil {
		return fmt.Errorf("timestamp timestamp")
	}
	if err := v.Payload.Validate(); err != nil {
		return err
	}
	return nil
}

type Heartbeat struct {
	ContractVersion string `json:"contract_version"`
	InstallationId  string `json:"installation_id"`
	BootId          string `json:"boot_id"`
}

func (v *Heartbeat) UnmarshalJSON(data []byte) error {
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(data, &fields); err != nil {
		return err
	}
	if _, present := fields["contract_version"]; !present {
		return fmt.Errorf("contract_version required")
	}
	if _, present := fields["installation_id"]; !present {
		return fmt.Errorf("installation_id required")
	}
	if _, present := fields["boot_id"]; !present {
		return fmt.Errorf("boot_id required")
	}
	type alias Heartbeat
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	return decoder.Decode((*alias)(v))
}

func (v Heartbeat) Validate() error {
	if v.ContractVersion != "bifrost.device-lifecycle/v1" {
		return fmt.Errorf("contract_version value")
	}
	if !regexp.MustCompile("^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$").MatchString(v.InstallationId) {
		return fmt.Errorf("installation_id pattern")
	}
	if !regexp.MustCompile("^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$").MatchString(v.BootId) {
		return fmt.Errorf("boot_id pattern")
	}
	return nil
}

type Assignment struct {
	OrganizationId   string `json:"organization_id"`
	OrganizationName string `json:"organization_name"`
	HaloTicketId     int64  `json:"halo_ticket_id"`
	DeviceId         string `json:"device_id"`
}

func (v *Assignment) UnmarshalJSON(data []byte) error {
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(data, &fields); err != nil {
		return err
	}
	if _, present := fields["organization_id"]; !present {
		return fmt.Errorf("organization_id required")
	}
	if _, present := fields["organization_name"]; !present {
		return fmt.Errorf("organization_name required")
	}
	if _, present := fields["halo_ticket_id"]; !present {
		return fmt.Errorf("halo_ticket_id required")
	}
	if _, present := fields["device_id"]; !present {
		return fmt.Errorf("device_id required")
	}
	type alias Assignment
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	return decoder.Decode((*alias)(v))
}

func (v Assignment) Validate() error {
	if !regexp.MustCompile("^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$").MatchString(v.OrganizationId) {
		return fmt.Errorf("organization_id pattern")
	}
	if len(v.OrganizationName) < 1 {
		return fmt.Errorf("organization_name length")
	}
	if len(v.OrganizationName) > 128 {
		return fmt.Errorf("organization_name length")
	}
	if v.HaloTicketId < 1 {
		return fmt.Errorf("halo_ticket_id range")
	}
	if len(v.DeviceId) < 1 {
		return fmt.Errorf("device_id length")
	}
	if len(v.DeviceId) > 200 {
		return fmt.Errorf("device_id length")
	}
	return nil
}

type OperatorMutation struct {
	ContractVersion  string     `json:"contract_version"`
	RequestId        string     `json:"request_id"`
	ActorId          string     `json:"actor_id"`
	SessionId        string     `json:"session_id"`
	ExpectedRevision int64      `json:"expected_revision"`
	IdentitySha256   string     `json:"identity_sha256"`
	PairingCode      string     `json:"pairing_code"`
	Assignment       Assignment `json:"assignment,omitempty"`
	Reason           string     `json:"reason,omitempty"`
}

func (v *OperatorMutation) UnmarshalJSON(data []byte) error {
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(data, &fields); err != nil {
		return err
	}
	if _, present := fields["contract_version"]; !present {
		return fmt.Errorf("contract_version required")
	}
	if _, present := fields["request_id"]; !present {
		return fmt.Errorf("request_id required")
	}
	if _, present := fields["actor_id"]; !present {
		return fmt.Errorf("actor_id required")
	}
	if _, present := fields["session_id"]; !present {
		return fmt.Errorf("session_id required")
	}
	if _, present := fields["expected_revision"]; !present {
		return fmt.Errorf("expected_revision required")
	}
	if _, present := fields["identity_sha256"]; !present {
		return fmt.Errorf("identity_sha256 required")
	}
	if _, present := fields["pairing_code"]; !present {
		return fmt.Errorf("pairing_code required")
	}
	type alias OperatorMutation
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	return decoder.Decode((*alias)(v))
}

func (v OperatorMutation) Validate() error {
	if v.ContractVersion != "bifrost.device-lifecycle/v1" {
		return fmt.Errorf("contract_version value")
	}
	if !regexp.MustCompile("^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$").MatchString(v.RequestId) {
		return fmt.Errorf("request_id pattern")
	}
	if !regexp.MustCompile("^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$").MatchString(v.ActorId) {
		return fmt.Errorf("actor_id pattern")
	}
	if !regexp.MustCompile("^[a-f0-9]{64}$").MatchString(v.SessionId) {
		return fmt.Errorf("session_id pattern")
	}
	if v.ExpectedRevision < 1 {
		return fmt.Errorf("expected_revision range")
	}
	if !regexp.MustCompile("^[a-f0-9]{64}$").MatchString(v.IdentitySha256) {
		return fmt.Errorf("identity_sha256 pattern")
	}
	if !regexp.MustCompile("^[A-Z2-9]{8}$").MatchString(v.PairingCode) {
		return fmt.Errorf("pairing_code pattern")
	}
	if err := v.Assignment.Validate(); err != nil {
		return err
	}
	if v.Reason != "" && (v.Reason != "technician_approved" && v.Reason != "technician_rejected" && v.Reason != "identity_collision" && v.Reason != "technician_cancelled") {
		return fmt.Errorf("reason value")
	}
	return nil
}

type Session struct {
	ContractVersion string          `json:"contract_version"`
	SessionId       string          `json:"session_id"`
	Registration    Registration    `json:"registration"`
	IdentitySha256  string          `json:"identity_sha256"`
	PairingCode     string          `json:"pairing_code"`
	State           string          `json:"state"`
	Revision        int64           `json:"revision"`
	LastSequence    int64           `json:"last_sequence"`
	CreatedAt       string          `json:"created_at"`
	ExpiresAt       string          `json:"expires_at"`
	LastSeenAt      string          `json:"last_seen_at"`
	Assignment      *Assignment     `json:"assignment"`
	ClaimedBy       *string         `json:"claimed_by"`
	ApprovedBy      *string         `json:"approved_by"`
	Events          []ClientEvent   `json:"events"`
	OperatorEvents  []OperatorEvent `json:"operator_events"`
	ResumeState     string          `json:"resume_state,omitempty"`
}

func (v *Session) UnmarshalJSON(data []byte) error {
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(data, &fields); err != nil {
		return err
	}
	if _, present := fields["contract_version"]; !present {
		return fmt.Errorf("contract_version required")
	}
	if _, present := fields["session_id"]; !present {
		return fmt.Errorf("session_id required")
	}
	if _, present := fields["registration"]; !present {
		return fmt.Errorf("registration required")
	}
	if _, present := fields["identity_sha256"]; !present {
		return fmt.Errorf("identity_sha256 required")
	}
	if _, present := fields["pairing_code"]; !present {
		return fmt.Errorf("pairing_code required")
	}
	if _, present := fields["state"]; !present {
		return fmt.Errorf("state required")
	}
	if _, present := fields["revision"]; !present {
		return fmt.Errorf("revision required")
	}
	if _, present := fields["last_sequence"]; !present {
		return fmt.Errorf("last_sequence required")
	}
	if _, present := fields["created_at"]; !present {
		return fmt.Errorf("created_at required")
	}
	if _, present := fields["expires_at"]; !present {
		return fmt.Errorf("expires_at required")
	}
	if _, present := fields["last_seen_at"]; !present {
		return fmt.Errorf("last_seen_at required")
	}
	if _, present := fields["assignment"]; !present {
		return fmt.Errorf("assignment required")
	}
	if _, present := fields["claimed_by"]; !present {
		return fmt.Errorf("claimed_by required")
	}
	if _, present := fields["approved_by"]; !present {
		return fmt.Errorf("approved_by required")
	}
	if _, present := fields["events"]; !present {
		return fmt.Errorf("events required")
	}
	if _, present := fields["operator_events"]; !present {
		return fmt.Errorf("operator_events required")
	}
	type alias Session
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	return decoder.Decode((*alias)(v))
}

func (v Session) Validate() error {
	if v.ContractVersion != "bifrost.device-lifecycle/v1" {
		return fmt.Errorf("contract_version value")
	}
	if !regexp.MustCompile("^[a-f0-9]{64}$").MatchString(v.SessionId) {
		return fmt.Errorf("session_id pattern")
	}
	if err := v.Registration.Validate(); err != nil {
		return err
	}
	if !regexp.MustCompile("^[a-f0-9]{64}$").MatchString(v.IdentitySha256) {
		return fmt.Errorf("identity_sha256 pattern")
	}
	if !regexp.MustCompile("^[A-Z2-9]{8}$").MatchString(v.PairingCode) {
		return fmt.Errorf("pairing_code pattern")
	}
	if v.State != "registered" && v.State != "awaiting_claim" && v.State != "claimed" && v.State != "approved" && v.State != "handoff_ready" && v.State != "completed" && v.State != "intervention_required" && v.State != "rejected" && v.State != "expired" && v.State != "cancelled" && v.State != "failed" {
		return fmt.Errorf("state value")
	}
	if v.Revision < 1 {
		return fmt.Errorf("revision range")
	}
	if v.LastSequence < 0 {
		return fmt.Errorf("last_sequence range")
	}
	if _, err := time.Parse(time.RFC3339, v.CreatedAt); err != nil {
		return fmt.Errorf("created_at timestamp")
	}
	if _, err := time.Parse(time.RFC3339, v.ExpiresAt); err != nil {
		return fmt.Errorf("expires_at timestamp")
	}
	if _, err := time.Parse(time.RFC3339, v.LastSeenAt); err != nil {
		return fmt.Errorf("last_seen_at timestamp")
	}
	if v.Assignment != nil {
		if err := v.Assignment.Validate(); err != nil {
			return err
		}
	}
	if v.ClaimedBy != nil && !regexp.MustCompile("^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$").MatchString(*v.ClaimedBy) {
		return fmt.Errorf("claimed_by uuid")
	}
	if v.ApprovedBy != nil && !regexp.MustCompile("^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$").MatchString(*v.ApprovedBy) {
		return fmt.Errorf("approved_by uuid")
	}
	if v.Events == nil {
		return fmt.Errorf("events required")
	}
	if len(v.Events) > 4096 {
		return fmt.Errorf("events limit")
	}
	for _, item := range v.Events {
		if err := item.Validate(); err != nil {
			return err
		}
	}
	if v.OperatorEvents == nil {
		return fmt.Errorf("operator_events required")
	}
	if len(v.OperatorEvents) > 256 {
		return fmt.Errorf("operator_events limit")
	}
	for _, item := range v.OperatorEvents {
		if err := item.Validate(); err != nil {
			return err
		}
	}
	if v.ResumeState != "" && (v.ResumeState != "registered" && v.ResumeState != "awaiting_claim" && v.ResumeState != "claimed" && v.ResumeState != "approved" && v.ResumeState != "handoff_ready" && v.ResumeState != "completed" && v.ResumeState != "intervention_required" && v.ResumeState != "rejected" && v.ResumeState != "expired" && v.ResumeState != "cancelled" && v.ResumeState != "failed") {
		return fmt.Errorf("resume_state value")
	}
	return nil
}

type OperatorEvent struct {
	EventId   string `json:"event_id"`
	Revision  int64  `json:"revision"`
	ActorId   string `json:"actor_id"`
	Operation string `json:"operation"`
	Timestamp string `json:"timestamp"`
	State     string `json:"state"`
}

func (v *OperatorEvent) UnmarshalJSON(data []byte) error {
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(data, &fields); err != nil {
		return err
	}
	if _, present := fields["event_id"]; !present {
		return fmt.Errorf("event_id required")
	}
	if _, present := fields["revision"]; !present {
		return fmt.Errorf("revision required")
	}
	if _, present := fields["actor_id"]; !present {
		return fmt.Errorf("actor_id required")
	}
	if _, present := fields["operation"]; !present {
		return fmt.Errorf("operation required")
	}
	if _, present := fields["timestamp"]; !present {
		return fmt.Errorf("timestamp required")
	}
	if _, present := fields["state"]; !present {
		return fmt.Errorf("state required")
	}
	type alias OperatorEvent
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	return decoder.Decode((*alias)(v))
}

func (v OperatorEvent) Validate() error {
	if !regexp.MustCompile("^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$").MatchString(v.EventId) {
		return fmt.Errorf("event_id pattern")
	}
	if v.Revision < 1 {
		return fmt.Errorf("revision range")
	}
	if !regexp.MustCompile("^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$").MatchString(v.ActorId) {
		return fmt.Errorf("actor_id pattern")
	}
	if v.Operation != "claim" && v.Operation != "approve" && v.Operation != "reject" && v.Operation != "cancel" && v.Operation != "resume" {
		return fmt.Errorf("operation value")
	}
	if _, err := time.Parse(time.RFC3339, v.Timestamp); err != nil {
		return fmt.Errorf("timestamp timestamp")
	}
	if v.State != "registered" && v.State != "awaiting_claim" && v.State != "claimed" && v.State != "approved" && v.State != "handoff_ready" && v.State != "completed" && v.State != "intervention_required" && v.State != "rejected" && v.State != "expired" && v.State != "cancelled" && v.State != "failed" {
		return fmt.Errorf("state value")
	}
	return nil
}
