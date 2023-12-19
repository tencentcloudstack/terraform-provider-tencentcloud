package cdb_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudMysqlCloneListDataSource_basic -v
func TestAccTencentCloudMysqlCloneListDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMysqlCloneListDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_mysql_clone_list.clone_list"),
				),
			},
		},
	})
}

const testAccMysqlCloneListDataSourceVar = `
variable "instance_id" {
  default = "` + tcacctest.DefaultDbBrainInstanceId + `"
}
`

const testAccMysqlCloneListDataSource = testAccMysqlCloneListDataSourceVar + `

data "tencentcloud_mysql_clone_list" "clone_list" {
  instance_id = var.instance_id
}

`
