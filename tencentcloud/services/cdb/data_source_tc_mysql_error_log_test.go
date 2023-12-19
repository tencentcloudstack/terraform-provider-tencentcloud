package cdb_test

import (
	"fmt"
	"testing"
	"time"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudMysqlErrorLogDataSource_basic -v
func TestAccTencentCloudMysqlErrorLogDataSource_basic(t *testing.T) {
	t.Parallel()

	startTime := time.Now().AddDate(0, 0, -29).Unix()
	endTime := time.Now().Unix()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccMysqlErrorLogDataSource, startTime, endTime),
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_mysql_error_log.error_log"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_error_log.error_log", "items.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_error_log.error_log", "items.0.content"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_error_log.error_log", "items.0.timestamp"),
				),
			},
		},
	})
}

const testAccMysqlErrorLogDataSourceVar = `
variable "instance_id" {
  default = "` + tcacctest.DefaultDbBrainInstanceId + `"
}
`

const testAccMysqlErrorLogDataSource = testAccMysqlErrorLogDataSourceVar + `

data "tencentcloud_mysql_error_log" "error_log" {
  instance_id = var.instance_id
  start_time = %v
  end_time = %v
}

`
