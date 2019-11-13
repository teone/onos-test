// Copyright 2019-present Open Networking Foundation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package config

import (
	"context"
	"github.com/google/uuid"
	"github.com/onosproject/onos-config/pkg/northbound/admin"
	"github.com/onosproject/onos-test/pkg/onit/env"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func (s *SmokeTestSuite) TestSnapshot(t *testing.T) {
	simulator1 := s.addSimulator(t)
	simulator2 := s.addSimulator(t)

	// Make a GNMI client to use for requests
	c, err := env.Config().NewGNMIClient()
	assert.NoError(t, err)
	assert.True(t, c != nil, "Fetching client returned nil")

	for i := 0; i < 100; i++ {
		setPath := makeDevicePath(simulator1.Name(), "/system/config/motd-banner")
		setPath[0].pathDataValue = uuid.New().String()
		setPath[0].pathDataType = StringVal
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		_, _, err := GNMISet(ctx, c, setPath, noPaths)
		cancel()
		assert.NoError(t, err)

		setPath = makeDevicePath(simulator2.Name(), "/system/config/motd-banner")
		setPath[0].pathDataValue = uuid.New().String()
		setPath[0].pathDataType = StringVal
		ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		_, _, err = GNMISet(ctx, c, setPath, noPaths)
		cancel()
		assert.NoError(t, err)
	}

	time.Sleep(10 * time.Second)

	adminClient, err := env.Config().NewAdminServiceClient()
	assert.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	retention := 5 * time.Second
	_, err = adminClient.CompactChanges(ctx, &admin.CompactChangesRequest{
		RetentionPeriod: &retention,
	})
	cancel()
	assert.NoError(t, err)
}
