package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudGaapCustomHeaderDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccGaapCustomHeaderDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_gaap_custom_header.custom_header"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_custom_header.custom_header", "headers.#"),
				),
			},
		},
	})
}

const testAccGaapCustomHeaderDataSource = `
data "tencentcloud_gaap_custom_header" "custom_header" {
	rule_id = "rule-9sdhv655"
}
`
