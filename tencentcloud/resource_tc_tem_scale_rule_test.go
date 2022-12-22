package tencentcloud

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// go test -i; go test -test.run TestAccTencentCloudNeedFixTemScaleRuleResource_basic -v
func TestAccTencentCloudNeedFixTemScaleRuleResource_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTemScaleRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTemScaleRule,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTemScaleRuleExists("tencentcloud_tem_scale_rule.scaleRule"),
					resource.TestCheckResourceAttrSet("tencentcloud_tem_scale_rule.scaleRule", "id"),
					resource.TestCheckResourceAttr("tencentcloud_tem_scale_rule.scaleRule", "environment_id", defaultEnvironmentId),
					resource.TestCheckResourceAttr("tencentcloud_tem_scale_rule.scaleRule", "application_id", defaultApplicationId),
					resource.TestCheckResourceAttr("tencentcloud_tem_scale_rule.scaleRule", "autoscaler.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_tem_scale_rule.scaleRule", "autoscaler.0.autoscaler_name", "test3123"),
				),
			},
			{
				ResourceName:      "tencentcloud_tem_scale_rule.scaleRule",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckTemScaleRuleDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	service := TemService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_tem_scale_rule" {
			continue
		}

		idSplit := strings.Split(rs.Primary.ID, FILED_SP)
		if len(idSplit) != 3 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		environmentId := idSplit[0]
		applicationId := idSplit[1]
		scaleRuleId := idSplit[2]

		res, err := service.DescribeTemScaleRule(ctx, environmentId, applicationId, scaleRuleId)
		if err != nil {
			return err
		}

		if res != nil {
			return fmt.Errorf("tem scale rule %s still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckTemScaleRuleExists(r string) resource.TestCheckFunc {
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
		environmentId := idSplit[0]
		applicationId := idSplit[1]
		scaleRuleId := idSplit[2]

		service := TemService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
		res, err := service.DescribeTemScaleRule(ctx, environmentId, applicationId, scaleRuleId)
		if err != nil {
			return err
		}

		if res == nil {
			return fmt.Errorf("tem scale rule %s is not found", rs.Primary.ID)
		}

		return nil
	}
}

const testAccTemScaleRuleVar = `
variable "environment_id" {
	default = "` + defaultEnvironmentId + `"
}

variable "application_id" {
	default = "` + defaultApplicationId + `"
}

variable "workload_id" {
	default = "` + defaultEnvironmentId + "#" + defaultApplicationId + `"
}
`

const testAccTemScaleRule = testAccTemScaleRuleVar + `

resource "tencentcloud_tem_scale_rule" "scaleRule" {
  environment_id = var.environment_id
  application_id = var.application_id
  workload_id = var.workload_id
  autoscaler {
    autoscaler_name = "test3123"
    description     = "test"
    enabled         = true
    min_replicas    = 1
    max_replicas    = 4
    cron_horizontal_autoscaler {
      name     = "test"
      period   = "* * *"
      priority = 1
      enabled  = true
      schedules {
        start_at        = "03:00"
        target_replicas = 1
      }
    }
    cron_horizontal_autoscaler {
      name     = "test123123"
      period   = "* * *"
      priority = 0
      enabled  = true
      schedules {
        start_at        = "04:13"
        target_replicas = 1
      }
    }
    horizontal_autoscaler {
      metrics      = "CPU"
      enabled      = true
      max_replicas = 4
      min_replicas = 1
      threshold    = 60
    }

  }
}
`
