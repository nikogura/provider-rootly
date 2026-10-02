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
	"rootly_team":                   "A Rootly team: a group of users that owns schedules, escalation policies and alert sources.",
	"rootly_escalation_policy":      "An escalation policy: who is paged for an alert, in what order, and how many times.",
	"rootly_escalation_level":       "One level of an escalation policy: who is notified at that step, and the delay before the next.",
	"rootly_schedule":               "An on-call schedule. Who is on call, and when, is set by the rotations that belong to it.",
	"rootly_schedule_rotation":      "A rotation within a schedule: the shift pattern, the hours it is active, and its members.",
	"rootly_schedule_rotation_user": "One user's position in a schedule rotation. Deprecated upstream in favour of the rotation's inline members.",
	"rootly_heartbeat":              "A heartbeat: Rootly raises an alert when an expected ping does not arrive within the interval.",
	"rootly_alerts_source":          "An alert source: an integration that sends alerts into Rootly, by webhook or by email.",
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
