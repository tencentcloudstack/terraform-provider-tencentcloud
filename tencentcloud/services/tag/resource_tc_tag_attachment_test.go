package tag_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
	tag "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tag/v20180813"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	svctag "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tag"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

// go test -i; go test -test.run TestAccTencentCloudTagAttachmentResource_basic -v
func TestAccTencentCloudTagAttachmentResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckTagAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTagResourceTag,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTagAttachmentExists("tencentcloud_tag_attachment.tag_attachment"),
					resource.TestCheckResourceAttr("tencentcloud_tag_attachment.tag_attachment", "tag_key", "test_terraform_tagAttachment_key"),
					resource.TestCheckResourceAttr("tencentcloud_tag_attachment.tag_attachment", "tag_value", "Terraform_tagAttachment_value"),
					resource.TestCheckResourceAttrSet("tencentcloud_tag_attachment.tag_attachment", "resource")),
			},
			{
				ResourceName:      "tencentcloud_tag_attachment.tag_attachment",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
func testAccCheckTagAttachmentDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_tag_attachment" {
			continue
		}
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service := svctag.NewTagService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())

		tags, err := service.DescribeTagTagAttachmentById(ctx, rs.Primary.Attributes["tag_key"],
			rs.Primary.Attributes["tag_value"], rs.Primary.Attributes["resource"])
		if err != nil {
			return err
		}
		if tags == nil {
			return nil
		}
		return fmt.Errorf("delete tagAttachment key %s fail, still on server", rs.Primary.Attributes["tag_key"])
	}
	return nil
}

func testAccCheckTagAttachmentExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}

		service := svctag.NewTagService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
		res, err := service.DescribeTagTagAttachmentById(ctx, rs.Primary.Attributes["tag_key"],
			rs.Primary.Attributes["tag_value"], rs.Primary.Attributes["resource"])
		if err != nil {
			return err
		}
		if res != nil && res.Resource != nil && res.Tags != nil {
			return nil
		}

		return fmt.Errorf("tagAttachment %s not found on server", rs.Primary.Attributes["tag_key"])
	}
}

const testAccTagResourceTag = tcacctest.DefaultCvmModificationVariable + `
data "tencentcloud_user_info" "info" {}

locals {
  uin = data.tencentcloud_user_info.info.uin
}

resource "tencentcloud_tag_attachment" "tag_attachment" {
  tag_key = "test_terraform_tagAttachment_key"
  tag_value = "Terraform_tagAttachment_value"
  resource = "qcs::cvm:ap-guangzhou:uin/${local.uin}:instance/${var.cvm_id}"
}

`

// --- Unit tests using gomonkey (no TF_ACC needed) ---

type mockMetaTagAttachment struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaTagAttachment) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaTagAttachment{}

func newMockMetaTagAttachment() *mockMetaTagAttachment {
	return &mockMetaTagAttachment{client: &connectivity.TencentCloudClient{}}
}

func ptrStringTagAttachment(s string) *string {
	return &s
}

// buildGetResourcesResponse builds a GetResources response that simulates a
// resource having a tag with the given tagKey and tagValue attached.
func buildGetResourcesResponse(resource, tagKey, tagValue string) *tag.GetResourcesResponse {
	resp := tag.NewGetResourcesResponse()
	resp.Response = &tag.GetResourcesResponseParams{
		ResourceTagMappingList: []*tag.ResourceTagMapping{
			{
				Resource: &resource,
				Tags: []*tag.Tag{
					{
						TagKey:   helper.String(tagKey),
						TagValue: helper.String(tagValue),
					},
				},
			},
		},
		RequestId: ptrStringTagAttachment("fake-request-id"),
	}
	return resp
}

