package bh_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudBhAclResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccBhAcl,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_bh_acl.example", "id"),
					resource.TestCheckResourceAttr("tencentcloud_bh_acl.example", "name", "tf-example"),
					resource.TestCheckResourceAttr("tencentcloud_bh_acl.example", "allow_disk_redirect", "true"),
					resource.TestCheckResourceAttr("tencentcloud_bh_acl.example", "allow_any_account", "false"),
				),
			},
			{
				Config: testAccBhAclUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_bh_acl.example", "id"),
					resource.TestCheckResourceAttr("tencentcloud_bh_acl.example", "name", "tf-example"),
					resource.TestCheckResourceAttr("tencentcloud_bh_acl.example", "allow_clip_file_up", "true"),
				),
			},
			{
				ResourceName:      "tencentcloud_bh_acl.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccBhAcl = `
resource "tencentcloud_bh_acl" "example" {
  name                    = "tf-example"
  allow_disk_redirect     = true
  allow_any_account       = false
  allow_clip_file_up      = true
  allow_clip_file_down    = false
  allow_clip_text_up      = false
  allow_clip_text_down    = true
  allow_file_up           = false
  allow_file_down         = true
  allow_disk_file_up      = true
  allow_disk_file_down    = false
  allow_shell_file_up     = false
  allow_shell_file_down   = false
  allow_file_del          = false
  allow_access_credential = true
  allow_keyboard_logger   = false
}
`

const testAccBhAclUpdate = `
resource "tencentcloud_bh_acl" "example" {
  name                    = "tf-example"
  allow_disk_redirect     = true
  allow_any_account       = false
  allow_clip_file_up      = false
  allow_clip_file_down    = false
  allow_clip_text_up      = false
  allow_clip_text_down    = false
  allow_file_up           = false
  allow_file_down         = false
  allow_disk_file_up      = false
  allow_disk_file_down    = false
  allow_shell_file_up     = false
  allow_shell_file_down   = false
  allow_file_del          = false
  allow_access_credential = true
  allow_keyboard_logger   = false
}
`
