package bh_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudBhUserDirectoryResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccBhUserDirectory,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_bh_user_directory.example", "id"),
				),
			},
			{
				Config: testAccBhUserDirectoryUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_bh_user_directory.example", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_bh_user_directory.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccBhUserDirectory = `
resource "tencentcloud_bh_user_directory" "example" {
  dir_id   = 895784
  dir_name = "tf-example"
  user_org_set {
    org_id        = 1576799
    org_name      = "orgName"
    org_id_path   = "819729.895784"
    org_name_path = "Root.demo"
    user_total    = 0
  }
  source      = 0
  source_name = "sourceName"
  user_count  = 3
}
`

const testAccBhUserDirectoryUpdate = `
resource "tencentcloud_bh_user_directory" "example" {
  dir_id   = 895784
  dir_name = "tf-example"
  user_org_set {
    org_id        = 1576799
    org_name      = "orgName"
    org_id_path   = "819729.895784"
    org_name_path = "Root.demo"
    user_total    = 0
  }
  source      = 0
  source_name = "sourceName"
  user_count  = 3
}
`
