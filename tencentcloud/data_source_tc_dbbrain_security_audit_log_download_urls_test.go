package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudDbbrainSecurityAuditLogDownloadUrlsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDbbrainSecurityAuditLogDownloadUrlsDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_dbbrain_security_audit_log_download_urls.security_audit_log_download_urls")),
			},
		},
	})
}

const testAccDbbrainSecurityAuditLogDownloadUrlsDataSource = `

data "tencentcloud_dbbrain_security_audit_log_download_urls" "security_audit_log_download_urls" {
  sec_audit_group_id = ""
  async_request_id = 
  product = ""
  }

`
