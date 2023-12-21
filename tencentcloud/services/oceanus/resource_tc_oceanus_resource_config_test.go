package oceanus_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudNeedFixOceanusResourceConfigResource_basic -v
func TestAccTencentCloudNeedFixOceanusResourceConfigResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOceanusResourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_oceanus_resource_config.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_oceanus_resource_config.example", "resource_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_oceanus_resource_config.example", "remark"),
					resource.TestCheckResourceAttrSet("tencentcloud_oceanus_resource_config.example", "work_space_id"),
				),
			},
		},
	})
}

const testAccOceanusResourceConfig = `
resource "tencentcloud_oceanus_resource_config" "example" {
  resource_id = "resource-8y9lzcuz"
  resource_loc {
    storage_type = 1
    param {
      bucket = "keep-terraform-1257058945"
      path   = "OceanusResource/junit-4.13.2.jar"
      region = "ap-guangzhou"
    }
  }

  remark        = "remark."
  work_space_id = "space-2idq8wbr"
}
`
