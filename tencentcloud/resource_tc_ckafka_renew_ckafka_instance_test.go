package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCkafkaRenewCkafkaInstanceResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCkafkaRenewCkafkaInstance,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_ckafka_renew_ckafka_instance.renew_ckafka_instance", "id")),
			},
			{
				ResourceName:      "tencentcloud_ckafka_renew_ckafka_instance.renew_ckafka_instance",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCkafkaRenewCkafkaInstance = `

resource "tencentcloud_ckafka_renew_ckafka_instance" "renew_ckafka_instance" {
  instance_id = "InstanceId"
  time_span = 1
  tags = {
    "createdBy" = "terraform"
  }
}

`
