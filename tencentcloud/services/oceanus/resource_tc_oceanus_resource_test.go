package oceanus_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudNeedFixOceanusResourceResource_basic -v
func TestAccTencentCloudNeedFixOceanusResourceResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOceanusResource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_oceanus_resource.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_oceanus_resource.example", "resource_type"),
					resource.TestCheckResourceAttrSet("tencentcloud_oceanus_resource.example", "remark"),
					resource.TestCheckResourceAttrSet("tencentcloud_oceanus_resource.example", "name"),
					resource.TestCheckResourceAttrSet("tencentcloud_oceanus_resource.example", "resource_config_remark"),
					resource.TestCheckResourceAttrSet("tencentcloud_oceanus_resource.example", "folder_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_oceanus_resource.example", "work_space_id"),
				),
			},
		},
	})
}

const testAccOceanusResource = `
resource "tencentcloud_oceanus_resource" "example" {
  resource_loc {
    storage_type = 1
    param {
      bucket = "keep-terraform-1257058945"
      path   = "OceanusResource/junit-4.13.2.jar"
      region = "ap-guangzhou"
    }
  }

  resource_type          = 1
  remark                 = "remark."
  name                   = "tf_example"
  resource_config_remark = "config remark."
  folder_id              = "folder-7ctl246z"
  work_space_id          = "space-2idq8wbr"
}
`
