package tse_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svctse "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tse"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
)

// go test -i; go test -test.run TestAccTencentCloudTseCngwServiceResource_basic -v
func TestAccTencentCloudTseCngwServiceResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckTseCngwServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTseCngwService,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTseCngwServiceExists("tencentcloud_tse_cngw_service.cngw_service"),
					resource.TestCheckResourceAttrSet("tencentcloud_tse_cngw_service.cngw_service", "id"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_service.cngw_service", "gateway_id", tcacctest.DefaultTseGatewayId),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_service.cngw_service", "name", "terraform-test"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_service.cngw_service", "path", "/test"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_service.cngw_service", "protocol", "http"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_service.cngw_service", "retries", "5"),
					resource.TestCheckResourceAttrSet("tencentcloud_tse_cngw_service.cngw_service", "service_id"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_service.cngw_service", "timeout", "60000"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_service.cngw_service", "upstream_type", "IPList"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_service.cngw_service", "upstream_info.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_service.cngw_service", "upstream_info.0.algorithm", "round-robin"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_service.cngw_service", "upstream_info.0.auto_scaling_cvm_port", "80"),
					resource.TestCheckResourceAttrSet("tencentcloud_tse_cngw_service.cngw_service", "upstream_info.0.auto_scaling_group_id"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_service.cngw_service", "upstream_info.0.auto_scaling_hook_status", "Normal"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_service.cngw_service", "upstream_info.0.auto_scaling_tat_cmd_status", "Normal"),
					resource.TestCheckResourceAttrSet("tencentcloud_tse_cngw_service.cngw_service", "upstream_info.0.host"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_service.cngw_service", "upstream_info.0.slow_start", "20"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_service.cngw_service", "upstream_info.0.targets.#", "1"),
					resource.TestCheckResourceAttrSet("tencentcloud_tse_cngw_service.cngw_service", "upstream_info.0.targets.0.created_time"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_service.cngw_service", "upstream_info.0.targets.0.host", "192.168.0.1"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_service.cngw_service", "upstream_info.0.targets.0.port", "80"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_service.cngw_service", "upstream_info.0.targets.0.weight", "100"),
				),
			},
			{
				ResourceName:            "tencentcloud_tse_cngw_service.cngw_service",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"upstream_info.0.targets.0.health"},
			},
			{
				Config: testAccTseCngwServiceUp,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTseCngwServiceExists("tencentcloud_tse_cngw_service.cngw_service"),
					resource.TestCheckResourceAttrSet("tencentcloud_tse_cngw_service.cngw_service", "id"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_service.cngw_service", "gateway_id", tcacctest.DefaultTseGatewayId),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_service.cngw_service", "name", "terraform-test"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_service.cngw_service", "path", "/test-1"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_service.cngw_service", "protocol", "http"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_service.cngw_service", "retries", "5"),
					resource.TestCheckResourceAttrSet("tencentcloud_tse_cngw_service.cngw_service", "service_id"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_service.cngw_service", "timeout", "6000"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_service.cngw_service", "upstream_type", "IPList"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_service.cngw_service", "upstream_info.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_service.cngw_service", "upstream_info.0.algorithm", "round-robin"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_service.cngw_service", "upstream_info.0.auto_scaling_cvm_port", "80"),
					resource.TestCheckResourceAttrSet("tencentcloud_tse_cngw_service.cngw_service", "upstream_info.0.auto_scaling_group_id"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_service.cngw_service", "upstream_info.0.auto_scaling_hook_status", "Normal"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_service.cngw_service", "upstream_info.0.auto_scaling_tat_cmd_status", "Normal"),
					resource.TestCheckResourceAttrSet("tencentcloud_tse_cngw_service.cngw_service", "upstream_info.0.host"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_service.cngw_service", "upstream_info.0.slow_start", "20"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_service.cngw_service", "upstream_info.0.targets.#", "1"),
					resource.TestCheckResourceAttrSet("tencentcloud_tse_cngw_service.cngw_service", "upstream_info.0.targets.0.created_time"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_service.cngw_service", "upstream_info.0.targets.0.host", "192.168.0.1"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_service.cngw_service", "upstream_info.0.targets.0.port", "80"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_service.cngw_service", "upstream_info.0.targets.0.weight", "80"),
				),
			},
		},
	})
}

func testAccCheckTseCngwServiceDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	service := svctse.NewTseService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_tse_cngw_service" {
			continue
		}

		idSplit := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
		if len(idSplit) != 2 {
			return fmt.Errorf("invalid ID %s", rs.Primary.ID)
		}
		gatewayId := idSplit[0]
		name := idSplit[1]

		res, err := service.DescribeTseCngwServiceById(ctx, gatewayId, name)
		if err != nil {
			if sdkErr, ok := err.(*sdkErrors.TencentCloudSDKError); ok {
				if sdkErr.Code == "ResourceNotFound.ResourceNotFound" {
					return nil
				}
			}
			return err
		}

		if res != nil {
			return fmt.Errorf("tse cngwService %s still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckTseCngwServiceExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}

		idSplit := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
		if len(idSplit) != 2 {
			return fmt.Errorf("invalid ID %s", rs.Primary.ID)
		}
		gatewayId := idSplit[0]
		name := idSplit[1]

		service := svctse.NewTseService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
		res, err := service.DescribeTseCngwServiceById(ctx, gatewayId, name)
		if err != nil {
			return err
		}

		if res == nil {
			return fmt.Errorf("tse cngwService %s is not found", rs.Primary.ID)
		}

		return nil
	}
}

const testAccTseCngwService = tcacctest.DefaultTseVar + `

resource "tencentcloud_tse_cngw_service" "cngw_service" {
	gateway_id = var.gateway_id
	name       = "terraform-test"
	path       = "/test"
	protocol   = "http"
	retries    = 5
	timeout       = 60000
	upstream_type = "IPList"
  
	upstream_info {
	  algorithm                   = "round-robin"
	  auto_scaling_cvm_port       = 80
	  auto_scaling_group_id       = "asg-519acdug"
	  auto_scaling_hook_status    = "Normal"
	  auto_scaling_tat_cmd_status = "Normal"
	  port                        = 0
	  slow_start                  = 20
  
	  targets {
		# health = "HEALTHCHECKS_OFF"
		host   = "192.168.0.1"
		port   = 80
		weight = 100
	  }
	}
}

`

const testAccTseCngwServiceUp = tcacctest.DefaultTseVar + `

resource "tencentcloud_tse_cngw_service" "cngw_service" {
	gateway_id = var.gateway_id
	name       = "terraform-test"
	path       = "/test-1"
	protocol   = "http"
	retries    = 5
	timeout       = 6000
	upstream_type = "IPList"
  
	upstream_info {
	  algorithm                   = "round-robin"
	  auto_scaling_cvm_port       = 80
	  auto_scaling_group_id       = "asg-519acdug"
	  auto_scaling_hook_status    = "Normal"
	  auto_scaling_tat_cmd_status = "Normal"
	  port                        = 0
	  slow_start                  = 20
  
	  targets {
		host   = "192.168.0.1"
		port   = 80
		weight = 80
	  }
	}
}

`
