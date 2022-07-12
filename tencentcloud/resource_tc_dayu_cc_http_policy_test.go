package tencentcloud

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var testDayuCCHttpPolicyResourceName = "tencentcloud_dayu_cc_http_policy"
var testDayuCCHttpPolicyResourceKey = testDayuCCHttpPolicyResourceName + ".test_policy"

func TestAccTencentCloudDayuCCHttpPolicyResource(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDayuCCHttpPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDayuCCHttpPolicy, defaultDayuBgpIp),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDayuCCHttpPolicyExists(testDayuCCHttpPolicyResourceKey),
					resource.TestCheckResourceAttrSet(testDayuCCHttpPolicyResourceKey, "create_time"),
					resource.TestCheckResourceAttrSet(testDayuCCHttpPolicyResourceKey, "policy_id"),
					resource.TestCheckResourceAttr(testDayuCCHttpPolicyResourceKey, "resource_type", "bgpip"),
					resource.TestCheckResourceAttr(testDayuCCHttpPolicyResourceKey, "name", "policy_match"),
					resource.TestCheckResourceAttr(testDayuCCHttpPolicyResourceKey, "smode", "matching"),
					resource.TestCheckResourceAttr(testDayuCCHttpPolicyResourceKey, "action", "drop"),
					resource.TestCheckResourceAttr(testDayuCCHttpPolicyResourceKey, "rule_list.#", "1"),
				),
			},
			{
				Config: fmt.Sprintf(testAccDayuCCHttpPolicyUpdate, defaultDayuBgpIp),
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
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDayuCCHttpPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDayuCCHttpPolicy_BGP, defaultDayuBgp),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDayuCCHttpPolicyExists(testDayuCCHttpPolicyResourceKey),
					resource.TestCheckResourceAttrSet(testDayuCCHttpPolicyResourceKey, "create_time"),
					resource.TestCheckResourceAttrSet(testDayuCCHttpPolicyResourceKey, "policy_id"),
					resource.TestCheckResourceAttr(testDayuCCHttpPolicyResourceKey, "resource_type", "bgp"),
					resource.TestCheckResourceAttr(testDayuCCHttpPolicyResourceKey, "name", "policy_match"),
					resource.TestCheckResourceAttr(testDayuCCHttpPolicyResourceKey, "smode", "matching"),
					resource.TestCheckResourceAttr(testDayuCCHttpPolicyResourceKey, "action", "alg"),
					resource.TestCheckResourceAttr(testDayuCCHttpPolicyResourceKey, "rule_list.#", "1"),
				),
			},
		},
	})
}

func TestAccTencentCloudDayuCCHttpPolicyResource_NET(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDayuCCHttpPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDayuCCHttpPolicy_NET, defaultDayuNet),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDayuCCHttpPolicyExists(testDayuCCHttpPolicyResourceKey),
					resource.TestCheckResourceAttrSet(testDayuCCHttpPolicyResourceKey, "create_time"),
					resource.TestCheckResourceAttrSet(testDayuCCHttpPolicyResourceKey, "policy_id"),
					resource.TestCheckResourceAttr(testDayuCCHttpPolicyResourceKey, "resource_type", "net"),
					resource.TestCheckResourceAttr(testDayuCCHttpPolicyResourceKey, "name", "policy_match"),
					resource.TestCheckResourceAttr(testDayuCCHttpPolicyResourceKey, "smode", "matching"),
					resource.TestCheckResourceAttr(testDayuCCHttpPolicyResourceKey, "action", "drop"),
					resource.TestCheckResourceAttr(testDayuCCHttpPolicyResourceKey, "rule_list.#", "1"),
				),
			},
		},
	})
}

func TestAccTencentCloudDayuCCHttpPolicyResource_BGPMUL(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDayuCCHttpPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDayuCCHttpPolicy_BGPMUL, defaultDayuBgpMul),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDayuCCHttpPolicyExists(testDayuCCHttpPolicyResourceKey),
					resource.TestCheckResourceAttrSet(testDayuCCHttpPolicyResourceKey, "create_time"),
					resource.TestCheckResourceAttrSet(testDayuCCHttpPolicyResourceKey, "policy_id"),
					resource.TestCheckResourceAttr(testDayuCCHttpPolicyResourceKey, "resource_type", "bgp-multip"),
					resource.TestCheckResourceAttr(testDayuCCHttpPolicyResourceKey, "name", "policy_match"),
					resource.TestCheckResourceAttr(testDayuCCHttpPolicyResourceKey, "smode", "matching"),
					resource.TestCheckResourceAttr(testDayuCCHttpPolicyResourceKey, "action", "alg"),
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
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

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
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

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
  resource_id 			= "%s"
  name					= "policy_match"
  smode					= "matching"
  action				= "drop"
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
  resource_id 			= "%s"
  name					= "policy_limit"
  smode					= "speedlimit"
  switch				= true
  frequency				= 100
}
`
const testAccDayuCCHttpPolicy_NET string = `
resource "tencentcloud_dayu_cc_http_policy" "test_policy" {
  resource_type         = "net"
  resource_id 			= "%s"
  name					= "policy_match"
  smode					= "matching"
  action				= "drop"
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
  resource_id 			= "%s"
  name					= "policy_match"
  smode					= "matching"
  action				= "alg"
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
  resource_id 			= "%s"
  name					= "policy_match"
  smode					= "matching"
  action				= "alg"
  switch				= true

  rule_list {
	skey 				= "ua"
	operator			= "not_include"
	value				= "123"
	}
}
`
