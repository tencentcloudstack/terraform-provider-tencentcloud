package teo_test

import (
	"context"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/teo"
)

// mockMetaOriginGroupHealthStatusDS implements tccommon.ProviderMeta
type mockMetaOriginGroupHealthStatusDS struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaOriginGroupHealthStatusDS) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaOriginGroupHealthStatusDS{}

func newMockMetaOriginGroupHealthStatusDS() *mockMetaOriginGroupHealthStatusDS {
	return &mockMetaOriginGroupHealthStatusDS{client: &connectivity.TencentCloudClient{}}
}

func ptrStringOriginGroupHealthStatusDS(s string) *string {
	return &s
}

// go test ./tencentcloud/services/teo/ -run "TestTeoOriginGroupHealthStatusDS" -v -count=1 -gcflags="all=-l"

// TestTeoOriginGroupHealthStatusDS_ReadWithRequiredParams tests data source Read with required parameters
func TestTeoOriginGroupHealthStatusDS_ReadWithRequiredParams(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaOriginGroupHealthStatusDS().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeOriginGroupHealthStatusWithContext", func(_ context.Context, request *teov20220901.DescribeOriginGroupHealthStatusRequest) (*teov20220901.DescribeOriginGroupHealthStatusResponse, error) {
		resp := teov20220901.NewDescribeOriginGroupHealthStatusResponse()
		resp.Response = &teov20220901.DescribeOriginGroupHealthStatusResponseParams{
			OriginGroupHealthStatusList: []*teov20220901.OriginGroupHealthStatusDetail{
				{
					OriginGroupId: ptrStringOriginGroupHealthStatusDS("origin-group-xxx1"),
					OriginHealthStatus: []*teov20220901.OriginHealthStatus{
						{
							Origin:  ptrStringOriginGroupHealthStatusDS("1.1.1.1"),
							Healthy: ptrStringOriginGroupHealthStatusDS("Healthy"),
						},
						{
							Origin:  ptrStringOriginGroupHealthStatusDS("2.2.2.2"),
							Healthy: ptrStringOriginGroupHealthStatusDS("Unhealthy"),
						},
					},
					CheckRegionHealthStatus: []*teov20220901.CheckRegionHealthStatus{
						{
							Region:  ptrStringOriginGroupHealthStatusDS("CN"),
							Healthy: ptrStringOriginGroupHealthStatusDS("Healthy"),
							OriginHealthStatus: []*teov20220901.OriginHealthStatus{
								{
									Origin:  ptrStringOriginGroupHealthStatusDS("1.1.1.1"),
									Healthy: ptrStringOriginGroupHealthStatusDS("Healthy"),
								},
							},
						},
					},
				},
			},
		}
		return resp, nil
	})

	meta := newMockMetaOriginGroupHealthStatusDS()
	res := teo.DataSourceTencentCloudTeoOriginGroupHealthStatus()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":        "zone-12345678",
		"lb_instance_id": "lb-instance-xxx",
	})

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.NotEmpty(t, d.Id())

	healthStatusList := d.Get("origin_group_health_status_list").([]interface{})
	assert.Len(t, healthStatusList, 1)
	healthStatusMap := healthStatusList[0].(map[string]interface{})
	assert.Equal(t, "origin-group-xxx1", healthStatusMap["origin_group_id"].(string))

	originHealthStatus := healthStatusMap["origin_health_status"].([]interface{})
	assert.Len(t, originHealthStatus, 2)
	originHealthStatusMap := originHealthStatus[0].(map[string]interface{})
	assert.Equal(t, "1.1.1.1", originHealthStatusMap["origin"].(string))
	assert.Equal(t, "Healthy", originHealthStatusMap["healthy"].(string))

	checkRegionHealthStatus := healthStatusMap["check_region_health_status"].([]interface{})
	assert.Len(t, checkRegionHealthStatus, 1)
	checkRegionMap := checkRegionHealthStatus[0].(map[string]interface{})
	assert.Equal(t, "CN", checkRegionMap["region"].(string))
	assert.Equal(t, "Healthy", checkRegionMap["healthy"].(string))

	regionOriginHealthStatus := checkRegionMap["origin_health_status"].([]interface{})
	assert.Len(t, regionOriginHealthStatus, 1)
	regionOriginMap := regionOriginHealthStatus[0].(map[string]interface{})
	assert.Equal(t, "1.1.1.1", regionOriginMap["origin"].(string))
	assert.Equal(t, "Healthy", regionOriginMap["healthy"].(string))
}

