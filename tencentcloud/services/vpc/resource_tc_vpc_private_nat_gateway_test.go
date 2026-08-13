package vpc_test

import (
	"context"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	vpcsdk "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	svcvpc "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/vpc"
)

func TestAccTencentCloudVpcPrivateNatGatewayResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcPrivateNatGateway,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_vpc_private_nat_gateway.private_nat_gateway", "id"),
					resource.TestCheckResourceAttr("tencentcloud_vpc_private_nat_gateway.private_nat_gateway", "nat_gateway_name", "private-nat-gateway"),
					resource.TestCheckResourceAttrSet("tencentcloud_vpc_private_nat_gateway.private_nat_gateway", "vpc_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_vpc_private_nat_gateway.private_nat_gateway", "cross_domain"),
					resource.TestCheckResourceAttrSet("tencentcloud_vpc_private_nat_gateway.private_nat_gateway", "vpc_type"),
					resource.TestCheckResourceAttr("tencentcloud_vpc_private_nat_gateway.private_nat_gateway", "ccn_id", ""),
				),
			},
			{
				Config: testAccVpcPrivateNatGatewayUpdateName,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_vpc_private_nat_gateway.private_nat_gateway", "id"),
					resource.TestCheckResourceAttr("tencentcloud_vpc_private_nat_gateway.private_nat_gateway", "nat_gateway_name", "private-nat-gateway-update"),
				),
			},
			{
				ResourceName:      "tencentcloud_vpc_private_nat_gateway.private_nat_gateway",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccVpcPrivateNatGateway = `
resource "tencentcloud_vpc" "foo" {
  name       = "private-nat-gateway-vpc"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_vpc_private_nat_gateway" "private_nat_gateway" {
  nat_gateway_name = "private-nat-gateway"
  vpc_id = tencentcloud_vpc.foo.id
}
`

const testAccVpcPrivateNatGatewayUpdateName = `
resource "tencentcloud_vpc" "foo" {
  name       = "private-nat-gateway-vpc"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_vpc_private_nat_gateway" "private_nat_gateway" {
  nat_gateway_name = "private-nat-gateway-update"
  vpc_id = tencentcloud_vpc.foo.id
}
`

// --- gomonkey mock unit tests for tags parameter ---

type mockMetaForVpcPrivateNatGateway struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaForVpcPrivateNatGateway) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaForVpcPrivateNatGateway{}

func newMockMetaForVpcPrivateNatGateway() *mockMetaForVpcPrivateNatGateway {
	return &mockMetaForVpcPrivateNatGateway{client: &connectivity.TencentCloudClient{}}
}

func ptrStrPNG(s string) *string    { return &s }
func ptrBoolPNG(v bool) *bool       { return &v }
func ptrUint64PNG(v uint64) *uint64 { return &v }

