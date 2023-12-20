package css_test

import (
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svccss "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/css"

	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccTencentCloudCssPlayDomainCertAttachmentResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		CheckDestroy: testAccCheckCssPlayDomainCertAttachmentDestroy,
		Providers:    tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCssPlayDomainCertAttachment, tcacctest.DefaultCSSBindingCertName, tcacctest.DefaultCSSPlayDomainName),
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_ssl_certificates.foo"),
					resource.TestCheckResourceAttr("data.tencentcloud_ssl_certificates.foo", "certificates.#", "1"),
					testAccCheckCssPlayDomainCertAttachmentExists("tencentcloud_css_play_domain_cert_attachment.play_domain_cert_attachment"),
					resource.TestCheckResourceAttrSet("tencentcloud_css_play_domain_cert_attachment.play_domain_cert_attachment", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_css_play_domain_cert_attachment.play_domain_cert_attachment", "cloud_cert_id"),
					resource.TestCheckResourceAttr("tencentcloud_css_play_domain_cert_attachment.play_domain_cert_attachment", "domain_info.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_css_play_domain_cert_attachment.play_domain_cert_attachment", "domain_info.0.domain_name", tcacctest.DefaultCSSPlayDomainName),
				),
			},
			{
				ResourceName:      "tencentcloud_css_play_domain_cert_attachment.play_domain_cert_attachment",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckCssPlayDomainCertAttachmentDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	cssService := svccss.NewCssService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_css_play_domain_cert_attachment" {
			continue
		}

		idSplit := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
		if len(idSplit) != 2 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		domainName := idSplit[0]
		cloudCertId := idSplit[1]

		ret, err := cssService.DescribeCssPlayDomainCertAttachmentById(ctx, domainName, cloudCertId)
		if err != nil {
			return err
		}

		if ret != nil {
			return fmt.Errorf("css cert attachment instance still exist, instanceId: %v", rs.Primary.ID)
		}

	}
	return nil
}

func testAccCheckCssPlayDomainCertAttachmentExists(re string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		cssService := svccss.NewCssService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())

		rs, ok := s.RootModule().Resources[re]
		if !ok {
			return fmt.Errorf("css cert attachment instance  %s is not found", re)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("css cert attachment instance id is not set")
		}

		idSplit := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
		if len(idSplit) != 2 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		domainName := idSplit[0]
		cloudCertId := idSplit[1]

		ret, err := cssService.DescribeCssPlayDomainCertAttachmentById(ctx, domainName, cloudCertId)
		if err != nil {
			return err
		}

		if ret != nil && *ret.DomainName == domainName {
			return nil
		}
		return fmt.Errorf("css cert attachment instance %v not found", rs.Primary.ID)
	}
}

const testAccCssPlayDomainCertAttachment = `

data "tencentcloud_ssl_certificates" "foo" {
	name = "%s"
}

resource "tencentcloud_css_play_domain_cert_attachment" "play_domain_cert_attachment" {
  cloud_cert_id = data.tencentcloud_ssl_certificates.foo.certificates.0.id
  domain_info {
    domain_name = "%s"
    status = 1
  }
}

`
