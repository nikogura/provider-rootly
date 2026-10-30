// SPDX-FileCopyrightText: 2026 Nik Ogura <nik.ogura@gmail.com>
//
// SPDX-License-Identifier: Apache-2.0

package v1alpha1

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestLateInitializeSkipsPathPosition guards against EscalationPaths that
// reorder each other on every reconcile.
//
// position is a path's place among all of its policy's paths, including the
// default one Rootly creates, so adding or removing a sibling path can move
// it. Copied into the spec, it would go stale with the next sibling, and the
// provider would send the old position back on every reconcile. It must not
// be late-initialized; it stays in status.atProvider.
func TestLateInitializeSkipsPathPosition(t *testing.T) {
	t.Parallel()

	name := "Critical during working hours"
	tr := &EscalationPath{}
	tr.Spec.ForProvider.Name = &name

	observed := []byte(`{
		"name": "Critical during working hours",
		"default": false,
		"notification_type": "audible",
		"path_type": "escalation",
		"match_mode": "match-all-rules",
		"position": 1,
		"repeat": false,
		"repeat_count": 2,
		"initial_delay": 0
	}`)

	changed, err := tr.LateInitialize(observed)
	require.NoError(t, err)

	assert.True(t, changed, "fields Rootly does not own are still late-initialized")
	require.NotNil(t, tr.Spec.ForProvider.NotificationType)
	assert.Equal(t, "audible", *tr.Spec.ForProvider.NotificationType)
	require.NotNil(t, tr.Spec.ForProvider.RepeatCount)
	assert.InDelta(t, 2, *tr.Spec.ForProvider.RepeatCount, 0)

	assert.Nil(t, tr.Spec.ForProvider.Position, "position must not be late-initialized")
}
