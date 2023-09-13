package tencentcloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// go test -test.run TestAccTencentCloudSesReceiverResource_basic -v
func TestAccTencentCloudSesReceiverResource_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccStepSetRegion(t, "ap-hongkong")
			testAccPreCheckBusiness(t, ACCOUNT_TYPE_SES)
		},
		Providers:    testAccProviders,
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
			// {
			// 	ResourceName:      "tencentcloud_ses_receiver.receiver",
			// 	ImportState:       true,
			// 	ImportStateVerify: true,
			// },
		},
	})
}

func testAccCheckSesReceiverDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	service := SesService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
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
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}

		service := SesService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
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
