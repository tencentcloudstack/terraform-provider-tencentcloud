package ckafka_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	localckafka "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/ckafka"
)

func TestAccTencentCloudCkafkaInstanceResource_prepaid(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccTencentCloudKafkaInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKafkaInstance,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKafkaInstanceExists("tencentcloud_ckafka_instance.kafka_instance"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_instance.kafka_instance", "instance_name", "ckafka-instance-prepaid"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_instance.kafka_instance", "zone_id", "100007"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_instance.kafka_instance", "charge_type", "PREPAID"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_instance.kafka_instance", "period", "1"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_instance.kafka_instance", "msg_retention_time", "1300"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_instance.kafka_instance", "max_message_byte", "1024"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_instance.kafka_instance", "renew_flag", "0"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_instance.kafka_instance", "kafka_version", "1.1.1"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_instance.kafka_instance", "disk_size", "200"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_instance.kafka_instance", "disk_type", "CLOUD_BASIC"),
					resource.TestCheckResourceAttrSet("tencentcloud_ckafka_instance.kafka_instance", "vip"),
					resource.TestCheckResourceAttrSet("tencentcloud_ckafka_instance.kafka_instance", "vport"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_instance.kafka_instance", "elastic_bandwidth_switch", "1"),
				),
			},
			{
				Config: testAccKafkaInstanceUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKafkaInstanceExists("tencentcloud_ckafka_instance.kafka_instance"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_instance.kafka_instance", "instance_name", "ckafka-instance-prepaid"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_instance.kafka_instance", "zone_id", "100007"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_instance.kafka_instance", "charge_type", "PREPAID"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_instance.kafka_instance", "period", "1"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_instance.kafka_instance", "msg_retention_time", "1200"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_instance.kafka_instance", "max_message_byte", "1025"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_instance.kafka_instance", "renew_flag", "0"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_instance.kafka_instance", "kafka_version", "1.1.1"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_instance.kafka_instance", "disk_size", "300"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_instance.kafka_instance", "disk_type", "CLOUD_BASIC"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_instance.kafka_instance", "elastic_bandwidth_switch", "1"),
				),
			},
			{
				PreConfig: func() {
					time.Sleep(2 * time.Minute)
				},
				ResourceName:            "tencentcloud_ckafka_instance.kafka_instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"period", "max_message_byte", "upgrade_strategy", "elastic_bandwidth_switch"},
			},
		},
	})
}

func TestAccTencentCloudCkafkaInstanceResource_postpaid(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccTencentCloudKafkaInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKafkaInstancePostpaid,
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
				Config: testAccKafkaInstanceUpdatePostpaid,
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
				Config: testAccKafkaInstanceUpdatePostpaidDiskSize,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKafkaInstanceExists("tencentcloud_ckafka_instance.kafka_instance_postpaid"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_instance.kafka_instance_postpaid", "disk_size", "400"),
				),
			},
			{
				ResourceName:            "tencentcloud_ckafka_instance.kafka_instance_postpaid",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"period", "max_message_byte", "upgrade_strategy", "elastic_bandwidth_switch"},
			},
		},
	})
}

func TestAccTencentCloudCkafkaInstanceResource_maz(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccTencentCloudKafkaInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKafkaInstanceMAZ,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKafkaInstanceExists("tencentcloud_ckafka_instance.kafka_instance"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_instance.kafka_instance", "instance_name", "ckafka-instance-maz-tf-test"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_instance.kafka_instance", "zone_id", "100007"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_instance.kafka_instance", "zone_ids.#", "2"),
				),
			},
			{
				PreConfig: func() {
					time.Sleep(2 * time.Minute)
				},
				ResourceName:            "tencentcloud_ckafka_instance.kafka_instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"period", "max_message_byte", "upgrade_strategy"},
			},
		},
	})
}

