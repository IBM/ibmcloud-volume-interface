/**
 * Copyright 2020 IBM Corp.
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

// Package provider ...
package provider

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

var (
	logger *zap.Logger
)

func init() {
	logger, _ = zap.NewDevelopment()
}

func TestForDefaultVolumeProviderTypeAndName(t *testing.T) {
	ccf := &DefaultVolumeProvider{sess: nil}

	volumeProvider := ccf.ProviderName()

	volumeType := ccf.Type()

	volumeProviderName := ccf.GetProviderDisplayName()

	assert.Empty(t, volumeProviderName)
	assert.Empty(t, volumeProvider)
	assert.Empty(t, volumeType)
}

func TestForCreateVolume(t *testing.T) {
	ccf := &DefaultVolumeProvider{sess: nil}

	providerVolume :=
		Volume{
			VolumeID: "16f293bf-test-4bff-816f-e199c0c65db5",
			Name:     String("test volume name"),
			Capacity: nil,
		}
	volume, _ := ccf.CreateVolume(providerVolume)
	assert.Nil(t, volume)
}

func TestForAttachDetachVolume(t *testing.T) {
	ccf := &DefaultVolumeProvider{sess: nil}

	attachRequest :=
		VolumeAttachmentRequest{
			VolumeID: "16f293bf-test-4bff-816f-e199c0c65db5",
		}

	volume, _ := ccf.AttachVolume(attachRequest)
	assert.Nil(t, volume)

	httpResponse, _ := ccf.DetachVolume(attachRequest)
	assert.Nil(t, httpResponse)

	volAttResponse, _ := ccf.WaitForAttachVolume(attachRequest)
	assert.Nil(t, volAttResponse)

	error := ccf.WaitForDetachVolume(attachRequest)
	assert.Nil(t, error)

	volAttachment, _ := ccf.GetVolumeAttachment(attachRequest)
	assert.Nil(t, volAttachment)
}

func TestForExpandVolume(t *testing.T) {
	ccf := &DefaultVolumeProvider{sess: nil}

	expandRequest :=
		ExpandVolumeRequest{
			VolumeID: "16f293bf-test-4bff-816f-e199c0c65db5",
			Name:     new(string),
			Capacity: 0,
		}

	res, _ := ccf.ExpandVolume(expandRequest)
	assert.Equal(t, int64(0), res)
}
func TestForSnapshots(t *testing.T) {
	ccf := &DefaultVolumeProvider{sess: nil}

	snapshot :=
		Snapshot{
			Volume:       Volume{},
			SnapshotID:   "",
			SnapshotTags: map[string]string{},
		}

	providerVolume :=
		Volume{
			VolumeID: "16f293bf-test-4bff-816f-e199c0c65db5",
			Name:     String("test volume name"),
			Capacity: nil,
		}
	volume, _ := ccf.CreateVolumeFromSnapshot(snapshot, nil)
	assert.Nil(t, volume)

	error := ccf.OrderSnapshot(providerVolume)
	assert.Nil(t, error)

	snap, _ := ccf.CreateSnapshot(&providerVolume, nil)
	assert.Nil(t, snap)

	errorDel := ccf.DeleteSnapshot(&snapshot)
	assert.Nil(t, errorDel)

	getSnap, _ := ccf.GetSnapshot(snapshot.SnapshotID)
	assert.Nil(t, getSnap)

	getSnapWithID, _ := ccf.GetSnapshotWithVolumeID(providerVolume.VolumeID, snapshot.SnapshotID)
	assert.Nil(t, getSnapWithID)

	listSnapWithID, _ := ccf.ListSnapshots()
	assert.Nil(t, listSnapWithID)

	listAllSnapWithID, _ := ccf.ListAllSnapshots(providerVolume.VolumeID)
	assert.Nil(t, listAllSnapWithID)
}

func TestForUpdateDeleteGetVolume(t *testing.T) {
	ccf := &DefaultVolumeProvider{sess: nil}

	providerVolume :=
		Volume{
			VolumeID: "16f293bf-test-4bff-816f-e199c0c65db5",
			Name:     String("test volume name"),
			Capacity: nil,
		}

	volume := ccf.UpdateVolume(providerVolume)
	assert.Nil(t, volume)

	delVolume := ccf.DeleteVolume(&providerVolume)
	assert.Nil(t, delVolume)

	getVolume, _ := ccf.GetVolume(providerVolume.VolumeID)
	assert.Nil(t, getVolume)

	getVolumeByName, _ := ccf.GetVolumeByName(*providerVolume.Name)
	assert.Nil(t, getVolumeByName)

	getVolumeByRequestID, _ := ccf.GetVolumeByRequestID("abc1234")
	assert.Nil(t, getVolumeByRequestID)

	getVolumeByRequestIDList, _ := ccf.ListVolumes(50, "1", nil)
	assert.Nil(t, getVolumeByRequestIDList)
}

func TestForAuthorizeVolume(t *testing.T) {
	ccf := &DefaultVolumeProvider{sess: nil}

	volumeAuthorization :=
		VolumeAuthorization{
			Volume:  Volume{},
			Subnets: []string{},
			HostIPs: []string{},
		}
	error := ccf.AuthorizeVolume(volumeAuthorization)
	assert.Nil(t, error)
}

// String returns a pointer to the string value provided
func String(v string) *string {
	return &v
}
