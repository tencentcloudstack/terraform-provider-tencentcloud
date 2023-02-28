package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudCssDomainsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCssDomainsDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_css_domains.domains"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_css_domains.domains", "domain_list.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_css_domains.domains", "domain_list.0.name"),
				),
			},
		},
	})
}

const testAccCssDomainsDataSource = `

data "tencentcloud_css_domains" "domains" {
  domain_type = 0
  play_type = 1
  is_delay_live = 0
}

`
