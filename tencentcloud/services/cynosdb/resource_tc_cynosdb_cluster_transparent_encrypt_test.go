package cynosdb_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudCynosdbClusterTransparentEncryptResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCynosdbClusterTransparentEncrypt,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_cluster_transparent_encrypt.cynosdb_cluster_transparent_encrypt", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_cluster_transparent_encrypt.cynosdb_cluster_transparent_encrypt", "key_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_cluster_transparent_encrypt.cynosdb_cluster_transparent_encrypt", "key_region"),
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_cluster_transparent_encrypt.cynosdb_cluster_transparent_encrypt", "key_type"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_cluster_transparent_encrypt.cynosdb_cluster_transparent_encrypt", "is_open_global_encryption", "false"),
				),
			},
			{
				ResourceName: "tencentcloud_cynosdb_cluster_transparent_encrypt.cynosdb_cluster_transparent_encrypt",
				ImportState:  true,
			},
		},
	})
}

const testAccCynosdbClusterTransparentEncrypt = testAccCynosdbCluster + `
resource "tencentcloud_cynosdb_cluster_transparent_encrypt" "cynosdb_cluster_transparent_encrypt" {
  cluster_id                = tencentcloud_cynosdb_cluster.foo.id
  is_open_global_encryption = false
  key_id                    = "f063c18b-654b-11ef-9d9f-525400d3a886"
  key_region                = "ap-guangzhou"
  key_type                  = "custom"
}
`
