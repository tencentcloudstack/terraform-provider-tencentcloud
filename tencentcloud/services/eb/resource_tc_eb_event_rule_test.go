package eb_test

import (
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svceb "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/eb"

	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
)

// go test -i; go test -test.run TestAccTencentCloudEbEventRuleResource_basic -v
func TestAccTencentCloudEbEventRuleResource_basic(t *testing.T) {
	// t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckEbEventRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEbEventRule,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEbEventRuleExists("tencentcloud_eb_event_rule.event_rule"),
					resource.TestCheckResourceAttrSet("tencentcloud_eb_event_rule.event_rule", "id"),
					resource.TestCheckResourceAttr("tencentcloud_eb_event_rule.event_rule", "rule_name", "tf-event_rule"),
					resource.TestCheckResourceAttr("tencentcloud_eb_event_rule.event_rule", "description", "event rule desc"),
					resource.TestCheckResourceAttr("tencentcloud_eb_event_rule.event_rule", "enable", "true"),
					resource.TestCheckResourceAttrSet("tencentcloud_eb_event_rule.event_rule", "event_pattern"),
					resource.TestCheckResourceAttr("tencentcloud_eb_event_rule.event_rule", "tags.createdBy", "terraform"),
				),
			},
			{
				ResourceName:      "tencentcloud_eb_event_rule.event_rule",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccEbEventRuleUp,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEbEventRuleExists("tencentcloud_eb_event_rule.event_rule"),
					resource.TestCheckResourceAttrSet("tencentcloud_eb_event_rule.event_rule", "id"),
					resource.TestCheckResourceAttr("tencentcloud_eb_event_rule.event_rule", "rule_name", "tf-event_rule"),
					resource.TestCheckResourceAttr("tencentcloud_eb_event_rule.event_rule", "description", "event rule"),
					resource.TestCheckResourceAttr("tencentcloud_eb_event_rule.event_rule", "enable", "false"),
					resource.TestCheckResourceAttrSet("tencentcloud_eb_event_rule.event_rule", "event_pattern"),
					resource.TestCheckResourceAttr("tencentcloud_eb_event_rule.event_rule", "tags.createdBy", "terraform-test"),
				),
			},
		},
	})
}

func testAccCheckEbEventRuleDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	service := svceb.NewEbService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_eb_event_rule" {
			continue
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("resource id is not set")
		}
		idSplit := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
		if len(idSplit) != 2 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		eventBusId := idSplit[0]
		ruleId := idSplit[1]

		rule, err := service.DescribeEbEventRuleById(ctx, eventBusId, ruleId)
		if err != nil {
			if err.(*sdkErrors.TencentCloudSDKError).Code == "ResourceNotFound" {
				return nil
			}
			return err
		}

		if rule != nil {
			return fmt.Errorf("eb eventRule %s still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckEbEventRuleExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("resource id is not set")
		}
		idSplit := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
		if len(idSplit) != 2 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		eventBusId := idSplit[0]
		ruleId := idSplit[1]

		service := svceb.NewEbService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
		rule, err := service.DescribeEbEventRuleById(ctx, eventBusId, ruleId)
		if err != nil {
			return err
		}

		if rule == nil {
			return fmt.Errorf("eb eventRule %s is not found", rs.Primary.ID)
		}

		return nil
	}
}

const testAccEbEventRuleVar = `
resource "tencentcloud_eb_event_bus" "foo" {
	event_bus_name = "tf-event_bus_rule"
	description    = "event bus desc"
	enable_store   = false
	save_days      = 1
	tags = {
	  "createdBy" = "terraform"
	}
  }
`

const testAccEbEventRule = testAccEbEventRuleVar + `

resource "tencentcloud_eb_event_rule" "event_rule" {
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

const testAccEbEventRuleUp = testAccEbEventRuleVar + `

resource "tencentcloud_eb_event_rule" "event_rule" {
  event_bus_id = tencentcloud_eb_event_bus.foo.id
  rule_name    = "tf-event_rule"
  description  = "event rule"
  enable       = false
  event_pattern = jsonencode(
    {
      source = "apigw.cloud.tencent"
      type = [
        "apigw:CloudEvent:ApiCall"
      ]
    }
  )
  tags = {
    "createdBy" = "terraform-test"
  }
}

`
