package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudSslDescribeCertificateBindResourceTaskDetailDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSslDescribeCertificateBindResourceTaskDetailDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_ssl_describe_certificate_bind_resource_task_detail.describe_certificate_bind_resource_task_detail")),
			},
		},
	})
}

const testAccSslDescribeCertificateBindResourceTaskDetailDataSource = `

data "tencentcloud_ssl_describe_certificate_bind_resource_task_detail" "describe_certificate_bind_resource_task_detail" {
  task_id = ""
  resource_types = 
  regions = 
                        }

`
