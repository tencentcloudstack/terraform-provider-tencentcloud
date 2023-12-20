package css_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudCssDomainsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCssDomainsDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_css_domains.domains"),
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
