package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudSqlserverRenewPostpaidDBInstanceResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSqlserverRenewPostpaidDBInstance,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_sqlserver_renew_postpaid_d_b_instance.renew_postpaid_d_b_instance", "id")),
			},
			{
				ResourceName:      "tencentcloud_sqlserver_renew_postpaid_d_b_instance.renew_postpaid_d_b_instance",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccSqlserverRenewPostpaidDBInstance = `

resource "tencentcloud_sqlserver_renew_postpaid_d_b_instance" "renew_postpaid_d_b_instance" {
  instance_id = "mssql-i1z41iwd"
}

`
