package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudPtsScenarioResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPtsScenario,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_pts_scenario.scenario", "id")),
			},
			{
				ResourceName:      "tencentcloud_pts_scenario.scenario",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccPtsScenario = `

resource "tencentcloud_pts_scenario" "scenario" {
  name = "pts"
  type = &lt;nil&gt;
  project_id = &lt;nil&gt;
  description = &lt;nil&gt;
  load {
		load_spec {
			concurrency {
				stages {
					duration_seconds = &lt;nil&gt;
					target_virtual_users = &lt;nil&gt;
				}
				iteration_count = &lt;nil&gt;
				max_requests_per_second = &lt;nil&gt;
				graceful_stop_seconds = &lt;nil&gt;
			}
			requests_per_second {
				max_requests_per_second = &lt;nil&gt;
				duration_seconds = &lt;nil&gt;
				resources = &lt;nil&gt;
				start_requests_per_second = &lt;nil&gt;
				target_requests_per_second = &lt;nil&gt;
				graceful_stop_seconds = &lt;nil&gt;
			}
			script_origin {
				machine_number = &lt;nil&gt;
				machine_specification = &lt;nil&gt;
				duration_seconds = &lt;nil&gt;
			}
		}
		vpc_load_distribution {
			region_id = &lt;nil&gt;
			region = &lt;nil&gt;
			vpc_id = &lt;nil&gt;
			subnet_ids = &lt;nil&gt;
		}
		geo_regions_load_distribution {
			region_id = &lt;nil&gt;
			region = &lt;nil&gt;
			percentage = &lt;nil&gt;
		}

  }
  datasets {
		name = &lt;nil&gt;
		split = &lt;nil&gt;
		header_in_file = &lt;nil&gt;
		header_columns = &lt;nil&gt;
		line_count = &lt;nil&gt;
		updated_at = &lt;nil&gt;
		size = &lt;nil&gt;
		head_lines = &lt;nil&gt;
		tail_lines = &lt;nil&gt;
		type = &lt;nil&gt;
		file_id = &lt;nil&gt;

  }
  cron_id = &lt;nil&gt;
  test_scripts {
		name = &lt;nil&gt;
		size = &lt;nil&gt;
		type = &lt;nil&gt;
		updated_at = &lt;nil&gt;
		encoded_content = &lt;nil&gt;
		encoded_http_archive = &lt;nil&gt;
		load_weight = &lt;nil&gt;

  }
  protocols {
		name = &lt;nil&gt;
		size = &lt;nil&gt;
		type = &lt;nil&gt;
		updated_at = &lt;nil&gt;
		file_id = &lt;nil&gt;

  }
  request_files {
		name = &lt;nil&gt;
		size = &lt;nil&gt;
		type = &lt;nil&gt;
		updated_at = &lt;nil&gt;
		file_id = &lt;nil&gt;

  }
  s_l_a_policy {
		s_l_a_rules {
			metric = &lt;nil&gt;
			aggregation = &lt;nil&gt;
			condition = &lt;nil&gt;
			value = &lt;nil&gt;
			label_filter {
				label_name = &lt;nil&gt;
				label_value = &lt;nil&gt;
			}
			abort_flag = &lt;nil&gt;
			for = &lt;nil&gt;
		}
		alert_channel {
			notice_id = &lt;nil&gt;
			a_m_p_consumer_id = &lt;nil&gt;
		}

  }
  plugins {
		name = &lt;nil&gt;
		size = &lt;nil&gt;
		type = &lt;nil&gt;
		updated_at = &lt;nil&gt;
		file_id = &lt;nil&gt;

  }
  domain_name_config {
		host_aliases {
			host_names = &lt;nil&gt;
			i_p = &lt;nil&gt;
		}
		d_n_s_config {
			nameservers = &lt;nil&gt;
		}

  }
              }

`
