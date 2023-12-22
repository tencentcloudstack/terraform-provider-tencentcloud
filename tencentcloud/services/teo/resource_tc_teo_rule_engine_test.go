package teo_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svcteo "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/teo"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// go test -i; go test -test.run TestAccTencentCloudTeoRuleEngine_basic -v
func TestAccTencentCloudTeoRuleEngine_basic(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PRIVATE) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckRuleEngineDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTeoRuleEngine,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRuleEngineExists("tencentcloud_teo_rule_engine.basic"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_rule_engine.basic", "zone_id"),
					resource.TestCheckResourceAttr("tencentcloud_teo_rule_engine.basic", "rule_name", "rule-1"),
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
			{
				Config: testAccTeoRuleEngineUp,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRuleEngineExists("tencentcloud_teo_rule_engine.basic"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_rule_engine.basic", "zone_id"),
					resource.TestCheckResourceAttr("tencentcloud_teo_rule_engine.basic", "rule_name", "rule-up"),
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
		},
	})
}

func testAccCheckRuleEngineDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	service := svcteo.NewTeoService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_teo_rule_engine" {
			continue
		}
		idSplit := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
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
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}
		idSplit := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
		if len(idSplit) != 2 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		zoneId := idSplit[0]
		ruleId := idSplit[1]

		service := svcteo.NewTeoService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
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

const testAccTeoRuleEngine = testAccTeoZone + `

resource "tencentcloud_teo_rule_engine" "basic" {
	rule_name = "rule-1"
	status    = "enable"
	zone_id   = tencentcloud_teo_zone.basic.id
  
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

const testAccTeoRuleEngineUp = testAccTeoZone + `

resource "tencentcloud_teo_rule_engine" "basic" {
	rule_name = "rule-up"
	status    = "enable"
	zone_id   = tencentcloud_teo_zone.basic.id
  
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
