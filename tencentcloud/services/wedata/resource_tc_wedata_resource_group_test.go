package wedata_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudWedataResourceGroupResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccWedataResourceGroup,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_resource_group.example", "id"),
				),
			},
			{
				Config: testAccWedataResourceGroupUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_resource_group.example", "id"),
				),
			},
		},
	})
}

const testAccWedataResourceGroup = `
resource "tencentcloud_wedata_resource_group" "example" {
  name = "tf_example"
  type {
    resource_group_type = "Integration"
    integration {
      real_time_data_sync {
        specification = "i32c"
        number        = 1
      }

      offline_data_sync {
        specification = "integrated"
        number        = 2
      }
    }
  }

  auto_renew_enabled = false
  purchase_period    = 1
  vpc_id             = "vpc-ds5rpnxh"
  subnet             = "subnet-fz7rw5zq"
  resource_region    = "ap-beijing-fsi"
  description        = "description."
}
`

const testAccWedataResourceGroupUpdate = `
resource "tencentcloud_wedata_resource_group" "example" {
  name = "tf_example"
  type {
    resource_group_type = "Integration"
    integration {
      real_time_data_sync {
        specification = "i32c"
        number        = 1
      }

      offline_data_sync {
        specification = "integrated"
        number        = 2
      }
    }
  }

  auto_renew_enabled = true
  purchase_period    = 2
  vpc_id             = "vpc-ds5rpnxh"
  subnet             = "subnet-fz7rw5zq"
  resource_region    = "ap-beijing-fsi"
  description        = "description."
}
`
