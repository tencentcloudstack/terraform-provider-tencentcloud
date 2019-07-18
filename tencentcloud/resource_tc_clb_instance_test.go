package tencentcloud

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccTencentCloudClbInstance_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckClbInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccClbInstance_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClbInstanceExists("tencentcloud_clb_instance.clb_basic"),
					resource.TestCheckResourceAttr("tencentcloud_clb_instance.clb_basic", "net_type", "OPEN"),
					resource.TestCheckResourceAttr("tencentcloud_clb_instance.clb_basic", "clb_name", "tf-clb-basic"),
				),
			},
			{
				ResourceName:      "tencentcloud_clb_instance.clb_basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccTencentCloudClbInstance_open(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckClbInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccClbInstance_open,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClbInstanceExists("tencentcloud_clb_instance.clb_open"),
					resource.TestCheckResourceAttr("tencentcloud_clb_instance.clb_open", "net_type", "OPEN"),
					resource.TestCheckResourceAttr("tencentcloud_clb_instance.clb_open", "clb_name", "tf-clb-open"),
					resource.TestCheckResourceAttr("tencentcloud_clb_instance.clb_open", "project_id", "0"),
					resource.TestCheckResourceAttr("tencentcloud_clb_instance.clb_open", "security_groups.#", "1"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_instance.clb_open", "security_groups.0"),
					resource.TestCheckResourceAttr("tencentcloud_clb_instance.clb_open", "target_region_info.region", "ap-guangzhou"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_instance.clb_open", "target_region_info.vpc_id"),
				),
			},
			{
				Config: testAccClbInstance_update_open,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClbInstanceExists("tencentcloud_clb_instance.clb_open"),
					resource.TestCheckResourceAttr("tencentcloud_clb_instance.clb_open", "clb_name", "tf-clb-update-open"),
					resource.TestCheckResourceAttr("tencentcloud_clb_instance.clb_open", "net_type", "OPEN"),
					resource.TestCheckResourceAttr("tencentcloud_clb_instance.clb_open", "project_id", "0"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_instance.clb_open", "vpc_id"),
					resource.TestCheckResourceAttr("tencentcloud_clb_instance.clb_open", "security_groups.#", "1"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_instance.clb_open", "security_groups.0"),
					resource.TestCheckResourceAttr("tencentcloud_clb_instance.clb_open", "target_region_info.region", "ap-guangzhou"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_instance.clb_open", "target_region_info.vpc_id"),
				),
			},
		},
	})
}

func TestAccTencentCloudClbInstance_internal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckClbInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccClbInstance_internal,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClbInstanceExists("tencentcloud_clb_instance.clb_internal"),
					resource.TestCheckResourceAttr("tencentcloud_clb_instance.clb_internal", "clb_name", "tf-clb-internal"),
					resource.TestCheckResourceAttr("tencentcloud_clb_instance.clb_internal", "net_type", "INTERNAL"),
					resource.TestCheckResourceAttr("tencentcloud_clb_instance.clb_internal", "project_id", "0"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_instance.clb_internal", "vpc_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_instance.clb_internal", "subnet_id"),
				),
			},
			{
				Config: testAccClbInstance_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClbInstanceExists("tencentcloud_clb_instance.clb_internal"),
					resource.TestCheckResourceAttr("tencentcloud_clb_instance.clb_internal", "clb_name", "tf-clb-update"),
					resource.TestCheckResourceAttr("tencentcloud_clb_instance.clb_internal", "net_type", "INTERNAL"),
					resource.TestCheckResourceAttr("tencentcloud_clb_instance.clb_internal", "project_id", "0"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_instance.clb_internal", "vpc_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_instance.clb_internal", "subnet_id"),
				),
			},
			{
				ResourceName:            "tencentcloud_clb_instance.clb_internal",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"target_region_info"},
			},
		},
	})
}

