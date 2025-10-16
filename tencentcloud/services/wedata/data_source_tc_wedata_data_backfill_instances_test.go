package wedata_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudWedataDataBackfillInstancesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccWedataDataBackfillInstancesDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_wedata_data_backfill_instances.wedata_data_backfill_instances"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_wedata_data_backfill_instances.wedata_data_backfill_instances", "id"),
					resource.TestCheckResourceAttr("data.tencentcloud_wedata_data_backfill_instances.wedata_data_backfill_instances", "data.#", "1"),
				),
			},
		},
	})
}

const testAccWedataDataBackfillInstancesDataSource = `

data "tencentcloud_wedata_data_backfill_instances" "wedata_data_backfill_instances" {
  project_id            = "1859317240494305280"
  data_backfill_plan_id = "deb71ea1-f708-47ab-8eb6-491ce5b9c011"
  task_id               = "20231011152006462"
}
`
