package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testAccPostgresqlReadonlyGroupsObject = "data.tencentcloud_postgresql_readonly_groups.read_only_groups"

func TestAccTencentCloudPostgresqlReadonlyGroupsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccStepSetRegion(t, "ap-guangzhou")
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPostgresqlReadonlyGroupsDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID(testAccPostgresqlReadonlyGroupsObject),
					resource.TestCheckResourceAttr(testAccPostgresqlReadonlyGroupsObject, "order_by", "CreateTime"),
					resource.TestCheckResourceAttr(testAccPostgresqlReadonlyGroupsObject, "order_by_type", "asc"),
					resource.TestCheckResourceAttrSet(testAccPostgresqlReadonlyGroupsObject, "filters.#"),
					resource.TestCheckResourceAttr(testAccPostgresqlReadonlyGroupsObject, "filters.0.name", "db-master-instance-id"),
					resource.TestCheckResourceAttrSet(testAccPostgresqlReadonlyGroupsObject, "filters.0.values.#"),
					resource.TestCheckResourceAttrSet(testAccPostgresqlReadonlyGroupsObject, "read_only_group_list.#"),
					resource.TestCheckResourceAttrSet(testAccPostgresqlReadonlyGroupsObject, "read_only_group_list.0.read_only_group_id"),
					resource.TestCheckResourceAttrSet(testAccPostgresqlReadonlyGroupsObject, "read_only_group_list.0.read_only_group_name"),
					resource.TestCheckResourceAttrSet(testAccPostgresqlReadonlyGroupsObject, "read_only_group_list.0.project_id"),
					resource.TestCheckResourceAttrSet(testAccPostgresqlReadonlyGroupsObject, "read_only_group_list.0.master_db_instance_id"),
					resource.TestCheckResourceAttrSet(testAccPostgresqlReadonlyGroupsObject, "read_only_group_list.0.min_delay_eliminate_reserve"),
					resource.TestCheckResourceAttrSet(testAccPostgresqlReadonlyGroupsObject, "read_only_group_list.0.replay_latency_eliminate"),
					resource.TestCheckResourceAttrSet(testAccPostgresqlReadonlyGroupsObject, "read_only_group_list.0.max_replay_lag"),
					resource.TestCheckResourceAttrSet(testAccPostgresqlReadonlyGroupsObject, "read_only_group_list.0.region"),
					resource.TestCheckResourceAttrSet(testAccPostgresqlReadonlyGroupsObject, "read_only_group_list.0.zone"),
					resource.TestCheckResourceAttrSet(testAccPostgresqlReadonlyGroupsObject, "read_only_group_list.0.status"),
					resource.TestCheckResourceAttrSet(testAccPostgresqlReadonlyGroupsObject, "read_only_group_list.0.read_only_db_instance_list.#"),
				),
			},
		},
	})
}

const testAccPostgresqlReadonlyGroupsDataSource = CommonPresetPGSQL + `

resource "tencentcloud_postgresql_readonly_group" "group" {
	master_db_instance_id = local.pgsql_id
	name = "test-datasource"
	project_id = 0
	vpc_id = "vpc-86v957zb"
	subnet_id = "subnet-enm92y0m"
	replay_lag_eliminate = 1
	replay_latency_eliminate =  1
	max_replay_lag = 100
	max_replay_latency = 512
	min_delay_eliminate_reserve = 1
}

data "tencentcloud_postgresql_readonly_groups" "read_only_groups" {
  filters {
	name = "db-master-instance-id"
	values = [tencentcloud_postgresql_readonly_group.group.master_db_instance_id]
  }
  order_by = "CreateTime"
  order_by_type = "asc"
}

`
