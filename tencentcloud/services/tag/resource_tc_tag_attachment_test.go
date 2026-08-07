package tag_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	svctag "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
	tag "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tag/v20180813"
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
				Config: testAccTagResourceTagUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTagAttachmentExists("tencentcloud_tag_attachment.tag_attachment"),
					resource.TestCheckResourceAttr("tencentcloud_tag_attachment.tag_attachment", "tag_key", "test_terraform_tagAttachment_key"),
					resource.TestCheckResourceAttr("tencentcloud_tag_attachment.tag_attachment", "tag_value", "Terraform_tagAttachment_value_updated"),
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

const testAccTagResourceTagUpdate = tcacctest.DefaultCvmModificationVariable + `
data "tencentcloud_user_info" "info" {}

locals {
  uin = data.tencentcloud_user_info.info.uin
}

resource "tencentcloud_tag_attachment" "tag_attachment" {
  tag_key = "test_terraform_tagAttachment_key"
  tag_value = "Terraform_tagAttachment_value_updated"
  resource = "qcs::cvm:ap-guangzhou:uin/${local.uin}:instance/${var.cvm_id}"
}

`

// mockMetaTagAttachment implements tccommon.ProviderMeta
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

// go test ./tencentcloud/services/tag/ -run "TestTagAttachmentUpdate" -v -count=1 -gcflags="all=-l"

