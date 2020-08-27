package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudAuditCosRegionsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudAuditCosRegionsDataSourceConfigWithWebsite,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_audit_cos_regions.all"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_audit_cos_regions.all", "cos_region_list.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_audit_cos_regions.all", "cos_region_list.0.cos_region"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_audit_cos_regions.all", "cos_region_list.0.cos_region_name"),
				),
			},
		},
	})
}

const testAccTencentCloudAuditCosRegionsDataSourceConfigWithWebsite = `
data "tencentcloud_audit_cos_regions" "all" {
}
`
