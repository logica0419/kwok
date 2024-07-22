/*
Copyright 2023 The Kubernetes Authors.

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

// Package binary_test is a test environment for kwok.
package binary_test

import (
	"flag"
	"os"
	"runtime"
	"testing"

	"sigs.k8s.io/e2e-framework/pkg/env"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/support/kwok"

	"sigs.k8s.io/kwok/pkg/consts"
	"sigs.k8s.io/kwok/pkg/utils/path"
	"sigs.k8s.io/kwok/test/e2e/helper"
)

var (
	runtimeEnv     = consts.RuntimeTypeBinary
	testEnv        env.Environment
	updateTestdata = false
	pwd            = os.Getenv("PWD")
	rootDir        = path.Join(pwd, "../../../..")
	logsDir        = path.Join(rootDir, "logs")
	clusterName    = envconf.RandomName("kwok-e2e-binary", 16)
	namespace      = envconf.RandomName("ns", 16)
	kwokPath       = path.Join(rootDir, "bin", runtime.GOOS, runtime.GOARCH, "kwok"+helper.BinSuffix)
	kwokctlPath    = path.Join(rootDir, "bin", runtime.GOOS, runtime.GOARCH, "kwokctl"+helper.BinSuffix)
	baseArgs       = []string{
		"--kwok-controller-binary=" + kwokPath,
		"--runtime=" + runtimeEnv,
		"--enable-metrics-server",
		"--wait=15m",
	}
)

func init() {
	_ = os.Setenv("KWOK_WORKDIR", path.Join(rootDir, "workdir"))
	flag.BoolVar(&updateTestdata, "update-testdata", false, "update all of testdata")
}

func TestMain(m *testing.M) {
	testEnv = helper.Environment()

	k := kwok.NewProvider().
		WithName(clusterName).
		WithPath(kwokctlPath)
	testEnv.Setup(
		helper.BuildKwokBinary(rootDir),
		helper.BuildKwokctlBinary(rootDir),
		helper.CreateCluster(k, append(baseArgs,
			"--controller-port=10247",
			"--config="+path.Join(rootDir, "test/e2e/port_forward.yaml"),
			"--config="+path.Join(rootDir, "test/e2e/logs.yaml"),
			"--config="+path.Join(rootDir, "test/e2e/attach.yaml"),
			"--config="+path.Join(rootDir, "test/e2e/exec.yaml"),
			"--config="+path.Join(rootDir, "kustomize/metrics/usage/usage-from-annotation.yaml"),
			"--config="+path.Join(rootDir, "kustomize/metrics/resource/metrics-resource.yaml"),
		)...),
		helper.CreateNamespace(namespace),
	)
	testEnv.Finish(
		helper.ExportLogs(k, logsDir),
		helper.DestroyCluster(k),
	)
	os.Exit(testEnv.Run(m))
}
