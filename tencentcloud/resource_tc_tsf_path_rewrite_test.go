package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudNeedFixTsfPathRewriteResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTsfPathRewrite,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_tsf_path_rewrite.path_rewrite", "id")),
			},
			{
				ResourceName:      "tencentcloud_tsf_path_rewrite.path_rewrite",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTsfPathRewrite = `

resource "tencentcloud_tsf_path_rewrite" "path_rewrite" {
    gateway_group_id = ""
  regex = ""
  replacement = ""
  blocked = ""
  order = 
}

`
