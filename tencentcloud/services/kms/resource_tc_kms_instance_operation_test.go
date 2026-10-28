package kms_test

import (
	"context"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	kms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/kms/v20190118"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/kms"
)

type mockMetaKmsInstanceOperation struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaKmsInstanceOperation) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaKmsInstanceOperation{}

func newMockMetaKmsInstanceOperation() *mockMetaKmsInstanceOperation {
	return &mockMetaKmsInstanceOperation{client: &connectivity.TencentCloudClient{}}
}

func ptrStringKmsInstanceOp(s string) *string { return &s }

// TestKmsInstanceOperation_Create_Success verifies the Create flow sets the id
// from the returned InstanceId and populates computed fields.
func TestKmsInstanceOperation_Create_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	kmsClient := &kms.Client{}
	patches.ApplyMethodReturn(newMockMetaKmsInstanceOperation().client, "UseKmsClient", kmsClient)

	var capturedRequest *kms.CreateKmsInstanceRequest
	patches.ApplyMethodFunc(kmsClient, "CreateKmsInstanceWithContext", func(_ context.Context, request *kms.CreateKmsInstanceRequest) (*kms.CreateKmsInstanceResponse, error) {
		capturedRequest = request
		assert.Equal(t, "PREPAID", *request.InstanceChargeType)
		assert.Equal(t, "professional", *request.ProductVersion)
		assert.NotNil(t, request.InstanceChargePrepaid)
		assert.Equal(t, uint64(1), *request.InstanceChargePrepaid.Period)
		assert.Equal(t, uint64(1), *request.InstanceChargePrepaid.RenewFlag)
		assert.Len(t, request.Tags, 1)
		assert.Equal(t, "env", *request.Tags[0].TagKey)
		assert.Equal(t, "prod", *request.Tags[0].TagValue)

		resp := kms.NewCreateKmsInstanceResponse()
		resp.Response = &kms.CreateKmsInstanceResponseParams{
			InstanceId: ptrStringKmsInstanceOp("kms-abc123"),
			DealName:   ptrStringKmsInstanceOp("20261028110001"),
			BillId:     ptrStringKmsInstanceOp("bill-001"),
			RequestId:  ptrStringKmsInstanceOp("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaKmsInstanceOperation()
	res := kms.ResourceTencentCloudKmsInstanceOperation()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_charge_type": "PREPAID",
		"product_version":      "professional",
		"client_token":         "tf-token-xxx",
		"instance_charge_prepaid": []interface{}{
			map[string]interface{}{
				"period":     1,
				"renew_flag": 1,
			},
		},
		"tags": []interface{}{
			map[string]interface{}{
				"tag_key":   "env",
				"tag_value": "prod",
			},
		},
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "kms-abc123", d.Id())
	assert.Equal(t, "kms-abc123", d.Get("instance_id"))
	assert.Equal(t, "20261028110001", d.Get("deal_name"))
	assert.Equal(t, "bill-001", d.Get("bill_id"))
	assert.NotNil(t, capturedRequest)
}

// TestKmsInstanceOperation_Create_EmptyResponse verifies that an empty response
// does not set an id and returns an error.
func TestKmsInstanceOperation_Create_EmptyResponse(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	kmsClient := &kms.Client{}
	patches.ApplyMethodReturn(newMockMetaKmsInstanceOperation().client, "UseKmsClient", kmsClient)

	patches.ApplyMethodFunc(kmsClient, "CreateKmsInstanceWithContext", func(_ context.Context, request *kms.CreateKmsInstanceRequest) (*kms.CreateKmsInstanceResponse, error) {
		resp := kms.NewCreateKmsInstanceResponse()
		resp.Response = &kms.CreateKmsInstanceResponseParams{
			InstanceId: ptrStringKmsInstanceOp(""),
			RequestId:  ptrStringKmsInstanceOp("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaKmsInstanceOperation()
	res := kms.ResourceTencentCloudKmsInstanceOperation()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_charge_type": "PREPAID",
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
	assert.Equal(t, "", d.Id())
}

// TestKmsInstanceOperation_Read_Delete verifies Read and Delete are no-ops.
func TestKmsInstanceOperation_Read_Delete(t *testing.T) {
	meta := newMockMetaKmsInstanceOperation()
	res := kms.ResourceTencentCloudKmsInstanceOperation()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_charge_type": "PREPAID",
	})
	d.SetId("kms-abc123")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "kms-abc123", d.Id())

	err = res.Delete(d, meta)
	assert.NoError(t, err)
}

// TestKmsInstanceOperation_Schema validates the schema definition.
func TestKmsInstanceOperation_Schema(t *testing.T) {
	res := kms.ResourceTencentCloudKmsInstanceOperation()

	assert.NotNil(t, res)
	assert.NotNil(t, res.Create)
	assert.NotNil(t, res.Read)
	assert.Nil(t, res.Update)
	assert.NotNil(t, res.Delete)
	assert.Nil(t, res.Importer)

	assert.Contains(t, res.Schema, "instance_charge_type")
	assert.Contains(t, res.Schema, "product_version")
	assert.Contains(t, res.Schema, "instance_charge_prepaid")
	assert.Contains(t, res.Schema, "client_token")
	assert.Contains(t, res.Schema, "tags")
	assert.Contains(t, res.Schema, "instance_id")
	assert.Contains(t, res.Schema, "deal_name")
	assert.Contains(t, res.Schema, "bill_id")

	instanceChargeType := res.Schema["instance_charge_type"]
	assert.Equal(t, schema.TypeString, instanceChargeType.Type)
	assert.True(t, instanceChargeType.Required)
	assert.True(t, instanceChargeType.ForceNew)

	instanceChargePrepaid := res.Schema["instance_charge_prepaid"]
	assert.Equal(t, schema.TypeList, instanceChargePrepaid.Type)
	assert.True(t, instanceChargePrepaid.Optional)
	assert.True(t, instanceChargePrepaid.ForceNew)

	tags := res.Schema["tags"]
	assert.Equal(t, schema.TypeList, tags.Type)

	instanceId := res.Schema["instance_id"]
	assert.Equal(t, schema.TypeString, instanceId.Type)
	assert.True(t, instanceId.Computed)

	dealName := res.Schema["deal_name"]
	assert.Equal(t, schema.TypeString, dealName.Type)
	assert.True(t, dealName.Computed)

	billId := res.Schema["bill_id"]
	assert.Equal(t, schema.TypeString, billId.Type)
	assert.True(t, billId.Computed)
}
