package tencentcloud

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	ssl "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ssl/v20191205"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

const sslInstance = "tencentcloud_ssl_pay_certificate.ssl"

func TestAccTencentCloudSSLInstance(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSSLInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testSSLConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSSLInstanceExists(sslInstance),
					resource.TestCheckResourceAttrSet(sslInstance, "certificate_id"),
					resource.TestCheckResourceAttrSet(sslInstance, "order_id"),
					resource.TestCheckResourceAttrSet(sslInstance, "status"),
					resource.TestCheckResourceAttr(sslInstance, "product_id", "33"),
					resource.TestCheckResourceAttr(sslInstance, "domain_num", "1"),
					resource.TestCheckResourceAttr(sslInstance, "time_span", "1"),
					resource.TestCheckResourceAttr(sslInstance, "alias", "test-ssl"),
					resource.TestCheckResourceAttr(sslInstance, "project_id", "0"),
					resource.TestCheckResourceAttr(sslInstance, "information.#", "1"),
					resource.TestCheckResourceAttr(sslInstance, "information.0.csr_type", "online"),
					resource.TestCheckResourceAttr(sslInstance, "information.0.organization_name", "test"),
					resource.TestCheckResourceAttr(sslInstance, "information.0.organization_division", "test"),
					resource.TestCheckResourceAttr(sslInstance, "information.0.organization_address", "test"),
					resource.TestCheckResourceAttr(sslInstance, "information.0.organization_country", "CN"),
					resource.TestCheckResourceAttr(sslInstance, "information.0.organization_city", "test"),
					resource.TestCheckResourceAttr(sslInstance, "information.0.organization_region", "test"),
					resource.TestCheckResourceAttr(sslInstance, "information.0.postal_code", "0755"),
					resource.TestCheckResourceAttr(sslInstance, "information.0.phone_area_code", "0755"),
					resource.TestCheckResourceAttr(sslInstance, "information.0.phone_number", "12345678901"),
					resource.TestCheckResourceAttr(sslInstance, "information.0.verify_type", "DNS"),
					resource.TestCheckResourceAttr(sslInstance, "information.0.admin_first_name", "test"),
					resource.TestCheckResourceAttr(sslInstance, "information.0.admin_last_name", "test"),
					resource.TestCheckResourceAttr(sslInstance, "information.0.admin_phone_num", "12345678901"),
					resource.TestCheckResourceAttr(sslInstance, "information.0.admin_email", "test@tencent.com"),
					resource.TestCheckResourceAttr(sslInstance, "information.0.admin_position", "dev"),
					resource.TestCheckResourceAttr(sslInstance, "information.0.contact_first_name", "test"),
					resource.TestCheckResourceAttr(sslInstance, "information.0.contact_last_name", "test"),
					resource.TestCheckResourceAttr(sslInstance, "information.0.contact_email", "test@tencent.com"),
					resource.TestCheckResourceAttr(sslInstance, "information.0.contact_number", "12345678901"),
					resource.TestCheckResourceAttr(sslInstance, "information.0.contact_position", "dev"),
					resource.TestCheckResourceAttr(sslInstance, "information.0.certificate_domain", "www.domain.com"),
				),
			},
			{
				ResourceName:            sslInstance,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"information"},
			},
			{
				Config: testSSLUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSSLInstanceExists(sslInstance),
					resource.TestCheckResourceAttrSet(sslInstance, "certificate_id"),
					resource.TestCheckResourceAttrSet(sslInstance, "order_id"),
					resource.TestCheckResourceAttrSet(sslInstance, "status"),
					resource.TestCheckResourceAttr(sslInstance, "product_id", "33"),
					resource.TestCheckResourceAttr(sslInstance, "domain_num", "1"),
					resource.TestCheckResourceAttr(sslInstance, "time_span", "1"),
					resource.TestCheckResourceAttr(sslInstance, "alias", "test-ssl-update"),
					resource.TestCheckResourceAttr(sslInstance, "project_id", "0"),
					resource.TestCheckResourceAttr(sslInstance, "information.#", "1"),
					resource.TestCheckResourceAttr(sslInstance, "information.0.csr_type", "online"),
					resource.TestCheckResourceAttr(sslInstance, "information.0.organization_name", "test-update"),
					resource.TestCheckResourceAttr(sslInstance, "information.0.organization_division", "test"),
					resource.TestCheckResourceAttr(sslInstance, "information.0.organization_address", "test"),
					resource.TestCheckResourceAttr(sslInstance, "information.0.organization_country", "CN"),
					resource.TestCheckResourceAttr(sslInstance, "information.0.organization_city", "test"),
					resource.TestCheckResourceAttr(sslInstance, "information.0.organization_region", "test"),
					resource.TestCheckResourceAttr(sslInstance, "information.0.postal_code", "0755"),
					resource.TestCheckResourceAttr(sslInstance, "information.0.phone_area_code", "0755"),
					resource.TestCheckResourceAttr(sslInstance, "information.0.phone_number", "12345678901"),
					resource.TestCheckResourceAttr(sslInstance, "information.0.verify_type", "DNS"),
					resource.TestCheckResourceAttr(sslInstance, "information.0.admin_first_name", "test"),
					resource.TestCheckResourceAttr(sslInstance, "information.0.admin_last_name", "test"),
					resource.TestCheckResourceAttr(sslInstance, "information.0.admin_phone_num", "12345678901"),
					resource.TestCheckResourceAttr(sslInstance, "information.0.admin_email", "test@tencent.com"),
					resource.TestCheckResourceAttr(sslInstance, "information.0.admin_position", "dev"),
					resource.TestCheckResourceAttr(sslInstance, "information.0.contact_first_name", "test"),
					resource.TestCheckResourceAttr(sslInstance, "information.0.contact_last_name", "test"),
					resource.TestCheckResourceAttr(sslInstance, "information.0.contact_email", "test@tencent.com"),
					resource.TestCheckResourceAttr(sslInstance, "information.0.contact_number", "12345678901"),
					resource.TestCheckResourceAttr(sslInstance, "information.0.contact_position", "dev"),
					resource.TestCheckResourceAttr(sslInstance, "information.0.certificate_domain", "www.domain.com"),
				),
			},
		},
	})
}

func testAccCheckSSLInstanceDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	sslService := SSLService{
		client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_ssl_pay_certificate" {
			continue
		}
		resourceId := rs.Primary.ID
		ids := strings.Split(resourceId, FILED_SP)
		if len(ids) != 4 {
			return fmt.Errorf("ids param is error. id:  %s", resourceId)
		}
		request := ssl.NewDescribeCertificateDetailRequest()
		request.CertificateId = helper.String(ids[0])
		var (
			response *ssl.DescribeCertificateDetailResponse
			err      error
		)
		response, err = sslService.DescribeCertificateDetail(ctx, request)
		if err != nil {
			response, err = sslService.DescribeCertificateDetail(ctx, request)
			if err != nil {
				return err
			}
		}

		if response != nil && response.Response != nil && response.Response.Status != nil && response.Response.CertificateId != nil {
			status := *response.Response.Status
			// 6  7
			if status != 7 && status != 6 {
				return fmt.Errorf("the SSL certificate [certificateId = %s] has not been cancelled and the order is successful", ids[0])
			}
		}
	}
	return nil
}

func testAccCheckSSLInstanceExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)
		sslService := SSLService{
			client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
		}

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("[TECENT_TERRAFORM_CHECK][SSL instance][Exists] check: SSL instance %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("[TECENT_TERRAFORM_CHECK][SSL instance][Exists] check: SSL instance certificateId is not set")
		}
		resourceId := rs.Primary.ID
		ids := strings.Split(resourceId, FILED_SP)
		if len(ids) != 4 {
			return fmt.Errorf("ids param is error. id:  %s", resourceId)
		}
		request := ssl.NewDescribeCertificateDetailRequest()
		request.CertificateId = helper.String(ids[0])
		var (
			response *ssl.DescribeCertificateDetailResponse
			err      error
		)
		response, err = sslService.DescribeCertificateDetail(ctx, request)
		if err != nil {
			response, err = sslService.DescribeCertificateDetail(ctx, request)
			if err != nil {
				return err
			}
		}
		if response == nil || response.Response == nil || response.Response.CertificateId == nil {
			return fmt.Errorf("certificateId %s does not exist", ids[0])
		}

		return nil
	}
}

const testSSLConfig = `
resource "tencentcloud_ssl_pay_certificate" "ssl" {
    product_id = 33
    domain_num = 1
    alias      = "test-ssl"
    project_id = 0
    information {
        csr_type              = "online"
        certificate_domain    = "www.domain.com"
        organization_name     = "test"
        organization_division = "test"
        organization_address  = "test"
        organization_country  = "CN"
        organization_city     = "test"
        organization_region   = "test"
        postal_code           = "0755"
        phone_area_code       = "0755"
        phone_number          = "12345678901"
        verify_type           = "DNS"
        admin_first_name      = "test"
        admin_last_name       = "test"
        admin_phone_num       = "12345678901"
        admin_email           = "test@tencent.com"
        admin_position        = "dev"
        contact_first_name    = "test"
        contact_last_name     = "test"
        contact_email         = "test@tencent.com"
        contact_number        = "12345678901"
        contact_position      = "dev"
    }
}
`
const testSSLUpdateConfig = `
resource "tencentcloud_ssl_pay_certificate" "ssl" {
    product_id = 33
    domain_num = 1
    alias      = "test-ssl-update"
    project_id = 0
    information {
        csr_type              = "online"
        certificate_domain    = "www.domain.com"
        organization_name     = "test-update"
        organization_division = "test"
        organization_address  = "test"
        organization_country  = "CN"
        organization_city     = "test"
        organization_region   = "test"
        postal_code           = "0755"
        phone_area_code       = "0755"
        phone_number          = "12345678901"
        verify_type           = "DNS"
        admin_first_name      = "test"
        admin_last_name       = "test"
        admin_phone_num       = "12345678901"
        admin_email           = "test@tencent.com"
        admin_position        = "dev"
        contact_first_name    = "test"
        contact_last_name     = "test"
        contact_email         = "test@tencent.com"
        contact_number        = "12345678901"
        contact_position      = "dev"
    }
}
`
