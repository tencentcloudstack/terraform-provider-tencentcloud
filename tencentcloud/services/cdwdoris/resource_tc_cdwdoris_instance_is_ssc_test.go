package cdwdoris_test

import (
	"context"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	cdwdorisv20211228 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdwdoris/v20211228"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	cdwdoris "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cdwdoris"
)

type mockMetaCdwdoris struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaCdwdoris) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaCdwdoris{}

func newMockMetaCdwdoris() *mockMetaCdwdoris {
	return &mockMetaCdwdoris{client: &connectivity.TencentCloudClient{}}
}

func ptrStringCdwdoris(s string) *string {
	return &s
}

func ptrInt64Cdwdoris(i int64) *int64 {
	return &i
}

func TestAccCdwdorisInstance_CreateWithIsSsc(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	cdwdorisClient := &cdwdorisv20211228.Client{}
	patches.ApplyMethodReturn(newMockMetaCdwdoris().client, "UseCdwdorisV20211228Client", cdwdorisClient)

	// Mock CreateInstanceNewWithContext - verify IsSSC is passed
	patches.ApplyMethodFunc(cdwdorisClient, "CreateInstanceNewWithContext", func(_ context.Context, request *cdwdorisv20211228.CreateInstanceNewRequest) (*cdwdorisv20211228.CreateInstanceNewResponse, error) {
		assert.NotNil(t, request.IsSSC)
		assert.True(t, *request.IsSSC)
		instanceId := "cdwdoris-test-123"
		resp := cdwdorisv20211228.NewCreateInstanceNewResponse()
		resp.Response = &cdwdorisv20211228.CreateInstanceNewResponseParams{
			InstanceId: &instanceId,
			RequestId:  ptrStringCdwdoris("fake-request-id"),
		}
		return resp, nil
	})

	// Mock DescribeInstanceStateWithContext - for wait loop
	patches.ApplyMethodFunc(cdwdorisClient, "DescribeInstanceStateWithContext", func(_ context.Context, request *cdwdorisv20211228.DescribeInstanceStateRequest) (*cdwdorisv20211228.DescribeInstanceStateResponse, error) {
		state := "Serving"
		resp := cdwdorisv20211228.NewDescribeInstanceStateResponse()
		resp.Response = &cdwdorisv20211228.DescribeInstanceStateResponseParams{
			InstanceState: &state,
			RequestId:     ptrStringCdwdoris("fake-request-id"),
		}
		return resp, nil
	})

	// Mock CdwdorisService.DescribeCdwdorisInstanceById for Read
	patches.ApplyMethodFunc(&cdwdoris.CdwdorisService{}, "DescribeCdwdorisInstanceById", func(_ context.Context, instanceId string) (*cdwdorisv20211228.InstanceInfo, error) {
		assert.Equal(t, "cdwdoris-test-123", instanceId)
		return &cdwdorisv20211228.InstanceInfo{
			InstanceId:   ptrStringCdwdoris("cdwdoris-test-123"),
			InstanceName: ptrStringCdwdoris("tf-example"),
			Version:      ptrStringCdwdoris("2.1"),
			Zone:         ptrStringCdwdoris("ap-guangzhou-6"),
			VpcId:        ptrStringCdwdoris("vpc-test"),
			SubnetId:     ptrStringCdwdoris("subnet-test"),
			HA:           ptrStringCdwdoris("true"),
			HaType:       ptrInt64Cdwdoris(1),
		}, nil
	})

	// Mock CdwdorisService.DescribeCdwdorisWorkloadGroupsById for Read
	patches.ApplyMethodFunc(&cdwdoris.CdwdorisService{}, "DescribeCdwdorisWorkloadGroupsById", func(_ context.Context, instanceId string) (*cdwdorisv20211228.DescribeWorkloadGroupResponseParams, error) {
		assert.Equal(t, "cdwdoris-test-123", instanceId)
		return &cdwdorisv20211228.DescribeWorkloadGroupResponseParams{
			Status: ptrStringCdwdoris("close"),
		}, nil
	})

	meta := newMockMetaCdwdoris()
	res := cdwdoris.ResourceTencentCloudCdwdorisInstance()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone":                  "ap-guangzhou-6",
		"user_vpc_id":           "vpc-test",
		"user_subnet_id":        "subnet-test",
		"product_version":       "2.1",
		"instance_name":         "tf-example",
		"doris_user_pwd":        "Password@test",
		"ha_flag":               true,
		"ha_type":               1,
		"case_sensitive":        0,
		"enable_multi_zones":    false,
		"is_ssc":                true,
		"workload_group_status": "close",
		"charge_properties": []interface{}{
			map[string]interface{}{
				"charge_type": "POSTPAID_BY_HOUR",
			},
		},
		"fe_spec": []interface{}{
			map[string]interface{}{
				"spec_name": "S_4_16_P",
				"count":     3,
				"disk_size": 200,
			},
		},
		"be_spec": []interface{}{
			map[string]interface{}{
				"spec_name": "S_4_16_P",
				"count":     3,
				"disk_size": 200,
			},
		},
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "cdwdoris-test-123", d.Id())
}

