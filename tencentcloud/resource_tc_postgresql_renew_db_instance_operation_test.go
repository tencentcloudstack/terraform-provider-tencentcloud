package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudPostgresqlRenewDbInstanceOperationResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPostgresqlRenewDbInstanceOperation,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_postgresql_renew_db_instance_operation.renew_db_instance_operation", "id")),
			},
			{
				ResourceName:      "tencentcloud_postgresql_renew_db_instance_operation.renew_db_instance_operation",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccPostgresqlRenewDbInstanceOperation = `

resource "tencentcloud_postgresql_renew_db_instance_operation" "renew_db_instance_operation" {
  db_instance_id = "postgres-6fego161"
  period = 12
  auto_voucher = 0
  voucher_ids = &lt;nil&gt;
}

`
