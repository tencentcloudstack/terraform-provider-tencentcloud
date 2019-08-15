package tencentcloud

import (
"testing"

"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceTencentCloudTke(t *testing.T) {

	key:="data.tencentcloud_kubernetes_clusters.name"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceTencentCloudTkeStr,
				Check: resource.ComposeTestCheckFunc(
					//name filter
					testAccCheckTencentCloudDataSourceID(key),
					resource.TestCheckResourceAttr(key, "cluster_name","terraform"),
					resource.TestCheckResourceAttrSet(key, "list.#"),

				),
			},
		},
	})
}

const testAccDataSourceTencentCloudTkeStr = `
data "tencentcloud_kubernetes_clusters"  "name" {
    cluster_name ="terraform"
}
`



