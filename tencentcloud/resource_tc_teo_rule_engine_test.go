package tencentcloud

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudTeoRuleEngine_basic -v
func TestAccTencentCloudTeoRuleEngine_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PRIVATE) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRuleEngineDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTeoRuleEngine,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRuleEngineExists("tencentcloud_teo_rule_engine.basic"),
					resource.TestCheckResourceAttr("tencentcloud_teo_rule_engine.basic", "zone_id", defaultZoneId),
					resource.TestCheckResourceAttr("tencentcloud_teo_rule_engine.basic", "status", "enable"),
					resource.TestCheckResourceAttr("tencentcloud_teo_rule_engine.basic", "rules.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_rule_engine.basic", "rules.0.actions.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_rule_engine.basic", "rules.0.or.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_rule_engine.basic", "rules.0.sub_rules.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_rule_engine.basic", "rules.0.sub_rules.0.tags.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_rule_engine.basic", "rules.0.sub_rules.0.rules.0.or.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_rule_engine.basic", "rules.0.sub_rules.0.rules.0.or.0.and.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_rule_engine.basic", "rules.0.sub_rules.0.rules.0.or.0.and.0.operator", "equal"),
					resource.TestCheckResourceAttr("tencentcloud_teo_rule_engine.basic", "rules.0.sub_rules.0.rules.0.or.0.and.0.target", "filename"),
					resource.TestCheckResourceAttr("tencentcloud_teo_rule_engine.basic", "rules.0.sub_rules.0.rules.0.or.0.and.0.ignore_case", "false"),
					resource.TestCheckResourceAttr("tencentcloud_teo_rule_engine.basic", "rules.0.sub_rules.0.rules.0.or.0.and.0.values.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_rule_engine.basic", "rules.0.sub_rules.0.rules.0.actions.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_rule_engine.basic", "rules.0.sub_rules.0.rules.0.actions.0.normal_action.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_rule_engine.basic", "rules.0.sub_rules.0.rules.0.actions.0.normal_action.0.action", "HostHeader"),
					resource.TestCheckResourceAttr("tencentcloud_teo_rule_engine.basic", "rules.0.sub_rules.0.rules.0.actions.0.normal_action.0.parameters.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_rule_engine.basic", "rules.0.sub_rules.0.rules.0.actions.0.normal_action.0.parameters.0.name", "ServerName"),
					resource.TestCheckResourceAttr("tencentcloud_teo_rule_engine.basic", "rules.0.sub_rules.0.rules.0.actions.0.normal_action.0.parameters.0.values.#", "1"),
				),
			},
			{
				ResourceName:      "tencentcloud_teo_rule_engine.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckRuleEngineDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	service := TeoService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_teo_rule_engine" {
			continue
		}
		idSplit := strings.Split(rs.Primary.ID, FILED_SP)
		if len(idSplit) != 2 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		zoneId := idSplit[0]
		ruleId := idSplit[1]

		originGroup, err := service.DescribeTeoRuleEngine(ctx, zoneId, ruleId)
		if originGroup != nil {
			return fmt.Errorf("zone ruleEngine %s still exists", rs.Primary.ID)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func testAccCheckRuleEngineExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}
		idSplit := strings.Split(rs.Primary.ID, FILED_SP)
		if len(idSplit) != 2 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		zoneId := idSplit[0]
		ruleId := idSplit[1]

		service := TeoService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
		originGroup, err := service.DescribeTeoRuleEngine(ctx, zoneId, ruleId)
		if originGroup == nil {
			return fmt.Errorf("zone ruleEngine %s is not found", rs.Primary.ID)
		}
		if err != nil {
			return err
		}

		return nil
	}
}

const testAccTeoRuleEngineVar = `
variable "zone_id" {
  default = "` + defaultZoneId + `"
}
`

const testAccTeoRuleEngine = testAccTeoRuleEngineVar + `

resource "tencentcloud_teo_rule_engine" "basic" {
	rule_name = "rule-1"
	status    = "enable"
	zone_id   = var.zone_id
  
	rules {
	  actions {
  
		rewrite_action {
		  action = "ResponseHeader"
  
		  parameters {
			action = "set"
			name   = "project"
			values = [
			  "1111",
			]
		  }
		}
	  }
  
	  or {
		and {
		  operator = "equal"
		  target   = "extension"
		  values   = [
			"11",
		  ]
		}
	  }
	  sub_rules {
		  tags = ["test-tag",]
		  rules {
			or {
			  and {
				operator = "equal"
				target = "filename"
				ignore_case = false
				values = ["test.txt"]
			  }
			}
			actions {
				normal_action {
					action = "HostHeader"
					parameters {
						name = "ServerName"
						values = ["terraform-test.com"]
					}
				}
			}
		  }
	  }
	}
  }
`
