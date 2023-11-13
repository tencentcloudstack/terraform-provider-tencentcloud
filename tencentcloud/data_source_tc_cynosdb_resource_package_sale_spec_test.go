package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCynosdbResourcePackageSaleSpecDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCynosdbResourcePackageSaleSpecDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_cynosdb_resource_package_sale_spec.resource_package_sale_spec")),
			},
		},
	})
}

const testAccCynosdbResourcePackageSaleSpecDataSource = `

data "tencentcloud_cynosdb_resource_package_sale_spec" "resource_package_sale_spec" {
  instance_type = "cynosdb-serverless"
  package_region = "china"
  package_type = "CCU"
  }

`