// TestTeoOriginGroupHealthStatusDS_ReadWithOriginGroupIds tests data source Read with origin_group_ids filter
func TestTeoOriginGroupHealthStatusDS_ReadWithOriginGroupIds(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaOriginGroupHealthStatusDS().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeOriginGroupHealthStatusWithContext", func(_ context.Context, request *teov20220901.DescribeOriginGroupHealthStatusRequest) (*teov20220901.DescribeOriginGroupHealthStatusResponse, error) {
		assert.NotNil(t, request.ZoneId)
		assert.Equal(t, "zone-12345678", *request.ZoneId)
		assert.NotNil(t, request.LBInstanceId)
		assert.Equal(t, "lb-instance-xxx", *request.LBInstanceId)
		assert.NotNil(t, request.OriginGroupIds)
		assert.Len(t, request.OriginGroupIds, 2)
		assert.Equal(t, "origin-group-xxx1", *request.OriginGroupIds[0])
		assert.Equal(t, "origin-group-xxx2", *request.OriginGroupIds[1])

		resp := teov20220901.NewDescribeOriginGroupHealthStatusResponse()
		resp.Response = &teov20220901.DescribeOriginGroupHealthStatusResponseParams{
			OriginGroupHealthStatusList: []*teov20220901.OriginGroupHealthStatusDetail{
				{
					OriginGroupId: ptrStringOriginGroupHealthStatusDS("origin-group-xxx1"),
					OriginHealthStatus: []*teov20220901.OriginHealthStatus{
						{
							Origin:  ptrStringOriginGroupHealthStatusDS("1.1.1.1"),
							Healthy: ptrStringOriginGroupHealthStatusDS("Healthy"),
						},
					},
					CheckRegionHealthStatus: []*teov20220901.CheckRegionHealthStatus{
						{
							Region:  ptrStringOriginGroupHealthStatusDS("US"),
							Healthy: ptrStringOriginGroupHealthStatusDS("Undetected"),
							OriginHealthStatus: []*teov20220901.OriginHealthStatus{
								{
									Origin:  ptrStringOriginGroupHealthStatusDS("1.1.1.1"),
									Healthy: ptrStringOriginGroupHealthStatusDS("Undetected"),
								},
							},
						},
					},
				},
			},
		}
		return resp, nil
	})

	meta := newMockMetaOriginGroupHealthStatusDS()
	res := teo.DataSourceTencentCloudTeoOriginGroupHealthStatus()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":          "zone-12345678",
		"lb_instance_id":   "lb-instance-xxx",
		"origin_group_ids": []interface{}{"origin-group-xxx1", "origin-group-xxx2"},
	})

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.NotEmpty(t, d.Id())

	healthStatusList := d.Get("origin_group_health_status_list").([]interface{})
	assert.Len(t, healthStatusList, 1)
	healthStatusMap := healthStatusList[0].(map[string]interface{})
	assert.Equal(t, "origin-group-xxx1", healthStatusMap["origin_group_id"].(string))
}

// TestTeoOriginGroupHealthStatusDS_ReadWithEmptyResult tests data source Read when API returns empty list
func TestTeoOriginGroupHealthStatusDS_ReadWithEmptyResult(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaOriginGroupHealthStatusDS().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeOriginGroupHealthStatusWithContext", func(_ context.Context, request *teov20220901.DescribeOriginGroupHealthStatusRequest) (*teov20220901.DescribeOriginGroupHealthStatusResponse, error) {
		resp := teov20220901.NewDescribeOriginGroupHealthStatusResponse()
		resp.Response = &teov20220901.DescribeOriginGroupHealthStatusResponseParams{
			OriginGroupHealthStatusList: []*teov20220901.OriginGroupHealthStatusDetail{},
		}
		return resp, nil
	})

	meta := newMockMetaOriginGroupHealthStatusDS()
	res := teo.DataSourceTencentCloudTeoOriginGroupHealthStatus()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":        "zone-12345678",
		"lb_instance_id": "lb-instance-xxx",
	})

	err := res.Read(d, meta)
	assert.Error(t, err)
	assert.Empty(t, d.Id())
}

// TestTeoOriginGroupHealthStatusDS_Schema tests the data source schema definition
func TestTeoOriginGroupHealthStatusDS_Schema(t *testing.T) {
	res := teo.DataSourceTencentCloudTeoOriginGroupHealthStatus()

	assert.NotNil(t, res)
	assert.Contains(t, res.Schema, "zone_id")
	assert.Contains(t, res.Schema, "lb_instance_id")
	assert.Contains(t, res.Schema, "origin_group_ids")
	assert.Contains(t, res.Schema, "origin_group_health_status_list")
	assert.Contains(t, res.Schema, "result_output_file")

	zoneIdSchema := res.Schema["zone_id"]
	assert.Equal(t, schema.TypeString, zoneIdSchema.Type)
	assert.True(t, zoneIdSchema.Required)

	lbInstanceIdSchema := res.Schema["lb_instance_id"]
	assert.Equal(t, schema.TypeString, lbInstanceIdSchema.Type)
	assert.True(t, lbInstanceIdSchema.Required)

	originGroupIdsSchema := res.Schema["origin_group_ids"]
	assert.Equal(t, schema.TypeList, originGroupIdsSchema.Type)
	assert.True(t, originGroupIdsSchema.Optional)

	healthStatusListSchema := res.Schema["origin_group_health_status_list"]
	assert.Equal(t, schema.TypeList, healthStatusListSchema.Type)
	assert.True(t, healthStatusListSchema.Computed)

	elemRes := healthStatusListSchema.Elem.(*schema.Resource)
	assert.Contains(t, elemRes.Schema, "origin_group_id")
	assert.Contains(t, elemRes.Schema, "origin_health_status")
	assert.Contains(t, elemRes.Schema, "check_region_health_status")
}
