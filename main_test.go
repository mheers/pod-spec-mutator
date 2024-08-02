package main

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
)

const patchJSONTest = `{"hostAliases":[{"ip":"192.168.1.100","hostnames":["foo.local"]}]}`

func TestUnmarshalPathJSON(t *testing.T) {
	var patchTemplate corev1.PodSpec
	err := json.Unmarshal([]byte(patchJSONTest), &patchTemplate)
	require.NoError(t, err)
}

func TestCreatePatch(t *testing.T) {
	var patchTemplate corev1.PodSpec
	err := json.Unmarshal([]byte(patchJSONTest), &patchTemplate)
	require.NoError(t, err)

	defaultPod := &corev1.Pod{
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  "nginx",
					Image: "nginx:latest",
				},
			},
			HostAliases: []corev1.HostAlias{
				{
					IP:        "10.1.70.3",
					Hostnames: []string{"demo"},
				},
			},
		},
	}

	newPod, err := applyPatch(defaultPod, []byte(patchJSONTest))
	require.NoError(t, err)
	require.NotEmpty(t, defaultPod.Spec.Containers)
	require.NotEmpty(t, newPod.Spec.Containers)
	require.Len(t, newPod.Spec.HostAliases, 2)

	patch, err := createPatch(*defaultPod, *newPod)
	require.NoError(t, err)
	require.NotEmpty(t, patch)
}
