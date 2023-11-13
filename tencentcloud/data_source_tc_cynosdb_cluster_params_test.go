package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCynosdbClusterParamsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCynosdbClusterParamsDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_cynosdb_cluster_params.cluster_params")),
			},
		},
	})
}

const testAccCynosdbClusterParamsDataSource = `

data "tencentcloud_cynosdb_cluster_params" "cluster_params" {
  cluster_id = &lt;nil&gt;
  param_name = &lt;nil&gt;
  total_count = &lt;nil&gt;
  items {
		current_value = &lt;nil&gt;
		default = &lt;nil&gt;
		enum_value = &lt;nil&gt;
		max = &lt;nil&gt;
		min = &lt;nil&gt;
		param_name = &lt;nil&gt;
		need_reboot = &lt;nil&gt;
		param_type = &lt;nil&gt;
		match_type = &lt;nil&gt;
		match_value = &lt;nil&gt;
		description = &lt;nil&gt;
		is_global = &lt;nil&gt;
		modifiable_info = &lt;nil&gt;
		is_func = &lt;nil&gt;
		func = &lt;nil&gt;

  }
}

`
