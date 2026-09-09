package dc_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	dc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dc/v20180410"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	dcservice "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/dc"
)

type mockMetaDcx struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaDcx) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaDcx{}

func newMockMetaDcx() *mockMetaDcx {
	return &mockMetaDcx{client: &connectivity.TencentCloudClient{}}
}

func ptrStringDcx(s string) *string {
	return &s
}

func ptrInt64Dcx(v int64) *int64 {
	return &v
}

// go test ./tencentcloud/services/dc/ -run "TestDcxInstance" -v -count=1 -gcflags="all=-l"

// TestDcxInstance_UpdateBandwidth_Success tests that updating bandwidth calls
// ModifyDirectConnectTunnelAttribute with the expected Bandwidth value and DirectConnectTunnelId.
func TestDcxInstance_UpdateBandwidth_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	dcClient := &dc.Client{}
	patches.ApplyMethodReturn(newMockMetaDcx().client, "UseDcClient", dcClient)

	var capturedRequest *dc.ModifyDirectConnectTunnelAttributeRequest
	patches.ApplyMethodFunc(dcClient, "ModifyDirectConnectTunnelAttribute", func(request *dc.ModifyDirectConnectTunnelAttributeRequest) (*dc.ModifyDirectConnectTunnelAttributeResponse, error) {
		capturedRequest = request
		resp := dc.NewModifyDirectConnectTunnelAttributeResponse()
		resp.Response = &dc.ModifyDirectConnectTunnelAttributeResponseParams{
			RequestId: ptrStringDcx("fake-request-id"),
		}
		return resp, nil
	})

	// Mock the service Describe method used by Read
	var service dcservice.DcService
	patches.ApplyMethodFunc(&service, "DescribeDirectConnectTunnel", func(ctx context.Context, dcxId string) (dc.DirectConnectTunnel, int64, error) {
		info := dc.DirectConnectTunnel{
			DirectConnectTunnelId:   ptrStringDcx("dcx-test-id"),
			DirectConnectTunnelName: ptrStringDcx("test-dcx"),
			Bandwidth:               ptrInt64Dcx(100),
			State:                   ptrStringDcx("AVAILABLE"),
		}
		return info, 1, nil
	})

	meta := newMockMetaDcx()
	res := dcservice.ResourceTencentCloudDcxInstance()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"dc_id":     "dc-test-id",
		"name":      "test-dcx",
		"dcg_id":    "dcg-test-id",
		"bandwidth": 100,
	})
	d.SetId("dcx-test-id")
	d.MarkNewResource()

	// Simulate bandwidth change from 50 to 100
	d.SetNew("bandwidth", 100)

	err := res.Update(d, meta)
	assert.NoError(t, err)

	// Verify the Modify API was called with expected Bandwidth and DirectConnectTunnelId
	assert.NotNil(t, capturedRequest)
	assert.NotNil(t, capturedRequest.Bandwidth)
	assert.Equal(t, int64(100), *capturedRequest.Bandwidth)
	assert.NotNil(t, capturedRequest.DirectConnectTunnelId)
	assert.Equal(t, "dcx-test-id", *capturedRequest.DirectConnectTunnelId)
}

// TestDcxInstance_UpdateBandwidth_APIError tests that an API error during bandwidth
// update is propagated correctly.
func TestDcxInstance_UpdateBandwidth_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	dcClient := &dc.Client{}
	patches.ApplyMethodReturn(newMockMetaDcx().client, "UseDcClient", dcClient)

	patches.ApplyMethodFunc(dcClient, "ModifyDirectConnectTunnelAttribute", func(request *dc.ModifyDirectConnectTunnelAttributeRequest) (*dc.ModifyDirectConnectTunnelAttributeResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=InvalidParameter, Message=Invalid bandwidth value")
	})

	// Mock the service Describe method used by Read (won't be reached due to update error)
	var service dcservice.DcService
	patches.ApplyMethodFunc(&service, "DescribeDirectConnectTunnel", func(ctx context.Context, dcxId string) (dc.DirectConnectTunnel, int64, error) {
		info := dc.DirectConnectTunnel{
			DirectConnectTunnelId: ptrStringDcx("dcx-test-id"),
			Bandwidth:             ptrInt64Dcx(100),
		}
		return info, 1, nil
	})

	meta := newMockMetaDcx()
	res := dcservice.ResourceTencentCloudDcxInstance()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"dc_id":     "dc-test-id",
		"name":      "test-dcx",
		"dcg_id":    "dcg-test-id",
		"bandwidth": 100,
	})
	d.SetId("dcx-test-id")
	d.MarkNewResource()
	d.SetNew("bandwidth", 100)

	err := res.Update(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "InvalidParameter")
}

// TestDcxInstance_Schema_BandwidthNotForceNew validates that the bandwidth schema field
// does NOT have ForceNew set (it should be updatable).
func TestDcxInstance_Schema_BandwidthNotForceNew(t *testing.T) {
	res := dcservice.ResourceTencentCloudDcxInstance()

	assert.Contains(t, res.Schema, "bandwidth")
	bandwidth := res.Schema["bandwidth"]
	assert.Equal(t, schema.TypeInt, bandwidth.Type)
	assert.True(t, bandwidth.Optional)
	assert.True(t, bandwidth.Computed)
	assert.False(t, bandwidth.ForceNew, "bandwidth should NOT have ForceNew=true (must be updatable)")
}
