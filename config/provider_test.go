// SPDX-FileCopyrightText: 2026 Nik Ogura <nik.ogura@gmail.com>
//
// SPDX-License-Identifier: Apache-2.0

package config

import (
	"testing"

	ujconfig "github.com/crossplane/upjet/v2/pkg/config"
)

// Both providers must expose exactly the kinds this provider exists for,
// under the names the custom configurators assign. This exercises the full
// pipeline: schema parse, include list, external-name defaults and the
// per-resource configurators.
func TestGetProvider(t *testing.T) {
	t.Parallel()

	wantKinds := map[string]struct {
		shortGroup string
		kind       string
	}{
		tfTeam:                                   {shortGroup: groupTeam, kind: "Team"},
		tfEscalationPolicy:                       {shortGroup: groupEscalation, kind: "EscalationPolicy"},
		tfEscalationLevel:                        {shortGroup: groupEscalation, kind: "EscalationLevel"},
		tfEscalationPath:                         {shortGroup: groupEscalation, kind: "EscalationPath"},
		tfSchedule:                               {shortGroup: groupSchedule, kind: "Schedule"},
		tfScheduleRotation:                       {shortGroup: groupSchedule, kind: "ScheduleRotation"},
		tfScheduleRotationUser:                   {shortGroup: groupSchedule, kind: "ScheduleRotationUser"},
		tfHeartbeat:                              {shortGroup: groupHeartbeat, kind: "Heartbeat"},
		tfAlertsSource:                           {shortGroup: groupAlerts, kind: "AlertsSource"},
		tfAlertRoute:                             {shortGroup: groupAlerts, kind: "AlertRoute"},
		tfAlertUrgency:                           {shortGroup: groupAlerts, kind: "AlertUrgency"},
		tfService:                                {shortGroup: groupService, kind: "Service"},
		tfFunctionality:                          {shortGroup: groupService, kind: "Functionality"},
		tfIncidentType:                           {shortGroup: groupIncident, kind: "IncidentType"},
		tfSeverity:                               {shortGroup: groupIncident, kind: "Severity"},
		tfWorkflowAlert:                          {shortGroup: groupWorkflow, kind: "WorkflowAlert"},
		tfWorkflowIncident:                       {shortGroup: groupWorkflow, kind: "WorkflowIncident"},
		tfWorkflowTaskCreateIncident:             {shortGroup: groupWorkflow, kind: "WorkflowTaskCreateIncident"},
		tfWorkflowTaskPublishIncident:            {shortGroup: groupWorkflow, kind: "WorkflowTaskPublishIncident"},
		tfWorkflowTaskPageRootlyOnCallResponders: {shortGroup: groupWorkflow, kind: "WorkflowTaskPageRootlyOnCallResponders"},
		tfStatusPage:                             {shortGroup: groupStatusPage, kind: "StatusPage"},
		tfStatusPageComponent:                    {shortGroup: groupStatusPage, kind: "StatusPageComponent"},
		tfStatusPageComponentGroup:               {shortGroup: groupStatusPage, kind: "StatusPageComponentGroup"},
	}

	tests := []struct {
		name string
		get  func() *ujconfig.Provider
	}{
		{
			name: "Cluster",
			get:  GetProvider,
		},
		{
			name: "Namespaced",
			get:  GetProviderNamespaced,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			p := tt.get()
			if p == nil {
				t.Fatal("provider is nil")
			}

			if len(p.Resources) != len(wantKinds) {
				got := make([]string, 0, len(p.Resources))
				for name := range p.Resources {
					got = append(got, name)
				}
				t.Fatalf("expected exactly %d resources, got %d: %v", len(wantKinds), len(p.Resources), got)
			}

			for tfName, want := range wantKinds {
				r, ok := p.Resources[tfName]
				if !ok {
					t.Errorf("resource %q not generated", tfName)
					continue
				}
				if r.ShortGroup != want.shortGroup {
					t.Errorf("%s: ShortGroup = %q, want %q", tfName, r.ShortGroup, want.shortGroup)
				}
				if r.Kind != want.kind {
					t.Errorf("%s: Kind = %q, want %q", tfName, r.Kind, want.kind)
				}
				// Rootly assigns UUIDs; every resource must take its external
				// name from the provider.
				if !r.ExternalName.DisableNameInitializer {
					t.Errorf("%s: external name is not IdentifierFromProvider", tfName)
				}
			}

			// A child is referenced to its parent by name, never by a pasted
			// Rootly UUID.
			wantReferences := []struct {
				resource string
				field    string
				target   string
			}{
				{resource: tfEscalationLevel, field: "escalation_policy_id", target: tfEscalationPolicy},
				{resource: tfEscalationLevel, field: "escalation_policy_path_id", target: tfEscalationPath},
				{resource: tfEscalationPath, field: "escalation_policy_id", target: tfEscalationPolicy},
				{resource: tfEscalationPath, field: "after_deferral_path_id", target: tfEscalationPath},
				{resource: tfScheduleRotation, field: "schedule_id", target: tfSchedule},
				{resource: tfScheduleRotationUser, field: "schedule_rotation_id", target: tfScheduleRotation},
				{resource: tfAlertRoute, field: "alerts_source_ids", target: tfAlertsSource},
				{resource: tfAlertRoute, field: "owning_team_ids", target: tfTeam},
				{resource: tfAlertRoute, field: "rules.destinations.target_id", target: tfService},
				{resource: tfService, field: "escalation_policy_id", target: tfEscalationPolicy},
				{resource: tfService, field: "owner_group_ids", target: tfTeam},
				{resource: tfFunctionality, field: "escalation_policy_id", target: tfEscalationPolicy},
				{resource: tfFunctionality, field: "owner_group_ids", target: tfTeam},
				{resource: tfFunctionality, field: "service_ids", target: tfService},
				{resource: tfWorkflowAlert, field: "group_ids", target: tfTeam},
				{resource: tfWorkflowAlert, field: "service_ids", target: tfService},
				{resource: tfWorkflowAlert, field: "functionality_ids", target: tfFunctionality},
				{resource: tfWorkflowAlert, field: "incident_type_ids", target: tfIncidentType},
				{resource: tfWorkflowAlert, field: "severity_ids", target: tfSeverity},
				{resource: tfWorkflowIncident, field: "group_ids", target: tfTeam},
				{resource: tfWorkflowIncident, field: "service_ids", target: tfService},
				{resource: tfWorkflowIncident, field: "functionality_ids", target: tfFunctionality},
				{resource: tfWorkflowIncident, field: "incident_type_ids", target: tfIncidentType},
				{resource: tfWorkflowIncident, field: "severity_ids", target: tfSeverity},
				{resource: tfWorkflowTaskCreateIncident, field: "workflow_id", target: tfWorkflowAlert},
				{resource: tfWorkflowTaskCreateIncident, field: "task_params.group_ids", target: tfTeam},
				{resource: tfWorkflowTaskCreateIncident, field: "task_params.service_ids", target: tfService},
				{resource: tfWorkflowTaskCreateIncident, field: "task_params.functionality_ids", target: tfFunctionality},
				{resource: tfWorkflowTaskCreateIncident, field: "task_params.incident_type_ids", target: tfIncidentType},
				{resource: tfWorkflowTaskCreateIncident, field: "task_params.severity_id", target: tfSeverity},
				{resource: tfWorkflowTaskPublishIncident, field: "workflow_id", target: tfWorkflowIncident},
				{resource: tfWorkflowTaskPublishIncident, field: "task_params.status_page_id", target: tfStatusPage},
				{resource: tfWorkflowTaskPublishIncident, field: "task_params.status_page_ids", target: tfStatusPage},
				{resource: tfWorkflowTaskPageRootlyOnCallResponders, field: "workflow_id", target: tfWorkflowAlert},
				{resource: tfWorkflowTaskPageRootlyOnCallResponders, field: "task_params.alert_urgency_id", target: tfAlertUrgency},
				{resource: tfStatusPage, field: "service_ids", target: tfService},
				{resource: tfStatusPage, field: "functionality_ids", target: tfFunctionality},
				{resource: tfStatusPageComponentGroup, field: "status_page_id", target: tfStatusPage},
				{resource: tfStatusPageComponent, field: "status_page_id", target: tfStatusPage},
				{resource: tfStatusPageComponent, field: "status_page_component_group_id", target: tfStatusPageComponentGroup},
				{resource: tfStatusPageComponent, field: "source_id", target: tfService},
			}

			for _, want := range wantReferences {
				child, ok := p.Resources[want.resource]
				if !ok {
					t.Fatalf("%s not generated", want.resource)
				}
				ref, ok := child.References[want.field]
				if !ok {
					t.Fatalf("%s.%s has no cross-resource reference", want.resource, want.field)
				}
				if ref.TerraformName != want.target {
					t.Errorf("%s.%s references %q, want %s", want.resource, want.field, ref.TerraformName, want.target)
				}
			}

			// Upstream leaves the status page password unflagged. It must
			// not be generated as a plain spec field.
			page := p.Resources[tfStatusPage]
			if !page.TerraformResource.Schema["authentication_password"].Sensitive {
				t.Error("rootly_status_page.authentication_password is not sensitive")
			}

			// Upstream leaves the heartbeat ping token unflagged. It must
			// not be published in the status.
			heartbeat := p.Resources[tfHeartbeat]
			if !heartbeat.TerraformResource.Schema["secret"].Sensitive {
				t.Error("rootly_heartbeat.secret is not sensitive")
			}
		})
	}
}
