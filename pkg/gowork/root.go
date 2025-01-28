/*
Copyright 2024 The Knative Authors

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

package gowork

import (
	"errors"
	"path"
)

// Root returns the root directory of the repo, either being a sole go module,
// or a go workspace project.
func Root(filesystem FileSystem, env Environment) (string, error) {
	gowork, err := findWorkfile(filesystem, env)
	if gowork == "" || errors.Is(err, ErrInvalidGowork) {
		curr, cerr := Current(filesystem, env)
		if cerr != nil {
			return "", cerr
		}
		return curr.Path, nil
	}
	return path.Dir(gowork), nil
}
