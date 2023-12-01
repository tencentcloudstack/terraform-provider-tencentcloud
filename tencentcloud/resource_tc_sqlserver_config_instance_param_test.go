package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudSqlserverConfigInstanceParamResource_basic -v
func TestAccTencentCloudSqlserverConfigInstanceParamResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		CheckDestroy: testAccCheckSqlserverInstanceDestroy,
		Providers:    testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSqlserverConfigInstanceParam,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_sqlserver_config_instance_param.example", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_sqlserver_config_instance_param.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccSqlserverConfigInstanceParam = defaultVpcSubnets + defaultSecurityGroupData + `
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

resource "tencentcloud_sqlserver_config_instance_param" "example" {
  instance_id = tencentcloud_sqlserver_basic_instance.example.id
  param_list {
    name          = "fill factor(%)"
    current_value = "90"
  }
}
`
