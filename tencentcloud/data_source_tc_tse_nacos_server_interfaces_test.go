package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudTseNacosServerInterfacesDataSource_basic -v
func TestAccTencentCloudTseNacosServerInterfacesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTseNacosServerInterfacesDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_tse_nacos_server_interfaces.nacos_server_interfaces"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_nacos_server_interfaces.nacos_server_interfaces", "instance_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_nacos_server_interfaces.nacos_server_interfaces", "content.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_nacos_server_interfaces.nacos_server_interfaces", "content.0.interface"),
				),
			},
		},
	})
}

const testAccTseNacosServerInterfacesDataSource = `

data "tencentcloud_tse_nacos_server_interfaces" "nacos_server_interfaces" {
	instance_id = "ins-15137c53"
}

`
