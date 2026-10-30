// SPDX-FileCopyrightText: 2026 Nik Ogura <nik.ogura@gmail.com>
//
// SPDX-License-Identifier: Apache-2.0

package config

import (
	"github.com/crossplane/upjet/v2/pkg/config"
	"github.com/crossplane/upjet/v2/pkg/registry"
)

// KindDescriptions says what each generated kind is. It becomes the kind's
// description in its CRD, which is what `kubectl explain` shows.
//
// The upstream documentation leaves most of these blank, so they are written
// here. Every resource in ExternalNameConfigs has an entry.
var KindDescriptions = map[string]string{
	tfTeam:                                   "A Rootly team: a group of users that owns schedules, escalation policies and alert sources.",
	tfEscalationPolicy:                       "An escalation policy: who is paged for an alert, in what order, and how many times.",
	tfEscalationLevel:                        "One level of an escalation policy: who is notified at that step, and the delay before the next.",
	tfEscalationPath:                         "A path through an escalation policy: the rules, such as alert urgency or working hours, that send an alert down it, and how often it repeats. Its levels are the EscalationLevel kinds that name it.",
	tfSchedule:                               "An on-call schedule. Who is on call, and when, is set by the rotations that belong to it.",
	tfScheduleRotation:                       "A rotation within a schedule: the shift pattern, the hours it is active, and its members.",
	tfScheduleRotationUser:                   "One user's position in a schedule rotation. Deprecated upstream in favour of the rotation's inline members.",
	tfHeartbeat:                              "A heartbeat: Rootly raises an alert when an expected ping does not arrive within the interval.",
	tfAlertsSource:                           "An alert source: an integration that sends alerts into Rootly, by webhook or by email.",
	tfAlertRoute:                             "An alert route: rules that match alerts arriving from one or more alert sources and send each to a service, team, functionality or escalation policy.",
	tfAlertUrgency:                           "An alert urgency: how pressing an alert is, such as Critical or Informational. A paging workflow step must name one.",
	tfService:                                "A service: a component of the system that alerts are routed to, incidents are attached to, and status page components report on.",
	tfFunctionality:                          "A functionality: a customer-facing capability, such as Login or Checkout, made up of one or more services.",
	tfIncidentType:                           "An incident type: a category an incident is filed under, such as Security or Outage.",
	tfSeverity:                               "A severity level for incidents, such as SEV0 or SEV1.",
	tfWorkflowAlert:                          "A workflow triggered by alert events. Its steps are the WorkflowTask kinds that reference it.",
	tfWorkflowIncident:                       "A workflow triggered by incident events, where the incident is in scope for its steps. Its steps are the WorkflowTask kinds that reference it.",
	tfWorkflowTaskCreateIncident:             "A workflow step that creates an incident, optionally attached to services, with a severity and type.",
	tfWorkflowTaskPublishIncident:            "A workflow step that publishes an incident update to a status page.",
	tfWorkflowTaskPageRootlyOnCallResponders: "A workflow step that pages the on-call responders of a service, team, functionality, escalation policy or user through Rootly On-Call, at a given alert urgency.",
	tfStatusPage:                             "A status page: a public or authenticated page reporting the state of services and functionalities.",
	tfStatusPageComponent:                    "A component shown on a status page, either ad hoc or mirroring a service or functionality.",
	tfStatusPageComponentGroup:               "A group that components are shown under on a status page.",
}

// DescribeKinds applies KindDescriptions to the resources this provider
// generates.
func DescribeKinds() (opt config.ResourceOption) {
	opt = func(r *config.Resource) {
		description, ok := KindDescriptions[r.Name]
		if !ok {
			return
		}

		if r.MetaResource == nil {
			r.MetaResource = &registry.Resource{}
		}

		r.MetaResource.Description = description
	}
	return opt
}
