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

func TestGetProviderType(t *testing.T) {
	ccf := &DefaultVolumeProvider{sess: nil}
	assert.Empty(t, ccf.Type())
}

func TestGetProviderName(t *testing.T) {
	ccf := &DefaultVolumeProvider{sess: nil}

	assert.Empty(t, ccf.ProviderName())
}

func TestGetProviderDisplayName(t *testing.T) {
	ccf := &DefaultVolumeProvider{sess: nil}

	assert.Empty(t, ccf.GetProviderDisplayName())
}

func TestCreateVolume(t *testing.T) {
	ccf := &DefaultVolumeProvider{sess: nil}

	volume, _ := ccf.CreateVolume(Volume{})
	assert.Nil(t, volume)
}

func TestDetachVolume(t *testing.T) {
	ccf := &DefaultVolumeProvider{sess: nil}

	volume, _ := ccf.DetachVolume(VolumeAttachmentRequest{})
	assert.Nil(t, volume)
	if volume != nil {
		volume.Body.Close() // #nosec G104
	}
}

func TestWaitForAttachVolume(t *testing.T) {
	ccf := &DefaultVolumeProvider{sess: nil}

	volume, _ := ccf.WaitForAttachVolume(VolumeAttachmentRequest{})
	assert.Nil(t, volume)
}
func TestAttachVolume(t *testing.T) {
	ccf := &DefaultVolumeProvider{sess: nil}

	volume, _ := ccf.AttachVolume(VolumeAttachmentRequest{})
	assert.Nil(t, volume)
}

func TestWaitForDetachVolume(t *testing.T) {
	ccf := &DefaultVolumeProvider{sess: nil}

	assert.Nil(t, ccf.WaitForDetachVolume(VolumeAttachmentRequest{}))
}

func TestGetVolumeAttachment(t *testing.T) {
	ccf := &DefaultVolumeProvider{sess: nil}

	volAttachment, _ := ccf.GetVolumeAttachment(VolumeAttachmentRequest{})
	assert.Nil(t, volAttachment)
}
func TestExpandVolume(t *testing.T) {
	ccf := &DefaultVolumeProvider{sess: nil}

	res, _ := ccf.ExpandVolume(ExpandVolumeRequest{})
	assert.Equal(t, int64(0), res)
}

func TestCreateVolumeFromSnapshot(t *testing.T) {
	ccf := &DefaultVolumeProvider{sess: nil}

	volume, _ := ccf.CreateVolumeFromSnapshot(Snapshot{}, nil)
	assert.Nil(t, volume)
}

func TestOrderSnapshot(t *testing.T) {
	ccf := &DefaultVolumeProvider{sess: nil}

	assert.Nil(t, ccf.OrderSnapshot(Volume{}))
}

func TestCreateSnapshot(t *testing.T) {
	ccf := &DefaultVolumeProvider{sess: nil}

	snapshot, _ := ccf.CreateSnapshot("vol-id", SnapshotParameters{})
	assert.Nil(t, snapshot)
}

func TestDeleteSnapshot(t *testing.T) {
	ccf := &DefaultVolumeProvider{sess: nil}

	assert.Nil(t, ccf.DeleteSnapshot(&Snapshot{}))
}

func TestGetSnapshot(t *testing.T) {
	ccf := &DefaultVolumeProvider{sess: nil}

	getSnap, _ := ccf.GetSnapshot("snap-id")
	assert.Nil(t, getSnap)
}

func TestListSnapshots(t *testing.T) {
	ccf := &DefaultVolumeProvider{sess: nil}

	listSnapWithID, _ := ccf.ListSnapshots(50, "1", nil)
	assert.Nil(t, listSnapWithID)
}

func TestUpdateVolume(t *testing.T) {
	ccf := &DefaultVolumeProvider{sess: nil}

	assert.Nil(t, ccf.UpdateVolume(Volume{}))
}

func TestDeleteVolume(t *testing.T) {
	ccf := &DefaultVolumeProvider{sess: nil}

	assert.Nil(t, ccf.DeleteVolume(&Volume{}))
}

func TestGetVolume(t *testing.T) {
	ccf := &DefaultVolumeProvider{sess: nil}

	getVolume, _ := ccf.GetVolume("VolumeID")
	assert.Nil(t, getVolume)
}

func TestGetVolumeByName(t *testing.T) {
	ccf := &DefaultVolumeProvider{sess: nil}

	getVolume, _ := ccf.GetVolumeByName("VolumeName")
	assert.Nil(t, getVolume)
}

func TestGetVolumeProfileByName(t *testing.T) {
	ccf := &DefaultVolumeProvider{sess: nil}

	getProfile, err := ccf.GetVolumeProfileByName("VolumeProfileName")
	assert.Nil(t, err)
	assert.Nil(t, getProfile)
}

func TestGetVolumeByRequestID(t *testing.T) {
	ccf := &DefaultVolumeProvider{sess: nil}

	getVolume, _ := ccf.GetVolumeByRequestID("abc1234")
	assert.Nil(t, getVolume)
}

func TestListVolumes(t *testing.T) {
	ccf := &DefaultVolumeProvider{sess: nil}

	getVolumeByRequestIDList, _ := ccf.ListVolumes(50, "1", nil)
	assert.Nil(t, getVolumeByRequestIDList)
}

func TestAuthorizeVolume(t *testing.T) {
	ccf := &DefaultVolumeProvider{sess: nil}

	assert.Nil(t, ccf.AuthorizeVolume(VolumeAuthorization{}))
}

func TestCreateVolumeAccessPoint(t *testing.T) {
	ccf := &DefaultVolumeProvider{sess: nil}

	accessPointResponse, _ := ccf.CreateVolumeAccessPoint(VolumeAccessPointRequest{})
	assert.Nil(t, accessPointResponse)
}

func TestDeleteVolumeAccessPoint(t *testing.T) {
	ccf := &DefaultVolumeProvider{sess: nil}

	accessPointResponse, _ := ccf.DeleteVolumeAccessPoint(VolumeAccessPointRequest{})
	assert.Nil(t, accessPointResponse)
	if accessPointResponse != nil {
		accessPointResponse.Body.Close() // #nosec G104
	}
}

func TestWaitForCreateVolumeAccessPoint(t *testing.T) {
	ccf := &DefaultVolumeProvider{sess: nil}

	accessPointResponse, _ := ccf.WaitForCreateVolumeAccessPoint(VolumeAccessPointRequest{})
	assert.Nil(t, accessPointResponse)
}

func TestWaitForDeleteVolumeAccessPoint(t *testing.T) {
	ccf := &DefaultVolumeProvider{sess: nil}

	err := ccf.WaitForDeleteVolumeAccessPoint(VolumeAccessPointRequest{})
	assert.Nil(t, err)
}

func TestGetVolumeAccessPoint(t *testing.T) {
	ccf := &DefaultVolumeProvider{sess: nil}

	accessPointResponse, _ := ccf.GetVolumeAccessPoint(VolumeAccessPointRequest{})
	assert.Nil(t, accessPointResponse)
}

func TestGetSubnetForVolumeAccessPoint(t *testing.T) {
	ccf := &DefaultVolumeProvider{sess: nil}

	subnetID, _ := ccf.GetSubnetForVolumeAccessPoint(SubnetRequest{})
	assert.Equal(t, subnetID, "")
}

func TestGetSecurityGroupForVolumeAccessPoint(t *testing.T) {
	ccf := &DefaultVolumeProvider{sess: nil}

	securityGroupID, _ := ccf.GetSecurityGroupForVolumeAccessPoint(SecurityGroupRequest{})
	assert.Equal(t, securityGroupID, "")
}
