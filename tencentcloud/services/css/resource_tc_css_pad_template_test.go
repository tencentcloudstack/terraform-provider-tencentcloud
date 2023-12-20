package css_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -test.run TestAccTencentCloudCssPadTemplateResource_basic -v
func TestAccTencentCloudCssPadTemplateResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCssPadTemplate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_css_pad_template.pad_template", "id"),
					resource.TestCheckResourceAttr("tencentcloud_css_pad_template.pad_template", "description", "pad template"),
					resource.TestCheckResourceAttr("tencentcloud_css_pad_template.pad_template", "max_duration", "120000"),
					resource.TestCheckResourceAttr("tencentcloud_css_pad_template.pad_template", "template_name", "tf-pad"),
					resource.TestCheckResourceAttr("tencentcloud_css_pad_template.pad_template", "type", "1"),
					resource.TestCheckResourceAttr("tencentcloud_css_pad_template.pad_template", "url", "https://livewatermark-1251132611.cos.ap-guangzhou.myqcloud.com/1308919341/watermark_img_1698736540399_1441698123618_.pic.jpg"),
					resource.TestCheckResourceAttr("tencentcloud_css_pad_template.pad_template", "wait_duration", "2000"),
				),
			},
			{
				ResourceName:      "tencentcloud_css_pad_template.pad_template",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccCssPadTemplateUp,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_css_pad_template.pad_template", "id"),
					resource.TestCheckResourceAttr("tencentcloud_css_pad_template.pad_template", "description", "pad template1"),
					resource.TestCheckResourceAttr("tencentcloud_css_pad_template.pad_template", "max_duration", "120001"),
					resource.TestCheckResourceAttr("tencentcloud_css_pad_template.pad_template", "template_name", "tf-pad1"),
					resource.TestCheckResourceAttr("tencentcloud_css_pad_template.pad_template", "type", "1"),
					resource.TestCheckResourceAttr("tencentcloud_css_pad_template.pad_template", "url", "https://livewatermark-1251132611.cos.ap-guangzhou.myqcloud.com/1308919341/watermark_img_1698736540399_1441698123618_.pic.jpg"),
					resource.TestCheckResourceAttr("tencentcloud_css_pad_template.pad_template", "wait_duration", "2001"),
				),
			},
		},
	})
}

const testAccCssPadTemplate = `

resource "tencentcloud_css_pad_template" "pad_template" {
  description   = "pad template"
  max_duration  = 120000
  template_name = "tf-pad"
  type          = 1
  url           = "https://livewatermark-1251132611.cos.ap-guangzhou.myqcloud.com/1308919341/watermark_img_1698736540399_1441698123618_.pic.jpg"
  wait_duration = 2000
}

`

const testAccCssPadTemplateUp = `

resource "tencentcloud_css_pad_template" "pad_template" {
  description   = "pad template1"
  max_duration  = 120001
  template_name = "tf-pad1"
  type          = 1
  url           = "https://livewatermark-1251132611.cos.ap-guangzhou.myqcloud.com/1308919341/watermark_img_1698736540399_1441698123618_.pic.jpg"
  wait_duration = 2001
}

`
