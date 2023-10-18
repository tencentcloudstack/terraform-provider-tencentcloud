package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudDnspodDomainAnalyticsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDnspodDomainAnalyticsDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_dnspod_domain_analytics.domain_analytics")),
			},
		},
	})
}

const testAccDnspodDomainAnalyticsDataSource = `

data "tencentcloud_dnspod_domain_analytics" "domain_analytics" {
  domain = "iac-tf.cloud"
  start_date = "2023-10-07"
  end_date = "2023-10-12"
  dns_format = "HOUR"
  # domain_id = 123
}

`
