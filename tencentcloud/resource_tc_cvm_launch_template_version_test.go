package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCvmLaunchTemplateVersionResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCvmLaunchTemplateVersion,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_cvm_launch_template_version.launch_template_version", "id")),
			},
			{
				ResourceName:      "tencentcloud_cvm_launch_template_version.launch_template_version",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCvmLaunchTemplateVersion = `

resource "tencentcloud_cvm_launch_template_version" "launch_template_version" {
  placement {
		zone = "ap-guangzhou-6"
		project_id = 0
		host_ids = 
		host_ips = 
		host_id = "host-dab6ejhx"

  }
  launch_template_id = "lt-lobxe2yo"
  launch_template_version = 1
  launch_template_version_description = ""
  instance_type = "S5.MEDIUM4"
  image_id = "img-eb30mz89"
  system_disk {
		disk_type = "CLOUD_PREMIUM"
		disk_id = ""
		disk_size = 50
		cdc_id = "cdc-b9pbd3px"

  }
  data_disks {
		disk_size = 50
		disk_type = "CLOUD_PREMIUM"
		disk_id = ""
		delete_with_instance = false
		snapshot_id = "snap-r9unnd89"
		encrypt = false
		kms_key_id = "kms-abcd1234"
		throughput_performance = 2
		cdc_id = "cdc-b9pbd3px"

  }
  virtual_private_cloud {
		vpc_id = "vpc-x2e4dam7"
		subnet_id = "subnet-g73bdf1r"
		as_vpc_gateway = false
		private_ip_addresses = 
		ipv6_address_count = 1

  }
  internet_accessible {
		internet_charge_type = "TRAFFIC_POSTPAID_BY_HOUR"
		internet_max_bandwidth_out = 0
		public_ip_assigned = false
		bandwidth_package_id = ""

  }
  instance_count = 1
  instance_name = ""
  login_settings {
		password = "1@34qwer"
		key_ids = 
		keep_image_login = "FALSE"

  }
  security_group_ids = 
  enhanced_service {
		security_service {
			enabled = true
		}
		monitor_service {
			enabled = true
		}
		automation_service {
			enabled = false
		}

  }
  client_token = "cg_1678809600005774197"
  host_name = ""
  action_timer {
		timer_action = "TerminateInstances"
		action_time = "2018-05-29T11:26:40Z"
		externals {
			release_address = false
			unsupport_networks = 
			storage_block_attr {
				type = "LOCAL_PRO"
				min_size = 10
				max_size = 100
			}
		}

  }
  disaster_recover_group_ids = 
  tag_specification {
		resource_type = "instance"
		tags {
			key = "tagKey"
			value = "tagValue"
		}

  }
  instance_market_options {
		spot_options {
			max_price = "1.99"
			spot_instance_type = "one-time"
		}
		market_type = "spot"

  }
  user_data = "IyEvdXNyL2Jpbi9lbnYgcHl0a"
  dry_run = false
  cam_role_name = ""
  hpc_cluster_id = "hpc-bwu6b3e2"
  instance_charge_type = "POSTPAID_BY_HOUR"
  instance_charge_prepaid {
		period = 1
		renew_flag = "NOTIFY_AND_AUTO_RENEW"

  }
  disable_api_termination = false
}

`
