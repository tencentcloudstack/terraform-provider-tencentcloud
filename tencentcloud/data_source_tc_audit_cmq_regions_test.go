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
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_audit_cmq_regions.all"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_audit_cmq_regions.all", "audit_cmq_region_list.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_audit_cmq_regions.all", "audit_cmq_region_list.0.cmq_region"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_audit_cmq_regions.all", "audit_cmq_region_list.0.cmq_region_name"),
				),
			},
		},
	})
}

const testAccTencentCloudAuditCmqRegionsDataSourceConfigWithWebsite = `
data "tencentcloud_audit_cmq_regions" "all" {
}
`
