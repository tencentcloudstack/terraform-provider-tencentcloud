package cynosdb_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudCynosdbResourcePackageListDataSource_basic -v
func TestAccTencentCloudCynosdbResourcePackageListDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PREPAY) },
				Config:    testAccCynosdbResourcePackageListDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_cynosdb_resource_package_list.resource_package_list"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cynosdb_resource_package_list.resource_package_list", "resource_package_list.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cynosdb_resource_package_list.resource_package_list", "resource_package_list.0.app_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cynosdb_resource_package_list.resource_package_list", "resource_package_list.0.package_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cynosdb_resource_package_list.resource_package_list", "resource_package_list.0.package_name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cynosdb_resource_package_list.resource_package_list", "resource_package_list.0.package_type"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cynosdb_resource_package_list.resource_package_list", "resource_package_list.0.package_region"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cynosdb_resource_package_list.resource_package_list", "resource_package_list.0.status"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cynosdb_resource_package_list.resource_package_list", "resource_package_list.0.package_total_spec"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cynosdb_resource_package_list.resource_package_list", "resource_package_list.0.package_used_spec"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cynosdb_resource_package_list.resource_package_list", "resource_package_list.0.has_quota"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cynosdb_resource_package_list.resource_package_list", "resource_package_list.0.bind_instance_infos.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cynosdb_resource_package_list.resource_package_list", "resource_package_list.0.start_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cynosdb_resource_package_list.resource_package_list", "resource_package_list.0.expire_time"),
				),
			},
		},
	})
}

const testAccCynosdbResourcePackageListDataSource = `
data "tencentcloud_cynosdb_resource_package_list" "resource_package_list" {
  package_id      = ["package-hy4d2ppl"]
  package_name    = ["keep-package-disk"]
  package_type    = ["DISK"]
  package_region  = ["china"]
  status          = ["using"]
  order_by        = ["startTime"]
  order_direction = "DESC"
}
`
