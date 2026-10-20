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
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestModifyVolume(t *testing.T) {
	testCases := []struct {
		name    string
		request ModifyVolumeRequest
	}{
		{
			name:    "empty request (nil session) returns nil, nil",
			request: ModifyVolumeRequest{},
		},
		{
			name:    "with VolumeID only returns nil, nil",
			request: ModifyVolumeRequest{VolumeID: "vol-id"},
		},
		{
			name:    "with IOPS set returns nil, nil",
			request: ModifyVolumeRequest{VolumeID: "vol-id", VolumePerformance: VolumePerformance{Iops: 3000}},
		},
		{
			name:    "with Bandwidth set returns nil, nil",
			request: ModifyVolumeRequest{VolumeID: "vol-id", VolumePerformance: VolumePerformance{Bandwidth: 5000}},
		},
		{
			name:    "with both IOPS and Bandwidth set returns nil, nil",
			request: ModifyVolumeRequest{VolumeID: "vol-id", VolumePerformance: VolumePerformance{Iops: 3000, Bandwidth: 5000}},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ccf := &DefaultVolumeProvider{sess: nil}
			resp, err := ccf.ModifyVolume(tc.request)
			assert.NoError(t, err)
			assert.Nil(t, resp)
		})
	}
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

// TestModifyVolumeRequest_VACFields verifies that ModifyVolumeRequest carries the
// VAC-relevant fields (VolumeID, Iops, Bandwidth) and round-trips through JSON.
func TestModifyVolumeRequest_VACFields(t *testing.T) {
	testCases := []struct {
		name     string
		req      ModifyVolumeRequest
		wantIops int64
		wantBW   int32
	}{
		{name: "zero-value request has no iops or bandwidth", req: ModifyVolumeRequest{}},
		{name: "dp2 VAC: IOPS only", req: ModifyVolumeRequest{VolumeID: "vol-dp2", VolumePerformance: VolumePerformance{Iops: 5000}}, wantIops: 5000},
		{name: "rfs VAC: bandwidth only", req: ModifyVolumeRequest{VolumeID: "vol-rfs", VolumePerformance: VolumePerformance{Bandwidth: 3000}}, wantBW: 3000},
		{name: "common VAC: both iops and bandwidth", req: ModifyVolumeRequest{VolumeID: "vol", VolumePerformance: VolumePerformance{Iops: 4000, Bandwidth: 2000}}, wantIops: 4000, wantBW: 2000},
		{name: "name is preserved", req: ModifyVolumeRequest{VolumeID: "vol", Name: "new-name", VolumePerformance: VolumePerformance{Iops: 1000}}, wantIops: 1000},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.wantIops, tc.req.Iops)
			assert.Equal(t, tc.wantBW, tc.req.Bandwidth)

			b, err := json.Marshal(tc.req)
			require.NoError(t, err)
			var decoded ModifyVolumeRequest
			require.NoError(t, json.Unmarshal(b, &decoded))
			assert.Equal(t, tc.req.VolumeID, decoded.VolumeID)
			assert.Equal(t, tc.req.Iops, decoded.Iops)
			assert.Equal(t, tc.req.Bandwidth, decoded.Bandwidth)
			assert.Equal(t, tc.req.Name, decoded.Name)
		})
	}
}

// TestModifyVolumeResponse_VACFields verifies ModifyVolumeResponse carries the
// updated attributes returned by a VAC modify call and round-trips through JSON.
func TestModifyVolumeResponse_VACFields(t *testing.T) {
	testCases := []struct {
		name     string
		resp     ModifyVolumeResponse
		wantIops int64
		wantBW   int32
	}{
		{name: "response carries volume ID, href and CRN", resp: ModifyVolumeResponse{VolumeID: "vol-id", Href: "https://example/vol-id", CRN: "crn:v1:vol-id", VolumePerformance: VolumePerformance{Iops: 6000}}, wantIops: 6000},
		{name: "dp2 response: IOPS updated", resp: ModifyVolumeResponse{VolumePerformance: VolumePerformance{Iops: 6000}}, wantIops: 6000},
		{name: "rfs response: bandwidth updated", resp: ModifyVolumeResponse{VolumePerformance: VolumePerformance{Bandwidth: 4000}}, wantBW: 4000},
		{name: "empty response has zero values", resp: ModifyVolumeResponse{}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.wantIops, tc.resp.Iops)
			assert.Equal(t, tc.wantBW, tc.resp.Bandwidth)

			b, err := json.Marshal(tc.resp)
			require.NoError(t, err)
			var decoded ModifyVolumeResponse
			require.NoError(t, json.Unmarshal(b, &decoded))
			assert.Equal(t, tc.resp, decoded)
		})
	}
}
