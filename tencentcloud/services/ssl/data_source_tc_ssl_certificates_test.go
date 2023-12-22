package ssl_test

import (
	"fmt"
	"regexp"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudSslCertificatesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: TestAccDataSourceTencentCloudSslCertificatesBasic,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_ssl_certificates.foo"),
					resource.TestMatchResourceAttr("data.tencentcloud_ssl_certificates.foo", "certificates.#", regexp.MustCompile(`^[1-9]\d*$`)),
					resource.TestCheckResourceAttrSet("data.tencentcloud_ssl_certificates.foo", "certificates.0.id"),
					resource.TestCheckResourceAttr("data.tencentcloud_ssl_certificates.foo", "certificates.0.name", "keep-ssl-ca"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_ssl_certificates.foo", "certificates.0.type"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_ssl_certificates.foo", "certificates.0.project_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_ssl_certificates.foo", "certificates.0.cert"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_ssl_certificates.foo", "certificates.0.product_zh_name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_ssl_certificates.foo", "certificates.0.domain"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_ssl_certificates.foo", "certificates.0.status"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_ssl_certificates.foo", "certificates.0.begin_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_ssl_certificates.foo", "certificates.0.end_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_ssl_certificates.foo", "certificates.0.create_time"),
				),
			},
		},
	})
}

func TestAccTencentCloudSslCertificatesDataSource_type(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: TestAccDataSourceTencentCloudSslCertificatesType,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_ssl_certificates.foo"),
					resource.TestMatchResourceAttr("data.tencentcloud_ssl_certificates.foo", "certificates.#", regexp.MustCompile(`^[1-9]\d*$`)),
					resource.TestCheckResourceAttrSet("data.tencentcloud_ssl_certificates.foo", "certificates.0.id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_ssl_certificates.foo", "certificates.0.name"),
					resource.TestCheckResourceAttr("data.tencentcloud_ssl_certificates.foo", "certificates.0.type", "CA"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_ssl_certificates.foo", "certificates.0.project_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_ssl_certificates.foo", "certificates.0.cert"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_ssl_certificates.foo", "certificates.0.product_zh_name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_ssl_certificates.foo", "certificates.0.domain"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_ssl_certificates.foo", "certificates.0.status"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_ssl_certificates.foo", "certificates.0.begin_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_ssl_certificates.foo", "certificates.0.end_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_ssl_certificates.foo", "certificates.0.create_time"),
				),
			},
		},
	})
}

func TestAccTencentCloudSslCertificatesDataSource_id(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: TestAccDataSourceTencentCloudSslCertificatesId,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_ssl_certificates.foo"),
					resource.TestCheckResourceAttr("data.tencentcloud_ssl_certificates.foo", "certificates.#", "1"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_ssl_certificates.foo", "certificates.0.id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_ssl_certificates.foo", "certificates.0.name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_ssl_certificates.foo", "certificates.0.type"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_ssl_certificates.foo", "certificates.0.project_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_ssl_certificates.foo", "certificates.0.cert"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_ssl_certificates.foo", "certificates.0.product_zh_name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_ssl_certificates.foo", "certificates.0.domain"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_ssl_certificates.foo", "certificates.0.status"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_ssl_certificates.foo", "certificates.0.begin_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_ssl_certificates.foo", "certificates.0.end_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_ssl_certificates.foo", "certificates.0.create_time"),
					resource.TestCheckResourceAttr("data.tencentcloud_ssl_certificates.foo", "subject_names.#", "0"),
				),
			},
		},
	})
}

var TestAccDataSourceTencentCloudSslCertificatesBasic = fmt.Sprintf(`
resource "tencentcloud_ssl_certificate" "foo" {
  type = "CA"
  cert = "%s"
  name = "keep-ssl-ca"
}

data "tencentcloud_ssl_certificates" "foo" {
  name = tencentcloud_ssl_certificate.foo.name
}
`, testAccSslCertificateCA)

var TestAccDataSourceTencentCloudSslCertificatesType = fmt.Sprintf(`
resource "tencentcloud_ssl_certificate" "foo" {
  type = "CA"
  cert = "%s"
  name = "keep-ssl-ca"
}

data "tencentcloud_ssl_certificates" "foo" {
  type = tencentcloud_ssl_certificate.foo.type
}
`, testAccSslCertificateCA)

var TestAccDataSourceTencentCloudSslCertificatesId = fmt.Sprintf(`
resource "tencentcloud_ssl_certificate" "foo" {
  type = "CA"
  cert = "%s"
  name = "keep-ssl-ca"
}

data "tencentcloud_ssl_certificates" "foo" {
  id = tencentcloud_ssl_certificate.foo.id
}
`, testAccSslCertificateCA)
