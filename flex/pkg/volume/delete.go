/*
Copyright 2016 The Kubernetes Authors.

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

package volume

import (
	"fmt"
	"os/exec"

	"github.com/golang/glog"
	"github.com/kubernetes-incubator/external-storage/lib/controller"
	"k8s.io/client-go/pkg/api/v1"
)

func (p *flexProvisioner) Delete(volume *v1.PersistentVolume) error {
	glog.Infof("Delete called for volume:", volume.Name)

	provisioned, err := p.provisioned(volume)
	if err != nil {
		return fmt.Errorf("error determining if this provisioner was the one to provision volume %q: %v", volume.Name, err)
	}
	if !provisioned {
		strerr := fmt.Sprintf("this provisioner id %s didn't provision volume %q and so can't delete it; id %s did & can", p.identity, volume.Name, volume.Annotations[annProvisionerId])
		return &controller.IgnoredError{strerr}
	}

	cmd := exec.Command(p.execCommand, "delete")
	output, err := cmd.CombinedOutput()
	if err != nil {
		glog.Errorf("Failed to delete volume %s, output: %s, error: %s", volume, output, err.Error())
		return err
	}
	return nil
}

func (p *flexProvisioner) provisioned(volume *v1.PersistentVolume) (bool, error) {
	provisionerId, ok := volume.Annotations[annProvisionerId]
	if !ok {
		return false, fmt.Errorf("PV doesn't have an annotation %s", annProvisionerId)
	}

	return provisionerId == string(p.identity), nil
}
