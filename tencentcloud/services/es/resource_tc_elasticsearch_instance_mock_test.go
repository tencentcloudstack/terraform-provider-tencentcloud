package es_test

import (
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	es "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/es/v20180416"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	svces "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/es"
)

type mockMetaEs struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaEs) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaEs{}

func newMockMetaEs() *mockMetaEs {
	return &mockMetaEs{client: &connectivity.TencentCloudClient{}}
}

func ptrStringEs(s string) *string {
	return &s
}

func ptrInt64Es(v int64) *int64 {
	return &v
}

// go test ./tencentcloud/services/es/ -run "TestEsInstanceDestroyProtection" -v -count=1 -gcflags="all=-l"

// TestEsInstanceDestroyProtection_Schema validates the enable_destroy_protection schema field definition
func TestEsInstanceDestroyProtection_Schema(t *testing.T) {
	res := svces.ResourceTencentCloudElasticsearchInstance()
	assert.NotNil(t, res)
	assert.Contains(t, res.Schema, "enable_destroy_protection")

	field := res.Schema["enable_destroy_protection"]
	assert.Equal(t, schema.TypeString, field.Type)
	assert.True(t, field.Optional)
	assert.True(t, field.Computed)
	assert.NotNil(t, field.ValidateFunc)
}

// TestEsInstanceDestroyProtection_Read_NonNil verifies the read flow sets enable_destroy_protection
// from InstanceInfo.EnableDestroyProtection when the API returns a non-nil value.
func TestEsInstanceDestroyProtection_Read_NonNil(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	esClient := &es.Client{}
	patches.ApplyMethodReturn(newMockMetaEs().client, "UseEsClient", esClient)

	patches.ApplyMethodFunc(esClient, "DescribeInstances", func(request *es.DescribeInstancesRequest) (*es.DescribeInstancesResponse, error) {
		resp := es.NewDescribeInstancesResponse()
		resp.Response = &es.DescribeInstancesResponseParams{
			InstanceList: []*es.InstanceInfo{
				{
					InstanceId:              ptrStringEs("es-destroy-protection-test"),
					InstanceName:            ptrStringEs("tf-test-instance"),
					Zone:                    ptrStringEs("ap-guangzhou-3"),
					EsVersion:               ptrStringEs("7.10.1"),
					VpcUid:                  ptrStringEs("vpc-test"),
					SubnetUid:               ptrStringEs("subnet-test"),
					ChargeType:              ptrStringEs("POSTPAID_BY_HOUR"),
					Status:                  ptrInt64Es(1),
					LicenseType:             ptrStringEs("platinum"),
					EnableDestroyProtection: ptrStringEs("OPEN"),
				},
			},
		}
		return resp, nil
	})

	meta := newMockMetaEs()
	res := svces.ResourceTencentCloudElasticsearchInstance()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"version":  "7.10.1",
		"vpc_id":   "vpc-test",
		"password": "Test1234",
	})
	d.SetId("es-destroy-protection-test")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "OPEN", d.Get("enable_destroy_protection").(string))
}

// TestEsInstanceDestroyProtection_Read_Nil verifies the read flow does not overwrite
// enable_destroy_protection state when the API returns a nil value.
func TestEsInstanceDestroyProtection_Read_Nil(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	esClient := &es.Client{}
	patches.ApplyMethodReturn(newMockMetaEs().client, "UseEsClient", esClient)

	patches.ApplyMethodFunc(esClient, "DescribeInstances", func(request *es.DescribeInstancesRequest) (*es.DescribeInstancesResponse, error) {
		resp := es.NewDescribeInstancesResponse()
		resp.Response = &es.DescribeInstancesResponseParams{
			InstanceList: []*es.InstanceInfo{
				{
					InstanceId:   ptrStringEs("es-destroy-protection-test"),
					InstanceName: ptrStringEs("tf-test-instance"),
					Zone:         ptrStringEs("ap-guangzhou-3"),
					EsVersion:    ptrStringEs("7.10.1"),
					VpcUid:       ptrStringEs("vpc-test"),
					SubnetUid:    ptrStringEs("subnet-test"),
					ChargeType:   ptrStringEs("POSTPAID_BY_HOUR"),
					Status:       ptrInt64Es(1),
					LicenseType:  ptrStringEs("platinum"),
				},
			},
		}
		return resp, nil
	})

	meta := newMockMetaEs()
	res := svces.ResourceTencentCloudElasticsearchInstance()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"version":                   "7.10.1",
		"vpc_id":                    "vpc-test",
		"password":                  "Test1234",
		"enable_destroy_protection": "OPEN",
	})
	d.SetId("es-destroy-protection-test")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	// state should be preserved (nil-safe read does not overwrite)
	assert.Equal(t, "OPEN", d.Get("enable_destroy_protection").(string))
}

