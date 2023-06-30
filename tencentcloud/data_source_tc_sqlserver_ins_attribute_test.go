package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudSqlserverInsAttributeDataSource_basic -v
func TestAccTencentCloudSqlserverInsAttributeDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSqlserverDatasourceInsAttributeDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_sqlserver_ins_attribute.ins_attribute"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_sqlserver_ins_attribute.ins_attribute", "instance_id"),
				),
			},
		},
	})
}

const testAccSqlserverDatasourceInsAttributeDataSource = `
data "tencentcloud_sqlserver_ins_attribute" "ins_attribute" {
  instance_id = "mssql-gyg9xycl"
}
`
