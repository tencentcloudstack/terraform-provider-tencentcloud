package gaap_test

import (
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svcgaap "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/gaap"

	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func init() {
	// go test -v ./tencentcloud -sweep=ap-guangzhou -sweep-run=tencentcloud_gaap_realserver
	resource.AddTestSweepers("tencentcloud_gaap_realserver", &resource.Sweeper{
		Name: "tencentcloud_gaap_realserver",
		F: func(r string) error {
			logId := tccommon.GetLogId(tccommon.ContextNil)
			ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
			sharedClient, err := tcacctest.SharedClientForRegion(r)
			if err != nil {
				return fmt.Errorf("getting tencentcloud client error: %s", err.Error())
			}
			client := sharedClient.(tccommon.ProviderMeta)
			service := svcgaap.NewGaapService(client.GetAPIV3Conn())

			realservers, err := service.DescribeRealservers(ctx, nil, nil, nil, -1)
			if err != nil {
				return err
			}
			for _, realserver := range realservers {
				instanceName := *realserver.RealServerName

				if strings.HasPrefix(instanceName, tcacctest.KeepResource) || strings.HasPrefix(instanceName, tcacctest.DefaultResource) {
					continue
				}

				ee := service.DeleteRealserver(ctx, *realserver.RealServerId)
				if ee != nil {
					continue
				}
			}

			return nil
		},
	})
}

func TestAccTencentCloudGaapRealserver_basic(t *testing.T) {
	t.Parallel()
	id := new(string)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PREPAY) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckGaapRealserverDestroy(id),
		Steps: []resource.TestStep{
			{
				Config: testAccGaapRealserverBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGaapRealserverExists("tencentcloud_gaap_realserver.foo", id),
					resource.TestCheckResourceAttr("tencentcloud_gaap_realserver.foo", "ip", "1.2.2.2"),
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
	t.Parallel()
	id := new(string)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PREPAY) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckGaapRealserverDestroy(id),
		Steps: []resource.TestStep{
			{
				Config: testAccGaapRealserverDomain,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGaapRealserverExists("tencentcloud_gaap_realserver.foo", id),
					resource.TestCheckResourceAttr("tencentcloud_gaap_realserver.foo", "domain", "www1.qq.com"),
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
	t.Parallel()
	id := new(string)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PREPAY) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckGaapRealserverDestroy(id),
		Steps: []resource.TestStep{
			{
				Config: testAccGaapRealserverBasic2,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGaapRealserverExists("tencentcloud_gaap_realserver.foo", id),
					resource.TestCheckResourceAttr("tencentcloud_gaap_realserver.foo", "ip", "1.2.2.3"),
					resource.TestCheckNoResourceAttr("tencentcloud_gaap_realserver.foo", "domain"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_realserver.foo", "name", "ci-test-gaap-realserver"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_realserver.foo", "project_id", "0"),
					resource.TestCheckNoResourceAttr("tencentcloud_gaap_realserver.foo", "tags"),
				),
			},
			{
				Config: testAccGaapRealserverNewName,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("tencentcloud_gaap_realserver.foo"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_realserver.foo", "name", "ci-test-gaap-realserver-new"),
				),
			},
		},
	})
}

func TestAccTencentCloudGaapRealserver_updateTags(t *testing.T) {
	t.Parallel()
	id := new(string)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PREPAY) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckGaapRealserverDestroy(id),
		Steps: []resource.TestStep{
			{
				Config: testAccGaapRealserverBasic3,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGaapRealserverExists("tencentcloud_gaap_realserver.foo", id),
					resource.TestCheckResourceAttr("tencentcloud_gaap_realserver.foo", "ip", "1.2.2.4"),
					resource.TestCheckNoResourceAttr("tencentcloud_gaap_realserver.foo", "domain"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_realserver.foo", "name", "ci-test-gaap-realserver"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_realserver.foo", "project_id", "0"),
					resource.TestCheckNoResourceAttr("tencentcloud_gaap_realserver.foo", "tags"),
				),
			},
			{
				Config: testAccGaapRealserverUpdateTags,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGaapRealserverExists("tencentcloud_gaap_realserver.foo", id),
					resource.TestCheckResourceAttr("tencentcloud_gaap_realserver.foo", "tags.test", "test"),
				),
			},
		},
	})
}

func testAccCheckGaapRealserverDestroy(id *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn()
		service := svcgaap.NewGaapService(client)

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

		projectIdStr := rs.Primary.Attributes["project_id"]
		projectId, err := strconv.Atoi(projectIdStr)
		if err != nil {
			return err
		}

		service := svcgaap.NewGaapService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())

		realservers, err := service.DescribeRealservers(context.TODO(), nil, nil, nil, projectId)
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
  ip   = "1.2.2.2"
  name = "ci-test-gaap-realserver"
}
`

const testAccGaapRealserverBasic2 = `
resource tencentcloud_gaap_realserver "foo" {
  ip   = "1.2.2.3"
  name = "ci-test-gaap-realserver"
}
`

const testAccGaapRealserverBasic3 = `
resource tencentcloud_gaap_realserver "foo" {
  ip   = "1.2.2.4"
  name = "ci-test-gaap-realserver"
}
`

const testAccGaapRealserverDomain = `
resource tencentcloud_gaap_realserver "foo" {
  domain = "www1.qq.com"
  name   = "ci-test-gaap-realserver"
}
`

const testAccGaapRealserverNewName = `
resource tencentcloud_gaap_realserver "foo" {
  ip   = "1.2.2.3"
  name = "ci-test-gaap-realserver-new"
}
`

const testAccGaapRealserverUpdateTags = `
resource tencentcloud_gaap_realserver "foo" {
  ip   = "1.2.2.4"
  name = "ci-test-gaap-realserver"

  tags = {
    "test" = "test"
  }
}
`