// TestEsInstanceDestroyProtection_Create verifies the create flow invokes UpdateInstance with
// EnableDestroyProtection set to the configured value after instance creation.
func TestEsInstanceDestroyProtection_Create(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	esClient := &es.Client{}
	patches.ApplyMethodReturn(newMockMetaEs().client, "UseEsClient", esClient)

	// mock CreateInstance
	patches.ApplyMethodFunc(esClient, "CreateInstance", func(request *es.CreateInstanceRequest) (*es.CreateInstanceResponse, error) {
		resp := es.NewCreateInstanceResponse()
		resp.Response = &es.CreateInstanceResponseParams{
			InstanceId: ptrStringEs("es-create-destroy-protection"),
		}
		return resp, nil
	})

	// track whether UpdateInstance was called with EnableDestroyProtection = OPEN
	var capturedEnableDestroyProtection *string
	patches.ApplyMethodFunc(esClient, "UpdateInstance", func(request *es.UpdateInstanceRequest) (*es.UpdateInstanceResponse, error) {
		capturedEnableDestroyProtection = request.EnableDestroyProtection
		resp := es.NewUpdateInstanceResponse()
		resp.Response = &es.UpdateInstanceResponseParams{}
		return resp, nil
	})

	// mock DescribeInstances for the post-create status waiting + final read
	callCount := 0
	patches.ApplyMethodFunc(esClient, "DescribeInstances", func(request *es.DescribeInstancesRequest) (*es.DescribeInstancesResponse, error) {
		callCount++
		resp := es.NewDescribeInstancesResponse()
		resp.Response = &es.DescribeInstancesResponseParams{
			InstanceList: []*es.InstanceInfo{
				{
					InstanceId:              ptrStringEs("es-create-destroy-protection"),
					InstanceName:            ptrStringEs("tf-test-instance"),
					Zone:                    ptrStringEs("ap-guangzhou-3"),
					EsVersion:               ptrStringEs("7.10.1"),
					VpcUid:                  ptrStringEs("vpc-test"),
					SubnetUid:               ptrStringEs("subnet-test"),
					ChargeType:              ptrStringEs("POSTPAID_BY_HOUR"),
					Status:                  ptrInt64Es(1),
					LicenseType:             ptrStringEs("platinum"),
					KibanaPublicAccess:      ptrStringEs("OPEN"),
					EnableDestroyProtection: ptrStringEs("OPEN"),
				},
			},
		}
		return resp, nil
	})

	meta := newMockMetaEs()
	res := svces.ResourceTencentCloudElasticsearchInstance()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_name":             "tf-test-instance",
		"availability_zone":         "ap-guangzhou-3",
		"version":                   "7.10.1",
		"vpc_id":                    "vpc-test",
		"subnet_id":                 "subnet-test",
		"password":                  "Test1234",
		"charge_type":               "POSTPAID_BY_HOUR",
		"license_type":              "platinum",
		"enable_destroy_protection": "OPEN",
		"node_info_list": []interface{}{
			map[string]interface{}{
				"node_num":  2,
				"node_type": "ES.S1.MEDIUM8",
				"type":      "hotData",
				"disk_type": "CLOUD_SSD",
				"disk_size": 100,
				"encrypt":   false,
			},
		},
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "es-create-destroy-protection", d.Id())
	// the UpdateInstance call for destroy protection must have carried EnableDestroyProtection = OPEN
	assert.NotNil(t, capturedEnableDestroyProtection)
	assert.Equal(t, "OPEN", *capturedEnableDestroyProtection)
	// final read should have populated state from DescribeInstances
	assert.Equal(t, "OPEN", d.Get("enable_destroy_protection").(string))
}

