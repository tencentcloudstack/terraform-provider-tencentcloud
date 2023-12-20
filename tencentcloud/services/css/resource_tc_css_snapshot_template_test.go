package css_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -test.run TestAccTencentCloudCssSnapshotTemplateResource_basic -v
func TestAccTencentCloudCssSnapshotTemplateResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCssSnapshotTemplate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_css_snapshot_template.snapshot_template", "id"),
					resource.TestCheckResourceAttr("tencentcloud_css_snapshot_template.snapshot_template", "cos_app_id", "1308919341"),
					resource.TestCheckResourceAttr("tencentcloud_css_snapshot_template.snapshot_template", "cos_bucket", "keep-bucket"),
					resource.TestCheckResourceAttr("tencentcloud_css_snapshot_template.snapshot_template", "cos_region", "ap-guangzhou"),
					resource.TestCheckResourceAttr("tencentcloud_css_snapshot_template.snapshot_template", "description", "snapshot template"),
					resource.TestCheckResourceAttr("tencentcloud_css_snapshot_template.snapshot_template", "height", "0"),
					resource.TestCheckResourceAttr("tencentcloud_css_snapshot_template.snapshot_template", "porn_flag", "0"),
					resource.TestCheckResourceAttr("tencentcloud_css_snapshot_template.snapshot_template", "width", "0"),
					resource.TestCheckResourceAttr("tencentcloud_css_snapshot_template.snapshot_template", "snapshot_interval", "2"),
					resource.TestCheckResourceAttr("tencentcloud_css_snapshot_template.snapshot_template", "template_name", "tf-snapshot-template"),
					resource.TestCheckResourceAttr("tencentcloud_css_snapshot_template.snapshot_template", "cos_prefix", "/test1-test2/"),
					resource.TestCheckResourceAttr("tencentcloud_css_snapshot_template.snapshot_template", "cos_file_name", "temp"),
				),
			},
			{
				ResourceName:      "tencentcloud_css_snapshot_template.snapshot_template",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccCssSnapshotTemplateUp,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_css_snapshot_template.snapshot_template", "id"),
					resource.TestCheckResourceAttr("tencentcloud_css_snapshot_template.snapshot_template", "cos_app_id", "1308919341"),
					resource.TestCheckResourceAttr("tencentcloud_css_snapshot_template.snapshot_template", "cos_bucket", "keep-bucket"),
					resource.TestCheckResourceAttr("tencentcloud_css_snapshot_template.snapshot_template", "cos_region", "ap-guangzhou"),
					resource.TestCheckResourceAttr("tencentcloud_css_snapshot_template.snapshot_template", "description", "snapshot template1"),
					resource.TestCheckResourceAttr("tencentcloud_css_snapshot_template.snapshot_template", "height", "2"),
					resource.TestCheckResourceAttr("tencentcloud_css_snapshot_template.snapshot_template", "porn_flag", "0"),
					resource.TestCheckResourceAttr("tencentcloud_css_snapshot_template.snapshot_template", "width", "2"),
					resource.TestCheckResourceAttr("tencentcloud_css_snapshot_template.snapshot_template", "snapshot_interval", "4"),
					resource.TestCheckResourceAttr("tencentcloud_css_snapshot_template.snapshot_template", "template_name", "tf-snapshot-template1"),
					resource.TestCheckResourceAttr("tencentcloud_css_snapshot_template.snapshot_template", "cos_prefix", "/test1-test2/"),
					resource.TestCheckResourceAttr("tencentcloud_css_snapshot_template.snapshot_template", "cos_file_name", "temp1"),
				),
			},
		},
	})
}

const testAccCssSnapshotTemplate = `

resource "tencentcloud_css_snapshot_template" "snapshot_template" {
  cos_app_id        = 1308919341
  cos_bucket        = "keep-bucket"
  cos_region        = "ap-guangzhou"
  description       = "snapshot template"
  height            = 0
  porn_flag         = 0
  snapshot_interval = 2
  template_name     = "tf-snapshot-template"
  width             = 0
  cos_prefix        = "/test1-test2/"
  cos_file_name     = "temp"
}

`

const testAccCssSnapshotTemplateUp = `

resource "tencentcloud_css_snapshot_template" "snapshot_template" {
  cos_app_id        = 1308919341
  cos_bucket        = "keep-bucket"
  cos_region        = "ap-guangzhou"
  description       = "snapshot template1"
  height            = 2
  porn_flag         = 0
  snapshot_interval = 4
  template_name     = "tf-snapshot-template1"
  width             = 2
  cos_prefix        = "/test1-test2/"
  cos_file_name     = "temp1"
}

`
