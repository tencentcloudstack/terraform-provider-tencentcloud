package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

var testDataTcaplusApplicationsName = "data.tencentcloud_tcaplus_applications.id_test"

func TestAccTencentCloudDataTcaplusApplications(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTcaplusApplicationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudDataTcaplusApplicationsBaic,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTcaplusApplicationExists("tencentcloud_tcaplus_application.test_app"),
					resource.TestCheckResourceAttrSet(testDataTcaplusApplicationsName, "app_id"),
					resource.TestCheckResourceAttr(testDataTcaplusApplicationsName, "list.#", "1"),
					resource.TestCheckResourceAttr(testDataTcaplusApplicationsName, "list.0.app_name", "tf_tcaplus_data_guagua"),
					resource.TestCheckResourceAttr(testDataTcaplusApplicationsName, "list.0.idl_type", "PROTO"),
					resource.TestCheckResourceAttr(testDataTcaplusApplicationsName, "list.0.password", "1qaA2k1wgvfa3ZZZ"),
					resource.TestCheckResourceAttrSet(testDataTcaplusApplicationsName, "list.0.network_type"),
					resource.TestCheckResourceAttrSet(testDataTcaplusApplicationsName, "list.0.create_time"),
					resource.TestCheckResourceAttrSet(testDataTcaplusApplicationsName, "list.0.password_status"),
					resource.TestCheckResourceAttrSet(testDataTcaplusApplicationsName, "list.0.api_access_id"),
					resource.TestCheckResourceAttrSet(testDataTcaplusApplicationsName, "list.0.api_access_ip"),
					resource.TestCheckResourceAttrSet(testDataTcaplusApplicationsName, "list.0.api_access_port"),
				),
			},
		},
	})
}

const testAccTencentCloudDataTcaplusApplicationsBaic = `
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
data "tencentcloud_tcaplus_applications" "id_test" {
  app_id = tencentcloud_tcaplus_application.test_app.id
}
`
