package ssl_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudSslDescribeCertificateDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_SSL)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSslDescribeCertificateDataSource,
				Check: resource.ComposeTestCheckFunc(tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_ssl_describe_certificate.describe_certificate"),
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
