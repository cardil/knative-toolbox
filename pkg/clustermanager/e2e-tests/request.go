/*
Copyright 2020 The Knative Authors

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

package e2e_tests

import (
	"fmt"

	clm "knative.dev/toolbox/pkg/clustermanager/e2e-tests/gke"
)

// RequestWrapper is a wrapper of the GKERequest.
type RequestWrapper struct {
	Request clm.GKERequest
	Regions []string
}

func (rw *RequestWrapper) acquire() (*clm.GKECluster, error) {
	gkeClient := clm.GKEClient{}
	clusterOps := gkeClient.Setup(rw.Request)
	gkeOps := clusterOps.(*clm.GKECluster)
	if err := gkeOps.Acquire(); err != nil || gkeOps.Cluster == nil {
		return nil, fmt.Errorf("failed acquiring GKE cluster: %w", err)
	}
	return gkeOps, nil
}
