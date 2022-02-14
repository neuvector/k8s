package v1alpha1

import "github.com/neuvector/k8s"

func init() {
	k8s.Register("scheduling.k8s.io", "v1alpha1", "priorityclasses", false, &PriorityClass{})

	k8s.RegisterList("scheduling.k8s.io", "v1alpha1", "priorityclasses", false, &PriorityClassList{})
}
