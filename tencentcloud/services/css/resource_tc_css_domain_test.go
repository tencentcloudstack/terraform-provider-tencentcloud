package css_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	svccss "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/css"
)

func init() {
	resource.AddTestSweepers("tencentcloud_css_domain", &resource.Sweeper{
		Name: "tencentcloud_css_domain",
		F:    testSweepCssDomainResource,
	})
}

// go test -v ./tencentcloud -sweep=ap-guangzhou -sweep-run=tencentcloud_css_domain
func testSweepCssDomainResource(r string) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	cli, _ := tcacctest.SharedClientForRegion(r)
	cssService := svccss.NewCssService(cli.(tccommon.ProviderMeta).GetAPIV3Conn())

	instances, err := cssService.DescribeCssDomainsByFilter(ctx, nil)
	if err != nil {
		return err
	}
	if instances == nil {
		return fmt.Errorf("css domain instance not exists.")
	}

	for _, v := range instances {
		delDomain := v.Name

		if strings.HasPrefix(*delDomain, "test_css_") {
			err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
				err := cssService.DeleteCssDomainById(ctx, delDomain, v.Type)
				if err != nil {
					return tccommon.RetryError(err)
				}
				return nil
			})
			if err != nil {
				return fmt.Errorf("[ERROR] delete css domain instance %s delDomain! reason:[%s]", *delDomain, err.Error())
			}
		}
	}
	return nil
}

func TestAccTencentCloudCssDomainResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckCssDomainDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCssDomain_push_enable,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCssDomainExists("tencentcloud_css_domain.domain"),
					resource.TestCheckResourceAttrSet("tencentcloud_css_domain.domain", "id"),
					resource.TestCheckResourceAttr("tencentcloud_css_domain.domain", "domain_name", "iac-tf.cloud"),
					resource.TestCheckResourceAttr("tencentcloud_css_domain.domain", "domain_type", "0"),
					resource.TestCheckResourceAttr("tencentcloud_css_domain.domain", "play_type", "1"),
					// resource.TestCheckResourceAttr("tencentcloud_css_domain.domain", "is_delay_live", "0"),
					resource.TestCheckResourceAttr("tencentcloud_css_domain.domain", "is_mini_program_live", "0"),
					resource.TestCheckResourceAttr("tencentcloud_css_domain.domain", "enable", "true"),
				),
			},
			{
				Config: testAccCssDomain_push_disable,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCssDomainExists("tencentcloud_css_domain.domain"),
					resource.TestCheckResourceAttrSet("tencentcloud_css_domain.domain", "id"),
					resource.TestCheckResourceAttr("tencentcloud_css_domain.domain", "domain_name", "iac-tf.cloud"),
					resource.TestCheckResourceAttr("tencentcloud_css_domain.domain", "domain_type", "0"),
					resource.TestCheckResourceAttr("tencentcloud_css_domain.domain", "play_type", "1"),
					// resource.TestCheckResourceAttr("tencentcloud_css_domain.domain", "is_delay_live", "0"),
					resource.TestCheckResourceAttr("tencentcloud_css_domain.domain", "is_mini_program_live", "0"),
					resource.TestCheckResourceAttr("tencentcloud_css_domain.domain", "enable", "false"),
				),
			},
			{
				ResourceName:            "tencentcloud_css_domain.domain",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"verify_owner_type"},
			},
		},
	})
}

func testAccCheckCssDomainDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	cssService := svccss.NewCssService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_css_domain" {
			continue
		}

		ret, err := cssService.DescribeCssDomainById(ctx, rs.Primary.ID)
		if err != nil {
			if strings.Contains(err.Error(), "ResourceNotFound.DomainNotExist") {
				return nil
			}
			return err
		}

		if ret != nil || *ret.Status == svccss.CSS_DOMAIN_STATUS_ACTIVATED {
			return fmt.Errorf("css domain instance still exist, name: %v", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckCssDomainExists(re string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[re]
		if !ok {
			return fmt.Errorf("css domain instance  %s is not found", re)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("css domain name is not set")
		}

		cssService := svccss.NewCssService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
		ret, err := cssService.DescribeCssDomainById(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}

		if ret == nil {
			return fmt.Errorf("css domain instance not found, name: %v", rs.Primary.ID)
		}

		return nil
	}
}

const testAccCssDomain_push_enable = `
resource "tencentcloud_css_domain" "domain" {
  domain_name = "iac-tf.cloud"
  domain_type = 0
  play_type = 1
  is_mini_program_live = 0
  verify_owner_type = "dbCheck"
  enable = true
}

`

const testAccCssDomain_push_disable = `
resource "tencentcloud_css_domain" "domain" {
  domain_name = "iac-tf.cloud"
  domain_type = 0
  play_type = 1
  is_mini_program_live = 0
  verify_owner_type = "dbCheck"
  enable = false
}

`
