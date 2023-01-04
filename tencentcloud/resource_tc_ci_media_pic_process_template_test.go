package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

func TestAccTencentCloudCiMediaPicProcessTemplateResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCiMediaPicProcessTemplate,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_ci_media_pic_process_template.media_pic_process_template", "id")),
			},
			{
				ResourceName:      "tencentcloud_ci_media_pic_process_template.media_pic_process_template",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCiMediaPicProcessTemplate = `

resource "tencentcloud_ci_media_pic_process_template" "media_pic_process_template" {
  pic_process {
		is_pic_info = &lt;nil&gt;
		process_rule = &lt;nil&gt;

  }
}

`
