package test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"k8gbterratest/utils"
)

func TestWeightsExistsInLocalDNSEndpoint(t *testing.T) {
	t.Parallel()
	const ingressHost = "terratest-roundrobin.cloud.example.com"
	const host = "weight-roundrobin.cloud.example.com"
	const gslbPath = "../examples/roundrobin_weight1.yaml"
	instanceEU, err := utils.NewWorkflow(t, "k3d-test-gslb1", 5053).
		WithGslb(gslbPath, host).
		WithTestApp("eu").
		Start()
	require.NoError(t, err)
	defer instanceEU.Kill()

	instanceUS, err := utils.NewWorkflow(t, "k3d-test-gslb1", 5053).
		WithGslb(gslbPath, host).
		WithTestApp("us").
		Start()
	require.NoError(t, err)
	defer instanceEU.Kill()

	err = instanceEU.WaitForAppIsRunning()
	require.NoError(t, err)
	err = instanceUS.WaitForAppIsRunning()
	require.NoError(t, err)

	epeu, err := instanceEU.GetExternalDNSEndpoint().GetEndpointByName(ingressHost)
	require.NoError(t, err, "missing EU endpoint", ingressHost)
	epus, err := instanceUS.GetExternalDNSEndpoint().GetEndpointByName(ingressHost)
	require.NoError(t, err, "missing US endpoint", ingressHost)

	require.Equal(t, "roundRobin", epeu.Labels["strategy"])
	require.Equal(t, "roundRobin", epus.Labels["strategy"])
	require.True(t, len(epus.Labels) > 1, "EU endpoint", ingressHost, " doesn't contain weight labels")
	require.True(t, len(epeu.Labels) > 1, "EU endpoint", ingressHost, " doesn't contain weight labels")
}
