// SPDX-FileCopyrightText: 2024 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/crossplane/upjet/v2/pkg/controller"

	alertroute "github.com/nikogura/provider-rootly/internal/controller/cluster/alerts/alertroute"
	alertssource "github.com/nikogura/provider-rootly/internal/controller/cluster/alerts/alertssource"
	alerturgency "github.com/nikogura/provider-rootly/internal/controller/cluster/alerts/alerturgency"
	escalationlevel "github.com/nikogura/provider-rootly/internal/controller/cluster/escalation/escalationlevel"
	escalationpath "github.com/nikogura/provider-rootly/internal/controller/cluster/escalation/escalationpath"
	escalationpolicy "github.com/nikogura/provider-rootly/internal/controller/cluster/escalation/escalationpolicy"
	heartbeat "github.com/nikogura/provider-rootly/internal/controller/cluster/heartbeat/heartbeat"
	incidenttype "github.com/nikogura/provider-rootly/internal/controller/cluster/incident/incidenttype"
	severity "github.com/nikogura/provider-rootly/internal/controller/cluster/incident/severity"
	providerconfig "github.com/nikogura/provider-rootly/internal/controller/cluster/providerconfig"
	schedule "github.com/nikogura/provider-rootly/internal/controller/cluster/schedule/schedule"
	schedulerotation "github.com/nikogura/provider-rootly/internal/controller/cluster/schedule/schedulerotation"
	schedulerotationuser "github.com/nikogura/provider-rootly/internal/controller/cluster/schedule/schedulerotationuser"
	functionality "github.com/nikogura/provider-rootly/internal/controller/cluster/service/functionality"
	service "github.com/nikogura/provider-rootly/internal/controller/cluster/service/service"
	statuspage "github.com/nikogura/provider-rootly/internal/controller/cluster/statuspage/statuspage"
	statuspagecomponent "github.com/nikogura/provider-rootly/internal/controller/cluster/statuspage/statuspagecomponent"
	statuspagecomponentgroup "github.com/nikogura/provider-rootly/internal/controller/cluster/statuspage/statuspagecomponentgroup"
	team "github.com/nikogura/provider-rootly/internal/controller/cluster/team/team"
	workflowalert "github.com/nikogura/provider-rootly/internal/controller/cluster/workflow/workflowalert"
	workflowincident "github.com/nikogura/provider-rootly/internal/controller/cluster/workflow/workflowincident"
	workflowtaskcreateincident "github.com/nikogura/provider-rootly/internal/controller/cluster/workflow/workflowtaskcreateincident"
	workflowtaskpagerootlyoncallresponders "github.com/nikogura/provider-rootly/internal/controller/cluster/workflow/workflowtaskpagerootlyoncallresponders"
	workflowtaskpublishincident "github.com/nikogura/provider-rootly/internal/controller/cluster/workflow/workflowtaskpublishincident"
)

// Setup creates all controllers with the supplied logger and adds them to
// the supplied manager.
func Setup(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		alertroute.Setup,
		alertssource.Setup,
		alerturgency.Setup,
		escalationlevel.Setup,
		escalationpath.Setup,
		escalationpolicy.Setup,
		heartbeat.Setup,
		incidenttype.Setup,
		severity.Setup,
		providerconfig.Setup,
		schedule.Setup,
		schedulerotation.Setup,
		schedulerotationuser.Setup,
		functionality.Setup,
		service.Setup,
		statuspage.Setup,
		statuspagecomponent.Setup,
		statuspagecomponentgroup.Setup,
		team.Setup,
		workflowalert.Setup,
		workflowincident.Setup,
		workflowtaskcreateincident.Setup,
		workflowtaskpagerootlyoncallresponders.Setup,
		workflowtaskpublishincident.Setup,
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
		alertroute.SetupGated,
		alertssource.SetupGated,
		alerturgency.SetupGated,
		escalationlevel.SetupGated,
		escalationpath.SetupGated,
		escalationpolicy.SetupGated,
		heartbeat.SetupGated,
		incidenttype.SetupGated,
		severity.SetupGated,
		providerconfig.SetupGated,
		schedule.SetupGated,
		schedulerotation.SetupGated,
		schedulerotationuser.SetupGated,
		functionality.SetupGated,
		service.SetupGated,
		statuspage.SetupGated,
		statuspagecomponent.SetupGated,
		statuspagecomponentgroup.SetupGated,
		team.SetupGated,
		workflowalert.SetupGated,
		workflowincident.SetupGated,
		workflowtaskcreateincident.SetupGated,
		workflowtaskpagerootlyoncallresponders.SetupGated,
		workflowtaskpublishincident.SetupGated,
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
		alertroute.SetupWebhookWithManager,
		alertssource.SetupWebhookWithManager,
		alerturgency.SetupWebhookWithManager,
		escalationlevel.SetupWebhookWithManager,
		escalationpath.SetupWebhookWithManager,
		escalationpolicy.SetupWebhookWithManager,
		heartbeat.SetupWebhookWithManager,
		incidenttype.SetupWebhookWithManager,
		severity.SetupWebhookWithManager,
		providerconfig.SetupWebhookWithManager,
		schedule.SetupWebhookWithManager,
		schedulerotation.SetupWebhookWithManager,
		schedulerotationuser.SetupWebhookWithManager,
		functionality.SetupWebhookWithManager,
		service.SetupWebhookWithManager,
		statuspage.SetupWebhookWithManager,
		statuspagecomponent.SetupWebhookWithManager,
		statuspagecomponentgroup.SetupWebhookWithManager,
		team.SetupWebhookWithManager,
		workflowalert.SetupWebhookWithManager,
		workflowincident.SetupWebhookWithManager,
		workflowtaskcreateincident.SetupWebhookWithManager,
		workflowtaskpagerootlyoncallresponders.SetupWebhookWithManager,
		workflowtaskpublishincident.SetupWebhookWithManager,
	} {
		if err := setup(mgr); err != nil {
			return err
		}
	}
	return nil
}
