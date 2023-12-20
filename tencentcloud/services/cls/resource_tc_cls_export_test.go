package cls_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudClsExportResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
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
  topic_id  = "7e34a3a7-635e-4da8-9005-88106c1fde69"
  log_count = 2
  query     = "select count(*) as count"
  from      = 1607499107000
  to        = 1607499108000
  order     = "desc"
  format    = "json"
}


`
