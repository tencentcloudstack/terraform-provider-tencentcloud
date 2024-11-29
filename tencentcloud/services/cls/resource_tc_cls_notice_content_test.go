package cls_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudClsNoticeContentResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccClsNoticeContent,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cls_notice_content.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_cls_notice_content.example", "name"),
					resource.TestCheckResourceAttrSet("tencentcloud_cls_notice_content.example", "type"),
					resource.TestCheckResourceAttrSet("tencentcloud_cls_notice_content.example", "notice_contents"),
				),
			},
			{
				Config: testAccClsNoticeContentUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cls_notice_content.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_cls_notice_content.example", "name"),
					resource.TestCheckResourceAttrSet("tencentcloud_cls_notice_content.example", "type"),
					resource.TestCheckResourceAttrSet("tencentcloud_cls_notice_content.example", "notice_contents"),
				),
			},
			{
				ResourceName:      "tencentcloud_cls_notice_content.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccClsNoticeContent = `
resource "tencentcloud_cls_notice_content" "example" {
  name = "tf-example"
  type = 0
  notice_contents {
    type = "Email"

    trigger_content {
      title   = "title"
      content = "This is content."
      headers = [
        "Content-Type:application/json"
      ]
    }

    recovery_content {
      title   = "title"
      content = "This is content."
      headers = [
        "Content-Type:application/json"
      ]
    }
  }
}
`

const testAccClsNoticeContentUpdate = `
resource "tencentcloud_cls_notice_content" "example" {
  name = "tf-example-update"
  type = 1
  notice_contents {
    type = "Sms"

    trigger_content {
      title   = "title"
      content = "This is content."
      headers = [
        "Content-Type:application/json"
      ]
    }

    recovery_content {
      title   = "title"
      content = "This is content."
      headers = [
        "Content-Type:application/json"
      ]
    }
  }
}
`
