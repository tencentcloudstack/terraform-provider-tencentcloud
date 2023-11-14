package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudMariadbRestartInstanceResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMariadbRestartInstance,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_mariadb_restart_instance.restart_instance", "id")),
			},
			{
				ResourceName:      "tencentcloud_mariadb_restart_instance.restart_instance",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccMariadbRestartInstance = `

resource "tencentcloud_mariadb_restart_instance" "restart_instance" {
  instance_ids = 
  restart_time = ""
}

`
