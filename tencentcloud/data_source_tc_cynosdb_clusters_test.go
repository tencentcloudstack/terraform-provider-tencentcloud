package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCynosdbClustersDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCynosdbClustersDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_cynosdb_clusters.clusters")),
			},
		},
	})
}

const testAccCynosdbClustersDataSource = `

data "tencentcloud_cynosdb_clusters" "clusters" {
  db_type = "MYSQL"
  limit = 20
  offset = 0
  order_by = &lt;nil&gt;
  order_by_type = &lt;nil&gt;
  filters {
		names = &lt;nil&gt;
		values = &lt;nil&gt;
		exact_match = &lt;nil&gt;
		name = &lt;nil&gt;
		operator = &lt;nil&gt;

  }
  total_count = &lt;nil&gt;
  cluster_set {
		status = ""
		update_time = &lt;nil&gt;
		zone = &lt;nil&gt;
		cluster_name = &lt;nil&gt;
		region = &lt;nil&gt;
		db_version = &lt;nil&gt;
		cluster_id = &lt;nil&gt;
		instance_num = &lt;nil&gt;
		uin = &lt;nil&gt;
		db_type = &lt;nil&gt;
		app_id = &lt;nil&gt;
		status_desc = &lt;nil&gt;
		create_time = ""
		pay_mode = &lt;nil&gt;
		period_end_time = &lt;nil&gt;
		vip = &lt;nil&gt;
		vport = &lt;nil&gt;
		project_i_d = &lt;nil&gt;
		vpc_id = &lt;nil&gt;
		subnet_id = &lt;nil&gt;
		cynos_version = &lt;nil&gt;
		storage_limit = &lt;nil&gt;
		renew_flag = &lt;nil&gt;
		processing_task = &lt;nil&gt;
		tasks {
			task_id = &lt;nil&gt;
			task_type = &lt;nil&gt;
			task_status = &lt;nil&gt;
			object_id = &lt;nil&gt;
			object_type = &lt;nil&gt;
		}
		resource_tags {
			tag_key = &lt;nil&gt;
			tag_value = &lt;nil&gt;
		}
		db_mode = &lt;nil&gt;
		serverless_status = &lt;nil&gt;
		storage = &lt;nil&gt;
		storage_id = &lt;nil&gt;
		storage_pay_mode = &lt;nil&gt;
		min_storage_size = &lt;nil&gt;
		max_storage_size = &lt;nil&gt;
		net_addrs {
			vip = &lt;nil&gt;
			vport = &lt;nil&gt;
			wan_domain = &lt;nil&gt;
			wan_port = &lt;nil&gt;
			net_type = &lt;nil&gt;
			uniq_subnet_id = &lt;nil&gt;
			uniq_vpc_id = &lt;nil&gt;
			description = &lt;nil&gt;
			wan_i_p = &lt;nil&gt;
			wan_status = &lt;nil&gt;
		}
		physical_zone = &lt;nil&gt;
		master_zone = &lt;nil&gt;
		has_slave_zone = &lt;nil&gt;
		slave_zones = &lt;nil&gt;
		business_type = &lt;nil&gt;
		is_freeze = &lt;nil&gt;
		order_source = &lt;nil&gt;
		ability {
			is_support_slave_zone = &lt;nil&gt;
			nonsupport_slave_zone_reason = &lt;nil&gt;
			is_support_ro = &lt;nil&gt;
			nonsupport_ro_reason = &lt;nil&gt;
		}

  }
}

`
