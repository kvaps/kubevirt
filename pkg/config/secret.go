/*
 * This file is part of the KubeVirt project
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
 *
 * Copyright 2018 Red Hat, Inc.
 *
 */

package config

import (
	"path/filepath"

	v1 "kubevirt.io/client-go/api/v1"
	ephemeraldiskutils "kubevirt.io/kubevirt/pkg/ephemeral-disk-utils"
)

// GetSecretSourcePath returns a path to Secret mounted on a pod
func GetSecretSourcePath(volumeName string) string {
	return filepath.Join(SecretSourceDir, volumeName)
}

// GetSecretDiskPath returns a path to Secret iso image created based on volume name
func GetSecretDiskPath(volumeName string) string {
	return filepath.Join(SecretDisksDir, volumeName+".iso")
}

// CreateSecretDisks creates Secret iso disks which are attached to vmis
func CreateSecretDisks(vmi *v1.VirtualMachineInstance, emptyIso bool) error {
	for _, volume := range vmi.Spec.Volumes {
		if volume.Secret != nil {

			var filesPath []string
			filesPath, err := getFilesLayout(GetSecretSourcePath(volume.Name))
			if err != nil {
				return err
			}

			disk := GetSecretDiskPath(volume.Name)
			vmiIsoSize, err := findIsoSize(vmi, &volume, emptyIso)
			if err != nil {
				return err
			}
			if err := createIsoConfigImage(disk, volume.Secret.VolumeLabel, filesPath, vmiIsoSize); err != nil {
				return err
			}

			if err := ephemeraldiskutils.DefaultOwnershipManager.SetFileOwnership(disk); err != nil {
				return err
			}
		}
	}

	return nil
}