// TestVpcPrivateNatGateway_Create_WithTags verifies that when tags is set,
// the CreatePrivateNatGateway request carries the correct Tags value.
func TestVpcPrivateNatGateway_Create_WithTags(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	vpcClient := &vpcsdk.Client{}
	patches.ApplyMethodReturn(newMockMetaForVpcPrivateNatGateway().client, "UseVpcClient", vpcClient)

	var capturedRequest *vpcsdk.CreatePrivateNatGatewayRequest
	patches.ApplyMethodFunc(vpcClient, "CreatePrivateNatGateway", func(request *vpcsdk.CreatePrivateNatGatewayRequest) (*vpcsdk.CreatePrivateNatGatewayResponse, error) {
		capturedRequest = request
		resp := vpcsdk.NewCreatePrivateNatGatewayResponse()
		resp.Response = &vpcsdk.CreatePrivateNatGatewayResponseParams{
			PrivateNatGatewaySet: []*vpcsdk.PrivateNatGateway{
				{
					NatGatewayId:   ptrStrPNG("intranat-abcdefgh"),
					NatGatewayName: ptrStrPNG("tf-test-private-nat-gateway"),
					VpcId:          ptrStrPNG("vpc-xxxxxxxx"),
					Status:         ptrStrPNG("AVAILABLE"),
					TagSet: []*vpcsdk.Tag{
						{Key: ptrStrPNG("key1"), Value: ptrStrPNG("value1")},
						{Key: ptrStrPNG("key2"), Value: ptrStrPNG("value2")},
					},
				},
			},
			TotalCount: ptrUint64PNG(1),
			RequestId:  ptrStrPNG("fake-request-id"),
		}
		return resp, nil
	})

	// Mock DescribeVpcPrivateNatGatewayById for the read-back after create
	patches.ApplyMethodFunc(&svcvpc.VpcService{}, "DescribeVpcPrivateNatGatewayById", func(_ context.Context, _ string) (*vpcsdk.PrivateNatGateway, error) {
		return &vpcsdk.PrivateNatGateway{
			NatGatewayId:   ptrStrPNG("intranat-abcdefgh"),
			NatGatewayName: ptrStrPNG("tf-test-private-nat-gateway"),
			VpcId:          ptrStrPNG("vpc-xxxxxxxx"),
			Status:         ptrStrPNG("AVAILABLE"),
			TagSet: []*vpcsdk.Tag{
				{Key: ptrStrPNG("key1"), Value: ptrStrPNG("value1")},
				{Key: ptrStrPNG("key2"), Value: ptrStrPNG("value2")},
			},
		}, nil
	})

	meta := newMockMetaForVpcPrivateNatGateway()
	res := svcvpc.ResourceTencentCloudVpcPrivateNatGateway()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"nat_gateway_name": "tf-test-private-nat-gateway",
		"vpc_id":           "vpc-xxxxxxxx",
		"tags": map[string]interface{}{
			"key1": "value1",
			"key2": "value2",
		},
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "intranat-abcdefgh", d.Id())

	// Verify Tags were set on the create request
	assert.NotNil(t, capturedRequest.Tags)
	assert.Len(t, capturedRequest.Tags, 2)

	tagMap := make(map[string]string)
	for _, tag := range capturedRequest.Tags {
		if tag.Key != nil {
			tagMap[*tag.Key] = ""
			if tag.Value != nil {
				tagMap[*tag.Key] = *tag.Value
			}
		}
	}
	assert.Equal(t, "value1", tagMap["key1"])
	assert.Equal(t, "value2", tagMap["key2"])

	// Verify tags are read back into state
	tags := d.Get("tags").(map[string]interface{})
	assert.Equal(t, "value1", tags["key1"])
	assert.Equal(t, "value2", tags["key2"])
}

// TestVpcPrivateNatGateway_Create_WithoutTags verifies that when tags is not set,
// the CreatePrivateNatGateway request does not carry Tags.
func TestVpcPrivateNatGateway_Create_WithoutTags(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	vpcClient := &vpcsdk.Client{}
	patches.ApplyMethodReturn(newMockMetaForVpcPrivateNatGateway().client, "UseVpcClient", vpcClient)

	var capturedRequest *vpcsdk.CreatePrivateNatGatewayRequest
	patches.ApplyMethodFunc(vpcClient, "CreatePrivateNatGateway", func(request *vpcsdk.CreatePrivateNatGatewayRequest) (*vpcsdk.CreatePrivateNatGatewayResponse, error) {
		capturedRequest = request
		resp := vpcsdk.NewCreatePrivateNatGatewayResponse()
		resp.Response = &vpcsdk.CreatePrivateNatGatewayResponseParams{
			PrivateNatGatewaySet: []*vpcsdk.PrivateNatGateway{
				{
					NatGatewayId:   ptrStrPNG("intranat-no-tags"),
					NatGatewayName: ptrStrPNG("tf-test-private-nat-gateway"),
					VpcId:          ptrStrPNG("vpc-xxxxxxxx"),
					Status:         ptrStrPNG("AVAILABLE"),
				},
			},
			TotalCount: ptrUint64PNG(1),
			RequestId:  ptrStrPNG("fake-request-id"),
		}
		return resp, nil
	})

	// Mock DescribeVpcPrivateNatGatewayById for the read-back after create
	patches.ApplyMethodFunc(&svcvpc.VpcService{}, "DescribeVpcPrivateNatGatewayById", func(_ context.Context, _ string) (*vpcsdk.PrivateNatGateway, error) {
		return &vpcsdk.PrivateNatGateway{
			NatGatewayId:   ptrStrPNG("intranat-no-tags"),
			NatGatewayName: ptrStrPNG("tf-test-private-nat-gateway"),
			VpcId:          ptrStrPNG("vpc-xxxxxxxx"),
			Status:         ptrStrPNG("AVAILABLE"),
		}, nil
	})

	meta := newMockMetaForVpcPrivateNatGateway()
	res := svcvpc.ResourceTencentCloudVpcPrivateNatGateway()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"nat_gateway_name": "tf-test-private-nat-gateway",
		"vpc_id":           "vpc-xxxxxxxx",
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "intranat-no-tags", d.Id())

	// Verify Tags were not set on the create request
	assert.Empty(t, capturedRequest.Tags)
}

