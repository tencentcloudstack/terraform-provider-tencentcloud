package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudClsConfig_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccClsConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_cls_config.config", "name", "config"),
				),
			},
		},
	})
}

func TestAccTencentCloudClsConfig_FullRegex(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccClsFullRegexConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_cls_config.config", "name", "tf-full-regex-config-test"),
				),
			},
		},
	})
}

const testAccClsConfig = `
resource "tencentcloud_cls_logset" "logset" {
  logset_name = "tf-config-test"
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
  topic_name           = "tf-config-test"
}

resource "tencentcloud_cls_config" "config" {
  name             = "config"
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
`
const testAccClsFullRegexConfig = `
resource "tencentcloud_cls_logset" "logset" {
  logset_name = "tf-full-regex-config-test"
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
  topic_name           = "tf-full-regex-config-test"
}

resource "tencentcloud_cls_config" "config" {
  name     = "tf-full-regex-config-test"
  output   = tencentcloud_cls_topic.topic.id
  path     = "/var/log/nginx/**/access.log"
  log_type = "fullregex_log"

  extract_rule {
    begin_regex = "\\d+\\.\\d+\\.\\d+\\.\\d+\\s+-\\s+.*"
    log_regex = "\\d+\\.\\d+\\.\\d+\\.\\d+\\s+-\\s+(.*)"
    keys = ["acd", "edf"]
  }
}
`