// TestTagAttachmentUpdate_TagValueChanged verifies that when tag_value changes,
// the Update function calls UpdateResourceTagValue with the correct arguments
// and rewrites the composite id with the new tag_value.
func TestTagAttachmentUpdate_TagValueChanged(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	tagClient := &tag.Client{}
	patches.ApplyMethodReturn(newMockMetaTagAttachment().client, "UseTagClient", tagClient)

	var capturedUpdateReq *tag.UpdateResourceTagValueRequest
	patches.ApplyMethodFunc(tagClient, "UpdateResourceTagValue", func(request *tag.UpdateResourceTagValueRequest) (*tag.UpdateResourceTagValueResponse, error) {
		capturedUpdateReq = request
		resp := tag.NewUpdateResourceTagValueResponse()
		resp.Response = &tag.UpdateResourceTagValueResponseParams{
			RequestId: ptrStringTagAttachment("fake-request-id"),
		}
		return resp, nil
	})

	// Mock GetResources called by the read after update, returning the new tag value.
	patches.ApplyMethodFunc(tagClient, "GetResources", func(request *tag.GetResourcesRequest) (*tag.GetResourcesResponse, error) {
		return buildGetResourcesResponse("qcs::cvm:ap-guangzhou:uin/100020512675:instance/ins-test", "env", "prod"), nil
	})

	meta := newMockMetaTagAttachment()
	res := svctag.ResourceTencentCloudTagAttachment()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"tag_key":   "env",
		"tag_value": "prod",
		"resource":  "qcs::cvm:ap-guangzhou:uin/100020512675:instance/ins-test",
	})
	// Set old id containing the old tag_value "dev".
	d.SetId("env" + tccommon.FILED_SP + "dev" + tccommon.FILED_SP + "qcs::cvm:ap-guangzhou:uin/100020512675:instance/ins-test")

	// Patch HasChange to simulate tag_value has changed.
	patches.ApplyMethodFunc(d, "HasChange", func(key string) bool {
		if key == "tag_value" {
			return true
		}
		return false
	})

	// Patch GetChange to return the old and new tag_value.
	patches.ApplyMethodFunc(d, "GetChange", func(key string) (interface{}, interface{}) {
		if key == "tag_value" {
			return "dev", "prod"
		}
		return nil, nil
	})

	err := res.Update(d, meta)
	assert.NoError(t, err)

	// Verify UpdateResourceTagValue was called with correct args.
	assert.NotNil(t, capturedUpdateReq)
	assert.Equal(t, "env", *capturedUpdateReq.TagKey)
	assert.Equal(t, "prod", *capturedUpdateReq.TagValue)
	assert.Equal(t, "qcs::cvm:ap-guangzhou:uin/100020512675:instance/ins-test", *capturedUpdateReq.Resource)

	// Verify the id was rewritten with the new tag_value.
	expectedId := "env" + tccommon.FILED_SP + "prod" + tccommon.FILED_SP + "qcs::cvm:ap-guangzhou:uin/100020512675:instance/ins-test"
	assert.Equal(t, expectedId, d.Id())
	// Verify state reflects new tag_value after read.
	assert.Equal(t, "prod", d.Get("tag_value").(string))
}

// TestTagAttachmentUpdate_TagValueUnchanged verifies that when tag_value does not
// change, the Update function does not call UpdateResourceTagValue.
func TestTagAttachmentUpdate_TagValueUnchanged(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	tagClient := &tag.Client{}
	patches.ApplyMethodReturn(newMockMetaTagAttachment().client, "UseTagClient", tagClient)

	updateCalled := false
	patches.ApplyMethodFunc(tagClient, "UpdateResourceTagValue", func(request *tag.UpdateResourceTagValueRequest) (*tag.UpdateResourceTagValueResponse, error) {
		updateCalled = true
		resp := tag.NewUpdateResourceTagValueResponse()
		resp.Response = &tag.UpdateResourceTagValueResponseParams{
			RequestId: ptrStringTagAttachment("fake-request-id"),
		}
		return resp, nil
	})

	// Mock GetResources called by read after update.
	patches.ApplyMethodFunc(tagClient, "GetResources", func(request *tag.GetResourcesRequest) (*tag.GetResourcesResponse, error) {
		return buildGetResourcesResponse("qcs::cvm:ap-guangzhou:uin/100020512675:instance/ins-test", "env", "dev"), nil
	})

	meta := newMockMetaTagAttachment()
	res := svctag.ResourceTencentCloudTagAttachment()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"tag_key":   "env",
		"tag_value": "dev",
		"resource":  "qcs::cvm:ap-guangzhou:uin/100020512675:instance/ins-test",
	})
	// id matches current tag_value.
	d.SetId("env" + tccommon.FILED_SP + "dev" + tccommon.FILED_SP + "qcs::cvm:ap-guangzhou:uin/100020512675:instance/ins-test")

	// Patch HasChange to simulate tag_value has NOT changed.
	patches.ApplyMethodFunc(d, "HasChange", func(key string) bool {
		return false
	})

	err := res.Update(d, meta)
	assert.NoError(t, err)

	// Verify UpdateResourceTagValue was NOT called.
	assert.False(t, updateCalled, "UpdateResourceTagValue should not be called when tag_value is unchanged")

	// Verify id unchanged.
	assert.Equal(t, "env"+tccommon.FILED_SP+"dev"+tccommon.FILED_SP+"qcs::cvm:ap-guangzhou:uin/100020512675:instance/ins-test", d.Id())
}

