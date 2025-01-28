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

package docs

import (
	"go/build"
	"log"
	"os"
	"path/filepath"

	"cuelang.org/go/pkg/strings"
	"knative.dev/toolbox/pkg/gowork"
)

func findGoModForPackage(pkg string) (*gowork.Module, error) {
	sys := gowork.RealSystem{}
	mods, err := gowork.List(sys, sys)
	if err != nil {
		return nil, err
	}
	var found *gowork.Module
	for _, mod := range mods {
		if strings.HasPrefix(pkg, mod.Name) {
			// find the most specific module for a package
			if found != nil {
				if len(mod.Name) > len(found.Name) {
					found = &mod
				}
			} else {
				found = &mod
			}
		}
	}
	return found, nil
}

func findRootGoMod() *gowork.Module {
	sys := gowork.RealSystem{}
	curr, err := gowork.Current(sys, sys)
	if err != nil {
		log.Fatal(err)
	}
	mods, err := gowork.List(sys, sys)
	if err != nil {
		log.Fatal(err)
	}
	var found *gowork.Module
	for _, mod := range mods {
		if strings.HasPrefix(curr.Name, mod.Name) {
			// find the shortest matching module -> the root module
			if found == nil {
				found = &mod
			} else {
				if len(mod.Name) < len(found.Name) {
					found = &mod
				}
			}
		}
	}
	return found
}

func gomodcache() string {
	// See: https://github.com/golang/go/blob/release-branch.go1.22/src/cmd/go/internal/cfg/cfg.go#L408
	return envOr("GOMODCACHE", gopathDir("pkg/mod"))
}

func gopathDir(rel string) string {
	list := filepath.SplitList(gopath())
	if len(list) == 0 || list[0] == "" {
		return ""
	}
	return filepath.Join(list[0], rel)
}

func gopath() string {
	gp := os.Getenv("GOPATH")
	if gp == "" {
		gp = build.Default.GOPATH
	}
	return gp
}

func envOr(key string, def string) string {
	val := os.Getenv(key)
	if val == "" {
		val = def
	}
	return val
}
