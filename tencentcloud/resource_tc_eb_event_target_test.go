package tencentcloud

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
)

// go test -i; go test -test.run TestAccTencentCloudEbEventTargetResource_basic -v
func TestAccTencentCloudEbEventTargetResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEbEventTargetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEbEventTarget,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEbEventTargetExists("tencentcloud_eb_event_target.event_target"),
					resource.TestCheckResourceAttrSet("tencentcloud_eb_event_target.event_target", "id"),
					resource.TestCheckResourceAttr("tencentcloud_eb_event_target.event_target", "type", "scf"),
					resource.TestCheckResourceAttr("tencentcloud_eb_event_target.event_target", "target_description.#", "1"),
					resource.TestCheckResourceAttrSet("tencentcloud_eb_event_target.event_target", "target_description.0.resource_description"),
					resource.TestCheckResourceAttr("tencentcloud_eb_event_target.event_target", "target_description.0.scf_params.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_eb_event_target.event_target", "target_description.0.scf_params.0.batch_event_count", "1"),
					resource.TestCheckResourceAttr("tencentcloud_eb_event_target.event_target", "target_description.0.scf_params.0.batch_timeout", "1"),
					resource.TestCheckResourceAttr("tencentcloud_eb_event_target.event_target", "target_description.0.scf_params.0.enable_batch_delivery", "true"),
				),
			},
			{
				ResourceName:      "tencentcloud_eb_event_target.event_target",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckEbEventTargetDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	service := EbService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_eb_event_target" {
			continue
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("resource id is not set")
		}
		idSplit := strings.Split(rs.Primary.ID, FILED_SP)
		if len(idSplit) != 3 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		eventBusId := idSplit[0]
		ruleId := idSplit[1]
		targetId := idSplit[2]

		target, err := service.DescribeEbEventTargetById(ctx, eventBusId, ruleId, targetId)
		if err != nil {
			if err.(*sdkErrors.TencentCloudSDKError).Code == "ResourceNotFound.Rule" {
				return nil
			}
			return err
		}

		if target != nil {
			return fmt.Errorf("eb target %s still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckEbEventTargetExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("resource id is not set")
		}
		idSplit := strings.Split(rs.Primary.ID, FILED_SP)
		if len(idSplit) != 3 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		eventBusId := idSplit[0]
		ruleId := idSplit[1]
		targetId := idSplit[2]

		service := EbService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
		target, err := service.DescribeEbEventTargetById(ctx, eventBusId, ruleId, targetId)
		if err != nil {
			return err
		}

		if target == nil {
			return fmt.Errorf("eb target %s is not found", rs.Primary.ID)
		}

		return nil
	}
}

const testAccEbEventTargetVar = `
variable "zone" {
  default = "ap-guangzhou"
}

variable "namespace" {
  default = "default"
}

variable "function" {
  default = "keep-1676351130"
}

variable "function_version" {
  default = "$LATEST"
}

data "tencentcloud_cam_users" "foo" {
}

resource "tencentcloud_eb_event_bus" "foo" {
  event_bus_name = "tf-event_bus"
  description    = "event bus desc"
  enable_store   = false
  save_days      = 1
  tags = {
    "createdBy" = "terraform"
  }
}

resource "tencentcloud_eb_event_rule" "foo" {
  event_bus_id = tencentcloud_eb_event_bus.foo.id
  rule_name    = "tf-event_rule"
  description  = "event rule desc"
  enable       = true
  event_pattern = jsonencode(
    {
      source = "apigw.cloud.tencent"
      type = [
        "connector:apigw",
      ]
    }
  )
  tags = {
    "createdBy" = "terraform"
  }
}
`

const testAccEbEventTarget = testAccEbEventTargetVar + `

resource "tencentcloud_eb_event_target" "event_target" {
    event_bus_id = tencentcloud_eb_event_bus.foo.id
    rule_id      = tencentcloud_eb_event_rule.foo.rule_id
    type         = "scf"

    target_description {
        resource_description = "qcs::scf:${var.zone}:uin/${data.tencentcloud_cam_users.foo.user_list.0.uin}:namespace/${var.namespace}/function/${var.function}/${var.function_version}"

        scf_params {
            batch_event_count     = 1
            batch_timeout         = 1
            enable_batch_delivery = true
        }
    }
}

`
