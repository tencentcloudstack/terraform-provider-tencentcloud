package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudPostgresqlRestartDbInstanceOperationResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPostgresqlRestartDbInstanceOperation,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_postgresql_restart_db_instance_operation.restart_db_instance_operation", "id")),
			},
			{
				ResourceName:      "tencentcloud_postgresql_restart_db_instance_operation.restart_db_instance_operation",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccPostgresqlRestartDbInstanceOperation = `

resource "tencentcloud_postgresql_restart_db_instance_operation" "restart_db_instance_operation" {
  db_instance_id = "postgres-6r233v55"
}

`