func TestAccTencentCloudCkafkaInstanceResource_type(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccTencentCloudKafkaInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKafkaInstanceType,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKafkaInstanceExists("tencentcloud_ckafka_instance.kafka_instance"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_instance.kafka_instance", "instance_name", "ckafka-instance-type-tf-test"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_instance.kafka_instance", "specifications_type", "profession"),
				),
			},
			{
				PreConfig: func() {
					time.Sleep(2 * time.Minute)
				},
				ResourceName:            "tencentcloud_ckafka_instance.kafka_instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"period", "max_message_byte", "upgrade_strategy"},
			},
		},
	})
}

func testAccTencentCloudKafkaInstanceDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	ckafkcService := localckafka.NewCkafkaService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, r := range s.RootModule().Resources {
		if r.Type != "tencentcloud_ckafka_instance" {
			continue
		}
		_, has, error := ckafkcService.DescribeCkafkaById(ctx, r.Primary.ID)
		if error != nil {
			return error
		}
		if !has {
			return nil
		}
		return fmt.Errorf("ckafka instance still exists: %s", r.Primary.ID)
	}
	return nil
}

func testAccCheckKafkaInstanceExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("ckafka instance %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("ckafka instance id is not set")
		}
		ckafkcService := localckafka.NewCkafkaService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())

		var exist bool
		outErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			_, has, inErr := ckafkcService.DescribeInstanceById(ctx, rs.Primary.ID)
			if inErr != nil {
				return tccommon.RetryError(inErr)
			}
			exist = has
			return nil
		})
		if outErr != nil {
			return outErr
		}
		if !exist {
			return fmt.Errorf("ckafka instance doesn't exist: %s", rs.Primary.ID)
		}
		return nil
	}
}

const testAccKafkaInstance = tcacctest.DefaultKafkaVariable + `
resource "tencentcloud_vpc" "kafka_vpc" {
	name       = "kafka-vpc"
	cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_subnet" "kafka_subnet" {
	vpc_id            = tencentcloud_vpc.kafka_vpc.id
	name              = "kafka-subnet"
	cidr_block        = "10.0.0.0/16"
	availability_zone = "ap-guangzhou-7"
	is_multicast      = false
}

resource "tencentcloud_ckafka_instance" "kafka_instance" {
  instance_name      = "ckafka-instance-prepaid"
  zone_id            = 100007
  period             = 1
  vpc_id             = tencentcloud_vpc.kafka_vpc.id
  subnet_id          = tencentcloud_subnet.kafka_subnet.id
  msg_retention_time = 1300
  max_message_byte   = 1024
  renew_flag         = 0
  kafka_version      = "1.1.1"
  disk_size          = 200
  band_width         = 20
  disk_type          = "CLOUD_BASIC"
  partition          = 400
  elastic_bandwidth_switch = 1

  config {
    auto_create_topic_enable   = true
    default_num_partitions     = 3
    default_replication_factor = 3
  }

  dynamic_retention_config {
    enable = 1
  }
}
`

const testAccKafkaInstanceUpdate = `
resource "tencentcloud_vpc" "kafka_vpc" {
	name       = "kafka-vpc"
	cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_subnet" "kafka_subnet" {
	vpc_id            = tencentcloud_vpc.kafka_vpc.id
	name              = "kafka-subnet"
	cidr_block        = "10.0.0.0/16"
	availability_zone = "ap-guangzhou-7"
	is_multicast      = false
}

resource "tencentcloud_ckafka_instance" "kafka_instance" {
  instance_name      = "ckafka-instance-prepaid"
  zone_id            = 100007
  period             = 1
  vpc_id             = tencentcloud_vpc.kafka_vpc.id
  subnet_id          = tencentcloud_subnet.kafka_subnet.id
  msg_retention_time = 1200
  max_message_byte   = 1025
  renew_flag         = 0
  kafka_version      = "1.1.1"
  disk_size          = 300
  band_width         = 20
  disk_type          = "CLOUD_BASIC"
  partition          = 400
  elastic_bandwidth_switch = 1

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
`

