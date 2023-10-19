package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudDnspodRecordLineListDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
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
  domain = "dnspod.cn"
  domain_grade = "DP_FREE"
  domain_id = 123
  tags = {
    "createdBy" = "terraform"
  }
}

`
