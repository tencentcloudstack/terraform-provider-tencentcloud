package tencentcloud

import (
	"context"
	"fmt"
	"strings"
	"testing"

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

// go test -v ./tencentcloud -sweep=ap-guangzhou -sweep-run=tencentcloud_tcr_immutable_tag_rule
func testSweepTcrImmutableTagRule(r string) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	cli, _ := sharedClientForRegion(r)
	tcrService := TCRService{client: cli.(*TencentCloudClient).apiV3Conn}

	// the non-keep namespace will be removed directly when run sweeper tencentcloud_tcr_namespace
	// so... only need to care about the rules under the keep namespace
	rules, err := tcrService.DescribeTcrImmutableTagRuleById(ctx, defaultTCRInstanceId, helper.String(defaultTCRNamespace), nil)

	if err != nil {
		return err
	}

	for _, rule := range rules {
		ruleId := helper.Int64ToStr(*rule.RuleId)

		err = tcrService.DeleteTcrImmutableTagRuleById(ctx, defaultTCRInstanceId, defaultTCRNamespace, ruleId)
		if err != nil {
			continue
		}
	}
	return nil
}

func TestAccTencentCloudTcrImmutableTagRuleResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTcrImmutableTagRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccTcrImmutableTagRule_basic, defaultTCRInstanceId, defaultTCRNamespace),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTcrImmutableTagRuleExists("tencentcloud_tcr_immutable_tag_rule.my_rule"),
					resource.TestCheckResourceAttrSet("tencentcloud_tcr_immutable_tag_rule.my_rule", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_tcr_immutable_tag_rule.my_rule", "rule.#"),
					resource.TestCheckResourceAttr("tencentcloud_tcr_immutable_tag_rule.my_rule", "rule.0.repository_pattern", "**"),
					resource.TestCheckResourceAttr("tencentcloud_tcr_immutable_tag_rule.my_rule", "rule.0.repository_decoration", "repoMatches"),
					resource.TestCheckResourceAttr("tencentcloud_tcr_immutable_tag_rule.my_rule", "rule.0.tag_pattern", "**"),
					resource.TestCheckResourceAttr("tencentcloud_tcr_immutable_tag_rule.my_rule", "rule.0.tag_decoration", "matches"),
					resource.TestCheckResourceAttr("tencentcloud_tcr_immutable_tag_rule.my_rule", "rule.0.disabled", "false"),
					resource.TestCheckResourceAttr("tencentcloud_tcr_immutable_tag_rule.my_rule", "registry_id", defaultTCRInstanceId),
					resource.TestCheckResourceAttr("tencentcloud_tcr_immutable_tag_rule.my_rule", "namespace_name", defaultTCRNamespace),
				),
			},
			{
				Config: fmt.Sprintf(testAccTcrImmutableTagRule_update, defaultTCRInstanceId, defaultTCRNamespace),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTcrImmutableTagRuleExists("tencentcloud_tcr_immutable_tag_rule.my_rule"),
					resource.TestCheckResourceAttrSet("tencentcloud_tcr_immutable_tag_rule.my_rule", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_tcr_immutable_tag_rule.my_rule", "rule.#"),
					resource.TestCheckResourceAttr("tencentcloud_tcr_immutable_tag_rule.my_rule", "rule.0.repository_pattern", "test/*"),
					resource.TestCheckResourceAttr("tencentcloud_tcr_immutable_tag_rule.my_rule", "rule.0.repository_decoration", "repoExcludes"),
					resource.TestCheckResourceAttr("tencentcloud_tcr_immutable_tag_rule.my_rule", "rule.0.tag_pattern", "test/*"),
					resource.TestCheckResourceAttr("tencentcloud_tcr_immutable_tag_rule.my_rule", "rule.0.tag_decoration", "excludes"),
					resource.TestCheckResourceAttr("tencentcloud_tcr_immutable_tag_rule.my_rule", "rule.0.disabled", "false"),
					resource.TestCheckResourceAttr("tencentcloud_tcr_immutable_tag_rule.my_rule", "registry_id", defaultTCRInstanceId),
					resource.TestCheckResourceAttr("tencentcloud_tcr_immutable_tag_rule.my_rule", "namespace_name", defaultTCRNamespace),
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
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TCRService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_tcr_immutable_tag_rule" {
			continue
		}

		idSplit := strings.Split(rs.Primary.ID, FILED_SP)
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
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)
		service := TCRService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}

		rs, ok := s.RootModule().Resources[re]
		if !ok {
			return fmt.Errorf("Tcr Immutable Tag Rule  %s is not found", re)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("Tcr Immutable Tag Rule id is not set")
		}

		idSplit := strings.Split(rs.Primary.ID, FILED_SP)
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

const testAccTcrImmutableTagRule_basic = defaultTCRInstanceData + `

resource "tencentcloud_tcr_immutable_tag_rule" "my_rule" {
  registry_id = "%s"
  namespace_name = "%s"
  rule {
	repository_pattern = "**"
	tag_pattern = "**"
	repository_decoration = "repoMatches"
	tag_decoration = "matches"
	disabled = false
  }
  tags = {
    "createdBy" = "terraform"
  }
}

`

const testAccTcrImmutableTagRule_update = defaultTCRInstanceData + `

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
