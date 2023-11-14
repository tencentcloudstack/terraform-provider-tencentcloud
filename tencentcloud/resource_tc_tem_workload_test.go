package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudTemWorkloadResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTemWorkload,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_tem_workload.workload", "id")),
			},
			{
				ResourceName:      "tencentcloud_tem_workload.workload",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTemWorkload = `

resource "tencentcloud_tem_workload" "workload" {
  application_id = "app-xxx"
  environment_id = "en-xxx"
  deploy_version = "hello-world"
  deploy_mode = "IMAGE"
  img_repo = "tem_demo/tem_demo"
  init_pod_num = 1
  cpu_spec = 
  memory_spec = 
  post_start = &lt;nil&gt;
  pre_stop = &lt;nil&gt;
  security_group_ids = 
  repo_type = 3
  repo_server = &lt;nil&gt;
  tcr_instance_id = &lt;nil&gt;
  env_conf {
		key = "key"
		value = "value"
		type = "default"
		config = "config-name"
		secret = "secret-name"

  }
  storage_confs {
		storage_vol_name = "xxx"
		storage_vol_ip = "0.0.0.0"
		storage_vol_path = "/"

  }
  storage_mount_confs {
		volume_name = "xxx"
		mount_path = "/"

  }
  liveness {
		type = "HttpGet"
		protocol = "HTTP"
		path = "/"
		exec = ""
		port = 80
		initial_delay_seconds = 0
		timeout_seconds = 1
		period_seconds = 10

  }
  readiness {
		type = "HttpGet"
		protocol = "HTTP"
		path = "/"
		exec = ""
		port = 80
		initial_delay_seconds = 0
		timeout_seconds = 1
		period_seconds = 10

  }
  startup_probe {
		type = "HttpGet"
		protocol = "HTTP"
		path = "/"
		exec = ""
		port = 80
		initial_delay_seconds = 0
		timeout_seconds = 1
		period_seconds = 10

  }
  deploy_strategy_conf {
		deploy_strategy_type = 0
		beta_batch_num = 0
		total_batch_count = 1
		batch_interval = 200
		min_available = -1
		force = true

  }
}

`
