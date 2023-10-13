package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudDnspodDomainListDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDnspodDomain_listDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_dnspod_domain_list.domain_list")),
			},
		},
	})
}

const testAccDnspodDomain_listDataSource = `

data "tencentcloud_dnspod_domain_list" "domain_list" {
  type = "ALL"
  #  group_id = 
  keyword = "keyword_demo"
  sort_field = "UPDATED_ON"
  sort_type = "DESC"
  # status = 
  # package = 
  remark = "remark_demo"
  updated_at_begin = "2021-05-01 03:00:00"
  updated_at_end = "2021-05-10 20:00:00"
  record_count_begin = 0
  record_count_end = 100
  project_id = -1
      tags = {
    "createdBy" = "terraform"
  }
}

`
