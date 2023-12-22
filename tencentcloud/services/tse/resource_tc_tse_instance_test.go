package tse_test

import (
	"context"
	"fmt"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svctse "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tse"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// go test -i; go test -test.run TestAccTencentCloudTseInstanceResource_basic -v
func TestAccTencentCloudTseInstanceResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers:    tcacctest.AccProviders,
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
					resource.TestCheckResourceAttr("tencentcloud_tse_instance.instance", "engine_resource_spec", tcacctest.DefaultEngineResourceSpec),
					resource.TestCheckResourceAttr("tencentcloud_tse_instance.instance", "engine_node_num", "3"),
					resource.TestCheckResourceAttr("tencentcloud_tse_instance.instance", "vpc_id", tcacctest.DefaultTseVpcId),
					resource.TestCheckResourceAttr("tencentcloud_tse_instance.instance", "subnet_id", tcacctest.DefaultTseSubnetId),
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
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	service := svctse.NewTseService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
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
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}

		service := svctse.NewTseService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
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
	default = "` + tcacctest.DefaultEngineResourceSpec + `"
}

variable "vpc_id" {
	default = "` + tcacctest.DefaultTseVpcId + `"
}

variable "subnet_id" {
	default = "` + tcacctest.DefaultTseSubnetId + `"
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
