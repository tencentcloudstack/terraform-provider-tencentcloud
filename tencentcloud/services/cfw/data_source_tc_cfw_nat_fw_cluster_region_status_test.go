package cfw_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	cfwv20190904 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cfw/v20190904"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cfw"
)

type mockMetaNatFwClusterRegionStatus struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaNatFwClusterRegionStatus) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaNatFwClusterRegionStatus{}

func newMockMetaNatFwClusterRegionStatus() *mockMetaNatFwClusterRegionStatus {
	return &mockMetaNatFwClusterRegionStatus{client: &connectivity.TencentCloudClient{}}
}

func ptrStringNatFwCluster(s string) *string {
	return &s
}

func ptrInt64NatFwCluster(i int64) *int64 {
	return &i
}

// TestCfwNatFwClusterRegionStatusDS_ReadSuccess tests data source Read with successful response
func TestCfwNatFwClusterRegionStatusDS_ReadSuccess(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	cfwClient := &cfwv20190904.Client{}
	patches.ApplyMethodReturn(newMockMetaNatFwClusterRegionStatus().client, "UseCfwClient", cfwClient)

	patches.ApplyMethodFunc(cfwClient, "DescribeNatFwClusterRegionStatusWithContext", func(_ context.Context, request *cfwv20190904.DescribeNatFwClusterRegionStatusRequest) (*cfwv20190904.DescribeNatFwClusterRegionStatusResponse, error) {
		resp := cfwv20190904.NewDescribeNatFwClusterRegionStatusResponse()
		resp.Response = &cfwv20190904.DescribeNatFwClusterRegionStatusResponseParams{
			Total: ptrInt64NatFwCluster(2),
			RegionFwStatus: []*cfwv20190904.NatFwClusterRegionStatus{
				{
					NatInsId:    ptrStringNatFwCluster("nat-xxxxxxxx"),
					CcnId:       ptrStringNatFwCluster("ccn-fkb9bo2v"),
					Region:      ptrStringNatFwCluster("ap-guangzhou"),
					Status:      ptrStringNatFwCluster("Auto"),
					Cidr:        ptrStringNatFwCluster("10.0.0.0/24"),
					RoutingMode: ptrInt64NatFwCluster(0),
				},
				{
					NatInsId:    ptrStringNatFwCluster("nat-yyyyyyyy"),
					CcnId:       ptrStringNatFwCluster("ccn-fkb9bo2v"),
					Region:      ptrStringNatFwCluster("ap-beijing"),
					Status:      ptrStringNatFwCluster("NotDeployed"),
					RoutingMode: ptrInt64NatFwCluster(1),
				},
			},
			RequestId: ptrStringNatFwCluster("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaNatFwClusterRegionStatus()
	res := cfw.DataSourceTencentCloudCfwNatFwClusterRegionStatus()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.NotEmpty(t, d.Id())
	assert.Equal(t, 2, d.Get("total"))

	regionFwStatus := d.Get("region_fw_status").([]interface{})
	assert.Equal(t, 2, len(regionFwStatus))

	item0 := regionFwStatus[0].(map[string]interface{})
	assert.Equal(t, "nat-xxxxxxxx", item0["nat_ins_id"])
	assert.Equal(t, "ccn-fkb9bo2v", item0["ccn_id"])
	assert.Equal(t, "ap-guangzhou", item0["region"])
	assert.Equal(t, "Auto", item0["status"])
	assert.Equal(t, "10.0.0.0/24", item0["cidr"])
	assert.Equal(t, int(0), item0["routing_mode"])

	item1 := regionFwStatus[1].(map[string]interface{})
	assert.Equal(t, "nat-yyyyyyyy", item1["nat_ins_id"])
	assert.Equal(t, "ap-beijing", item1["region"])
	assert.Equal(t, "NotDeployed", item1["status"])
}

// TestCfwNatFwClusterRegionStatusDS_ReadNilResponse tests data source Read with nil response
func TestCfwNatFwClusterRegionStatusDS_ReadNilResponse(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	cfwClient := &cfwv20190904.Client{}
	patches.ApplyMethodReturn(newMockMetaNatFwClusterRegionStatus().client, "UseCfwClient", cfwClient)

	patches.ApplyMethodFunc(cfwClient, "DescribeNatFwClusterRegionStatusWithContext", func(_ context.Context, request *cfwv20190904.DescribeNatFwClusterRegionStatusRequest) (*cfwv20190904.DescribeNatFwClusterRegionStatusResponse, error) {
		resp := cfwv20190904.NewDescribeNatFwClusterRegionStatusResponse()
		resp.Response = nil
		return resp, nil
	})

	meta := newMockMetaNatFwClusterRegionStatus()
	res := cfw.DataSourceTencentCloudCfwNatFwClusterRegionStatus()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})

	err := res.Read(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "response is nil")
}

// TestCfwNatFwClusterRegionStatusDS_ReadAPIError tests data source Read with API error
func TestCfwNatFwClusterRegionStatusDS_ReadAPIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	cfwClient := &cfwv20190904.Client{}
	patches.ApplyMethodReturn(newMockMetaNatFwClusterRegionStatus().client, "UseCfwClient", cfwClient)

	patches.ApplyMethodFunc(cfwClient, "DescribeNatFwClusterRegionStatusWithContext", func(_ context.Context, request *cfwv20190904.DescribeNatFwClusterRegionStatusRequest) (*cfwv20190904.DescribeNatFwClusterRegionStatusResponse, error) {
		return nil, fmt.Errorf("api error: internal server error")
	})

	meta := newMockMetaNatFwClusterRegionStatus()
	res := cfw.DataSourceTencentCloudCfwNatFwClusterRegionStatus()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})

	err := res.Read(d, meta)
	assert.Error(t, err)
}

// TestCfwNatFwClusterRegionStatusDS_ReadEmptyRegionFwStatus tests data source Read with empty RegionFwStatus
func TestCfwNatFwClusterRegionStatusDS_ReadEmptyRegionFwStatus(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	cfwClient := &cfwv20190904.Client{}
	patches.ApplyMethodReturn(newMockMetaNatFwClusterRegionStatus().client, "UseCfwClient", cfwClient)

	patches.ApplyMethodFunc(cfwClient, "DescribeNatFwClusterRegionStatusWithContext", func(_ context.Context, request *cfwv20190904.DescribeNatFwClusterRegionStatusRequest) (*cfwv20190904.DescribeNatFwClusterRegionStatusResponse, error) {
		resp := cfwv20190904.NewDescribeNatFwClusterRegionStatusResponse()
		resp.Response = &cfwv20190904.DescribeNatFwClusterRegionStatusResponseParams{
			Total:          ptrInt64NatFwCluster(0),
			RegionFwStatus: nil,
			RequestId:      ptrStringNatFwCluster("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaNatFwClusterRegionStatus()
	res := cfw.DataSourceTencentCloudCfwNatFwClusterRegionStatus()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.NotEmpty(t, d.Id())
	assert.Equal(t, 0, d.Get("total"))

	regionFwStatus := d.Get("region_fw_status").([]interface{})
	assert.Equal(t, 0, len(regionFwStatus))
}
