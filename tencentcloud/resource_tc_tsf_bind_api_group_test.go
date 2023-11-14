package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudTsfBindApiGroupResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTsfBindApiGroup,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_tsf_bind_api_group.bind_api_group", "id")),
			},
			{
				ResourceName:      "tencentcloud_tsf_bind_api_group.bind_api_group",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTsfBindApiGroup = `

resource "tencentcloud_tsf_bind_api_group" "bind_api_group" {
  group_gateway_list {
		gateway_deploy_group_id = "group-vzd97zpy"
		group_id = "grp-qp0rj3zi"

  }
}

`
