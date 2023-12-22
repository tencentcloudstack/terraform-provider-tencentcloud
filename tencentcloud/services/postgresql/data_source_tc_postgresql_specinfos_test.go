package postgresql_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudPostgresqlSpecinfosDataSource(t *testing.T) {
	// t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccStepSetRegion(t, "ap-guangzhou")
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudPostgresqlSpecinfosConfigBasic,
				PreConfig: func() {
					tcacctest.AccStepSetRegion(t, "ap-guangzhou")
					tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_COMMON)
				},
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_postgresql_specinfos.foo"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_postgresql_specinfos.foo", "list.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_postgresql_specinfos.foo", "list.0.id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_postgresql_specinfos.foo", "list.0.memory"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_postgresql_specinfos.foo", "list.0.storage_min"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_postgresql_specinfos.foo", "list.0.storage_max"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_postgresql_specinfos.foo", "list.0.cpu"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_postgresql_specinfos.foo", "list.0.qps"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_postgresql_specinfos.foo", "list.0.engine_version"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_postgresql_specinfos.foo", "list.0.engine_version_name"),
				),
			},
		},
	})
}

const testAccTencentCloudPostgresqlSpecinfosConfigBasic = `
data "tencentcloud_postgresql_specinfos" "foo"{
  availability_zone = "ap-guangzhou-3"
}
`
