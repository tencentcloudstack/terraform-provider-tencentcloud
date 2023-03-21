package tencentcloud

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func NeedFixTestAccTencentCloudDbbrainSecurityAuditLogDownloadUrlsDataSource_basic(t *testing.T) {
	t.Parallel()
	loc, _ := time.LoadLocation("Asia/Chongqing")
	startTime := time.Now().Add(-30 * time.Minute).In(loc).Format("2006-01-02 15:04:05")
	endTime := time.Now().In(loc).Format("2006-01-02 15:04:05")
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDbbrainSecurityAuditLogDownloadUrlsDataSource, defaultDbBrainsagId, startTime, endTime, defaultDbBrainsagId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_dbbrain_security_audit_log_download_urls.test"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dbbrain_security_audit_log_download_urls.test", "urls.#"),
				),
			},
		},
	})
}

const testAccDbbrainSecurityAuditLogDownloadUrlsDataSource = `

resource "tencentcloud_dbbrain_security_audit_log_export_task" "task" {
	sec_audit_group_id = "%s"
	start_time = "%s"
	end_time = "%s"
	product = "mysql"
	danger_levels = [0,1,2]
}

data "tencentcloud_dbbrain_security_audit_log_download_urls" "test" {
	sec_audit_group_id = "%s"
	async_request_id = tencentcloud_dbbrain_security_audit_log_export_task.task.async_request_id
	product = "mysql"
}

`
