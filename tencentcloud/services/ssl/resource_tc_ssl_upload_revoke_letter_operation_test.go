package ssl_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudSslUploadRevokeLetterResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_SSL)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSslUploadRevokeLetter,
				Check: resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_ssl_upload_revoke_letter_operation.upload_revoke_letter", "id"),
					resource.TestCheckResourceAttr("tencentcloud_ssl_upload_revoke_letter_operation.upload_revoke_letter", "certificate_id", "8xRYdDlc"),
					resource.TestCheckResourceAttrSet("tencentcloud_ssl_upload_revoke_letter_operation.upload_revoke_letter", "revoke_letter"),
				),
			},
		},
	})
}

const testAccSslUploadRevokeLetter = `

resource "tencentcloud_ssl_upload_revoke_letter_operation" "upload_revoke_letter" {
  certificate_id = "8xRYdDlc"
  revoke_letter = filebase64("./c.pdf")
}

`
