package tencentcloud

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccTencentCloudAPIGateWayCustomDomain(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCustomDomainDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCustomDomain,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCustomDomainExists("tencentcloud_api_gateway_custom_domain.foo"),
					resource.TestCheckResourceAttr("tencentcloud_api_gateway_custom_domain.foo", "service_id", "service-ohxqslqe"),
					resource.TestCheckResourceAttr("tencentcloud_api_gateway_custom_domain.foo", "sub_domain", "tic-test.dnsv1.com"),
					resource.TestCheckResourceAttr("tencentcloud_api_gateway_custom_domain.foo", "protocol", "http"),
					resource.TestCheckResourceAttr("tencentcloud_api_gateway_custom_domain.foo", "net_type", "OUTER"),
					resource.TestCheckResourceAttr("tencentcloud_api_gateway_custom_domain.foo", "is_default_mapping", "false"),
					resource.TestCheckResourceAttrSet("tencentcloud_api_gateway_custom_domain.foo", "default_domain"),
					resource.TestCheckResourceAttr("tencentcloud_api_gateway_custom_domain.foo", "path_mappings.#", "2"),
				),
			},
			{
				Config: testAccCustomDomainUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCustomDomainExists("tencentcloud_api_gateway_custom_domain.foo"),
					resource.TestCheckResourceAttr("tencentcloud_api_gateway_custom_domain.foo", "service_id", "service-ohxqslqe"),
					resource.TestCheckResourceAttr("tencentcloud_api_gateway_custom_domain.foo", "sub_domain", "tic-test.dnsv1.com"),
					resource.TestCheckResourceAttr("tencentcloud_api_gateway_custom_domain.foo", "protocol", "http"),
					resource.TestCheckResourceAttr("tencentcloud_api_gateway_custom_domain.foo", "net_type", "OUTER"),
					resource.TestCheckResourceAttr("tencentcloud_api_gateway_custom_domain.foo", "is_default_mapping", "false"),
					resource.TestCheckResourceAttrSet("tencentcloud_api_gateway_custom_domain.foo", "default_domain"),
					resource.TestCheckResourceAttr("tencentcloud_api_gateway_custom_domain.foo", "path_mappings.#", "1"),
				),
			},
		},
	})
}

func testAccCheckCustomDomainDestroy(s *terraform.State) error {
	var (
		logId             = getLogId(contextNil)
		ctx               = context.WithValue(context.TODO(), logIdKey, logId)
		apigatewayService = APIGatewayService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
	)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_api_gateway_custom_domain" {
			continue
		}

		params := strings.Split(rs.Primary.ID, FILED_SP)
		if len(params) != 2 {
			return fmt.Errorf("ids param is error. id:  %s", rs.Primary.ID)
		}
		serviceId := params[0]
		subDomain := params[1]

		resultList, err := apigatewayService.DescribeServiceSubDomainsService(ctx, serviceId, subDomain)
		if err != nil {
			return err
		}
		if len(resultList) > 0 {
			return fmt.Errorf("custom domain: %s still exist", subDomain)
		}
	}
	return nil
}

func testAccCheckCustomDomainExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		var (
			logId             = getLogId(contextNil)
			ctx               = context.WithValue(context.TODO(), logIdKey, logId)
			apigatewayService = APIGatewayService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
		)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("API getway custom domain %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("API getway custom domain id is not set")
		}

		params := strings.Split(rs.Primary.ID, FILED_SP)
		if len(params) != 2 {
			return fmt.Errorf("ids param is error. id:  %s", rs.Primary.ID)
		}
		serviceId := params[0]
		subDomain := params[1]

		resultList, err := apigatewayService.DescribeServiceSubDomainsService(ctx, serviceId, subDomain)
		if err != nil {
			return err
		}
		if len(resultList) == 0 {
			return fmt.Errorf("custom domain: %s create failed", subDomain)
		}
		return nil
	}
}

const testAccCustomDomain = `
resource "tencentcloud_api_gateway_custom_domain" "foo" {
	service_id         = "service-ohxqslqe"
	sub_domain         = "tic-test.dnsv1.com"
	protocol           = "http"
	net_type           = "OUTER"
	is_default_mapping = "false"
	default_domain     = "service-ohxqslqe-1259649581.gz.apigw.tencentcs.com"
	path_mappings      = ["/good#test","/root#release"]
}
`

const testAccCustomDomainUpdate = `
resource "tencentcloud_api_gateway_custom_domain" "foo" {
	service_id         = "service-ohxqslqe"
	sub_domain         = "tic-test.dnsv1.com"
	protocol           = "http"
	net_type           = "OUTER"
	is_default_mapping = "false"
	default_domain     = "service-ohxqslqe-1259649581.gz.apigw.tencentcs.com"
	path_mappings      = ["/good#test"]
}
`
