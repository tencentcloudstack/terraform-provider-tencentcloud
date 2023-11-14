package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudPostgresModifyDBInstanceChargeTypeResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPostgresModifyDBInstanceChargeType,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_postgres_modify_d_b_instance_charge_type.modify_d_b_instance_charge_type", "id")),
			},
			{
				ResourceName:      "tencentcloud_postgres_modify_d_b_instance_charge_type.modify_d_b_instance_charge_type",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccPostgresModifyDBInstanceChargeType = `

resource "tencentcloud_postgres_modify_d_b_instance_charge_type" "modify_d_b_instance_charge_type" {
  d_b_instance_id = "postgres-6r233v55"
  instance_charge_type = "PREPAID"
  period = 1
  auto_renew_flag = 0
  auto_voucher = 0
  tags = {
    "createdBy" = "terraform"
  }
}

`
