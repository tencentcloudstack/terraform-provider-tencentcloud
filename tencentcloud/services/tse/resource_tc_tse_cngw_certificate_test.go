package tse_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svctse "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tse"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
)

// go test -test.run TestAccTencentCloudTseCngwCertificateResource_basic -v
func TestAccTencentCloudTseCngwCertificateResource_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckTseCngwCertificateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTseCngwCertificate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTseCngwCngwCertificateExists("tencentcloud_tse_cngw_certificate.cngw_certificate"),
					resource.TestCheckResourceAttrSet("tencentcloud_tse_cngw_certificate.cngw_certificate", "id"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_certificate.cngw_certificate", "name", "tf-certificate"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_certificate.cngw_certificate", "cert_id", tcacctest.DefaultTseCertId),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_certificate.cngw_certificate", "gateway_id", tcacctest.DefaultTseGatewayId),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_certificate.cngw_certificate", "bind_domains.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_certificate.cngw_certificate", "bind_domains.0", "example.com"),
				),
			},
			{
				ResourceName:      "tencentcloud_tse_cngw_certificate.cngw_certificate",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccTseCngwCertificateUp,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTseCngwCngwCertificateExists("tencentcloud_tse_cngw_certificate.cngw_certificate"),
					resource.TestCheckResourceAttrSet("tencentcloud_tse_cngw_certificate.cngw_certificate", "id"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_certificate.cngw_certificate", "name", "tf-certificate-up"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_certificate.cngw_certificate", "cert_id", tcacctest.DefaultTseCertId),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_certificate.cngw_certificate", "gateway_id", tcacctest.DefaultTseGatewayId),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_certificate.cngw_certificate", "bind_domains.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_certificate.cngw_certificate", "bind_domains.0", "example-up.com"),
				),
			},
		},
	})
}

func testAccCheckTseCngwCertificateDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	service := svctse.NewTseService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_tse_cngw_certificate" {
			continue
		}

		idSplit := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
		if len(idSplit) != 2 {
			return fmt.Errorf("invalid ID %s", rs.Primary.ID)
		}
		gatewayId := idSplit[0]
		certificateId := idSplit[1]

		res, err := service.DescribeTseCngwCertificateById(ctx, gatewayId, certificateId)
		if err != nil {
			if sdkErr, ok := err.(*sdkErrors.TencentCloudSDKError); ok {
				if sdkErr.Code == "ResourceNotFound.ResourceNotFound" {
					return nil
				}
			}
			return err
		}

		if res != nil {
			return fmt.Errorf("tse certificate %s still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckTseCngwCngwCertificateExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}

		idSplit := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
		if len(idSplit) != 2 {
			return fmt.Errorf("invalid ID %s", rs.Primary.ID)
		}
		gatewayId := idSplit[0]
		certificateId := idSplit[1]

		service := svctse.NewTseService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
		res, err := service.DescribeTseCngwCertificateById(ctx, gatewayId, certificateId)
		if err != nil {
			return err
		}

		if res == nil {
			return fmt.Errorf("tse certificate %s is not found", rs.Primary.ID)
		}

		return nil
	}
}

const testAccTseCngwCertificate = tcacctest.DefaultTseVar + `

resource "tencentcloud_tse_cngw_certificate" "cngw_certificate" {
  gateway_id   = var.gateway_id
  bind_domains = ["example.com"]
  cert_id      = var.cert_id
  name         = "tf-certificate"
}

`

const testAccTseCngwCertificateUp = tcacctest.DefaultTseVar + `

resource "tencentcloud_tse_cngw_certificate" "cngw_certificate" {
  gateway_id   = var.gateway_id
  bind_domains = ["example-up.com"]
  cert_id      = var.cert_id
  name         = "tf-certificate-up"
}

`
