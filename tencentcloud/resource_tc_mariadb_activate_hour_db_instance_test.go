package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudMariadbActivateHourDbInstanceResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMariadbActivateHourDbInstance,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_mariadb_activate_hour_db_instance.activate_hour_db_instance", "id")),
			},
			{
				ResourceName:      "tencentcloud_mariadb_activate_hour_db_instance.activate_hour_db_instance",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccMariadbActivateHourDbInstance = `

resource "tencentcloud_mariadb_activate_hour_db_instance" "activate_hour_db_instance" {
  instance_ids = 
}

`
