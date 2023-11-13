package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudClsExportResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccClsExport,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_cls_export.export", "id")),
			},
			{
				ResourceName:      "tencentcloud_cls_export.export",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccClsExport = `

resource "tencentcloud_cls_export" "export" {
  topic_id = "5cd3a17e-fb0b-418c-afd7-77b365397426"
  query = "* | select count(*) as count"
  from = 1607499107000
  order = "desc"
  format = "json"
}

`
