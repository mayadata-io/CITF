/*
Copyright 2018 The OpenEBS Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package k8s

import (
	openebs_v1 "github.com/openebs/CITF/pkg/apis/openebs.io/v1alpha1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GetStoragePoolClaim returns the StoragePoolClaim object for given spcName.
// :return: *openebs_v1.StoragePoolClaim: Pointer to StoragePoolClaim objects.
func (k8s K8S) GetStoragePoolClaim(spcName string) (*openebs_v1.StoragePoolClaim, error) {
	spcClient := k8s.OpenebsClientSet.Openebs().StoragePoolClaims()
	return spcClient.Get(spcName, meta_v1.GetOptions{})
}

// ListStoragePoolClaims returns all the StoragePoolClaim objects
func (k8s K8S) ListStoragePoolClaims() ([]openebs_v1.StoragePoolClaim, error) {
	spcCient := k8s.OpenebsClientSet.Openebs().StoragePoolClaims()
	spcs, err := spcCient.List(meta_v1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return spcs.Items, nil
}

// DeleteStoragePoolClaim deletes a StoragePoolClaim with the given name
func (k8s K8S) DeleteStoragePoolClaim(spcName string) error {
	spcClient := k8s.OpenebsClientSet.Openebs().StoragePoolClaims()
	return spcClient.Delete(spcName, &meta_v1.DeleteOptions{})
}

// GetCStorVolumeReplica returns the CStorVolumeReplica object for given cvrName, namespace.
// :return: *openebs_v1.CStorVolumeReplica: Pointer to CStorVolumeReplica objects.
func (k8s K8S) GetCStorVolumeReplica(cvrName, namespace string) (*openebs_v1.CStorVolumeReplica, error) {
	cvrClient := k8s.OpenebsClientSet.OpenebsV1alpha1().CStorVolumeReplicas(namespace)
	return cvrClient.Get(cvrName, meta_v1.GetOptions{})
}

// ListCStorVolumeReplica returns all the CStorVolumeReplicaList for given namespace
func (k8s K8S) ListCStorVolumeReplica(namespace string) (*openebs_v1.CStorVolumeReplicaList, error) {
	cvrClient := k8s.OpenebsClientSet.OpenebsV1alpha1().CStorVolumeReplicas(namespace)
	return cvrClient.List(meta_v1.ListOptions{})
}

// DeleteCStorVolumeReplica deletes a CStorVolumeReplica with the given cvrName, namespace
func (k8s K8S) DeleteCStorVolumeReplica(cvrName, namespace string) error {
	cvrClient := k8s.OpenebsClientSet.OpenebsV1alpha1().CStorVolumeReplicas(namespace)
	return cvrClient.Delete(cvrName, &meta_v1.DeleteOptions{})
}
