package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCamServiceLinkedRoleResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCamServiceLinkedRole,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_cam_service_linked_role.service_linked_role", "id")),
			},
			{
				ResourceName:      "tencentcloud_cam_service_linked_role.service_linked_role",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCamServiceLinkedRole = `

resource "tencentcloud_cam_service_linked_role" "service_linked_role" {
  qcs_service_name = "recordlistdnspod.cdn.cloud.tencent.com"
  custom_suffix = ""
  description = "desc cam"
  tags = {
    "createdBy" = "terraform"
  }
}

`
