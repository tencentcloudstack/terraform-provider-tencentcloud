package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudMysqlSwitchRecordDataSource_basic -v
func TestAccTencentCloudMysqlSwitchRecordDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMysqlSwitchRecordDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_mysql_switch_record.switch_record"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_switch_record.switch_record", "id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_switch_record.switch_record", "items.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_switch_record.switch_record", "items.0.switch_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_switch_record.switch_record", "items.0.switch_type"),
				),
			},
		},
	})
}

const testAccMysqlSwitchRecordDataSourceVar = `
variable "instance_id" {
	default = "` + defaultDbBrainInstanceId + `"
}
`

const testAccMysqlSwitchRecordDataSource = testAccMysqlSwitchRecordDataSourceVar + `

data "tencentcloud_mysql_switch_record" "switch_record" {
  instance_id = var.instance_id
  }

`
