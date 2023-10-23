package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudDnspodDomainLogListDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY) },
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
  domain = "iac-tf.cloud"
  # domain_id = 123
}

`
