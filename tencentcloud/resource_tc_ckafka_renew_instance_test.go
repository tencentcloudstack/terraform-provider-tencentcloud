package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudCkafkaRenewInstanceResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCkafkaRenewInstance,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_ckafka_renew_instance.renew_instance", "id")),
			},
		},
	})
}

const testAccCkafkaRenewInstance = defaultKafkaVariable + `
resource "tencentcloud_ckafka_instance" "renew_instance" {
	instance_name      = "ckafka-instance-renew-test"
	zone_id            = 100003
	period             = 1
	vpc_id             = var.vpc_id
	subnet_id          = var.subnet_id
	msg_retention_time = 1300
	max_message_byte   = 1024
	renew_flag         = 0
	kafka_version      = "1.1.1"
	disk_size          = 500
	disk_type          = "CLOUD_BASIC"
  
  
	config {
	  auto_create_topic_enable   = true
	  default_num_partitions     = 3
	  default_replication_factor = 3
	}
  
	dynamic_retention_config {
	  enable = 1
	}
}

resource "tencentcloud_ckafka_renew_instance" "renew_instance" {
  instance_id = tencentcloud_ckafka_instance.renew_instance.id
  time_span = 1
}

`
