package mongodb_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudMongodbInstanceTransparentDataEncryptionResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMongodbInstanceTransparentDataEncryption,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_mongodb_instance_transparent_data_encryption.encryption", "id"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_instance_transparent_data_encryption.encryption", "transparent_data_encryption_status", "open"),
					resource.TestCheckResourceAttrSet("tencentcloud_mongodb_instance_transparent_data_encryption.encryption", "key_info_list.#"),
				),
			},
			{
				ResourceName: "tencentcloud_mongodb_instance_transparent_data_encryption.encryption",
				ImportState:  true,
			},
		},
	})
}

const testAccMongodbInstanceTransparentDataEncryption = tcacctest.DefaultMongoDBSpec + `
resource "tencentcloud_vpc" "vpc" {
	name       = "mongodb-instance-vpc"
	cidr_block = "10.0.0.0/16"
  }
  
resource "tencentcloud_subnet" "subnet" {
	vpc_id            = tencentcloud_vpc.vpc.id
	name              = "mongodb-instance-subnet"
	cidr_block        = "10.0.0.0/16"
	availability_zone = "ap-guangzhou-3"
}

resource "tencentcloud_mongodb_instance" "mongodb" {
  instance_name  = "tf-mongodb-encryption-test"
  memory         = local.memory
  volume         = local.volume
  engine_version = local.engine_version
  machine_type   = local.machine_type
  security_groups = [local.security_group_id]
  available_zone = "ap-guangzhou-3"
  project_id     = 0
  password       = "test1234"
  vpc_id         = tencentcloud_vpc.vpc.id
  subnet_id      = tencentcloud_subnet.subnet.id
}

resource "tencentcloud_mongodb_instance_transparent_data_encryption" "encryption" {
    instance_id = tencentcloud_mongodb_instance.mongodb.id
    kms_region = "ap-guangzhou"
}
`
