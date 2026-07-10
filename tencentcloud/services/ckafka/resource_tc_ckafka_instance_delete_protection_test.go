package ckafka_test

import (
	"context"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	ckafka "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ckafka/v20190819"
	tag "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tag/v20180813"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	localckafka "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/ckafka"
)

type mockMetaCkafkaInstanceDeleteProtection struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaCkafkaInstanceDeleteProtection) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaCkafkaInstanceDeleteProtection{}

func newMockMetaCkafkaInstanceDeleteProtection() *mockMetaCkafkaInstanceDeleteProtection {
	return &mockMetaCkafkaInstanceDeleteProtection{client: &connectivity.TencentCloudClient{}}
}

func dpPtrString(s string) *string { return &s }
func dpPtrInt64(v int64) *int64    { return &v }
func dpPtrUint64(v uint64) *uint64 { return &v }

// buildCkafkaInstanceDetail builds a minimal InstanceDetail that is safe for the
// Read path (fields dereferenced in resourceTencentCloudCkafkaInstanceRead are
// non-nil).
func buildCkafkaInstanceDetail(instanceId string, status int64) *ckafka.InstanceDetail {
	bandwidth := int64(20)
	return &ckafka.InstanceDetail{
		InstanceId:   dpPtrString(instanceId),
		InstanceName: dpPtrString("tf-example"),
		Vip:          dpPtrString("10.0.0.1"),
		Vport:        dpPtrString("9092"),
		Status:       dpPtrInt64(status),
		Bandwidth:    dpPtrInt64(bandwidth),
		ZoneId:       dpPtrInt64(100007),
		RenewFlag:    dpPtrInt64(0),
		Version:      dpPtrString("2.8.1"),
		DiskSize:     dpPtrInt64(200),
		DiskType:     dpPtrString("CLOUD_BASIC"),
		ExpireTime:   dpPtrInt64(0),
		InstanceType: dpPtrString("profession"),
	}
}

