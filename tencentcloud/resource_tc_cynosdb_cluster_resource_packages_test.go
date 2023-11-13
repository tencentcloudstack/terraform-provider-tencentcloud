package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCynosdbClusterResourcePackagesResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCynosdbClusterResourcePackages,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_cluster_resource_packages.cluster_resource_packages", "id")),
			},
			{
				ResourceName:      "tencentcloud_cynosdb_cluster_resource_packages.cluster_resource_packages",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCynosdbClusterResourcePackages = `

resource "tencentcloud_cynosdb_cluster_resource_packages" "cluster_resource_packages" {
  package_ids = 
  cluster_id = "cynosdb-qwerty"
}

`
