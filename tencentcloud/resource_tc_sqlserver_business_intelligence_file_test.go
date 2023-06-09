package tencentcloud

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// go test -i; go test -test.run TestAccTencentCloudSqlserverBusinessIntelligenceFileResource_basic -v
func TestAccTencentCloudSqlserverBusinessIntelligenceFileResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		CheckDestroy: testAccCheckSqlserverBusinessIntelligenceFileDestroy,
		Providers:    testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSqlserverBusinessIntelligenceFile,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSqlserverBusinessIntelligenceFileExists("tencentcloud_sqlserver_business_intelligence_file.business_intelligence_file"),
					resource.TestCheckResourceAttrSet("tencentcloud_sqlserver_business_intelligence_file.business_intelligence_file", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_sqlserver_business_intelligence_file.business_intelligence_file",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckSqlserverBusinessIntelligenceFileDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_sqlserver_business_intelligence_file" {
			continue
		}
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)
		service := SqlserverService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}

		idSplit := strings.Split(rs.Primary.ID, FILED_SP)
		if len(idSplit) != 2 {
			return fmt.Errorf("id is broken, id is %s", rs.Primary.ID)
		}

		instanceId := idSplit[0]
		fileName := idSplit[1]

		result, err := service.DescribeSqlserverBusinessIntelligenceFileById(ctx, instanceId, fileName)
		if err != nil {
			return err
		}

		if result != nil {
			return fmt.Errorf("sqlserver businessIntelligenceFile %s still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckSqlserverBusinessIntelligenceFileExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource %s is not found", n)
		}

		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)
		service := SqlserverService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}

		idSplit := strings.Split(rs.Primary.ID, FILED_SP)
		if len(idSplit) != 2 {
			return fmt.Errorf("id is broken, id is %s", rs.Primary.ID)
		}

		instanceId := idSplit[0]
		fileName := idSplit[1]

		result, err := service.DescribeSqlserverBusinessIntelligenceFileById(ctx, instanceId, fileName)
		if err != nil {
			return err
		}

		if result == nil {
			return fmt.Errorf("sqlserver businessIntelligenceFile %s is not found", rs.Primary.ID)
		} else {
			return nil
		}
	}
}

const testAccSqlserverBusinessIntelligenceFile = testAccSqlserverBusinessIntelligenceInstance + `
resource "tencentcloud_sqlserver_business_intelligence_file" "business_intelligence_file" {
  instance_id = tencentcloud_sqlserver_business_intelligence_instance.business_intelligence_instance.id
  file_url = "https://keep-sqlserver-1308919341.cos.ap-guangzhou.myqcloud.com/test.xlsx"
  file_type = "FLAT"
  remark = "test case."
}
`
