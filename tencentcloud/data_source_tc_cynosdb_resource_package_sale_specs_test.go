package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudCynosdbResourcePackageSaleSpecsDataSource_basic -v
func TestAccTencentCloudCynosdbResourcePackageSaleSpecsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCynosdbResourcePackageSaleSpecsDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_cynosdb_resource_package_sale_specs.resource_package_sale_specs"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cynosdb_resource_package_sale_specs.resource_package_sale_specs", "detail.#"),
				),
			},
		},
	})
}

const testAccCynosdbResourcePackageSaleSpecsDataSource = `
data "tencentcloud_cynosdb_resource_package_sale_specs" "resource_package_sale_specs" {
  instance_type  = "cynosdb-serverless"
  package_region = "china"
  package_type   = "CCU"
}
`
