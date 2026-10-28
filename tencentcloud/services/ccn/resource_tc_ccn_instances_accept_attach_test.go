package ccn_test

import (
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	vpcsdk "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	svcccn "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/ccn"
)

type mockMetaForCcnInstancesAcceptAttach struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaForCcnInstancesAcceptAttach) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaForCcnInstancesAcceptAttach{}

func newMockMetaForCcnInstancesAcceptAttach() *mockMetaForCcnInstancesAcceptAttach {
	return &mockMetaForCcnInstancesAcceptAttach{client: &connectivity.TencentCloudClient{}}
}

func ptrStrCIAA(s string) *string { return &s }

// TestCcnInstancesAcceptAttach_Create_WithOrderType verifies that when order_type is set,
// the AcceptAttachCcnInstances request carries the correct OrderType value on the CcnInstance.
func TestCcnInstancesAcceptAttach_Create_WithOrderType(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	vpcClient := &vpcsdk.Client{}
	patches.ApplyMethodReturn(newMockMetaForCcnInstancesAcceptAttach().client, "UseVpcClient", vpcClient)

	var capturedRequest *vpcsdk.AcceptAttachCcnInstancesRequest
	patches.ApplyMethodFunc(vpcClient, "AcceptAttachCcnInstances", func(request *vpcsdk.AcceptAttachCcnInstancesRequest) (*vpcsdk.AcceptAttachCcnInstancesResponse, error) {
		capturedRequest = request
		resp := vpcsdk.NewAcceptAttachCcnInstancesResponse()
		resp.Response = &vpcsdk.AcceptAttachCcnInstancesResponseParams{
			RequestId: ptrStrCIAA("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForCcnInstancesAcceptAttach()
	res := svcccn.ResourceTencentCloudCcnInstancesAcceptAttach()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"ccn_id": "ccn-xxxxxxxx",
		"instances": []interface{}{
			map[string]interface{}{
				"instance_id":     "vpc-yyyyyyyy",
				"instance_region": "ap-guangzhou",
				"instance_type":   "VPC",
				"order_type":      "PayByCcnOwner",
				"route_table_id":  "rtb-zzzzzzzz",
				"description":     "test",
			},
		},
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "ccn-xxxxxxxx", d.Id())

	// Verify the request captured the OrderType
	assert.NotNil(t, capturedRequest)
	assert.Len(t, capturedRequest.Instances, 1)
	assert.NotNil(t, capturedRequest.Instances[0].OrderType)
	assert.Equal(t, "PayByCcnOwner", *capturedRequest.Instances[0].OrderType)
}

// TestCcnInstancesAcceptAttach_Create_WithoutOrderType verifies that when order_type is not set,
// the AcceptAttachCcnInstances request does not carry OrderType on the CcnInstance.
func TestCcnInstancesAcceptAttach_Create_WithoutOrderType(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	vpcClient := &vpcsdk.Client{}
	patches.ApplyMethodReturn(newMockMetaForCcnInstancesAcceptAttach().client, "UseVpcClient", vpcClient)

	var capturedRequest *vpcsdk.AcceptAttachCcnInstancesRequest
	patches.ApplyMethodFunc(vpcClient, "AcceptAttachCcnInstances", func(request *vpcsdk.AcceptAttachCcnInstancesRequest) (*vpcsdk.AcceptAttachCcnInstancesResponse, error) {
		capturedRequest = request
		resp := vpcsdk.NewAcceptAttachCcnInstancesResponse()
		resp.Response = &vpcsdk.AcceptAttachCcnInstancesResponseParams{
			RequestId: ptrStrCIAA("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForCcnInstancesAcceptAttach()
	res := svcccn.ResourceTencentCloudCcnInstancesAcceptAttach()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"ccn_id": "ccn-xxxxxxxx",
		"instances": []interface{}{
			map[string]interface{}{
				"instance_id":     "vpc-yyyyyyyy",
				"instance_region": "ap-guangzhou",
				"instance_type":   "VPC",
			},
		},
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "ccn-xxxxxxxx", d.Id())

	// Verify OrderType was not set on the request
	assert.NotNil(t, capturedRequest)
	assert.Len(t, capturedRequest.Instances, 1)
	assert.Nil(t, capturedRequest.Instances[0].OrderType)
}

// TestCcnInstancesAcceptAttach_Create_WithEmptyOrderType verifies that when order_type is empty string,
// the AcceptAttachCcnInstances request does not carry OrderType on the CcnInstance.
func TestCcnInstancesAcceptAttach_Create_WithEmptyOrderType(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	vpcClient := &vpcsdk.Client{}
	patches.ApplyMethodReturn(newMockMetaForCcnInstancesAcceptAttach().client, "UseVpcClient", vpcClient)

	var capturedRequest *vpcsdk.AcceptAttachCcnInstancesRequest
	patches.ApplyMethodFunc(vpcClient, "AcceptAttachCcnInstances", func(request *vpcsdk.AcceptAttachCcnInstancesRequest) (*vpcsdk.AcceptAttachCcnInstancesResponse, error) {
		capturedRequest = request
		resp := vpcsdk.NewAcceptAttachCcnInstancesResponse()
		resp.Response = &vpcsdk.AcceptAttachCcnInstancesResponseParams{
			RequestId: ptrStrCIAA("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForCcnInstancesAcceptAttach()
	res := svcccn.ResourceTencentCloudCcnInstancesAcceptAttach()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"ccn_id": "ccn-xxxxxxxx",
		"instances": []interface{}{
			map[string]interface{}{
				"instance_id":     "vpc-yyyyyyyy",
				"instance_region": "ap-guangzhou",
				"instance_type":   "VPC",
				"order_type":      "",
			},
		},
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "ccn-xxxxxxxx", d.Id())

	// Verify OrderType was not set on the request when empty
	assert.NotNil(t, capturedRequest)
	assert.Len(t, capturedRequest.Instances, 1)
	assert.Nil(t, capturedRequest.Instances[0].OrderType)
}

// TestCcnInstancesAcceptAttach_Schema_OrderType validates the order_type schema definition.
func TestCcnInstancesAcceptAttach_Schema_OrderType(t *testing.T) {
	res := svcccn.ResourceTencentCloudCcnInstancesAcceptAttach()
	assert.Contains(t, res.Schema, "instances")

	instancesSchema := res.Schema["instances"]
	assert.Equal(t, schema.TypeList, instancesSchema.Type)

	elemRes, ok := instancesSchema.Elem.(*schema.Resource)
	assert.True(t, ok)
	assert.Contains(t, elemRes.Schema, "order_type")

	orderTypeSchema := elemRes.Schema["order_type"]
	assert.Equal(t, schema.TypeString, orderTypeSchema.Type)
	assert.True(t, orderTypeSchema.Optional)
	assert.False(t, orderTypeSchema.Required)
}
