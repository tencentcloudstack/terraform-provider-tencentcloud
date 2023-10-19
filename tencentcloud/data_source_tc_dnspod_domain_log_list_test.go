package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudDnspodDomainLogListDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDnspodDomainLogListDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_dnspod_domain_log_list.domain_log_list")),
			},
		},
	})
}

const testAccDnspodDomainLogListDataSource = `

data "tencentcloud_dnspod_domain_log_list" "domain_log_list" {
  domain = "dnspod.cn"
  domain_id = 123
  tags = {
    "createdBy" = "terraform"
  }
}

`
