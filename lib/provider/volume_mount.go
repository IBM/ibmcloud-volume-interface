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
	"net/http"
	"time"
)

// VolumeMountManager ...
type VolumeMountManager interface {
	//Mount method mount a volume/fileset to a server
	CreateVolumeMount(mountRequest VolumeMountRequest) (*VolumeMountResponse, error)

	//DeleteVolumeMount method delete a volume/fileset target
	DeleteVolumeMount(deleteMountRequest VolumeMountRequest) (*http.Response, error)

	//WaitForCreateMountVolume waits for the volume mount to be created
	//Return error if wait is timed out OR there is other error
	WaitForCreateVolumeMount(mountRequest VolumeMountRequest) (*VolumeMountResponse, error)

	//WaitForDeleteMountVolume waits for the volume mount to be deleted
	//Return error if wait is timed out OR there is other error
	WaitForDeleteVolumeMount(deleteMountRequest VolumeMountRequest) error

	//GetVolumeMount retrieves the current status of given volume mount request
	GetVolumeMount(mountRequest VolumeMountRequest) (*VolumeMountResponse, error)
}

// VolumeMountRequest  used for both mount and unmount operation
type VolumeMountRequest struct {

	// Volume provider type i.e  Endurance or Performance or any other name
	ProviderType VolumeProviderType `json:"providerType,omitempty"`

	//Target name for the mount
	TargetName string `json:"name,omitempty"`

	// Volume to create the mount for
	VolumeID string `json:"volumeID"`

	//TargetID to search
	TargetID string `json:"targetID,omitempty"`

	//Subnet to create mount for
	SubnetID string `json:"subnet_id,omitempty"`

	//VPC to create mount for
	VPCID string `json:"vpc_id,omitempty"`
}

// VolumeMountResponse used for both mount and unmount operation
type VolumeMountResponse struct {
	VolumeID  string     `json:"volumeID"`
	TargetID  string     `json:"targetID"`
	Status    string     `json:"status"`
	Server    string     `json:"server"`
	MountPath string     `json:"mount_path"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
}
