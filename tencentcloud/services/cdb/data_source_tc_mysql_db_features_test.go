package cdb_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudMysqlDbFeaturesDataSource_basic -v
func TestAccTencentCloudMysqlDbFeaturesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMysqlDbFeaturesDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_mysql_db_features.db_features"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_db_features.db_features", "id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_db_features.db_features", "audit_need_upgrade"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_db_features.db_features", "current_sub_version"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_db_features.db_features", "encryption_need_upgrade"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_db_features.db_features", "instance_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_db_features.db_features", "is_remote_ro"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_db_features.db_features", "is_support_audit"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_db_features.db_features", "is_support_encryption"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_db_features.db_features", "is_support_update_sub_version"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_db_features.db_features", "target_sub_version"),
				),
			},
		},
	})
}

const testAccMysqlDbFeaturesDataSourceVar = `
variable "instance_id" {
	default = "` + tcacctest.DefaultDbBrainInstanceId + `"
  }
`

const testAccMysqlDbFeaturesDataSource = testAccMysqlDbFeaturesDataSourceVar + `

data "tencentcloud_mysql_db_features" "db_features" {
  instance_id = var.instance_id
}

`
