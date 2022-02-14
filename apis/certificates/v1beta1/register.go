package v1beta1

import "github.com/neuvector/k8s"

func init() {
	k8s.Register("certificates.k8s.io", "v1beta1", "certificatesigningrequests", false, &CertificateSigningRequest{})

	k8s.RegisterList("certificates.k8s.io", "v1beta1", "certificatesigningrequests", false, &CertificateSigningRequestList{})
}
