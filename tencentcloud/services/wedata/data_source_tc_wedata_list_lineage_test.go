package wedata_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudWedataListLineageDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{{
			Config: testAccWedataListLineageDataSource,
			Check: resource.ComposeTestCheckFunc(
				tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_wedata_list_lineage.example"),
			),
		}},
	})
}

const testAccWedataListLineageDataSource = `
data "tencentcloud_wedata_list_lineage" "example" {
  resource_unique_id = "fM8OgzE-AM2h4aaJmdXoPg"
  resource_type      = "TABLE"
  direction          = "INPUT"
  platform           = "WEDATA"
}
`
