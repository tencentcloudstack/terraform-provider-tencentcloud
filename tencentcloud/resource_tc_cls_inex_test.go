package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudClsIndex_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccClsIndex,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_cls_index.index", "status", "true"),
				),
			},
			{
				ResourceName:      "tencentcloud_cls_index.index",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccClsIndex = `
resource "tencentcloud_cls_logset" "logset" {
  logset_name = "tf-topic-index"
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
  topic_name           = "tf-topic-index"
}

resource "tencentcloud_cls_index" "index" {
  topic_id = tencentcloud_cls_topic.topic.id

  rule {
    full_text {
      case_sensitive = true
      tokenizer      = "@&?|#()='\",;:<>[]{}"
      contain_z_h    = true
    }

    key_value {
      case_sensitive = true
      key_values {
        key = "hello"
        value {
          contain_z_h = true
          sql_flag    = true
          tokenizer   = "@&?|#()='\",;:<>[]{}"
          type        = "text"
        }
      }

      key_values {
        key = "world"
        value {
          contain_z_h = true
          sql_flag    = true
          tokenizer   = "@&?|#()='\",;:<>[]{}"
          type        = "text"
        }
      }
    }

    tag {
      case_sensitive = true
      key_values {
        key = "terraform"
        value {
          contain_z_h = true
          sql_flag    = true
          tokenizer   = "@&?|#()='\",;:<>[]{}"
          type        = "text"
        }
      }
    }
  }
  status                  = true
  include_internal_fields = true
  metadata_flag           = 1
}





`
