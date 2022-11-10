package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudPtsScenario_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPtsScenario,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_pts_scenario.scenario", "id"),
				),
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
  type = ""
  project_id = ""
  description = ""
  load {
		load_spec {
			concurrency {
				stages {
						duration_seconds = ""
						target_virtual_users = ""
				}
					iteration_count = ""
					max_requests_per_second = ""
					graceful_stop_seconds = ""
			}
			requests_per_second {
					max_requests_per_second = ""
					duration_seconds = ""
					target_virtual_users = ""
					resources = ""
					start_requests_per_second = ""
					target_requests_per_second = ""
					graceful_stop_seconds = ""
			}
			script_origin {
					machine_number = ""
					machine_specification = ""
					duration_seconds = ""
			}
		}
		vpc_load_distribution {
				region_id = ""
				region = ""
				vpc_id = ""
				subnet_ids = ""
		}
		geo_regions_load_distribution {
				region_id = ""
				region = ""
				percentage = ""
		}

  }
  configs = ""
  datasets {
			name = ""
			split = ""
			header_in_file = ""
			header_columns = ""
			line_count = ""
			updated_at = ""
			size = ""
			head_lines = ""
			tail_lines = ""
			type = ""
			file_id = ""

  }
  extensions = ""
  s_l_a_id = ""
  cron_id = ""
  scripts = ""
  test_scripts {
			name = ""
			size = ""
			type = ""
			updated_at = ""
			encoded_content = ""
			encoded_http_archive = ""
			load_weight = ""

  }
  protocols {
			name = ""
			size = ""
			type = ""
			updated_at = ""
			file_id = ""

  }
  request_files {
			name = ""
			size = ""
			type = ""
			updated_at = ""
			file_id = ""

  }
  s_l_a_policy {
		s_l_a_rules {
				metric = ""
				aggregation = ""
				condition = ""
				value = ""
			label_filter {
					label_name = ""
					label_value = ""
			}
				abort_flag = ""
				for = ""
		}
		alert_channel {
				notice_id = ""
				a_m_p_consumer_id = ""
		}

  }
  plugins {
			name = ""
			size = ""
			type = ""
			updated_at = ""
			file_id = ""

  }
  domain_name_config {
		host_aliases {
				host_names = ""
				i_p = ""
		}
		d_n_s_config {
				nameservers = ""
		}

  }
                }

`
