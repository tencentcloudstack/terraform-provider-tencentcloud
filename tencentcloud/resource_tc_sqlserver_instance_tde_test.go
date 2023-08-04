package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudNeedFixSqlserverInstanceTDEResource_basic -v
func TestAccTencentCloudNeedFixSqlserverInstanceTDEResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSqlserverInstanceTDE,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_sqlserver_instance_tde.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_sqlserver_instance_tde.example", "instance_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_sqlserver_instance_tde.example", "certificate_attribution"),
					resource.TestCheckResourceAttrSet("tencentcloud_sqlserver_instance_tde.example", "quote_uin"),
				),
			},
			{
				ResourceName:      "tencentcloud_sqlserver_instance_tde.instance_tde",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccSqlserverInstanceTDE = defaultVpcSubnets + defaultSecurityGroupData + `
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

resource "tencentcloud_sqlserver_instance_tde" "instance_tde" {
  instance_id             = tencentcloud_sqlserver_basic_instance.example.id
  certificate_attribution = "self"
}
`
