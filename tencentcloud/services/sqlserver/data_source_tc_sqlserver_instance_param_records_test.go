package sqlserver_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudSqlserverInstanceParamRecordsDataSource_basic -v
func TestAccTencentCloudSqlserverInstanceParamRecordsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSqlserverInstanceParamRecordsDataSource,
				Check:  resource.ComposeTestCheckFunc(tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_sqlserver_instance_param_records.example")),
			},
		},
	})
}

const testAccSqlserverInstanceParamRecordsDataSource = `
data "tencentcloud_sqlserver_instance_param_records" "example" {
  instance_id = "mssql-qelbzgwf"
}
`
