package teo_test

import (
	"context"
	"fmt"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svcteo "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/teo"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccTencentCloudTeoL7AccRuleResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PRIVATE) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckTeoL7AccRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTeoL7AccRule_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTeoL7AccRuleExists("tencentcloud_teo_l7_acc_rule.example"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_l7_acc_rule.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_l7_acc_rule.example", "zone_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_l7_acc_rule.example", "rules.#"),
					// Verify Rules fields are correctly read
					resource.TestCheckResourceAttrSet("tencentcloud_teo_l7_acc_rule.example", "rules.0.rule_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_l7_acc_rule.example", "rules.0.rule_name"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_l7_acc_rule.example", "rules.0.description"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_l7_acc_rule.example", "rules.0.rule_priority"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_l7_acc_rule.example", "rules.0.branches.#"),
				),
			},
			{
				ResourceName:      "tencentcloud_teo_l7_acc_rule.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccTeoL7AccRule_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTeoL7AccRuleExists("tencentcloud_teo_l7_acc_rule.example"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_l7_acc_rule.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_l7_acc_rule.example", "rules.#"),
					// Verify updated rule name is read correctly
					resource.TestCheckResourceAttrSet("tencentcloud_teo_l7_acc_rule.example", "rules.0.rule_name"),
				),
			},
		},
	})
}

func TestAccTencentCloudTeoL7AccRuleResource_emptyRules(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PRIVATE) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckTeoL7AccRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTeoL7AccRule_emptyRules,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTeoL7AccRuleExists("tencentcloud_teo_l7_acc_rule.example"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_l7_acc_rule.example", "id"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule.example", "rules.#", "0"),
				),
			},
		},
	})
}

func testAccCheckTeoL7AccRuleDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	service := svcteo.NewTeoService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_teo_l7_acc_rule" {
			continue
		}

		zoneId := rs.Primary.ID

		resp, err := service.DescribeTeoL7AccRuleById(ctx, zoneId, "")
		if err != nil {
			return err
		}

		if resp != nil && resp.Rules != nil && len(resp.Rules) > 0 {
			return fmt.Errorf("L7 Acc Rule still exists for zone %s", zoneId)
		}
	}
	return nil
}

func testAccCheckTeoL7AccRuleExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}

		zoneId := rs.Primary.ID

		service := svcteo.NewTeoService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
		resp, err := service.DescribeTeoL7AccRuleById(ctx, zoneId, "")
		if err != nil {
			return err
		}

		if resp == nil {
			return fmt.Errorf("L7 Acc Rule does not exist for zone %s", zoneId)
		}

		return nil
	}
}

const testAccTeoL7AccRule_basic = `
resource "tencentcloud_teo_l7_acc_rule" "example" {
  zone_id = "zone-xxxxxxxx"

  rules {
    rule_name = "test-rule-basic"
    description = ["test rule for verification"]

    branches {
      condition = "test-condition"
      actions {
        name = "test-action"
      }
    }
  }
}
`

const testAccTeoL7AccRule_update = `
resource "tencentcloud_teo_l7_acc_rule" "example" {
  zone_id = "zone-xxxxxxxx"

  rules {
    rule_name = "test-rule-updated"
    description = ["updated rule description"]

    branches {
      condition = "test-condition-updated"
      actions {
        name = "test-action-updated"
      }
    }
  }
}
`

const testAccTeoL7AccRule_emptyRules = `
resource "tencentcloud_teo_l7_acc_rule" "example" {
  zone_id = "zone-xxxxxxxx"

  rules = []
}
`
