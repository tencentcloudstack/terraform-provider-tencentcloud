package tencentcloud

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

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
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	cli, _ := sharedClientForRegion(r)
	tcrService := TCRService{client: cli.(*TencentCloudClient).apiV3Conn}

	domains, err := tcrService.DescribeTcrCustomizedDomainById(ctx, defaultTCRInstanceId, nil)
	if err != nil {
		return err
	}
	if domains == nil {
		return fmt.Errorf("tcr customize domain instance not exists.")
	}

	for _, v := range domains {
		delName := *v.DomainName

		if strings.HasPrefix(delName, "test") {
			err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
				err := tcrService.DeleteTcrCustomizedDomainById(ctx, defaultTCRInstanceId, delName)
				if err != nil {
					return retryError(err)
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
		PreCheck:  func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_COMMON) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				SkipFunc: func() (bool, error) {
					if os.Getenv(E2ETEST_ENV_REGION) != "" || os.Getenv(E2ETEST_ENV_AZ) != "" {
						fmt.Printf("[International station]skip TestAccTencentCloudTcrCustomizedDomainResource_basic, because the international station did not support this feature yet!\n")
						return true, nil
					}
					return false, nil
				},
				Config: fmt.Sprintf(testAccTcrCustomizedDomain, defaultTCRSSL),
				PreConfig: func() {
					// testAccStepSetRegion(t, "ap-shanghai")
					testAccPreCheckCommon(t, ACCOUNT_TYPE_COMMON)
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_tcr_customized_domain.my_domain", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_tcr_customized_domain.my_domain", "registry_id"),
					resource.TestCheckResourceAttr("tencentcloud_tcr_customized_domain.my_domain", "domain_name", "www.test.com"),
					resource.TestCheckResourceAttrSet("tencentcloud_tcr_customized_domain.my_domain", "certificate_id"),
				),
			},
			{
				SkipFunc: func() (bool, error) {
					if os.Getenv(E2ETEST_ENV_REGION) != "" || os.Getenv(E2ETEST_ENV_AZ) != "" {
						fmt.Printf("[International station]skip TestAccTencentCloudTcrCustomizedDomainResource_basic, because the international station did not support this feature yet!\n")
						return true, nil
					}
					return false, nil
				},
				ResourceName:      "tencentcloud_tcr_customized_domain.my_domain",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTcrCustomizedDomain = TCRDataSource + `
variable "env_default_ssl_id" {
  default = ""
  type = string
}


resource "tencentcloud_tcr_customized_domain" "my_domain" {
  registry_id = local.tcr_id
  domain_name = "www.test.com"
  certificate_id = var.env_default_ssl_id != "" ? var.env_default_ssl_id : "%s"
  tags = {
    "createdBy" = "terraform"
  }
}

`
