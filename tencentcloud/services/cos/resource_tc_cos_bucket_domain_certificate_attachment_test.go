package cos_test

import (
	"context"
	"fmt"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/pkg/errors"

	localcos "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cos"
)

func TestAccTencentCloudCosBucketDomainCertificate_basic(t *testing.T) {
	id := new(string)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
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
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service := localcos.NewCosService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())

		ret, _, err := service.DescribeCosBucketDomainCertificate(ctx, *configId)

		if err != nil {
			return err
		}

		if ret == nil {
			return fmt.Errorf("cosBucketDoaminCertificate's status can not found! id:[%s]", *configId)
		}

		if ret.Status == localcos.CERT_ENABLED {
			return errors.New("the cosBucketDoaminCertificate still exists")
		}
		return nil
	}
}

func testAccCheckCosBucketDoaminCertificateExists(n string, id *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service := localcos.NewCosService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())

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

		if ret.Status == localcos.CERT_DISABLED {
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
  name = "` + tcacctest.DefaultCosCertificateName + `"
}

locals {
  publicKey = data.tencentcloud_ssl_certificates.foo1.certificates.0.cert
  privateKey = data.tencentcloud_ssl_certificates.foo1.certificates.0.key
}

`

const testAccCosBucketDomainCertificate_basic = tcacctest.UserInfoData + testAccCosDomain + testAccSSLCertificate + `
resource "tencentcloud_cos_bucket_domain_certificate_attachment" "basic" {
  bucket = "` + tcacctest.DefaultCosCertificateBucketPrefix + `-${local.app_id}"
  domain_certificate {
	domain = "` + tcacctest.DefaultCosCertDomainName + `"
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
