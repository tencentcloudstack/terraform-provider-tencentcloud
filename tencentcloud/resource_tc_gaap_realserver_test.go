package tencentcloud

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccTencentCloudGaapRealserver_basic(t *testing.T) {
	id := new(string)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGaapRealserverDestroy(id),
		Steps: []resource.TestStep{
			{
				Config: testAccGaapRealserverBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGaapRealserverExists("tencentcloud_gaap_realserver.foo", id),
					resource.TestCheckResourceAttr("tencentcloud_gaap_realserver.foo", "ip", "1.1.1.1"),
					resource.TestCheckNoResourceAttr("tencentcloud_gaap_realserver.foo", "domain"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_realserver.foo", "name", "ci-test-gaap-realserver"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_realserver.foo", "project_id", "0"),
					resource.TestCheckNoResourceAttr("tencentcloud_gaap_realserver.foo", "tags"),
				),
			},
			{
				ResourceName:      "tencentcloud_gaap_realserver.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccTencentCloudGaapRealserver_domain(t *testing.T) {
	id := new(string)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGaapRealserverDestroy(id),
		Steps: []resource.TestStep{
			{
				Config: testAccGaapRealserverDomain,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGaapRealserverExists("tencentcloud_gaap_realserver.foo", id),
					resource.TestCheckResourceAttr("tencentcloud_gaap_realserver.foo", "domain", "www.qq.com"),
					resource.TestCheckNoResourceAttr("tencentcloud_gaap_realserver.foo", "ip"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_realserver.foo", "name", "ci-test-gaap-realserver"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_realserver.foo", "project_id", "0"),
					resource.TestCheckNoResourceAttr("tencentcloud_gaap_realserver.foo", "tags"),
				),
			},
		},
	})
}

func TestAccTencentCloudGaapRealserver_updateName(t *testing.T) {
	id := new(string)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGaapRealserverDestroy(id),
		Steps: []resource.TestStep{
			{
				Config: testAccGaapRealserverBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGaapRealserverExists("tencentcloud_gaap_realserver.foo", id),
					resource.TestCheckResourceAttr("tencentcloud_gaap_realserver.foo", "ip", "1.1.1.1"),
					resource.TestCheckNoResourceAttr("tencentcloud_gaap_realserver.foo", "domain"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_realserver.foo", "name", "ci-test-gaap-realserver"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_realserver.foo", "project_id", "0"),
					resource.TestCheckNoResourceAttr("tencentcloud_gaap_realserver.foo", "tags"),
				),
			},
			{
				Config: testAccGaapRealserverNewName,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("tencentcloud_gaap_realserver.foo"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_realserver.foo", "name", "ci-test-gaap-realserver-new"),
				),
			},
		},
	})
}

func testAccCheckGaapRealserverDestroy(id *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*TencentCloudClient).apiV3Conn
		service := GaapService{client: client}

		realservers, err := service.DescribeRealservers(context.TODO(), id, nil, nil, -1)
		if err != nil {
			return err
		}

		if len(realservers) != 0 {
			return errors.New("realserver still exists")
		}

		return nil
	}
}

func testAccCheckGaapRealserverExists(n string, id *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("no realserver ID is set")
		}

		service := GaapService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}

		realservers, err := service.DescribeRealservers(context.TODO(), id, nil, nil, -1)
		if err != nil {
			return err
		}

		if len(realservers) == 0 {
			return fmt.Errorf("realserver not found: %s", rs.Primary.ID)
		}

		for _, realserver := range realservers {
			if realserver.RealServerId == nil {
				return errors.New("realserver id is nil")
			}
			if *realserver.RealServerId == rs.Primary.ID {
				*id = rs.Primary.ID
				break
			}
		}

		if *id == "" {
			return fmt.Errorf("realserver not found: %s", rs.Primary.ID)
		}

		return nil
	}
}

const testAccGaapRealserverBasic = `
resource tencentcloud_gaap_realserver "foo" {
  ip   = "1.1.1.1"
  name = "ci-test-gaap-realserver"
}
`

const testAccGaapRealserverDomain = `
resource tencentcloud_gaap_realserver "foo" {
  domain = "www.qq.com"
  name   = "ci-test-gaap-realserver"
}
`

const testAccGaapRealserverNewName = `
resource tencentcloud_gaap_realserver "foo" {
  ip   = "1.1.1.1"
  name = "ci-test-gaap-realserver-new"
}
`
