package sqlserver_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var testDataSqlserverReadonlyGroupsName = "data.tencentcloud_sqlserver_readonly_groups.test"

// go test -i; go test -test.run TestAccDataSourceTencentCloudSqlserverReadonlyGroups -v
func TestAccDataSourceTencentCloudSqlserverReadonlyGroups(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckLBDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudDataSqlserverReadonlyGroupsBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(testDataSqlserverReadonlyGroupsName, "list.#"),
					resource.TestCheckResourceAttrSet(testDataSqlserverReadonlyGroupsName, "list.0.vport"),
					resource.TestCheckResourceAttrSet(testDataSqlserverReadonlyGroupsName, "list.0.vip"),
					resource.TestCheckResourceAttrSet(testDataSqlserverReadonlyGroupsName, "list.0.min_instances"),
					resource.TestCheckResourceAttrSet(testDataSqlserverReadonlyGroupsName, "list.0.is_offline_delay"),
					resource.TestCheckResourceAttrSet(testDataSqlserverReadonlyGroupsName, "list.0.max_delay_time"),
					resource.TestCheckResourceAttrSet(testDataSqlserverReadonlyGroupsName, "list.0.name"),
					resource.TestCheckResourceAttrSet(testDataSqlserverReadonlyGroupsName, "list.0.id"),
					resource.TestCheckResourceAttrSet(testDataSqlserverReadonlyGroupsName, "list.0.master_instance_id"),
					resource.TestCheckResourceAttrSet(testDataSqlserverReadonlyGroupsName, "list.0.status"),
					resource.TestCheckResourceAttrSet(testDataSqlserverReadonlyGroupsName, "list.0.readonly_instance_set.0"),
				),
			},
		},
	})
}

const testAccTencentCloudDataSqlserverReadonlyGroupsBasic = tcacctest.CommonPresetSQLServer + `
data "tencentcloud_sqlserver_readonly_groups" "test"{
  master_instance_id = local.sqlserver_id
}
`