const testAccKafkaInstancePostpaid = `
resource "tencentcloud_vpc" "kafka_vpc" {
	name       = "postpaid-kafka-vpc"
	cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_subnet" "kafka_subnet" {
	vpc_id            = tencentcloud_vpc.kafka_vpc.id
	name              = "postpaid-kafka-subnet"
	cidr_block        = "10.0.0.0/16"
	availability_zone = "ap-guangzhou-7"
	is_multicast      = false
}

resource "tencentcloud_ckafka_instance" "kafka_instance_postpaid" {
  instance_name      = "ckafka-instance-postpaid"
  zone_id            = 100007
  vpc_id             = tencentcloud_vpc.kafka_vpc.id
  subnet_id          = tencentcloud_subnet.kafka_subnet.id
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

const testAccKafkaInstanceUpdatePostpaid = `
resource "tencentcloud_vpc" "kafka_vpc" {
	name       = "postpaid-kafka-vpc"
	cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_subnet" "kafka_subnet" {
	vpc_id            = tencentcloud_vpc.kafka_vpc.id
	name              = "postpaid-kafka-subnet"
	cidr_block        = "10.0.0.0/16"
	availability_zone = "ap-guangzhou-7"
	is_multicast      = false
}

resource "tencentcloud_ckafka_instance" "kafka_instance_postpaid" {
  instance_name      = "ckafka-instance-postpaid"
  zone_id            = 100007
  vpc_id             = tencentcloud_vpc.kafka_vpc.id
  subnet_id          = tencentcloud_subnet.kafka_subnet.id
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

const testAccKafkaInstanceUpdatePostpaidDiskSize = `
resource "tencentcloud_vpc" "kafka_vpc" {
	name       = "postpaid-kafka-vpc"
	cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_subnet" "kafka_subnet" {
	vpc_id            = tencentcloud_vpc.kafka_vpc.id
	name              = "postpaid-kafka-subnet"
	cidr_block        = "10.0.0.0/16"
	availability_zone = "ap-guangzhou-7"
	is_multicast      = false
}

resource "tencentcloud_ckafka_instance" "kafka_instance_postpaid" {
  instance_name      = "ckafka-instance-postpaid"
  zone_id            = 100007
  vpc_id             = tencentcloud_vpc.kafka_vpc.id
  subnet_id          = tencentcloud_subnet.kafka_subnet.id
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

const testAccKafkaInstanceMAZ = `
resource "tencentcloud_vpc" "kafka_vpc" {
	name       = "maz-kafka-vpc"
	cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_subnet" "kafka_subnet" {
	vpc_id            = tencentcloud_vpc.kafka_vpc.id
	name              = "maz-kafka-subnet"
	cidr_block        = "10.0.0.0/16"
	availability_zone = "ap-guangzhou-7"
	is_multicast      = false
}

resource "tencentcloud_ckafka_instance" "kafka_instance" {
  instance_name      = "ckafka-instance-maz-tf-test"
  zone_id            = 100007
  multi_zone_flag    = true
  zone_ids           = [100007, 100006]
  period             = 1
  vpc_id             = tencentcloud_vpc.kafka_vpc.id
  subnet_id          = tencentcloud_subnet.kafka_subnet.id
  msg_retention_time = 1300
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
`

const testAccKafkaInstanceType = `
resource "tencentcloud_vpc" "vpc" {
  name       = "kafka-type-vpc"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_subnet" "subnet" {
  vpc_id            = tencentcloud_vpc.vpc.id
  name              = "kafka-type-subnet"
  cidr_block        = "10.0.0.0/16"
  availability_zone = "ap-guangzhou-7"
}

resource "tencentcloud_ckafka_instance" "kafka_instance" {
  instance_name      = "ckafka-instance-type-tf-test"
  zone_id            = 100007
  vpc_id             = tencentcloud_vpc.vpc.id
  subnet_id          = tencentcloud_subnet.subnet.id
  msg_retention_time = 1300
  kafka_version      = "1.1.1"
  specifications_type = "profession"
  disk_size          = 200
  disk_type          = "CLOUD_BASIC"
  band_width         = 20
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
`
