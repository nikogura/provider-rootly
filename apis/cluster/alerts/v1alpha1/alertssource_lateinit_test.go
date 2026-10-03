// SPDX-FileCopyrightText: 2026 Nik Ogura <nik.ogura@gmail.com>
//
// SPDX-License-Identifier: Apache-2.0

package v1alpha1

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestLateInitializeSkipsConflictingAlertBlocks guards against the provider
// writing an AlertsSource into a state it can never leave.
//
// The upstream Read always stores alert_template_attributes as a one-element
// list, with empty strings when the API has none, beside whatever
// alert_source_fields_attributes the API returned. The two blocks conflict,
// so if both were late-initialized into the spec, every later refresh would be
// rejected and the resource could neither be observed nor updated. Neither
// block may be late-initialized; what the user declares is what is sent.
func TestLateInitializeSkipsConflictingAlertBlocks(t *testing.T) {
	t.Parallel()

	name := "alertmanager-prod"
	tr := &AlertsSource{}
	tr.Spec.ForProvider.Name = &name

	observed := []byte(`{
		"name": "alertmanager-prod",
		"source_type": "alertmanager",
		"alert_template_attributes": [{"title": "", "description": "", "external_url": ""}],
		"alert_source_fields_attributes": [{"alert_field_id": "50dc639c", "template_body": ""}]
	}`)

	changed, err := tr.LateInitialize(observed)
	require.NoError(t, err)

	assert.True(t, changed, "fields outside the conflicting pair are still late-initialized")
	require.NotNil(t, tr.Spec.ForProvider.SourceType)
	assert.Equal(t, "alertmanager", *tr.Spec.ForProvider.SourceType)

	assert.Nil(t, tr.Spec.ForProvider.AlertTemplateAttributes, "alert_template_attributes must not be late-initialized")
	assert.Nil(t, tr.Spec.ForProvider.AlertSourceFieldsAttributes, "alert_source_fields_attributes must not be late-initialized")
}
