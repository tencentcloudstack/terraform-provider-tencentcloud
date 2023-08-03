package tencentcloud

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccTencentCloudScfProvisionedConcurrencyConfigResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		CheckDestroy: testAccCheckScfProvisionedConcurrencyConfigDestroy,
		Providers:    testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccScfProvisionedConcurrencyConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckScfProvisionedConcurrencyConfigExists("tencentcloud_scf_provisioned_concurrency_config.provisioned_concurrency_config"),
					resource.TestCheckResourceAttrSet("tencentcloud_scf_provisioned_concurrency_config.provisioned_concurrency_config", "function_name"),
					resource.TestCheckResourceAttrSet("tencentcloud_scf_provisioned_concurrency_config.provisioned_concurrency_config", "qualifier"),
					resource.TestCheckResourceAttrSet("tencentcloud_scf_provisioned_concurrency_config.provisioned_concurrency_config", "version_provisioned_concurrency_num"),
					resource.TestCheckResourceAttrSet("tencentcloud_scf_provisioned_concurrency_config.provisioned_concurrency_config", "trigger_actions.#"),
					resource.TestCheckResourceAttrSet("tencentcloud_scf_provisioned_concurrency_config.provisioned_concurrency_config", "trigger_actions.0.trigger_name"),
					resource.TestCheckResourceAttrSet("tencentcloud_scf_provisioned_concurrency_config.provisioned_concurrency_config", "trigger_actions.0.trigger_provisioned_concurrency_num"),
					resource.TestCheckResourceAttrSet("tencentcloud_scf_provisioned_concurrency_config.provisioned_concurrency_config", "trigger_actions.0.trigger_cron_config"),
					resource.TestCheckResourceAttrSet("tencentcloud_scf_provisioned_concurrency_config.provisioned_concurrency_config", "trigger_actions.0.provisioned_type"),
					resource.TestCheckResourceAttrSet("tencentcloud_scf_provisioned_concurrency_config.provisioned_concurrency_config", "provisioned_type"),
					resource.TestCheckResourceAttrSet("tencentcloud_scf_provisioned_concurrency_config.provisioned_concurrency_config", "tracking_target"),
					resource.TestCheckResourceAttrSet("tencentcloud_scf_provisioned_concurrency_config.provisioned_concurrency_config", "min_capacity"),
					resource.TestCheckResourceAttrSet("tencentcloud_scf_provisioned_concurrency_config.provisioned_concurrency_config", "max_capacity"),
				),
			},
		},
	})
}
func testAccCheckScfProvisionedConcurrencyConfigDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	service := ScfService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_scf_provisioned_concurrency_config" {
			continue
		}

		idSplit := strings.Split(rs.Primary.ID, FILED_SP)
		if len(idSplit) != 3 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		functionName := idSplit[0]
		qualifier := idSplit[1]
		namespace := idSplit[2]

		provisionedConcurrencyConfig, err := service.DescribeScfProvisionedConcurrencyConfigById(ctx, functionName, qualifier, namespace)
		if err != nil {
			return err
		}

		if provisionedConcurrencyConfig != nil {
			return fmt.Errorf("ScfProvisionedConcurrencyConfig Resource %s still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckScfProvisionedConcurrencyConfigExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}

		idSplit := strings.Split(rs.Primary.ID, FILED_SP)
		if len(idSplit) != 3 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		functionName := idSplit[0]
		qualifier := idSplit[1]
		namespace := idSplit[2]

		service := ScfService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
		provisionedConcurrencyConfig, err := service.DescribeScfProvisionedConcurrencyConfigById(ctx, functionName, qualifier, namespace)
		if err != nil {
			return err
		}
		if provisionedConcurrencyConfig == nil {
			return fmt.Errorf("ScfProvisionedConcurrencyConfig Resource %s is not found", rs.Primary.ID)
		}
		return nil
	}
}

const testAccScfProvisionedConcurrencyConfig = `

resource "tencentcloud_scf_provisioned_concurrency_config" "provisioned_concurrency_config" {
  function_name                       = "keep-1676351130"
  qualifier                           = "2"
  version_provisioned_concurrency_num = 2
  namespace                           = "default"
  trigger_actions {
    trigger_name                        = "test"
    trigger_provisioned_concurrency_num = 2
    trigger_cron_config                 = "29 45 12 29 05 * 2023"
    provisioned_type                    = "Default"
  }
  provisioned_type                    = "Default"
  tracking_target                     = 0.5
  min_capacity                        = 1
  max_capacity                        = 2
}

`
