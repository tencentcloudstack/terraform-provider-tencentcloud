package cvm_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudCvmLaunchTemplateVersionResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCvmLaunchTemplateVersion,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_cvm_launch_template_version.test_launch_tpl_v2", "id")),
			},
		},
	})
}

const testAccCvmLaunchTemplateVersion = `
data "tencentcloud_images" "default" {
  image_type = ["PUBLIC_IMAGE"]
  image_name_regex = "Final"
}
data "tencentcloud_vpc_instances" "vpc" {
  is_default = true
}
data "tencentcloud_vpc_subnets" "subnets" {
  availability_zone = "ap-guangzhou-7"
  vpc_id = data.tencentcloud_vpc_instances.vpc.instance_list.0.vpc_id
}

resource "tencentcloud_cvm_launch_template" "test_launch_tpl" {
  launch_template_name = "test"
  image_id             = data.tencentcloud_images.default.images.0.image_id
  placement {
    zone = "ap-guangzhou-7"
  }
  instance_name = "v1"
}

resource "tencentcloud_cvm_launch_template_version" "test_launch_tpl_v2" {
  launch_template_id = tencentcloud_cvm_launch_template.test_launch_tpl.id
  placement {
    zone = "ap-guangzhou-7"
    project_id = 0
  }
  launch_template_version_description = "test"
  instance_type = "S5.MEDIUM2"
  image_id             = data.tencentcloud_images.default.images.0.image_id
#   system_disk {
#     disk_type = "CLOUD_PREMIUM"
#     disk_size = 20
#   }
#   data_disks {
#     disk_type = "CLOUD_PREMIUM"
#     disk_size = 200
#   }
  virtual_private_cloud {
    vpc_id    = data.tencentcloud_vpc_subnets.subnets.instance_list.0.vpc_id
    subnet_id = data.tencentcloud_vpc_subnets.subnets.instance_list.0.subnet_id
  }
#   internet_accessible {
#     internet_charge_type = "TRAFFIC_POSTPAID_BY_HOUR"
#     internet_max_bandwidth_out = 20
#   }
  instance_count = 2
  instance_name = "v2"
  security_group_ids = ["sg-5275dorp","sg-cm7fbbf3"]
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
  client_token = "123"
  host_name = "test"
  disaster_recover_group_ids = []
  tag_specification {
    resource_type = "instance"
    tags {
      key = "key"
      value = "value2"
    }
    tags {
      key = "key2"
      value = "value1"
    }
  }
  user_data = "aGhoCg=="
  cam_role_name = ""
  hpc_cluster_id = ""
  instance_charge_type = "PREPAID"
  instance_charge_prepaid {
    period = 3
    renew_flag = "DISABLE_NOTIFY_AND_MANUAL_RENEW"
  }
  disable_api_termination = true
}
`
