package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudMariadbDbInstance_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMariadbDbInstance,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_mariadb_db_instance.db_instance", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_mariadb_db_instance.dbInstance",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccMariadbDbInstance = `

resource "tencentcloud_mariadb_db_instance" "db_instance" {
  zones = ""
  node_count = ""
  memory = ""
  storage = ""
  count = ""
  vpc_id = ""
  subnet_id = ""
  db_version_id = ""
  instance_name = ""
  tags = {
    "createdBy" = "terraform"
  }
}

`
