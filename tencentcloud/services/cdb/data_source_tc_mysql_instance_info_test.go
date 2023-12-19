package cdb_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudMysqlInstanceInfoDataSource_basic -v
func TestAccTencentCloudMysqlInstanceInfoDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMysqlInstanceInfoDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_mysql_instance_info.instance_info"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_instance_info.instance_info", "id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_instance_info.instance_info", "instance_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_instance_info.instance_info", "instance_name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_instance_info.instance_info", "encryption"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_instance_info.instance_info", "default_kms_region"),
				),
			},
		},
	})
}

const testAccMysqlInstanceInfoDataSourceVar = `
variable "instance_id" {
	default = "` + tcacctest.DefaultDbBrainInstanceId + `"
}
`

const testAccMysqlInstanceInfoDataSource = testAccMysqlInstanceInfoDataSourceVar + `

data "tencentcloud_mysql_instance_info" "instance_info" {
	instance_id = var.instance_id
}

`
