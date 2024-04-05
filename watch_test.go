package k8s_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/neuvector/k8s"
    corev1 "k8s.io/api/core/v1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// configMapJSON is used to test the JSON serialization watch.
type configMapJSON struct {
	metav1.ObjectMeta `json:"metadata"`
	Data     map[string]string  `json:"data"`
}

func init() {
	k8s.Register("", "v1", "configmaps", true, &configMapJSON{})
}

func testWatch(t *testing.T, client *k8s.Client, namespace string, newCM func() metav1.Object, update func(cm metav1.Object)) {
	w, err := client.Watch(context.TODO(), namespace, newCM())
	if err != nil {
		t.Errorf("watch configmaps: %v", err)
	}
	defer w.Close()

	cm := newCM()
	want := func(eventType string) {
		got := newCM()
		eT, err := w.Next(got)
		if err != nil {
			t.Errorf("decode watch event: %v", err)
			return
		}
		if eT != eventType {
			t.Errorf("expected event type %q got %q", eventType, eT)
		}
		cm.SetResourceVersion("")
		got.SetResourceVersion("")
		if !reflect.DeepEqual(got, cm) {
			t.Errorf("configmaps didn't match")
			t.Errorf("want: %#v", cm)
			t.Errorf(" got: %#v", got)
		}
	}

	if err := client.Create(context.TODO(), cm); err != nil {
		t.Errorf("create configmap: %v", err)
		return
	}
	want(k8s.EventAdded)

	update(cm)

	if err := client.Update(context.TODO(), cm); err != nil {
		t.Errorf("update configmap: %v", err)
		return
	}
	want(k8s.EventModified)

	if err := client.Delete(context.TODO(), cm); err != nil {
		t.Errorf("Delete configmap: %v", err)
		return
	}
	want(k8s.EventDeleted)
}

func TestWatchConfigMapJSON(t *testing.T) {
	withNamespace(t, func(client *k8s.Client, namespace string) {
		newCM := func() metav1.Object {
			return &configMapJSON{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "my-configmap",
					Namespace: namespace,
				},
			}
		}

		updateCM := func(cm metav1.Object) {
			(cm.(*configMapJSON)).Data = map[string]string{"hello": "world"}
		}
		testWatch(t, client, namespace, newCM, updateCM)
	})
}

func TestWatchConfigMapProto(t *testing.T) {
	withNamespace(t, func(client *k8s.Client, namespace string) {
		newCM := func() metav1.Object {
			return &corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "my-configmap",
					Namespace: namespace,
				},
			}
		}

		updateCM := func(cm metav1.Object) {
			(cm.(*corev1.ConfigMap)).Data = map[string]string{"hello": "world"}
		}
		testWatch(t, client, namespace, newCM, updateCM)
	})
}
