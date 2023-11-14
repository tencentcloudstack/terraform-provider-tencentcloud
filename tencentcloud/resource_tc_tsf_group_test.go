package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudTsfGroupResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTsfGroup,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_tsf_group.group", "id")),
			},
			{
				ResourceName:      "tencentcloud_tsf_group.group",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTsfGroup = `

resource "tencentcloud_tsf_group" "group" {
  application_id = ""
  namespace_id = ""
  group_name = ""
  cluster_id = ""
  group_desc = ""
  group_resource_type = ""
  alias = ""
  tags = {
    "createdBy" = "terraform"
  }
}

`
