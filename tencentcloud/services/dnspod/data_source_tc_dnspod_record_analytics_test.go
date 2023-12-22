package dnspod_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudNeedFixDnspodRecordAnalyticsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PREPAY) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDnspodRecordAnalyticsDataSource,
				Check:  resource.ComposeTestCheckFunc(tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_dnspod_record_analytics.record_analytics")),
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
