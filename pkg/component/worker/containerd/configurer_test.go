// SPDX-FileCopyrightText: 2026 k0s authors
// SPDX-License-Identifier: Apache-2.0

package containerd

import (
	"testing"

	"github.com/pelletier/go-toml"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMarshalContainerdConfig_RuntimeBinary(t *testing.T) {
	binaryPath := []string{"plugins", "io.containerd.cri.v1.runtime", "containerd", "runtimes", "runc", "options", "BinaryName"}

	data, err := marshalContainerdConfig("/imports", "pause:1", "")
	require.NoError(t, err)
	tree, err := toml.LoadBytes(data)
	require.NoError(t, err)
	assert.False(t, tree.HasPath(binaryPath), "bundled runc is the default")

	data, err = marshalContainerdConfig("/imports", "pause:1", "/usr/bin/crun")
	require.NoError(t, err)
	tree, err = toml.LoadBytes(data)
	require.NoError(t, err)
	assert.Equal(t, "/usr/bin/crun", tree.GetPath(binaryPath))
	assert.NoError(t, validateConfigTree(tree))
}
