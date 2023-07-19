package tencentcloud

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudPostgresqlXlogsDataSource(t *testing.T) {
	t.Parallel()

	startTime := time.Now().AddDate(0, 0, -7).Format("2006-01-02 15:04:05")
	endTime := time.Now().Format("2006-01-02 15:04:05")
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccStepSetRegion(t, "ap-guangzhou")
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourcePostgresqlXlogsBasic(startTime, endTime),
				PreConfig: func() {
					testAccStepSetRegion(t, "ap-guangzhou")
					testAccPreCheckCommon(t, ACCOUNT_TYPE_COMMON)
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.tencentcloud_postgresql_xlogs.foo", "start_time", startTime),
					resource.TestCheckResourceAttr("data.tencentcloud_postgresql_xlogs.foo", "end_time", endTime),
					resource.TestCheckResourceAttrSet("data.tencentcloud_postgresql_xlogs.foo", "list.#"),
				),
			},
		},
	})
}

func testAccDataSourcePostgresqlXlogsBasic(startTime, endTime string) string {
	return fmt.Sprintf(`
%s
data "tencentcloud_postgresql_xlogs" "foo" {
	instance_id = local.pgsql_id
	start_time = "%s"
	end_time = "%s"
}

data "tencentcloud_postgresql_xlogs" "bar" {
	instance_id = local.pgsql_id
}	
`, CommonPresetPGSQL, startTime, endTime)
}
