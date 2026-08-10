package tag_test

import (
	"context"
	"fmt"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svctag "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tag"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
	tag "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tag/v20180813"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
)

type mockMetaTagAttachment struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaTagAttachment) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaTagAttachment{}

func newMockMetaTagAttachment() *mockMetaTagAttachment {
	return &mockMetaTagAttachment{client: &connectivity.TencentCloudClient{Region: "ap-guangzhou"}}
}

func ptrStringTagAttachment(s string) *string { return &s }

// go test ./tencentcloud/services/tag/ -run "TestUnitTagAttachmentUpdate" -v -count=1 -gcflags="all=-l"
func TestUnitTagAttachmentUpdateTagValue(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	meta := newMockMetaTagAttachment()
	tagClient := &tag.Client{}
	patches.ApplyMethodReturn(meta.client, "UseTagClient", tagClient)

	updateCalled := false
	patches.ApplyMethodFunc(tagClient, "UpdateResourceTagValue", func(request *tag.UpdateResourceTagValueRequest) (*tag.UpdateResourceTagValueResponse, error) {
		updateCalled = true
		assert.Equal(t, "test_key", *request.TagKey)
		assert.Equal(t, "new_value", *request.TagValue)
		assert.Equal(t, "qcs::cvm:ap-guangzhou:uin/100020512675:instance/ins-kfrlvcp4", *request.Resource)
		resp := &tag.UpdateResourceTagValueResponse{}
		resp.Response = &tag.UpdateResourceTagValueResponseParams{
			RequestId: ptrStringTagAttachment("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(tagClient, "GetResources", func(request *tag.GetResourcesRequest) (*tag.GetResourcesResponse, error) {
		resp := &tag.GetResourcesResponse{}
		resp.Response = &tag.GetResourcesResponseParams{
			ResourceTagMappingList: []*tag.ResourceTagMapping{
				{
					Resource: ptrStringTagAttachment("qcs::cvm:ap-guangzhou:uin/100020512675:instance/ins-kfrlvcp4"),
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

	res := svctag.ResourceTencentCloudTagAttachment()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"tag_key":   "test_key",
		"tag_value": "new_value",
		"resource":  "qcs::cvm:ap-guangzhou:uin/100020512675:instance/ins-kfrlvcp4",
	})
	d.SetId("test_key" + tccommon.FILED_SP + "old_value" + tccommon.FILED_SP + "qcs::cvm:ap-guangzhou:uin/100020512675:instance/ins-kfrlvcp4")

	// force only tag_value to be detected as changed
	patches.ApplyMethodFunc(d, "HasChange", func(key string) bool {
		return key == "tag_value"
	})

	err := res.Update(d, meta)
	assert.NoError(t, err)
	assert.True(t, updateCalled)
	assert.Equal(t, "test_key"+tccommon.FILED_SP+"new_value"+tccommon.FILED_SP+"qcs::cvm:ap-guangzhou:uin/100020512675:instance/ins-kfrlvcp4", d.Id())
}

// go test ./tencentcloud/services/tag/ -run "TestUnitTagAttachmentUpdateNoChange" -v -count=1 -gcflags="all=-l"
func TestUnitTagAttachmentUpdateNoChange(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	meta := newMockMetaTagAttachment()
	tagClient := &tag.Client{}
	patches.ApplyMethodReturn(meta.client, "UseTagClient", tagClient)

	updateCalled := false
	patches.ApplyMethodFunc(tagClient, "UpdateResourceTagValue", func(request *tag.UpdateResourceTagValueRequest) (*tag.UpdateResourceTagValueResponse, error) {
		updateCalled = true
		resp := &tag.UpdateResourceTagValueResponse{}
		resp.Response = &tag.UpdateResourceTagValueResponseParams{
			RequestId: ptrStringTagAttachment("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(tagClient, "GetResources", func(request *tag.GetResourcesRequest) (*tag.GetResourcesResponse, error) {
		resp := &tag.GetResourcesResponse{}
		resp.Response = &tag.GetResourcesResponseParams{
			ResourceTagMappingList: []*tag.ResourceTagMapping{
				{
					Resource: ptrStringTagAttachment("qcs::cvm:ap-guangzhou:uin/100020512675:instance/ins-kfrlvcp4"),
					Tags: []*tag.Tag{
						{
							TagKey:   ptrStringTagAttachment("test_key"),
							TagValue: ptrStringTagAttachment("same_value"),
						},
					},
				},
			},
			RequestId: ptrStringTagAttachment("fake-request-id"),
		}
		return resp, nil
	})

	res := svctag.ResourceTencentCloudTagAttachment()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"tag_key":   "test_key",
		"tag_value": "same_value",
		"resource":  "qcs::cvm:ap-guangzhou:uin/100020512675:instance/ins-kfrlvcp4",
	})
	d.SetId("test_key" + tccommon.FILED_SP + "same_value" + tccommon.FILED_SP + "qcs::cvm:ap-guangzhou:uin/100020512675:instance/ins-kfrlvcp4")

	// no changes detected
	patches.ApplyMethodFunc(d, "HasChange", func(key string) bool {
		return false
	})

	err := res.Update(d, meta)
	assert.NoError(t, err)
	assert.False(t, updateCalled)
	assert.Equal(t, "test_key"+tccommon.FILED_SP+"same_value"+tccommon.FILED_SP+"qcs::cvm:ap-guangzhou:uin/100020512675:instance/ins-kfrlvcp4", d.Id())
}

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
