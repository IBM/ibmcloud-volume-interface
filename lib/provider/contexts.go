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

import "net/http"

// Context represents the volume provider management API for individual account, user ID, etc.
//go:generate counterfeiter -o fakes/context.go --fake-name Context . Context
type Context interface {
	VolumeManager
	VolumeAttachManager
	SnapshotManager
}

// Session is an Context that is notified when it is no longer required
//go:generate counterfeiter -o fakes/session.go --fake-name Session . Session
type Session interface {
	Context

	// GetProviderDisplayName returns the name of the provider that is being used
	GetProviderDisplayName() VolumeProvider

	// Close is called when the Session is nolonger required
	Close()
}

//DefaultSession Implementation
type DefaultSession struct {
	Session
}

var _ Session = &DefaultSession{}

//NewDefaultSession object
func NewDefaultSession() Session {
	return &DefaultSession{}
}

//ProviderName returns provider
func (tes *DefaultSession) ProviderName() VolumeProvider {
	return ""
}

//Type returns the underlying volume type
func (tes *DefaultSession) Type() VolumeType {
	return ""
}

//CreateVolume creates a volume
func (tes *DefaultSession) CreateVolume(VolumeRequest Volume) (*Volume, error) {
	return nil, nil
}

//AttachVolume attaches a volume
func (tes *DefaultSession) AttachVolume(attachRequest VolumeAttachmentRequest) (*VolumeAttachmentResponse, error) {
	return nil, nil
}

//CreateVolumeFromSnapshot creates a volume from snapshot
func (tes *DefaultSession) CreateVolumeFromSnapshot(snapshot Snapshot, tags map[string]string) (*Volume, error) {
	return nil, nil
}

//UpdateVolume the volume
func (tes *DefaultSession) UpdateVolume(Volume) error {
	return nil
}

//DeleteVolume deletes the volume
func (tes *DefaultSession) DeleteVolume(*Volume) error {
	return nil
}

//GetVolume by using ID
func (tes *DefaultSession) GetVolume(id string) (*Volume, error) {
	return nil, nil
}

// GetVolumeByName gets volume by name,
// actually some of providers(like VPC) has the capability to provide volume
// details by usig user provided volume name
func (tes *DefaultSession) GetVolumeByName(name string) (*Volume, error) {
	return nil, nil
}

//ListVolumes Get volume lists by using filters
func (tes *DefaultSession) ListVolumes(limit int, start string, tags map[string]string) (*VolumeList, error) {
	return nil, nil
}

// GetVolumeByRequestID fetch the volume by request ID.
// Request Id is the one that is returned when volume is provsioning request is
// placed with Iaas provider.
func (tes *DefaultSession) GetVolumeByRequestID(requestID string) (*Volume, error) {
	return nil, nil
}

//AuthorizeVolume allows aceess to volume  based on given authorization
func (tes *DefaultSession) AuthorizeVolume(volumeAuthorization VolumeAuthorization) error {
	return nil
}

// DetachVolume  by passing required information in the volume object
func (tes *DefaultSession) DetachVolume(detachRequest VolumeAttachmentRequest) (*http.Response, error) {
	return nil, nil
}

//WaitForAttachVolume waits for the volume to be attached to the host
//Return error if wait is timed out OR there is other error
func (tes *DefaultSession) WaitForAttachVolume(attachRequest VolumeAttachmentRequest) (*VolumeAttachmentResponse, error) {
	return nil, nil
}

//WaitForDetachVolume waits for the volume to be detached from the host
//Return error if wait is timed out OR there is other error
func (tes *DefaultSession) WaitForDetachVolume(detachRequest VolumeAttachmentRequest) error {
	return nil
}

//GetVolumeAttachment retirves the current status of given volume attach request
func (tes *DefaultSession) GetVolumeAttachment(attachRequest VolumeAttachmentRequest) (*VolumeAttachmentResponse, error) {
	return nil, nil
}

//OrderSnapshot orders the snapshot
func (tes *DefaultSession) OrderSnapshot(VolumeRequest Volume) error {
	return nil
}

// CreateSnapshot on the volume
func (tes *DefaultSession) CreateSnapshot(volume *Volume, tags map[string]string) (*Snapshot, error) {
	return nil, nil
}

//DeleteSnapshot deletes the snapshot
func (tes *DefaultSession) DeleteSnapshot(*Snapshot) error {
	return nil
}

//GetSnapshot gets the snapshot
func (tes *DefaultSession) GetSnapshot(snapshotID string) (*Snapshot, error) {
	return nil, nil
}

//GetSnapshotWithVolumeID gets the snapshot with volumeID
func (tes *DefaultSession) GetSnapshotWithVolumeID(volumeID string, snapshotID string) (*Snapshot, error) {
	return nil, nil
}

//ListSnapshots list the snapshots
func (tes *DefaultSession) ListSnapshots() ([]*Snapshot, error) {
	return nil, nil
}

//ListAllSnapshots list all the snapshots
func (tes *DefaultSession) ListAllSnapshots(volumeID string) ([]*Snapshot, error) {
	return nil, nil
}

//GetProviderDisplayName gets provider by displayname
func (tes *DefaultSession) GetProviderDisplayName() VolumeProvider {
	return ""
}

//Close is called when the Session is nolonger required
func (tes *DefaultSession) Close() {
}
