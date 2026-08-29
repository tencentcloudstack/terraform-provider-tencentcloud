package tdmysql_test

import (
	"context"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	tdmysqlv20211122 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tdmysql/v20211122"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	svctdmysql "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tdmysql"
)

// mockMetaTdmysqlDbInstance implements tccommon.ProviderMeta
type mockMetaTdmysqlDbInstance struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaTdmysqlDbInstance) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaTdmysqlDbInstance{}

func ptrTdmysqlInt64(i int64) *int64 {
	return &i
}

func ptrTdmysqlString(s string) *string {
	return &s
}

// go test ./tencentcloud/services/tdmysql/ -run "TestTdmysqlDbInstanceCreate" -v -count=1 -gcflags="all=-l"

func TestTdmysqlDbInstanceCreate(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	tdmysqlClient := &tdmysqlv20211122.Client{}
	mockClient := &connectivity.TencentCloudClient{}
	patches.ApplyMethodReturn(mockClient, "UseTdmysqlV20211122Client", tdmysqlClient)

	var capturedCreateRequest *tdmysqlv20211122.CreateDBInstancesRequest
	patches.ApplyMethodFunc(tdmysqlClient, "CreateDBInstancesWithContext", func(_ context.Context, request *tdmysqlv20211122.CreateDBInstancesRequest) (*tdmysqlv20211122.CreateDBInstancesResponse, error) {
		capturedCreateRequest = request
		resp := tdmysqlv20211122.NewCreateDBInstancesResponse()
		resp.Response = &tdmysqlv20211122.CreateDBInstancesResponseParams{
			InstanceIds: []*string{ptrTdmysqlString("tdmysqldb-test-instance-id")},
			FlowId:      ptrTdmysqlInt64(1001),
			RequestId:   ptrTdmysqlString("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(tdmysqlClient, "DescribeFlowWithContext", func(_ context.Context, request *tdmysqlv20211122.DescribeFlowRequest) (*tdmysqlv20211122.DescribeFlowResponse, error) {
		resp := tdmysqlv20211122.NewDescribeFlowResponse()
		resp.Response = &tdmysqlv20211122.DescribeFlowResponseParams{
			Status:    ptrTdmysqlString("success"),
			RequestId: ptrTdmysqlString("fake-request-id"),
		}
		return resp, nil
	})

	// Patch DescribeDBInstanceDetailWithContext for the Read call after Create
	patches.ApplyMethodFunc(tdmysqlClient, "DescribeDBInstanceDetailWithContext", func(_ context.Context, request *tdmysqlv20211122.DescribeDBInstanceDetailRequest) (*tdmysqlv20211122.DescribeDBInstanceDetailResponse, error) {
		resp := tdmysqlv20211122.NewDescribeDBInstanceDetailResponse()
		resp.Response = &tdmysqlv20211122.DescribeDBInstanceDetailResponseParams{
			InstanceId:     ptrTdmysqlString("tdmysqldb-test-instance-id"),
			InstanceName:   ptrTdmysqlString("tf-example"),
			Zone:           ptrTdmysqlString("ap-guangzhou-3"),
			VpcId:          ptrTdmysqlString("vpc-xxxxxxxx"),
			SubnetId:       ptrTdmysqlString("subnet-xxxxxxxx"),
			SpecCode:       ptrTdmysqlString("TDSQL-C-LS001"),
			Disk:           ptrTdmysqlInt64(200),
			StorageNodeNum: ptrTdmysqlInt64(2),
			Replications:   ptrTdmysqlInt64(3),
			Status:         ptrTdmysqlString("running"),
			Vip:            ptrTdmysqlString("10.0.1.10"),
			Vport:          ptrTdmysqlInt64(3306),
			RequestId:      ptrTdmysqlString("fake-request-id"),
		}
		return resp, nil
	})

	meta := &mockMetaTdmysqlDbInstance{client: mockClient}
	res := svctdmysql.ResourceTencentCloudTdmysqlDbInstance()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone":             "ap-guangzhou-3",
		"vpc_id":           "vpc-xxxxxxxx",
		"subnet_id":        "subnet-xxxxxxxx",
		"spec_code":        "TDSQL-C-LS001",
		"disk":             200,
		"storage_node_num": 2,
		"replications":     3,
		"instance_name":    "tf-example",
		"instance_count":   1,
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "tdmysqldb-test-instance-id", d.Id())
	assert.NotNil(t, capturedCreateRequest)
	assert.Equal(t, "ap-guangzhou-3", *capturedCreateRequest.Zone)
	assert.Equal(t, "vpc-xxxxxxxx", *capturedCreateRequest.VpcId)
	assert.Equal(t, "subnet-xxxxxxxx", *capturedCreateRequest.SubnetId)
	assert.Equal(t, "TDSQL-C-LS001", *capturedCreateRequest.SpecCode)
	assert.Equal(t, int64(200), *capturedCreateRequest.Disk)
	assert.Equal(t, int64(2), *capturedCreateRequest.StorageNodeNum)
	assert.Equal(t, int64(3), *capturedCreateRequest.Replications)
	assert.Equal(t, int64(1), *capturedCreateRequest.InstanceCount)
	assert.Equal(t, "tf-example", *capturedCreateRequest.InstanceName)
	assert.Equal(t, 1001, d.Get("flow_id").(int))
	instanceIds := d.Get("instance_ids").([]interface{})
	assert.Len(t, instanceIds, 1)
	assert.Equal(t, "tdmysqldb-test-instance-id", instanceIds[0].(string))
}

func TestTdmysqlDbInstanceRead(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	tdmysqlClient := &tdmysqlv20211122.Client{}
	mockClient := &connectivity.TencentCloudClient{}
	patches.ApplyMethodReturn(mockClient, "UseTdmysqlV20211122Client", tdmysqlClient)

	patches.ApplyMethodFunc(tdmysqlClient, "DescribeDBInstanceDetailWithContext", func(_ context.Context, request *tdmysqlv20211122.DescribeDBInstanceDetailRequest) (*tdmysqlv20211122.DescribeDBInstanceDetailResponse, error) {
		assert.Equal(t, "tdmysqldb-test-instance-id", *request.InstanceId)
		resp := tdmysqlv20211122.NewDescribeDBInstanceDetailResponse()
		resp.Response = &tdmysqlv20211122.DescribeDBInstanceDetailResponseParams{
			InstanceId:       ptrTdmysqlString("tdmysqldb-test-instance-id"),
			InstanceName:     ptrTdmysqlString("tf-example"),
			Zone:             ptrTdmysqlString("ap-guangzhou-3"),
			VpcId:            ptrTdmysqlString("vpc-xxxxxxxx"),
			SubnetId:         ptrTdmysqlString("subnet-xxxxxxxx"),
			SpecCode:         ptrTdmysqlString("TDSQL-C-LS001"),
			Disk:             ptrTdmysqlInt64(200),
			StorageNodeNum:   ptrTdmysqlInt64(2),
			Replications:     ptrTdmysqlInt64(3),
			Status:           ptrTdmysqlString("running"),
			Vip:              ptrTdmysqlString("10.0.1.10"),
			Vport:            ptrTdmysqlInt64(3306),
			CharSet:          ptrTdmysqlString("utf8"),
			Region:           ptrTdmysqlString("ap-guangzhou"),
			StatusDesc:       ptrTdmysqlString("运行中"),
			EncryptionEnable: ptrTdmysqlInt64(0),
			RequestId:        ptrTdmysqlString("fake-request-id"),
		}
		return resp, nil
	})

	meta := &mockMetaTdmysqlDbInstance{client: mockClient}
	res := svctdmysql.ResourceTencentCloudTdmysqlDbInstance()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone":             "ap-guangzhou-3",
		"vpc_id":           "vpc-xxxxxxxx",
		"subnet_id":        "subnet-xxxxxxxx",
		"spec_code":        "TDSQL-C-LS001",
		"disk":             200,
		"storage_node_num": 2,
		"replications":     3,
		"instance_name":    "tf-example",
	})
	d.SetId("tdmysqldb-test-instance-id")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "tdmysqldb-test-instance-id", d.Id())
	assert.Equal(t, "running", d.Get("status").(string))
	assert.Equal(t, "10.0.1.10", d.Get("vip").(string))
	assert.Equal(t, 3306, d.Get("vport").(int))
	assert.Equal(t, "utf8", d.Get("char_set").(string))
	assert.Equal(t, "ap-guangzhou", d.Get("region").(string))
	assert.Equal(t, "运行中", d.Get("status_desc").(string))
	assert.Equal(t, "tdmysqldb-test-instance-id", d.Get("instance_id").(string))
}

func TestTdmysqlDbInstanceReadNotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	tdmysqlClient := &tdmysqlv20211122.Client{}
	mockClient := &connectivity.TencentCloudClient{}
	patches.ApplyMethodReturn(mockClient, "UseTdmysqlV20211122Client", tdmysqlClient)

	patches.ApplyMethodFunc(tdmysqlClient, "DescribeDBInstanceDetailWithContext", func(_ context.Context, request *tdmysqlv20211122.DescribeDBInstanceDetailRequest) (*tdmysqlv20211122.DescribeDBInstanceDetailResponse, error) {
		resp := tdmysqlv20211122.NewDescribeDBInstanceDetailResponse()
		resp.Response = nil
		return resp, nil
	})

	meta := &mockMetaTdmysqlDbInstance{client: mockClient}
	res := svctdmysql.ResourceTencentCloudTdmysqlDbInstance()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone":             "ap-guangzhou-3",
		"vpc_id":           "vpc-xxxxxxxx",
		"subnet_id":        "subnet-xxxxxxxx",
		"spec_code":        "TDSQL-C-LS001",
		"disk":             200,
		"storage_node_num": 2,
		"replications":     3,
		"instance_name":    "tf-example",
	})
	d.SetId("tdmysqldb-test-instance-id")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

func TestTdmysqlDbInstanceUpdate(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	tdmysqlClient := &tdmysqlv20211122.Client{}
	mockClient := &connectivity.TencentCloudClient{}
	patches.ApplyMethodReturn(mockClient, "UseTdmysqlV20211122Client", tdmysqlClient)

	var capturedRequest *tdmysqlv20211122.ModifyInstanceNameRequest
	patches.ApplyMethodFunc(tdmysqlClient, "ModifyInstanceNameWithContext", func(_ context.Context, request *tdmysqlv20211122.ModifyInstanceNameRequest) (*tdmysqlv20211122.ModifyInstanceNameResponse, error) {
		capturedRequest = request
		resp := tdmysqlv20211122.NewModifyInstanceNameResponse()
		resp.Response = &tdmysqlv20211122.ModifyInstanceNameResponseParams{
			RequestId: ptrTdmysqlString("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(tdmysqlClient, "DescribeDBInstanceDetailWithContext", func(_ context.Context, request *tdmysqlv20211122.DescribeDBInstanceDetailRequest) (*tdmysqlv20211122.DescribeDBInstanceDetailResponse, error) {
		resp := tdmysqlv20211122.NewDescribeDBInstanceDetailResponse()
		resp.Response = &tdmysqlv20211122.DescribeDBInstanceDetailResponseParams{
			InstanceId:     ptrTdmysqlString("tdmysqldb-test-instance-id"),
			InstanceName:   ptrTdmysqlString("tf-example-update"),
			Zone:           ptrTdmysqlString("ap-guangzhou-3"),
			VpcId:          ptrTdmysqlString("vpc-xxxxxxxx"),
			SubnetId:       ptrTdmysqlString("subnet-xxxxxxxx"),
			SpecCode:       ptrTdmysqlString("TDSQL-C-LS001"),
			Disk:           ptrTdmysqlInt64(200),
			StorageNodeNum: ptrTdmysqlInt64(2),
			Replications:   ptrTdmysqlInt64(3),
			Status:         ptrTdmysqlString("running"),
			RequestId:      ptrTdmysqlString("fake-request-id"),
		}
		return resp, nil
	})

	meta := &mockMetaTdmysqlDbInstance{client: mockClient}
	res := svctdmysql.ResourceTencentCloudTdmysqlDbInstance()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone":             "ap-guangzhou-3",
		"vpc_id":           "vpc-xxxxxxxx",
		"subnet_id":        "subnet-xxxxxxxx",
		"spec_code":        "TDSQL-C-LS001",
		"disk":             200,
		"storage_node_num": 2,
		"replications":     3,
		"instance_name":    "tf-example-update",
	})
	d.SetId("tdmysqldb-test-instance-id")

	// Patch HasChange to simulate only instance_name change
	patches.ApplyMethodFunc(d, "HasChange", func(key string) bool {
		return key == "instance_name"
	})

	err := res.Update(d, meta)
	assert.NoError(t, err)
	assert.NotNil(t, capturedRequest)
	assert.Equal(t, "tdmysqldb-test-instance-id", *capturedRequest.InstanceId)
	assert.Equal(t, "tf-example-update", *capturedRequest.InstanceName)
}

