package dnspod_test

import (
	"context"
	"fmt"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	svcdnspod "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/dnspod"
)

const Domain = "terraformer.com"

func init() {
	// go test -v ./tencentcloud -sweep=ap-guangzhou -sweep-run=tencentcloud_dnspod_domain_instance
	resource.AddTestSweepers("tencentcloud_dnspod_domain_instance", &resource.Sweeper{
		Name: "tencentcloud_dnspod_domain_instance",
		F:    testSweepDnspodDoamin,
	})
}
func testSweepDnspodDoamin(region string) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	cli, _ := tcacctest.SharedClientForRegion(region)
	client := cli.(tccommon.ProviderMeta).GetAPIV3Conn()
	service := svcdnspod.NewDnspodService(client)

	response, err := service.DescribeDomain(ctx, Domain)
	if err != nil {
		return err
	}

	domainInfo := response.Response.DomainInfo
	if domainInfo == nil {
		return nil
	}

	err = service.DeleteDomain(ctx, Domain)
	if err != nil {
		return err
	}

	return nil
}

// go test -i; go test -test.run TestAccTencentCloudDnspodDoamin -v
func TestAccTencentCloudDnspodDoamin(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_COMMON) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckDnspodDomainDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudDnspodDomain,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDnspodDomainExists("tencentcloud_dnspod_domain_instance.domain"),
					resource.TestCheckResourceAttr("tencentcloud_dnspod_domain_instance.domain", "domain", "terraformer.com"),
					resource.TestCheckResourceAttr("tencentcloud_dnspod_domain_instance.domain", "is_mark", "no"),
				),
			},
		},
	})
}

func testAccCheckDnspodDomainDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	dnspodService := svcdnspod.NewDnspodService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_dnspod_domain_instance" {
			continue
		}

		response, err := dnspodService.DescribeDomain(ctx, rs.Primary.ID)
		if err != nil {
			return nil
		}
		if response.Response.DomainInfo != nil {
			return fmt.Errorf("record rule %s still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckDnspodDomainExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("domain id is not set")
		}

		dnspodService := svcdnspod.NewDnspodService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
		response, err := dnspodService.DescribeDomain(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}
		if response.Response.DomainInfo == nil {
			return fmt.Errorf("dnspod domain %s is not found", rs.Primary.ID)
		}

		return nil
	}
}

const testAccTencentCloudDnspodDomain = `
resource "tencentcloud_dnspod_domain_instance" "domain" {
  domain      = "` + Domain + `"
  is_mark     = "no"
}
`
