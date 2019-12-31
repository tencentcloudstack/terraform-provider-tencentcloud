package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"strings"
	"testing"
)

var testDayuCCHttpPolicyResourceName = "tencentcloud_dayu_cc_http_policy"
var testDayuCCHttpPolicyResourceKey = testDayuCCHttpPolicyResourceName + ".test_policy"

func TestAccTencentCloudDayuCCHttpPolicyResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDayuCCHttpPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDayuCCHttpPolicy,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDayuCCHttpPolicyExists(testDayuCCHttpPolicyResourceKey),
					resource.TestCheckResourceAttrSet(testDayuCCHttpPolicyResourceKey, "create_time"),
					resource.TestCheckResourceAttrSet(testDayuCCHttpPolicyResourceKey, "policy_id"),
					resource.TestCheckResourceAttr(testDayuCCHttpPolicyResourceKey, "resource_type", "bgpip"),
					resource.TestCheckResourceAttr(testDayuCCHttpPolicyResourceKey, "name", "policy_match"),
					resource.TestCheckResourceAttr(testDayuCCHttpPolicyResourceKey, "smode", "matching"),
					resource.TestCheckResourceAttr(testDayuCCHttpPolicyResourceKey, "exe_mode", "drop"),
					resource.TestCheckResourceAttr(testDayuCCHttpPolicyResourceKey, "rule_list.#", "1"),
				),
			},
			{
				Config: testAccDayuCCHttpPolicyUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDayuCCHttpPolicyExists(testDayuCCHttpPolicyResourceKey),
					resource.TestCheckResourceAttrSet(testDayuCCHttpPolicyResourceKey, "create_time"),
					resource.TestCheckResourceAttr(testDayuCCHttpPolicyResourceKey, "resource_type", "bgpip"),
					resource.TestCheckResourceAttr(testDayuCCHttpPolicyResourceKey, "name", "policy_limit"),
					resource.TestCheckResourceAttr(testDayuCCHttpPolicyResourceKey, "smode", "speedlimit"),
					resource.TestCheckResourceAttr(testDayuCCHttpPolicyResourceKey, "frequency", "100"),
				),
			},
		},
	})
}

func TestAccTencentCloudDayuCCHttpPolicyResource_BGP(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDayuCCHttpPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDayuCCHttpPolicy_BGP,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDayuCCHttpPolicyExists(testDayuCCHttpPolicyResourceKey),
					resource.TestCheckResourceAttrSet(testDayuCCHttpPolicyResourceKey, "create_time"),
					resource.TestCheckResourceAttrSet(testDayuCCHttpPolicyResourceKey, "policy_id"),
					resource.TestCheckResourceAttr(testDayuCCHttpPolicyResourceKey, "resource_type", "bgp"),
					resource.TestCheckResourceAttr(testDayuCCHttpPolicyResourceKey, "name", "policy_match"),
					resource.TestCheckResourceAttr(testDayuCCHttpPolicyResourceKey, "smode", "matching"),
					resource.TestCheckResourceAttr(testDayuCCHttpPolicyResourceKey, "exe_mode", "alg"),
					resource.TestCheckResourceAttr(testDayuCCHttpPolicyResourceKey, "rule_list.#", "1"),
				),
			},
		},
	})
}

func TestAccTencentCloudDayuCCHttpPolicyResource_NET(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDayuCCHttpPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDayuCCHttpPolicy_NET,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDayuCCHttpPolicyExists(testDayuCCHttpPolicyResourceKey),
					resource.TestCheckResourceAttrSet(testDayuCCHttpPolicyResourceKey, "create_time"),
					resource.TestCheckResourceAttrSet(testDayuCCHttpPolicyResourceKey, "policy_id"),
					resource.TestCheckResourceAttr(testDayuCCHttpPolicyResourceKey, "resource_type", "net"),
					resource.TestCheckResourceAttr(testDayuCCHttpPolicyResourceKey, "name", "policy_match"),
					resource.TestCheckResourceAttr(testDayuCCHttpPolicyResourceKey, "smode", "matching"),
					resource.TestCheckResourceAttr(testDayuCCHttpPolicyResourceKey, "exe_mode", "drop"),
					resource.TestCheckResourceAttr(testDayuCCHttpPolicyResourceKey, "rule_list.#", "1"),
				),
			},
		},
	})
}

