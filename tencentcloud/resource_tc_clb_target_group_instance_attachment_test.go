package tencentcloud

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccTencentCloudClbTGAttachmentInstance_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckClbTGAttachmentInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccClbTGAttachmentInstance_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClbTGAttachmentInstanceExists("tencentcloud_clb_target_group_instance_attachment.test"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_target_group_instance_attachment.test", "target_group_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_target_group_instance_attachment.test", "bind_ip"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_target_group_instance_attachment.test", "port"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_target_group_instance_attachment.test", "weight"),
				),
			},
			{
				Config: testAccClbTGAttachmentInstance_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClbTGAttachmentInstanceExists("tencentcloud_clb_target_group_instance_attachment.test"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_target_group_instance_attachment.test", "target_group_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_target_group_instance_attachment.test", "bind_ip"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_target_group_instance_attachment.test", "port"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_target_group_instance_attachment.test", "weight"),
				),
			},
		},
	})
}

func testAccCheckClbTGAttachmentInstanceDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	clbService := ClbService{
		client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_clb_target_group_instance_attachment" {
			continue
		}
		time.Sleep(5 * time.Second)
		idSplit := strings.Split(rs.Primary.ID, FILED_SP)
		if len(idSplit) != 3 {
			return fmt.Errorf("target group instance attachment id is not set")
		}
		targetGroupId := idSplit[0]
		bindIp := idSplit[1]
		port, err := strconv.ParseUint(idSplit[2], 0, 64)
		if err != nil {
			return err
		}

		filters := make(map[string]string)
		filters["TargetGroupId"] = targetGroupId
		filters["BindIP"] = bindIp
		targetGroupInstances, err := clbService.DescribeTargetGroupInstances(ctx, filters)
		if err != nil {
			return err
		}
		for _, tgInstance := range targetGroupInstances {
			if *tgInstance.Port == port {
				return fmt.Errorf("[CHECK][CLB target group instance attachment][Destroy] check: CLB target group instance attachment still exists: %s", rs.Primary.ID)
			}
		}
	}
	return nil
}

func testAccCheckClbTGAttachmentInstanceExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("[CHECK][CLB target group instance attachment][Exists] check: CLB target group %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("[CHECK][CLB target group instance attachment][Exists] check: CLB target group id is not set")
		}
		clbService := ClbService{
			client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
		}
		idSplit := strings.Split(rs.Primary.ID, FILED_SP)
		if len(idSplit) != 3 {
			return fmt.Errorf("target group instance attachment id is not set")
		}
		targetGroupId := idSplit[0]
		bindIp := idSplit[1]
		port, err := strconv.ParseUint(idSplit[2], 0, 64)
		if err != nil {
			return err
		}

		filters := make(map[string]string)
		filters["TargetGroupId"] = targetGroupId
		filters["BindIP"] = bindIp
		targetGroupInstances, err := clbService.DescribeTargetGroupInstances(ctx, filters)
		if err != nil {
			return err
		}
		for _, tgInstance := range targetGroupInstances {
			if *tgInstance.Port == port {
				return nil
			}
		}
		return fmt.Errorf("[CHECK][CLB target group instance attachment][Exists] id %s is not exist", rs.Primary.ID)
	}
}

