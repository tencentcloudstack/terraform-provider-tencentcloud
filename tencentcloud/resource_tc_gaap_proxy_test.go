package tencentcloud

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccTencentCloudGaapProxy_basic(t *testing.T) {
	id := new(string)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGaapProxyDestroy(id),
		Steps: []resource.TestStep{
			{
				Config: testAccGaapProxyBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGaapProxyExists("tencentcloud_gaap_proxy.foo", id),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "name", "ci-test-gaap-proxy"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "bandwidth", "10"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "concurrent", "2"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "project_id", "0"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "access_region", "SouthChina"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "realserver_region", "NorthChina"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "enable", "true"),
					resource.TestCheckNoResourceAttr("tencentcloud_gaap_proxy.foo", "tags"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_proxy.foo", "create_time"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_proxy.foo", "status"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_proxy.foo", "access_domain"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_proxy.foo", "access_ip"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_proxy.foo", "scalable"),
					resource.TestMatchResourceAttr("tencentcloud_gaap_proxy.foo", "support_protocols.#", regexp.MustCompile(`^[1-9]\d*$`)),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_proxy.foo", "forward_ip"),
				),
			},
		},
	})
}

func TestAccTencentCloudGaapProxy_updateName(t *testing.T) {
	id := new(string)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGaapProxyDestroy(id),
		Steps: []resource.TestStep{
			{
				Config: testAccGaapProxyBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGaapProxyExists("tencentcloud_gaap_proxy.foo", id),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "name", "ci-test-gaap-proxy"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "bandwidth", "10"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "concurrent", "2"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "project_id", "0"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "access_region", "SouthChina"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "realserver_region", "NorthChina"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "enable", "true"),
					resource.TestCheckNoResourceAttr("tencentcloud_gaap_proxy.foo", "tags"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_proxy.foo", "create_time"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_proxy.foo", "status"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_proxy.foo", "access_domain"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_proxy.foo", "access_ip"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_proxy.foo", "scalable"),
					resource.TestMatchResourceAttr("tencentcloud_gaap_proxy.foo", "support_protocols.#", regexp.MustCompile(`^[1-9]\d*$`)),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_proxy.foo", "forward_ip"),
				),
			},
			{
				Config: testAccGaapProxyNewName,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("tencentcloud_gaap_proxy.foo"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "name", "ci-test-gaap-proxy-new"),
				),
			},
		},
	})
}

func TestAccTencentCloudGaapProxy_updateBandwidth(t *testing.T) {
	id := new(string)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGaapProxyDestroy(id),
		Steps: []resource.TestStep{
			{
				Config: testAccGaapProxyBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGaapProxyExists("tencentcloud_gaap_proxy.foo", id),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "name", "ci-test-gaap-proxy"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "bandwidth", "10"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "concurrent", "2"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "project_id", "0"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "access_region", "SouthChina"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "realserver_region", "NorthChina"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "enable", "true"),
					resource.TestCheckNoResourceAttr("tencentcloud_gaap_proxy.foo", "tags"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_proxy.foo", "create_time"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_proxy.foo", "status"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_proxy.foo", "access_domain"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_proxy.foo", "access_ip"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_proxy.foo", "scalable"),
					resource.TestMatchResourceAttr("tencentcloud_gaap_proxy.foo", "support_protocols.#", regexp.MustCompile(`^[1-9]\d*$`)),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_proxy.foo", "forward_ip"),
				),
			},
			{
				Config: testAccGaapProxyNewBandwidth,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("tencentcloud_gaap_proxy.foo"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "bandwidth", "20"),
				),
			},
		},
	})
}

func TestAccTencentCloudGaapProxy_updateConcurrent(t *testing.T) {
	id := new(string)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGaapProxyDestroy(id),
		Steps: []resource.TestStep{
			{
				Config: testAccGaapProxyBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGaapProxyExists("tencentcloud_gaap_proxy.foo", id),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "name", "ci-test-gaap-proxy"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "bandwidth", "10"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "concurrent", "2"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "project_id", "0"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "access_region", "SouthChina"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "realserver_region", "NorthChina"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "enable", "true"),
					resource.TestCheckNoResourceAttr("tencentcloud_gaap_proxy.foo", "tags"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_proxy.foo", "create_time"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_proxy.foo", "status"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_proxy.foo", "access_domain"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_proxy.foo", "access_ip"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_proxy.foo", "scalable"),
					resource.TestMatchResourceAttr("tencentcloud_gaap_proxy.foo", "support_protocols.#", regexp.MustCompile(`^[1-9]\d*$`)),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_proxy.foo", "forward_ip"),
				),
			},
			{
				Config: testAccGaapProxyNewConcurrent,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("tencentcloud_gaap_proxy.foo"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "concurrent", "10"),
				),
			},
		},
	})
}

