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

// Implements /bin/build for ruby/rails buildpack.
package main

import (
	"path/filepath"

	gcp "github.com/GoogleCloudPlatform/buildpacks/pkg/gcpbuildpack"
)

func main() {
	gcp.Main(detectFn, buildFn)
}

func detectFn(ctx *gcp.Context) error {
	if !ctx.FileExists("bin", "rails") {
		ctx.OptOut("bin/rails not found.")
	}
	if !needsRailsAssetPrecompile(ctx) {
		ctx.OptOut("Rails assets do not need precompilation.")
	}
	return nil
}

func needsRailsAssetPrecompile(ctx *gcp.Context) bool {
	if !ctx.FileExists("app", "assets") {
		return false
	}

	if ctx.FileExists("public", "assets", "manifest.yml") {
		return false
	}

	matches := ctx.Glob(filepath.Join(ctx.ApplicationRoot(), "public/assets/manifest-*.json"))
	if matches != nil {
		return false
	}

	matches = ctx.Glob(filepath.Join(ctx.ApplicationRoot(), "public/assets/.sprockets-manifest-*.json"))
	if matches != nil {
		return false
	}

	return true
}

func buildFn(ctx *gcp.Context) error {
	ctx.Logf("Running Rails asset precompilation")

	// It is common practise in Ruby asset precompilation to ignore non-zero exit codes.
	params := gcp.ExecParams{
		Cmd: []string{"bundle", "exec", "bin/rails", "assets:precompile"},
		Env: []string{"RAILS_ENV=production"},
	}
	result, err := ctx.ExecWithErrWithParams(params)
	if err != nil && result != nil && result.ExitCode != 0 {
		ctx.Logf("WARNING: Asset precompilation returned non-zero exit code %d. Ignoring.", result.ExitCode)
	} else if err != nil && result != nil {
		ctx.Exit(1, gcp.UserErrorf(result.Combined))
	} else if err != nil {
		ctx.Exit(1, gcp.InternalErrorf("asset precompilation failed: %v", err))
	}

	return nil
}