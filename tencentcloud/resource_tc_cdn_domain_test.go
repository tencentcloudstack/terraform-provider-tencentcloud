package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	domain "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/domain/v20180808"
)

var testAccCdnDomain = ""

func init() {
	log.Printf("initialize domain testcase")
	cli, _ := sharedClientForRegion(defaultRegion)
	client := cli.(*TencentCloudClient).apiV3Conn
	request := domain.NewDescribeDomainNameListRequest()
	response, err := client.UseDomainClient().DescribeDomainNameList(request)
	if err != nil {
		log.Printf("[DescribeDomainNameList] error: %s", err.Error())
		return
	}

	domains := response.Response.DomainSet

	if len(domains) == 0 {
		log.Printf("[WARN] no domain on your account")
		return
	}

	testAccCdnDomain = *domains[0].DomainName
}

func TestAccTencentCloudCdnDomainResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCdnDomainDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCdnDomainBasic("www." + testAccCdnDomain),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCdnDomainExists("tencentcloud_cdn_domain.foo"),
					resource.TestCheckResourceAttrSet("tencentcloud_cdn_domain.foo", "domain"),
					resource.TestCheckResourceAttrSet("tencentcloud_cdn_domain.foo", "service_type"),
					resource.TestCheckResourceAttrSet("tencentcloud_cdn_domain.foo", "area"),
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "origin.0.origin_type", "cos"),
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "origin.0.origin_list.#", "1"),
				),
			},
			{
				ResourceName:            "tencentcloud_cdn_domain.foo",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"https_config"},
			},
		},
	})
}

func TestAccTencentCloudCdnDomainWithHTTPs(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCdnDomainDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCdnDomainFull("c." + testAccCdnDomain),
				PreConfig: func() {

				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cdn_domain.foo", "domain"),
					resource.TestCheckResourceAttrSet("tencentcloud_cdn_domain.foo", "service_type"),
					resource.TestCheckResourceAttrSet("tencentcloud_cdn_domain.foo", "area"),
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "full_url_cache", "false"),
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "range_origin_switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "rule_cache.0.cache_time", "10000"),
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "rule_cache.0.rule_paths.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "rule_cache.0.rule_type", "default"),
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "rule_cache.0.switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "rule_cache.0.compare_max_age", "off"),
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "rule_cache.0.ignore_cache_control", "off"),
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "rule_cache.0.ignore_set_cookie", "off"),
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "rule_cache.0.no_cache_switch", "on"),
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "rule_cache.0.re_validate", "on"),
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "rule_cache.0.follow_origin_switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "request_header.0.switch", "on"),
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "request_header.0.header_rules.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "origin.0.origin_type", "ip"),
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "origin.0.origin_list.#", "1"),
					resource.TestCheckResourceAttrSet("tencentcloud_cdn_domain.foo", "origin.0.server_name"),
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "origin.0.origin_pull_protocol", "follow"),
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "https_config.0.https_switch", "on"),
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "https_config.0.http2_switch", "on"),
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "https_config.0.ocsp_stapling_switch", "on"),
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "https_config.0.spdy_switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "https_config.0.verify_client", "off"),
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "https_config.0.server_certificate_config.0.message", "test"),
					resource.TestCheckResourceAttrSet("tencentcloud_cdn_domain.foo", "https_config.0.server_certificate_config.0.deploy_time"),
					resource.TestCheckResourceAttrSet("tencentcloud_cdn_domain.foo", "https_config.0.server_certificate_config.0.expire_time"),
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "tags.hello", "world"),
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "https_config.0.force_redirect.0.switch", "on"),
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "https_config.0.force_redirect.0.redirect_type", "https"),
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "https_config.0.force_redirect.0.redirect_status_code", "302"),
				),
			},
			{
				Config: testAccCdnDomainFullUpdate("c." + testAccCdnDomain),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cdn_domain.foo", "domain"),
					resource.TestCheckResourceAttrSet("tencentcloud_cdn_domain.foo", "service_type"),
					resource.TestCheckResourceAttrSet("tencentcloud_cdn_domain.foo", "area"),
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "full_url_cache", "false"),
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "range_origin_switch", "on"),
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "rule_cache.0.cache_time", "20000"),
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "rule_cache.0.rule_paths.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "rule_cache.0.rule_type", "all"),
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "rule_cache.0.switch", "on"),
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "rule_cache.0.compare_max_age", "on"),
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "rule_cache.0.ignore_cache_control", "on"),
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "rule_cache.0.ignore_set_cookie", "off"),
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "rule_cache.0.no_cache_switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "rule_cache.0.re_validate", "off"),
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "rule_cache.0.follow_origin_switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "request_header.0.switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "origin.0.origin_type", "cos"),
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "origin.0.origin_list.#", "1"),
					resource.TestCheckResourceAttrSet("tencentcloud_cdn_domain.foo", "origin.0.server_name"),
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "origin.0.origin_pull_protocol", "follow"),
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "https_config.0.https_switch", "on"),
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "https_config.0.http2_switch", "on"),
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "https_config.0.ocsp_stapling_switch", "on"),
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "https_config.0.spdy_switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "https_config.0.verify_client", "off"),
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "https_config.0.server_certificate_config.0.message", "test"),
					resource.TestCheckResourceAttrSet("tencentcloud_cdn_domain.foo", "https_config.0.server_certificate_config.0.deploy_time"),
					resource.TestCheckResourceAttrSet("tencentcloud_cdn_domain.foo", "https_config.0.server_certificate_config.0.expire_time"),
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "tags.hello", "world"),
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "https_config.0.force_redirect.0.switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "https_config.0.force_redirect.0.redirect_type", "http"),
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "https_config.0.force_redirect.0.redirect_status_code", "302"),
				),
			},
			{
				ResourceName:            "tencentcloud_cdn_domain.foo",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"https_config"},
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

