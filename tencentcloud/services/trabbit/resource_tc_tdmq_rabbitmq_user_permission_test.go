package trabbit_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svctdmq "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tdmq"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// go test -i; go test -test.run TestAccTencentCloudTdmqRabbitmqUserPermissionResource_basic -v
func TestAccTencentCloudTdmqRabbitmqUserPermissionResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckTdmqRabbitmqUserPermissionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTdmqRabbitmqUserPermission,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTdmqRabbitmqUserPermissionExists("tencentcloud_tdmq_rabbitmq_user_permission.example"),
					resource.TestCheckResourceAttrSet("tencentcloud_tdmq_rabbitmq_user_permission.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_tdmq_rabbitmq_user_permission.example", "instance_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_tdmq_rabbitmq_user_permission.example", "user"),
					resource.TestCheckResourceAttrSet("tencentcloud_tdmq_rabbitmq_user_permission.example", "virtual_host"),
					resource.TestCheckResourceAttr("tencentcloud_tdmq_rabbitmq_user_permission.example", "config_regexp", ".*"),
					resource.TestCheckResourceAttr("tencentcloud_tdmq_rabbitmq_user_permission.example", "write_regexp", ".*"),
					resource.TestCheckResourceAttr("tencentcloud_tdmq_rabbitmq_user_permission.example", "read_regexp", ".*"),
				),
			},
			{
				ResourceName:      "tencentcloud_tdmq_rabbitmq_user_permission.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccTdmqRabbitmqUserPermissionUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTdmqRabbitmqUserPermissionExists("tencentcloud_tdmq_rabbitmq_user_permission.example"),
					resource.TestCheckResourceAttrSet("tencentcloud_tdmq_rabbitmq_user_permission.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_tdmq_rabbitmq_user_permission.example", "instance_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_tdmq_rabbitmq_user_permission.example", "user"),
					resource.TestCheckResourceAttrSet("tencentcloud_tdmq_rabbitmq_user_permission.example", "virtual_host"),
					resource.TestCheckResourceAttr("tencentcloud_tdmq_rabbitmq_user_permission.example", "config_regexp", "^tf-.*"),
					resource.TestCheckResourceAttr("tencentcloud_tdmq_rabbitmq_user_permission.example", "write_regexp", "^tf-.*"),
					resource.TestCheckResourceAttr("tencentcloud_tdmq_rabbitmq_user_permission.example", "read_regexp", "^tf-.*"),
				),
			},
		},
	})
}

func testAccCheckTdmqRabbitmqUserPermissionDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	service := svctdmq.NewTdmqService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_tdmq_rabbitmq_user_permission" {
			continue
		}

		idSplit := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
		if len(idSplit) != 3 {
			return fmt.Errorf("id is broken,%s", idSplit)
		}
		instanceId := idSplit[0]
		user := idSplit[1]
		virtualHost := idSplit[2]

		permission, err := service.DescribeTdmqRabbitmqPermissionById(ctx, instanceId, user, virtualHost)
		if permission != nil {
			return fmt.Errorf("tdmq rabbitmqUserPermission %s still exists", rs.Primary.ID)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func testAccCheckTdmqRabbitmqUserPermissionExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}

		idSplit := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
		if len(idSplit) != 3 {
			return fmt.Errorf("id is broken,%s", idSplit)
		}
		instanceId := idSplit[0]
		user := idSplit[1]
		virtualHost := idSplit[2]

		service := svctdmq.NewTdmqService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
		permission, err := service.DescribeTdmqRabbitmqPermissionById(ctx, instanceId, user, virtualHost)
		if permission == nil {
			return fmt.Errorf("tdmq rabbitmqUserPermission %s is not found", rs.Primary.ID)
		}
		if err != nil {
			return err
		}

		return nil
	}
}

