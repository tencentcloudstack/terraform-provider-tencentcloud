package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudClsConfigAttachment_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccClsConfigAttachment,
			},
		},
	})
}

const testAccClsConfigAttachment = `
resource "tencentcloud_cls_logset" "logset" {
  logset_name = "tf-attach-test"
  tags        = {
    "test" = "test"
  }
}

resource "tencentcloud_cls_topic" "topic" {
  auto_split           = true
  logset_id            = tencentcloud_cls_logset.logset.id
  max_split_partitions = 20
  partition_count      = 1
  period               = 10
  storage_type         = "hot"
  tags                 = {
    "test" = "test"
  }
  topic_name           = "tf-attach-test"
}

resource "tencentcloud_cls_config" "config" {
  name             = "attach"
  output           = tencentcloud_cls_topic.topic.id
  path             = "/var/log/kubernetes/**/kubernetes.audit"
  log_type         = "json_log"
  extract_rule {
    filter_key_regex {
      key   = "key1"
      regex = "value1"
    }
    filter_key_regex {
      key   = "key2"
      regex = "value2"
    }
    un_match_up_load_switch = true
    un_match_log_key        = "config"
    backtracking            = -1
  }
  exclude_paths {
    type  = "Path"
    value = "/data"
  }
  exclude_paths {
    type  = "File"
    value = "/file"
  }
}

resource "tencentcloud_cls_machine_group" "group" {
  group_name        = "tf-attach-group"
  service_logging   = true
  auto_update       = true
  update_end_time   = "19:05:00"
  update_start_time = "17:05:00"

  machine_group_type {
    type   = "ip"
    values = [
      "192.168.1.1",
      "192.168.1.2",
    ]
  }
}

resource "tencentcloud_cls_config_attachment" "attach" {
  config_id = tencentcloud_cls_config.config.id
  group_id = tencentcloud_cls_machine_group.group.id
}
`