func testAccCdnDomainBasic(name string) string {
	return fmt.Sprintf(`
resource "tencentcloud_cdn_domain" "foo" {
  domain = "%s"
  service_type = "web"
  area = "overseas"
  origin {
	origin_type = "ip"
	origin_list = ["43.133.14.92"]
  }
}
`, name)
}

const testAccSSLForCDN = `
data "tencentcloud_ssl_certificates" "foo" {
  name = "keep-c-ssl"
}

data "tencentcloud_cos_buckets" "bucket" {
  bucket_prefix = "keep-cdn-test"
}


data "tencentcloud_user_info" "info" {}

locals {
  certId = data.tencentcloud_ssl_certificates.foo.certificates.0.id
  bucket_url = "keep-cdn-test-${data.tencentcloud_user_info.info.app_id}.cos.ap-singapore.myqcloud.com"
}
`

func testAccCdnDomainFull(name string) string {
	return fmt.Sprintf(testAccSSLForCDN+`
resource "tencentcloud_cdn_domain" "foo" {
  domain         = "%[1]v"
  service_type   = "web"
  area           = "overseas"
  full_url_cache = false
  range_origin_switch = "off"
  
  rule_cache{
	cache_time = 10000
	no_cache_switch="on"
	re_validate="on"
  }

  request_header{
	switch = "on"

	header_rules {
		header_mode = "add"
		header_name = "tf-header-name"
		header_value = "tf-header-value"
		rule_type = "all"
		rule_paths = ["*"]
	}
  }

  origin {
	origin_type          = "cos"
	origin_list          = [local.bucket_url]
	server_name			 = local.bucket_url
    origin_pull_protocol = "follow"
  }

  https_config {
    https_switch         = "on"
    http2_switch         = "on"
    ocsp_stapling_switch = "on"
    spdy_switch          = "off"
	verify_client        = "off"

    force_redirect {
        switch               = "on"
        redirect_type        = "https"
        redirect_status_code = 302
    }

	server_certificate_config {
      certificate_id = local.certId
      message = "test"
    }
}

  tags = {
    hello = "world"
  }
}
`, name)
}

func testAccCdnDomainFullUpdate(name string) string {
	return fmt.Sprintf(testAccSSLForCDN+`
resource "tencentcloud_cdn_domain" "foo" {
  domain         = "%[1]v"
  service_type   = "web"
  area           = "overseas"
  full_url_cache = false
  range_origin_switch = "on"

  rule_cache {
	cache_time = 20000
	rule_paths=["*"]
	rule_type="all"
	switch="on"
    compare_max_age="on"
    ignore_cache_control="on"
	ignore_set_cookie="off"
    no_cache_switch="off"
	re_validate="off"
	follow_origin_switch="off"
  }

  request_header{
	switch = "off"
  }

  origin {
	origin_type          = "cos"
	origin_list          = [local.bucket_url]
	server_name			 = local.bucket_url
    origin_pull_protocol = "follow"
  }

  https_config {
    https_switch         = "on"
    http2_switch         = "on"
    ocsp_stapling_switch = "on"
    spdy_switch          = "off"
	verify_client        = "off"

    force_redirect {
        switch               = "off"
    }

	server_certificate_config {
      certificate_id = local.certId
      message = "test"
    }
}

  tags = {
    hello = "world"
  }
}
`, name)
}
