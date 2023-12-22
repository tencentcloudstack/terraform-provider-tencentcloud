package postgresql_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudPostgresqlDbInstanceClassesDataSource_basic(t *testing.T) {
	// t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPostgresqlDbInstanceClassesDataSource,
				PreConfig: func() {
					tcacctest.AccStepSetRegion(t, "ap-guangzhou")
					tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_COMMON)
				},
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_postgresql_db_instance_classes.db_instance_classes"),
					resource.TestCheckResourceAttr("data.tencentcloud_postgresql_db_instance_classes.db_instance_classes", "zone", "ap-guangzhou-7"),
					resource.TestCheckResourceAttr("data.tencentcloud_postgresql_db_instance_classes.db_instance_classes", "db_engine", "postgresql"),
					resource.TestCheckResourceAttr("data.tencentcloud_postgresql_db_instance_classes.db_instance_classes", "db_major_version", "13"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_postgresql_db_instance_classes.db_instance_classes", "class_info_set.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_postgresql_db_instance_classes.db_instance_classes", "class_info_set.0.spec_code"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_postgresql_db_instance_classes.db_instance_classes", "class_info_set.0.cpu"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_postgresql_db_instance_classes.db_instance_classes", "class_info_set.0.memory"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_postgresql_db_instance_classes.db_instance_classes", "class_info_set.0.max_storage"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_postgresql_db_instance_classes.db_instance_classes", "class_info_set.0.min_storage"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_postgresql_db_instance_classes.db_instance_classes", "class_info_set.0.qps"),
				),
			},
		},
	})
}

const testAccPostgresqlDbInstanceClassesDataSource = `
data "tencentcloud_postgresql_db_instance_classes" "db_instance_classes" {
  zone = "ap-guangzhou-7"
  db_engine = "postgresql"
  db_major_version = "13"
}

`
