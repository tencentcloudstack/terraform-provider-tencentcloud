package wedata_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudWedataListColumnLineageDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{{
			Config: testAccWedataListColumnLineageDataSource,
			Check: resource.ComposeTestCheckFunc(
				tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_wedata_list_column_lineage.example"),
			),
		}},
	})
}

const testAccWedataListColumnLineageDataSource = `
data "tencentcloud_wedata_list_column_lineage" "example" {
  table_unique_id = "B_CRyO4-3rMvNFPH_7aTaw"
  direction       = "INPUT"
  column_name     = "a"
  platform        = "WEDATA"
}
`
