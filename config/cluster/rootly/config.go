// SPDX-FileCopyrightText: 2026 Nik Ogura <nik.ogura@gmail.com>
//
// SPDX-License-Identifier: Apache-2.0

package rootly

import (
	ujconfig "github.com/crossplane/upjet/v2/pkg/config"
)

// Configure shapes the generated kinds. Groups follow the resource family and
// kinds carry the full name -- upjet's default would split
// rootly_escalation_policy into group "escalation", kind "Policy", and a kind
// named Policy tells you nothing at a glance in a cluster full of policies.
//
// Wherever a field holds the Rootly UUID of another kind this provider
// generates, it is made a reference so the object can be named instead.
// Where the target kind varies with a sibling type field, the reference binds
// the kind the alerting flow needs, and the other targets take the raw id.
func Configure(p *ujconfig.Provider) {
	configureOnCall(p)
	configureAlerting(p)
	configureCatalog(p)
	configureWorkflows(p)
	configureStatusPages(p)
}

// configureOnCall covers teams, escalation and schedules.
func configureOnCall(p *ujconfig.Provider) {
	p.AddResourceConfigurator(tfTeam, func(r *ujconfig.Resource) {
		r.ShortGroup = groupTeam
		r.Kind = "Team"
	})
	p.AddResourceConfigurator(tfEscalationPolicy, func(r *ujconfig.Resource) {
		r.ShortGroup = groupEscalation
		r.Kind = "EscalationPolicy"
	})
	p.AddResourceConfigurator(tfEscalationLevel, func(r *ujconfig.Resource) {
		r.ShortGroup = groupEscalation
		r.Kind = "EscalationLevel"
		// A level is meaningless outside its policy, and the id is a Rootly
		// UUID nobody should be copy-pasting: reference the policy by name.
		r.References["escalation_policy_id"] = ujconfig.Reference{
			TerraformName: tfEscalationPolicy,
		}
	})
	p.AddResourceConfigurator(tfSchedule, func(r *ujconfig.Resource) {
		r.ShortGroup = groupSchedule
		r.Kind = "Schedule"
	})
	p.AddResourceConfigurator(tfScheduleRotation, func(r *ujconfig.Resource) {
		r.ShortGroup = groupSchedule
		r.Kind = "ScheduleRotation"
		// A rotation belongs to exactly one schedule: reference it by name.
		r.References["schedule_id"] = ujconfig.Reference{
			TerraformName: tfSchedule,
		}
	})
	p.AddResourceConfigurator(tfScheduleRotationUser, func(r *ujconfig.Resource) {
		r.ShortGroup = groupSchedule
		r.Kind = "ScheduleRotationUser"
		// One user's place in one rotation. Upstream deprecates this resource
		// in favour of the rotation's inline scheduleRotationMembers, which
		// needs the account's schedule-nesting feature; this kind is the way
		// to place users without it. Use one or the other on a rotation, never
		// both -- they would each reconcile the other's members away.
		r.References["schedule_rotation_id"] = ujconfig.Reference{
			TerraformName: tfScheduleRotation,
		}
	})
	p.AddResourceConfigurator(tfHeartbeat, func(r *ujconfig.Resource) {
		r.ShortGroup = groupHeartbeat
		r.Kind = "Heartbeat"
		// The secret is the bearer token that authenticates a ping, so anyone
		// holding it can keep the heartbeat alive. Upstream does not flag it
		// as sensitive, which would publish it in clear in the status of
		// every Heartbeat. Flagging it here drops it from the status and
		// writes it to the connection secret instead.
		r.TerraformResource.Schema["secret"].Sensitive = true
	})
}

