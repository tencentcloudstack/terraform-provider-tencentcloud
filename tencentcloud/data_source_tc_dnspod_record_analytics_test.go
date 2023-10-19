package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudDnspodRecordAnalyticsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDnspodRecordAnalyticsDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_dnspod_record_analytics.record_analytics")),
			},
		},
	})
}

const testAccDnspodRecordAnalyticsDataSource = `

data "tencentcloud_dnspod_record_analytics" "record_analytics" {
  domain = "dnspod.cn"
  start_date = "2023-09-07 00:00:00 +0000 UTC"
  end_date = "2023-09-07 00:00:00 +0000 UTC"
  subdomain = "www"
  dns_format = "HOUR"
  domain_id = 123
  tags = {
    "createdBy" = "terraform"
  }
}

`
