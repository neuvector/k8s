package k8s

import (
	"testing"
	"time"

    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Redefine types since all API groups import "github.com/neuvector/k8s"
// We can't use them here because it'll create a circular import cycle.

type Pod struct {
	metav1.ObjectMeta
}

type PodList struct {
	metav1.ListMeta
}

type Deployment struct {
	metav1.ObjectMeta
}

type DeploymentList struct {
	metav1.ListMeta
}

type ClusterRole struct {
	metav1.ObjectMeta
}

type ClusterRoleList struct {
	metav1.ListMeta
}

func init() {
	Register("", "v1", "pods", true, &Pod{})
	RegisterList("", "v1", "pods", true, &PodList{})

	Register("apps", "v1beta2", "deployments", true, &Deployment{})
	RegisterList("apps", "v1beta2", "deployments", true, &DeploymentList{})

	Register("rbac.authorization.k8s.io", "v1", "clusterroles", false, &ClusterRole{})
	RegisterList("rbac.authorization.k8s.io", "v1", "clusterroles", false, &ClusterRoleList{})
}

func TestResourceURL(t *testing.T) {
	tests := []struct {
		name     string
		endpoint string
		resource metav1.Object
		withName bool
		options  []Option
		want     string
		wantErr  bool
	}{
		{
			name:     "pod",
			endpoint: "https://example.com",
			resource: &Pod{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "my-namespace",
					Name:      "my-pod",
				},
			},
			want: "https://example.com/api/v1/namespaces/my-namespace/pods",
		},
		{
			name:     "deployment",
			endpoint: "https://example.com",
			resource: &Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "my-namespace",
					Name:      "my-deployment",
				},
			},
			want: "https://example.com/apis/apps/v1beta2/namespaces/my-namespace/deployments",
		},
		{
			name:     "deployment-with-name",
			endpoint: "https://example.com",
			resource: &Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "my-namespace",
					Name:      "my-deployment",
				},
			},
			withName: true,
			want:     "https://example.com/apis/apps/v1beta2/namespaces/my-namespace/deployments/my-deployment",
		},
		{
			name:     "deployment-with-subresource",
			endpoint: "https://example.com",
			resource: &Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "my-namespace",
					Name:      "my-deployment",
				},
			},
			withName: true,
			options: []Option{
				Subresource("status"),
			},
			want: "https://example.com/apis/apps/v1beta2/namespaces/my-namespace/deployments/my-deployment/status",
		},
		{
			name:     "pod-with-timeout",
			endpoint: "https://example.com",
			resource: &Pod{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "my-namespace",
					Name:      "my-pod",
				},
			},
			options: []Option{
				Timeout(time.Minute),
			},
			want: "https://example.com/api/v1/namespaces/my-namespace/pods?timeoutSeconds=60",
		},
		{
			name:     "pod-with-resource-version",
			endpoint: "https://example.com",
			resource: &Pod{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "my-namespace",
					Name:      "my-pod",
				},
			},
			options: []Option{
				ResourceVersion("foo"),
			},
			want: "https://example.com/api/v1/namespaces/my-namespace/pods?resourceVersion=foo",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := resourceURL(test.endpoint, test.resource, test.withName, test.options...)
			if err != nil {
				if test.wantErr {
					return
				}
				t.Fatalf("constructing resource URL: %v", err)
			}
			if test.wantErr {
				t.Fatal("expected error")
			}
			if test.want != got {
				t.Errorf("wanted=%q, got=%q", test.want, got)
			}
		})
	}
}
