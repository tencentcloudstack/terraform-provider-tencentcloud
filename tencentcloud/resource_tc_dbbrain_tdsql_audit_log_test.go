package tencentcloud

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func init() {
	resource.AddTestSweepers("tencentcloud_dbbrain_tdsql_audit_log", &resource.Sweeper{
		Name: "tencentcloud_dbbrain_tdsql_audit_log",
		F:    testSweepDbbrainTdsqlAuditLog,
	})
}

// go test -v ./tencentcloud -sweep=ap-guangzhou -sweep-run=tencentcloud_dbbrain_tdsql_audit_log
func testSweepDbbrainTdsqlAuditLog(r string) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	cli, _ := sharedClientForRegion(r)
	dbbrainService := DbbrainService{client: cli.(*TencentCloudClient).apiV3Conn}
	products := []string{
		"dcdb",
	}

	for _, product := range products {
		auditLogs, err := dbbrainService.DescribeDbbrainTdsqlAuditLogById(ctx, nil, defaultDcdbInstanceId, product)
		if err != nil {
			return err
		}
		if len(auditLogs) == 0 {
			return fmt.Errorf("Dbbrain tdsql audit log not exists.")
		}

		for _, audit := range auditLogs {
			err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
				err := dbbrainService.DeleteDbbrainTdsqlAuditLogById(ctx, helper.Int64ToStr(*audit.AsyncRequestId), defaultDcdbInstanceId, product)
				if err != nil {
					return retryError(err)
				}
				return nil
			})

			if err != nil {
				return fmt.Errorf("[ERROR] sweeper Dbbrain tdsql audit logs:[AsyncRequestId:%v, instanceId:%s, product:%s] failed! reason:[%s]", *audit.AsyncRequestId, defaultDbBrainInstanceId, product, err.Error())
			}
		}

	}
	return nil
}

func TestAccTencentCloudDbbrainTdsqlAuditLogResource_basic(t *testing.T) {
	t.Parallel()
	loc, _ := time.LoadLocation("Asia/Chongqing")
	startTime := time.Now().Add(-24 * time.Hour).In(loc).Format("2006-01-02 15:04:05")
	endTime := time.Now().Add(-2 * time.Hour).In(loc).Format("2006-01-02 15:04:05")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDbbrainTdsqlAuditLogDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDbbrainTdsqlAuditLog, defaultDcdbInstanceId, startTime, endTime),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDbbrainTdsqlAuditLogExists("tencentcloud_dbbrain_tdsql_audit_log.my_log"),
					resource.TestCheckResourceAttrSet("tencentcloud_dbbrain_tdsql_audit_log.my_log", "id"),
					resource.TestCheckResourceAttr("tencentcloud_dbbrain_tdsql_audit_log.my_log", "product", "dcdb"),
					resource.TestCheckResourceAttr("tencentcloud_dbbrain_tdsql_audit_log.my_log", "node_request_type", "dcdb"),
					resource.TestCheckResourceAttr("tencentcloud_dbbrain_tdsql_audit_log.my_log", "instance_id", defaultDcdbInstanceId),
					resource.TestCheckResourceAttr("tencentcloud_dbbrain_tdsql_audit_log.my_log", "start_time", startTime),
					resource.TestCheckResourceAttr("tencentcloud_dbbrain_tdsql_audit_log.my_log", "end_time", endTime),
					resource.TestCheckResourceAttr("tencentcloud_dbbrain_tdsql_audit_log.my_log", "filter.#", "1"),
					resource.TestCheckTypeSetElemAttr("tencentcloud_dbbrain_tdsql_audit_log.my_log", "filter.0.host.*", "%"),
					resource.TestCheckTypeSetElemAttr("tencentcloud_dbbrain_tdsql_audit_log.my_log", "filter.0.host.*", "127.0.0.1"),
					resource.TestCheckTypeSetElemAttr("tencentcloud_dbbrain_tdsql_audit_log.my_log", "filter.0.user.*", "tf_test"),
					resource.TestCheckTypeSetElemAttr("tencentcloud_dbbrain_tdsql_audit_log.my_log", "filter.0.user.*", "mysql"),
				),
			},
		},
	})
}

func testAccCheckDbbrainTdsqlAuditLogDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	dbbrainService := DbbrainService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_dbbrain_tdsql_audit_log" {
			continue
		}

		ids := strings.Split(rs.Primary.ID, FILED_SP)
		if len(ids) != 3 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		asyncRequestId := ids[0]
		instanceId := ids[1]
		product := ids[2]

		auditLog, err := dbbrainService.DescribeDbbrainTdsqlAuditLogById(ctx, &asyncRequestId, instanceId, product)
		if err != nil {
			return err
		}

		if auditLog != nil {
			return fmt.Errorf("Dbbrain tdsql audit log, Id: %v", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckDbbrainTdsqlAuditLogExists(re string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)
		dbbrainService := DbbrainService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}

		rs, ok := s.RootModule().Resources[re]
		if !ok {
			return fmt.Errorf("Dbbrain tdsql audit log %s is not found", re)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("Dbbrain tdsql audit log id is not set")
		}

		ids := strings.Split(rs.Primary.ID, FILED_SP)
		if len(ids) != 3 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		asyncRequestId := ids[0]
		instanceId := ids[1]
		product := ids[2]

		auditLog, err := dbbrainService.DescribeDbbrainTdsqlAuditLogById(ctx, &asyncRequestId, instanceId, product)
		if err != nil {
			return err
		}

		if auditLog == nil {
			return fmt.Errorf("Dbbrain tdsql audit log, Id: %v", rs.Primary.ID)
		}
		return nil
	}
}

const testAccDbbrainTdsqlAuditLog = `

resource "tencentcloud_dbbrain_tdsql_audit_log" "my_log" {
  product = "dcdb"
  node_request_type = "dcdb"
  instance_id = "%s"
  start_time = "%s"
  end_time = "%s"
  filter {
		host = ["%%", "127.0.0.1"]
		user = ["tf_test", "mysql"]
  }
}

`
