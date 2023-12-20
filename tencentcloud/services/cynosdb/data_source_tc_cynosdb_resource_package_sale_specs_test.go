package cynosdb_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudCynosdbResourcePackageSaleSpecsDataSource_basic -v
func TestAccTencentCloudCynosdbResourcePackageSaleSpecsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCynosdbResourcePackageSaleSpecsDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_cynosdb_resource_package_sale_specs.resource_package_sale_specs"),
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
