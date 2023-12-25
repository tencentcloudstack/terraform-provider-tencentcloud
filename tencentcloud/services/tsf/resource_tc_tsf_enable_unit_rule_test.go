package tsf_test

import (
	"context"
	"fmt"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svctsf "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tsf"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// go test -i; go test -test.run TestAccTencentCloudTsfEnableUnitRuleResource_basic -v
func TestAccTencentCloudTsfEnableUnitRuleResource_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_TSF) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckTsfEnableUnitRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTsfEnableUnitRule,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTsfEnableUnitRuleExists("tencentcloud_tsf_enable_unit_rule.enable_unit_rule"),
					resource.TestCheckResourceAttrSet("tencentcloud_tsf_enable_unit_rule.enable_unit_rule", "id"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_enable_unit_rule.enable_unit_rule", "switch", "enabled"),
				),
			},
			{
				ResourceName:      "tencentcloud_tsf_enable_unit_rule.enable_unit_rule",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccTsfEnableUnitRuleUp,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTsfEnableUnitRuleExists("tencentcloud_tsf_enable_unit_rule.enable_unit_rule"),
					resource.TestCheckResourceAttrSet("tencentcloud_tsf_enable_unit_rule.enable_unit_rule", "id"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_enable_unit_rule.enable_unit_rule", "switch", "disabled"),
				),
			},
		},
	})
}

func testAccCheckTsfEnableUnitRuleDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	service := svctsf.NewTsfService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_tsf_enable_unit_rule" {
			continue
		}

		res, err := service.DescribeTsfEnableUnitRuleById(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}

		if *res.Status != "disabled" {
			return fmt.Errorf("tsf enable unitRule %s still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckTsfEnableUnitRuleExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}

		service := svctsf.NewTsfService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
		res, err := service.DescribeTsfEnableUnitRuleById(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}

		if res == nil {
			return fmt.Errorf("tsf enable unitRule %s is not found", rs.Primary.ID)
		}

		return nil
	}
}

const testAccTsfEnableUnitRule = `

resource "tencentcloud_tsf_enable_unit_rule" "enable_unit_rule" {
	rule_id = "unit-rl-za8fcg7b"
	switch = "enabled"
}

`

const testAccTsfEnableUnitRuleUp = `

resource "tencentcloud_tsf_enable_unit_rule" "enable_unit_rule" {
	rule_id = "unit-rl-za8fcg7b"
	switch = "disabled"
}

`
