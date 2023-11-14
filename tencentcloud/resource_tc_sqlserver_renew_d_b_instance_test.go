package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudSqlserverRenewDBInstanceResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSqlserverRenewDBInstance,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_sqlserver_renew_d_b_instance.renew_d_b_instance", "id")),
			},
			{
				ResourceName:      "tencentcloud_sqlserver_renew_d_b_instance.renew_d_b_instance",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccSqlserverRenewDBInstance = `

resource "tencentcloud_sqlserver_renew_d_b_instance" "renew_d_b_instance" {
  instance_id = "mssql-i1z41iwd"
  period = 1
  auto_voucher = 1
  voucher_ids = 
  auto_renew_flag = 0
}

`
