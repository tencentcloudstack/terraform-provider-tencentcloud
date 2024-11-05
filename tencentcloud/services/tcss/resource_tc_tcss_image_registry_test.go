package tcss_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudTcssImageRegistryResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTcssImageRegistry,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_tcss_image_registry.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_tcss_image_registry.example", "name"),
					resource.TestCheckResourceAttrSet("tencentcloud_tcss_image_registry.example", "username"),
					resource.TestCheckResourceAttrSet("tencentcloud_tcss_image_registry.example", "password"),
					resource.TestCheckResourceAttrSet("tencentcloud_tcss_image_registry.example", "url"),
					resource.TestCheckResourceAttrSet("tencentcloud_tcss_image_registry.example", "registry_type"),
					resource.TestCheckResourceAttrSet("tencentcloud_tcss_image_registry.example", "net_type"),
					resource.TestCheckResourceAttrSet("tencentcloud_tcss_image_registry.example", "registry_version"),
					resource.TestCheckResourceAttrSet("tencentcloud_tcss_image_registry.example", "registry_region"),
					resource.TestCheckResourceAttrSet("tencentcloud_tcss_image_registry.example", "need_scan"),
				),
			},
		},
	})
}

const testAccTcssImageRegistry = `
resource "tencentcloud_tcss_image_registry" "example" {
  name             = "terraform"
  username         = "root"
  password         = "Password@demo"
  url              = "https://example.com"
  registry_type    = "harbor"
  net_type         = "public"
  registry_version = "V1"
  registry_region  = "default"
  need_scan        = true
  conn_detect_config {
    quuid = "backend"
    uuid  = "backend"
  }
}
`