// TestEsInstanceDestroyProtection_Update verifies the update flow invokes UpdateInstance with
// EnableDestroyProtection when d.HasChange("enable_destroy_protection") is true.
func TestEsInstanceDestroyProtection_Update(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	esClient := &es.Client{}
	patches.ApplyMethodReturn(newMockMetaEs().client, "UseEsClient", esClient)

	var capturedEnableDestroyProtection *string
	patches.ApplyMethodFunc(esClient, "UpdateInstance", func(request *es.UpdateInstanceRequest) (*es.UpdateInstanceResponse, error) {
		if request.EnableDestroyProtection != nil {
			capturedEnableDestroyProtection = request.EnableDestroyProtection
		}
		resp := es.NewUpdateInstanceResponse()
		resp.Response = &es.UpdateInstanceResponseParams{}
		return resp, nil
	})

	// mock DescribeInstances for the upgrade-wait helper + final read
	patches.ApplyMethodFunc(esClient, "DescribeInstances", func(request *es.DescribeInstancesRequest) (*es.DescribeInstancesResponse, error) {
		resp := es.NewDescribeInstancesResponse()
		resp.Response = &es.DescribeInstancesResponseParams{
			InstanceList: []*es.InstanceInfo{
				{
					InstanceId:              ptrStringEs("es-update-destroy-protection"),
					InstanceName:            ptrStringEs("tf-test-instance"),
					Zone:                    ptrStringEs("ap-guangzhou-3"),
					EsVersion:               ptrStringEs("7.10.1"),
					VpcUid:                  ptrStringEs("vpc-test"),
					SubnetUid:               ptrStringEs("subnet-test"),
					ChargeType:              ptrStringEs("POSTPAID_BY_HOUR"),
					Status:                  ptrInt64Es(1),
					LicenseType:             ptrStringEs("platinum"),
					EnableDestroyProtection: ptrStringEs("OPEN"),
				},
			},
		}
		return resp, nil
	})

	meta := newMockMetaEs()
	res := svces.ResourceTencentCloudElasticsearchInstance()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_name":             "tf-test-instance",
		"availability_zone":         "ap-guangzhou-3",
		"version":                   "7.10.1",
		"vpc_id":                    "vpc-test",
		"subnet_id":                 "subnet-test",
		"password":                  "Test1234",
		"charge_type":               "POSTPAID_BY_HOUR",
		"license_type":              "platinum",
		"enable_destroy_protection": "OPEN",
	})
	d.SetId("es-update-destroy-protection")

	// force only enable_destroy_protection to be detected as changed
	patches.ApplyMethodFunc(d, "HasChange", func(key string) bool {
		return key == "enable_destroy_protection"
	})

	err := res.Update(d, meta)
	assert.NoError(t, err)
	assert.NotNil(t, capturedEnableDestroyProtection)
	assert.Equal(t, "OPEN", *capturedEnableDestroyProtection)
	assert.Equal(t, "OPEN", d.Get("enable_destroy_protection").(string))
}

// TestEsInstanceDestroyProtection_Update_NoChange verifies the update flow does NOT call
// UpdateInstance for destroy protection when d.HasChange("enable_destroy_protection") is false.
func TestEsInstanceDestroyProtection_Update_NoChange(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	esClient := &es.Client{}
	patches.ApplyMethodReturn(newMockMetaEs().client, "UseEsClient", esClient)

	updateCalled := false
	patches.ApplyMethodFunc(esClient, "UpdateInstance", func(request *es.UpdateInstanceRequest) (*es.UpdateInstanceResponse, error) {
		if request.EnableDestroyProtection != nil {
			updateCalled = true
		}
		resp := es.NewUpdateInstanceResponse()
		resp.Response = &es.UpdateInstanceResponseParams{}
		return resp, nil
	})

	// mock DescribeInstances for the final read
	patches.ApplyMethodFunc(esClient, "DescribeInstances", func(request *es.DescribeInstancesRequest) (*es.DescribeInstancesResponse, error) {
		resp := es.NewDescribeInstancesResponse()
		resp.Response = &es.DescribeInstancesResponseParams{
			InstanceList: []*es.InstanceInfo{
				{
					InstanceId:              ptrStringEs("es-update-destroy-protection"),
					InstanceName:            ptrStringEs("tf-test-instance"),
					Zone:                    ptrStringEs("ap-guangzhou-3"),
					EsVersion:               ptrStringEs("7.10.1"),
					VpcUid:                  ptrStringEs("vpc-test"),
					SubnetUid:               ptrStringEs("subnet-test"),
					ChargeType:              ptrStringEs("POSTPAID_BY_HOUR"),
					Status:                  ptrInt64Es(1),
					LicenseType:             ptrStringEs("platinum"),
					EnableDestroyProtection: ptrStringEs("CLOSE"),
				},
			},
		}
		return resp, nil
	})

	meta := newMockMetaEs()
	res := svces.ResourceTencentCloudElasticsearchInstance()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_name":             "tf-test-instance",
		"availability_zone":         "ap-guangzhou-3",
		"version":                   "7.10.1",
		"vpc_id":                    "vpc-test",
		"subnet_id":                 "subnet-test",
		"password":                  "Test1234",
		"charge_type":               "POSTPAID_BY_HOUR",
		"license_type":              "platinum",
		"enable_destroy_protection": "CLOSE",
	})
	d.SetId("es-update-destroy-protection")

	// no changes detected
	patches.ApplyMethodFunc(d, "HasChange", func(key string) bool {
		return false
	})

	err := res.Update(d, meta)
	assert.NoError(t, err)
	assert.False(t, updateCalled)
}
