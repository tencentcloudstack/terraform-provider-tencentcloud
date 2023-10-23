package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudNeedFixDnspodRecordAnalyticsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY) },
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
  domain = "iac-tf.cloud"
  start_date = "2023-09-07"
  end_date = "2023-11-07"
  subdomain = "www"
  dns_format = "HOUR"
  # domain_id = 123
}

`
