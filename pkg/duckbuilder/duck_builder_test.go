package duckbuilder

import (
	"github.com/sclevine/spec"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	duckv1alpha1 "knative.dev/pkg/apis/duck/v1alpha1"

	"github.com/pivotal/kpack/pkg/apis/build/v1alpha1"
	"testing"
)

func TestDuckBuilder(t *testing.T) {
	spec.Run(t, "testDuckBuilder", testDuckBuilder)
}

func testDuckBuilder(t *testing.T, when spec.G, it spec.S) {
	duckBuilder := &DuckBuilder{
		ObjectMeta: metav1.ObjectMeta{
			Generation: 1,
		},
		Spec: DuckBuilderSpec{
			ImagePullSecrets: []corev1.LocalObjectReference{
				{
					Name: "test-secret",
				},
			},
		},
		Status: v1alpha1.BuilderStatus{
			Status: duckv1alpha1.Status{
				ObservedGeneration: 1,
				Conditions: duckv1alpha1.Conditions{
					{
						Type:   duckv1alpha1.ConditionReady,
						Status: corev1.ConditionTrue,
					},
				},
			},
			BuilderMetadata: v1alpha1.BuildpackMetadataList{
				{
					ID:      "test.builder",
					Version: "test.version",
				},
			},
			Stack: v1alpha1.BuildStack{
				RunImage: "some/run@sha256:12345678",
			},
			LatestImage: "some/builder@sha256:12345678",
		},
	}

	when("Ready", func() {

		it("ready when ready condition is true", func() {
			require.True(t, duckBuilder.Ready())
		})

		it("not ready without conditions", func() {
			duckBuilder.Status.Conditions = nil

			require.False(t, duckBuilder.Ready())
		})

		it("not ready when not ready", func() {
			duckBuilder.Status.Conditions = duckv1alpha1.Conditions{
				{
					Type:   duckv1alpha1.ConditionReady,
					Status: corev1.ConditionUnknown,
				},
			}

			require.False(t, duckBuilder.Ready())
		})

		it("not ready when generation does not match observed generation", func() {
			duckBuilder.Generation = duckBuilder.Status.ObservedGeneration + 1

			require.False(t, duckBuilder.Ready())
		})
	})

	it("BuildBuilderSpec provides latest image and pull secrets", func() {
		require.Equal(t, v1alpha1.BuildBuilderSpec{
			Image: "some/builder@sha256:12345678",
			ImagePullSecrets: []corev1.LocalObjectReference{
				{
					Name: "test-secret",
				},
			},
		}, duckBuilder.BuildBuilderSpec())
	})

	it("BuildpackMetadata provides buildpack metadata", func() {
		require.Equal(t, v1alpha1.BuildpackMetadataList{
			{
				ID:      "test.builder",
				Version: "test.version",
			},
		}, duckBuilder.BuildpackMetadata())
	})

	it("RunImage provides latest runimage", func() {
		require.Equal(t, "some/run@sha256:12345678", duckBuilder.RunImage())
	})

}