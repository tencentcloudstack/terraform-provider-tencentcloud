package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudSslDescribeCertificateBindResourceTaskResultDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSslDescribeCertificateBindResourceTaskResultDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_ssl_describe_certificate_bind_resource_task_result.describe_certificate_bind_resource_task_result")),
			},
		},
	})
}

const testAccSslDescribeCertificateBindResourceTaskResultDataSource = `

data "tencentcloud_ssl_describe_certificate_bind_resource_task_result" "describe_certificate_bind_resource_task_result" {
  task_ids = ""
  }

`
