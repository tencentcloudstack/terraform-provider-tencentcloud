package ckafka_test

import (
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudInternationalCkafkaResource_instance(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccInternationalKafkaInstancePostpaid,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKafkaInstanceExists("tencentcloud_ckafka_instance.kafka_instance_postpaid"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_instance.kafka_instance_postpaid", "instance_name", "ckafka-instance-postpaid"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_instance.kafka_instance_postpaid", "zone_id", "100007"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_instance.kafka_instance_postpaid", "msg_retention_time", "1300"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_instance.kafka_instance_postpaid", "kafka_version", "1.1.1"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_instance.kafka_instance_postpaid", "disk_size", "500"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_instance.kafka_instance_postpaid", "disk_type", "CLOUD_BASIC"),
					resource.TestCheckResourceAttrSet("tencentcloud_ckafka_instance.kafka_instance_postpaid", "vip"),
					resource.TestCheckResourceAttrSet("tencentcloud_ckafka_instance.kafka_instance_postpaid", "vport"),
				),
			},
			{
				Config: testAccInternationalKafkaInstanceUpdatePostpaid,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKafkaInstanceExists("tencentcloud_ckafka_instance.kafka_instance_postpaid"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_instance.kafka_instance_postpaid", "instance_name", "ckafka-instance-postpaid"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_instance.kafka_instance_postpaid", "zone_id", "100007"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_instance.kafka_instance_postpaid", "msg_retention_time", "1200"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_instance.kafka_instance_postpaid", "kafka_version", "1.1.1"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_instance.kafka_instance_postpaid", "disk_size", "500"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_instance.kafka_instance_postpaid", "disk_type", "CLOUD_BASIC"),
				),
			},
			{
				PreConfig: func() {
					time.Sleep(60 * time.Second)
				},
				ResourceName:            "tencentcloud_ckafka_instance.kafka_instance_postpaid",
				ImportState:             true,
				ImportStateVerifyIgnore: []string{"period", "max_message_byte", "charge_type", "upgrade_strategy"},
			},
		},
	})
}

const testAccInternationalKafkaInstancePostpaid = tcacctest.DefaultKafkaVariable + `
resource "tencentcloud_ckafka_instance" "kafka_instance_postpaid" {
  instance_name      = "ckafka-instance-postpaid"
  zone_id            = 100007
  vpc_id             = var.international_vpc_id
  subnet_id          = var.international_subnet_id
  msg_retention_time = 1300
  kafka_version      = "1.1.1"
  disk_size          = 500
  band_width         = 20
  disk_type          = "CLOUD_BASIC"
  partition          = 400
  charge_type        = "POSTPAID_BY_HOUR"

  config {
    auto_create_topic_enable   = true
    default_num_partitions     = 3
    default_replication_factor = 3
  }

  dynamic_retention_config {
    enable = 1
  }
}

resource "tencentcloud_ckafka_topic" "foo" {
	instance_id                     = tencentcloud_ckafka_instance.kafka_instance_postpaid.id
	topic_name                      = "tmp"
	note                            = "topic note"
	replica_num                     = 2
	partition_num                   = 1
	clean_up_policy                 = "delete"
	sync_replica_min_num            = 1
	unclean_leader_election_enable  = false
	retention                       = 60000
}
`

const testAccInternationalKafkaInstanceUpdatePostpaid = tcacctest.DefaultKafkaVariable + `
resource "tencentcloud_ckafka_instance" "kafka_instance_postpaid" {
  instance_name      = "ckafka-instance-postpaid"
  zone_id            = 100007
  vpc_id             =  var.international_vpc_id
  subnet_id          =  var.international_subnet_id
  msg_retention_time = 1200
  kafka_version      = "1.1.1"
  disk_type          = "CLOUD_BASIC"
  disk_size          = 500
  band_width         = 20
  charge_type        = "POSTPAID_BY_HOUR"

  config {
    auto_create_topic_enable   = true
    default_num_partitions     = 3
    default_replication_factor = 3
  }

  dynamic_retention_config {
    enable = 1
  }

  tag_set = {
    createdBy = "terraform"
  }
}

resource "tencentcloud_ckafka_topic" "foo" {
	instance_id                     = tencentcloud_ckafka_instance.kafka_instance_postpaid.id
	topic_name                      = "tmp"
	note                            = "topic note"
	replica_num                     = 2
	partition_num                   = 1
	clean_up_policy                 = "delete"
	sync_replica_min_num            = 1
	unclean_leader_election_enable  = false
	retention                       = 60000
}
`

const testAccInternationalKafkaInstanceUpdatePostpaidDiskSize = tcacctest.DefaultKafkaVariable + `
resource "tencentcloud_ckafka_instance" "kafka_instance_postpaid" {
  instance_name      = "ckafka-instance-postpaid"
  zone_id            = 100007
  vpc_id             =  var.international_vpc_id
  subnet_id          =  var.international_subnet_id
  msg_retention_time = 1200
  kafka_version      = "1.1.1"
  disk_type          = "CLOUD_BASIC"
  disk_size          = 400
  band_width         = 20
  charge_type        = "POSTPAID_BY_HOUR"

  config {
    auto_create_topic_enable   = true
    default_num_partitions     = 3
    default_replication_factor = 3
  }

  dynamic_retention_config {
    enable = 1
  }

  tag_set = {
    createdBy = "terraform"
  }
}

resource "tencentcloud_ckafka_topic" "foo" {
	instance_id                     = tencentcloud_ckafka_instance.kafka_instance_postpaid.id
	topic_name                      = "tmp"
	note                            = "topic note"
	replica_num                     = 2
	partition_num                   = 1
	clean_up_policy                 = "delete"
	sync_replica_min_num            = 1
	unclean_leader_election_enable  = false
	retention                       = 60000
}
`
