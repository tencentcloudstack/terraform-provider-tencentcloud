package tencentcloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccTencentCloudCkafkaInstance(t *testing.T) {
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
					resource.TestCheckResourceAttr("tencentcloud_ckafka_instance.kafka_instance", "instance_name", "ckafka-instance-tf-test"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_instance.kafka_instance", "zone_id", "100003"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_instance.kafka_instance", "period", "1"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_instance.kafka_instance", "msg_retention_time", "1300"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_instance.kafka_instance", "renew_flag", "0"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_instance.kafka_instance", "kafka_version", "1.1.1"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_instance.kafka_instance", "disk_size", "500"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_instance.kafka_instance", "disk_type", "CLOUD_BASIC"),
					resource.TestCheckResourceAttrSet("tencentcloud_ckafka_instance.kafka_instance", "vip"),
					resource.TestCheckResourceAttrSet("tencentcloud_ckafka_instance.kafka_instance", "vport"),
				),
			},
			{
				Config: testAccKafkaInstanceUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKafkaInstanceExists("tencentcloud_ckafka_instance.kafka_instance"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_instance.kafka_instance", "instance_name", "ckafka-instance-tf-test"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_instance.kafka_instance", "zone_id", "100003"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_instance.kafka_instance", "period", "1"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_instance.kafka_instance", "msg_retention_time", "1200"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_instance.kafka_instance", "renew_flag", "0"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_instance.kafka_instance", "kafka_version", "1.1.1"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_instance.kafka_instance", "disk_size", "500"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_instance.kafka_instance", "disk_type", "CLOUD_BASIC"),
				),
			},
			{
				ResourceName:            "tencentcloud_ckafka_instance.kafka_instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"period"},
			},
		},
	})
}

func TestAccTencentCloudKafkaInstanceMAZ(t *testing.T) {
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
					resource.TestCheckResourceAttr("tencentcloud_ckafka_instance.kafka_instance", "zone_id", "100003"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_instance.kafka_instance", "zone_ids.#", "2"),
				),
			},
			{
				ResourceName:            "tencentcloud_ckafka_instance.kafka_instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"period"},
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
  instance_name      = "ckafka-instance-tf-test"
  zone_id            = 100003
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

const testAccKafkaInstanceUpdate = defaultKafkaVariable + `
resource "tencentcloud_ckafka_instance" "kafka_instance" {
  instance_name      = "ckafka-instance-tf-test"
  zone_id            = 100003
  period             = 1
  vpc_id             =  var.vpc_id
  subnet_id          =  var.subnet_id
  msg_retention_time = 1200
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

  tag_set = {
    createdBy = "terraform"
  }
}
`

const testAccKafkaInstanceMAZ = defaultKafkaVariable + `
resource "tencentcloud_ckafka_instance" "kafka_instance" {
  instance_name      = "ckafka-instance-maz-tf-test"
  zone_id            = 100003
  multi_zone_flag    = true
  zone_ids           = [100003, 100006]
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
