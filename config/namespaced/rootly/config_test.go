// SPDX-FileCopyrightText: 2026 Nik Ogura <nik.ogura@gmail.com>
//
// SPDX-License-Identifier: Apache-2.0

package rootly

import (
	"os"
	"testing"

	ujconfig "github.com/crossplane/upjet/v2/pkg/config"
	"github.com/stretchr/testify/assert"
)

// generated is every resource this provider generates, as the include list
// the provider builds from its external-name table.
var generated = []string{
	tfTeam, tfEscalationPolicy, tfEscalationLevel,
	tfSchedule, tfScheduleRotation, tfScheduleRotationUser,
	tfHeartbeat, tfAlertsSource, tfAlertRoute, tfAlertUrgency,
	tfService, tfFunctionality, tfIncidentType, tfSeverity,
	tfWorkflowAlert, tfWorkflowIncident, tfWorkflowTaskCreateIncident, tfWorkflowTaskPublishIncident,
	tfWorkflowTaskPageRootlyOnCallResponders,
	tfStatusPage, tfStatusPageComponent, tfStatusPageComponentGroup,
}

// configuredProvider runs the same pipeline the provider does, on the real
// committed schema. Configure only registers configurator callbacks; they
// run inside ConfigureResources.
func configuredProvider(t *testing.T) (p *ujconfig.Provider) {
	t.Helper()

	schema, err := os.ReadFile("../../schema.json")
	if err != nil {
		t.Fatalf("cannot read provider schema: %v", err)
	}
	meta, err := os.ReadFile("../../provider-metadata.yaml")
	if err != nil {
		t.Fatalf("cannot read provider metadata: %v", err)
	}

	include := make([]string, len(generated))
	for i, name := range generated {
		include[i] = name + "$"
	}

	p = ujconfig.NewProvider(schema, "rootly", "github.com/nikogura/provider-rootly", meta,
		ujconfig.WithRootGroup("rootly.crossplane.io"),
		ujconfig.WithIncludeList(include))

	Configure(p)
	p.ConfigureResources()

	return p
}

func TestConfigureNamesEveryKind(t *testing.T) {
	t.Parallel()

	p := configuredProvider(t)

	tests := []struct {
		name       string
		resource   string
		shortGroup string
		kind       string
	}{
		{name: "Team", resource: tfTeam, shortGroup: groupTeam, kind: "Team"},
		{name: "EscalationPolicy", resource: tfEscalationPolicy, shortGroup: groupEscalation, kind: "EscalationPolicy"},
		{name: "EscalationLevel", resource: tfEscalationLevel, shortGroup: groupEscalation, kind: "EscalationLevel"},
		{name: "Schedule", resource: tfSchedule, shortGroup: groupSchedule, kind: "Schedule"},
		{name: "ScheduleRotation", resource: tfScheduleRotation, shortGroup: groupSchedule, kind: "ScheduleRotation"},
		{name: "ScheduleRotationUser", resource: tfScheduleRotationUser, shortGroup: groupSchedule, kind: "ScheduleRotationUser"},
		{name: "Heartbeat", resource: tfHeartbeat, shortGroup: groupHeartbeat, kind: "Heartbeat"},
		{name: "AlertsSource", resource: tfAlertsSource, shortGroup: groupAlerts, kind: "AlertsSource"},
		{name: "AlertRoute", resource: tfAlertRoute, shortGroup: groupAlerts, kind: "AlertRoute"},
		{name: "AlertUrgency", resource: tfAlertUrgency, shortGroup: groupAlerts, kind: "AlertUrgency"},
		{name: "Service", resource: tfService, shortGroup: groupService, kind: "Service"},
		{name: "Functionality", resource: tfFunctionality, shortGroup: groupService, kind: "Functionality"},
		{name: "IncidentType", resource: tfIncidentType, shortGroup: groupIncident, kind: "IncidentType"},
		{name: "Severity", resource: tfSeverity, shortGroup: groupIncident, kind: "Severity"},
		{name: "WorkflowAlert", resource: tfWorkflowAlert, shortGroup: groupWorkflow, kind: "WorkflowAlert"},
		{name: "WorkflowIncident", resource: tfWorkflowIncident, shortGroup: groupWorkflow, kind: "WorkflowIncident"},
		{name: "WorkflowTaskCreateIncident", resource: tfWorkflowTaskCreateIncident, shortGroup: groupWorkflow, kind: "WorkflowTaskCreateIncident"},
		{name: "WorkflowTaskPublishIncident", resource: tfWorkflowTaskPublishIncident, shortGroup: groupWorkflow, kind: "WorkflowTaskPublishIncident"},
		{name: "WorkflowTaskPageRootlyOnCallResponders", resource: tfWorkflowTaskPageRootlyOnCallResponders, shortGroup: groupWorkflow, kind: "WorkflowTaskPageRootlyOnCallResponders"},
		{name: "StatusPage", resource: tfStatusPage, shortGroup: groupStatusPage, kind: "StatusPage"},
		{name: "StatusPageComponent", resource: tfStatusPageComponent, shortGroup: groupStatusPage, kind: "StatusPageComponent"},
		{name: "StatusPageComponentGroup", resource: tfStatusPageComponentGroup, shortGroup: groupStatusPage, kind: "StatusPageComponentGroup"},
	}

	if len(tests) != len(generated) {
		t.Fatalf("%d kinds named, %d generated", len(tests), len(generated))
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r, ok := p.Resources[tt.resource]
			if !ok {
				t.Fatalf("resource %q not present", tt.resource)
			}
			if r.ShortGroup != tt.shortGroup {
				t.Errorf("ShortGroup = %q, want %q", r.ShortGroup, tt.shortGroup)
			}
			if r.Kind != tt.kind {
				t.Errorf("Kind = %q, want %q", r.Kind, tt.kind)
			}
		})
	}
}

// Every id field that names another generated kind resolves it by name,
// never by a pasted Rootly UUID. A polymorphic field resolves the kind the
// alerting flow needs.
func TestConfigureReferences(t *testing.T) {
	t.Parallel()

	p := configuredProvider(t)

	references := []struct {
		resource string
		field    string
		target   string
	}{
		{resource: tfEscalationLevel, field: "escalation_policy_id", target: tfEscalationPolicy},
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

	for _, want := range references {
		ref, ok := p.Resources[want.resource].References[want.field]
		if !ok {
			t.Errorf("%s.%s has no cross-resource reference", want.resource, want.field)

			continue
		}
		if ref.TerraformName != want.target {
			t.Errorf("%s.%s references %q, want %s", want.resource, want.field, ref.TerraformName, want.target)
		}
	}
}

func TestConfigureFieldHandling(t *testing.T) {
	t.Parallel()

	p := configuredProvider(t)

	ignored := p.Resources[tfAlertsSource].LateInitializer.IgnoredFields
	assert.ElementsMatch(t, []string{"alert_template_attributes", "alert_source_fields_attributes"}, ignored,
		"the conflicting AlertsSource blocks must not be late-initialized")

	password := p.Resources[tfStatusPage].TerraformResource.Schema["authentication_password"]
	assert.True(t, password.Sensitive, "the status page password must be a secret reference, not a spec field")

	secret := p.Resources[tfHeartbeat].TerraformResource.Schema["secret"]
	assert.True(t, secret.Sensitive, "the heartbeat ping token must go to the connection secret, not the status")
}
