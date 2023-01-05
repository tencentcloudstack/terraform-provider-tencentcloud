package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"

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

func TestAccTencentCloudCdnDomainResource_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY)
			if err := testAccCdnDomainVerify("www2"); err != nil {
				log.Printf("[TestAccTencentCloudCdnDomainResource] Domain Verify failed: %s", err)
				t.Fatalf("[TestAccTencentCloudCdnDomainResource] Domain Verify failed: %s", err)
			}
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCdnDomainDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCdnDomainBasic,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckCdnDomainExists("tencentcloud_cdn_domain.foo"),
					resource.TestCheckResourceAttrSet("tencentcloud_cdn_domain.foo", "domain"),
					resource.TestCheckResourceAttrSet("tencentcloud_cdn_domain.foo", "service_type"),
					resource.TestCheckResourceAttrSet("tencentcloud_cdn_domain.foo", "area"),
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "origin.0.origin_type", "cos"),
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "origin.0.origin_list.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "authentication.0.switch", "on"),
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "authentication.0.type_a.#", "1"),
					// extends
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "cache_key.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "cache_key.0.full_url_cache", "off"),
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "cache_key.0.ignore_case", "off"),
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "cache_key.0.query_string.0.switch", "on"),
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "cache_key.0.query_string.0.value", "t"),
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "cache_key.0.key_rules.#", "2"),
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "cache_key.0.key_rules.0.rule_paths.0", "/"),
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "cache_key.0.key_rules.1.rule_paths.#", "3"),
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "cache_key.0.key_rules.0.query_string.0.value", "sign;s"),
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "cache_key.0.key_rules.1.query_string.0.value", "t;timestamp"),
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "status_code_cache.0.switch", "on"),
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "status_code_cache.0.cache_rules.#", "2"),
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "status_code_cache.0.cache_rules.0.status_code", "403"),
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "status_code_cache.0.cache_rules.0.cache_time", "5"),
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "status_code_cache.0.cache_rules.1.status_code", "404"),
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "status_code_cache.0.cache_rules.1.cache_time", "10"),
				),
			},
			{
				Config: testAccCdnDomainBasicUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCdnDomainExists("tencentcloud_cdn_domain.foo"),
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "authentication.0.switch", "on"),
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "authentication.0.type_c.#", "1"),
					// extends
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "ip_freq_limit.0.switch", "on"),
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "status_code_cache.0.switch", "on"),
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "status_code_cache.0.cache_rules.#", "2"),
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "status_code_cache.0.cache_rules.0.status_code", "403"),
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "status_code_cache.0.cache_rules.0.cache_time", "10"),
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "status_code_cache.0.cache_rules.1.status_code", "404"),
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "status_code_cache.0.cache_rules.1.cache_time", "10"),
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "post_max_size.0.switch", "on"),
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "post_max_size.0.max_size", "63"),
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "cache_key.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "cache_key.0.full_url_cache", "off"),
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "cache_key.0.ignore_case", "off"),
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "cache_key.0.query_string.0.switch", "on"),
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "cache_key.0.query_string.0.action", "excludeCustom"),
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "cache_key.0.query_string.0.value", "timestamp"),
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "cache_key.0.key_rules.#", "2"),
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "cache_key.0.key_rules.0.rule_paths.0", "/"),
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "cache_key.0.key_rules.1.rule_paths.0", "/vendor.js"),
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "cache_key.0.key_rules.0.query_string.0.switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "cache_key.0.key_rules.1.query_string.0.value", "t"),
				),
			},
			{
				ResourceName:      "tencentcloud_cdn_domain.foo",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"full_url_cache",
					"https_config",
					"ip_filter",
					"ip_freq_limit",
					"status_code_cache",
					"compression",
					"band_width_alert",
					"error_page",
					"response_header",
					"downstream_capping",
					"origin_pull_optimization",
					"post_max_size",
					"referer",
					"max_age",
					"cache_key",
					"aws_private_access",
					"oss_private_access",
					"hw_private_access",
					"qn_private_access",
				},
			},
		},
	})
}