func TestAccTencentCloudGaapProxy_enable(t *testing.T) {
	id := new(string)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGaapProxyDestroy(id),
		Steps: []resource.TestStep{
			{
				Config: testAccGaapProxyDisable,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGaapProxyExists("tencentcloud_gaap_proxy.foo", id),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "name", "ci-test-gaap-proxy"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "bandwidth", "10"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "concurrent", "2"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "project_id", "0"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "access_region", "SouthChina"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "realserver_region", "NorthChina"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "enable", "false"),
					resource.TestCheckNoResourceAttr("tencentcloud_gaap_proxy.foo", "tags"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_proxy.foo", "create_time"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_proxy.foo", "status"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_proxy.foo", "access_domain"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_proxy.foo", "access_ip"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_proxy.foo", "scalable"),
					resource.TestMatchResourceAttr("tencentcloud_gaap_proxy.foo", "support_protocols.#", regexp.MustCompile(`^[1-9]\d*$`)),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_proxy.foo", "forward_ip"),
				),
			},
			{
				Config: testAccGaapProxyEnable,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGaapProxyExists("tencentcloud_gaap_proxy.foo", id),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "enable", "true"),
				),
			},
		},
	})
}

const testAccGaapProxyBasic = `
resource tencentcloud_gaap_proxy "foo" {
  name              = "ci-test-gaap-proxy"
  bandwidth         = 10
  concurrent        = 2
  access_region     = "SouthChina"
  realserver_region = "NorthChina"
}
`

func testAccCheckGaapProxyDestroy(id *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*TencentCloudClient).apiV3Conn
		service := GaapService{client: client}

		proxies, err := service.DescribeProxies(context.TODO(), []string{*id}, nil, nil, nil, nil)
		if err != nil {
			return err
		}

		if len(proxies) != 0 {
			return fmt.Errorf("proxy still exists")
		}

		return nil
	}
}

func testAccCheckGaapProxyExists(n string, id *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no proxy ID is set")
		}

		service := GaapService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}

		proxies, err := service.DescribeProxies(context.TODO(), []string{rs.Primary.ID}, nil, nil, nil, nil)
		if err != nil {
			return err
		}

		if len(proxies) == 0 {
			return fmt.Errorf("proxy not found: %s", rs.Primary.ID)
		}

		for _, proxy := range proxies {
			if proxy.ProxyId == nil {
				return errors.New("realserver id is nil")
			}
			if *proxy.ProxyId == rs.Primary.ID {
				*id = rs.Primary.ID
				break
			}
		}

		if *id == "" {
			return fmt.Errorf("proxy not found: %s", rs.Primary.ID)
		}

		return nil
	}
}

const testAccGaapProxyNewName = `
resource tencentcloud_gaap_proxy "foo" {
  name              = "ci-test-gaap-proxy-new"
  bandwidth         = 10
  concurrent        = 2
  access_region     = "SouthChina"
  realserver_region = "NorthChina"
}
`

const testAccGaapProxyNewBandwidth = `
resource tencentcloud_gaap_proxy "foo" {
  name              = "ci-test-gaap-proxy"
  bandwidth         = 20
  concurrent        = 2
  access_region     = "SouthChina"
  realserver_region = "NorthChina"
}
`

const testAccGaapProxyNewConcurrent = `
resource tencentcloud_gaap_proxy "foo" {
  name              = "ci-test-gaap-proxy"
  bandwidth         = 10
  concurrent        = 10
  access_region     = "SouthChina"
  realserver_region = "NorthChina"
}
`

const testAccGaapProxyDisable = `
resource tencentcloud_gaap_proxy "foo" {
  name              = "ci-test-gaap-proxy"
  bandwidth         = 10
  concurrent        = 2
  access_region     = "SouthChina"
  realserver_region = "NorthChina"
  enable            = false
}
`

const testAccGaapProxyEnable = `
resource tencentcloud_gaap_proxy "foo" {
  name              = "ci-test-gaap-proxy"
  bandwidth         = 10
  concurrent        = 2
  access_region     = "SouthChina"
  realserver_region = "NorthChina"
  enable            = true
}
`
