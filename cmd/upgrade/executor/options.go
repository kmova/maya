/*
Copyright 2019 The OpenEBS Authors.

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

package executor

import (
	errors "github.com/openebs/maya/pkg/errors/v1alpha1"
	"github.com/spf13/cobra"
	"strings"
)

// UpgradeOptions stores information required for upgrade
type UpgradeOptions struct {
	resourceName      string
	resourceKind      string
	fromVersion       string
	toVersion         string
	namespace         string
	imageURLPrefix    string
	toVersionImageTag string
}

var (
	options = &UpgradeOptions{
		namespace:      "openebs",
		imageURLPrefix: "quay.io/openebs/",
	}
)

// RunPreFlightChecks will ensure the sanity of the common upgrade options
func (u *UpgradeOptions) RunPreFlightChecks(cmd *cobra.Command) error {
	if len(strings.TrimSpace(u.namespace)) != 0 {
		return errors.Errorf("Cannot execute upgrade job: namespace is missing")
	}

	if len(strings.TrimSpace(u.fromVersion)) == 0 {
		return errors.Errorf("Cannot execute upgrade job: from-version is missing")
	}

	if len(strings.TrimSpace(u.toVersion)) == 0 {
		return errors.Errorf("Cannot execute upgrade job: to-version is missing")
	}
	return nil
}