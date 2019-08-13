package tencentcloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
)

func TestAccTencentCloudAsAttachment(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAsAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAsAttachment(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAsAttachmentExists("tencentcloud_as_attachment.attachment"),
					resource.TestCheckResourceAttrSet("tencentcloud_as_attachment.attachment", "scaling_group_id"),
					resource.TestCheckResourceAttr("tencentcloud_as_attachment.attachment", "instance_ids.#", "1"),
				),
			},
			// test update case
			{
				Config: testAccAsAttachment_update(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAsAttachmentExists("tencentcloud_as_attachment.attachment"),
					resource.TestCheckResourceAttrSet("tencentcloud_as_attachment.attachment", "scaling_group_id"),
					resource.TestCheckResourceAttr("tencentcloud_as_attachment.attachment", "instance_ids.#", "2"),
				),
			},
		},
	})
}

func testAccCheckAsAttachmentExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), "logId", logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("auto scaling attachment %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("auto scaling attachment id is not set")
		}
		asService := AsService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
		_, err := asService.DescribeAutoScalingAttachment(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}
		return nil
	}
}

func testAccCheckAsAttachmentDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	asService := AsService{
		client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_as_attachment" {
			continue
		}

		instances, err := asService.DescribeAutoScalingAttachment(ctx, rs.Primary.ID)
		if err != nil {
			if sdkErr, ok := err.(*errors.TencentCloudSDKError); ok {
				if sdkErr.Code == AsScalingGroupNotFound {
					return nil
				}
			}
			return err
		}
		if len(instances) > 0 {
			return fmt.Errorf("auto scaling attachment still exists: %s", rs.Primary.ID)
		}
	}
	return nil
}

func testAccAsAttachment() string {
	return `
resource "tencentcloud_vpc" "vpc" {
  name       = "tf-as-vpc"
  cidr_block = "10.2.0.0/16"
}

resource "tencentcloud_subnet" "subnet" {
  vpc_id            = "${tencentcloud_vpc.vpc.id}"
  name              = "tf-as-subnet"
  cidr_block        = "10.2.11.0/24"
  availability_zone = "ap-guangzhou-3"
}

resource "tencentcloud_as_scaling_config" "launch_configuration" {
  configuration_name = "tf-as-attachment-config"
  image_id           = "img-9qabwvbn"
  instance_types     = ["SA1.SMALL1"]
}

resource "tencentcloud_as_scaling_group" "scaling_group" {
  scaling_group_name = "tf-as-attachment-group"
  configuration_id   = "${tencentcloud_as_scaling_config.launch_configuration.id}"
  max_size           = 5
  min_size           = 0
  vpc_id             = "${tencentcloud_vpc.vpc.id}"
  subnet_ids         = ["${tencentcloud_subnet.subnet.id}"]
}

resource "tencentcloud_instance" "cvm_instance" {
  instance_name     = "tf_as_instance"
  availability_zone = "ap-guangzhou-3"
  image_id          = "img-9qabwvbn"
  instance_type     = "SA1.SMALL1"
  system_disk_type  = "CLOUD_SSD"
  vpc_id            = "${tencentcloud_vpc.vpc.id}"
  subnet_id         = "${tencentcloud_subnet.subnet.id}"
}

resource "tencentcloud_as_attachment" "attachment" {
  scaling_group_id = "${tencentcloud_as_scaling_group.scaling_group.id}"
  instance_ids     = ["${tencentcloud_instance.cvm_instance.id}"]
}
`
}

func testAccAsAttachment_update() string {
	return `
resource "tencentcloud_vpc" "vpc" {
  name       = "tf-as-vpc"
  cidr_block = "10.2.0.0/16"
}

resource "tencentcloud_subnet" "subnet" {
  vpc_id            = "${tencentcloud_vpc.vpc.id}"
  name              = "tf-as-subnet"
  cidr_block        = "10.2.11.0/24"
  availability_zone = "ap-guangzhou-3"
}

resource "tencentcloud_as_scaling_config" "launch_configuration" {
  configuration_name = "tf-as-attachment-config"
  image_id           = "img-9qabwvbn"
  instance_types     = ["SA1.SMALL1"]
}

resource "tencentcloud_as_scaling_group" "scaling_group" {
  scaling_group_name = "tf-as-attachment-group"
  configuration_id   = "${tencentcloud_as_scaling_config.launch_configuration.id}"
  max_size           = 5
  min_size           = 0
  vpc_id             = "${tencentcloud_vpc.vpc.id}"
  subnet_ids         = ["${tencentcloud_subnet.subnet.id}"]
}

resource "tencentcloud_instance" "cvm_instance" {
  instance_name     = "tf_as_instance"
  availability_zone = "ap-guangzhou-3"
  image_id          = "img-9qabwvbn"
  instance_type     = "SA1.SMALL1"
  system_disk_type  = "CLOUD_SSD"
  vpc_id            = "${tencentcloud_vpc.vpc.id}"
  subnet_id         = "${tencentcloud_subnet.subnet.id}"
}

resource "tencentcloud_instance" "cvm_instance_1" {
  instance_name     = "tf_as_instance_1"
  availability_zone = "ap-guangzhou-3"
  image_id          = "img-9qabwvbn"
  instance_type     = "SA1.SMALL1"
  system_disk_type  = "CLOUD_SSD"
  vpc_id            = "${tencentcloud_vpc.vpc.id}"
  subnet_id         = "${tencentcloud_subnet.subnet.id}"
}

resource "tencentcloud_as_attachment" "attachment" {
  scaling_group_id = "${tencentcloud_as_scaling_group.scaling_group.id}"
  instance_ids     = ["${tencentcloud_instance.cvm_instance.id}", "${tencentcloud_instance.cvm_instance_1.id}"]
}
`
}
