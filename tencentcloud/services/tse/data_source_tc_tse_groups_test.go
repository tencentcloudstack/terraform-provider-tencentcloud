package tse_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudTseGroupsDataSource_basic -v
func TestAccTencentCloudTseGroupsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTseGroupsDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_tse_groups.groups"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_groups.groups", "result.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_groups.groups", "result.0.gateway_group_list.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_groups.groups", "result.0.gateway_group_list.0.binding_strategy.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_groups.groups", "result.0.gateway_group_list.0.binding_strategy.0.config.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_groups.groups", "result.0.gateway_group_list.0.binding_strategy.0.config.0.enabled"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_groups.groups", "result.0.gateway_group_list.0.binding_strategy.0.cron_config.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_groups.groups", "result.0.gateway_group_list.0.binding_strategy.0.cron_config.0.enabled"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_groups.groups", "result.0.gateway_group_list.0.create_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_groups.groups", "result.0.gateway_group_list.0.gateway_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_groups.groups", "result.0.gateway_group_list.0.group_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_groups.groups", "result.0.gateway_group_list.0.internet_max_bandwidth_out"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_groups.groups", "result.0.gateway_group_list.0.is_first_group"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_groups.groups", "result.0.gateway_group_list.0.modify_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_groups.groups", "result.0.gateway_group_list.0.name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_groups.groups", "result.0.gateway_group_list.0.node_config.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_groups.groups", "result.0.gateway_group_list.0.node_config.0.number"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_groups.groups", "result.0.gateway_group_list.0.node_config.0.specification"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_groups.groups", "result.0.gateway_group_list.0.status"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_groups.groups", "result.0.gateway_group_list.0.subnet_ids"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_groups.groups", "result.0.total_count"),
				),
			},
		},
	})
}

const testAccTseGroupsDataSource = tcacctest.DefaultTseVar + `

data "tencentcloud_tse_groups" "groups" {
  gateway_id = var.gateway_id
  filters {
    name   = "GroupId"
    values = ["group-013c0d8e"]
  }
}

`
