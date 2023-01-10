package tencentcloud

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccTencentCloudCynosdbAuditLogFileResource(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCynosdbAuditLogFileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCynosdbAuditLogFile,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCynosdbCynosdbAuditLogFileExists("tencentcloud_cynosdb_audit_log_file.audit_log_file"),
				),
			},
		},
	})
}

func testAccCheckCynosdbAuditLogFileDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CynosdbService{
		client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_cynosdb_audit_log_file" {
			continue
		}

		idSplit := strings.Split(rs.Primary.ID, FILED_SP)
		if len(idSplit) != 2 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		instanceId := idSplit[0]
		fileName := idSplit[1]

		auditLogFile, err := service.DescribeCynosdbAuditLogFileById(ctx, instanceId, fileName)
		if err != nil {
			return err
		}
		if auditLogFile == nil {
			return nil
		}
		return fmt.Errorf("cynosdb audit log file still exists: %s", rs.Primary.ID)
	}
	return nil
}

func testAccCheckCynosdbCynosdbAuditLogFileExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("cynosdb audit log file %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("cynosdb audit log file id is not set")
		}
		service := CynosdbService{
			client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
		}
		idSplit := strings.Split(rs.Primary.ID, FILED_SP)
		if len(idSplit) != 2 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		instanceId := idSplit[0]
		fileName := idSplit[1]

		auditLogFile, err := service.DescribeCynosdbAuditLogFileById(ctx, instanceId, fileName)
		if err != nil {
			return err
		}
		if auditLogFile == nil {
			return fmt.Errorf("cynosdb audit log file doesn't exist: %s", rs.Primary.ID)
		}
		return nil
	}
}

const testAccCynosdbAuditLogFile = CommonCynosdb + `
resource "tencentcloud_cynosdb_audit_log_file" "audit_log_file" {
  instance_id = var.cynosdb_cluster_instance_id
  start_time = "2023-01-04 16:54:20"
  end_time =  "2023-01-04 16:55:00"
}
`