// TestTagAttachmentUpdate_APIError verifies that when UpdateResourceTagValue
// returns an error, the Update function returns an error and does not update the id.
func TestTagAttachmentUpdate_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	tagClient := &tag.Client{}
	patches.ApplyMethodReturn(newMockMetaTagAttachment().client, "UseTagClient", tagClient)

	patches.ApplyMethodFunc(tagClient, "UpdateResourceTagValue", func(request *tag.UpdateResourceTagValueRequest) (*tag.UpdateResourceTagValueResponse, error) {
		return nil, fmt.Errorf("internal error")
	})

	meta := newMockMetaTagAttachment()
	res := svctag.ResourceTencentCloudTagAttachment()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"tag_key":   "env",
		"tag_value": "prod",
		"resource":  "qcs::cvm:ap-guangzhou:uin/100020512675:instance/ins-test",
	})
	// old id with old tag_value "dev".
	d.SetId("env" + tccommon.FILED_SP + "dev" + tccommon.FILED_SP + "qcs::cvm:ap-guangzhou:uin/100020512675:instance/ins-test")

	// Patch HasChange to simulate tag_value has changed.
	patches.ApplyMethodFunc(d, "HasChange", func(key string) bool {
		if key == "tag_value" {
			return true
		}
		return false
	})

	// Patch GetChange to return the old and new tag_value.
	patches.ApplyMethodFunc(d, "GetChange", func(key string) (interface{}, interface{}) {
		if key == "tag_value" {
			return "dev", "prod"
		}
		return nil, nil
	})

	err := res.Update(d, meta)
	assert.Error(t, err)

	// Verify id was NOT updated (still the old id with "dev").
	assert.Equal(t, "env"+tccommon.FILED_SP+"dev"+tccommon.FILED_SP+"qcs::cvm:ap-guangzhou:uin/100020512675:instance/ins-test", d.Id())
}

// TestTagAttachmentSchema verifies the schema definition after the change:
// tag_value should not be ForceNew, tag_key and resource should remain ForceNew.
func TestTagAttachmentSchema(t *testing.T) {
	res := svctag.ResourceTencentCloudTagAttachment()

	tagKeyField, ok := res.Schema["tag_key"]
	assert.True(t, ok, "tag_key should exist in schema")
	assert.True(t, tagKeyField.Required, "tag_key should be required")
	assert.True(t, tagKeyField.ForceNew, "tag_key should be ForceNew")

	tagValueField, ok := res.Schema["tag_value"]
	assert.True(t, ok, "tag_value should exist in schema")
	assert.True(t, tagValueField.Required, "tag_value should be required")
	assert.False(t, tagValueField.ForceNew, "tag_value should not be ForceNew")

	resourceField, ok := res.Schema["resource"]
	assert.True(t, ok, "resource should exist in schema")
	assert.True(t, resourceField.Required, "resource should be required")
	assert.True(t, resourceField.ForceNew, "resource should be ForceNew")

	// Verify Update is registered.
	assert.NotNil(t, res.Update, "Update function should be registered")
}
