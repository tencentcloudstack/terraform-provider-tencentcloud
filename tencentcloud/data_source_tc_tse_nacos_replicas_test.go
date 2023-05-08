package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudTseNacosReplicasDataSource_basic -v
func TestAccTencentCloudTseNacosReplicasDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTseNacosReplicasDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_tse_nacos_replicas.nacos_replicas"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_nacos_replicas.nacos_replicas", "replicas.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_nacos_replicas.nacos_replicas", "replicas.0.name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_nacos_replicas.nacos_replicas", "replicas.0.status"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_nacos_replicas.nacos_replicas", "replicas.0.zone"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_nacos_replicas.nacos_replicas", "replicas.0.zone_id"),
				),
			},
		},
	})
}

const testAccTseNacosReplicasDataSource = `

data "tencentcloud_tse_nacos_replicas" "nacos_replicas" {
  instance_id = "ins-15137c53"
}

`
