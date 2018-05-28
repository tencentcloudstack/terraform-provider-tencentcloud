package tencentcloud

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccTencentCloudRouteEntry_basic(t *testing.T) {
	var reId string

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreSetRegion("ap-guangzhou")
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRouteEntryDestroy(&reId),
		Steps: []resource.TestStep{
			{
				Config: testAccRouteEntryConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRouteEntryExists("tencentcloud_route_entry.instance", &reId),
					resource.TestCheckResourceAttr("tencentcloud_route_entry.instance", "next_type", "instance"),
				),
			},
		},
	})
}

func testAccCheckRouteEntryDestroy(id *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := testAccProvider.Meta().(*TencentCloudClient).commonConn
		params, _ := routeIdDecode(*id)
		params["Action"] = "DescribeRoutes"
		response, err := conn.SendRequest("vpc", params)
		if err != nil {
			return err
		}
		var jsonresp struct {
			Code     int    `json:"code"`
			Message  string `json:"message"`
			CodeDesc string `json:"codeDesc"`
			Data     struct {
				TotalNum int `json:"totalNum"`
			} `json:"data"`
		}
		err = json.Unmarshal([]byte(response), &jsonresp)
		if err != nil {
			return err
		}
		if jsonresp.Data.TotalNum != 1 {
			return fmt.Errorf("Route entry still exists: %s", *id)
		}
		return nil
	}
}

func testAccCheckRouteEntryExists(n string, id *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No route entry ID is set")
		}

		conn := testAccProvider.Meta().(*TencentCloudClient).commonConn
		params, _ := routeIdDecode(rs.Primary.ID)
		params["Action"] = "DescribeRoutes"
		response, err := conn.SendRequest("vpc", params)
		if err != nil {
			return err
		}
		var jsonresp struct {
			Code     int    `json:"code"`
			Message  string `json:"message"`
			CodeDesc string `json:"codeDesc"`
			Data     struct {
				TotalNum int `json:"totalNum"`
			} `json:"data"`
		}
		err = json.Unmarshal([]byte(response), &jsonresp)
		if err != nil {
			return err
		}
		if jsonresp.Data.TotalNum != 2 {
			return fmt.Errorf("Route entry not found: %s", rs.Primary.ID)
		}

		*id = rs.Primary.ID
		return nil
	}
}

const testAccRouteEntryConfig = `
data "tencentcloud_image" "foo" {
  filter {
    name   = "image-type"
    values = ["PUBLIC_IMAGE"]
  }
}

resource "tencentcloud_vpc" "foo" {
    name = "ci-temp-test"
    cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_subnet" "foo" {
  vpc_id            = "${tencentcloud_vpc.foo.id}"
  name              = "terraform test subnet"
  cidr_block        = "10.0.12.0/24"
  availability_zone = "ap-guangzhou-3"
}

resource "tencentcloud_route_table" "foo" {
    vpc_id = "${tencentcloud_vpc.foo.id}"
    name = "ci-temp-test-rt"
}

resource "tencentcloud_instance" "foo" {
  availability_zone = "ap-guangzhou-3"
  image_id          = "${data.tencentcloud_image.foo.image_id}"
  vpc_id            = "${tencentcloud_vpc.foo.id}"
  subnet_id         = "${tencentcloud_subnet.foo.id}"
  instance_type     = "S2.SMALL2"
}

resource "tencentcloud_route_entry" "instance" {
  vpc_id = "${tencentcloud_vpc.foo.id}"
  route_table_id = "${tencentcloud_route_table.foo.id}"
  cidr_block     = "10.4.4.0/24"
  next_type      = "instance"
  next_hub       = "${tencentcloud_instance.foo.private_ip}"
}
`
