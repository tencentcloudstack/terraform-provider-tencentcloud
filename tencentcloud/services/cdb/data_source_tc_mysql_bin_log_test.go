package cdb_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudMysqlBinLogDataSource_basic -v
func TestAccTencentCloudMysqlBinLogDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMysqlBinLogDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_mysql_bin_log.bin_log"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_bin_log.bin_log", "items.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_bin_log.bin_log", "items.0.name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_bin_log.bin_log", "items.0.region"),
					// resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_bin_log.bin_log", "items.0.remote_info"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_bin_log.bin_log", "items.0.size"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_bin_log.bin_log", "items.0.status"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_bin_log.bin_log", "items.0.type"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_bin_log.bin_log", "items.0.intranet_url"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_bin_log.bin_log", "items.0.internet_url"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_bin_log.bin_log", "items.0.date"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_bin_log.bin_log", "items.0.cos_storage_type"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_bin_log.bin_log", "items.0.binlog_start_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_bin_log.bin_log", "items.0.binlog_finish_time"),
				),
			},
		},
	})
}

const testAccMysqlBinLogDataSourceVar = `
variable "instance_id" {
  default = "` + tcacctest.DefaultDbBrainInstanceId + `"
}
`
const testAccMysqlBinLogDataSource = testAccMysqlBinLogDataSourceVar + `

data "tencentcloud_mysql_bin_log" "bin_log" {
  instance_id = var.instance_id
}

`
