package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudNeedFixCynosdbResourcePackageResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCynosdbResourcePackage,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_resource_package.resource_package", "id")),
			},
			{
				ResourceName:      "tencentcloud_cynosdb_resource_package.resource_package",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCynosdbResourcePackage = `

resource "tencentcloud_cynosdb_resource_package" "resource_package" {
  instance_type = "cdb"
  package_region = "china"
  package_type = "CCU"
  package_version = "base"
  package_spec = 
  expire_day = 180
  package_count = 1
  package_name = "PackageName"
}

`
