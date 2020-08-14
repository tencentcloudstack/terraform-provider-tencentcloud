package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

var testDataSqlserverReadonlyGroupsName = "data.tencentcloud_sqlserver_readonly_groups.test"

func TestAccTencentCloudDataSqlserverReadonlyGroups(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLBDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudDataSqlserverReadonlyGroupsBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testDataSqlserverReadonlyGroupsName, "list.#", "2"),
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

const testAccTencentCloudDataSqlserverReadonlyGroupsBasic = `
data "tencentcloud_sqlserver_readonly_groups" "test"{
	master_instance_id = "mssql-ixq78we9"
}
`
