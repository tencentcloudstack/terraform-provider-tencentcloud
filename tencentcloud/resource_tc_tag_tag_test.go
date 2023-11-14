package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudTagTagResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTagTag,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_tag_tag.tag", "id")),
			},
			{
				ResourceName:      "tencentcloud_tag_tag.tag",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTagTag = `

resource "tencentcloud_tag_tag" "tag" {
  tag_key = "test"
  tag_value = "Terraform"
}

`
