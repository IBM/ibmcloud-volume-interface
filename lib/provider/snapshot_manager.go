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

// SnapshotManager ...
type SnapshotManager interface {
	// Create the snapshot on the volume
	CreateSnapshot(snapshotRequest SnapshotRequest) (*Snapshot, error)

	// Delete the snapshot
	DeleteSnapshot(*Snapshot) error

	// Get the snapshot
	GetSnapshot(snapshotID string) (*Snapshot, error)

	// Get the snapshot By name
	GetSnapshotByName(snapshotName string) (*Snapshot, error)

	// Get the snapshot with volume ID
	GetSnapshotWithVolumeID(volumeID string, snapshotID string) (*Snapshot, error)

	// Snapshot list by using tags
	ListSnapshots(limit int, start string, tags map[string]string) (*SnapshotList, error)

	//List all the  snapshots for a given volume
	ListAllSnapshots(volumeID string) ([]*Snapshot, error)
}
