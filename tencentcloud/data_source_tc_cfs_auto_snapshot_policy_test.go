package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCfsAutoSnapshotPolicyDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCfsAutoSnapshotPolicyDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_cfs_auto_snapshot_policy.auto_snapshot_policy")),
			},
		},
	})
}

const testAccCfsAutoSnapshotPolicyDataSource = `

data "tencentcloud_cfs_auto_snapshot_policy" "auto_snapshot_policy" {
  auto_snapshot_policy_id = &lt;nil&gt;
  offset = &lt;nil&gt;
  limit = &lt;nil&gt;
  filters {
		values = &lt;nil&gt;
		name = &lt;nil&gt;

  }
  order = &lt;nil&gt;
  order_field = &lt;nil&gt;
  total_count = &lt;nil&gt;
  auto_snapshot_policies {
		auto_snapshot_policy_id = &lt;nil&gt;
		policy_name = &lt;nil&gt;
		creation_time = &lt;nil&gt;
		file_system_nums = &lt;nil&gt;
		day_of_week = &lt;nil&gt;
		hour = &lt;nil&gt;
		is_activated = &lt;nil&gt;
		next_active_time = &lt;nil&gt;
		status = &lt;nil&gt;
		app_id = &lt;nil&gt;
		alive_days = &lt;nil&gt;
		region_name = &lt;nil&gt;
		file_systems {
			creation_token = &lt;nil&gt;
			file_system_id = &lt;nil&gt;
			size_byte = &lt;nil&gt;
			storage_type = &lt;nil&gt;
			total_snapshot_size = &lt;nil&gt;
			creation_time = &lt;nil&gt;
			zone_id = &lt;nil&gt;
		}

  }
  request_id = &lt;nil&gt;
}

`
