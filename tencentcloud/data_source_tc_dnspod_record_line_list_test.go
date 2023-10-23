package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudDnspodRecordLineListDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDnspodRecordLineListDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_dnspod_record_line_list.record_line_list")),
			},
		},
	})
}

const testAccDnspodRecordLineListDataSource = `

data "tencentcloud_dnspod_record_line_list" "record_line_list" {
  domain = "iac-tf.cloud"
  domain_grade = "DP_FREE"
  # domain_id = 123
}

`
