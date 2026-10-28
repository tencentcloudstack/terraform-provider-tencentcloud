package bh_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudBhBindDeviceAccountKubeconfigResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccBhBindDeviceAccountKubeconfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_bh_bind_device_account_kubeconfig.example", "id"),
					resource.TestCheckResourceAttr("tencentcloud_bh_bind_device_account_kubeconfig.example", "account_id", "12345"),
					resource.TestCheckResourceAttr("tencentcloud_bh_bind_device_account_kubeconfig.example", "manage_dimension", "1"),
				),
			},
			{
				Config: testAccBhBindDeviceAccountKubeconfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_bh_bind_device_account_kubeconfig.example", "kubeconfig", "test-kubeconfig-v2"),
				),
			},
		},
	})
}

const testAccBhBindDeviceAccountKubeconfig = `
resource "tencentcloud_bh_bind_device_account_kubeconfig" "example" {
  account_id       = 12345
  kubeconfig       = "test-kubeconfig-v1"
  manage_dimension = 1
}
`

const testAccBhBindDeviceAccountKubeconfigUpdate = `
resource "tencentcloud_bh_bind_device_account_kubeconfig" "example" {
  account_id       = 12345
  kubeconfig       = "test-kubeconfig-v2"
  manage_dimension = 1
}
`
