// Copyright 2020 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Run these tests using:
//   apphosting/runtime/titanium/buildpacks/tools/run-acceptance-tests.sh --runtime=gcpbase
//
package acceptance

import (
	"testing"

	"github.com/GoogleCloudPlatform/buildpacks/pkg/acceptance"
)

func init() {
	acceptance.DefineFlags()
}

func TestAcceptance(t *testing.T) {
	builder, cleanup := acceptance.CreateBuilder(t)
	t.Cleanup(cleanup)

	testCases := []acceptance.Test{
		{
			Name:    "simple Python application",
			App:     "python/simple",
			Env:     []string{"GOOGLE_ENTRYPOINT=gunicorn -b :8080 main:app"},
			MustUse: []string{pythonRuntime, pythonPIP, entrypoint},
		},
		{
			Name: "Python runtime version respected",
			App:  "python/simple",
			Path: "/version?want=3.8.0",
			Env: []string{
				"GOOGLE_ENTRYPOINT=gunicorn -b :8080 main:app",
				"GOOGLE_RUNTIME_VERSION=3.8.0",
			},
			MustUse: []string{pythonRuntime, pythonPIP, entrypoint},
		},
		{
			Name:       "Python selected via GOOGLE_RUNTIME",
			App:        "override",
			Env:        []string{"GOOGLE_RUNTIME=python", "GOOGLE_ENTRYPOINT=gunicorn -b :8080 main:app"},
			MustUse:    []string{pythonRuntime},
			MustNotUse: []string{goRuntime, javaRuntime, nodeRuntime},
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()

			acceptance.TestApp(t, builder, tc)
		})
	}
}

func TestFailures(t *testing.T) {
	builder, cleanup := acceptance.CreateBuilder(t)
	t.Cleanup(cleanup)

	testCases := []acceptance.FailureTest{
		{
			Name:      "bad runtime version",
			App:       "python/simple",
			Env:       []string{"GOOGLE_RUNTIME_VERSION=BAD_NEWS_BEARS", "GOOGLE_ENTRYPOINT=gunicorn -b :8080 main:app"},
			MustMatch: "Runtime version BAD_NEWS_BEARS does not exist",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()

			acceptance.TestBuildFailure(t, builder, tc)
		})
	}
}