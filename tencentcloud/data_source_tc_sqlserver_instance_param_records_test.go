package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudSqlserverInstanceParamRecordsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSqlserverInstanceParamRecordsDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_sqlserver_instance_param_records.instance_param_records")),
			},
		},
	})
}

const testAccSqlserverInstanceParamRecordsDataSource = `

data "tencentcloud_sqlserver_instance_param_records" "instance_param_records" {
  instance_id = "mssql-j8kv137v"
  }

`
