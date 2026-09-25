// SPDX-FileCopyrightText: 2024 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/crossplane/upjet/v2/pkg/controller"

	alertssource "github.com/nikogura/provider-rootly/internal/controller/namespaced/alerts/alertssource"
	escalationlevel "github.com/nikogura/provider-rootly/internal/controller/namespaced/escalation/escalationlevel"
	escalationpolicy "github.com/nikogura/provider-rootly/internal/controller/namespaced/escalation/escalationpolicy"
	heartbeat "github.com/nikogura/provider-rootly/internal/controller/namespaced/heartbeat/heartbeat"
	providerconfig "github.com/nikogura/provider-rootly/internal/controller/namespaced/providerconfig"
	schedule "github.com/nikogura/provider-rootly/internal/controller/namespaced/schedule/schedule"
	team "github.com/nikogura/provider-rootly/internal/controller/namespaced/team/team"
)

// Setup creates all controllers with the supplied logger and adds them to
// the supplied manager.
func Setup(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		alertssource.Setup,
		escalationlevel.Setup,
		escalationpolicy.Setup,
		heartbeat.Setup,
		providerconfig.Setup,
		schedule.Setup,
		team.Setup,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}

// SetupGated creates all controllers with the supplied logger and adds them to
// the supplied manager gated.
func SetupGated(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		alertssource.SetupGated,
		escalationlevel.SetupGated,
		escalationpolicy.SetupGated,
		heartbeat.SetupGated,
		providerconfig.SetupGated,
		schedule.SetupGated,
		team.SetupGated,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}

// SetupWebhookWithManager registers conversion webhooks for all resource kinds in the group.
func SetupWebhookWithManager(mgr ctrl.Manager) error {
	for _, setup := range []func(ctrl.Manager) error{
		alertssource.SetupWebhookWithManager,
		escalationlevel.SetupWebhookWithManager,
		escalationpolicy.SetupWebhookWithManager,
		heartbeat.SetupWebhookWithManager,
		providerconfig.SetupWebhookWithManager,
		schedule.SetupWebhookWithManager,
		team.SetupWebhookWithManager,
	} {
		if err := setup(mgr); err != nil {
			return err
		}
	}
	return nil
}
