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

package artifact

import (
	"fmt"
	"log"
	"path"
	"strings"

	"github.com/google/ko/pkg/build"
	"github.com/google/ko/pkg/commands"
	"github.com/google/ko/pkg/commands/options"
	"github.com/google/ko/pkg/publish"
	"knative.dev/toolbox/magetasks/config"
	artifactimage "knative.dev/toolbox/magetasks/pkg/artifact/image"
	"knative.dev/toolbox/magetasks/pkg/output/color"
	"knative.dev/toolbox/magetasks/pkg/version"
)

const koPublishResult = "ko.publish.result"

// KoPublisherConfigurator is used to configure the publish options for KO.
type KoPublisherConfigurator func(*options.PublishOptions) error

// KoPublisher publishes images with Google's KO.
type KoPublisher struct {
	Configurators []KoPublisherConfigurator
}

func (kp KoPublisher) Accepts(artifact config.Artifact) bool {
	_, ok := artifact.(Image)
	return ok
}

func (kp KoPublisher) Publish(artifact config.Artifact, notifier config.Notifier) config.Result {
	image, ok := artifact.(Image)
	if !ok {
		return config.Result{Error: ErrInvalidArtifact}
	}
	buildResult, ok := config.Actual().Context.Value(BuildKey(image)).(config.Result)
	if !ok || buildResult.Failed() {
		return config.Result{Error: fmt.Errorf(
			"%w: can't find successful KO build result", ErrInvalidArtifact)}
	}
	result, ok := buildResult.Info[koBuildResult].(build.Result)
	if !ok {
		return config.Result{Error: fmt.Errorf(
			"%w: can't find successful KO build result", ErrInvalidArtifact)}
	}
	po, err := kp.publishOptions()
	if err != nil {
		return resultErrKoFailed(err)
	}
	ctx := config.Actual().Context
	publisher, err := commands.NewPublisher(po)
	if err != nil {
		return resultErrKoFailed(err)
	}
	ref, err := publisher.Publish(ctx, result, image.GetName())
	if err != nil {
		return resultErrKoFailed(err)
	}
	notifier.Notify(fmt.Sprintf("pushed image: %s", color.Blue(ref)))
	notifier.Notify(fmt.Sprintf("image tags: %s", color.Blue(fmt.Sprintf("%+q", po.Tags))))
	return config.Result{Info: map[string]interface{}{
		koPublishResult: ref,
	}}
}

func (kp KoPublisher) publishOptions() (*options.PublishOptions, error) {
	opts := &options.PublishOptions{
		DockerRepo: artifactimage.BaseName(),
		Push:       true,
		ImageNamer: customSeparatorImageNamer{
			artifactimage.BaseNameSeparator(),
		}.name,
	}
	if ver := config.Actual().Version; ver != nil {
		if err := resolveTags(opts, ver.Resolver); err != nil {
			return nil, err
		}
	}
	for _, configurator := range kp.Configurators {
		if err := configurator(opts); err != nil {
			return nil, err
		}
	}
	return opts, nil
}

func resolveTags(opts *options.PublishOptions, resolver version.Resolver) error {
	tags, err := artifactimage.Tags(resolver)
	if err != nil {
		return err
	}
	opts.Tags = tags
	return nil
}

func closePublisher(publisher publish.Interface) {
	if err := publisher.Close(); err != nil {
		log.Fatal(err)
	}
}

type customSeparatorImageNamer struct {
	sep string
}

func (n customSeparatorImageNamer) name(base, importpath string) string {
	return n.join(base, path.Base(importpath))
}

func (n customSeparatorImageNamer) join(paths ...string) string {
	for i, e := range paths {
		if e != "" {
			return path.Clean(strings.Join(paths[i:], n.sep))
		}
	}
	return ""
}
