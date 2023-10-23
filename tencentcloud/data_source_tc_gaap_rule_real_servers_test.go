package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudGaapRuleRealServersDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccGaapRuleRealServersDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_gaap_rule_real_servers.rule_real_servers"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_rule_real_servers.rule_real_servers", "bind_real_server_set.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_rule_real_servers.rule_real_servers", "real_server_set.#"),
				),
			},
		},
	})
}

const testAccGaapRuleRealServersDataSource = `
data "tencentcloud_gaap_rule_real_servers" "rule_real_servers" {
	rule_id = "rule-9sdhv655"
}
`
