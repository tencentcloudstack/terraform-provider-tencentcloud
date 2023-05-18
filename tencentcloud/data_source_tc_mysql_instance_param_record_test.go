package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudMysqlInstanceParamRecordDataSource_basic -v
func TestAccTencentCloudMysqlInstanceParamRecordDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMysqlInstanceParamRecordDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_mysql_instance_param_record.instance_param_record"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_instance_param_record.instance_param_record", "id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_instance_param_record.instance_param_record", "items.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_instance_param_record.instance_param_record", "items.0.is_success"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_instance_param_record.instance_param_record", "items.0.modify_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_instance_param_record.instance_param_record", "items.0.new_value"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_instance_param_record.instance_param_record", "items.0.old_value"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_instance_param_record.instance_param_record", "items.0.param_name"),
				),
			},
		},
	})
}

const testAccMysqlInstanceParamRecordDataSourceVar = `
variable "instance_id" {
	default = "` + defaultDbBrainInstanceId + `"
}
`

const testAccMysqlInstanceParamRecordDataSource = testAccMysqlInstanceParamRecordDataSourceVar + `

data "tencentcloud_mysql_instance_param_record" "instance_param_record" {
  instance_id = var.instance_id
  }

`
