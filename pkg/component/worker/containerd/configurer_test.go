// SPDX-FileCopyrightText: 2026 k0s authors
// SPDX-License-Identifier: Apache-2.0

package containerd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/k0sproject/k0s/pkg/apis/k0s/v1beta1"
	workerconfig "github.com/k0sproject/k0s/pkg/component/worker/config"

	"github.com/pelletier/go-toml"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var runtimeBinaryPath = []string{"plugins", "io.containerd.cri.v1.runtime", "containerd", "runtimes", "runc", "options", "BinaryName"}

func TestMarshalContainerdConfig_RuntimeBinary(t *testing.T) {
	t.Run("default keeps the bundled runc", func(t *testing.T) {
		data, err := marshalContainerdConfig("/imports", "pause:1", "")
		require.NoError(t, err)
		tree, err := toml.LoadBytes(data)
		require.NoError(t, err)
		assert.False(t, tree.HasPath(runtimeBinaryPath))
	})

	t.Run("runtime binary sets BinaryName", func(t *testing.T) {
		data, err := marshalContainerdConfig("/imports", "pause:1", "/usr/bin/crun")
		require.NoError(t, err)
		tree, err := toml.LoadBytes(data)
		require.NoError(t, err)
		assert.Equal(t, "/usr/bin/crun", tree.GetPath(runtimeBinaryPath))
		assert.NoError(t, validateConfigTree(tree))
	})
}

func TestSetupConfig_RuntimeBinary(t *testing.T) {
	newComponent := func(confPath string) *Component {
		return &Component{
			RuntimeBinary: "/usr/bin/crun",
			confPath:      confPath,
			importsPath:   t.TempDir(),
			Profile:       &workerconfig.Profile{PauseImage: &v1beta1.ImageSpec{}},
		}
	}

	t.Run("managed config gets BinaryName", func(t *testing.T) {
		confPath := filepath.Join(t.TempDir(), "containerd.toml")
		require.NoError(t, newComponent(confPath).setupConfig())

		tree, err := toml.LoadFile(confPath)
		require.NoError(t, err)
		assert.Equal(t, "/usr/bin/crun", tree.GetPath(runtimeBinaryPath))
	})

	t.Run("user-managed config stays unchanged", func(t *testing.T) {
		confPath := filepath.Join(t.TempDir(), "containerd.toml")
		userConfig := []byte("version = 3\n")
		require.NoError(t, os.WriteFile(confPath, userConfig, 0644))
		require.NoError(t, newComponent(confPath).setupConfig())

		data, err := os.ReadFile(confPath)
		require.NoError(t, err)
		assert.Equal(t, userConfig, data)
	})
}
