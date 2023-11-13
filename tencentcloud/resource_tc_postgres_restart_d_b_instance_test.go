package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudPostgresRestartDBInstanceResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPostgresRestartDBInstance,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_postgres_restart_d_b_instance.restart_d_b_instance", "id")),
			},
			{
				ResourceName:      "tencentcloud_postgres_restart_d_b_instance.restart_d_b_instance",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccPostgresRestartDBInstance = `

resource "tencentcloud_postgres_restart_d_b_instance" "restart_d_b_instance" {
  d_b_instance_id = "postgres-6r233v55"
  tags = {
    "createdBy" = "terraform"
  }
}

`