const testAccClbTGAttachmentInstance_basic = `
data "tencentcloud_images" "my_favorite_image" {
  image_type = ["PUBLIC_IMAGE"]
  os_name    = "centos"
}

data "tencentcloud_instance_types" "my_favorite_instance_types" {
  filter {
    name   = "instance-family"
    values = ["S3"]
  }

  cpu_core_count = 1
  memory_size    = 1
}

data "tencentcloud_availability_zones" "default" {
}

resource "tencentcloud_vpc" "app" {
  cidr_block = "10.0.0.0/16"
  name       = "awesome_app_vpc"
}

resource "tencentcloud_subnet" "app" {
  vpc_id            = tencentcloud_vpc.app.id
  availability_zone = data.tencentcloud_availability_zones.default.zones.0.name
  name              = "awesome_app_subnet"
  cidr_block        = "10.0.1.0/24"
}

resource "tencentcloud_instance" "my_awesome_app" {
  instance_name              = "awesome_app"
  availability_zone          = data.tencentcloud_availability_zones.default.zones.0.name
  image_id                   = data.tencentcloud_images.my_favorite_image.images.0.image_id
  instance_type              = data.tencentcloud_instance_types.my_favorite_instance_types.instance_types.0.instance_type
  system_disk_type           = "CLOUD_PREMIUM"
  system_disk_size           = 50
  hostname                   = "user"
  project_id                 = 0
  vpc_id                     = tencentcloud_vpc.app.id
  subnet_id                  = tencentcloud_subnet.app.id
  internet_max_bandwidth_out = 20

  data_disks {
    data_disk_type = "CLOUD_PREMIUM"
    data_disk_size = 50
	encrypt        = false
  }

  tags = {
    tagKey = "tagValue"
  }
}

data "tencentcloud_instances" "foo" {
  instance_id = tencentcloud_instance.my_awesome_app.id
}

resource "tencentcloud_clb_target_group" "test"{
    target_group_name = "test"
    vpc_id            = tencentcloud_vpc.app.id
}

resource "tencentcloud_clb_target_group_instance_attachment" "test"{
    target_group_id = tencentcloud_clb_target_group.test.id
    bind_ip         = data.tencentcloud_instances.foo.instance_list[0].private_ip 
    port            = 88
    weight          = 3
}
`

const testAccClbTGAttachmentInstance_update = `
data "tencentcloud_images" "my_favorite_image" {
  image_type = ["PUBLIC_IMAGE"]
  os_name    = "centos"
}

data "tencentcloud_instance_types" "my_favorite_instance_types" {
  filter {
    name   = "instance-family"
    values = ["S3"]
  }

  cpu_core_count = 1
  memory_size    = 1
}

data "tencentcloud_availability_zones" "default" {
}

resource "tencentcloud_vpc" "app" {
  cidr_block = "10.0.0.0/16"
  name       = "awesome_app_vpc"
}

resource "tencentcloud_subnet" "app" {
  vpc_id            = tencentcloud_vpc.app.id
  availability_zone = data.tencentcloud_availability_zones.default.zones.0.name
  name              = "awesome_app_subnet"
  cidr_block        = "10.0.1.0/24"
}

resource "tencentcloud_instance" "my_awesome_app" {
  instance_name              = "awesome_app"
  availability_zone          = data.tencentcloud_availability_zones.default.zones.0.name
  image_id                   = data.tencentcloud_images.my_favorite_image.images.0.image_id
  instance_type              = data.tencentcloud_instance_types.my_favorite_instance_types.instance_types.0.instance_type
  system_disk_type           = "CLOUD_PREMIUM"
  system_disk_size           = 50
  hostname                   = "user"
  project_id                 = 0
  vpc_id                     = tencentcloud_vpc.app.id
  subnet_id                  = tencentcloud_subnet.app.id
  internet_max_bandwidth_out = 20

  data_disks {
    data_disk_type = "CLOUD_PREMIUM"
    data_disk_size = 50
	encrypt        = false
  }

  tags = {
    tagKey = "tagValue"
  }
}

data "tencentcloud_instances" "foo" {
  instance_id = tencentcloud_instance.my_awesome_app.id
}

resource "tencentcloud_clb_target_group" "test"{
    target_group_name = "test"
    vpc_id            = tencentcloud_vpc.app.id
}

resource "tencentcloud_clb_target_group_instance_attachment" "test"{
    target_group_id = tencentcloud_clb_target_group.test.id
    bind_ip         = data.tencentcloud_instances.foo.instance_list[0].private_ip 
    port            = 88
    weight          = 5
}
`
