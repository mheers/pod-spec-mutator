package models

import corev1 "k8s.io/api/core/v1"

type Patch struct {
	PodNameRegex string     `json:"podNameRegex"`
	Pod          corev1.Pod `json:"pod"`
}
