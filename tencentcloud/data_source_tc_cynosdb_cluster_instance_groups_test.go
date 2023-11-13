package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCynosdbClusterInstanceGroupsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCynosdbClusterInstanceGroupsDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_cynosdb_cluster_instance_groups.cluster_instance_groups")),
			},
		},
	})
}

const testAccCynosdbClusterInstanceGroupsDataSource = `

data "tencentcloud_cynosdb_cluster_instance_groups" "cluster_instance_groups" {
  cluster_id = &lt;nil&gt;
  total_count = &lt;nil&gt;
  instance_grp_info_list {
		app_id = &lt;nil&gt;
		cluster_id = &lt;nil&gt;
		created_time = &lt;nil&gt;
		deleted_time = &lt;nil&gt;
		instance_grp_id = &lt;nil&gt;
		status = &lt;nil&gt;
		type = &lt;nil&gt;
		updated_time = &lt;nil&gt;
		vip = &lt;nil&gt;
		vport = &lt;nil&gt;
		wan_domain = &lt;nil&gt;
		wan_i_p = &lt;nil&gt;
		wan_port = &lt;nil&gt;
		wan_status = &lt;nil&gt;
		instance_set {
			uin = &lt;nil&gt;
			app_id = &lt;nil&gt;
			cluster_id = &lt;nil&gt;
			cluster_name = &lt;nil&gt;
			instance_id = &lt;nil&gt;
			instance_name = &lt;nil&gt;
			project_id = &lt;nil&gt;
			region = &lt;nil&gt;
			zone = &lt;nil&gt;
			status = &lt;nil&gt;
			status_desc = &lt;nil&gt;
			db_type = &lt;nil&gt;
			db_version = &lt;nil&gt;
			cpu = &lt;nil&gt;
			memory = &lt;nil&gt;
			storage = &lt;nil&gt;
			instance_type = &lt;nil&gt;
			instance_role = &lt;nil&gt;
			update_time = &lt;nil&gt;
			create_time = &lt;nil&gt;
			vpc_id = &lt;nil&gt;
			subnet_id = &lt;nil&gt;
			vip = &lt;nil&gt;
			vport = &lt;nil&gt;
			pay_mode = &lt;nil&gt;
			period_end_time = &lt;nil&gt;
			destroy_deadline_text = &lt;nil&gt;
			isolate_time = &lt;nil&gt;
			net_type = &lt;nil&gt;
			wan_domain = &lt;nil&gt;
			wan_i_p = &lt;nil&gt;
			wan_port = &lt;nil&gt;
			wan_status = &lt;nil&gt;
			destroy_time = &lt;nil&gt;
			cynos_version = &lt;nil&gt;
			processing_task = &lt;nil&gt;
			renew_flag = &lt;nil&gt;
			min_cpu = &lt;nil&gt;
			max_cpu = &lt;nil&gt;
			serverless_status = &lt;nil&gt;
			storage_id = &lt;nil&gt;
			storage_pay_mode = &lt;nil&gt;
			physical_zone = &lt;nil&gt;
			business_type = &lt;nil&gt;
			tasks {
				task_id = &lt;nil&gt;
				task_type = &lt;nil&gt;
				task_status = &lt;nil&gt;
				object_id = &lt;nil&gt;
				object_type = &lt;nil&gt;
			}
			is_freeze = &lt;nil&gt;
			resource_tags {
				tag_key = &lt;nil&gt;
				tag_value = &lt;nil&gt;
			}
		}

  }
}

`
