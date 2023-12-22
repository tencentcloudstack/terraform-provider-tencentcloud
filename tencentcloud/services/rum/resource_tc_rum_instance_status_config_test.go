package rum_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -test.run TestAccTencentCloudRumInstanceStatusConfigResource_basic -v
func TestAccTencentCloudRumInstanceStatusConfigResource_basic(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccRumInstanceStatusConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_rum_instance_status_config.instance_status_config", "id"),
					resource.TestCheckResourceAttr("tencentcloud_rum_instance_status_config.instance_status_config", "instance_status", "6"),
				),
			},
			{
				Config: testAccRumInstanceStatusConfigUp,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_rum_instance_status_config.instance_status_config", "id"),
					resource.TestCheckResourceAttr("tencentcloud_rum_instance_status_config.instance_status_config", "instance_status", "2"),
				),
			},
			{
				ResourceName:            "tencentcloud_rum_instance_status_config.instance_status_config",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"operate"},
			},
		},
	})
}

const testAccRumInstanceStatusConfigVar = `
variable "instance_id" {
  default = "` + tcacctest.DefaultRumInstanceId + `"
}
`

const testAccRumInstanceStatusConfig = testAccRumInstanceStatusConfigVar + `

resource "tencentcloud_rum_instance_status_config" "instance_status_config" {
	instance_id = var.instance_id
	operate = "stop"
}

`

const testAccRumInstanceStatusConfigUp = testAccRumInstanceStatusConfigVar + `

resource "tencentcloud_rum_instance_status_config" "instance_status_config" {
	instance_id = var.instance_id
	operate = "resume"
}

`
