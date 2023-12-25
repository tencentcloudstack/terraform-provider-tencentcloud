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
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
)

// go test -i; go test -test.run TestAccTencentCloudTseCngwGatewayResource_basic -v
func TestAccTencentCloudTseCngwGatewayResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckTseCngwGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTseCngwGateway,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTseCngwCngwGatewayExists("tencentcloud_tse_cngw_gateway.cngw_gateway"),
					resource.TestCheckResourceAttrSet("tencentcloud_tse_cngw_gateway.cngw_gateway", "id"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_gateway.cngw_gateway", "description", "terraform test"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_gateway.cngw_gateway", "enable_cls", "false"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_gateway.cngw_gateway", "engine_region", "ap-guangzhou"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_gateway.cngw_gateway", "feature_version", "STANDARD"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_gateway.cngw_gateway", "gateway_version", "2.5.1"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_gateway.cngw_gateway", "ingress_class_name", "tse-nginx-ingress"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_gateway.cngw_gateway", "internet_max_bandwidth_out", "0"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_gateway.cngw_gateway", "name", "terraform-gateway"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_gateway.cngw_gateway", "trade_type", "0"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_gateway.cngw_gateway", "node_config.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_gateway.cngw_gateway", "node_config.0.number", "2"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_gateway.cngw_gateway", "node_config.0.specification", "1c2g"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_gateway.cngw_gateway", "vpc_config.#", "1"),
					resource.TestCheckResourceAttrSet("tencentcloud_tse_cngw_gateway.cngw_gateway", "vpc_config.0.subnet_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_tse_cngw_gateway.cngw_gateway", "vpc_config.0.vpc_id"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_gateway.cngw_gateway", "tags.createdBy", "terraform"),
				),
			},
			{
				ResourceName:      "tencentcloud_tse_cngw_gateway.cngw_gateway",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccTseCngwGatewayUp,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTseCngwCngwGatewayExists("tencentcloud_tse_cngw_gateway.cngw_gateway"),
					resource.TestCheckResourceAttrSet("tencentcloud_tse_cngw_gateway.cngw_gateway", "id"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_gateway.cngw_gateway", "description", "terraform test update"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_gateway.cngw_gateway", "enable_cls", "true"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_gateway.cngw_gateway", "engine_region", "ap-guangzhou"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_gateway.cngw_gateway", "feature_version", "STANDARD"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_gateway.cngw_gateway", "gateway_version", "2.5.1"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_gateway.cngw_gateway", "ingress_class_name", "tse-nginx-ingress"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_gateway.cngw_gateway", "internet_max_bandwidth_out", "0"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_gateway.cngw_gateway", "name", "terraform-gateway-update"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_gateway.cngw_gateway", "trade_type", "0"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_gateway.cngw_gateway", "node_config.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_gateway.cngw_gateway", "node_config.0.number", "2"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_gateway.cngw_gateway", "node_config.0.specification", "1c2g"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_gateway.cngw_gateway", "vpc_config.#", "1"),
					resource.TestCheckResourceAttrSet("tencentcloud_tse_cngw_gateway.cngw_gateway", "vpc_config.0.subnet_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_tse_cngw_gateway.cngw_gateway", "vpc_config.0.vpc_id"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_gateway.cngw_gateway", "tags.createdBy", "terraform"),
				),
			},
		},
	})
}

func testAccCheckTseCngwGatewayDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	service := svctse.NewTseService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_tse_cngw_gateway" {
			continue
		}

		res, err := service.DescribeTseCngwGatewayById(ctx, rs.Primary.ID)
		if err != nil {
			if sdkErr, ok := err.(*sdkErrors.TencentCloudSDKError); ok {
				if sdkErr.Code == "ResourceNotFound.InstanceNotFound" {
					return nil
				}
			}
			return err
		}

		if res != nil {
			return fmt.Errorf("tse gateway %s still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckTseCngwCngwGatewayExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}

		service := svctse.NewTseService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
		res, err := service.DescribeTseCngwGatewayById(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}

		if res == nil {
			return fmt.Errorf("tse gateway %s is not found", rs.Primary.ID)
		}

		return nil
	}
}

const testAccTseCngwGatewayVar = `
variable "availability_zone" {
  default = "ap-guangzhou-4"
}

resource "tencentcloud_vpc" "vpc" {
  cidr_block = "10.0.0.0/16"
  name       = "tf_tse_vpc"
}

resource "tencentcloud_subnet" "subnet" {
  vpc_id            = tencentcloud_vpc.vpc.id
  availability_zone = var.availability_zone
  name              = "tf_tse_subnet"
  cidr_block        = "10.0.1.0/24"
}
`

const testAccTseCngwGateway = testAccTseCngwGatewayVar + `

resource "tencentcloud_tse_cngw_gateway" "cngw_gateway" {
  description                = "terraform test"
  enable_cls                 = false
  engine_region              = "ap-guangzhou"
  feature_version            = "STANDARD"
  gateway_version            = "2.5.1"
  ingress_class_name         = "tse-nginx-ingress"
  internet_max_bandwidth_out = 0
  name                       = "terraform-gateway"
  trade_type                 = 0
  type                       = "kong"

  node_config {
    number        = 2
    specification = "1c2g"
  }

  vpc_config {
    subnet_id = tencentcloud_subnet.subnet.id
    vpc_id    = tencentcloud_vpc.vpc.id
  }

  tags = {
    "createdBy" = "terraform"
  }
}
`

const testAccTseCngwGatewayUp = testAccTseCngwGatewayVar + `

resource "tencentcloud_tse_cngw_gateway" "cngw_gateway" {
  description                = "terraform test update"
  enable_cls                 = true
  engine_region              = "ap-guangzhou"
  feature_version            = "STANDARD"
  gateway_version            = "2.5.1"
  ingress_class_name         = "tse-nginx-ingress"
  internet_max_bandwidth_out = 0
  name                       = "terraform-gateway-update"
  trade_type                 = 0
  type                       = "kong"

  node_config {
    number        = 2
    specification = "1c2g"
  }

  vpc_config {
    subnet_id = tencentcloud_subnet.subnet.id
    vpc_id    = tencentcloud_vpc.vpc.id
  }

  tags = {
    "createdBy" = "terraform"
  }
}
`
