package tcr_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svctcr "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tcr"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func init() {
	resource.AddTestSweepers("tencentcloud_tcr_customized_domain", &resource.Sweeper{
		Name: "tencentcloud_tcr_customized_domain",
		F:    testSweepTcrCustomizedDomain,
	})
}

// go test -v ./tencentcloud -sweep=ap-shanghai -sweep-run=tencentcloud_tcr_customized_domain
func testSweepTcrCustomizedDomain(r string) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	cli, _ := tcacctest.SharedClientForRegion(r)
	tcrService := svctcr.NewTCRService(cli.(tccommon.ProviderMeta).GetAPIV3Conn())

	domains, err := tcrService.DescribeTcrCustomizedDomainById(ctx, tcacctest.DefaultTCRInstanceId, nil)
	if err != nil {
		return err
	}
	if domains == nil {
		return nil
	}

	for _, v := range domains {
		delName := *v.DomainName

		if strings.HasPrefix(delName, "test") {
			err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
				err := tcrService.DeleteTcrCustomizedDomainById(ctx, tcacctest.DefaultTCRInstanceId, delName)
				if err != nil {
					return tccommon.RetryError(err)
				}
				return nil
			})
			if err != nil {
				return fmt.Errorf("[ERROR] delete tcr customize domain instance %s failed! reason:[%s]", delName, err.Error())
			}
		}
	}
	return nil
}

func TestAccTencentCloudTcrCustomizedDomainResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_COMMON) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccTcrCustomizedDomain, tcacctest.DefaultTCRSSL),
				PreConfig: func() {
					tcacctest.AccStepSetRegion(t, "ap-shanghai")
					tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_COMMON)
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_tcr_customized_domain.my_domain", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_tcr_customized_domain.my_domain", "registry_id"),
					resource.TestCheckResourceAttr("tencentcloud_tcr_customized_domain.my_domain", "domain_name", "www.test.com"),
					resource.TestCheckResourceAttr("tencentcloud_tcr_customized_domain.my_domain", "certificate_id", tcacctest.DefaultTCRSSL),
				),
			},
			{
				ResourceName:      "tencentcloud_tcr_customized_domain.my_domain",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTcrCustomizedDomain = tcacctest.DefaultTCRInstanceData + `

resource "tencentcloud_tcr_customized_domain" "my_domain" {
  registry_id = local.tcr_id
  domain_name = "www.test.com"
  certificate_id = "%s"
  tags = {
    "createdBy" = "terraform"
  }
}

`