// configureAlerting covers where alerts come from and where they go.
func configureAlerting(p *ujconfig.Provider) {
	p.AddResourceConfigurator(tfAlertsSource, func(r *ujconfig.Resource) {
		r.ShortGroup = groupAlerts
		r.Kind = "AlertsSource"
		// alert_template_attributes and alert_source_fields_attributes conflict
		// upstream, and the upstream Read always stores the template block,
		// as one empty element when the API has none. Late-initializing
		// either would put both in the spec and have every later refresh
		// rejected. What the user declares is what is sent.
		r.LateInitializer = ujconfig.LateInitializer{
			IgnoredFields: []string{"alert_template_attributes", "alert_source_fields_attributes"},
		}
	})
	p.AddResourceConfigurator(tfAlertRoute, func(r *ujconfig.Resource) {
		r.ShortGroup = groupAlerts
		r.Kind = "AlertRoute"
		r.References["alerts_source_ids"] = ujconfig.Reference{
			TerraformName: tfAlertsSource,
		}
		r.References["owning_team_ids"] = ujconfig.Reference{
			TerraformName: tfTeam,
		}
		// A rule's destination is a Service, Group, Functionality or
		// EscalationPolicy, chosen by targetType. The reference resolves a
		// Service, the target that ties an alert to an incident and a status
		// page; the other three take the Rootly id in targetId.
		r.References["rules.destinations.target_id"] = ujconfig.Reference{
			TerraformName: tfService,
		}
	})
	p.AddResourceConfigurator(tfAlertUrgency, func(r *ujconfig.Resource) {
		r.ShortGroup = groupAlerts
		r.Kind = "AlertUrgency"
	})
}

// configureCatalog covers the service catalog and incident taxonomy.
func configureCatalog(p *ujconfig.Provider) {
	p.AddResourceConfigurator(tfService, func(r *ujconfig.Resource) {
		r.ShortGroup = groupService
		r.Kind = "Service"
		r.References["escalation_policy_id"] = ujconfig.Reference{
			TerraformName: tfEscalationPolicy,
		}
		// Rootly calls teams "groups" in every id field.
		r.References["owner_group_ids"] = ujconfig.Reference{
			TerraformName: tfTeam,
		}
	})
	p.AddResourceConfigurator(tfFunctionality, func(r *ujconfig.Resource) {
		r.ShortGroup = groupService
		r.Kind = "Functionality"
		r.References["escalation_policy_id"] = ujconfig.Reference{
			TerraformName: tfEscalationPolicy,
		}
		r.References["owner_group_ids"] = ujconfig.Reference{
			TerraformName: tfTeam,
		}
		r.References["service_ids"] = ujconfig.Reference{
			TerraformName: tfService,
		}
	})
	p.AddResourceConfigurator(tfIncidentType, func(r *ujconfig.Resource) {
		r.ShortGroup = groupIncident
		r.Kind = "IncidentType"
	})
	p.AddResourceConfigurator(tfSeverity, func(r *ujconfig.Resource) {
		r.ShortGroup = groupIncident
		r.Kind = "Severity"
	})
}

