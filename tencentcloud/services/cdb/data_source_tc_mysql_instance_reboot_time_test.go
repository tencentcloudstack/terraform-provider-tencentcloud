package cdb_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudMysqlInstanceRebootTimeDataSource_basic -v
func TestAccTencentCloudMysqlInstanceRebootTimeDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMysqlInstanceRebootTimeDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_mysql_instance_reboot_time.instance_reboot_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_instance_reboot_time.instance_reboot_time", "id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_instance_reboot_time.instance_reboot_time", "items.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_instance_reboot_time.instance_reboot_time", "items.0.instance_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_instance_reboot_time.instance_reboot_time", "items.0.time_in_seconds"),
				),
			},
		},
	})
}

const testAccMysqlInstanceRebootTimeDataSourceVar = `
variable "instance_id" {
	default = "` + tcacctest.DefaultDbBrainInstanceId + `"
}
`

const testAccMysqlInstanceRebootTimeDataSource = testAccMysqlInstanceRebootTimeDataSourceVar + `

data "tencentcloud_mysql_instance_reboot_time" "instance_reboot_time" {
	instance_ids = [var.instance_id]
}

`
