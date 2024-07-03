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

// go test -i; go test -test.run TestAccTencentCloudTdmqRabbitmqVirtualHostResource_basic -v
func TestAccTencentCloudTdmqRabbitmqVirtualHostResource_basic(t *testing.T) {
	//t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		CheckDestroy: testAccCheckTdmqRabbitmqVirtualHostDestroy,
		Providers:    tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTdmqRabbitmqVirtualHost,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTdmqRabbitmqVirtualHostExists("tencentcloud_tdmq_rabbitmq_virtual_host.example"),
					resource.TestCheckResourceAttrSet("tencentcloud_tdmq_rabbitmq_virtual_host.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_tdmq_rabbitmq_virtual_host.example", "instance_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_tdmq_rabbitmq_virtual_host.example", "virtual_host"),
					resource.TestCheckResourceAttrSet("tencentcloud_tdmq_rabbitmq_virtual_host.example", "description"),
					resource.TestCheckResourceAttrSet("tencentcloud_tdmq_rabbitmq_virtual_host.example", "trace_flag"),
				),
			},
			{
				ResourceName:      "tencentcloud_tdmq_rabbitmq_virtual_host.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccTdmqRabbitmqVirtualHostUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTdmqRabbitmqVirtualHostExists("tencentcloud_tdmq_rabbitmq_virtual_host.example"),
					resource.TestCheckResourceAttrSet("tencentcloud_tdmq_rabbitmq_virtual_host.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_tdmq_rabbitmq_virtual_host.example", "instance_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_tdmq_rabbitmq_virtual_host.example", "virtual_host"),
					resource.TestCheckResourceAttrSet("tencentcloud_tdmq_rabbitmq_virtual_host.example", "description"),
					resource.TestCheckResourceAttrSet("tencentcloud_tdmq_rabbitmq_virtual_host.example", "trace_flag"),
				),
			},
		},
	})
}

func testAccCheckTdmqRabbitmqVirtualHostDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := svctdmq.NewTdmqService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_tdmq_rabbitmq_virtual_host" {
			continue
		}

		ids := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
		if len(ids) != 2 {
			return fmt.Errorf("invalid ID %s", rs.Primary.ID)
		}
		instanceId := ids[0]
		virtualHost := ids[1]

		ret, err := service.DescribeTdmqRabbitmqVirtualHostById(ctx, instanceId, virtualHost)
		if err != nil {
			return err
		}

		if ret != nil {
			return fmt.Errorf("tdmq rabbitmq virtual_host still exist, id: %v", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckTdmqRabbitmqVirtualHostExists(re string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[re]
		if !ok {
			return fmt.Errorf("tdcpg instance  %s is not found", re)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("tdcpg instance id is not set")
		}

		service := svctdmq.NewTdmqService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
		ids := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
		if len(ids) != 2 {
			return fmt.Errorf("invalid ID %s", rs.Primary.ID)
		}
		instanceId := ids[0]
		virtualHost := ids[1]

		ret, err := service.DescribeTdmqRabbitmqVirtualHostById(ctx, instanceId, virtualHost)
		if err != nil {
			return err
		}

		if ret == nil {
			return fmt.Errorf("tdmq rabbitmq virtual_host not found, id: %v", rs.Primary.ID)
		}

		return nil
	}
}

const testAccTdmqRabbitmqVirtualHost = `
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
  cluster_name                          = "tf-example-rabbitmq-vip-instance"
  node_spec                             = "rabbit-vip-basic-1"
  node_num                              = 1
  storage_size                          = 200
  enable_create_default_ha_mirror_queue = false
  auto_renew_flag                       = true
  time_span                             = 1
}

# create virtual host
resource "tencentcloud_tdmq_rabbitmq_virtual_host" "example" {
  instance_id  = tencentcloud_tdmq_rabbitmq_vip_instance.example.id
  virtual_host = "tf-example-vhost"
  description  = "desc."
  trace_flag   = true
}
`

const testAccTdmqRabbitmqVirtualHostUpdate = `
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
  cluster_name                          = "tf-example-rabbitmq-vip-instance"
  node_spec                             = "rabbit-vip-basic-1"
  node_num                              = 1
  storage_size                          = 200
  enable_create_default_ha_mirror_queue = false
  auto_renew_flag                       = true
  time_span                             = 1
}

# create virtual host
resource "tencentcloud_tdmq_rabbitmq_virtual_host" "example" {
  instance_id  = tencentcloud_tdmq_rabbitmq_vip_instance.example.id
  virtual_host = "tf-example-vhost"
  description  = "desc update."
  trace_flag   = false
}
`
