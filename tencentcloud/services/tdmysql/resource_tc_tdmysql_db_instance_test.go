package tdmysql_test

import (
	"context"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	tdmysql_sdk "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tdmysql/v20211122"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tdmysql"
)

type mockMetaForTdmysqlDbInstance struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaForTdmysqlDbInstance) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaForTdmysqlDbInstance{}

func newMockMetaForTdmysqlDbInstance() *mockMetaForTdmysqlDbInstance {
	return &mockMetaForTdmysqlDbInstance{client: &connectivity.TencentCloudClient{}}
}

func ptrStringTdmysql(s string) *string  { return &s }
func ptrInt64Tdmysql(v int64) *int64     { return &v }
func ptrBoolTdmysql(b bool) *bool        { return &b }
func ptrFloat64Tdmysql(v float64) *float64 { return &v }

func TestTdmysqlDbInstance_Create_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	tdmysqlClient := &tdmysql_sdk.Client{}
	patches.ApplyMethodReturn(newMockMetaForTdmysqlDbInstance().client, "UseTdmysqlV20211122Client", tdmysqlClient)

	patches.ApplyMethodFunc(tdmysqlClient, "CreateDBInstancesWithContext", func(_ context.Context, request *tdmysql_sdk.CreateDBInstancesRequest) (*tdmysql_sdk.CreateDBInstancesResponse, error) {
		resp := tdmysql_sdk.NewCreateDBInstancesResponse()
		resp.Response = &tdmysql_sdk.CreateDBInstancesResponseParams{
			InstanceIds: []*string{ptrStringTdmysql("tdsqlshard-test1234")},
			FlowId:      ptrInt64Tdmysql(123456),
			RequestId:   ptrStringTdmysql("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(tdmysqlClient, "DescribeFlowWithContext", func(_ context.Context, request *tdmysql_sdk.DescribeFlowRequest) (*tdmysql_sdk.DescribeFlowResponse, error) {
		resp := tdmysql_sdk.NewDescribeFlowResponse()
		resp.Response = &tdmysql_sdk.DescribeFlowResponseParams{
			Status:    ptrStringTdmysql("success"),
			RequestId: ptrStringTdmysql("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(tdmysqlClient, "DescribeDBInstanceDetailWithContext", func(_ context.Context, request *tdmysql_sdk.DescribeDBInstanceDetailRequest) (*tdmysql_sdk.DescribeDBInstanceDetailResponse, error) {
		resp := tdmysql_sdk.NewDescribeDBInstanceDetailResponse()
		resp.Response = &tdmysql_sdk.DescribeDBInstanceDetailResponseParams{
			InstanceId:    ptrStringTdmysql("tdsqlshard-test1234"),
			InstanceName:  ptrStringTdmysql("tf-tdmysql-example"),
			Zone:          ptrStringTdmysql("ap-guangzhou-3"),
			VpcId:         ptrStringTdmysql("vpc-xxxxxxxx"),
			SubnetId:      ptrStringTdmysql("subnet-xxxxxxxx"),
			SpecCode:      ptrStringTdmysql("spec-code"),
			Disk:          ptrInt64Tdmysql(100),
			StorageNodeNum: ptrInt64Tdmysql(2),
			Replications:  ptrInt64Tdmysql(3),
			Status:        ptrStringTdmysql("running"),
			RequestId:     ptrStringTdmysql("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForTdmysqlDbInstance()
	res := tdmysql.ResourceTencentCloudTdmysqlDbInstance()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone":             "ap-guangzhou-3",
		"vpc_id":           "vpc-xxxxxxxx",
		"subnet_id":        "subnet-xxxxxxxx",
		"spec_code":        "spec-code",
		"disk":             100,
		"storage_node_num": 2,
		"replications":     3,
		"instance_name":    "tf-tdmysql-example",
		"pay_mode":         "0",
		"instance_type":    "separate",
		"storage_type":     "CLOUD_HSSD",
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "tdsqlshard-test1234", d.Id())
	assert.Equal(t, "tdsqlshard-test1234", d.Get("instance_id"))
	assert.Equal(t, 123456, d.Get("flow_id"))
}

func TestTdmysqlDbInstance_Read_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	tdmysqlClient := &tdmysql_sdk.Client{}
	patches.ApplyMethodReturn(newMockMetaForTdmysqlDbInstance().client, "UseTdmysqlV20211122Client", tdmysqlClient)

	patches.ApplyMethodFunc(tdmysqlClient, "DescribeDBInstanceDetailWithContext", func(_ context.Context, request *tdmysql_sdk.DescribeDBInstanceDetailRequest) (*tdmysql_sdk.DescribeDBInstanceDetailResponse, error) {
		resp := tdmysql_sdk.NewDescribeDBInstanceDetailResponse()
		resp.Response = &tdmysql_sdk.DescribeDBInstanceDetailResponseParams{
			InstanceId:     ptrStringTdmysql("tdsqlshard-test1234"),
			InstanceName:   ptrStringTdmysql("tf-tdmysql-example"),
			Zone:           ptrStringTdmysql("ap-guangzhou-3"),
			VpcId:          ptrStringTdmysql("vpc-xxxxxxxx"),
			SubnetId:       ptrStringTdmysql("subnet-xxxxxxxx"),
			SpecCode:       ptrStringTdmysql("spec-code"),
			Disk:           ptrInt64Tdmysql(100),
			StorageNodeNum: ptrInt64Tdmysql(2),
			Replications:   ptrInt64Tdmysql(3),
			Vip:            ptrStringTdmysql("10.0.0.1"),
			Vport:          ptrInt64Tdmysql(3306),
			Status:         ptrStringTdmysql("running"),
			StatusDesc:     ptrStringTdmysql("运行中"),
			Region:         ptrStringTdmysql("ap-guangzhou"),
			CharSet:        ptrStringTdmysql("utf8"),
			PayMode:        ptrStringTdmysql("0"),
			InstanceType:   ptrStringTdmysql("separate"),
			StorageType:    ptrStringTdmysql("CLOUD_HSSD"),
			RequestId:      ptrStringTdmysql("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForTdmysqlDbInstance()
	res := tdmysql.ResourceTencentCloudTdmysqlDbInstance()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
	d.SetId("tdsqlshard-test1234")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "tdsqlshard-test1234", d.Id())
	assert.Equal(t, "tdsqlshard-test1234", d.Get("instance_id"))
	assert.Equal(t, "tf-tdmysql-example", d.Get("instance_name"))
	assert.Equal(t, "ap-guangzhou-3", d.Get("zone"))
	assert.Equal(t, "vpc-xxxxxxxx", d.Get("vpc_id"))
	assert.Equal(t, "subnet-xxxxxxxx", d.Get("subnet_id"))
	assert.Equal(t, "spec-code", d.Get("spec_code"))
	assert.Equal(t, 100, d.Get("disk"))
	assert.Equal(t, 2, d.Get("storage_node_num"))
	assert.Equal(t, 3, d.Get("replications"))
	assert.Equal(t, "10.0.0.1", d.Get("vip"))
	assert.Equal(t, 3306, d.Get("vport"))
	assert.Equal(t, "running", d.Get("status"))
	assert.Equal(t, "运行中", d.Get("status_desc"))
	assert.Equal(t, "ap-guangzhou", d.Get("region"))
	assert.Equal(t, "utf8", d.Get("char_set"))
	assert.Equal(t, "0", d.Get("pay_mode"))
	assert.Equal(t, "separate", d.Get("instance_type"))
	assert.Equal(t, "CLOUD_HSSD", d.Get("storage_type"))
}

func TestTdmysqlDbInstance_Read_NotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	tdmysqlClient := &tdmysql_sdk.Client{}
	patches.ApplyMethodReturn(newMockMetaForTdmysqlDbInstance().client, "UseTdmysqlV20211122Client", tdmysqlClient)

	patches.ApplyMethodFunc(tdmysqlClient, "DescribeDBInstanceDetailWithContext", func(_ context.Context, request *tdmysql_sdk.DescribeDBInstanceDetailRequest) (*tdmysql_sdk.DescribeDBInstanceDetailResponse, error) {
		resp := tdmysql_sdk.NewDescribeDBInstanceDetailResponse()
		resp.Response = &tdmysql_sdk.DescribeDBInstanceDetailResponseParams{
			InstanceId: nil,
			RequestId:  ptrStringTdmysql("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForTdmysqlDbInstance()
	res := tdmysql.ResourceTencentCloudTdmysqlDbInstance()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
	d.SetId("tdsqlshard-test1234")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

func TestTdmysqlDbInstance_Update_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	tdmysqlClient := &tdmysql_sdk.Client{}
	patches.ApplyMethodReturn(newMockMetaForTdmysqlDbInstance().client, "UseTdmysqlV20211122Client", tdmysqlClient)

	patches.ApplyMethodFunc(tdmysqlClient, "ModifyInstanceNameWithContext", func(_ context.Context, request *tdmysql_sdk.ModifyInstanceNameRequest) (*tdmysql_sdk.ModifyInstanceNameResponse, error) {
		resp := tdmysql_sdk.NewModifyInstanceNameResponse()
		resp.Response = &tdmysql_sdk.ModifyInstanceNameResponseParams{
			RequestId: ptrStringTdmysql("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(tdmysqlClient, "DescribeDBInstanceDetailWithContext", func(_ context.Context, request *tdmysql_sdk.DescribeDBInstanceDetailRequest) (*tdmysql_sdk.DescribeDBInstanceDetailResponse, error) {
		resp := tdmysql_sdk.NewDescribeDBInstanceDetailResponse()
		resp.Response = &tdmysql_sdk.DescribeDBInstanceDetailResponseParams{
			InstanceId:   ptrStringTdmysql("tdsqlshard-test1234"),
			InstanceName: ptrStringTdmysql("tf-tdmysql-updated"),
			RequestId:    ptrStringTdmysql("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForTdmysqlDbInstance()
	res := tdmysql.ResourceTencentCloudTdmysqlDbInstance()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_name": "tf-tdmysql-updated",
	})
	d.SetId("tdsqlshard-test1234")

	err := res.Update(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "tf-tdmysql-updated", d.Get("instance_name"))
}

func TestTdmysqlDbInstance_Delete_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	tdmysqlClient := &tdmysql_sdk.Client{}
	patches.ApplyMethodReturn(newMockMetaForTdmysqlDbInstance().client, "UseTdmysqlV20211122Client", tdmysqlClient)

	patches.ApplyMethodFunc(tdmysqlClient, "IsolateDBInstanceWithContext", func(_ context.Context, request *tdmysql_sdk.IsolateDBInstanceRequest) (*tdmysql_sdk.IsolateDBInstanceResponse, error) {
		resp := tdmysql_sdk.NewIsolateDBInstanceResponse()
		resp.Response = &tdmysql_sdk.IsolateDBInstanceResponseParams{
			SuccessInstanceIds: []*string{ptrStringTdmysql("tdsqlshard-test1234")},
			RequestId:          ptrStringTdmysql("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForTdmysqlDbInstance()
	res := tdmysql.ResourceTencentCloudTdmysqlDbInstance()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
	d.SetId("tdsqlshard-test1234")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
}
