package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudMysqlSlowLogDataSource_basic -v
func TestAccTencentCloudMysqlSlowLogDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMysqlSlowLogDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_mysql_slow_log.slow_log"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_slow_log.slow_log", "id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_slow_log.slow_log", "items.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_slow_log.slow_log", "items.0.date"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_slow_log.slow_log", "items.0.internet_url"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_slow_log.slow_log", "items.0.intranet_url"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_slow_log.slow_log", "items.0.name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_slow_log.slow_log", "items.0.size"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_slow_log.slow_log", "items.0.type"),
				),
			},
		},
	})
}

const testAccMysqlSlowLogDataSourceVar = `
variable "instance_id" {
	default = "` + defaultDbBrainInstanceId + `"
}
`

const testAccMysqlSlowLogDataSource = testAccMysqlSlowLogDataSourceVar + `

data "tencentcloud_mysql_slow_log" "slow_log" {
  instance_id = var.instance_id
}

`
