package dnspod_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudDnspodRecordListDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PREPAY) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDnspodRecordListDataSource,
				Check:  resource.ComposeTestCheckFunc(tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_dnspod_record_list.record_list")),
			},
		},
	})
}

const testAccDnspodRecordListDataSource = `

data "tencentcloud_dnspod_record_list" "record_list" {
  domain = "iac-tf.cloud"
  # domain_id = 123
  # sub_domain = "www"
  record_type = ["A", "NS", "CNAME", "NS", "AAAA"]
  # record_line = [""]
  group_id = []
  keyword = ""
  sort_field = "UPDATED_ON"
  sort_type = "DESC"
  record_value = "bicycle.dnspod.net"
  record_status = ["ENABLE"]
  weight_begin = 0
  weight_end = 100
  mx_begin = 0
  mx_end = 10
  ttl_begin = 1
  ttl_end = 864000
  updated_at_begin = "2021-09-07"
  updated_at_end = "2023-12-07"
  remark = ""
  is_exact_sub_domain = true
  # project_id = -1
}

`