// TestVpcPrivateNatGateway_Read_WithTags verifies that Read converts TagSet to map.
func TestVpcPrivateNatGateway_Read_WithTags(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	patches.ApplyMethodFunc(&svcvpc.VpcService{}, "DescribeVpcPrivateNatGatewayById", func(_ context.Context, _ string) (*vpcsdk.PrivateNatGateway, error) {
		return &vpcsdk.PrivateNatGateway{
			NatGatewayId:   ptrStrPNG("intranat-abcdefgh"),
			NatGatewayName: ptrStrPNG("tf-test-private-nat-gateway"),
			VpcId:          ptrStrPNG("vpc-xxxxxxxx"),
			Status:         ptrStrPNG("AVAILABLE"),
			TagSet: []*vpcsdk.Tag{
				{Key: ptrStrPNG("env"), Value: ptrStrPNG("prod")},
				{Key: ptrStrPNG("owner"), Value: ptrStrPNG("devops")},
			},
		}, nil
	})

	meta := newMockMetaForVpcPrivateNatGateway()
	res := svcvpc.ResourceTencentCloudVpcPrivateNatGateway()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
	d.SetId("intranat-abcdefgh")

	err := res.Read(d, meta)
	assert.NoError(t, err)

	tags := d.Get("tags").(map[string]interface{})
	assert.Equal(t, "prod", tags["env"])
	assert.Equal(t, "devops", tags["owner"])
}

// TestVpcPrivateNatGateway_Read_WithoutTags verifies that Read does not set tags
// when TagSet is nil or empty.
func TestVpcPrivateNatGateway_Read_WithoutTags(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	patches.ApplyMethodFunc(&svcvpc.VpcService{}, "DescribeVpcPrivateNatGatewayById", func(_ context.Context, _ string) (*vpcsdk.PrivateNatGateway, error) {
		return &vpcsdk.PrivateNatGateway{
			NatGatewayId:   ptrStrPNG("intranat-abcdefgh"),
			NatGatewayName: ptrStrPNG("tf-test-private-nat-gateway"),
			VpcId:          ptrStrPNG("vpc-xxxxxxxx"),
			Status:         ptrStrPNG("AVAILABLE"),
		}, nil
	})

	meta := newMockMetaForVpcPrivateNatGateway()
	res := svcvpc.ResourceTencentCloudVpcPrivateNatGateway()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
	d.SetId("intranat-abcdefgh")

	err := res.Read(d, meta)
	assert.NoError(t, err)

	tags := d.Get("tags")
	tagsMap, ok := tags.(map[string]interface{})
	assert.True(t, ok)
	assert.Empty(t, tagsMap)
}

// TestVpcPrivateNatGateway_Update_TagsImmutable verifies that changing tags
// in Update returns an error because tags is immutable.
func TestVpcPrivateNatGateway_Update_TagsImmutable(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	meta := newMockMetaForVpcPrivateNatGateway()
	res := svcvpc.ResourceTencentCloudVpcPrivateNatGateway()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"nat_gateway_name": "tf-test-private-nat-gateway",
		"tags": map[string]interface{}{
			"key1": "value1",
		},
	})
	d.SetId("intranat-abcdefgh")

	err := res.Update(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "tags")
	assert.Contains(t, err.Error(), "cannot be changed")
}

// TestVpcPrivateNatGateway_Schema_Tags validates the tags schema definition.
func TestVpcPrivateNatGateway_Schema_Tags(t *testing.T) {
	res := svcvpc.ResourceTencentCloudVpcPrivateNatGateway()
	assert.Contains(t, res.Schema, "tags")

	tagsSchema := res.Schema["tags"]
	assert.Equal(t, schema.TypeMap, tagsSchema.Type)
	assert.True(t, tagsSchema.Optional)
	assert.False(t, tagsSchema.Required)
}
