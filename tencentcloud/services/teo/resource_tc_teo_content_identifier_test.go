package teo_test

import (
	"fmt"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/teo"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

// mockMetaForContentIdentifier implements tccommon.ProviderMeta
type mockMetaForContentIdentifier struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaForContentIdentifier) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaForContentIdentifier{}

func newMockMetaForContentIdentifier() *mockMetaForContentIdentifier {
	return &mockMetaForContentIdentifier{client: &connectivity.TencentCloudClient{}}
}

func ptrStrContentIdentifier(s string) *string {
	return &s
}

func TestAccTencentCloudTeoContentIdentifierResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTeoContentIdentifier,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_teo_content_identifier.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_content_identifier.example", "plan_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_content_identifier.example", "description"),
				),
			},
			{
				Config: testAccTeoContentIdentifierUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_teo_content_identifier.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_content_identifier.example", "plan_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_content_identifier.example", "description"),
				),
			},
			{
				ResourceName:      "tencentcloud_teo_content_identifier.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTeoContentIdentifier = `
resource "tencentcloud_teo_content_identifier" "example" {
  plan_id     = "edgeone-3bzvsgjkfw6g"
  description = "example"
  tags {
    tag_key   = "tagKey"
    tag_value = "tagValue"
  }
}
`

const testAccTeoContentIdentifierUpdate = `
resource "tencentcloud_teo_content_identifier" "example" {
  plan_id     = "edgeone-3bzvsgjkfw6g"
  description = "example update"
  tags {
    tag_key   = "tagKey"
    tag_value = "tagValue"
  }
}
`

// go test ./tencentcloud/services/teo/ -run "TestContentIdentifierStatus" -v -count=1 -gcflags="all=-l"

// TestContentIdentifierStatus_Read_Success tests Read populates status attribute
func TestContentIdentifierStatus_Read_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForContentIdentifier().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeContentIdentifiers", func(request *teov20220901.DescribeContentIdentifiersRequest) (*teov20220901.DescribeContentIdentifiersResponse, error) {
		resp := teov20220901.NewDescribeContentIdentifiersResponse()
		resp.Response = &teov20220901.DescribeContentIdentifiersResponseParams{
			ContentIdentifiers: []*teov20220901.ContentIdentifier{
				{
					ContentId: ptrStrContentIdentifier("eocontent-3dy8iyfq8dba"),
					PlanId:    ptrStrContentIdentifier("edgeone-3bzvsgjkfw6g"),
					Status:    ptrStrContentIdentifier("active"),
				},
			},
			RequestId: ptrStrContentIdentifier("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForContentIdentifier()
	res := teo.ResourceTencentCloudTeoContentIdentifier()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"description": "test",
		"plan_id":     "edgeone-3bzvsgjkfw6g",
	})
	d.SetId("eocontent-3dy8iyfq8dba")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "eocontent-3dy8iyfq8dba", d.Get("content_id"))
	assert.Equal(t, "active", d.Get("status"))
}

// TestContentIdentifierStatus_Read_NilStatus tests Read handles nil status
func TestContentIdentifierStatus_Read_NilStatus(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForContentIdentifier().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeContentIdentifiers", func(request *teov20220901.DescribeContentIdentifiersRequest) (*teov20220901.DescribeContentIdentifiersResponse, error) {
		resp := teov20220901.NewDescribeContentIdentifiersResponse()
		resp.Response = &teov20220901.DescribeContentIdentifiersResponseParams{
			ContentIdentifiers: []*teov20220901.ContentIdentifier{
				{
					ContentId: ptrStrContentIdentifier("eocontent-3dy8iyfq8dba"),
					PlanId:    ptrStrContentIdentifier("edgeone-3bzvsgjkfw6g"),
					Status:    nil,
				},
			},
			RequestId: ptrStrContentIdentifier("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForContentIdentifier()
	res := teo.ResourceTencentCloudTeoContentIdentifier()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"description": "test",
		"plan_id":     "edgeone-3bzvsgjkfw6g",
	})
	d.SetId("eocontent-3dy8iyfq8dba")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Get("status"))
}

// TestContentIdentifierStatus_Read_DeletedStatus tests Read handles deleted status
func TestContentIdentifierStatus_Read_DeletedStatus(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForContentIdentifier().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeContentIdentifiers", func(request *teov20220901.DescribeContentIdentifiersRequest) (*teov20220901.DescribeContentIdentifiersResponse, error) {
		resp := teov20220901.NewDescribeContentIdentifiersResponse()
		resp.Response = &teov20220901.DescribeContentIdentifiersResponseParams{
			ContentIdentifiers: []*teov20220901.ContentIdentifier{
				{
					ContentId: ptrStrContentIdentifier("eocontent-3dy8iyfq8dba"),
					PlanId:    ptrStrContentIdentifier("edgeone-3bzvsgjkfw6g"),
					Status:    ptrStrContentIdentifier("deleted"),
				},
			},
			RequestId: ptrStrContentIdentifier("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForContentIdentifier()
	res := teo.ResourceTencentCloudTeoContentIdentifier()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"description": "test",
		"plan_id":     "edgeone-3bzvsgjkfw6g",
	})
	d.SetId("eocontent-3dy8iyfq8dba")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "deleted", d.Get("status"))
}

// TestContentIdentifierStatus_Read_NotFound tests Read handles resource not found
func TestContentIdentifierStatus_Read_NotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForContentIdentifier().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeContentIdentifiers", func(request *teov20220901.DescribeContentIdentifiersRequest) (*teov20220901.DescribeContentIdentifiersResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=ResourceNotFound, Message=Content identifier not found")
	})

	meta := newMockMetaForContentIdentifier()
	res := teo.ResourceTencentCloudTeoContentIdentifier()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"description": "test",
		"plan_id":     "edgeone-3bzvsgjkfw6g",
	})
	d.SetId("eocontent-3dy8iyfq8dba")

	err := res.Read(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ResourceNotFound")
}

// TestContentIdentifierStatus_Schema tests the status schema definition
func TestContentIdentifierStatus_Schema(t *testing.T) {
	res := teo.ResourceTencentCloudTeoContentIdentifier()

	assert.NotNil(t, res)

	assert.Contains(t, res.Schema, "status")
	statusSchema := res.Schema["status"]
	assert.Equal(t, schema.TypeString, statusSchema.Type)
	assert.True(t, statusSchema.Computed)
	assert.False(t, statusSchema.Optional)
	assert.False(t, statusSchema.Required)
}
