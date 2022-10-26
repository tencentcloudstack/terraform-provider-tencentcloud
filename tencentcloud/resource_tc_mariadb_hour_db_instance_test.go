package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudMariadbHourDbInstance_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMariadbHourDbInstance,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_mariadb_hour_db_instance.hour_db_instance", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_mariadb_hour_db_instance.hourDbInstance",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccMariadbHourDbInstance = `

resource "tencentcloud_mariadb_hour_db_instance" "hour_db_instance" {
  zones = ""
  node_count = ""
  memory = ""
  storage = ""
  vpc_id = ""
  subnet_id = ""
  db_version_id = ""
  instance_name = ""
  tags = {
    "createdBy" = "terraform"
  }
}

`
