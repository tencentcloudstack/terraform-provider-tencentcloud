package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCynosdbRestartInstanceResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCynosdbRestartInstance,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_restart_instance.restart_instance", "id")),
			},
			{
				ResourceName:      "tencentcloud_cynosdb_restart_instance.restart_instance",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCynosdbRestartInstance = `

resource "tencentcloud_cynosdb_restart_instance" "restart_instance" {
  instance_id = "cynosdbmysql-ins-xxxxxxxx"
}

`
