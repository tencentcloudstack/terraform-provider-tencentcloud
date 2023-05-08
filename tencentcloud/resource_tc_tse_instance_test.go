package tencentcloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// go test -i; go test -test.run TestAccTencentCloudTseInstanceResource_basic -v
func TestAccTencentCloudTseInstanceResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTseInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTseInstance,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTseInstanceExists("tencentcloud_tse_instance.instance"),
					resource.TestCheckResourceAttrSet("tencentcloud_tse_instance.instance", "id"),
					resource.TestCheckResourceAttr("tencentcloud_tse_instance.instance", "engine_type", "zookeeper"),
					resource.TestCheckResourceAttr("tencentcloud_tse_instance.instance", "engine_version", "3.5.9.4"),
					resource.TestCheckResourceAttr("tencentcloud_tse_instance.instance", "engine_product_version", "STANDARD"),
					resource.TestCheckResourceAttr("tencentcloud_tse_instance.instance", "engine_region", "ap-guangzhou"),
					resource.TestCheckResourceAttr("tencentcloud_tse_instance.instance", "engine_name", "zookeeper-test"),
					resource.TestCheckResourceAttr("tencentcloud_tse_instance.instance", "trade_type", "0"),
					resource.TestCheckResourceAttr("tencentcloud_tse_instance.instance", "engine_resource_spec", defaultEngineResourceSpec),
					resource.TestCheckResourceAttr("tencentcloud_tse_instance.instance", "engine_node_num", "3"),
					resource.TestCheckResourceAttr("tencentcloud_tse_instance.instance", "vpc_id", defaultTseVpcId),
					resource.TestCheckResourceAttr("tencentcloud_tse_instance.instance", "subnet_id", defaultTseSubnetId),
				),
			},
			{
				ResourceName:      "tencentcloud_tse_instance.instance",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckTseInstanceDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	service := TseService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_tse_instance" {
			continue
		}

		res, err := service.DescribeTseInstanceById(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}

		if res != nil {
			return fmt.Errorf("tse instance %s still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckTseInstanceExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}

		service := TseService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
		res, err := service.DescribeTseInstanceById(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}

		if res == nil {
			return fmt.Errorf("tse instance %s is not found", rs.Primary.ID)
		}

		return nil
	}
}

const testAccTseInstanceVar = `
variable "engine_resource_spec" {
	default = "` + defaultEngineResourceSpec + `"
}

variable "vpc_id" {
	default = "` + defaultTseVpcId + `"
}

variable "subnet_id" {
	default = "` + defaultTseSubnetId + `"
}
`

const testAccTseInstance = testAccTseInstanceVar + `

resource "tencentcloud_tse_instance" "instance" {
	engine_type = "zookeeper"
	engine_version = "3.5.9.4"
	engine_product_version = "STANDARD"
	engine_region = "ap-guangzhou"
	engine_name = "zookeeper-test"
	trade_type = 0
	engine_resource_spec = var.engine_resource_spec
	engine_node_num = 3
	vpc_id = var.vpc_id
	subnet_id = var.subnet_id
  
	tags = {
	  "createdBy" = "terraform"
	}
}

`
