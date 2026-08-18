package cdb_test

import (
	"context"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
	cdb_sdk "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cdb"
)

type mockMetaForCdbCloneInstance struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaForCdbCloneInstance) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaForCdbCloneInstance{}

func newMockMetaForCdbCloneInstance() *mockMetaForCdbCloneInstance {
	return &mockMetaForCdbCloneInstance{client: &connectivity.TencentCloudClient{}}
}

func ptrStringClone(s string) *string { return &s }
func ptrInt64Clone(v int64) *int64    { return &v }

func patchAsyncRequestSuccessClone(patches *gomonkey.Patches, cdbClient *cdb_sdk.Client) {
	patches.ApplyMethodFunc(cdbClient, "DescribeAsyncRequestInfo", func(request *cdb_sdk.DescribeAsyncRequestInfoRequest) (*cdb_sdk.DescribeAsyncRequestInfoResponse, error) {
		resp := cdb_sdk.NewDescribeAsyncRequestInfoResponse()
		resp.Response = &cdb_sdk.DescribeAsyncRequestInfoResponseParams{
			Status:    ptrStringClone("SUCCESS"),
			Info:      ptrStringClone("ok"),
			RequestId: ptrStringClone("fake-request-id"),
		}
		return resp, nil
	})
}

func patchDescribeDBInstancesClone(patches *gomonkey.Patches, cdbClient *cdb_sdk.Client, instanceId string, status int64, memory, volume, cpu int64) {
	patches.ApplyMethodFunc(cdbClient, "DescribeDBInstances", func(request *cdb_sdk.DescribeDBInstancesRequest) (*cdb_sdk.DescribeDBInstancesResponse, error) {
		resp := cdb_sdk.NewDescribeDBInstancesResponse()
		resp.Response = &cdb_sdk.DescribeDBInstancesResponseParams{
			Items: []*cdb_sdk.InstanceInfo{
				{
					InstanceId:   ptrStringClone(instanceId),
					Status:       ptrInt64Clone(status),
					Memory:       ptrInt64Clone(memory),
					Volume:       ptrInt64Clone(volume),
					Cpu:          ptrInt64Clone(cpu),
					InstanceName: ptrStringClone("tf-clone-example"),
					ProtectMode:  ptrInt64Clone(0),
					DeployMode:   ptrInt64Clone(0),
					DeviceType:   ptrStringClone("UNIVERSAL"),
					Zone:         ptrStringClone("ap-guangzhou-3"),
					ProjectId:    ptrInt64Clone(0),
					UniqVpcId:    ptrStringClone("vpc-i5yyodl9"),
					UniqSubnetId: ptrStringClone("subnet-hhi88a58"),
					Vip:          ptrStringClone("10.0.0.1"),
					Vport:        ptrInt64Clone(3306),
				},
			},
			RequestId: ptrStringClone("fake-request-id"),
		}
		return resp, nil
	})
}

func TestCdbCloneInstance_Create_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	cdbClient := &cdb_sdk.Client{}
	patches.ApplyMethodReturn(newMockMetaForCdbCloneInstance().client, "UseMysqlClient", cdbClient)

	patches.ApplyMethodFunc(cdbClient, "CreateCloneInstanceWithContext", func(_ context.Context, request *cdb_sdk.CreateCloneInstanceRequest) (*cdb_sdk.CreateCloneInstanceResponse, error) {
		resp := cdb_sdk.NewCreateCloneInstanceResponse()
		resp.Response = &cdb_sdk.CreateCloneInstanceResponseParams{
			AsyncRequestId: ptrStringClone("async-request-id-123"),
			RequestId:      ptrStringClone("fake-request-id"),
		}
		return resp, nil
	})

	patchAsyncRequestSuccessClone(patches, cdbClient)

	patches.ApplyMethodFunc(cdbClient, "DescribeCloneList", func(request *cdb_sdk.DescribeCloneListRequest) (*cdb_sdk.DescribeCloneListResponse, error) {
		resp := cdb_sdk.NewDescribeCloneListResponse()
		resp.Response = &cdb_sdk.DescribeCloneListResponseParams{
			TotalCount: ptrInt64Clone(1),
			Items: []*cdb_sdk.CloneItem{
				{
					SrcInstanceId:    ptrStringClone("cdb-source1234"),
					DstInstanceId:    ptrStringClone("cdb-clone1234"),
					CloneJobId:       ptrInt64Clone(1001),
					RollbackStrategy: ptrStringClone("timepoint"),
					TaskStatus:       ptrStringClone("success"),
				},
			},
			RequestId: ptrStringClone("fake-request-id"),
		}
		return resp, nil
	})

	patchDescribeDBInstancesClone(patches, cdbClient, "cdb-clone1234", 1, 4000, 200, 2)

	meta := newMockMetaForCdbCloneInstance()
	res := cdb.ResourceTencentCloudCdbCloneInstance()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_id":             "cdb-source1234",
		"specified_rollback_time": "2024-01-01 12:00:00",
		"memory":                  4000,
		"volume":                  200,
		"cpu":                     2,
		"instance_name":           "tf-clone-example",
		"uniq_vpc_id":             "vpc-i5yyodl9",
		"uniq_subnet_id":          "subnet-hhi88a58",
		"device_type":             "UNIVERSAL",
		"pay_type":                "USED_PAID",
		"zone":                    "ap-guangzhou-3",
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "cdb-clone1234", d.Id())
	assert.Equal(t, "async-request-id-123", d.Get("async_request_id"))
}

