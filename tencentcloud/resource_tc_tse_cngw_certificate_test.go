package tencentcloud

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
)

// go test -test.run TestAccTencentCloudTseCngwCertificateResource_basic -v
func TestAccTencentCloudTseCngwCertificateResource_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTseCngwCertificateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTseCngwCertificate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTseCngwCngwCertificateExists("tencentcloud_tse_cngw_certificate.cngw_certificate"),
					resource.TestCheckResourceAttrSet("tencentcloud_tse_cngw_certificate.cngw_certificate", "id"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_certificate.cngw_certificate", "name", "tf-certificate"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_certificate.cngw_certificate", "cert_id", defaultTseCertId),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_certificate.cngw_certificate", "gateway_id", defaultTseGatewayId),
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
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_certificate.cngw_certificate", "cert_id", defaultTseCertId),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_certificate.cngw_certificate", "gateway_id", defaultTseGatewayId),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_certificate.cngw_certificate", "bind_domains.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_tse_cngw_certificate.cngw_certificate", "bind_domains.0", "example-up.com"),
				),
			},
		},
	})
}

func testAccCheckTseCngwCertificateDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	service := TseService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_tse_cngw_certificate" {
			continue
		}

		idSplit := strings.Split(rs.Primary.ID, FILED_SP)
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
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}

		idSplit := strings.Split(rs.Primary.ID, FILED_SP)
		if len(idSplit) != 2 {
			return fmt.Errorf("invalid ID %s", rs.Primary.ID)
		}
		gatewayId := idSplit[0]
		certificateId := idSplit[1]

		service := TseService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
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

const testAccTseCngwCertificate = DefaultTseVar + `

resource "tencentcloud_tse_cngw_certificate" "cngw_certificate" {
  gateway_id   = var.gateway_id
  bind_domains = ["example.com"]
  cert_id      = var.cert_id
  name         = "tf-certificate"
}

`

const testAccTseCngwCertificateUp = DefaultTseVar + `

resource "tencentcloud_tse_cngw_certificate" "cngw_certificate" {
  gateway_id   = var.gateway_id
  bind_domains = ["example-up.com"]
  cert_id      = var.cert_id
  name         = "tf-certificate-up"
}

`
