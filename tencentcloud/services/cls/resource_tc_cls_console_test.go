package cls_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudClsConsoleResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccClsConsole,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cls_console.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_cls_console.example", "console_id"),
					resource.TestCheckResourceAttr("tencentcloud_cls_console.example", "login_mode", "0"),
					resource.TestCheckResourceAttr("tencentcloud_cls_console.example", "domain_prefix", "tf-example"),
					resource.TestCheckResourceAttr("tencentcloud_cls_console.example", "remarks", "tf example"),
					resource.TestCheckResourceAttr("tencentcloud_cls_console.example", "access_mode.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_cls_console.example", "access_mode.0", "public"),
				),
			},
			{
				Config: testAccClsConsoleUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_cls_console.example", "remarks", "tf example updated"),
				),
			},
			{
				ResourceName:            "tencentcloud_cls_console.example",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"accounts"},
			},
		},
	})
}

const testAccClsConsole = `
resource "tencentcloud_cls_console" "example" {
  access_mode   = ["public"]
  login_mode    = 0
  domain_prefix = "tf-example"
  remarks       = "tf example"

  accounts {
    user_name = "tf_admin"
    password  = "Pass#1234"
    email     = "tf_admin@example.com"
  }

  tags {
    key   = "createdBy"
    value = "terraform"
  }
}
`

const testAccClsConsoleUpdate = `
resource "tencentcloud_cls_console" "example" {
  access_mode   = ["public"]
  login_mode    = 0
  domain_prefix = "tf-example"
  remarks       = "tf example updated"

  accounts {
    user_name = "tf_admin"
    password  = "Pass#1234"
    email     = "tf_admin@example.com"
  }

  tags {
    key   = "createdBy"
    value = "terraform"
  }
}
`
