package audit_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudAuditCosRegionsDataSource(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudAuditCosRegionsDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_audit_cos_regions.all"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_audit_cos_regions.all", "audit_cos_region_list.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_audit_cos_regions.all", "audit_cos_region_list.0.cos_region"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_audit_cos_regions.all", "audit_cos_region_list.0.cos_region_name"),
				),
			},
		},
	})
}

const testAccTencentCloudAuditCosRegionsDataSource = `
data "tencentcloud_audit_cos_regions" "all" {
}
`
