package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudCynosdbResourcePackageListDataSource_basic -v
func TestAccTencentCloudCynosdbResourcePackageListDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCynosdbResourcePackageListDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_cynosdb_resource_package_list.resource_package_list"),
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