const testAccTdmqRabbitmqUserPermission = `
data "tencentcloud_availability_zones" "zones" {
  name = "ap-guangzhou-6"
}

# create vpc
resource "tencentcloud_vpc" "vpc" {
  name       = "vpc"
  cidr_block = "10.0.0.0/16"
}

# create vpc subnet
resource "tencentcloud_subnet" "subnet" {
  name              = "subnet"
  vpc_id            = tencentcloud_vpc.vpc.id
  availability_zone = "ap-guangzhou-6"
  cidr_block        = "10.0.20.0/28"
  is_multicast      = false
}

# create rabbitmq instance
resource "tencentcloud_tdmq_rabbitmq_vip_instance" "example" {
  zone_ids                              = [data.tencentcloud_availability_zones.zones.zones.0.id]
  vpc_id                                = tencentcloud_vpc.vpc.id
  subnet_id                             = tencentcloud_subnet.subnet.id
  cluster_name                          = "tf-example-rabbitmq-permission"
  node_spec                             = "rabbit-vip-basic-1"
  node_num                              = 1
  storage_size                          = 200
  enable_create_default_ha_mirror_queue = false
  auto_renew_flag                       = true
  time_span                             = 1
}

# create rabbitmq user
resource "tencentcloud_tdmq_rabbitmq_user" "example" {
  instance_id = tencentcloud_tdmq_rabbitmq_vip_instance.example.id
  user        = "tf-example-user"
  password    = "Password@123"
  description = "test user"
  tags        = ["management"]
}

# create virtual host
resource "tencentcloud_tdmq_rabbitmq_virtual_host" "example" {
  instance_id  = tencentcloud_tdmq_rabbitmq_vip_instance.example.id
  virtual_host = "tf-example-vhost"
  description  = "test virtual host"
  trace_flag   = false
}

# create user permission
resource "tencentcloud_tdmq_rabbitmq_user_permission" "example" {
  instance_id    = tencentcloud_tdmq_rabbitmq_vip_instance.example.id
  user           = tencentcloud_tdmq_rabbitmq_user.example.user
  virtual_host   = tencentcloud_tdmq_rabbitmq_virtual_host.example.virtual_host
  config_regexp  = ".*"
  write_regexp   = ".*"
  read_regexp    = ".*"
}
`

const testAccTdmqRabbitmqUserPermissionUpdate = `
data "tencentcloud_availability_zones" "zones" {
  name = "ap-guangzhou-6"
}

# create vpc
resource "tencentcloud_vpc" "vpc" {
  name       = "vpc"
  cidr_block = "10.0.0.0/16"
}

# create vpc subnet
resource "tencentcloud_subnet" "subnet" {
  name              = "subnet"
  vpc_id            = tencentcloud_vpc.vpc.id
  availability_zone = "ap-guangzhou-6"
  cidr_block        = "10.0.20.0/28"
  is_multicast      = false
}

# create rabbitmq instance
resource "tencentcloud_tdmq_rabbitmq_vip_instance" "example" {
  zone_ids                              = [data.tencentcloud_availability_zones.zones.zones.0.id]
  vpc_id                                = tencentcloud_vpc.vpc.id
  subnet_id                             = tencentcloud_subnet.subnet.id
  cluster_name                          = "tf-example-rabbitmq-permission"
  node_spec                             = "rabbit-vip-basic-1"
  node_num                              = 1
  storage_size                          = 200
  enable_create_default_ha_mirror_queue = false
  auto_renew_flag                       = true
  time_span                             = 1
}

# create rabbitmq user
resource "tencentcloud_tdmq_rabbitmq_user" "example" {
  instance_id = tencentcloud_tdmq_rabbitmq_vip_instance.example.id
  user        = "tf-example-user"
  password    = "Password@123"
  description = "test user"
  tags        = ["management"]
}

# create virtual host
resource "tencentcloud_tdmq_rabbitmq_virtual_host" "example" {
  instance_id  = tencentcloud_tdmq_rabbitmq_vip_instance.example.id
  virtual_host = "tf-example-vhost"
  description  = "test virtual host"
  trace_flag   = false
}

# create user permission (updated)
resource "tencentcloud_tdmq_rabbitmq_user_permission" "example" {
  instance_id    = tencentcloud_tdmq_rabbitmq_vip_instance.example.id
  user           = tencentcloud_tdmq_rabbitmq_user.example.user
  virtual_host   = tencentcloud_tdmq_rabbitmq_virtual_host.example.virtual_host
  config_regexp  = "^tf-.*"
  write_regexp   = "^tf-.*"
  read_regexp    = "^tf-.*"
}
`
