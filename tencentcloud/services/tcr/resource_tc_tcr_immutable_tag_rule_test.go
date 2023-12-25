package tcr_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svctcr "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tcr"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func init() {
	resource.AddTestSweepers("tencentcloud_tcr_immutable_tag_rule", &resource.Sweeper{
		Name: "tencentcloud_tcr_immutable_tag_rule",
		F:    testSweepTcrImmutableTagRule,
	})
}

// go test -v ./tencentcloud -sweep=ap-shanghai -sweep-run=tencentcloud_tcr_immutable_tag_rule
func testSweepTcrImmutableTagRule(r string) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	cli, _ := tcacctest.SharedClientForRegion(r)
	tcrService := svctcr.NewTCRService(cli.(tccommon.ProviderMeta).GetAPIV3Conn())

	// the non-keep namespace will be removed directly when run sweeper tencentcloud_tcr_namespace
	// so... only need to care about the rules under the keep namespace
	rules, err := tcrService.DescribeTcrImmutableTagRuleById(ctx, tcacctest.DefaultTCRInstanceId, helper.String(tcacctest.DefaultTCRNamespace), nil)
	if err != nil {
		return err
	}

	for _, rule := range rules {
		ruleId := helper.Int64ToStr(*rule.RuleId)

		err = tcrService.DeleteTcrImmutableTagRuleById(ctx, tcacctest.DefaultTCRInstanceId, tcacctest.DefaultTCRNamespace, ruleId)
		if err != nil {
			continue
		}
	}
	return nil
}

func TestAccTencentCloudTcrImmutableTagRuleResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_COMMON) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckTcrImmutableTagRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccTcrImmutableTagRule_basic, tcacctest.DefaultTCRInstanceId, tcacctest.DefaultTCRNamespace),
				PreConfig: func() {
					tcacctest.AccStepSetRegion(t, "ap-shanghai")
					tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_COMMON)
				},
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTcrImmutableTagRuleExists("tencentcloud_tcr_immutable_tag_rule.my_rule"),
					resource.TestCheckResourceAttrSet("tencentcloud_tcr_immutable_tag_rule.my_rule", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_tcr_immutable_tag_rule.my_rule", "rule.#"),
					resource.TestCheckResourceAttr("tencentcloud_tcr_immutable_tag_rule.my_rule", "rule.0.repository_pattern", "**"),
					resource.TestCheckResourceAttr("tencentcloud_tcr_immutable_tag_rule.my_rule", "rule.0.repository_decoration", "repoMatches"),
					resource.TestCheckResourceAttr("tencentcloud_tcr_immutable_tag_rule.my_rule", "rule.0.tag_pattern", "v1"),
					resource.TestCheckResourceAttr("tencentcloud_tcr_immutable_tag_rule.my_rule", "rule.0.tag_decoration", "matches"),
					resource.TestCheckResourceAttr("tencentcloud_tcr_immutable_tag_rule.my_rule", "rule.0.disabled", "false"),
					resource.TestCheckResourceAttr("tencentcloud_tcr_immutable_tag_rule.my_rule", "registry_id", tcacctest.DefaultTCRInstanceId),
					resource.TestCheckResourceAttr("tencentcloud_tcr_immutable_tag_rule.my_rule", "namespace_name", tcacctest.DefaultTCRNamespace),
				),
			},
			{
				Config: fmt.Sprintf(testAccTcrImmutableTagRule_update, tcacctest.DefaultTCRInstanceId, tcacctest.DefaultTCRNamespace),
				PreConfig: func() {
					tcacctest.AccStepSetRegion(t, "ap-shanghai")
					tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_COMMON)
				},
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTcrImmutableTagRuleExists("tencentcloud_tcr_immutable_tag_rule.my_rule"),
					resource.TestCheckResourceAttrSet("tencentcloud_tcr_immutable_tag_rule.my_rule", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_tcr_immutable_tag_rule.my_rule", "rule.#"),
					resource.TestCheckResourceAttr("tencentcloud_tcr_immutable_tag_rule.my_rule", "rule.0.repository_pattern", "test/*"),
					resource.TestCheckResourceAttr("tencentcloud_tcr_immutable_tag_rule.my_rule", "rule.0.repository_decoration", "repoExcludes"),
					resource.TestCheckResourceAttr("tencentcloud_tcr_immutable_tag_rule.my_rule", "rule.0.tag_pattern", "test/*"),
					resource.TestCheckResourceAttr("tencentcloud_tcr_immutable_tag_rule.my_rule", "rule.0.tag_decoration", "excludes"),
					resource.TestCheckResourceAttr("tencentcloud_tcr_immutable_tag_rule.my_rule", "rule.0.disabled", "false"),
					resource.TestCheckResourceAttr("tencentcloud_tcr_immutable_tag_rule.my_rule", "registry_id", tcacctest.DefaultTCRInstanceId),
					resource.TestCheckResourceAttr("tencentcloud_tcr_immutable_tag_rule.my_rule", "namespace_name", tcacctest.DefaultTCRNamespace),
				),
			},
			{
				ResourceName:      "tencentcloud_tcr_immutable_tag_rule.my_rule",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckTcrImmutableTagRuleDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := svctcr.NewTCRService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_tcr_immutable_tag_rule" {
			continue
		}

		idSplit := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
		if len(idSplit) != 3 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		registryId := idSplit[0]
		namespaceName := idSplit[1]
		ruleId := idSplit[2]

		rule, err := service.DescribeTcrImmutableTagRuleById(ctx, registryId, &namespaceName, &ruleId)
		if err != nil {
			return err
		}

		if rule != nil || len(rule) > 0 {
			return fmt.Errorf("Tcr Immutable Tag Rule still exist, Id: %v", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckTcrImmutableTagRuleExists(re string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service := svctcr.NewTCRService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())

		rs, ok := s.RootModule().Resources[re]
		if !ok {
			return fmt.Errorf("Tcr Immutable Tag Rule  %s is not found", re)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("Tcr Immutable Tag Rule id is not set")
		}

		idSplit := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
		if len(idSplit) != 3 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		registryId := idSplit[0]
		namespaceName := idSplit[1]
		ruleId := idSplit[2]

		rule, err := service.DescribeTcrImmutableTagRuleById(ctx, registryId, &namespaceName, &ruleId)
		if err != nil {
			return err
		}

		if len(rule) == 0 {
			return fmt.Errorf("Tcr Immutable Tag Rule not found, Id: %v", rs.Primary.ID)
		}
		return nil
	}
}

const testAccTcrImmutableTagRule_basic = tcacctest.DefaultTCRInstanceData + `

resource "tencentcloud_tcr_immutable_tag_rule" "my_rule" {
  registry_id = "%s"
  namespace_name = "%s"
  rule {
	repository_pattern = "**"
	tag_pattern = "v1"
	repository_decoration = "repoMatches"
	tag_decoration = "matches"
	disabled = false
  }
  tags = {
    "createdBy" = "terraform"
  }
}

`

const testAccTcrImmutableTagRule_update = tcacctest.DefaultTCRInstanceData + `

resource "tencentcloud_tcr_immutable_tag_rule" "my_rule" {
  registry_id = "%s"
  namespace_name = "%s"
  rule {
		repository_pattern = "test/*"
		tag_pattern = "test/*"
		repository_decoration = "repoExcludes"
		tag_decoration = "excludes"
		disabled = false
  }
  tags = {
    "createdBy" = "terraform"
  }
}

`
