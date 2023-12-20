package audit_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudAuditsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudAuditsDataSourceConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_audits.all"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_audits.all", "audit_list.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_audits.all", "audit_list.0.name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_audits.all", "audit_list.0.cos_bucket"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_audits.all", "audit_list.0.log_file_prefix"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_audits.all", "audit_list.0.audit_switch"),
				),
			},
			{
				Config: testAccTencentCloudAuditsDataSourceConfigName,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_audits.name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_audits.name", "audit_list.#"),
				),
			},
		},
	})
}

const testAccTencentCloudAuditsDataSourceConfigBasic = `
data "tencentcloud_audits" "all" {
}
`

const testAccTencentCloudAuditsDataSourceConfigName = `
data "tencentcloud_audits" "name" {
  name = "test"
}
`
