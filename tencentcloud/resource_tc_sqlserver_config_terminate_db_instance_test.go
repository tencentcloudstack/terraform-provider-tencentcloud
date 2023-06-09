package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudSqlserverConfigTerminateDBInstanceResource_basic -v
func TestAccTencentCloudSqlserverConfigTerminateDBInstanceResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		CheckDestroy: testAccCheckSqlserverInstanceDestroy,
		Providers:    testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSqlserverConfigTerminateDBInstance,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_sqlserver_config_terminate_db_instance.config_terminate_db_instance", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_sqlserver_config_terminate_db_instance.config_terminate_db_instance",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccSqlserverConfigTerminateDBInstance = testAccSqlserverBasicInstanceNetwork + `
resource "tencentcloud_sqlserver_instance" "test" {
  name                          = "tf_sqlserver_instance"
  availability_zone             = var.default_az
  charge_type                   = "POSTPAID_BY_HOUR"
  vpc_id                        = local.vpc_id
  subnet_id                     = local.subnet_id
  security_groups               = [local.sg_id]
  project_id                    = 0
  memory                        = 2
  storage                       = 10
  maintenance_week_set          = [1,2,3]
  maintenance_start_time        = "09:00"
  maintenance_time_span         = 3
  tags = {
    "test"                      = "test"
  }
}

resource "tencentcloud_sqlserver_config_terminate_db_instance" "config_terminate_db_instance" {
  instance_id = tencentcloud_sqlserver_instance.test.id
}

resource "tencentcloud_sqlserver_renew_postpaid_db_instance" "renew_postpaid_db_instance" {
  instance_id = tencentcloud_sqlserver_instance.test.id
}
`
