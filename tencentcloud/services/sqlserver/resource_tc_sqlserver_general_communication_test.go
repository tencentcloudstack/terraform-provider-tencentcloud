package sqlserver_test

import (
	"context"
	"fmt"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svcsqlserver "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/sqlserver"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// go test -i; go test -test.run TestAccTencentCloudSqlserverGeneralCommunicationResource_basic -v
func TestAccTencentCloudSqlserverGeneralCommunicationResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		CheckDestroy: testAccCheckSqlserverGeneralCommunicationDestroy,
		Providers:    tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSqlserverGeneralCommunication,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSqlserverGeneralCommunicationExists("tencentcloud_sqlserver_general_communication.example"),
					resource.TestCheckResourceAttrSet("tencentcloud_sqlserver_general_communication.example", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_sqlserver_general_communication.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckSqlserverGeneralCommunicationDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_sqlserver_general_communication" {
			continue
		}
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service := svcsqlserver.NewSqlserverService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())

		instanceId := rs.Primary.ID
		result, err := service.DescribeSqlserverGeneralCommunicationById(ctx, instanceId)
		if err != nil {
			return err
		}

		if result != nil {
			return fmt.Errorf("sqlserver general communicationinstance %s still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckSqlserverGeneralCommunicationExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource %s is not found", n)
		}

		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service := svcsqlserver.NewSqlserverService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())

		instanceId := rs.Primary.ID
		result, err := service.DescribeSqlserverGeneralCommunicationById(ctx, instanceId)
		if err != nil {
			return err
		}

		if result == nil {
			return fmt.Errorf("sqlserver general communicationinstance %s is not found", rs.Primary.ID)
		} else {
			return nil
		}
	}
}

const testAccSqlserverGeneralCommunication = tcacctest.DefaultVpcSubnets + tcacctest.DefaultSecurityGroupData + `
data "tencentcloud_availability_zones_by_product" "zones" {
  product = "sqlserver"
}

resource "tencentcloud_sqlserver_basic_instance" "example" {
  name                   = "tf-example"
  availability_zone      = data.tencentcloud_availability_zones_by_product.zones.zones.4.name
  charge_type            = "POSTPAID_BY_HOUR"
  vpc_id                 = local.vpc_id
  subnet_id              = local.subnet_id
  project_id             = 0
  memory                 = 4
  storage                = 100
  cpu                    = 2
  machine_type           = "CLOUD_PREMIUM"
  maintenance_week_set   = [1, 2, 3]
  maintenance_start_time = "09:00"
  maintenance_time_span  = 3
  security_groups        = [local.sg_id]

  tags = {
    "test" = "test"
  }
}

resource "tencentcloud_sqlserver_general_communication" "example" {
  instance_id = tencentcloud_sqlserver_basic_instance.example.id
}
`
