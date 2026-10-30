// SPDX-FileCopyrightText: 2026 Nik Ogura <nik.ogura@gmail.com>
//
// SPDX-License-Identifier: Apache-2.0

package config

// Terraform resource names of the kinds this provider generates.
const (
	tfTeam                                   = "rootly_team"
	tfEscalationPolicy                       = "rootly_escalation_policy"
	tfEscalationLevel                        = "rootly_escalation_level"
	tfEscalationPath                         = "rootly_escalation_path"
	tfSchedule                               = "rootly_schedule"
	tfScheduleRotation                       = "rootly_schedule_rotation"
	tfScheduleRotationUser                   = "rootly_schedule_rotation_user"
	tfHeartbeat                              = "rootly_heartbeat"
	tfAlertsSource                           = "rootly_alerts_source"
	tfAlertRoute                             = "rootly_alert_route"
	tfAlertUrgency                           = "rootly_alert_urgency"
	tfService                                = "rootly_service"
	tfFunctionality                          = "rootly_functionality"
	tfIncidentType                           = "rootly_incident_type"
	tfSeverity                               = "rootly_severity"
	tfWorkflowAlert                          = "rootly_workflow_alert"
	tfWorkflowIncident                       = "rootly_workflow_incident"
	tfWorkflowTaskCreateIncident             = "rootly_workflow_task_create_incident"
	tfWorkflowTaskPublishIncident            = "rootly_workflow_task_publish_incident"
	tfWorkflowTaskPageRootlyOnCallResponders = "rootly_workflow_task_page_rootly_on_call_responders"
	tfStatusPage                             = "rootly_status_page"
	tfStatusPageComponent                    = "rootly_status_page_component"
	tfStatusPageComponentGroup               = "rootly_status_page_component_group"
)

// API groups the kinds are generated into, by resource family.
const (
	groupTeam       = "team"
	groupEscalation = "escalation"
	groupSchedule   = "schedule"
	groupHeartbeat  = "heartbeat"
	groupAlerts     = "alerts"
	groupService    = "service"
	groupIncident   = "incident"
	groupWorkflow   = "workflow"
	groupStatusPage = "statuspage"
)
