package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var testDataSqlserverReadonlyGroupsName = "data.tencentcloud_sqlserver_readonly_groups.test"

func TestAccDataSourceTencentCloudSqlserverReadonlyGroups(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
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

const testAccTencentCloudDataSqlserverReadonlyGroupsBasic = CommonPresetSQLServer + `
data "tencentcloud_sqlserver_readonly_groups" "test"{
	master_instance_id = local.sqlserver_id
}
`
