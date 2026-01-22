package mongodb_test

import (
	"context"
	"fmt"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svcmongodb "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/mongodb"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// go test -i; go test -test.run TestAccTencentCloudMongodbInstanceSslResource_basic -v
func TestAccTencentCloudMongodbInstanceSslResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckMongodbInstanceSslDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongodbInstanceSsl_enable,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongodbInstanceSslExists("tencentcloud_mongodb_instance_ssl.ssl"),
					resource.TestCheckResourceAttrSet("tencentcloud_mongodb_instance_ssl.ssl", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_mongodb_instance_ssl.ssl", "instance_id"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_instance_ssl.ssl", "enable", "true"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_instance_ssl.ssl", "status", "1"),
					resource.TestCheckResourceAttrSet("tencentcloud_mongodb_instance_ssl.ssl", "expired_time"),
					resource.TestCheckResourceAttrSet("tencentcloud_mongodb_instance_ssl.ssl", "cert_url"),
				),
			},
			{
				ResourceName:      "tencentcloud_mongodb_instance_ssl.ssl",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccMongodbInstanceSsl_disable,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongodbInstanceSslExists("tencentcloud_mongodb_instance_ssl.ssl"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_instance_ssl.ssl", "enable", "false"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_instance_ssl.ssl", "status", "0"),
				),
			},
			{
				Config: testAccMongodbInstanceSsl_enable,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongodbInstanceSslExists("tencentcloud_mongodb_instance_ssl.ssl"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_instance_ssl.ssl", "enable", "true"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_instance_ssl.ssl", "status", "1"),
				),
			},
		},
	})
}

func testAccCheckMongodbInstanceSslDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	service := svcmongodb.NewMongodbService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_mongodb_instance_ssl" {
			continue
		}

		sslStatus, err := service.DescribeMongodbInstanceSSLById(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}

		if sslStatus != nil && sslStatus.Response != nil && sslStatus.Response.Status != nil {
			if *sslStatus.Response.Status == 1 {
				return fmt.Errorf("mongodb instance ssl %s still enabled", rs.Primary.ID)
			}
		}
	}

	return nil
}

func testAccCheckMongodbInstanceSslExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource %s is not found", n)
		}

		service := svcmongodb.NewMongodbService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
		sslStatus, err := service.DescribeMongodbInstanceSSLById(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}

		if sslStatus == nil || sslStatus.Response == nil {
			return fmt.Errorf("mongodb instance ssl %s is not found", rs.Primary.ID)
		}

		return nil
	}
}

const testAccMongodbInstanceSsl_base = tcacctest.DefaultMongoDBSpec + `
resource "tencentcloud_vpc" "vpc" {
  name       = "mongodb-ssl-vpc"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_subnet" "subnet" {
  vpc_id            = tencentcloud_vpc.vpc.id
  name              = "mongodb-ssl-subnet"
  cidr_block        = "10.0.0.0/16"
  availability_zone = "ap-guangzhou-3"
}

resource "tencentcloud_mongodb_instance" "mongodb" {
  instance_name  = "tf-mongodb-ssl-test"
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
`

const testAccMongodbInstanceSsl_enable = testAccMongodbInstanceSsl_base + `
resource "tencentcloud_mongodb_instance_ssl" "ssl" {
  instance_id = tencentcloud_mongodb_instance.mongodb.id
  enable      = true
}
`

const testAccMongodbInstanceSsl_disable = testAccMongodbInstanceSsl_base + `
resource "tencentcloud_mongodb_instance_ssl" "ssl" {
  instance_id = tencentcloud_mongodb_instance.mongodb.id
  enable      = false
}
`
