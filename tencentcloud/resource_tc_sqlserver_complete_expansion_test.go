package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudSqlserverCompleteExpansionResource_basic -v
// go test -v -run TestAccTencentCloudSqlserverCompleteExpansionResource_basic -timeout=0
func TestAccTencentCloudSqlserverCompleteExpansionResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		CheckDestroy: testAccCheckSqlserverInstanceDestroy,
		Providers:    testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNewSqlserverInstance,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSqlserverInstanceExists(testSqlserverInstanceResourceKey),
					resource.TestCheckResourceAttrSet(testSqlserverInstanceResourceKey, "id"),
				),
			},
			{
				Config: testAccUpdateNewSqlserverInstance,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSqlserverInstanceExists(testSqlserverInstanceResourceKey),
					resource.TestCheckResourceAttrSet(testSqlserverInstanceResourceKey, "id"),
				),
			},
		},
	})
}

const testAccNewSqlserverInstance string = defaultVpcVariable + defaultSecurityGroupData + `
data "tencentcloud_availability_zones_by_product" "zones" {
  product = "sqlserver"
}

resource "tencentcloud_sqlserver_instance" "example" {
  name                   = "tf_example_sql"
  availability_zone      = data.tencentcloud_availability_zones_by_product.zones.zones.4.name
  charge_type            = "POSTPAID_BY_HOUR"
  vpc_id                 = var.vpc_id
  subnet_id              = var.subnet_id
  security_groups        = [local.sg_id]
  project_id             = 0
  memory                 = 2
  storage                = 20
  maintenance_week_set   = [1, 2, 3]
  maintenance_start_time = "01:00"
  maintenance_time_span  = 3
  tags                   = {
    "createBy" = "tfExample"
  }
}
`

const testAccUpdateNewSqlserverInstance string = defaultVpcSubnets + defaultSecurityGroupData + `
data "tencentcloud_availability_zones_by_product" "zones" {
  product = "sqlserver"
}

resource "tencentcloud_sqlserver_instance" "example" {
  name                   = "tf_example_sql"
  availability_zone      = data.tencentcloud_availability_zones_by_product.zones.zones.4.name
  charge_type            = "POSTPAID_BY_HOUR"
  vpc_id                 = var.vpc_id
  subnet_id              = var.subnet_id
  security_groups        = [local.sg_id]
  project_id             = 0
  memory                 = 4
  storage                = 40
  maintenance_week_set   = [1, 2, 3]
  maintenance_start_time = "01:00"
  maintenance_time_span  = 3
  wait_switch            = 1
  tags                   = {
    "createBy" = "tfExample"
  }
}

resource "tencentcloud_sqlserver_complete_expansion" "example" {
  instance_id = tencentcloud_sqlserver_instance.example.id
}
`
