package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudSslDescribeCertificateDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCommon(t, ACCOUNT_TYPE_SSL)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSslDescribeCertificateDataSource,
				Check: resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_ssl_describe_certificate.describe_certificate"),
					resource.TestCheckResourceAttr("data.tencentcloud_ssl_describe_certificate.describe_certificate", "certificate_id", "8mCN3eKd")),
			},
		},
	})
}

const testAccSslDescribeCertificateDataSource = `

data "tencentcloud_ssl_describe_certificate" "describe_certificate" {
  certificate_id = "8mCN3eKd"
}
`
