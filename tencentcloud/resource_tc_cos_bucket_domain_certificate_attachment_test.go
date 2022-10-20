package tencentcloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/pkg/errors"
)

func TestAccTencentCloudCosBucketDomainCertificate_basic(t *testing.T) {
	id := new(string)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCosBucketDoaminCertificateDestroy(id),
		Steps: []resource.TestStep{
			{
				Config: testAccCosBucketDomainCertificate_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCosBucketDoaminCertificateExists("tencentcloud_cos_bucket_domain_certificate_attachment.basic", id),
					resource.TestCheckResourceAttrSet("tencentcloud_cos_bucket_domain_certificate_attachment.basic", "bucket"),
					resource.TestCheckResourceAttrSet("tencentcloud_cos_bucket_domain_certificate_attachment.basic", "domain_certificate.0.certificate.0.custom_cert.0.cert"),
					resource.TestCheckResourceAttrSet("tencentcloud_cos_bucket_domain_certificate_attachment.basic", "domain_certificate.0.certificate.0.custom_cert.0.private_key"),
					resource.TestCheckResourceAttr("tencentcloud_cos_bucket_domain_certificate_attachment.basic", "domain_certificate.0.certificate.0.cert_type", "CustomCert"),
				),
			},
		},
	})
}

func testAccCheckCosBucketDoaminCertificateDestroy(configId *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)
		service := CosService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}

		ret, _, err := service.DescribeCosBucketDomainCertificate(ctx, *configId)

		if err != nil {
			return err
		}

		if ret == nil {
			return fmt.Errorf("cosBucketDoaminCertificate's status can not found! id:[%s]", *configId)
		}

		if ret.Status == CERT_ENABLED {
			return errors.New("the cosBucketDoaminCertificate still exists")
		}
		return nil
	}
}

func testAccCheckCosBucketDoaminCertificateExists(n string, id *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)
		service := CosService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		configId := rs.Primary.ID
		if configId == "" {
			return errors.New("no cosBucketDoaminCertificate configID is set")
		}

		ret, _, err := service.DescribeCosBucketDomainCertificate(ctx, configId)

		if err != nil {
			return err
		}

		if ret == nil {
			return fmt.Errorf("cosBucketDoaminCertificate's status can not found: %s", rs.Primary.ID)
		}

		if ret.Status == CERT_DISABLED {
			return fmt.Errorf("cosBucketDoaminCertificate does not exist: %s", rs.Primary.ID)
		}

		*id = configId

		return nil
	}
}

const testAccCosDomain = `
provider "tencentcloud" {
  region = "ap-singapore"
}

`

// name = "keep-c-ssl"keep-cos-domain-cert
const testAccSSLCertificate = `
data "tencentcloud_ssl_certificates" "foo1" {
  name = "` + defaultCosCertificateName + `"
}

locals {
  publicKey = data.tencentcloud_ssl_certificates.foo1.certificates.0.cert
  privateKey = data.tencentcloud_ssl_certificates.foo1.certificates.0.key
}

`

const testAccCosBucketDomainCertificate_basic = userInfoData + testAccCosDomain + testAccSSLCertificate + `
resource "tencentcloud_cos_bucket_domain_certificate_attachment" "basic" {
  bucket = "` + defaultCosCertificateBucketPrefix + `-${local.app_id}"
  domain_certificate {
	domain = "` + defaultCosCertDomainName + `"
    certificate {
      cert_type = "CustomCert"
      custom_cert {
        cert        = "${local.publicKey}"
        private_key = "${local.privateKey}"
      }
    }
  }
}
`
