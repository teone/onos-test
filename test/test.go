package test

import (
	"fmt"
	"github.com/onosproject/onos-test/pkg/helm"
	"github.com/onosproject/onos-test/pkg/test"
	"github.com/stretchr/testify/assert"
	"testing"
)

type ChartTestSuite struct {
	test.Suite
}

func (s *ChartTestSuite) TestLocalInstall(t *testing.T) {
	namespace := helm.Namespace()
	atomix := namespace.Chart("/etc/charts/atomix-controller").
		Release("atomix-controller").
		Set("namespace", namespace.Namespace())
	err := atomix.Install(true)
	assert.NoError(t, err)

	topo := helm.Namespace().
		Chart("/etc/charts/onos-topo").
		Release("onos-topo").
		Set("store.controller", fmt.Sprintf("atomix-controller.%s.svc.cluster.local:5679", namespace.Namespace()))
	err = topo.Install(true)
	assert.NoError(t, err)

	deployment, err := topo.Apps().V1().Deployments().Get("onos-topo")
	assert.NoError(t, err)

	pods, err := deployment.Pods().List()
	assert.NoError(t, err)
	assert.Len(t, pods, 1)
}

func (s *ChartTestSuite) TestRemoteInstall(t *testing.T) {
	kafka := helm.Namespace().
		Chart("kafka").
		SetRepository("http://storage.googleapis.com/kubernetes-charts-incubator").
		Release("device-simulator-test").
		Set("replicas", 1).
		Set("zookeeper.replicaCount", 1)
	err := kafka.Install(true)
	assert.NoError(t, err)
}