func TestAccTencentCloudCdnDomainDryRun(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCdnDomainDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCndDryRunDomainResource(1),
				Check: func(state *terraform.State) error {
					for _, rs := range state.RootModule().Resources {
						if rs.Type != "tencentcloud_cdn_domain" {
							continue
						}

						request := testAccCdnDryRunAddDomainRequest()
						dryRunResult, ok := rs.Primary.Attributes["dry_run_create_result"]
						if requestJson := request.ToJsonString(); ok && dryRunResult != requestJson {
							log.Printf("[ERROR]: unexpected argument, top is expected, bottom is actual: \n%s\n%s", requestJson, dryRunResult)
							return fmt.Errorf("resource %s.dry_run is no match request json", rs.Type)
						}
					}

					return nil
				},
			},
			{
				ExpectNonEmptyPlan: true,
				Config:             testAccCndDryRunDomainResource(2),
				Check: func(state *terraform.State) error {
					for _, rs := range state.RootModule().Resources {
						if rs.Type != "tencentcloud_cdn_domain" {
							continue
						}

						request := cdn.NewUpdateDomainConfigRequest()
						request.IpFilter = &cdn.IpFilter{
							Switch: helper.String("off"),
						}
						request.ErrorPage = &cdn.ErrorPage{
							Switch: helper.String("on"),
							PageRules: []*cdn.ErrorPageRule{
								{
									StatusCode:   helper.IntInt64(403),
									RedirectCode: helper.IntInt64(302),
									RedirectUrl:  helper.String("https://www.test.com/error3.html"),
								},
							},
						}
						dryRunResult, ok := rs.Primary.Attributes["dry_run_update_result"]
						if requestJson := request.ToJsonString(); ok && dryRunResult != requestJson {
							log.Printf("[ERROR]: unexpected argument, top is expected, bottom is actual: \n%s\n%s", requestJson, dryRunResult)
							return fmt.Errorf("resource %s.dry_run_update is no match request json", rs.Type)
						}
					}

					return nil
				},
			},
		},
	})
}

func TestAccTencentCloudCdnDomainResource_HTTPs(t *testing.T) {
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
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cdn_domain.foo", "domain"),
					resource.TestCheckResourceAttrSet("tencentcloud_cdn_domain.foo", "service_type"),
					resource.TestCheckResourceAttrSet("tencentcloud_cdn_domain.foo", "area"),
					//resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "full_url_cache", "false"),
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
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "rule_cache.1.follow_origin_switch", "on"),
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "rule_cache.1.heuristic_cache_switch", "on"),
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "rule_cache.1.heuristic_cache_time", "3600"),
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
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "https_config.0.tls_versions.0", "TLSv1"),
				),
			},
			{
				Config: testAccCdnDomainFullUpdate,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cdn_domain.foo", "domain"),
					resource.TestCheckResourceAttrSet("tencentcloud_cdn_domain.foo", "service_type"),
					resource.TestCheckResourceAttrSet("tencentcloud_cdn_domain.foo", "area"),
					//resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "full_url_cache", "false"),
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
					resource.TestCheckResourceAttr("tencentcloud_cdn_domain.foo", "https_config.0.tls_versions.#", "2"),
				),
			},
			{
				ResourceName:      "tencentcloud_cdn_domain.foo",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"https_config",
					"authentication",
					"full_url_cache",
					"cache_key",
				},
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
	for i := range domains {
		item := domains[i]
		if v := *item.DomainName; strings.HasPrefix(v, "tencent") {
			return v, nil
		}
	}
	return "", fmt.Errorf("no available domain")
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
  domains = data.tencentcloud_domains.domains.list.*.domain_name
  domain = [for i in local.domains: i if length(regexall("^tencent", i)) > 0][0]
  bucket_url = "keep-cdn-test-${data.tencentcloud_user_info.info.app_id}.cos.ap-singapore.myqcloud.com"
}
`

const testAccCdnDomainBasic = testAccDomainCosForCDN + `

