package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudTagResourceTagResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTagResourceTag,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_tag_resource_tag.resource_tag", "id")),
			},
			{
				ResourceName:      "tencentcloud_tag_resource_tag.resource_tag",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTagResourceTag = `

resource "tencentcloud_tag_resource_tag" "resource_tag" {
  tag_key = "test3"
  tag_value = "Terraform3"
  resource = "qcs::cvm:ap-guangzhou:uin/100020512675:instance/ins-kfrlvcp4"
}

`
