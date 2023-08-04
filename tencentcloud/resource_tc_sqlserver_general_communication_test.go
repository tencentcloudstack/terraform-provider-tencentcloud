package tencentcloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// go test -i; go test -test.run TestAccTencentCloudSqlserverGeneralCommunicationResource_basic -v
func TestAccTencentCloudSqlserverGeneralCommunicationResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		CheckDestroy: testAccCheckSqlserverGeneralCommunicationDestroy,
		Providers:    testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSqlserverGeneralCommunication,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSqlserverGeneralCommunicationExists("tencentcloud_sqlserver_general_communication.general_communication"),
					resource.TestCheckResourceAttrSet("tencentcloud_sqlserver_general_communication.general_communication", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_sqlserver_general_communication.general_communication",
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
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)
		service := SqlserverService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}

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

		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)
		service := SqlserverService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}

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

const testAccSqlserverGeneralCommunication = `
data "tencentcloud_availability_zones_by_product" "zones" {
  product = "sqlserver"
}

resource "tencentcloud_vpc" "vpc" {
  name       = "vpc-example"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_subnet" "subnet" {
  availability_zone = data.tencentcloud_availability_zones_by_product.zones.zones.4.name
  name              = "subnet-example"
  vpc_id            = tencentcloud_vpc.vpc.id
  cidr_block        = "10.0.0.0/16"
  is_multicast      = false
}

resource "tencentcloud_security_group" "security_group" {
  name        = "sg-example"
  description = "desc."
}

resource "tencentcloud_sqlserver_basic_instance" "example" {
  name                   = "tf-example"
  availability_zone      = data.tencentcloud_availability_zones_by_product.zones.zones.4.name
  charge_type            = "POSTPAID_BY_HOUR"
  vpc_id                 = tencentcloud_vpc.vpc.id
  subnet_id              = tencentcloud_subnet.subnet.id
  project_id             = 0
  memory                 = 4
  storage                = 100
  cpu                    = 2
  machine_type           = "CLOUD_PREMIUM"
  maintenance_week_set   = [1, 2, 3]
  maintenance_start_time = "09:00"
  maintenance_time_span  = 3
  security_groups        = [tencentcloud_security_group.security_group.id]

  tags = {
    "test" = "test"
  }
}

resource "tencentcloud_sqlserver_general_communication" "example" {
  instance_id = tencentcloud_sqlserver_basic_instance.example.id
}
`
