package controllers

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

const (
	timeout  = time.Second * 10
	interval = time.Millisecond * 600
)

var _ = Describe("Pod controller", func() {
	When("pods are reconciled with keep=0", func() {
		It("gc collects the necessary pods", func() {
			var pods = []corev1.Pod{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "pending-pod",
						Labels: map[string]string{
							"test-expect-keep": "true",
						},
						Namespace: "default",
					},
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{
							{
								Name:  "x",
								Image: "x",
							},
						},
					},
					Status: corev1.PodStatus{
						Phase: corev1.PodFailed,
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "evicted-unmanaged-pod",
						Labels: map[string]string{
							"test-expect-keep": "true",
						},
						Namespace: "default",
					},
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{
							{
								Name:  "x",
								Image: "x",
							},
						},
					},
					Status: corev1.PodStatus{
						Reason: "Evicted",
						Phase:  corev1.PodFailed,
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "evicted-managed-pod",
						Labels: map[string]string{
							"test-expect-keep": "false",
						},
						OwnerReferences: []metav1.OwnerReference{
							{
								APIVersion: "v1",
								Kind:       "ReplicaSet",
								Name:       "parent",
								UID:        "x",
							},
						},
						Namespace: "default",
					},
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{
							{
								Name:  "x",
								Image: "x",
							},
						},
					},
					Status: corev1.PodStatus{
						Reason: "Evicted",
						Phase:  corev1.PodFailed,
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "running-managed-pod",
						Labels: map[string]string{
							"test-expect-keep": "true",
						},
						OwnerReferences: []metav1.OwnerReference{
							{
								APIVersion: "v1",
								Kind:       "ReplicaSet",
								Name:       "parent",
								UID:        "x",
							},
						},
						Namespace: "default",
					},
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{
							{
								Name:  "x",
								Image: "x",
							},
						},
					},
					Status: corev1.PodStatus{
						Phase: corev1.PodRunning,
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "another-evicted-managed-pod",
						Labels: map[string]string{
							"test-expect-keep": "false",
						},
						OwnerReferences: []metav1.OwnerReference{
							{
								APIVersion: "v1",
								Kind:       "ReplicaSet",
								Name:       "parent",
								UID:        "x",
							},
						},
						Namespace: "default",
					},
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{
							{
								Name:  "x",
								Image: "x",
							},
						},
					},
					Status: corev1.PodStatus{
						Reason: "Evicted",
						Phase:  corev1.PodFailed,
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "yet-another-evicted-managed-pod",
						Labels: map[string]string{
							"test-expect-keep": "false",
						},
						OwnerReferences: []metav1.OwnerReference{
							{
								APIVersion: "v1",
								Kind:       "ReplicaSet",
								Name:       "parent",
								UID:        "x",
							},
						},
						Namespace: "default",
					},
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{
							{
								Name:  "x",
								Image: "x",
							},
						},
					},
					Status: corev1.PodStatus{
						Reason: "Evicted",
						Phase:  corev1.PodFailed,
					},
				},
			}

			for i, pod := range pods {
				Expect(k8sClient.Create(ctx, &pod)).Should(Succeed())
				pod.Status = pods[i].Status
				Expect(k8sClient.Status().Update(ctx, &pod)).Should(Succeed())

			}

			Eventually(func() bool {
				for _, pod := range pods {
					var p corev1.Pod
					err := k8sClient.Get(ctx, types.NamespacedName{
						Name:      pod.Name,
						Namespace: pod.Namespace,
					}, &p)

					if pod.Labels["test-expect-keep"] == "true" && err != nil {
						return false
					}

					if pod.Labels["test-expect-keep"] == "false" && err == nil {
						return false
					}
				}

				return true
			}, timeout, interval).Should(BeTrue())
		})
	})

	When("pods are reconciled with keep=2", func() {
		It("gc collects the necessary pods", func() {
			podController.Keep = 2
			var pods = []corev1.Pod{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "keep2-pending-pod",
						Labels: map[string]string{
							"test-expect-keep": "true",
						},
						Namespace: "default",
					},
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{
							{
								Name:  "x",
								Image: "x",
							},
						},
					},
					Status: corev1.PodStatus{
						Phase: corev1.PodFailed,
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "keep2-evicted-unmanaged-pod",
						Labels: map[string]string{
							"test-expect-keep": "true",
						},
						Namespace: "default",
					},
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{
							{
								Name:  "x",
								Image: "x",
							},
						},
					},
					Status: corev1.PodStatus{
						Reason: "Evicted",
						Phase:  corev1.PodFailed,
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "keep2-evicted-managed-pod",
						Labels: map[string]string{
							"test-expect-keep": "false",
						},
						OwnerReferences: []metav1.OwnerReference{
							{
								APIVersion: "v1",
								Kind:       "ReplicaSet",
								Name:       "parent",
								UID:        "x",
							},
						},
						Namespace: "default",
					},
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{
							{
								Name:  "x",
								Image: "x",
							},
						},
					},
					Status: corev1.PodStatus{
						Reason: "Evicted",
						Phase:  corev1.PodFailed,
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "keep2-running-managed-pod",
						Labels: map[string]string{
							"test-expect-keep": "true",
						},
						OwnerReferences: []metav1.OwnerReference{
							{
								APIVersion: "v1",
								Kind:       "ReplicaSet",
								Name:       "parent",
								UID:        "x",
							},
						},
						Namespace: "default",
					},
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{
							{
								Name:  "x",
								Image: "x",
							},
						},
					},
					Status: corev1.PodStatus{
						Phase: corev1.PodRunning,
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "keep2-another-evicted-managed-pod",
						Labels: map[string]string{
							"test-expect-keep": "true",
						},
						OwnerReferences: []metav1.OwnerReference{
							{
								APIVersion: "v1",
								Kind:       "ReplicaSet",
								Name:       "parent",
								UID:        "x",
							},
						},
						Namespace: "default",
					},
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{
							{
								Name:  "x",
								Image: "x",
							},
						},
					},
					Status: corev1.PodStatus{
						Reason: "Evicted",
						Phase:  corev1.PodFailed,
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "keep2-yet-another-evicted-managed-pod",
						Labels: map[string]string{
							"test-expect-keep": "true",
						},
						OwnerReferences: []metav1.OwnerReference{
							{
								APIVersion: "v1",
								Kind:       "ReplicaSet",
								Name:       "parent",
								UID:        "x",
							},
						},
						Namespace: "default",
					},
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{
							{
								Name:  "x",
								Image: "x",
							},
						},
					},
					Status: corev1.PodStatus{
						Reason: "Evicted",
						Phase:  corev1.PodFailed,
					},
				},
			}

			for i, pod := range pods {
				//we can't change .metadata.creationTime, so just sleep 1 sec between creations
				time.Sleep(time.Second)
				Expect(k8sClient.Create(ctx, &pod)).Should(Succeed())
				pod.Status = pods[i].Status
				Expect(k8sClient.Status().Update(ctx, &pod)).Should(Succeed())
			}

			Eventually(func() bool {
				for _, pod := range pods {
					var p corev1.Pod
					err := k8sClient.Get(ctx, types.NamespacedName{
						Name:      pod.Name,
						Namespace: pod.Namespace,
					}, &p)

					if pod.Labels["test-expect-keep"] == "true" && err != nil {
						return false
					}

					if pod.Labels["test-expect-keep"] == "false" && err == nil {
						return false
					}
				}

				return true
			}, timeout, interval).Should(BeTrue())
		})
	})

	When("evicted pods older than maxAge are always deleted", func() {
		It("gc collects the necessary pods", func() {
			podController.Keep = 1
			podController.MaxAge = time.Second

			var pods = []corev1.Pod{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "maxage-pending-pod",
						Labels: map[string]string{
							"test-expect-keep": "true",
						},
						Namespace: "default",
					},
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{
							{
								Name:  "x",
								Image: "x",
							},
						},
					},
					Status: corev1.PodStatus{
						Phase: corev1.PodFailed,
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "maxage-evicted-unmanaged-pod",
						Labels: map[string]string{
							"test-expect-keep": "false",
						},
						Namespace: "default",
					},
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{
							{
								Name:  "x",
								Image: "x",
							},
						},
					},
					Status: corev1.PodStatus{
						Reason: "Evicted",
						Phase:  corev1.PodFailed,
					},
				},
			}

			for i, pod := range pods {
				Expect(k8sClient.Create(ctx, &pod)).Should(Succeed())
				pod.Status = pods[i].Status
				Expect(k8sClient.Status().Update(ctx, &pod)).Should(Succeed())
			}

			Eventually(func() bool {
				for _, pod := range pods {
					var p corev1.Pod
					err := k8sClient.Get(ctx, types.NamespacedName{
						Name:      pod.Name,
						Namespace: pod.Namespace,
					}, &p)

					if pod.Labels["test-expect-keep"] == "true" && err != nil {
						return false
					}

					if pod.Labels["test-expect-keep"] == "false" && err == nil {
						return false
					}
				}

				return true
			}, timeout, interval).Should(BeTrue())
		})
	})
})
