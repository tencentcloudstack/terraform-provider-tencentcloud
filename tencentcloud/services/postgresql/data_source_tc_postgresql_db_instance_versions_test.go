package postgresql_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudPostgresqlDbInstanceVersionsDataSource_basic(t *testing.T) {
	// t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccStepSetRegion(t, "ap-guangzhou")
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPostgresqlDbInstanceVersionsDataSource,
				PreConfig: func() {
					tcacctest.AccStepSetRegion(t, "ap-guangzhou")
					tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_COMMON)
				},
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_postgresql_db_instance_versions.db_instance_versions"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_postgresql_db_instance_versions.db_instance_versions", "version_set.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_postgresql_db_instance_versions.db_instance_versions", "version_set.0.db_engine"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_postgresql_db_instance_versions.db_instance_versions", "version_set.0.db_version"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_postgresql_db_instance_versions.db_instance_versions", "version_set.0.db_major_version"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_postgresql_db_instance_versions.db_instance_versions", "version_set.0.db_kernel_version"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_postgresql_db_instance_versions.db_instance_versions", "version_set.0.supported_feature_names.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_postgresql_db_instance_versions.db_instance_versions", "version_set.0.status"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_postgresql_db_instance_versions.db_instance_versions", "version_set.0.available_upgrade_target.#"),
				),
			},
		},
	})
}

const testAccPostgresqlDbInstanceVersionsDataSource = `

data "tencentcloud_postgresql_db_instance_versions" "db_instance_versions" {}

`
