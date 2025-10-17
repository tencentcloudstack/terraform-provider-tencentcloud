package wedata_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudWedataResourceGroupToProjectAttachmentResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccWedataResourceGroupToProjectAttachment,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_resource_group_to_project_attachment.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_resource_group_to_project_attachment.example", "resource_group_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_resource_group_to_project_attachment.example", "project_id"),
				),
			},
			{
				ResourceName:      "tencentcloud_wedata_resource_group_to_project_attachment.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccWedataResourceGroupToProjectAttachment = `
resource "tencentcloud_wedata_resource_group_to_project_attachment" "example" {
  resource_group_id  = "20250909161820129828"
  project_id         = "2983848457986924544"
}
`
