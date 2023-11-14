package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudRumProjectStatusAttachmentResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccRumProjectStatusAttachment,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_rum_project_status_attachment.project_status_attachment", "id")),
			},
			{
				ResourceName:      "tencentcloud_rum_project_status_attachment.project_status_attachment",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccRumProjectStatusAttachment = `

resource "tencentcloud_rum_project_status_attachment" "project_status_attachment" {
  project_id = 101
}

`
