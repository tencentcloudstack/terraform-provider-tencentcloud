package ses_test

import (
	"context"
	"fmt"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svcses "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/ses"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// go test -test.run TestAccTencentCloudSesReceiverResource_basic -v
func TestAccTencentCloudSesReceiverResource_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccStepSetRegion(t, "ap-hongkong")
			tcacctest.AccPreCheckBusiness(t, tcacctest.ACCOUNT_TYPE_SES)
		},
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckSesReceiverDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSesReceiver,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSesReceiverExists("tencentcloud_ses_receiver.receiver"),
					resource.TestCheckResourceAttrSet("tencentcloud_ses_receiver.receiver", "id"),
					resource.TestCheckResourceAttr("tencentcloud_ses_receiver.receiver", "receivers_name", "terraform_test"),
					resource.TestCheckResourceAttr("tencentcloud_ses_receiver.receiver", "desc", "description"),
					resource.TestCheckResourceAttr("tencentcloud_ses_receiver.receiver", "data.#", "2"),
					resource.TestCheckResourceAttr("tencentcloud_ses_receiver.receiver", "data.0.email", "abc@abc.com"),
					resource.TestCheckResourceAttr("tencentcloud_ses_receiver.receiver", "data.0.template_data", "{\"name\":\"xxx\",\"age\":\"xx\"}"),
				),
			},
			{
				ResourceName:      "tencentcloud_ses_receiver.receiver",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckSesReceiverDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	service := svcses.NewSesService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_ses_receiver" {
			continue
		}

		res, err := service.DescribeSesReceiverById(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}

		if res != nil {
			return fmt.Errorf("ses receiver %s still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckSesReceiverExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}

		service := svcses.NewSesService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
		res, err := service.DescribeSesReceiverById(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}

		if res == nil {
			return fmt.Errorf("ses receiver %s is not found", rs.Primary.ID)
		}

		return nil
	}
}

const testAccSesReceiver = `

resource "tencentcloud_ses_receiver" "receiver" {
  receivers_name = "terraform_test"
  desc = "description"

  data {
    email = "abc@abc.com"
    template_data = "{\"name\":\"xxx\",\"age\":\"xx\"}"
  }

  data {
    email = "abcd@abcd.com"
    template_data = "{\"name\":\"xxx\",\"age\":\"xx\"}"
  }
}

`