// configureWorkflows covers the alert- and incident-triggered workflows and
// the steps that page responders and turn an alert into a published incident.
// Upstream has one workflow resource per trigger genre, and a task is a step
// of exactly one workflow. Only the incident genre runs its steps with the
// incident in scope, so a task that acts on the incident, such as publishing
// it, belongs to a WorkflowIncident. A task's workflowId reference resolves
// the genre its example uses; the other genre takes the raw id.
func configureWorkflows(p *ujconfig.Provider) {
	p.AddResourceConfigurator(tfWorkflowAlert, func(r *ujconfig.Resource) {
		r.ShortGroup = groupWorkflow
		r.Kind = "WorkflowAlert"
		r.References["group_ids"] = ujconfig.Reference{
			TerraformName: tfTeam,
		}
		r.References["service_ids"] = ujconfig.Reference{
			TerraformName: tfService,
		}
		r.References["functionality_ids"] = ujconfig.Reference{
			TerraformName: tfFunctionality,
		}
		r.References["incident_type_ids"] = ujconfig.Reference{
			TerraformName: tfIncidentType,
		}
		r.References["severity_ids"] = ujconfig.Reference{
			TerraformName: tfSeverity,
		}
	})
	p.AddResourceConfigurator(tfWorkflowIncident, func(r *ujconfig.Resource) {
		r.ShortGroup = groupWorkflow
		r.Kind = "WorkflowIncident"
		r.References["group_ids"] = ujconfig.Reference{
			TerraformName: tfTeam,
		}
		r.References["service_ids"] = ujconfig.Reference{
			TerraformName: tfService,
		}
		r.References["functionality_ids"] = ujconfig.Reference{
			TerraformName: tfFunctionality,
		}
		r.References["incident_type_ids"] = ujconfig.Reference{
			TerraformName: tfIncidentType,
		}
		r.References["severity_ids"] = ujconfig.Reference{
			TerraformName: tfSeverity,
		}
	})
	p.AddResourceConfigurator(tfWorkflowTaskCreateIncident, func(r *ujconfig.Resource) {
		r.ShortGroup = groupWorkflow
		r.Kind = "WorkflowTaskCreateIncident"
		r.References["workflow_id"] = ujconfig.Reference{
			TerraformName: tfWorkflowAlert,
		}
		r.References["task_params.group_ids"] = ujconfig.Reference{
			TerraformName: tfTeam,
		}
		r.References["task_params.service_ids"] = ujconfig.Reference{
			TerraformName: tfService,
		}
		r.References["task_params.functionality_ids"] = ujconfig.Reference{
			TerraformName: tfFunctionality,
		}
		r.References["task_params.incident_type_ids"] = ujconfig.Reference{
			TerraformName: tfIncidentType,
		}
		r.References["task_params.severity_id"] = ujconfig.Reference{
			TerraformName: tfSeverity,
		}
	})
	p.AddResourceConfigurator(tfWorkflowTaskPublishIncident, func(r *ujconfig.Resource) {
		r.ShortGroup = groupWorkflow
		r.Kind = "WorkflowTaskPublishIncident"
		// Publishing acts on the incident, which only an incident-triggered
		// workflow has in scope; in an alert workflow this step runs and
		// posts nothing.
		r.References["workflow_id"] = ujconfig.Reference{
			TerraformName: tfWorkflowIncident,
		}
		r.References["task_params.status_page_id"] = ujconfig.Reference{
			TerraformName: tfStatusPage,
		}
		r.References["task_params.status_page_ids"] = ujconfig.Reference{
			TerraformName: tfStatusPage,
		}
	})
	p.AddResourceConfigurator(tfWorkflowTaskPageRootlyOnCallResponders, func(r *ujconfig.Resource) {
		r.ShortGroup = groupWorkflow
		r.Kind = "WorkflowTaskPageRootlyOnCallResponders"
		r.References["workflow_id"] = ujconfig.Reference{
			TerraformName: tfWorkflowAlert,
		}
		r.References["task_params.alert_urgency_id"] = ujconfig.Reference{
			TerraformName: tfAlertUrgency,
		}
		// Who is paged is one of five {id, name} maps -- serviceTarget,
		// groupTarget, functionalityTarget, escalationPolicyTarget,
		// userTarget. A reference resolves a scalar id, not a map value, so
		// each target takes the Rootly id and name as written.
	})
}

// configureStatusPages covers the page and what is shown on it.
func configureStatusPages(p *ujconfig.Provider) {
	p.AddResourceConfigurator(tfStatusPage, func(r *ujconfig.Resource) {
		r.ShortGroup = groupStatusPage
		r.Kind = "StatusPage"
		// Upstream does not flag the page password as sensitive, which would
		// put it in the spec and the status in clear. Flagging it here makes
		// it a secret reference in the spec and keeps it out of the status.
		r.TerraformResource.Schema["authentication_password"].Sensitive = true
		r.References["service_ids"] = ujconfig.Reference{
			TerraformName: tfService,
		}
		r.References["functionality_ids"] = ujconfig.Reference{
			TerraformName: tfFunctionality,
		}
	})
	p.AddResourceConfigurator(tfStatusPageComponentGroup, func(r *ujconfig.Resource) {
		r.ShortGroup = groupStatusPage
		r.Kind = "StatusPageComponentGroup"
		r.References["status_page_id"] = ujconfig.Reference{
			TerraformName: tfStatusPage,
		}
	})
	p.AddResourceConfigurator(tfStatusPageComponent, func(r *ujconfig.Resource) {
		r.ShortGroup = groupStatusPage
		r.Kind = "StatusPageComponent"
		r.References["status_page_id"] = ujconfig.Reference{
			TerraformName: tfStatusPage,
		}
		r.References["status_page_component_group_id"] = ujconfig.Reference{
			TerraformName: tfStatusPageComponentGroup,
		}
		// A component mirrors a Service or a Functionality, chosen by
		// sourceType. The reference resolves a Service; a Functionality takes
		// the Rootly id in sourceId.
		r.References["source_id"] = ujconfig.Reference{
			TerraformName: tfService,
		}
	})
}
