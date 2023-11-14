package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

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
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_tse_nacos_server_interfaces.nacos_server_interfaces")),
			},
		},
	})
}

const testAccTseNacosServerInterfacesDataSource = `

data "tencentcloud_tse_nacos_server_interfaces" "nacos_server_interfaces" {
  instance_id = "ins-xxxxxx"
  }

`
