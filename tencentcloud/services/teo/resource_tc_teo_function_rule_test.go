package teo_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudTeoFunctionRuleResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTeoFunctionRule,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_teo_function_rule.teo_function_rule", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_function_rule.teo_function_rule", "function_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_function_rule.teo_function_rule", "zone_id"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function_rule.teo_function_rule", "remark", "aaa"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function_rule.teo_function_rule", "function_rule_conditions.#", "2"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function_rule.teo_function_rule", "function_rule_conditions.0.rule_conditions.#", "2"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function_rule.teo_function_rule", "function_rule_conditions.0.rule_conditions.0.ignore_case", "false"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function_rule.teo_function_rule", "function_rule_conditions.0.rule_conditions.0.operator", "equal"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function_rule.teo_function_rule", "function_rule_conditions.0.rule_conditions.0.target", "host"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function_rule.teo_function_rule", "function_rule_conditions.0.rule_conditions.0.values.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function_rule.teo_function_rule", "function_rule_conditions.0.rule_conditions.0.values.0", "aaa.makn.cn"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function_rule.teo_function_rule", "function_rule_conditions.0.rule_conditions.1.ignore_case", "false"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function_rule.teo_function_rule", "function_rule_conditions.0.rule_conditions.1.operator", "equal"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function_rule.teo_function_rule", "function_rule_conditions.0.rule_conditions.1.target", "extension"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function_rule.teo_function_rule", "function_rule_conditions.0.rule_conditions.1.values.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function_rule.teo_function_rule", "function_rule_conditions.0.rule_conditions.1.values.0", ".txt"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function_rule.teo_function_rule", "function_rule_conditions.1.rule_conditions.#", "2"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function_rule.teo_function_rule", "function_rule_conditions.1.rule_conditions.0.ignore_case", "false"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function_rule.teo_function_rule", "function_rule_conditions.1.rule_conditions.0.operator", "notequal"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function_rule.teo_function_rule", "function_rule_conditions.1.rule_conditions.0.target", "host"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function_rule.teo_function_rule", "function_rule_conditions.1.rule_conditions.0.values.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function_rule.teo_function_rule", "function_rule_conditions.1.rule_conditions.0.values.0", "aaa.makn.cn"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function_rule.teo_function_rule", "function_rule_conditions.1.rule_conditions.1.ignore_case", "false"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function_rule.teo_function_rule", "function_rule_conditions.1.rule_conditions.1.operator", "equal"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function_rule.teo_function_rule", "function_rule_conditions.1.rule_conditions.1.target", "extension"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function_rule.teo_function_rule", "function_rule_conditions.1.rule_conditions.1.values.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function_rule.teo_function_rule", "function_rule_conditions.1.rule_conditions.1.values.0", ".png"),
				),
			},
			{
				ResourceName:      "tencentcloud_teo_function_rule.teo_function_rule",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccTeoFunctionRuleUp,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_teo_function_rule.teo_function_rule", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_function_rule.teo_function_rule", "function_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_function_rule.teo_function_rule", "zone_id"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function_rule.teo_function_rule", "remark", "bbb"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function_rule.teo_function_rule", "function_rule_conditions.#", "2"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function_rule.teo_function_rule", "function_rule_conditions.0.rule_conditions.#", "2"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function_rule.teo_function_rule", "function_rule_conditions.0.rule_conditions.0.ignore_case", "false"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function_rule.teo_function_rule", "function_rule_conditions.0.rule_conditions.0.operator", "notequal"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function_rule.teo_function_rule", "function_rule_conditions.0.rule_conditions.0.target", "host"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function_rule.teo_function_rule", "function_rule_conditions.0.rule_conditions.0.values.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function_rule.teo_function_rule", "function_rule_conditions.0.rule_conditions.0.values.0", "aaa.makn.cn"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function_rule.teo_function_rule", "function_rule_conditions.0.rule_conditions.1.ignore_case", "false"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function_rule.teo_function_rule", "function_rule_conditions.0.rule_conditions.1.operator", "equal"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function_rule.teo_function_rule", "function_rule_conditions.0.rule_conditions.1.target", "extension"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function_rule.teo_function_rule", "function_rule_conditions.0.rule_conditions.1.values.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function_rule.teo_function_rule", "function_rule_conditions.0.rule_conditions.1.values.0", ".txt"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function_rule.teo_function_rule", "function_rule_conditions.1.rule_conditions.#", "2"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function_rule.teo_function_rule", "function_rule_conditions.1.rule_conditions.0.ignore_case", "false"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function_rule.teo_function_rule", "function_rule_conditions.1.rule_conditions.0.operator", "notequal"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function_rule.teo_function_rule", "function_rule_conditions.1.rule_conditions.0.target", "host"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function_rule.teo_function_rule", "function_rule_conditions.1.rule_conditions.0.values.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function_rule.teo_function_rule", "function_rule_conditions.1.rule_conditions.0.values.0", "aaa.makn.cn"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function_rule.teo_function_rule", "function_rule_conditions.1.rule_conditions.1.ignore_case", "false"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function_rule.teo_function_rule", "function_rule_conditions.1.rule_conditions.1.operator", "equal"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function_rule.teo_function_rule", "function_rule_conditions.1.rule_conditions.1.target", "extension"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function_rule.teo_function_rule", "function_rule_conditions.1.rule_conditions.1.values.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function_rule.teo_function_rule", "function_rule_conditions.1.rule_conditions.1.values.0", ".png"),
				),
			},
		},
	})
}

const testAccTeoFunctionRule = `

resource "tencentcloud_teo_function_rule" "teo_function_rule" {
    function_id   = "ef-txx7fnua"
    remark        = "aaa"
    zone_id       = "zone-2qtuhspy7cr6"

    function_rule_conditions {
        rule_conditions {
            ignore_case = false
            name        = null
            operator    = "equal"
            target      = "host"
            values      = [
                "aaa.makn.cn",
            ]
        }
        rule_conditions {
            ignore_case = false
            name        = null
            operator    = "equal"
            target      = "extension"
            values      = [
                ".txt",
            ]
        }
    }
    function_rule_conditions {
        rule_conditions {
            ignore_case = false
            name        = null
            operator    = "notequal"
            target      = "host"
            values      = [
                "aaa.makn.cn",
            ]
        }
        rule_conditions {
            ignore_case = false
            name        = null
            operator    = "equal"
            target      = "extension"
            values      = [
                ".png",
            ]
        }
    }
}
`

const testAccTeoFunctionRuleUp = `

resource "tencentcloud_teo_function_rule" "teo_function_rule" {
    function_id   = "ef-txx7fnua"
    remark        = "bbb"
    zone_id       = "zone-2qtuhspy7cr6"

    function_rule_conditions {
        rule_conditions {
            ignore_case = false
            name        = null
            operator    = "notequal"
            target      = "host"
            values      = [
                "aaa.makn.cn",
            ]
        }
        rule_conditions {
            ignore_case = false
            name        = null
            operator    = "equal"
            target      = "extension"
            values      = [
                ".txt",
            ]
        }
    }
    function_rule_conditions {
        rule_conditions {
            ignore_case = false
            name        = null
            operator    = "notequal"
            target      = "host"
            values      = [
                "aaa.makn.cn",
            ]
        }
        rule_conditions {
            ignore_case = false
            name        = null
            operator    = "equal"
            target      = "extension"
            values      = [
                ".png",
            ]
        }
    }
}
`