// patchCkafkaReadSide mocks the Read-side APIs (DescribeInstancesDetail for both
// CheckCkafkaInstanceReady and DescribeCkafkaInstanceById, DescribeInstanceAttributes
// and the tag service) so that resourceTencentCloudCkafkaInstanceRead can run
// end-to-end after Create/Update.
func patchCkafkaReadSide(patches *gomonkey.Patches, ckafkaClient *ckafka.Client, instanceId string, deleteProtectionEnable *int64) {
	patches.ApplyMethodFunc(ckafkaClient, "DescribeInstancesDetail", func(request *ckafka.DescribeInstancesDetailRequest) (*ckafka.DescribeInstancesDetailResponse, error) {
		resp := ckafka.NewDescribeInstancesDetailResponse()
		resp.Response = &ckafka.DescribeInstancesDetailResponseParams{
			Result: &ckafka.InstanceDetailResponse{
				TotalCount:   dpPtrInt64(1),
				InstanceList: []*ckafka.InstanceDetail{buildCkafkaInstanceDetail(instanceId, 1)},
			},
			RequestId: dpPtrString("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(ckafkaClient, "DescribeInstanceAttributes", func(request *ckafka.DescribeInstanceAttributesRequest) (*ckafka.DescribeInstanceAttributesResponse, error) {
		resp := ckafka.NewDescribeInstanceAttributesResponse()
		resp.Response = &ckafka.DescribeInstanceAttributesResponseParams{
			Result: &ckafka.InstanceAttributesResponse{
				InstanceId:             dpPtrString(instanceId),
				InstanceName:           dpPtrString("tf-example"),
				MsgRetentionTime:       dpPtrInt64(1300),
				PublicNetwork:          dpPtrInt64(0),
				DeleteProtectionEnable: deleteProtectionEnable,
				RetentionTimeConfig: &ckafka.DynamicRetentionTime{
					Enable: dpPtrInt64(1),
				},
			},
			RequestId: dpPtrString("fake-request-id"),
		}
		return resp, nil
	})

	tagClient := &tag.Client{}
	patches.ApplyMethodReturn(newMockMetaCkafkaInstanceDeleteProtection().client, "UseTagClient", tagClient)
	patches.ApplyMethodFunc(tagClient, "DescribeResourceTagsByResourceIds", func(request *tag.DescribeResourceTagsByResourceIdsRequest) (*tag.DescribeResourceTagsByResourceIdsResponse, error) {
		resp := tag.NewDescribeResourceTagsByResourceIdsResponse()
		resp.Response = &tag.DescribeResourceTagsByResourceIdsResponseParams{
			TotalCount: dpPtrUint64(0),
			Offset:     dpPtrUint64(0),
			Limit:      dpPtrUint64(15),
			Tags:       []*tag.TagResource{},
			RequestId:  dpPtrString("fake-request-id"),
		}
		return resp, nil
	})
}

// go test ./tencentcloud/services/ckafka/ -run "TestCkafkaInstanceDeleteProtection" -v -count=1 -gcflags="all=-l"

// TestCkafkaInstanceDeleteProtection_Create_Enable verifies that creating an
// instance with delete_protection_enable=1 fills DeleteProtectionEnable=1 in the
// ModifyInstanceAttributes request.
func TestCkafkaInstanceDeleteProtection_Create_Enable(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ckafkaClient := &ckafka.Client{}
	patches.ApplyMethodReturn(newMockMetaCkafkaInstanceDeleteProtection().client, "UseCkafkaClient", ckafkaClient)

	instanceId := "ckafka-test-create-enable"
	var capturedDeleteProtectionEnable *int64

	patches.ApplyMethodFunc(ckafkaClient, "CreatePostPaidInstance", func(request *ckafka.CreatePostPaidInstanceRequest) (*ckafka.CreatePostPaidInstanceResponse, error) {
		resp := ckafka.NewCreatePostPaidInstanceResponse()
		resp.Response = &ckafka.CreatePostPaidInstanceResponseParams{
			Result: &ckafka.CreateInstancePostResp{
				ReturnCode: dpPtrString("0"),
				Data: &ckafka.CreateInstancePostData{
					InstanceId: dpPtrString(instanceId),
				},
			},
			RequestId: dpPtrString("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(ckafkaClient, "ModifyInstanceAttributes", func(request *ckafka.ModifyInstanceAttributesRequest) (*ckafka.ModifyInstanceAttributesResponse, error) {
		capturedDeleteProtectionEnable = request.DeleteProtectionEnable
		resp := ckafka.NewModifyInstanceAttributesResponse()
		resp.Response = &ckafka.ModifyInstanceAttributesResponseParams{
			RequestId: dpPtrString("fake-request-id"),
		}
		return resp, nil
	})

	patchCkafkaReadSide(patches, ckafkaClient, instanceId, dpPtrInt64(1))

	meta := newMockMetaCkafkaInstanceDeleteProtection()
	res := localckafka.ResourceTencentCloudCkafkaInstance()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_name":            "tf-example",
		"zone_id":                  100007,
		"charge_type":              "POSTPAID_BY_HOUR",
		"delete_protection_enable": 1,
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, instanceId, d.Id())
	assert.NotNil(t, capturedDeleteProtectionEnable)
	assert.Equal(t, int64(1), *capturedDeleteProtectionEnable)
	assert.Equal(t, 1, d.Get("delete_protection_enable").(int))
}

// TestCkafkaInstanceDeleteProtection_Create_Disable verifies that creating an
// instance with an explicit delete_protection_enable=0 fills
// DeleteProtectionEnable=0 (GetOkExists must not skip the zero value).
func TestCkafkaInstanceDeleteProtection_Create_Disable(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ckafkaClient := &ckafka.Client{}
	patches.ApplyMethodReturn(newMockMetaCkafkaInstanceDeleteProtection().client, "UseCkafkaClient", ckafkaClient)

	instanceId := "ckafka-test-create-disable"
	var capturedDeleteProtectionEnable *int64

	patches.ApplyMethodFunc(ckafkaClient, "CreatePostPaidInstance", func(request *ckafka.CreatePostPaidInstanceRequest) (*ckafka.CreatePostPaidInstanceResponse, error) {
		resp := ckafka.NewCreatePostPaidInstanceResponse()
		resp.Response = &ckafka.CreatePostPaidInstanceResponseParams{
			Result: &ckafka.CreateInstancePostResp{
				ReturnCode: dpPtrString("0"),
				Data: &ckafka.CreateInstancePostData{
					InstanceId: dpPtrString(instanceId),
				},
			},
			RequestId: dpPtrString("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(ckafkaClient, "ModifyInstanceAttributes", func(request *ckafka.ModifyInstanceAttributesRequest) (*ckafka.ModifyInstanceAttributesResponse, error) {
		capturedDeleteProtectionEnable = request.DeleteProtectionEnable
		resp := ckafka.NewModifyInstanceAttributesResponse()
		resp.Response = &ckafka.ModifyInstanceAttributesResponseParams{
			RequestId: dpPtrString("fake-request-id"),
		}
		return resp, nil
	})

	patchCkafkaReadSide(patches, ckafkaClient, instanceId, dpPtrInt64(0))

	meta := newMockMetaCkafkaInstanceDeleteProtection()
	res := localckafka.ResourceTencentCloudCkafkaInstance()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_name":            "tf-example",
		"zone_id":                  100007,
		"charge_type":              "POSTPAID_BY_HOUR",
		"delete_protection_enable": 0,
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, instanceId, d.Id())
	assert.NotNil(t, capturedDeleteProtectionEnable)
	assert.Equal(t, int64(0), *capturedDeleteProtectionEnable)
	assert.Equal(t, 0, d.Get("delete_protection_enable").(int))
}

// TestCkafkaInstanceDeleteProtection_Update_Enable verifies that updating
// delete_protection_enable from 0 to 1 fills DeleteProtectionEnable=1.
func TestCkafkaInstanceDeleteProtection_Update_Enable(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ckafkaClient := &ckafka.Client{}
	patches.ApplyMethodReturn(newMockMetaCkafkaInstanceDeleteProtection().client, "UseCkafkaClient", ckafkaClient)

	instanceId := "ckafka-test-update-enable"
	var capturedDeleteProtectionEnable *int64

	patches.ApplyMethodFunc(ckafkaClient, "ModifyInstanceAttributes", func(request *ckafka.ModifyInstanceAttributesRequest) (*ckafka.ModifyInstanceAttributesResponse, error) {
		capturedDeleteProtectionEnable = request.DeleteProtectionEnable
		resp := ckafka.NewModifyInstanceAttributesResponse()
		resp.Response = &ckafka.ModifyInstanceAttributesResponseParams{
			RequestId: dpPtrString("fake-request-id"),
		}
		return resp, nil
	})

	patchCkafkaReadSide(patches, ckafkaClient, instanceId, dpPtrInt64(1))

	meta := newMockMetaCkafkaInstanceDeleteProtection()
	res := localckafka.ResourceTencentCloudCkafkaInstance()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_name":            "tf-example",
		"zone_id":                  100007,
		"charge_type":              "POSTPAID_BY_HOUR",
		"delete_protection_enable": 1,
	})
	d.SetId(instanceId)

	// force only delete_protection_enable to be detected as changed (so immutableArgs
	// and other update branches are not triggered)
	patches.ApplyMethodFunc(d, "HasChange", func(key string) bool {
		return key == "delete_protection_enable"
	})

	err := res.Update(d, meta)
	assert.NoError(t, err)
	assert.NotNil(t, capturedDeleteProtectionEnable)
	assert.Equal(t, int64(1), *capturedDeleteProtectionEnable)
	assert.Equal(t, 1, d.Get("delete_protection_enable").(int))
}

// TestCkafkaInstanceDeleteProtection_Update_Disable verifies that updating
// delete_protection_enable from 1 to 0 fills DeleteProtectionEnable=0.
func TestCkafkaInstanceDeleteProtection_Update_Disable(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ckafkaClient := &ckafka.Client{}
	patches.ApplyMethodReturn(newMockMetaCkafkaInstanceDeleteProtection().client, "UseCkafkaClient", ckafkaClient)

	instanceId := "ckafka-test-update-disable"
	var capturedDeleteProtectionEnable *int64

	patches.ApplyMethodFunc(ckafkaClient, "ModifyInstanceAttributes", func(request *ckafka.ModifyInstanceAttributesRequest) (*ckafka.ModifyInstanceAttributesResponse, error) {
		capturedDeleteProtectionEnable = request.DeleteProtectionEnable
		resp := ckafka.NewModifyInstanceAttributesResponse()
		resp.Response = &ckafka.ModifyInstanceAttributesResponseParams{
			RequestId: dpPtrString("fake-request-id"),
		}
		return resp, nil
	})

	patchCkafkaReadSide(patches, ckafkaClient, instanceId, dpPtrInt64(0))

	meta := newMockMetaCkafkaInstanceDeleteProtection()
	res := localckafka.ResourceTencentCloudCkafkaInstance()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_name":            "tf-example",
		"zone_id":                  100007,
		"charge_type":              "POSTPAID_BY_HOUR",
		"delete_protection_enable": 0,
	})
	d.SetId(instanceId)

	// force only delete_protection_enable to be detected as changed (so immutableArgs
	// and other update branches are not triggered)
	patches.ApplyMethodFunc(d, "HasChange", func(key string) bool {
		return key == "delete_protection_enable"
	})

	err := res.Update(d, meta)
	assert.NoError(t, err)
	assert.NotNil(t, capturedDeleteProtectionEnable)
	assert.Equal(t, int64(0), *capturedDeleteProtectionEnable)
	assert.Equal(t, 0, d.Get("delete_protection_enable").(int))
}