func TestAccCdwdorisInstance_CreateWithoutIsSsc(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	cdwdorisClient := &cdwdorisv20211228.Client{}
	patches.ApplyMethodReturn(newMockMetaCdwdoris().client, "UseCdwdorisV20211228Client", cdwdorisClient)

	// Mock CreateInstanceNewWithContext - verify IsSSC is NOT set
	patches.ApplyMethodFunc(cdwdorisClient, "CreateInstanceNewWithContext", func(_ context.Context, request *cdwdorisv20211228.CreateInstanceNewRequest) (*cdwdorisv20211228.CreateInstanceNewResponse, error) {
		assert.Nil(t, request.IsSSC, "IsSSC should be nil when is_ssc is not set")
		instanceId := "cdwdoris-test-456"
		resp := cdwdorisv20211228.NewCreateInstanceNewResponse()
		resp.Response = &cdwdorisv20211228.CreateInstanceNewResponseParams{
			InstanceId: &instanceId,
			RequestId:  ptrStringCdwdoris("fake-request-id"),
		}
		return resp, nil
	})

	// Mock DescribeInstanceStateWithContext - for wait loop
	patches.ApplyMethodFunc(cdwdorisClient, "DescribeInstanceStateWithContext", func(_ context.Context, request *cdwdorisv20211228.DescribeInstanceStateRequest) (*cdwdorisv20211228.DescribeInstanceStateResponse, error) {
		state := "Serving"
		resp := cdwdorisv20211228.NewDescribeInstanceStateResponse()
		resp.Response = &cdwdorisv20211228.DescribeInstanceStateResponseParams{
			InstanceState: &state,
			RequestId:     ptrStringCdwdoris("fake-request-id"),
		}
		return resp, nil
	})

	// Mock CdwdorisService.DescribeCdwdorisInstanceById for Read
	patches.ApplyMethodFunc(&cdwdoris.CdwdorisService{}, "DescribeCdwdorisInstanceById", func(_ context.Context, instanceId string) (*cdwdorisv20211228.InstanceInfo, error) {
		return &cdwdorisv20211228.InstanceInfo{
			InstanceId:   ptrStringCdwdoris("cdwdoris-test-456"),
			InstanceName: ptrStringCdwdoris("tf-example"),
			Version:      ptrStringCdwdoris("2.1"),
			Zone:         ptrStringCdwdoris("ap-guangzhou-6"),
			VpcId:        ptrStringCdwdoris("vpc-test"),
			SubnetId:     ptrStringCdwdoris("subnet-test"),
			HA:           ptrStringCdwdoris("true"),
			HaType:       ptrInt64Cdwdoris(1),
		}, nil
	})

	// Mock CdwdorisService.DescribeCdwdorisWorkloadGroupsById for Read
	patches.ApplyMethodFunc(&cdwdoris.CdwdorisService{}, "DescribeCdwdorisWorkloadGroupsById", func(_ context.Context, instanceId string) (*cdwdorisv20211228.DescribeWorkloadGroupResponseParams, error) {
		return &cdwdorisv20211228.DescribeWorkloadGroupResponseParams{
			Status: ptrStringCdwdoris("close"),
		}, nil
	})

	meta := newMockMetaCdwdoris()
	res := cdwdoris.ResourceTencentCloudCdwdorisInstance()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone":                  "ap-guangzhou-6",
		"user_vpc_id":           "vpc-test",
		"user_subnet_id":        "subnet-test",
		"product_version":       "2.1",
		"instance_name":         "tf-example",
		"doris_user_pwd":        "Password@test",
		"ha_flag":               true,
		"ha_type":               1,
		"case_sensitive":        0,
		"enable_multi_zones":    false,
		"workload_group_status": "close",
		"charge_properties": []interface{}{
			map[string]interface{}{
				"charge_type": "POSTPAID_BY_HOUR",
			},
		},
		"fe_spec": []interface{}{
			map[string]interface{}{
				"spec_name": "S_4_16_P",
				"count":     3,
				"disk_size": 200,
			},
		},
		"be_spec": []interface{}{
			map[string]interface{}{
				"spec_name": "S_4_16_P",
				"count":     3,
				"disk_size": 200,
			},
		},
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "cdwdoris-test-456", d.Id())
}

func TestAccCdwdorisInstance_UpdateIsSscError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	meta := newMockMetaCdwdoris()
	res := cdwdoris.ResourceTencentCloudCdwdorisInstance()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone":                  "ap-guangzhou-6",
		"user_vpc_id":           "vpc-test",
		"user_subnet_id":        "subnet-test",
		"product_version":       "2.1",
		"instance_name":         "tf-example",
		"doris_user_pwd":        "Password@test",
		"ha_flag":               true,
		"ha_type":               1,
		"case_sensitive":        0,
		"enable_multi_zones":    false,
		"is_ssc":                true,
		"workload_group_status": "close",
		"charge_properties": []interface{}{
			map[string]interface{}{
				"charge_type": "POSTPAID_BY_HOUR",
			},
		},
		"fe_spec": []interface{}{
			map[string]interface{}{
				"spec_name": "S_4_16_P",
				"count":     3,
				"disk_size": 200,
			},
		},
		"be_spec": []interface{}{
			map[string]interface{}{
				"spec_name": "S_4_16_P",
				"count":     3,
				"disk_size": 200,
			},
		},
	})
	d.SetId("cdwdoris-test-789")

	// Mock HasChange to simulate that is_ssc is being changed
	patches.ApplyMethodFunc(d, "HasChange", func(key string) bool {
		return key == "is_ssc"
	})

	err := res.Update(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "argument `is_ssc` cannot be changed")
}