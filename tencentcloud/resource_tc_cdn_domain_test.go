package tencentcloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccTencentCloudCdnDomain(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCdnDomainDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCdnDomain,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCdnDomainExists("tencentcloud_cdn_domain.foo"),
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "domain", "test.zhaoshaona.com"),
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "service_type", "web"),
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "area", "mainland"),
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "origin.0.origin_type", "ip"),
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "origin.0.origin_list.#", "1"),
				),
			},
			{
				Config: testAccCdnDomainFull,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "domain", "test.zhaoshaona.com"),
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "service_type", "web"),
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "area", "mainland"),
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "full_url_cache", "false"),
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "origin.0.origin_type", "ip"),
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "origin.0.origin_list.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "origin.0.server_name", "test.zhaoshaona.com"),
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "origin.0.origin_pull_protocol", "follow"),
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "https_config.0.https_switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "https_config.0.http2_switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "https_config.0.ocsp_stapling_switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "https_config.0.spdy_switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "https_config.0.verify_client", "off"),
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "tags.hello", "world"),
				),
			},
			{
				ResourceName:      "tencentcloud_cdn_domain.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckCdnDomainDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	cdnService := CdnService{
		client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_cdn_domain" {
			continue
		}

		domainConfig, err := cdnService.DescribeDomainsConfigByDomain(ctx, rs.Primary.ID)
		if err != nil {
			err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
				domainConfig, err = cdnService.DescribeDomainsConfigByDomain(ctx, rs.Primary.ID)
				if err != nil {
					return retryError(err)
				}
				return nil
			})
		}
		if err != nil {
			return err
		}
		if domainConfig != nil {
			return fmt.Errorf("cdn domain still exists: %s", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckCdnDomainExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("cdn domain %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("cdn domain id is not set")
		}
		cdnService := CdnService{
			client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
		}
		domainConfig, err := cdnService.DescribeDomainsConfigByDomain(ctx, rs.Primary.ID)
		if err != nil {
			err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
				domainConfig, err = cdnService.DescribeDomainsConfigByDomain(ctx, rs.Primary.ID)
				if err != nil {
					return retryError(err)
				}
				return nil
			})
		}
		if err != nil {
			return err
		}
		if domainConfig == nil {
			return fmt.Errorf("cdn domain is not found")
		}
		return nil
	}
}

const testAccCdnDomain = `
resource "tencentcloud_cdn_domain" "foo" {
  domain = "test.zhaoshaona.com"
  service_type = "web"
  area = "mainland"
  origin {
	origin_type = "ip"
	origin_list = ["139.199.199.140"]
  }
}
`

const testAccCdnDomainFull = `
resource "tencentcloud_cdn_domain" "foo" {
  domain = "test.zhaoshaona.com"
  service_type = "web"
  area = "mainland"
  full_url_cache = false

  origin {
	origin_type = "ip"
	origin_list = ["139.199.199.140"]
    server_name = "test.zhaoshaona.com"
    origin_pull_protocol = "follow"
  }

  https_config {
    https_switch = "off"
    http2_switch = "off"
    ocsp_stapling_switch = "off"
    spdy_switch = "off"
	verify_client = "off"
  }

  tags = {
    hello = "world"
  }
}
`
