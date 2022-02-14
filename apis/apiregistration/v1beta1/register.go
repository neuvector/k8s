package v1beta1

import "github.com/neuvector/k8s"

func init() {
	k8s.Register("apiregistration.k8s.io", "v1beta1", "apiservices", false, &APIService{})

	k8s.RegisterList("apiregistration.k8s.io", "v1beta1", "apiservices", false, &APIServiceList{})
}
