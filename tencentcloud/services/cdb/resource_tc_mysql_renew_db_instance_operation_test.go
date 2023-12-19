package cdb_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudMysqlRenewDbInstanceOperationResource_basic -v
func TestAccTencentCloudMysqlRenewDbInstanceOperationResource_basic(t *testing.T) {
	// t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PREPAY) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckMysqlMasterInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMysqlRenewDbInstanceOperation,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMysqlMasterInstanceExists("tencentcloud_mysql_renew_db_instance_operation.renew_db_instance_operation"),
					resource.TestCheckResourceAttrSet("tencentcloud_mysql_renew_db_instance_operation.renew_db_instance_operation", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_mysql_renew_db_instance_operation.renew_db_instance_operation", "deal_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_mysql_renew_db_instance_operation.renew_db_instance_operation", "deadline_time"),
				),
			},
		},
	})
}

const testAccMysqlRenewDbInstanceOperation = testAccMySQLPrepaid + `

resource "tencentcloud_mysql_renew_db_instance_operation" "renew_db_instance_operation" {
	instance_id = tencentcloud_mysql_instance.prepaid.id
	time_span = 1
}

`
