package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
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
				Config: testAccDnspodDomainListDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_dnspod_domain_list.domain_list")),
			},
		},
	})
}

const testAccDnspodDomainListDataSource = `

data "tencentcloud_dnspod_domain_list" "domain_list" {
	type = "ALL"
	group_id = [1]
	keyword = ""
	sort_field = "UPDATED_ON"
	sort_type = "DESC"
	status = ["PAUSE"]
	package = [""]
	remark = ""
	updated_at_begin = "2021-05-01 03:00:00"
	updated_at_end = "2024-05-10 20:00:00"
	record_count_begin = 0
	record_count_end = 100
	project_id = -1
	tags {
		tag_key = "created_by"
		tag_value = ["terraform"]
	}
}

`
