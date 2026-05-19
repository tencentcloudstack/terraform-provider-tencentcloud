package bh_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudBhBindDeviceResourceConfigResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccBhBindDeviceResourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_bh_bind_device_resource_config.example", "id"),
					resource.TestCheckResourceAttr("tencentcloud_bh_bind_device_resource_config.example", "resource_id", "bh-saas-jn2p3"),
					resource.TestCheckResourceAttr("tencentcloud_bh_bind_device_resource_config.example", "device_id_set.#", "2"),
				),
			},
			{
				Config: testAccBhBindDeviceResourceConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_bh_bind_device_resource_config.example", "device_id_set.#", "1"),
				),
			},
			{
				ResourceName:      "tencentcloud_bh_bind_device_resource_config.example",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"manage_kubeconfig"},
			},
		},
	})
}

const testAccBhBindDeviceResourceConfig = `
resource "tencentcloud_bh_bind_device_resource_config" "example" {
  resource_id   = "bh-saas-jn2p3"
  device_id_set = [5186, 5187]
  domain_id     = "net-4sovwr11w7"
}
`

const testAccBhBindDeviceResourceConfigUpdate = `
resource "tencentcloud_bh_bind_device_resource_config" "example" {
  resource_id   = "bh-saas-jn2p3"
  device_id_set = [5186]
  domain_id     = "net-4sovwr11w7"
}
`