func TestCdbCloneInstance_Read_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	cdbClient := &cdb_sdk.Client{}
	patches.ApplyMethodReturn(newMockMetaForCdbCloneInstance().client, "UseMysqlClient", cdbClient)

	patchDescribeDBInstancesClone(patches, cdbClient, "cdb-clone1234", 1, 4000, 200, 2)

	meta := newMockMetaForCdbCloneInstance()
	res := cdb.ResourceTencentCloudCdbCloneInstance()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_id": "cdb-source1234",
	})
	d.SetId("cdb-clone1234")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "cdb-clone1234", d.Id())
	assert.Equal(t, int64(4000), d.Get("memory"))
	assert.Equal(t, int64(200), d.Get("volume"))
	assert.Equal(t, int64(2), d.Get("cpu"))
}

func TestCdbCloneInstance_Read_NotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	cdbClient := &cdb_sdk.Client{}
	patches.ApplyMethodReturn(newMockMetaForCdbCloneInstance().client, "UseMysqlClient", cdbClient)

	patches.ApplyMethodFunc(cdbClient, "DescribeDBInstances", func(request *cdb_sdk.DescribeDBInstancesRequest) (*cdb_sdk.DescribeDBInstancesResponse, error) {
		resp := cdb_sdk.NewDescribeDBInstancesResponse()
		resp.Response = &cdb_sdk.DescribeDBInstancesResponseParams{
			Items:     []*cdb_sdk.InstanceInfo{},
			RequestId: ptrStringClone("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForCdbCloneInstance()
	res := cdb.ResourceTencentCloudCdbCloneInstance()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_id": "cdb-source1234",
	})
	d.SetId("cdb-clone1234")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

func TestCdbCloneInstance_Update_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	cdbClient := &cdb_sdk.Client{}
	patches.ApplyMethodReturn(newMockMetaForCdbCloneInstance().client, "UseMysqlClient", cdbClient)

	patches.ApplyMethodFunc(cdbClient, "UpgradeDBInstanceWithContext", func(_ context.Context, request *cdb_sdk.UpgradeDBInstanceRequest) (*cdb_sdk.UpgradeDBInstanceResponse, error) {
		resp := cdb_sdk.NewUpgradeDBInstanceResponse()
		resp.Response = &cdb_sdk.UpgradeDBInstanceResponseParams{
			AsyncRequestId: ptrStringClone("async-upgrade-123"),
			RequestId:      ptrStringClone("fake-request-id"),
		}
		return resp, nil
	})

	patchAsyncRequestSuccessClone(patches, cdbClient)
	patchDescribeDBInstancesClone(patches, cdbClient, "cdb-clone1234", 1, 8000, 400, 4)

	meta := newMockMetaForCdbCloneInstance()
	res := cdb.ResourceTencentCloudCdbCloneInstance()

	state := &terraform.InstanceState{
		ID: "cdb-clone1234",
		Attributes: map[string]string{
			"id":               "cdb-clone1234",
			"instance_id":      "cdb-source1234",
			"memory":           "4000",
			"volume":           "200",
			"cpu":              "2",
			"protect_mode":     "0",
			"deploy_mode":      "0",
			"slave_zone":       "ap-guangzhou-3",
			"device_type":      "UNIVERSAL",
			"async_request_id": "async-request-id-123",
		},
	}

	rawConfig := terraform.NewResourceConfigRaw(map[string]interface{}{
		"instance_id":  "cdb-source1234",
		"memory":       8000,
		"volume":       400,
		"cpu":          4,
		"protect_mode": 0,
		"deploy_mode":  0,
		"slave_zone":   "ap-guangzhou-3",
		"device_type":  "UNIVERSAL",
	})

	diff, err := res.Diff(nil, state, rawConfig, newMockMetaForCdbCloneInstance())
	assert.NoError(t, err)
	assert.NotNil(t, diff)

	d, err := schema.InternalMap(res.Schema).Data(state, diff)
	assert.NoError(t, err)

	updateErr := res.Update(d, newMockMetaForCdbCloneInstance())
	assert.NoError(t, updateErr)
	assert.Equal(t, "cdb-clone1234", d.Id())
	assert.Equal(t, "async-upgrade-123", d.Get("async_request_id"))
}

func TestCdbCloneInstance_Update_ImmutableChange(t *testing.T) {
	res := cdb.ResourceTencentCloudCdbCloneInstance()

	state := &terraform.InstanceState{
		ID: "cdb-clone1234",
		Attributes: map[string]string{
			"id":            "cdb-clone1234",
			"instance_id":   "cdb-source1234",
			"instance_name": "old-name",
			"memory":        "4000",
			"volume":        "200",
		},
	}

	rawConfig := terraform.NewResourceConfigRaw(map[string]interface{}{
		"instance_id":   "cdb-source1234",
		"instance_name": "new-name",
		"memory":        4000,
		"volume":        200,
	})

	diff, err := res.Diff(nil, state, rawConfig, newMockMetaForCdbCloneInstance())
	assert.NoError(t, err)
	assert.NotNil(t, diff)

	d, err := schema.InternalMap(res.Schema).Data(state, diff)
	assert.NoError(t, err)

	updateErr := res.Update(d, newMockMetaForCdbCloneInstance())
	assert.Error(t, updateErr)
	assert.Contains(t, updateErr.Error(), "cannot be modified")
}

func TestCdbCloneInstance_Delete_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	cdbClient := &cdb_sdk.Client{}
	patches.ApplyMethodReturn(newMockMetaForCdbCloneInstance().client, "UseMysqlClient", cdbClient)

	patches.ApplyMethodFunc(cdbClient, "IsolateDBInstance", func(request *cdb_sdk.IsolateDBInstanceRequest) (*cdb_sdk.IsolateDBInstanceResponse, error) {
		resp := cdb_sdk.NewIsolateDBInstanceResponse()
		resp.Response = &cdb_sdk.IsolateDBInstanceResponseParams{
			RequestId: ptrStringClone("fake-request-id"),
		}
		return resp, nil
	})

	statusCallCount := 0
	patches.ApplyMethodFunc(cdbClient, "DescribeDBInstances", func(request *cdb_sdk.DescribeDBInstancesRequest) (*cdb_sdk.DescribeDBInstancesResponse, error) {
		resp := cdb_sdk.NewDescribeDBInstancesResponse()
		if statusCallCount == 0 {
			resp.Response = &cdb_sdk.DescribeDBInstancesResponseParams{
				Items: []*cdb_sdk.InstanceInfo{
					{
						InstanceId: ptrStringClone("cdb-clone1234"),
						Status:     ptrInt64Clone(5),
					},
				},
				RequestId: ptrStringClone("fake-request-id"),
			}
		} else {
			resp.Response = &cdb_sdk.DescribeDBInstancesResponseParams{
				Items:     []*cdb_sdk.InstanceInfo{},
				RequestId: ptrStringClone("fake-request-id"),
			}
		}
		statusCallCount++
		return resp, nil
	})

	patches.ApplyMethodFunc(cdbClient, "OfflineIsolatedInstances", func(request *cdb_sdk.OfflineIsolatedInstancesRequest) (*cdb_sdk.OfflineIsolatedInstancesResponse, error) {
		resp := cdb_sdk.NewOfflineIsolatedInstancesResponse()
		resp.Response = &cdb_sdk.OfflineIsolatedInstancesResponseParams{
			RequestId: ptrStringClone("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForCdbCloneInstance()
	res := cdb.ResourceTencentCloudCdbCloneInstance()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_id": "cdb-source1234",
	})
	d.SetId("cdb-clone1234")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
}

func TestCdbCloneInstance_Delete_AlreadyGone(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	cdbClient := &cdb_sdk.Client{}
	patches.ApplyMethodReturn(newMockMetaForCdbCloneInstance().client, "UseMysqlClient", cdbClient)

	patches.ApplyMethodFunc(cdbClient, "IsolateDBInstance", func(request *cdb_sdk.IsolateDBInstanceRequest) (*cdb_sdk.IsolateDBInstanceResponse, error) {
		resp := cdb_sdk.NewIsolateDBInstanceResponse()
		resp.Response = &cdb_sdk.IsolateDBInstanceResponseParams{
			RequestId: ptrStringClone("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(cdbClient, "DescribeDBInstances", func(request *cdb_sdk.DescribeDBInstancesRequest) (*cdb_sdk.DescribeDBInstancesResponse, error) {
		resp := cdb_sdk.NewDescribeDBInstancesResponse()
		resp.Response = &cdb_sdk.DescribeDBInstancesResponseParams{
			Items:     []*cdb_sdk.InstanceInfo{},
			RequestId: ptrStringClone("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForCdbCloneInstance()
	res := cdb.ResourceTencentCloudCdbCloneInstance()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_id": "cdb-source1234",
	})
	d.SetId("cdb-clone1234")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
}
