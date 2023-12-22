package dayuv2_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudAntiddosOverviewDdosEventListDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PREPAY) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAntiddosOverviewDdosEventListDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_antiddos_overview_ddos_event_list.overview_ddos_event_list"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_antiddos_overview_ddos_event_list.overview_ddos_event_list", "event_list.#"),
				),
			},
		},
	})
}

const testAccAntiddosOverviewDdosEventListDataSource = `
data "tencentcloud_antiddos_overview_ddos_event_list" "overview_ddos_event_list" {
	start_time = "2023-11-20 00:00:00"
	end_time = "2023-11-21 00:00:00"
	attack_status = "end"
}
`
