package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var testAPIGatewayDomainSourceName = "data.tencentcloud_api_gateway_customer_domains"

// go test -i; go test -test.run TestAccTencentAPIGatewayCustomerDomain_basic -v
func TestAccTencentAPIGatewayCustomerDomain_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCustomDomainDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTestAccTencentAPIGatewayDomain(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckCustomDomainExists("tencentcloud_api_gateway_custom_domain.foo"),
					resource.TestCheckResourceAttrSet(testAPIGatewayDomainSourceName+".id", "list.#"),
					resource.TestCheckResourceAttrSet(testAPIGatewayDomainSourceName+".id", "list.0.domain_name"),
					resource.TestCheckResourceAttrSet(testAPIGatewayDomainSourceName+".id", "list.0.is_status_on"),
					resource.TestCheckResourceAttrSet(testAPIGatewayDomainSourceName+".id", "list.0.is_default_mapping"),
					resource.TestCheckResourceAttrSet(testAPIGatewayDomainSourceName+".id", "list.0.net_type"),
					resource.TestCheckResourceAttrSet(testAPIGatewayDomainSourceName+".id", "list.0.path_mappings.#"),
				),
			},
		},
	})
}

func testAccTestAccTencentAPIGatewayDomain() string {
	return `
resource "tencentcloud_api_gateway_custom_domain" "foo" {
	service_id         = "service-7lybgojo"
	sub_domain         = "tic-test.dnsv1.com"
	protocol           = "http"
	net_type           = "OUTER"
	is_default_mapping = "false"
	default_domain     = "service-ohxqslqe-1259649581.gz.apigw.tencentcs.com"
	path_mappings      = ["/good#test","/root#release"]
}

data "tencentcloud_api_gateway_customer_domains" "id" {
	service_id = tencentcloud_api_gateway_custom_domain.foo.service_id 
}
`
}
