package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudMongodbInstancePropertyResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMongodbInstanceProperty,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_mongodb_instance_property.instance_property", "id")),
			},
			{
				ResourceName:      "tencentcloud_mongodb_instance_property.instance_property",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccMongodbInstanceProperty = `

resource "tencentcloud_mongodb_instance_property" "instance_property" {
  instance_id = "cmgo-9d0p6umb"
  user_name = "test_account"
  password = "Abc@123..."
  mongo_user_password = "Abc@123."
  user_desc = "test account"
  auth_role {
		mask = 
		name_space = ""

  }
}

`