// TestTagAttachmentUpdate_Success tests successful tag_value update via UpdateResourceTagValue
func TestTagAttachmentUpdate_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	tagClient := &tag.Client{}
	patches.ApplyMethodReturn(newMockMetaTagAttachment().client, "UseTagClient", tagClient)

	patches.ApplyMethodFunc(tagClient, "UpdateResourceTagValue", func(request *tag.UpdateResourceTagValueRequest) (*tag.UpdateResourceTagValueResponse, error) {
		assert.NotNil(t, request.TagKey)
		assert.Equal(t, "test_key", *request.TagKey)
		assert.NotNil(t, request.TagValue)
		assert.Equal(t, "new_value", *request.TagValue)
		assert.NotNil(t, request.Resource)
		assert.Equal(t, "qcs::cvm:ap-guangzhou:uin/100020512675:instance/ins-test", *request.Resource)

		resp := tag.NewUpdateResourceTagValueResponse()
		resp.Response = &tag.UpdateResourceTagValueResponseParams{
			RequestId: ptrStringTagAttachment("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(tagClient, "GetResources", func(request *tag.GetResourcesRequest) (*tag.GetResourcesResponse, error) {
		resp := tag.NewGetResourcesResponse()
		resp.Response = &tag.GetResourcesResponseParams{
			ResourceTagMappingList: []*tag.ResourceTagMapping{
				{
					Resource: ptrStringTagAttachment("qcs::cvm:ap-guangzhou:uin/100020512675:instance/ins-test"),
					Tags: []*tag.Tag{
						{
							TagKey:   ptrStringTagAttachment("test_key"),
							TagValue: ptrStringTagAttachment("new_value"),
						},
					},
				},
			},
			RequestId: ptrStringTagAttachment("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaTagAttachment()
	res := svctag.ResourceTencentCloudTagAttachment()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"tag_key":   "test_key",
		"tag_value": "new_value",
		"resource":  "qcs::cvm:ap-guangzhou:uin/100020512675:instance/ins-test",
	})
	d.SetId("test_key" + tccommon.FILED_SP + "old_value" + tccommon.FILED_SP + "qcs::cvm:ap-guangzhou:uin/100020512675:instance/ins-test")

	patches.ApplyMethodFunc(d, "HasChange", func(key string) bool {
		return key == "tag_value"
	})

	err := res.Update(d, meta)
	assert.NoError(t, err)
	expectedId := "test_key" + tccommon.FILED_SP + "new_value" + tccommon.FILED_SP + "qcs::cvm:ap-guangzhou:uin/100020512675:instance/ins-test"
	assert.Equal(t, expectedId, d.Id())
	assert.Equal(t, "new_value", d.Get("tag_value").(string))
}

// TestTagAttachmentUpdate_APIError tests error propagation when UpdateResourceTagValue fails
func TestTagAttachmentUpdate_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	tagClient := &tag.Client{}
	patches.ApplyMethodReturn(newMockMetaTagAttachment().client, "UseTagClient", tagClient)

	patches.ApplyMethodFunc(tagClient, "UpdateResourceTagValue", func(request *tag.UpdateResourceTagValueRequest) (*tag.UpdateResourceTagValueResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=ResourceNotFound.AttachedTagKeyNotFound, Message=tag key not attached")
	})

	meta := newMockMetaTagAttachment()
	res := svctag.ResourceTencentCloudTagAttachment()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"tag_key":   "test_key",
		"tag_value": "new_value",
		"resource":  "qcs::cvm:ap-guangzhou:uin/100020512675:instance/ins-test",
	})
	d.SetId("test_key" + tccommon.FILED_SP + "old_value" + tccommon.FILED_SP + "qcs::cvm:ap-guangzhou:uin/100020512675:instance/ins-test")

	patches.ApplyMethodFunc(d, "HasChange", func(key string) bool {
		return key == "tag_value"
	})

	err := res.Update(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ResourceNotFound")
	originalId := "test_key" + tccommon.FILED_SP + "old_value" + tccommon.FILED_SP + "qcs::cvm:ap-guangzhou:uin/100020512675:instance/ins-test"
	assert.Equal(t, originalId, d.Id(), "ID should not be reset on update failure")
}

// TestTagAttachmentUpdate_Schema verifies that tag_value is no longer ForceNew
func TestTagAttachmentUpdate_Schema(t *testing.T) {
	res := svctag.ResourceTencentCloudTagAttachment()

	assert.NotNil(t, res)
	assert.NotNil(t, res.Update)
	assert.NotNil(t, res.Create)
	assert.NotNil(t, res.Read)
	assert.NotNil(t, res.Delete)

	assert.Contains(t, res.Schema, "tag_key")
	assert.Contains(t, res.Schema, "tag_value")
	assert.Contains(t, res.Schema, "resource")

	tagKey := res.Schema["tag_key"]
	assert.True(t, tagKey.Required)
	assert.True(t, tagKey.ForceNew)

	tagValue := res.Schema["tag_value"]
	assert.True(t, tagValue.Required)
	assert.False(t, tagValue.ForceNew, "tag_value should not be ForceNew")

	resourceField := res.Schema["resource"]
	assert.True(t, resourceField.Required)
	assert.True(t, resourceField.ForceNew)
}

// TestTagAttachmentUpdate_NoChange tests that UpdateResourceTagValue is not called when tag_value is unchanged
func TestTagAttachmentUpdate_NoChange(t *testing.T) {
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

	patches.ApplyMethodFunc(tagClient, "GetResources", func(request *tag.GetResourcesRequest) (*tag.GetResourcesResponse, error) {
		resp := tag.NewGetResourcesResponse()
		resp.Response = &tag.GetResourcesResponseParams{
			ResourceTagMappingList: []*tag.ResourceTagMapping{
				{
					Resource: ptrStringTagAttachment("qcs::cvm:ap-guangzhou:uin/100020512675:instance/ins-test"),
					Tags: []*tag.Tag{
						{
							TagKey:   ptrStringTagAttachment("test_key"),
							TagValue: ptrStringTagAttachment("old_value"),
						},
					},
				},
			},
			RequestId: ptrStringTagAttachment("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaTagAttachment()
	res := svctag.ResourceTencentCloudTagAttachment()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"tag_key":   "test_key",
		"tag_value": "old_value",
		"resource":  "qcs::cvm:ap-guangzhou:uin/100020512675:instance/ins-test",
	})
	d.SetId("test_key" + tccommon.FILED_SP + "old_value" + tccommon.FILED_SP + "qcs::cvm:ap-guangzhou:uin/100020512675:instance/ins-test")

	patches.ApplyMethodFunc(d, "HasChange", func(key string) bool {
		return false
	})

	err := res.Update(d, meta)
	assert.NoError(t, err)
	assert.False(t, updateCalled, "UpdateResourceTagValue should not be called when tag_value is unchanged")
	originalId := "test_key" + tccommon.FILED_SP + "old_value" + tccommon.FILED_SP + "qcs::cvm:ap-guangzhou:uin/100020512675:instance/ins-test"
	assert.Equal(t, originalId, d.Id(), "ID should not change when tag_value is unchanged")
}
