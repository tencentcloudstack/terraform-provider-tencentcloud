package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

var testDataTcaplusZonesName = "data.tencentcloud_tcaplus_zones.id_test"

func TestAccTencentCloudDataTcaplusZones(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTcaplusZoneDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudDataTcaplusZonesBaic,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTcaplusZoneExists("tencentcloud_tcaplus_zone.test_zone"),
					resource.TestCheckResourceAttrSet(testDataTcaplusZonesName, "app_id"),
					resource.TestCheckResourceAttrSet(testDataTcaplusZonesName, "zone_id"),
					resource.TestCheckResourceAttr(testDataTcaplusZonesName, "list.#", "1"),
					resource.TestCheckResourceAttr(testDataTcaplusZonesName, "list.0.zone_name", "tf_test_zone_name_guagua"),
					resource.TestCheckResourceAttr(testDataTcaplusZonesName, "list.0.table_count", "0"),
					resource.TestCheckResourceAttrSet(testDataTcaplusZonesName, "list.0.zone_id"),
					resource.TestCheckResourceAttrSet(testDataTcaplusZonesName, "list.0.total_size"),
					resource.TestCheckResourceAttrSet(testDataTcaplusZonesName, "list.0.create_time"),
				),
			},
		},
	})
}

const testAccTencentCloudDataTcaplusZonesBaic = `
variable "availability_zone" {
default = "ap-shanghai-2"
}

variable "instance_name" {
default = "` + defaultInsName + `"
}
variable "vpc_cidr" {
default = "` + defaultVpcCidr + `"
}
variable "subnet_cidr" {
default = "` + defaultSubnetCidr + `"
}

resource "tencentcloud_vpc" "foo" {
name       = var.instance_name
cidr_block = var.vpc_cidr
}

resource "tencentcloud_subnet" "subnet" {
name              = var.instance_name
vpc_id            = tencentcloud_vpc.foo.id
availability_zone = var.availability_zone
cidr_block        = var.subnet_cidr
is_multicast      = false
}
resource "tencentcloud_tcaplus_application" "test_app" {
  idl_type                 = "PROTO"
  app_name                 = "tf_tcaplus_data_guagua"
  vpc_id                   = tencentcloud_vpc.foo.id
  subnet_id                = tencentcloud_subnet.subnet.id
  password                 = "1qaA2k1wgvfa3ZZZ"
  old_password_expire_last = 3600
}
resource "tencentcloud_tcaplus_zone" "test_zone" {
  app_id    = tencentcloud_tcaplus_application.test_app.id
  zone_name = "tf_test_zone_name_guagua"
}

data "tencentcloud_tcaplus_zones" "id_test" {
   app_id    = tencentcloud_tcaplus_application.test_app.id
   zone_id   = tencentcloud_tcaplus_zone.test_zone.id
}
`
