package bh_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudNeedFixDasbBindDeviceResourceResource_basic -v
func TestAccTencentCloudNeedFixDasbBindDeviceResourceResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDasbBindDeviceResource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_dasb_bind_device_resource.example", "id"),
					resource.TestCheckResourceAttr("tencentcloud_dasb_bind_device_resource.example", "resource_id", "bh-saas-ocmzo6lgxiv"),
				),
			},
		},
	})
}

const testAccDasbBindDeviceResource = `
resource "tencentcloud_dasb_bind_device_resource" "example" {
  resource_id   = "bh-saas-ocmzo6lgxiv"
  device_id_set = [17, 18]
}
`
