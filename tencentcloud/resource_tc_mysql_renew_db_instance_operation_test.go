package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudMysqlRenewDbInstanceOperationResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMysqlRenewDbInstanceOperation,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_mysql_renew_db_instance_operation.renew_db_instance_operation", "id")),
			},
			{
				ResourceName:      "tencentcloud_mysql_renew_db_instance_operation.renew_db_instance_operation",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccMysqlRenewDbInstanceOperation = `

resource "tencentcloud_mysql_renew_db_instance_operation" "renew_db_instance_operation" {
	instance_id = "cdb-c1nl9rpv"
	time_span = 1
	modify_pay_type = "PREPAID"
}

`
