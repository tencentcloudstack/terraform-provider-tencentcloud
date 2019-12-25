package tencentcloud

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceTencentCloudSslCertificates_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: TestAccDataSourceTencentCloudSslCertificatesBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_ssl_certificates.foo"),
					resource.TestMatchResourceAttr("data.tencentcloud_ssl_certificates.foo", "certificates.#", regexp.MustCompile(`^[1-9]\d*$`)),
					resource.TestCheckResourceAttrSet("data.tencentcloud_ssl_certificates.foo", "certificates.0.id"),
					resource.TestCheckResourceAttr("data.tencentcloud_ssl_certificates.foo", "certificates.0.name", "ci-test-ssl-ca"),
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

func TestAccDataSourceTencentCloudSslCertificates_type(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: TestAccDataSourceTencentCloudSslCertificatesType,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_ssl_certificates.foo"),
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

func TestAccDataSourceTencentCloudSslCertificates_id(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: TestAccDataSourceTencentCloudSslCertificatesId,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_ssl_certificates.foo"),
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
  name = "ci-test-ssl-ca"
}

data "tencentcloud_ssl_certificates" "foo" {
  name = tencentcloud_ssl_certificate.foo.name
}
`, testAccSslCertificateCA)

var TestAccDataSourceTencentCloudSslCertificatesType = fmt.Sprintf(`
resource "tencentcloud_ssl_certificate" "foo" {
  type = "CA"
  cert = "%s"
  name = "ci-test-ssl-ca"
}

data "tencentcloud_ssl_certificates" "foo" {
  type = tencentcloud_ssl_certificate.foo.type
}
`, testAccSslCertificateCA)

var TestAccDataSourceTencentCloudSslCertificatesId = fmt.Sprintf(`
resource "tencentcloud_ssl_certificate" "foo" {
  type = "CA"
  cert = "%s"
  name = "ci-test-ssl-ca"
}

data "tencentcloud_ssl_certificates" "foo" {
  id = tencentcloud_ssl_certificate.foo.id
}
`, testAccSslCertificateCA)
