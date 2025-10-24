package wedata_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudWedataListProcessLineageDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{{
			Config: testAccWedataListProcessLineageDataSource,
			Check: resource.ComposeTestCheckFunc(
				tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_wedata_list_process_lineage.example"),
			),
		}},
	})
}

const testAccWedataListProcessLineageDataSource = `
data "tencentcloud_wedata_list_process_lineage" "example" {
  process_id   = "20241107221758402"
  process_type = "SCHEDULE_TASK"
  platform     = "WEDATA"
}
`
