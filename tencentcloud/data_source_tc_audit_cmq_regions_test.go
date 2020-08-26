package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudAuditCmqRegionsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudAuditCmqRegionsDataSourceConfigWithWebsite,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_audit_cmq_regions.filter"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_audit_cmq_regions.filter", "cmq_region_list.#"),
				),
			},
		},
	})
}

const testAccTencentCloudAuditCmqRegionsDataSourceConfigWithWebsite = `
data "tencentcloud_audit_cmq_regions" "filter" {
	website_type = "zh"
}
`
