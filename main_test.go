package main

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const patchJSONTestUpperFirst = `
{
"HostAliases":[{"Ip":"192.168.1.100","Hostnames":["foo.local"]}],
"Containers":[{"Name":"postgres","ImagePullPolicy":"Never","Image":""}]
}
`

const patchJSONTest = `{"containers":[{"imagePullPolicy":"Never","name":"postgres"}],"hostAliases":[{"hostnames":["foo.local"],"ip":"192.168.1.100"}]}`

const patchMultiple = `
{
	"patches": [
		{
			"podNameRegex": "demo",
			"pod": {
				"spec": {
					"containers": [
						{
							"name": "postgres",
							"imagePullPolicy": "Never",
							"image": ""
						}
					],
					"hostAliases": [
						{
							"ip": "192.168.1.100",
							"hostnames": [
								"foo.local"
							]
						}
					]
				}
			}
		}
	]
}
`

func TestUnmarshalPathJSON(t *testing.T) {
	var patchTemplate corev1.PodSpec
	err := json.Unmarshal([]byte(patchJSONTest), &patchTemplate)
	require.NoError(t, err)
}

func TestCreatePatch(t *testing.T) {
	defaultPod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: "demo",
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:            "postgres",
					Image:           "postgres:latest",
					ImagePullPolicy: corev1.PullIfNotPresent,
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

	require.Equal(t, defaultPod.Spec.Containers[0].ImagePullPolicy, corev1.PullIfNotPresent)

	newPod, matched, err := applyPatchMultipleFromJSONString(defaultPod, patchMultiple)
	require.NoError(t, err)
	require.True(t, matched)
	require.NotEmpty(t, defaultPod.Spec.Containers)
	require.NotEmpty(t, newPod.Spec.Containers)
	require.Len(t, newPod.Spec.HostAliases, 2)
	require.Len(t, newPod.Spec.Containers, 1)
	require.Equal(t, corev1.PullNever, newPod.Spec.Containers[0].ImagePullPolicy)
	require.Equal(t, "postgres:latest", newPod.Spec.Containers[0].Image)

	patch, err := createPatch(*defaultPod, *newPod)
	require.NoError(t, err)
	require.NotEmpty(t, patch)
}

func TestMatchRegex(t *testing.T) {
	regex := "^(logical-backup|default-).*$"
	podName := "logical-backup-default-28715940-dbq9b"
	match := matchRegex(regex, podName)
	require.True(t, match)

	podName2 := "default-0"
	match2 := matchRegex(regex, podName2)
	require.True(t, match2)
}