func TestAccTencentCloudDayuCCHttpPolicyResource_BGPMUL(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDayuCCHttpPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDayuCCHttpPolicy_BGPMUL,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDayuCCHttpPolicyExists(testDayuCCHttpPolicyResourceKey),
					resource.TestCheckResourceAttrSet(testDayuCCHttpPolicyResourceKey, "create_time"),
					resource.TestCheckResourceAttrSet(testDayuCCHttpPolicyResourceKey, "policy_id"),
					resource.TestCheckResourceAttr(testDayuCCHttpPolicyResourceKey, "resource_type", "bgp-multip"),
					resource.TestCheckResourceAttr(testDayuCCHttpPolicyResourceKey, "name", "policy_match"),
					resource.TestCheckResourceAttr(testDayuCCHttpPolicyResourceKey, "smode", "matching"),
					resource.TestCheckResourceAttr(testDayuCCHttpPolicyResourceKey, "exe_mode", "alg"),
					resource.TestCheckResourceAttr(testDayuCCHttpPolicyResourceKey, "rule_list.#", "1"),
				),
			},
		},
	})
}
func testAccCheckDayuCCHttpPolicyDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != testDayuCCHttpPolicyResourceName {
			continue
		}
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), "logId", logId)

		items := strings.Split(rs.Primary.ID, FILED_SP)
		if len(items) < 3 {
			return fmt.Errorf("broken ID of DDos policy case")
		}
		resourceType := items[0]
		resourceId := items[1]
		policyId := items[2]

		service := DayuService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}

		_, has, err := service.DescribeCCSelfdefinePolicy(ctx, resourceType, resourceId, policyId)
		if err != nil {
			_, has, err = service.DescribeCCSelfdefinePolicy(ctx, resourceType, resourceId, policyId)
		}
		if err != nil {
			return err
		}
		if !has {
			return nil
		} else {
			return fmt.Errorf("delete CC http policy %s fail, still on server", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckDayuCCHttpPolicyExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource %s is not found", n)
		}
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), "logId", logId)

		items := strings.Split(rs.Primary.ID, FILED_SP)
		if len(items) < 3 {
			return fmt.Errorf("broken ID of DDos policy case")
		}
		resourceType := items[0]
		resourceId := items[1]
		policyId := items[2]

		service := DayuService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}

		_, has, err := service.DescribeCCSelfdefinePolicy(ctx, resourceType, resourceId, policyId)
		if err != nil {
			_, has, err = service.DescribeCCSelfdefinePolicy(ctx, resourceType, resourceId, policyId)
		}
		if err != nil {
			return err
		}
		if has {
			return nil
		} else {
			return fmt.Errorf("CC http policy %s not found on server", rs.Primary.ID)

		}
	}
}

const testAccDayuCCHttpPolicy string = `
resource "tencentcloud_dayu_cc_http_policy" "test_policy" {
  resource_type         = "bgpip"
  resource_id 			= "bgpip-00000294"
  name					= "policy_match"
  smode					= "matching"
  exe_mode				= "drop"
  switch				= true
  rule_list {
	skey 				= "host"
	operator			= "include"
	value				= "123"
	}
}
`
const testAccDayuCCHttpPolicyUpdate string = `
resource "tencentcloud_dayu_cc_http_policy" "test_policy" {
  resource_type         = "bgpip"
  resource_id 			= "bgpip-00000294"
  name					= "policy_limit"
  smode					= "speedlimit"
  switch				= true
  frequency				= 100
}
`
const testAccDayuCCHttpPolicy_NET string = `
resource "tencentcloud_dayu_cc_http_policy" "test_policy" {
  resource_type         = "net"
  resource_id 			= "net-0000007e"
  name					= "policy_match"
  smode					= "matching"
  exe_mode				= "drop"
  switch				= true
  rule_list {
	skey 				= "cgi"
	operator			= "equal"
	value				= "123"
	}
}
`

const testAccDayuCCHttpPolicy_BGPMUL string = `
resource "tencentcloud_dayu_cc_http_policy" "test_policy" {
  resource_type         = "bgp-multip"
  resource_id 			= "bgp-0000008o"
  name					= "policy_match"
  smode					= "matching"
  exe_mode				= "alg"
  switch				= true
  ip					= "111.230.178.25"

  rule_list {
	skey 				= "referer"
	operator			= "not_include"
	value				= "123"
	}
}
`

const testAccDayuCCHttpPolicy_BGP string = `
resource "tencentcloud_dayu_cc_http_policy" "test_policy" {
  resource_type         = "bgp"
  resource_id 			= "bgp-000006mq"
  name					= "policy_match"
  smode					= "matching"
  exe_mode				= "alg"
  switch				= true

  rule_list {
	skey 				= "ua"
	operator			= "not_include"
	value				= "123"
	}
}
`
