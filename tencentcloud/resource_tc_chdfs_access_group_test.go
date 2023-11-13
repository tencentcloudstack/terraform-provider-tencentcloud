package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudChdfsAccessGroupResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccChdfsAccessGroup,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_chdfs_access_group.access_group", "id")),
			},
			{
				ResourceName:      "tencentcloud_chdfs_access_group.access_group",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccChdfsAccessGroup = `

resource "tencentcloud_chdfs_access_group" "access_group" {
  access_group_name = "testAccessGroup"
  vpc_type = 1
  vpc_id = "test-vpd-id"
  description = &lt;nil&gt;
}

`
