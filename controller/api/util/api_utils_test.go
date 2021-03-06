package util

import (
	"testing"

	pb "github.com/runconduit/conduit/controller/gen/public"
	"github.com/runconduit/conduit/pkg/k8s"
)

func TestGetWindow(t *testing.T) {
	t.Run("Returns valid windows", func(t *testing.T) {
		expectations := map[string]pb.TimeWindow{
			"10s": pb.TimeWindow_TEN_SEC,
			"1m":  pb.TimeWindow_ONE_MIN,
			"10m": pb.TimeWindow_TEN_MIN,
			"1h":  pb.TimeWindow_ONE_HOUR,
		}

		for windowFriendlyName, expectedTimeWindow := range expectations {
			actualTimeWindow, err := GetWindow(windowFriendlyName)
			if err != nil {
				t.Fatalf("Unexpected error when resolving time window friendly name [%s]: %v",
					windowFriendlyName, err)
			}

			if actualTimeWindow != expectedTimeWindow {
				t.Fatalf("Expected resolving friendly name [%s] to return timw window [%v], but got [%v]",
					windowFriendlyName, expectedTimeWindow, actualTimeWindow)
			}
		}
	})

	t.Run("Returns error and default value if unknown friendly name for TimeWindow", func(t *testing.T) {
		invalidNames := []string{
			"10seconds", "10sec", "9s",
			"10minutes", "10min", "9m",
			"1minute", "1min", "0s", "2s",
			"1hour", "0h", "2h",
			"10", ""}
		defaultTimeWindow := pb.TimeWindow_ONE_MIN

		for _, invalidName := range invalidNames {
			window, err := GetWindow(invalidName)
			if err == nil {
				t.Fatalf("Expected invalid friendly name [%s] to generate error, but got no error and result [%v]",
					invalidName, window)
			}

			if window != defaultTimeWindow {
				t.Fatalf("Expected invalid friendly name resolution to return default window [%v], but got [%v]",
					defaultTimeWindow, window)
			}
		}
	})
}

func TestGetWindowString(t *testing.T) {
	t.Run("Returns names for valid windows", func(t *testing.T) {
		expectations := map[pb.TimeWindow]string{
			pb.TimeWindow_TEN_SEC:  "10s",
			pb.TimeWindow_ONE_MIN:  "1m",
			pb.TimeWindow_TEN_MIN:  "10m",
			pb.TimeWindow_ONE_HOUR: "1h",
		}

		for window, expectedName := range expectations {
			actualName, err := GetWindowString(window)
			if err != nil {
				t.Fatalf("Unexpected error when resolving name for window [%v]: %v", window, err)
			}

			if actualName != expectedName {
				t.Fatalf("Expected window [%v] to resolve to name [%s], but got [%s]", window, expectedName, actualName)
			}
		}
	})
}

func TestBuildStatSummaryRequest(t *testing.T) {
	t.Run("Maps Kubernetes friendly names to canonical names", func(t *testing.T) {
		expectations := map[string]string{
			"deployments": k8s.KubernetesDeployments,
			"deployment":  k8s.KubernetesDeployments,
			"deploy":      k8s.KubernetesDeployments,
			"pods":        k8s.KubernetesPods,
			"pod":         k8s.KubernetesPods,
			"po":          k8s.KubernetesPods,
		}

		for friendly, canonical := range expectations {
			statSummaryRequest, err := BuildStatSummaryRequest(
				StatSummaryRequestParams{
					ResourceType: friendly,
				},
			)
			if err != nil {
				t.Fatalf("Unexpected error from BuildStatSummaryRequest [%s => %s]: %s", friendly, canonical, err)
			}
			if statSummaryRequest.Selector.Resource.Type != canonical {
				t.Fatalf("Unexpected resource type from BuildStatSummaryRequest [%s => %s]: %s", friendly, canonical, statSummaryRequest.Selector.Resource.Type)
			}
		}
	})

	t.Run("Rejects invalid Kubernetes resource types", func(t *testing.T) {
		expectations := map[string]string{
			"foo": "cannot find Kubernetes canonical name from friendly name [foo]",
			"":    "cannot find Kubernetes canonical name from friendly name []",
		}

		for input, msg := range expectations {
			_, err := BuildStatSummaryRequest(
				StatSummaryRequestParams{
					ResourceType: input,
				},
			)
			if err == nil {
				t.Fatalf("BuildStatSummaryRequest(%s) unexpectedly succeeded, should have returned %s", input, msg)
			}
			if err.Error() != msg {
				t.Fatalf("BuildStatSummaryRequest(%s) should have returned: %s but got unexpected message: %s", input, msg, err)
			}
		}
	})
}
