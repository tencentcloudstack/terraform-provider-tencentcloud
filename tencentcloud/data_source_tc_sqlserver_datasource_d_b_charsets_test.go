package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudSqlserverDatasourceDBCharsetsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSqlserverDatasourceDBCharsetsDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_sqlserver_datasource_d_b_charsets.datasource_d_b_charsets")),
			},
		},
	})
}

const testAccSqlserverDatasourceDBCharsetsDataSource = `

data "tencentcloud_sqlserver_datasource_d_b_charsets" "datasource_d_b_charsets" {
  instance_id = "mssql-j8kv137v"
  }

`
