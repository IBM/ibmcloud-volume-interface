/**
 * Copyright 2021 IBM Corp.
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

func TestForGetProviderType(t *testing.T) {
	ccf := &DefaultVolumeProvider{sess: nil}
	assert.Empty(t, ccf.Type())
}

func TestForGetProviderName(t *testing.T) {
	ccf := &DefaultVolumeProvider{sess: nil}

	assert.Empty(t, ccf.ProviderName())
}

func TestForGetProviderDisplayName(t *testing.T) {
	ccf := &DefaultVolumeProvider{sess: nil}

	assert.Empty(t, ccf.GetProviderDisplayName())
}

func TestForCreateVolume(t *testing.T) {
	ccf := &DefaultVolumeProvider{sess: nil}

	volume, _ := ccf.CreateVolume(Volume{})
	assert.Nil(t, volume)
}

func TestForDetachVolume(t *testing.T) {
	ccf := &DefaultVolumeProvider{sess: nil}

	volume, _ := ccf.DetachVolume(VolumeAttachmentRequest{})
	assert.Nil(t, volume)
}

func TestForWaitForAttachVolume(t *testing.T) {
	ccf := &DefaultVolumeProvider{sess: nil}

	volume, _ := ccf.WaitForAttachVolume(VolumeAttachmentRequest{})
	assert.Nil(t, volume)
}
func TestForAttachVolume(t *testing.T) {
	ccf := &DefaultVolumeProvider{sess: nil}

	volume, _ := ccf.AttachVolume(VolumeAttachmentRequest{})
	assert.Nil(t, volume)
}

func TestForWaitForDetachVolume(t *testing.T) {
	ccf := &DefaultVolumeProvider{sess: nil}

	assert.Nil(t, ccf.WaitForDetachVolume(VolumeAttachmentRequest{}))
}

func TestForGetVolumeAttachment(t *testing.T) {
	ccf := &DefaultVolumeProvider{sess: nil}

	volAttachment, _ := ccf.GetVolumeAttachment(VolumeAttachmentRequest{})
	assert.Nil(t, volAttachment)
}
func TestForExpandVolume(t *testing.T) {
	ccf := &DefaultVolumeProvider{sess: nil}

	res, _ := ccf.ExpandVolume(ExpandVolumeRequest{})
	assert.Equal(t, int64(0), res)
}

func TestForCreateVolumeFromSnapshot(t *testing.T) {
	ccf := &DefaultVolumeProvider{sess: nil}

	volume, _ := ccf.CreateVolumeFromSnapshot(Snapshot{}, nil)
	assert.Nil(t, volume)
}

func TestForOrderSnapshot(t *testing.T) {
	ccf := &DefaultVolumeProvider{sess: nil}

	assert.Nil(t, ccf.OrderSnapshot(Volume{}))
}

func TestForCreateSnapshot(t *testing.T) {
	ccf := &DefaultVolumeProvider{sess: nil}

	volume, _ := ccf.CreateSnapshot(&Volume{}, nil)
	assert.Nil(t, volume)
}

func TestForDeleteSnapshot(t *testing.T) {
	ccf := &DefaultVolumeProvider{sess: nil}

	assert.Nil(t, ccf.DeleteSnapshot(&Snapshot{}))
}

func TestForGetSnapshot(t *testing.T) {
	ccf := &DefaultVolumeProvider{sess: nil}

	getSnap, _ := ccf.GetSnapshot("snap-id")
	assert.Nil(t, getSnap)
}

func TestForGetSnapshotWithVolumeID(t *testing.T) {
	ccf := &DefaultVolumeProvider{sess: nil}

	getSnapWithID, _ := ccf.GetSnapshotWithVolumeID("VolumeID", "SnapshotID")
	assert.Nil(t, getSnapWithID)
}

func TestForListSnapshots(t *testing.T) {
	ccf := &DefaultVolumeProvider{sess: nil}

	listSnapWithID, _ := ccf.ListSnapshots()
	assert.Nil(t, listSnapWithID)
}

func TestForListAllSnapshots(t *testing.T) {
	ccf := &DefaultVolumeProvider{sess: nil}

	listAllSnapWithID, _ := ccf.ListAllSnapshots("VolumeID")
	assert.Nil(t, listAllSnapWithID)
}

func TestForUpdateVolume(t *testing.T) {
	ccf := &DefaultVolumeProvider{sess: nil}

	assert.Nil(t, ccf.UpdateVolume(Volume{}))
}

func TestForDeleteVolume(t *testing.T) {
	ccf := &DefaultVolumeProvider{sess: nil}

	assert.Nil(t, ccf.DeleteVolume(&Volume{}))
}

func TestForGetVolume(t *testing.T) {
	ccf := &DefaultVolumeProvider{sess: nil}

	getVolume, _ := ccf.GetVolume("VolumeID")
	assert.Nil(t, getVolume)
}

func TestForGetVolumeByName(t *testing.T) {
	ccf := &DefaultVolumeProvider{sess: nil}

	getVolume, _ := ccf.GetVolumeByName("VolumeName")
	assert.Nil(t, getVolume)
}
func TestForGetVolumeByRequestID(t *testing.T) {
	ccf := &DefaultVolumeProvider{sess: nil}

	getVolume, _ := ccf.GetVolumeByRequestID("abc1234")
	assert.Nil(t, getVolume)
}

func TestForListVolumes(t *testing.T) {
	ccf := &DefaultVolumeProvider{sess: nil}

	getVolumeByRequestIDList, _ := ccf.ListVolumes(50, "1", nil)
	assert.Nil(t, getVolumeByRequestIDList)
}

func TestForAuthorizeVolume(t *testing.T) {
	ccf := &DefaultVolumeProvider{sess: nil}

	assert.Nil(t, ccf.AuthorizeVolume(VolumeAuthorization{}))
}
