package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudDnspodRecordListDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDnspodRecordListDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_dnspod_record_list.record_list")),
			},
		},
	})
}

const testAccDnspodRecordListDataSource = `

data "tencentcloud_dnspod_record_list" "record_list" {
  domain = "dnspod.cn"
  domain_id = 123
  sub_domain = "www"
  record_type = 
  record_line = 
  group_id = 
  keyword = "keyword_demo"
  sort_field = "UPDATED_ON"
  sort_type = "DESC"
  record_value = "129.29.29.29"
  record_status = 
  weight_begin = 0
  weight_end = 100
  mx_begin = 0
  mx_end = 10
  ttl_begin = 1
  ttl_end = 600
  updated_at_begin = "2023-09-07 00:00:00 +0000 UTC"
  updated_at_end = "2023-09-07 00:00:00 +0000 UTC"
  remark = "remark_demo"
  is_exact_sub_domain = true
  project_id = -1
      tags = {
    "createdBy" = "terraform"
  }
}

`
