package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"testing"

	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"

	dnspod "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dnspod/v20210323"
	domain "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/domain/v20180808"

	cdn "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdn/v20180606"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func init() {
	// go test -v ./tencentcloud -sweep=ap-guangzhou -sweep-run=cdn_domain
	resource.AddTestSweepers("cdn_domain", &resource.Sweeper{
		Name: "cdn_domain",
		F: func(r string) error {
			logId := getLogId(contextNil)
			ctx := context.WithValue(context.TODO(), logIdKey, logId)
			cli, _ := sharedClientForRegion(r)
			client := cli.(*TencentCloudClient).apiV3Conn

			service := CdnService{client}
			domains, err := service.DescribeDomainsConfigByFilters(ctx, nil)
			if err != nil {
				return err
			}

			for i := range domains {
				item := domains[i]
				name := *item.Domain

				if isResourcePersist(name, nil) {
					continue
				}

				if *item.Status != "offline" {
					_ = service.StopDomain(ctx, name)
				}

				err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
					inErr := service.DeleteDomain(ctx, name)
					if inErr != nil {
						retryError(err, cdn.RESOURCEUNAVAILABLE_CDNHOSTISNOTOFFLINE)
					}
					return nil
				})
				if err != nil {
					break
				}
			}

			return nil
		},
	})
}

func TestAccTencentCloudCdnDomainResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY)
			if err := testAccCdnDomainVerify("www"); err != nil {
				log.Printf("[TestAccTencentCloudCdnDomainResource] Domain Verify failed: %s", err)
				t.Fatalf("[TestAccTencentCloudCdnDomainResource] Domain Verify failed: %s", err)
			}
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCdnDomainDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCdnDomainBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCdnDomainExists("tencentcloud_cdn_domain.foo"),
					resource.TestCheckResourceAttrSet("tencentcloud_cdn_domain.foo", "domain"),
					resource.TestCheckResourceAttrSet("tencentcloud_cdn_domain.foo", "service_type"),
					resource.TestCheckResourceAttrSet("tencentcloud_cdn_domain.foo", "area"),
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "origin.0.origin_type", "cos"),
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "origin.0.origin_list.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "authentication.0.switch", "on"),
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "authentication.0.type_a.#", "1"),
				),
			},
			{
				Config: testAccCdnDomainBasicUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCdnDomainExists("tencentcloud_cdn_domain.foo"),
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "authentication.0.switch", "on"),
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "authentication.0.type_c.#", "1"),
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
		PreCheck: func() {
			testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY)
			if err := testAccCdnDomainVerify("c"); err != nil {
				log.Printf("[TestAccCentcentCloudCdnDomainWithHTTPs] Domain Verify failed: %s", err)
				t.Fatalf("[TestAccCentcentCloudCdnDomainWithHTTPs] Domain Verify failed: %s", err)
			}
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCdnDomainDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCdnDomainFull,
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
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "https_config.0.force_redirect.0.switch", "on"),
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "https_config.0.force_redirect.0.redirect_type", "https"),
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "https_config.0.force_redirect.0.redirect_status_code", "302"),
				),
			},
			{
				Config: testAccCdnDomainFullUpdate,
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
				ImportStateVerifyIgnore: []string{"https_config", "authentication"},
			},
		},
	})
}

func testAccGetTestingDomain() (string, error) {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	cli, _ := sharedClientForRegion("ap-guangzhou")
	client := cli.(*TencentCloudClient).apiV3Conn
	service := DomainService{client}
	request := domain.NewDescribeDomainNameListRequest()
	domains, err := service.DescribeDomainNameList(ctx, request)
	if err != nil {
		return "", err
	}
	if len(domains) == 0 {
		return "", nil
	}
	return *domains[0].DomainName, nil
}

