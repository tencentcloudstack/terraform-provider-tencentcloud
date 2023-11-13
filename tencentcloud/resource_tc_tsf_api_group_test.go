package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudTsfApiGroupResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTsfApiGroup,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_tsf_api_group.api_group", "id")),
			},
			{
				ResourceName:      "tencentcloud_tsf_api_group.api_group",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTsfApiGroup = `

resource "tencentcloud_tsf_api_group" "api_group" {
  group_name = ""
  group_context = ""
  auth_type = ""
  description = ""
  group_type = ""
  gateway_instance_id = ""
  namespace_name_key = ""
  service_name_key = ""
  namespace_name_key_position = ""
  service_name_key_position = ""
  }

`
