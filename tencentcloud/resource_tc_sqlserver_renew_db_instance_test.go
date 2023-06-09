package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudSqlserverRenewDBInstanceResource_basic -v
func TestAccTencentCloudSqlserverRenewDBInstanceResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY) },
		CheckDestroy: testAccCheckSqlserverInstanceDestroy,
		Providers:    testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSqlserverRenewDBInstance,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_sqlserver_renew_db_instance.renew_db_instance", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_sqlserver_renew_db_instance.renew_db_instance",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccSqlserverRenewDBInstance = testAccSqlserverInstanceBasicPrepaid + `
resource "tencentcloud_sqlserver_instance" "test" {
  name                          = "test_sqlserver_instance"
  availability_zone             = local.az
  charge_type                   = "PREPAID"
  period                        = 1
  vpc_id                        = local.vpc_id
  subnet_id                     = local.vpc_subnet_id
  project_id                    = 0
  memory                        = 2
  storage                       = 10
  maintenance_week_set          = [1,2,3]
  maintenance_start_time        = "09:00"
  maintenance_time_span         = 3
  security_groups               = [local.sg]
}

resource "tencentcloud_sqlserver_renew_db_instance" "renew_db_instance" {
  instance_id = tencentcloud_sqlserver_instance.test.id
  period = 1
}
`
