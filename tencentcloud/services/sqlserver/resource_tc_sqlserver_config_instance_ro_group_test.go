package sqlserver_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudSqlserverConfigInstanceRoGroupResource_basic -v
func TestAccTencentCloudSqlserverConfigInstanceRoGroupResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSqlserverConfigInstanceRoGroup,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_sqlserver_config_instance_ro_group.example", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_sqlserver_config_instance_ro_group.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccSqlserverConfigInstanceRoGroup = tcacctest.DefaultVpcSubnets + tcacctest.DefaultSecurityGroupData + `
data "tencentcloud_availability_zones_by_product" "zones" {
  product = "sqlserver"
}

resource "tencentcloud_sqlserver_basic_instance" "example" {
  name                   = "tf-example"
  availability_zone      = data.tencentcloud_availability_zones_by_product.zones.zones.4.name
  charge_type            = "POSTPAID_BY_HOUR"
  vpc_id                 = local.vpc_id
  subnet_id              = local.subnet_id
  project_id             = 0
  memory                 = 4
  storage                = 100
  cpu                    = 2
  machine_type           = "CLOUD_PREMIUM"
  maintenance_week_set   = [1, 2, 3]
  maintenance_start_time = "09:00"
  maintenance_time_span  = 3
  security_groups        = [local.sg_id]

  tags = {
    "test" = "test"
  }
}

resource "tencentcloud_sqlserver_readonly_instance" "example" {
  name                     = "tf_example"
  availability_zone        = data.tencentcloud_availability_zones_by_product.zones.zones.4.name
  charge_type              = "POSTPAID_BY_HOUR"
  vpc_id                   = local.vpc_id
  subnet_id                = local.subnet_id
  memory                   = 4
  storage                  = 20
  master_instance_id       = tencentcloud_sqlserver_basic_instance.example.id
  readonly_group_type      = 2
  read_only_group_name     = "tf_example_ro"
  is_offline_delay         = 1
  read_only_max_delay_time = 10
  min_read_only_in_group   = 0
  force_upgrade            = true
}

resource "tencentcloud_sqlserver_config_instance_ro_group" "example" {
  instance_id              = tencentcloud_sqlserver_readonly_instance.example.master_instance_id
  read_only_group_id       = tencentcloud_sqlserver_readonly_instance.example.readonly_group_id
  read_only_group_name     = "tf_example_ro_update"
  is_offline_delay         = 1
  read_only_max_delay_time = 5
  min_read_only_in_group   = 1
}
`