resource "tencentcloud_cdn_domain" "foo" {
  domain = "www2.${local.domain}"
  service_type = "web"
  area = "overseas"
  origin {
	origin_type          = "cos"
	origin_list          = [local.bucket_url]
	server_name			 = local.bucket_url
    origin_pull_protocol = "follow"
  }
  cache_key {
    full_url_cache = "off"
    ignore_case = "off"
    query_string {
      switch = "on"
      action = "includeCustom"
      value = "t"
    }
    key_rules {
      rule_paths = ["/"]
      rule_type = "index"
      full_url_cache = "off"
      ignore_case = "off"
      query_string {
        action = "excludeCustom"
        switch = "on"
        value = "sign;s"
      }
    }
    key_rules {
      rule_paths = ["css", "js", "map"]
      rule_type = "file"
      full_url_cache = "off"
      ignore_case = "off"
      query_string {
        action = "excludeCustom"
        switch = "on"
        value = "t;timestamp"
      }
    }
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
  status_code_cache {
    switch = "on"
    cache_rules {
      status_code = "403"
      cache_time = 5
    }
    cache_rules {
      status_code = "404"
      cache_time = 10
    }
  }
}
`

const testAccCdnDomainBasicUpdate = testAccDomainCosForCDN + `

resource "tencentcloud_cdn_domain" "foo" {
  domain = "www2.${local.domain}"
  service_type = "web"
  area = "overseas"
  origin {
	origin_type          = "cos"
	origin_list          = [local.bucket_url]
	server_name			 = local.bucket_url
    origin_pull_protocol = "follow"
  }
  cache_key {
    full_url_cache = "off"
    ignore_case = "off"
    query_string {
      switch = "on"
      action = "excludeCustom"
      value = "timestamp"
    }
    key_rules {
      rule_paths = ["/"]
      rule_type = "index"
      query_string {
        switch = "off"
      }
    }
    key_rules {
      rule_paths = ["/vendor.js"]
      rule_type = "path"
      full_url_cache = "off"
      ignore_case = "off"
      query_string {
        action = "includeCustom"
        switch = "on"
        value = "t"
      }
    }
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
  ip_freq_limit {
    switch = "on"
    qps = 50
  }
  status_code_cache {
    switch = "on"
    cache_rules {
      status_code = "403"
      cache_time = 10
    }
    cache_rules {
      status_code = "404"
      cache_time = 10
    }
  }
  post_max_size {
    switch = "on"
    max_size = 63
  }
}
`

func testAccCndDryRunDomainResource(stage int) string {
	ipFilter := `
switch = "on"
filter_type = "blacklist"
filters = ["10.0.3.0/24"]
filter_rules {
  filter_type = "whitelist"
  filters = ["10.0.2.0/24"]
  rule_type = "all"
  rule_paths = ["*"]
}
filter_rules {
  filter_type = "blacklist"
  filters = ["10.0.1.0/24"]
  rule_type = "file"
  rule_paths = ["txt"]
}
return_code = 404
	`
	errorPage := `
switch = "on"
page_rules {
  status_code = 403
  redirect_code = 302
  redirect_url = "https://www.test.com/error.html"
}
page_rules {
  status_code = 404
  redirect_code = 301
  redirect_url = "https://www.test2.com/error2.html"
}
	`

	if stage == 2 {
		ipFilter = `switch = "off"`
		errorPage = `
switch = "on"
page_rules {
  status_code = 403
  redirect_code = 302
  redirect_url = "https://www.test.com/error3.html"
}
`
	}
	return fmt.Sprintf(`
%s
resource "tencentcloud_cdn_domain" "dry_run" {

  explicit_using_dry_run = true

  domain = "www.example1.com"
  service_type = "web"
  area = "overseas"
  origin {
    origin_type          = "cos"
    origin_list          = ["test-1231231231.mycloud.com"]
    server_name			 = "test-1231231231.mycloud.com"
    origin_pull_protocol = "follow"
  }
  ip_filter { %s }
  ip_freq_limit {
    switch = "on"
    qps = 50
  }
  status_code_cache {
    switch = "on"
    cache_rules {
      status_code = "403"
      cache_time = 5
    }
    cache_rules {
      status_code = "404"
      cache_time = 10
    }
  }
  compression {
    switch = "on"
    compression_rules {
      compress = true
      min_length = 200
      max_length = 2000
      algorithms = ["gzip"]
      file_extensions = ["png", "js"]
      rule_type = "directory"
      rule_paths = ["/assets"]
    }
    compression_rules {
      compress = true
      min_length = 300
      max_length = 4000
      algorithms = ["gzip"]
      file_extensions = [".css", "js"]
      rule_type = "directory"
      rule_paths = ["/assets2"]
    }
  }
  band_width_alert {
    switch = "on"
    bps_threshold = 1024*1024*30
    counter_measure = "RETURN_404"
    alert_switch = "on"
    alert_percentage = 90
    metric = "bandwidth"
  }
  error_page { %s }
  response_header {
    switch = "on"
    header_rules {
      header_mode = "set"
      header_name = "x-foo-header"
      header_value = "v1"
      rule_type = "directory"
      rule_paths = ["/xxx/test/"]
    }
    header_rules {
      header_mode = "add"
      header_name = "x-bar-header"
      header_value = "v2"
      rule_type = "file"
      rule_paths = ["jpg"]
    }
  }
  downstream_capping {
    switch = "on"
    capping_rules {
      rule_type = "all"
      rule_paths = ["*"]
      kbps_threshold = 300
    }
    capping_rules {
      rule_type = "path"
      rule_paths = ["/data/d.html"]
      kbps_threshold = 400
    }
  }
  origin_pull_optimization {
    switch = "on"
    optimization_type = "CNToOV"
  }
  referer {
    switch = "on"
    referer_rules {
      rule_type = "all"
      rule_paths = ["*"]
      referer_type = "blacklist"
      referers = ["www.example123.com"]
      allow_empty = true
    }
  }
  max_age {
    switch = "on"
    max_age_rules {
      max_age_type = "all"
      max_age_contents = ["*"]
      max_age_time = 3600
      follow_origin = "on"
    }
    max_age_rules {
      max_age_type = "path"
      max_age_contents = ["/a/b.html"]
      max_age_time = 7200
      follow_origin = "off"
    }
  }

  specific_config_mainland = <<-EOF
{
  "Seo": { "Switch": "on" }
}
EOF
  specific_config_overseas = <<-EOF
{
  "FollowRedirect": { "Switch": "on" }
}
EOF

  origin_pull_timeout {
    connect_timeout = 10
    receive_timeout = 10
  }

  response_header_cache_switch = "on"
  seo_switch = "on"
  video_seek_switch = "on"
  offline_cache_switch = "on"
  quic_switch = "on"
}
`, testAccDomainCosForCDN, ipFilter, errorPage)
}

func testAccCdnDryRunAddDomainRequest() *cdn.AddCdnDomainRequest {
	on := helper.String("on")
	off := helper.String("off")
	request := &cdn.AddCdnDomainRequest{
		Domain:      helper.String("www.example1.com"),
		ServiceType: helper.String("web"),
		Area:        helper.String("overseas"),
		ProjectId:   helper.IntInt64(0),
		Origin: &cdn.Origin{
			OriginType: helper.String("cos"),
			Origins: []*string{
				helper.String("test-1231231231.mycloud.com"),
			},
			ServerName:         helper.String("test-1231231231.mycloud.com"),
			OriginPullProtocol: helper.String("follow"),
			CosPrivateAccess:   off,
		},
		IpFilter: &cdn.IpFilter{
			Switch:     on,
			ReturnCode: helper.IntInt64(404),
			FilterType: helper.String("blacklist"),
			Filters: []*string{
				helper.String("10.0.3.0/24"),
			},
			FilterRules: []*cdn.IpFilterPathRule{
				{
					FilterType: helper.String("whitelist"),
					Filters: []*string{
						helper.String("10.0.2.0/24"),
					},
					RuleType: helper.String("all"),
					RulePaths: []*string{
						helper.String("*"),
					},
				},
				{
					FilterType: helper.String("blacklist"),
					Filters: []*string{
						helper.String("10.0.1.0/24"),
					},
					RuleType: helper.String("file"),
					RulePaths: []*string{
						helper.String("txt"),
					},
				},
			},
		},
		IpFreqLimit: &cdn.IpFreqLimit{
			Switch: on,
			Qps:    helper.IntInt64(50),
		},
		StatusCodeCache: &cdn.StatusCodeCache{
			Switch: on,
			CacheRules: []*cdn.StatusCodeCacheRule{
				{
					StatusCode: helper.String("403"),
					CacheTime:  helper.IntInt64(5),
				},
				{
					StatusCode: helper.String("404"),
					CacheTime:  helper.IntInt64(10),
				},
			},
		},
		Compression: &cdn.Compression{
			Switch: on,
			CompressionRules: []*cdn.CompressionRule{
				{
					Compress:  helper.Bool(true),
					MinLength: helper.IntInt64(200),
					MaxLength: helper.IntInt64(2000),
					Algorithms: []*string{
						helper.String("gzip"),
					},
					FileExtensions: []*string{
						helper.String("png"),
						helper.String("js"),
					},
					RuleType: helper.String("directory"),
					RulePaths: []*string{
						helper.String("/assets"),
					},
				},
				{
					Compress:  helper.Bool(true),
					MinLength: helper.IntInt64(300),
					MaxLength: helper.IntInt64(4000),
					Algorithms: []*string{
						helper.String("gzip"),
					},
					FileExtensions: []*string{
						helper.String(".css"),
						helper.String("js"),
					},
					RuleType: helper.String("directory"),
					RulePaths: []*string{
						helper.String("/assets2"),
					},
				},
			},
		},
		BandwidthAlert: &cdn.BandwidthAlert{
			Switch:          on,
			BpsThreshold:    helper.IntInt64(1024 * 1024 * 30),
			CounterMeasure:  helper.String("RETURN_404"),
			AlertSwitch:     on,
			AlertPercentage: helper.IntInt64(90),
			Metric:          helper.String("bandwidth"),
		},
		ErrorPage: &cdn.ErrorPage{
			Switch: on,
			PageRules: []*cdn.ErrorPageRule{
				{
					StatusCode:   helper.IntInt64(403),
					RedirectCode: helper.IntInt64(302),
					RedirectUrl:  helper.String("https://www.test.com/error.html"),
				},
				{
					StatusCode:   helper.IntInt64(404),
					RedirectCode: helper.IntInt64(301),
					RedirectUrl:  helper.String("https://www.test2.com/error2.html"),
				},
			},
		},
		ResponseHeader: &cdn.ResponseHeader{
			Switch: on,
			HeaderRules: []*cdn.HttpHeaderPathRule{
				{
					HeaderMode:  helper.String("set"),
					HeaderName:  helper.String("x-foo-header"),
					HeaderValue: helper.String("v1"),
					RuleType:    helper.String("directory"),
					RulePaths: []*string{
						helper.String("/xxx/test/"),
					},
				},
				{
					HeaderMode:  helper.String("add"),
					HeaderName:  helper.String("x-bar-header"),
					HeaderValue: helper.String("v2"),
					RuleType:    helper.String("file"),
					RulePaths: []*string{
						helper.String("jpg"),
					},
				},
			},
		},
		DownstreamCapping: &cdn.DownstreamCapping{
			Switch: on,
			CappingRules: []*cdn.CappingRule{
				{
					RuleType: helper.String("all"),
					RulePaths: []*string{
						helper.String("*"),
					},
					KBpsThreshold: helper.IntInt64(300),
				},
				{
					RuleType: helper.String("path"),
					RulePaths: []*string{
						helper.String("/data/d.html"),
					},
					KBpsThreshold: helper.IntInt64(400),
				},
			},
		},
		OriginPullOptimization: &cdn.OriginPullOptimization{
			Switch:           on,
			OptimizationType: helper.String("CNToOV"),
		},
		Referer: &cdn.Referer{
			Switch: on,
			RefererRules: []*cdn.RefererRule{
				{
					RuleType: helper.String("all"),
					RulePaths: []*string{
						helper.String("*"),
					},
					RefererType: helper.String("blacklist"),
					Referers: []*string{
						helper.String("www.example123.com"),
					},
					AllowEmpty: helper.Bool(true),
				},
			},
		},
		MaxAge: &cdn.MaxAge{
			Switch: on,
			MaxAgeRules: []*cdn.MaxAgeRule{
				{
					MaxAgeType: helper.String("all"),
					MaxAgeContents: []*string{
						helper.String("*"),
					},
					MaxAgeTime:   helper.IntInt64(3600),
					FollowOrigin: on,
				},
				{
					MaxAgeType: helper.String("path"),
					MaxAgeContents: []*string{
						helper.String("/a/b.html"),
					},
					MaxAgeTime:   helper.IntInt64(7200),
					FollowOrigin: off,
				},
			},
		},
		SpecificConfig: &cdn.SpecificConfig{
			Overseas: &cdn.OverseaConfig{
				FollowRedirect: &cdn.FollowRedirect{Switch: on},
			},
			Mainland: &cdn.MainlandConfig{
				Seo: &cdn.Seo{Switch: on},
			},
		},
		OriginPullTimeout: &cdn.OriginPullTimeout{
			ConnectTimeout: helper.IntUint64(10),
			ReceiveTimeout: helper.IntUint64(10),
		},
		ResponseHeaderCache: &cdn.ResponseHeaderCache{Switch: on},
		Seo:                 &cdn.Seo{Switch: on},
		VideoSeek:           &cdn.VideoSeek{Switch: on},
		OfflineCache:        &cdn.OfflineCache{Switch: on},
		Quic:                &cdn.Quic{Switch: on},
		Ipv6Access:          &cdn.Ipv6Access{Switch: off},
		RangeOriginPull:     &cdn.RangeOriginPull{Switch: on}, // by argument default
		FollowRedirect:      &cdn.FollowRedirect{Switch: off}, // by argument default
		CacheKey:            &cdn.CacheKey{FullUrlCache: on},  // by argument default
	}
	return request
}

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
  range_origin_switch = "off"
  
  rule_cache {
	cache_time = 10000
	no_cache_switch="on"
	re_validate="on"
  }
  rule_cache {
	cache_time = 10000
    follow_origin_switch="on"
    rule_paths=["jpg", "png"]
	rule_type="file"
    heuristic_cache_switch = "on"
	heuristic_cache_time = 3600
  }

  request_header {
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
	tls_versions = ["TLSv1"]
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

  rule_cache {
    cache_time = 10000
    rule_paths=["jpg", "png"]
	rule_type="file"
    follow_origin_switch = "on"
    heuristic_cache_switch = "off"
  }

  request_header {
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
	tls_versions = ["TLSv1", "TLSv1.1"]
  }

  tags = {
    hello = "world"
  }
}
`