func testAccCheckClbInstanceDestroy(s *terraform.State) error {
	logId := GetLogId(nil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	clbService := ClbService{
		client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_clb_instance" {
			continue
		}
		time.Sleep(5 * time.Second)

		_, err := clbService.DescribeLoadBalancerById(ctx, rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("clb instance still exists: %s", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckClbInstanceExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := GetLogId(nil)
		ctx := context.WithValue(context.TODO(), "logId", logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("clb instance %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("clb instance id is not set")
		}
		clbService := ClbService{
			client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
		}
		_, err := clbService.DescribeLoadBalancerById(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}
		return nil
	}
}

const testAccClbInstance_basic = `
resource "tencentcloud_clb_instance" "clb_basic" {
	net_type      = "OPEN"
	clb_name      = "tf-clb-basic"

}
`
const testAccClbInstance_internal = `
variable "availability_zone" {
	default = "ap-guangzhou-3"
}

resource "tencentcloud_vpc" "foo" {
    name="guagua-ci-temp-test"
    cidr_block="10.0.0.0/16"
}
resource "tencentcloud_subnet" "subnet" {
	availability_zone="${var.availability_zone}"
	name="guagua-ci-temp-test"
	vpc_id="${tencentcloud_vpc.foo.id}"
	cidr_block="10.0.20.0/28"
	is_multicast=false
}

resource "tencentcloud_clb_instance" "clb_internal" {
	net_type      = "INTERNAL"
	clb_name      = "tf-clb-internal"
  	vpc_id    		  = "${tencentcloud_vpc.foo.id}"
  	subnet_id		  = "${tencentcloud_subnet.subnet.id}"
	project_id = 0
	
}

`
const testAccClbInstance_open = `
resource "tencentcloud_security_group" "foo" {
  name = "ci-temp-test-sg"
}

resource "tencentcloud_vpc" "foo" {
    name="guagua-ci-temp-test"
    cidr_block="10.0.0.0/16"
}

resource "tencentcloud_clb_instance" "clb_open" {
	net_type      = "OPEN"
	clb_name      = "tf-clb-open"
	project_id = 0
	vpc_id="${tencentcloud_vpc.foo.id}"
	target_region_info = {
		region = "ap-guangzhou"
		vpc_id = "${tencentcloud_vpc.foo.id}"
	}

	security_groups = ["${tencentcloud_security_group.foo.id}"]
}
`

const testAccClbInstance_update = `
variable "availability_zone" {
	default = "ap-guangzhou-3"
}

resource "tencentcloud_vpc" "foo" {
    name="guagua-ci-temp-test"
    cidr_block="10.0.0.0/16"
}
resource "tencentcloud_subnet" "subnet" {
	availability_zone="${var.availability_zone}"
	name="guagua-ci-temp-test"
	vpc_id="${tencentcloud_vpc.foo.id}"
	cidr_block="10.0.20.0/28"
	is_multicast=false
}

resource "tencentcloud_clb_instance" "clb_internal" {
	net_type      = "INTERNAL"
	clb_name      = "tf-clb-update"
	vpc_id    		  = "${tencentcloud_vpc.foo.id}"
  	subnet_id		  = "${tencentcloud_subnet.subnet.id}"
	project_id = 0

}
`
const testAccClbInstance_update_open = `
resource "tencentcloud_security_group" "foo" {
  name = "ci-temp-test1-sg"
}

resource "tencentcloud_vpc" "foo" {
    name="guagua-ci-temp-test"
    cidr_block="10.0.0.0/16"
}
resource "tencentcloud_clb_instance" "clb_open" {
	net_type      = "OPEN"
	clb_name      = "tf-clb-update-open"
	vpc_id    		  = "${tencentcloud_vpc.foo.id}"
	project_id = 0
	target_region_info = {
		region = "ap-guangzhou"
		vpc_id = "${tencentcloud_vpc.foo.id}"
	}

	security_groups = ["${tencentcloud_security_group.foo.id}"]

}
`
