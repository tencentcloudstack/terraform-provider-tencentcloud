package emr_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudEmrAutoScaleStrategyResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{{
			Config: testAccEmrAutoScaleStrategy,
			Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_emr_auto_scale_strategy.emr_auto_scale_strategy", "id")),
		}, {
			ResourceName:      "tencentcloud_emr_auto_scale_strategy.emr_auto_scale_strategy",
			ImportState:       true,
			ImportStateVerify: true,
		}},
	})
}

const testAccEmrAutoScaleStrategy = `
resource "tencentcloud_emr_auto_scale_strategy" "emr_auto_scale_strategy" {
  instance_id   = "emr-rzrochgp"
  strategy_type = 2
  time_auto_scale_strategy {
    strategy_name    = "tf-test1"
    interval_time    = 100
    scale_action     = 1
    scale_num        = 1
    strategy_status  = 1
    retry_valid_time = 60
    repeat_strategy {
      repeat_type = "DAY"
      day_repeat {
        execute_at_time_of_day = "16:30:00"
        step                   = 1
      }
      expire = "2026-02-20 23:59:59"
    }
    grace_down_flag = false
    tags {
      tag_key   = "createBy"
      tag_value = "terraform"
    }
    config_group_assigned = "{\"HDFS-2.8.5\":-1,\"YARN-2.8.5\":-1}"
    measure_method        = "INSTANCE"
    terminate_policy      = "DEFAULT"
    soft_deploy_info      = [1, 2]
    service_node_info     = [7]
    priority              = 1
  }
  time_auto_scale_strategy {
    strategy_name    = "tf-test2"
    interval_time    = 100
    scale_action     = 1
    scale_num        = 1
    strategy_status  = 1
    retry_valid_time = 60
    repeat_strategy {
      repeat_type = "DAY"
      day_repeat {
        execute_at_time_of_day = "17:30:00"
        step                   = 1
      }
      expire = "2026-02-20 23:59:59"
    }
    grace_down_flag = false
    tags {
      tag_key   = "createBy"
      tag_value = "terraform"
    }
    config_group_assigned = "{\"HDFS-2.8.5\":-1,\"YARN-2.8.5\":-1}"
    measure_method        = "INSTANCE"
    terminate_policy      = "DEFAULT"
    soft_deploy_info      = [1, 2]
    service_node_info     = [7]
    priority              = 2
  }
}
`