func testAccCdnDomainVerify(domainPrefix string) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	cli, _ := sharedClientForRegion("ap-guangzhou")
	client := cli.(*TencentCloudClient).apiV3Conn
	service := CdnService{client}
	continueCode := []string{
		// no record
		cdn.UNAUTHORIZEDOPERATION_CDNDOMAINRECORDNOTVERIFIED,
		// has record but need modify
		cdn.UNAUTHORIZEDOPERATION_CDNTXTRECORDVALUENOTMATCH,
	}

	domainName, err := testAccGetTestingDomain()
	l3domain := fmt.Sprintf("%s.%s", domainPrefix, domainName)

	if err != nil {
		return err
	}

	result, err := service.VerifyDomainRecord(ctx, l3domain)

	if err != nil {

		code := err.(*sdkErrors.TencentCloudSDKError).Code

		if !IsContains(continueCode, code) {
			return err
		}
	}

	log.Printf("[Precheck] Domain Record Verify: %t", result)
	if result {
		return nil
	}

	cRes, err := service.CreateVerifyRecord(ctx, l3domain)
	if err != nil {
		return err
	}

	recordType := *cRes.RecordType
	record := *cRes.Record

	err = testAccSetDnsPodRecord(domainName, recordType, record)

	if err != nil {
		return err
	}

	err = resource.Retry(readRetryTimeout*3, func() *resource.RetryError {
		result, err = service.VerifyDomainRecord(ctx, l3domain)
		if err != nil {
			return retryError(err, continueCode...)
		}
		if result {
			return nil
		}
		return resource.RetryableError(fmt.Errorf("verifying domain, retry"))
	})

	return nil
}

func testAccSetDnsPodRecord(domainName, recordType, record string) error {
	cli, _ := sharedClientForRegion("ap-guangzhou")
	client := cli.(*TencentCloudClient).apiV3Conn
	recordLine := "默认"
	subDomain := "_cdnauth"

	request := dnspod.NewDescribeRecordListRequest()
	request.Domain = &domainName
	request.Subdomain = &subDomain
	response, err := client.UseDnsPodClient().DescribeRecordList(request)

	if err != nil {
		code := err.(*sdkErrors.TencentCloudSDKError).Code
		if code != dnspod.RESOURCENOTFOUND_NODATAOFRECORD {
			return err
		}
	}

	if response.Response != nil && len(response.Response.RecordList) > 0 {
		for i := range response.Response.RecordList {
			recordInfo := response.Response.RecordList[i]
			if *recordInfo.Value == record {
				return nil
			}
		}
	}

	cRequest := dnspod.NewCreateRecordRequest()
	cRequest.Domain = &domainName
	cRequest.SubDomain = &subDomain
	cRequest.RecordType = &recordType
	cRequest.RecordLine = &recordLine
	cRequest.Value = &record

	if _, err := client.UseDnsPodClient().CreateRecord(cRequest); err != nil {
		return err
	}
	return nil
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

const testAccDomainCosForCDN = `
data "tencentcloud_domains" "domains" {}
data "tencentcloud_user_info" "info" {}

locals {
  domain = data.tencentcloud_domains.domains.list.0.domain_name
  bucket_url = "keep-cdn-test-${data.tencentcloud_user_info.info.app_id}.cos.ap-singapore.myqcloud.com"
}
`

const testAccCdnDomainBasic = testAccDomainCosForCDN + `

resource "tencentcloud_cdn_domain" "foo" {
  domain = "www.${local.domain}"
  service_type = "web"
  area = "overseas"
  origin {
	origin_type          = "cos"
	origin_list          = [local.bucket_url]
	server_name			 = local.bucket_url
    origin_pull_protocol = "follow"
  }
  authentication {
    switch = "on"
    type_a {
      secret_key = "kXNgx2625Rre"
      expire_time = 60
      sign_param = "sign"
      file_extensions = ["/test/1.jpg"]
      filter_type = "whitelist"
      backup_secret_key = "ujU4vsH3jbLzAg2DTUoTLj"
    }
  }
}
`

const testAccCdnDomainBasicUpdate = testAccDomainCosForCDN + `

resource "tencentcloud_cdn_domain" "foo" {
  domain = "www.${local.domain}"
  service_type = "web"
  area = "overseas"
  origin {
	origin_type          = "cos"
	origin_list          = [local.bucket_url]
	server_name			 = local.bucket_url
    origin_pull_protocol = "follow"
  }
  authentication {
    switch = "on"
    type_c {
      secret_key = "MLjUov41L9aXO311"
      expire_time = 60
      file_extensions = ["/test/1.jpg"]
      filter_type = "whitelist"
      time_format = "dec"
      backup_secret_key = "BQ6yVr"
    }
  }
}
`

const testAccSSLForCDN = `
data "tencentcloud_ssl_certificates" "foo" {
  name = "keep-c-ssl"
}

locals {
  certId = data.tencentcloud_ssl_certificates.foo.certificates.0.id
}
`

const testAccCdnDomainFull = testAccDomainCosForCDN + testAccSSLForCDN + `
resource "tencentcloud_cdn_domain" "foo" {
  domain         = "c.${local.domain}"
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
`

const testAccCdnDomainFullUpdate = testAccDomainCosForCDN + testAccSSLForCDN + `
resource "tencentcloud_cdn_domain" "foo" {
  domain         = "c.${local.domain}"
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
`