func TestTdmysqlDbInstanceUpdateImmutable(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	tdmysqlClient := &tdmysqlv20211122.Client{}
	mockClient := &connectivity.TencentCloudClient{}
	patches.ApplyMethodReturn(mockClient, "UseTdmysqlV20211122Client", tdmysqlClient)

	meta := &mockMetaTdmysqlDbInstance{client: mockClient}
	res := svctdmysql.ResourceTencentCloudTdmysqlDbInstance()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone":             "ap-guangzhou-3",
		"vpc_id":           "vpc-xxxxxxxx",
		"subnet_id":        "subnet-xxxxxxxx",
		"spec_code":        "TDSQL-C-LS001",
		"disk":             200,
		"storage_node_num": 2,
		"replications":     3,
		"instance_name":    "tf-example",
	})
	d.SetId("tdmysqldb-test-instance-id")

	// Patch HasChange to simulate an immutable field (disk) change
	patches.ApplyMethodFunc(d, "HasChange", func(key string) bool {
		return key == "disk"
	})

	err := res.Update(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "disk")
}

func TestTdmysqlDbInstanceDelete(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	tdmysqlClient := &tdmysqlv20211122.Client{}
	mockClient := &connectivity.TencentCloudClient{}
	patches.ApplyMethodReturn(mockClient, "UseTdmysqlV20211122Client", tdmysqlClient)

	var capturedRequest *tdmysqlv20211122.IsolateDBInstanceRequest
	patches.ApplyMethodFunc(tdmysqlClient, "IsolateDBInstanceWithContext", func(_ context.Context, request *tdmysqlv20211122.IsolateDBInstanceRequest) (*tdmysqlv20211122.IsolateDBInstanceResponse, error) {
		capturedRequest = request
		resp := tdmysqlv20211122.NewIsolateDBInstanceResponse()
		resp.Response = &tdmysqlv20211122.IsolateDBInstanceResponseParams{
			SuccessInstanceIds: []*string{ptrTdmysqlString("tdmysqldb-test-instance-id")},
			RequestId:          ptrTdmysqlString("fake-request-id"),
		}
		return resp, nil
	})

	meta := &mockMetaTdmysqlDbInstance{client: mockClient}
	res := svctdmysql.ResourceTencentCloudTdmysqlDbInstance()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone":             "ap-guangzhou-3",
		"vpc_id":           "vpc-xxxxxxxx",
		"subnet_id":        "subnet-xxxxxxxx",
		"spec_code":        "TDSQL-C-LS001",
		"disk":             200,
		"storage_node_num": 2,
		"replications":     3,
		"instance_name":    "tf-example",
	})
	d.SetId("tdmysqldb-test-instance-id")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
	assert.NotNil(t, capturedRequest)
	assert.Len(t, capturedRequest.InstanceIds, 1)
	assert.Equal(t, "tdmysqldb-test-instance-id", *capturedRequest.InstanceIds[0])
}

