package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -test.run TestAccTencentCloudRumProjectStatusAttachmentResource_basic -v
func TestAccTencentCloudRumProjectStatusAttachmentResource_basic(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccRumProjectStatusAttachment,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_rum_project_status_attachment.project_status_attachment", "id"),
				),
			},
			{
				Config: testAccRumProjectStatusAttachmentUp,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_rum_project_status_attachment.project_status_attachment", "id"),
				),
			},
			{
				ResourceName:            "tencentcloud_rum_project_status_attachment.project_status_attachment",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"operate"},
			},
		},
	})
}

const testAccRumProjectStatusAttachmentVar = `
variable "project_id" {
  default = "` + defaultRumProjectId + `"
}
`

const testAccRumProjectStatusAttachment = testAccRumProjectStatusAttachmentVar + `

resource "tencentcloud_rum_project_status_attachment" "project_status_attachment" {
	project_id = var.project_id
	operate = "stop"
}

`

const testAccRumProjectStatusAttachmentUp = testAccRumProjectStatusAttachmentVar + `


resource "tencentcloud_rum_project_status_attachment" "project_status_attachment" {
	project_id = var.project_id
	operate = "resume"
}

`
