package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudPostgresqlDisisolateDbInstanceOperationResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPostgresqlDisisolateDbInstanceOperation,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_postgresql_disisolate_db_instance_operation.disisolate_db_instance_operation", "id")),
			},
			{
				ResourceName:      "tencentcloud_postgresql_disisolate_db_instance_operation.disisolate_db_instance_operation",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccPostgresqlDisisolateDbInstanceOperation = `

resource "tencentcloud_postgresql_disisolate_db_instance_operation" "disisolate_db_instance_operation" {
  d_b_instance_id_set = &lt;nil&gt;
  period = 12
  auto_voucher = false
  voucher_ids = &lt;nil&gt;
}

`
