package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudNeedFixOceanusWorkSpaceResource_basic -v
func TestAccTencentCloudNeedFixOceanusWorkSpaceResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOceanusWorkSpace,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_oceanus_work_space.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_oceanus_work_space.example", "work_space_name"),
					resource.TestCheckResourceAttrSet("tencentcloud_oceanus_work_space.example", "description"),
				),
			},
			{
				ResourceName:      "tencentcloud_oceanus_work_space.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccOceanusWorkSpaceUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_oceanus_work_space.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_oceanus_work_space.example", "work_space_name"),
					resource.TestCheckResourceAttrSet("tencentcloud_oceanus_work_space.example", "description"),
				),
			},
		},
	})
}

const testAccOceanusWorkSpace = `
resource "tencentcloud_oceanus_work_space" "example" {
  work_space_name = "tf_example"
  description     = "example description."
}
`

const testAccOceanusWorkSpaceUpdate = `
resource "tencentcloud_oceanus_work_space" "example" {
  work_space_name = "tf_example_update"
  description     = "example description update."
}
`
