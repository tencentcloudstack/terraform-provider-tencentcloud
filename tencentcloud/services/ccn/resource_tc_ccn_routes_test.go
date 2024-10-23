package ccn_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudCcnRoutesResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcCcnRoutes,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_ccn_routes.ccn_routes", "id")),
			},
			{
				Config: testAccVpcCcnRoutesUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_ccn_routes.ccn_routes", "id"),
					resource.TestCheckResourceAttr("tencentcloud_ccn_routes.ccn_routes", "switch", "on"),
				),
			},
			{
				ResourceName: "tencentcloud_ccn_routes.ccn_routes",
				ImportState:  true,
			},
		},
	})
}

const testAccVpcCcnRoutes = `
resource "tencentcloud_ccn" "ccn" {
  name                 = "tf-test-routes"
  description          = "tf-test-routes-des"
  qos                  = "AG"
}

resource "tencentcloud_ccn_route_table" "table" {
  ccn_id      = tencentcloud_ccn.ccn.id
  name        = "tf-test-routes"
  description = "desc."
}

data "tencentcloud_availability_zones" "zones" {}

resource "tencentcloud_vpc" "vpc" {
  name       = "vpc-ccn-routes"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_subnet" "subnet" {
  vpc_id            = tencentcloud_vpc.vpc.id
  name              = "subnet-ccn-routes"
  cidr_block        = "10.0.0.0/16"
  availability_zone = data.tencentcloud_availability_zones.zones.zones.0.name
  is_multicast      = false
}

resource "tencentcloud_ccn_attachment" "attachment" {
  ccn_id          = tencentcloud_ccn.ccn.id
  instance_type   = "VPC"
  instance_id     = tencentcloud_vpc.vpc.id
  instance_region = "ap-guangzhou"
  route_table_id = tencentcloud_ccn_route_table.table.id
}

data "tencentcloud_ccn_routes" "routes" {
  ccn_id = tencentcloud_ccn.ccn.id
  depends_on = [tencentcloud_ccn_attachment.attachment]
}

resource "tencentcloud_ccn_routes" "ccn_routes" {
  ccn_id = tencentcloud_ccn.ccn.id
  route_id = data.tencentcloud_ccn_routes.routes.route_list.0.route_id
  switch = "off"
}

`

const testAccVpcCcnRoutesUpdate = `
resource "tencentcloud_ccn" "ccn" {
  name                 = "tf-test-routes"
  description          = "tf-test-routes-des"
  qos                  = "AG"
}

resource "tencentcloud_ccn_route_table" "table" {
  ccn_id      = tencentcloud_ccn.ccn.id
  name        = "tf-test-routes"
  description = "desc."
}

data "tencentcloud_availability_zones" "zones" {}

resource "tencentcloud_vpc" "vpc" {
  name       = "vpc-ccn-routes"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_subnet" "subnet" {
  vpc_id            = tencentcloud_vpc.vpc.id
  name              = "subnet-ccn-routes"
  cidr_block        = "10.0.0.0/16"
  availability_zone = data.tencentcloud_availability_zones.zones.zones.0.name
  is_multicast      = false
}

resource "tencentcloud_ccn_attachment" "attachment" {
  ccn_id          = tencentcloud_ccn.ccn.id
  instance_type   = "VPC"
  instance_id     = tencentcloud_vpc.vpc.id
  instance_region = "ap-guangzhou"
  route_table_id = tencentcloud_ccn_route_table.table.id
}

data "tencentcloud_ccn_routes" "routes" {
  ccn_id = tencentcloud_ccn.ccn.id
  depends_on = [tencentcloud_ccn_attachment.attachment]
}

resource "tencentcloud_ccn_routes" "ccn_routes" {
  ccn_id = tencentcloud_ccn.ccn.id
  route_id = data.tencentcloud_ccn_routes.routes.route_list.0.route_id
  switch = "on"
}

`
