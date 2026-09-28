/**
 * Copyright 2026 IBM Corp.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package provider_test

import (
	"testing"

	"github.com/IBM/ibmcloud-volume-interface/lib/provider"
	providerfake "github.com/IBM/ibmcloud-volume-interface/lib/provider/fake"
	providerfakes "github.com/IBM/ibmcloud-volume-interface/lib/provider/fakes"
	"github.com/stretchr/testify/assert"
)

func TestGroupSnapshotFakesImplementManager(t *testing.T) {
	var _ provider.GroupSnapshotManager = &providerfake.FakeSession{}
	var _ provider.GroupSnapshotManager = &providerfakes.Context{}
}

func TestFakeSessionDeleteGroupSnapshotCapturesSnapshotIDs(t *testing.T) {
	fakeSession := &providerfake.FakeSession{}
	snapshotIDs := []string{"snapshot-id-1", "snapshot-id-2"}

	assert.Nil(t, fakeSession.DeleteGroupSnapshot("group-snapshot-id", snapshotIDs))
	snapshotIDs[0] = "mutated-snapshot-id"

	groupSnapshotID, snapshotIDs := fakeSession.DeleteGroupSnapshotArgsForCall(0)
	assert.Equal(t, "group-snapshot-id", groupSnapshotID)
	assert.Equal(t, []string{"snapshot-id-1", "snapshot-id-2"}, snapshotIDs)
}

func TestFakeContextDeleteGroupSnapshotCapturesSnapshotIDs(t *testing.T) {
	fakeContext := &providerfakes.Context{}

	assert.Nil(t, fakeContext.DeleteGroupSnapshot("group-snapshot-id", []string{"snapshot-id-1", "snapshot-id-2"}))

	groupSnapshotID, snapshotIDs := fakeContext.DeleteGroupSnapshotArgsForCall(0)
	assert.Equal(t, "group-snapshot-id", groupSnapshotID)
	assert.Equal(t, []string{"snapshot-id-1", "snapshot-id-2"}, snapshotIDs)
}
