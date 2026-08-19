package emr_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudEmrClusterV2_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccEmrClusterV2_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_emr_cluster_v2.example", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_emr_cluster_v2.example",
				ImportState:       true,
				ImportStateVerify: true,
				// Fields not round-tripped by DescribeInstances (secrets, idempotency token, custom conf).
				ImportStateVerifyIgnore: []string{
					"login_settings",
					"meta_db_info",
					"custom_conf",
					"zone_resource_configuration",
					"script_bootstrap_action_config",
					"node_marks",
				},
			},
		},
	})
}

const testAccEmrClusterV2_basic = `
resource "tencentcloud_emr_cluster_v2" "example" {
  product_version        = "EMR-V3.5.0"
  enable_support_ha_flag = false
  instance_name          = "tf-acc-emr-v2"
  instance_charge_type   = "POSTPAID_BY_HOUR"
  need_master_wan        = "NEED_MASTER_WAN"

  login_settings {
    password = "Tencent@2026"
  }

  scene_software_config {
    scene_name = "Hadoop-Default"
    software   = [
      "hdfs-3.2.2",
      "yarn-3.2.2",
      "zookeeper-3.6.3",
      "openldap-2.4.44",
      "knox-1.6.1",
      "hive-3.1.3",
    ]
  }

  meta_db_info {
    meta_type = "EMR_DEFAULT_META"
  }

  tags {
    tag_key   = "createBy"
    tag_value = "terraform"
  }

  zone_resource_configuration {
    virtual_private_cloud {
      vpc_id    = "vpc-xxxxxxxx"
      subnet_id = "subnet-xxxxxxxx"
    }

    placement {
      zone       = "ap-guangzhou-7"
      project_id = 0
    }

    all_node_resource_spec {
      # One master_resource_spec block = 1 master node.
      master_resource_spec {
        instance_type = "S6.2XLARGE32"
        system_disk {
          count     = 1
          disk_size = 50
          disk_type = "CLOUD_HSSD"
        }
      }

      # Two core_resource_spec blocks = 2 core nodes (first block is the template).
      core_resource_spec {
        instance_type = "SA4.8XLARGE64"
        system_disk {
          count     = 1
          disk_size = 50
          disk_type = "CLOUD_SSD"
        }
        data_disk {
          count     = 1
          disk_size = 300
          disk_type = "CLOUD_TSSD"
        }
      }

      core_resource_spec {
        instance_type = "SA4.8XLARGE64"
        system_disk {
          count     = 1
          disk_size = 50
          disk_type = "CLOUD_SSD"
        }
        data_disk {
          count     = 1
          disk_size = 300
          disk_type = "CLOUD_TSSD"
        }
      }
    }
  }
}
`