func TestTdmysqlDbInstanceSchema(t *testing.T) {
	res := svctdmysql.ResourceTencentCloudTdmysqlDbInstance()

	assert.NotNil(t, res)
	assert.Contains(t, res.Schema, "zone")
	assert.Contains(t, res.Schema, "vpc_id")
	assert.Contains(t, res.Schema, "subnet_id")
	assert.Contains(t, res.Schema, "spec_code")
	assert.Contains(t, res.Schema, "disk")
	assert.Contains(t, res.Schema, "storage_node_num")
	assert.Contains(t, res.Schema, "replications")
	assert.Contains(t, res.Schema, "instance_name")
	assert.Contains(t, res.Schema, "instance_id")
	assert.Contains(t, res.Schema, "flow_id")
	assert.Contains(t, res.Schema, "instance_ids")

	// instance_name should be Required and not ForceNew
	instanceName := res.Schema["instance_name"]
	assert.Equal(t, schema.TypeString, instanceName.Type)
	assert.True(t, instanceName.Required)
	assert.False(t, instanceName.ForceNew)

	// zone should be Required and ForceNew
	zone := res.Schema["zone"]
	assert.True(t, zone.Required)
	assert.True(t, zone.ForceNew)

	// instance_id should be Computed
	instanceId := res.Schema["instance_id"]
	assert.True(t, instanceId.Computed)

	// flow_id should be Computed
	flowId := res.Schema["flow_id"]
	assert.True(t, flowId.Computed)

	// Importer should be set
	assert.NotNil(t, res.Importer)
}
