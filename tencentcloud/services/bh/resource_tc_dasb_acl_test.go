package bh_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudNeedFixDasbAclResource_basic -v
func TestAccTencentCloudNeedFixDasbAclResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDasbAcl,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_dasb_acl.example", "id"),
					resource.TestCheckResourceAttr("tencentcloud_dasb_acl.example", "name", "tf_example"),
					resource.TestCheckResourceAttrSet("tencentcloud_dasb_acl.example", "allow_disk_redirect"),
					resource.TestCheckResourceAttrSet("tencentcloud_dasb_acl.example", "allow_any_account"),
					resource.TestCheckResourceAttrSet("tencentcloud_dasb_acl.example", "allow_clip_file_up"),
					resource.TestCheckResourceAttrSet("tencentcloud_dasb_acl.example", "allow_clip_file_down"),
					resource.TestCheckResourceAttrSet("tencentcloud_dasb_acl.example", "allow_clip_text_up"),
					resource.TestCheckResourceAttrSet("tencentcloud_dasb_acl.example", "allow_clip_text_down"),
					resource.TestCheckResourceAttrSet("tencentcloud_dasb_acl.example", "allow_file_up"),
					resource.TestCheckResourceAttrSet("tencentcloud_dasb_acl.example", "allow_file_down"),
					resource.TestCheckResourceAttrSet("tencentcloud_dasb_acl.example", "max_file_up_size"),
					resource.TestCheckResourceAttrSet("tencentcloud_dasb_acl.example", "max_file_down_size"),
					resource.TestCheckResourceAttrSet("tencentcloud_dasb_acl.example", "allow_disk_file_up"),
					resource.TestCheckResourceAttrSet("tencentcloud_dasb_acl.example", "allow_disk_file_down"),
					resource.TestCheckResourceAttrSet("tencentcloud_dasb_acl.example", "allow_shell_file_up"),
					resource.TestCheckResourceAttrSet("tencentcloud_dasb_acl.example", "allow_file_del"),
					resource.TestCheckResourceAttrSet("tencentcloud_dasb_acl.example", "allow_access_credential"),
					resource.TestCheckResourceAttr("tencentcloud_dasb_acl.example", "department_id", "1.2"),
					resource.TestCheckResourceAttr("tencentcloud_dasb_acl.example", "validate_from", "2023-09-22T00:00:00+08:00"),
					resource.TestCheckResourceAttr("tencentcloud_dasb_acl.example", "validate_to", "2024-09-23T00:00:00+08:00"),
				),
			},
			{
				ResourceName:      "tencentcloud_dasb_acl.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDasbAclUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_dasb_acl.example", "id"),
					resource.TestCheckResourceAttr("tencentcloud_dasb_acl.example", "name", "tf_example_update"),
					resource.TestCheckResourceAttrSet("tencentcloud_dasb_acl.example", "allow_disk_redirect"),
					resource.TestCheckResourceAttrSet("tencentcloud_dasb_acl.example", "allow_any_account"),
					resource.TestCheckResourceAttrSet("tencentcloud_dasb_acl.example", "allow_clip_file_up"),
					resource.TestCheckResourceAttrSet("tencentcloud_dasb_acl.example", "allow_clip_file_down"),
					resource.TestCheckResourceAttrSet("tencentcloud_dasb_acl.example", "allow_clip_text_up"),
					resource.TestCheckResourceAttrSet("tencentcloud_dasb_acl.example", "allow_clip_text_down"),
					resource.TestCheckResourceAttrSet("tencentcloud_dasb_acl.example", "allow_file_up"),
					resource.TestCheckResourceAttrSet("tencentcloud_dasb_acl.example", "allow_file_down"),
					resource.TestCheckResourceAttrSet("tencentcloud_dasb_acl.example", "max_file_up_size"),
					resource.TestCheckResourceAttrSet("tencentcloud_dasb_acl.example", "max_file_down_size"),
					resource.TestCheckResourceAttrSet("tencentcloud_dasb_acl.example", "allow_disk_file_up"),
					resource.TestCheckResourceAttrSet("tencentcloud_dasb_acl.example", "allow_disk_file_down"),
					resource.TestCheckResourceAttrSet("tencentcloud_dasb_acl.example", "allow_shell_file_up"),
					resource.TestCheckResourceAttrSet("tencentcloud_dasb_acl.example", "allow_file_del"),
					resource.TestCheckResourceAttrSet("tencentcloud_dasb_acl.example", "allow_access_credential"),
					resource.TestCheckResourceAttr("tencentcloud_dasb_acl.example", "department_id", "1.2.3"),
					resource.TestCheckResourceAttr("tencentcloud_dasb_acl.example", "validate_from", "2023-09-22T00:00:00+08:00"),
					resource.TestCheckResourceAttr("tencentcloud_dasb_acl.example", "validate_to", "2024-09-23T00:00:00+08:00"),
				),
			},
		},
	})
}

const testAccDasbAcl = `
resource "tencentcloud_dasb_acl" "example" {
  name                    = "tf_example"
  allow_disk_redirect     = true
  allow_any_account       = false
  allow_clip_file_up      = true
  allow_clip_file_down    = true
  allow_clip_text_up      = true
  allow_clip_text_down    = true
  allow_file_up           = true
  allow_file_down         = true
  max_file_up_size        = 0
  max_file_down_size      = 0
  user_id_set             = ["6", "2"]
  user_group_id_set       = ["6", "36"]
  device_id_set           = ["39", "81"]
  device_group_id_set     = ["2", "3"]
  account_set             = ["root"]
  cmd_template_id_set     = ["1", "7"]
  ac_template_id_set      = []
  allow_disk_file_up      = true
  allow_disk_file_down    = true
  allow_shell_file_up     = true
  allow_shell_file_down   = true
  allow_file_del          = true
  allow_access_credential = true
  department_id           = "1.2"
  validate_from           = "2023-09-22T00:00:00+08:00"
  validate_to             = "2024-09-23T00:00:00+08:00"
}
`

const testAccDasbAclUpdate = `
resource "tencentcloud_dasb_acl" "example" {
  name                    = "tf_example_update"
  allow_disk_redirect     = true
  allow_any_account       = false
  allow_clip_file_up      = true
  allow_clip_file_down    = true
  allow_clip_text_up      = true
  allow_clip_text_down    = true
  allow_file_up           = true
  allow_file_down         = true
  max_file_up_size        = 0
  max_file_down_size      = 0
  user_id_set             = ["6", "2"]
  user_group_id_set       = ["6", "36"]
  device_id_set           = ["39", "81"]
  device_group_id_set     = ["2", "3"]
  account_set             = ["root"]
  cmd_template_id_set     = ["1", "7"]
  ac_template_id_set      = []
  allow_disk_file_up      = true
  allow_disk_file_down    = true
  allow_shell_file_up     = true
  allow_shell_file_down   = true
  allow_file_del          = true
  allow_access_credential = true
  department_id           = "1.2.3"
  validate_from           = "2023-09-22T00:00:00+08:00"
  validate_to             = "2024-09-23T00:00:00+08:00"
}
`
