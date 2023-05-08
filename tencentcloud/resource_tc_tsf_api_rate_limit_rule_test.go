package tencentcloud

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// go test -i; go test -test.run TestAccTencentCloudTsfApiRateLimitRuleResource_basic -v
func TestAccTencentCloudTsfApiRateLimitRuleResource_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_TSF) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTsfApiRateLimitRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTsfApiRateLimitRule,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTsfApiRateLimitRuleExists("tencentcloud_tsf_api_rate_limit_rule.api_rate_limit_rule"),
					resource.TestCheckResourceAttrSet("tencentcloud_tsf_api_rate_limit_rule.api_rate_limit_rule", "id"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_api_rate_limit_rule.api_rate_limit_rule", "api_id", defaultTsfApiId),
					resource.TestCheckResourceAttr("tencentcloud_tsf_api_rate_limit_rule.api_rate_limit_rule", "max_qps", "10"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_api_rate_limit_rule.api_rate_limit_rule", "usable_status", "enabled"),
				),
			},
			{
				ResourceName:      "tencentcloud_tsf_api_rate_limit_rule.api_rate_limit_rule",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckTsfApiRateLimitRuleDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	service := TsfService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_tsf_api_rate_limit_rule" {
			continue
		}

		idSplit := strings.Split(rs.Primary.ID, FILED_SP)
		if len(idSplit) != 2 {
			return fmt.Errorf("invalid ID %s", rs.Primary.ID)
		}
		apiId := idSplit[0]
		ruleId := idSplit[1]

		res, err := service.DescribeTsfApiRateLimitRuleById(ctx, apiId, ruleId)
		if err != nil {
			return err
		}

		if res != nil {
			return fmt.Errorf("tsf ApiRateLimitRule %s still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckTsfApiRateLimitRuleExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}
		idSplit := strings.Split(rs.Primary.ID, FILED_SP)
		if len(idSplit) != 2 {
			return fmt.Errorf("invalid ID %s", rs.Primary.ID)
		}
		apiId := idSplit[0]
		ruleId := idSplit[1]

		service := TsfService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
		res, err := service.DescribeTsfApiRateLimitRuleById(ctx, apiId, ruleId)
		if err != nil {
			return err
		}

		if res == nil {
			return fmt.Errorf("tsf ApiRateLimitRule %s is not found", rs.Primary.ID)
		}

		return nil
	}
}

const testAccTsfApiRateLimitRuleVar = `
variable "api_id" {
	default = "` + defaultTsfApiId + `"
}
`

const testAccTsfApiRateLimitRule = testAccTsfApiRateLimitRuleVar + `

resource "tencentcloud_tsf_api_rate_limit_rule" "api_rate_limit_rule" {
	api_id = var.api_id
	max_qps = 10
	usable_status = "enabled"
}

`
