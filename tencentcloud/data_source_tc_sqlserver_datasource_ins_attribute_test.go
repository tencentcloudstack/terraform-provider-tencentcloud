package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudSqlserverDatasourceInsAttributeDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSqlserverDatasourceInsAttributeDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_sqlserver_datasource_ins_attribute.datasource_ins_attribute")),
			},
		},
	})
}

const testAccSqlserverDatasourceInsAttributeDataSource = `

data "tencentcloud_sqlserver_datasource_ins_attribute" "datasource_ins_attribute" {
                  }

`
