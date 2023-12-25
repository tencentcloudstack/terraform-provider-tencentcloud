package sqlserver_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudSqlserverInsAttributeDataSource_basic -v
func TestAccTencentCloudSqlserverInsAttributeDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSqlserverDatasourceInsAttributeDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_sqlserver_ins_attribute.example"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_sqlserver_ins_attribute.example", "instance_id"),
				),
			},
		},
	})
}

const testAccSqlserverDatasourceInsAttributeDataSource = `
data "tencentcloud_sqlserver_ins_attribute" "example" {
  instance_id = "mssql-gyg9xycl"
}
`
