package wedata_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudWedataSqlFolderResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{{
			Config: testAccWedataSqlFolder,
			Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_wedata_sql_folder.wedata_sql_folder", "id")),
		}, {
			ResourceName:      "tencentcloud_wedata_sql_folder.wedata_sql_folder",
			ImportState:       true,
			ImportStateVerify: true,
		}},
	})
}

const testAccWedataSqlFolder = `

resource "tencentcloud_wedata_sql_folder" "wedata_sql_folder" {
}
`
