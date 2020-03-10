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

package nodejs

import (
	"fmt"

	gcp "github.com/GoogleCloudPlatform/buildpacks/pkg/gcpbuildpack"
	"github.com/blang/semver"
)

const (
	// PackageLock is the name of the npm lock file.
	PackageLock = "package-lock.json"
)

var (
	// minCIVersion is the minimium npm version where `ci` does not delete the node_modules directory.
	minCIVersion = semver.MustParse("6.12.1")
)

// EnsurePackageLock generates a package-lock.json if necessary.
func EnsurePackageLock(ctx *gcp.Context) {
	if !ctx.FileExists(PackageLock) {
		ctx.Logf("Generating %s.", PackageLock)
		ctx.Warnf("*** Improve build performance by generating and committing %s.", PackageLock)
		ctx.ExecUser([]string{"npm", "install", "--package-lock-only", "--quiet"})
	}
}

// NPMInstallCommand returns the correct install commmand based on the version of Node.js.
func NPMInstallCommand(ctx *gcp.Context) (string, error) {
	// Use npm install instead of npm ci for Node.js 10 as it did not launch with npm v6.
	raw := ctx.Exec([]string{"npm", "--version"}).Stdout
	version, err := semver.Parse(raw)
	if err != nil {
		return "", fmt.Errorf("parsing npm version %q: %v", raw, err)
	}
	if version.GTE(minCIVersion) {
		return "ci", nil
	}
	return "install", nil
}