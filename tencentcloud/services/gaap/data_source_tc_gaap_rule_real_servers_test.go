package gaap_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudGaapRuleRealServersDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PREPAY) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccGaapRuleRealServersDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_gaap_rule_real_servers.rule_real_servers"),
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
