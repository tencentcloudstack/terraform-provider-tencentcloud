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

// go test -i; go test -test.run TestAccTencentCloudNeedFixTdmqRabbitmqVirtualHostResource_basic -v
func TestAccTencentCloudNeedFixTdmqRabbitmqVirtualHostResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PREPAY)
		},
		CheckDestroy: testAccCheckTdmqRabbitmqVirtualHostDestroy,
		Providers:    tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTdmqRabbitmqVirtualHost,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTdmqRabbitmqVirtualHostExists("tencentcloud_tdmq_rabbitmq_virtual_host.rabbitmq_virtual_host"),
					resource.TestCheckResourceAttrSet("tencentcloud_tdmq_rabbitmq_virtual_host.rabbitmq_virtual_host", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_tdmq_rabbitmq_virtual_host.rabbitmq_virtual_host", "instance_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_tdmq_rabbitmq_virtual_host.rabbitmq_virtual_host", "virtual_host"),
					resource.TestCheckResourceAttrSet("tencentcloud_tdmq_rabbitmq_virtual_host.rabbitmq_virtual_host", "description"),
					resource.TestCheckResourceAttrSet("tencentcloud_tdmq_rabbitmq_virtual_host.rabbitmq_virtual_host", "trace_flag"),
				),
			},
			{
				Config: testAccTdmqRabbitmqVirtualHostUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTdmqRabbitmqVirtualHostExists("tencentcloud_tdmq_rabbitmq_virtual_host.rabbitmq_virtual_host"),
					resource.TestCheckResourceAttrSet("tencentcloud_tdmq_rabbitmq_virtual_host.rabbitmq_virtual_host", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_tdmq_rabbitmq_virtual_host.rabbitmq_virtual_host", "instance_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_tdmq_rabbitmq_virtual_host.rabbitmq_virtual_host", "virtual_host"),
					resource.TestCheckResourceAttrSet("tencentcloud_tdmq_rabbitmq_virtual_host.rabbitmq_virtual_host", "description"),
					resource.TestCheckResourceAttrSet("tencentcloud_tdmq_rabbitmq_virtual_host.rabbitmq_virtual_host", "trace_flag"),
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
resource "tencentcloud_tdmq_rabbitmq_virtual_host" "rabbitmq_virtual_host" {
  instance_id  = "amqp-kzbe8p3n"
  virtual_host = "vh-test-1"
  description  = "desc"
  trace_flag   = false
}
`

const testAccTdmqRabbitmqVirtualHostUpdate = `
resource "tencentcloud_tdmq_rabbitmq_virtual_host" "rabbitmq_virtual_host" {
  instance_id  = "amqp-kzbe8p3n"
  virtual_host = "vh-test-1"
  description  = "desc update"
  trace_flag   = true
}
`
