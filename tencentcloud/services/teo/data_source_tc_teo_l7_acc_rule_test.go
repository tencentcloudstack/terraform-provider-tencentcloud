package teo_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudTeoL7AccRuleDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{{
			Config: testAccTeoL7AccRuleDataSource,
			Check: resource.ComposeTestCheckFunc(
				tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_teo_l7_acc_rule.example"),
				resource.TestCheckResourceAttrSet("data.tencentcloud_teo_l7_acc_rule.example", "zone_id"),
			),
		}},
	})
}

func TestAccTencentCloudTeoL7AccRuleDataSource_offset_zero(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{{
			Config: testAccTeoL7AccRuleDataSourceOffsetZero,
			Check: resource.ComposeTestCheckFunc(
				tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_teo_l7_acc_rule.example"),
				resource.TestCheckResourceAttr("data.tencentcloud_teo_l7_acc_rule.example", "zone_id", "zone-3fkff38fyw8s"),
				resource.TestCheckResourceAttr("data.tencentcloud_teo_l7_acc_rule.example", "offset", "0"),
			),
		}},
	})
}

func TestAccTencentCloudTeoL7AccRuleDataSource_offset_positive(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{{
			Config: testAccTeoL7AccRuleDataSourceOffsetPositive,
			Check: resource.ComposeTestCheckFunc(
				tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_teo_l7_acc_rule.example"),
				resource.TestCheckResourceAttr("data.tencentcloud_teo_l7_acc_rule.example", "zone_id", "zone-3fkff38fyw8s"),
				resource.TestCheckResourceAttr("data.tencentcloud_teo_l7_acc_rule.example", "offset", "10"),
			),
		}},
	})
}

func TestAccTencentCloudTeoL7AccRuleDataSource_offset_negative(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{{
			Config:      testAccTeoL7AccRuleDataSourceOffsetNegative,
			ExpectError: tcacctest.RegexpContains(`offset: value must be at least 0`),
		}},
	})
}

func TestAccTencentCloudTeoL7AccRuleDataSource_byRuleId(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{{
			Config: testAccTeoL7AccRuleDataSourceByRuleId,
			Check: resource.ComposeTestCheckFunc(
				tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_teo_l7_acc_rule.example"),
				resource.TestCheckResourceAttrSet("data.tencentcloud_teo_l7_acc_rule.example", "zone_id"),
				resource.TestCheckResourceAttrSet("data.tencentcloud_teo_l7_acc_rule.example", "rule_id"),
			),
		}},
	})
}

const testAccTeoL7AccRuleDataSource = `
data "tencentcloud_teo_l7_acc_rule" "example" {
  zone_id = "zone-3fkff38fyw8s"
}
`

const testAccTeoL7AccRuleDataSourceOffsetZero = `
data "tencentcloud_teo_l7_acc_rule" "example" {
  zone_id = "zone-3fkff38fyw8s"
  offset  = 0
}
`

const testAccTeoL7AccRuleDataSourceOffsetPositive = `
data "tencentcloud_teo_l7_acc_rule" "example" {
  zone_id = "zone-3fkff38fyw8s"
  offset  = 10
}
`

const testAccTeoL7AccRuleDataSourceOffsetNegative = `
data "tencentcloud_teo_l7_acc_rule" "example" {
  zone_id = "zone-3fkff38fyw8s"
  offset  = -5
}
`

const testAccTeoL7AccRuleDataSourceByRuleId = `
data "tencentcloud_teo_l7_acc_rule" "example" {
  zone_id = "zone-3fkff38fyw8s"
  rule_id = "rule-test-id"
}
`
