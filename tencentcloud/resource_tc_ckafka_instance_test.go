package tencentcloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccTencentCloudCkafkaInstanceResource_prepaid(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY) },
		Providers:    testAccProviders,
		CheckDestroy: testAccTencentCloudKafkaInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKafkaInstance,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKafkaInstanceExists("tencentcloud_ckafka_instance.kafka_instance"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_instance.kafka_instance", "instance_name", "ckafka-instance-prepaid"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_instance.kafka_instance", "zone_id", "100007"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_instance.kafka_instance", "period", "1"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_instance.kafka_instance", "msg_retention_time", "1300"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_instance.kafka_instance", "max_message_byte", "1024"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_instance.kafka_instance", "renew_flag", "0"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_instance.kafka_instance", "kafka_version", "1.1.1"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_instance.kafka_instance", "disk_size", "200"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_instance.kafka_instance", "disk_type", "CLOUD_BASIC"),
					resource.TestCheckResourceAttrSet("tencentcloud_ckafka_instance.kafka_instance", "vip"),
					resource.TestCheckResourceAttrSet("tencentcloud_ckafka_instance.kafka_instance", "vport"),
				),
			},
			{
				Config: testAccKafkaInstanceUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKafkaInstanceExists("tencentcloud_ckafka_instance.kafka_instance"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_instance.kafka_instance", "instance_name", "ckafka-instance-prepaid"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_instance.kafka_instance", "zone_id", "100007"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_instance.kafka_instance", "period", "1"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_instance.kafka_instance", "msg_retention_time", "1200"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_instance.kafka_instance", "max_message_byte", "1025"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_instance.kafka_instance", "renew_flag", "0"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_instance.kafka_instance", "kafka_version", "1.1.1"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_instance.kafka_instance", "disk_size", "200"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_instance.kafka_instance", "disk_type", "CLOUD_BASIC"),
				),
			},
			{
				ResourceName:            "tencentcloud_ckafka_instance.kafka_instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"period", "max_message_byte", "charge_type"},
			},
		},
	})
}

func TestAccTencentCloudCkafkaInstanceResource_postpaid(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY) },
		Providers:    testAccProviders,
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
				ImportStateVerifyIgnore: []string{"period", "max_message_byte", "charge_type", "upgrade_strategy"},
			},
		},
	})
}

func TestAccTencentCloudCkafkaInstanceResource_maz(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY) },
		Providers:    testAccProviders,
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
				ResourceName:            "tencentcloud_ckafka_instance.kafka_instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"period", "max_message_byte", "charge_type"},
			},
		},
	})
}

func TestAccTencentCloudCkafkaInstanceResource_type(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY) },
		Providers:    testAccProviders,
		CheckDestroy: testAccTencentCloudKafkaInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKafkaInstanceType,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKafkaInstanceExists("tencentcloud_ckafka_instance.kafka_instance"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_instance.kafka_instance", "instance_name", "ckafka-instance-type-tf-test"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_instance.kafka_instance", "specifications_type", "standard"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_instance.kafka_instance", "instance_type", "2"),
				),
			},
			{
				ResourceName:            "tencentcloud_ckafka_instance.kafka_instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"period", "max_message_byte", "charge_type"},
			},
		},
	})
}

func testAccTencentCloudKafkaInstanceDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	ckafkcService := CkafkaService{
		client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
	}
	for _, r := range s.RootModule().Resources {
		if r.Type != "tencentcloud_ckafka_instance" {
			continue
		}
		_, has, error := ckafkcService.DescribeInstanceById(ctx, r.Primary.ID)
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
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("ckafka instance %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("ckafka instance id is not set")
		}
		ckafkcService := CkafkaService{
			client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
		}

		var exist bool
		outErr := resource.Retry(readRetryTimeout, func() *resource.RetryError {
			_, has, inErr := ckafkcService.DescribeInstanceById(ctx, rs.Primary.ID)
			if inErr != nil {
				return retryError(inErr)
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

const testAccKafkaInstance = defaultKafkaVariable + `
resource "tencentcloud_ckafka_instance" "kafka_instance" {
  instance_name      = "ckafka-instance-prepaid"
  zone_id            = 100007
  period             = 1
  vpc_id             = var.vpc_id
  subnet_id          = var.subnet_id
  msg_retention_time = 1300
  max_message_byte   = 1024
  renew_flag         = 0
  kafka_version      = "1.1.1"
  disk_size          = 200
  band_width         = 20
  disk_type          = "CLOUD_BASIC"
  partition          = 400


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

const testAccKafkaInstanceUpdate = defaultKafkaVariable + `
resource "tencentcloud_ckafka_instance" "kafka_instance" {
  instance_name      = "ckafka-instance-prepaid"
  zone_id            = 100007
  period             = 1
  vpc_id             =  var.vpc_id
  subnet_id          =  var.subnet_id
  msg_retention_time = 1200
  max_message_byte   = 1025
  renew_flag         = 0
  kafka_version      = "1.1.1"
  disk_size          = 200
  band_width         = 20
  disk_type          = "CLOUD_BASIC"
  partition          = 400


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

const testAccKafkaInstancePostpaid = defaultKafkaVariable + `
resource "tencentcloud_ckafka_instance" "kafka_instance_postpaid" {
  instance_name      = "ckafka-instance-postpaid"
  zone_id            = 100007
  vpc_id             = var.vpc_id
  subnet_id          = var.subnet_id
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

const testAccKafkaInstanceUpdatePostpaid = defaultKafkaVariable + `
resource "tencentcloud_ckafka_instance" "kafka_instance_postpaid" {
  instance_name      = "ckafka-instance-postpaid"
  zone_id            = 100007
  vpc_id             =  var.vpc_id
  subnet_id          =  var.subnet_id
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

const testAccKafkaInstanceUpdatePostpaidDiskSize = defaultKafkaVariable + `
resource "tencentcloud_ckafka_instance" "kafka_instance_postpaid" {
  instance_name      = "ckafka-instance-postpaid"
  zone_id            = 100007
  vpc_id             =  var.vpc_id
  subnet_id          =  var.subnet_id
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

const testAccKafkaInstanceMAZ = defaultKafkaVariable + `
resource "tencentcloud_ckafka_instance" "kafka_instance" {
  instance_name      = "ckafka-instance-maz-tf-test"
  zone_id            = 100007
  multi_zone_flag    = true
  zone_ids           = [100007, 100006]
  period             = 1
  vpc_id             = var.vpc_id
  subnet_id          = var.subnet_id
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
  name       = "tmp"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_subnet" "subnet" {
  vpc_id            = tencentcloud_vpc.vpc.id
  name              = "subnet-example"
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
  specifications_type = "standard"
  instance_type       = 2
  disk_size          = 1000
  disk_type          = "CLOUD_BASIC"
  band_width         = 100
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
