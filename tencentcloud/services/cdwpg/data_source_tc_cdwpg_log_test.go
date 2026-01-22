package cdwpg_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudCdwpgLogDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{{
			Config: testAccCdwpgLogDataSource,
			Check:  resource.ComposeTestCheckFunc(tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_cdwpg_log.cdwpg_log")),
		}},
	})
}

const testAccCdwpgLogDataSource = `
data "tencentcloud_cdwpg_log" "cdwpg_log" {
	instance_id = "cdwpg-gexy9tue"
	start_time = "2025-03-21 00:00:00"
	end_time = "2025-03-21 23:59:59"
}
`
