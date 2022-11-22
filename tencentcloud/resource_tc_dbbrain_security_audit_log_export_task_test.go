package tencentcloud

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func TestAccTencentCloudDbbrainSecurityAuditLogExportTaskResource_basic(t *testing.T) {
	t.Parallel()
	baseTime := time.Now().UTC()
	startTime := baseTime.Add(-2 * time.Hour).Format(time.RFC3339)
	endTime := baseTime.Add(2 * time.Hour).Format(time.RFC3339)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDbbrainSecurityAuditLogExportTaskDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDbbrainSecurityAuditLogExportTask(startTime, endTime),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckDbbrainSqlFilterExists("tencentcloud_dbbrain_security_audit_log_export_task.task"),
					resource.TestCheckResourceAttrSet("tencentcloud_dbbrain_security_audit_log_export_task.task", "id"),
					resource.TestCheckResourceAttr("tencentcloud_dbbrain_security_audit_log_export_task.task", "sec_audit_group_id", defaultDbBrainsagId),
					resource.TestCheckResourceAttr("tencentcloud_dbbrain_security_audit_log_export_task.task", "start_time", startTime),
					resource.TestCheckResourceAttr("tencentcloud_dbbrain_security_audit_log_export_task.task", "end_time", endTime),
					resource.TestCheckResourceAttr("tencentcloud_dbbrain_security_audit_log_export_task.task", "product", "mysql"),
					resource.TestCheckResourceAttr("tencentcloud_dbbrain_security_audit_log_export_task.task", "danger_levels.#", "3"),
				),
			},
			{
				ResourceName:      "tencentcloud_dbbrain_security_audit_log_export_task.task",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckDbbrainSecurityAuditLogExportTaskDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	dbbrainService := DbbrainService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_dbbrain_security_audit_log_export_task" {
			continue
		}

		idSplit := strings.Split(rs.Primary.ID, FILED_SP)
		if len(idSplit) != 2 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		secAuditGroupId := helper.String(idSplit[0])
		asyncRequestId := helper.String(idSplit[1])

		filter, err := dbbrainService.DescribeDbbrainSecurityAuditLogExportTask(ctx, secAuditGroupId, asyncRequestId, nil)
		if err != nil {
			return err
		}

		if filter != nil {
			return fmt.Errorf("Dbbrain security audit log export task still exist, Id: %v", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckDbbrainSecurityAuditLogExportTaskExists(re string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)
		dbbrainService := DbbrainService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}

		rs, ok := s.RootModule().Resources[re]
		if !ok {
			return fmt.Errorf("Dbbrain security audit log export task %s is not found", re)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("Dbbrain security audit log export task id is not set")
		}

		idSplit := strings.Split(rs.Primary.ID, FILED_SP)
		if len(idSplit) != 2 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		secAuditGroupId := helper.String(idSplit[0])
		asyncRequestId := helper.String(idSplit[1])

		filter, err := dbbrainService.DescribeDbbrainSecurityAuditLogExportTask(ctx, secAuditGroupId, asyncRequestId, nil)
		if err != nil {
			return err
		}

		if filter == nil {
			return fmt.Errorf("Dbbrain security audit log export task not found, Id: %v", rs.Primary.ID)
		}
		return nil
	}
}

func testAccDbbrainSecurityAuditLogExportTask(st, et string) string {
	return fmt.Sprintf(`%s

resource "tencentcloud_dbbrain_security_audit_log_export_task" "task" {
  sec_audit_group_id = "%s"
  start_time = "%s"
  end_time = "%s"
  product = "mysql"
  danger_levels = [0,1,2]
}

`, CommonPresetMysql, defaultDbBrainsagId, st, et)
}