// TestCkafkaInstanceDeleteProtection_Read_NotNil verifies that Read populates
// delete_protection_enable from the DescribeInstanceAttributes response when the
// field is not nil.
func TestCkafkaInstanceDeleteProtection_Read_NotNil(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ckafkaClient := &ckafka.Client{}
	patches.ApplyMethodReturn(newMockMetaCkafkaInstanceDeleteProtection().client, "UseCkafkaClient", ckafkaClient)

	instanceId := "ckafka-test-read-notnil"
	patchCkafkaReadSide(patches, ckafkaClient, instanceId, dpPtrInt64(1))

	meta := newMockMetaCkafkaInstanceDeleteProtection()
	res := localckafka.ResourceTencentCloudCkafkaInstance()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_name": "tf-example",
		"zone_id":       100007,
		"charge_type":   "POSTPAID_BY_HOUR",
	})
	d.SetId(instanceId)

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, instanceId, d.Id())
	assert.Equal(t, 1, d.Get("delete_protection_enable").(int))
}

// TestCkafkaInstanceDeleteProtection_Read_Nil verifies that Read does not panic
// and skips setting delete_protection_enable when the response field is nil.
func TestCkafkaInstanceDeleteProtection_Read_Nil(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ckafkaClient := &ckafka.Client{}
	patches.ApplyMethodReturn(newMockMetaCkafkaInstanceDeleteProtection().client, "UseCkafkaClient", ckafkaClient)

	instanceId := "ckafka-test-read-nil"
	patchCkafkaReadSide(patches, ckafkaClient, instanceId, nil)

	meta := newMockMetaCkafkaInstanceDeleteProtection()
	res := localckafka.ResourceTencentCloudCkafkaInstance()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_name": "tf-example",
		"zone_id":       100007,
		"charge_type":   "POSTPAID_BY_HOUR",
	})
	d.SetId(instanceId)

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, instanceId, d.Id())
	assert.Equal(t, 0, d.Get("delete_protection_enable").(int))
}

// Ensure the context import is used in case of future helpers referencing it.
var _ = context.TODO
