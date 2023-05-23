package tencentcloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// go test -i; go test -test.run TestAccTencentCloudSqlserverBusinessIntelligenceInstanceResource_basic -v
func TestAccTencentCloudSqlserverBusinessIntelligenceInstanceResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		CheckDestroy: testAccCheckSqlserverBusinessIntelligenceInstanceDestroy,
		Providers:    testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSqlserverBusinessIntelligenceInstance,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSqlserverBusinessIntelligenceInstanceExists("tencentcloud_sqlserver_business_intelligence_instance.business_intelligence_instance"),
					resource.TestCheckResourceAttrSet("tencentcloud_sqlserver_business_intelligence_instance.business_intelligence_instance", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_sqlserver_business_intelligence_instance.business_intelligence_instance",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccSqlserverBusinessIntelligenceInstanceUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSqlserverBusinessIntelligenceInstanceExists("tencentcloud_sqlserver_business_intelligence_instance.business_intelligence_instance"),
					resource.TestCheckResourceAttrSet("tencentcloud_sqlserver_business_intelligence_instance.business_intelligence_instance", "id"),
				),
			},
		},
	})
}

func testAccCheckSqlserverBusinessIntelligenceInstanceDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_sqlserver_business_intelligence_instance" {
			continue
		}
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)
		service := SqlserverService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}

		instanceId := rs.Primary.ID
		result, err := service.DescribeSqlserverBusinessIntelligenceInstanceById(ctx, instanceId)
		if err != nil {
			return err
		}

		if result != nil {
			return fmt.Errorf("sqlserver business intelligence instance %s still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckSqlserverBusinessIntelligenceInstanceExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource %s is not found", n)
		}

		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)
		service := SqlserverService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}

		instanceId := rs.Primary.ID
		result, err := service.DescribeSqlserverBusinessIntelligenceInstanceById(ctx, instanceId)
		if err != nil {
			return err
		}

		if result == nil {
			return fmt.Errorf("sqlserver business intelligence instance %s is not found", rs.Primary.ID)
		} else {
			return nil
		}
	}
}

const testAccSqlserverBusinessIntelligenceInstance = `
resource "tencentcloud_sqlserver_business_intelligence_instance" "business_intelligence_instance" {
  zone = "ap-guangzhou-6"
  memory = 4
  storage = 20
  cpu = 2
  machine_type = "CLOUD_PREMIUM"
  project_id = 0
  subnet_id = "subnet-dwj7ipnc"
  vpc_id = "vpc-4owdpnwr"
  db_version = "201603"
  security_group_list = []
  weekly = [1, 2, 3, 4, 5, 6, 7]
  start_time = "00:00"
  span = 6
  instance_name = "create_db_name"
}
`

const testAccSqlserverBusinessIntelligenceInstanceUpdate = `
resource "tencentcloud_sqlserver_business_intelligence_instance" "business_intelligence_instance" {
  zone = "ap-guangzhou-6"
  memory = 4
  storage = 20
  cpu = 2
  machine_type = "CLOUD_PREMIUM"
  project_id = 0
  subnet_id = "subnet-dwj7ipnc"
  vpc_id = "vpc-4owdpnwr"
  db_version = "201603"
  security_group_list = []
  weekly = [1, 2, 3, 4, 5, 6, 7]
  start_time = "00:00"
  span = 6
  instance_name = "update_db_name"
}
`
