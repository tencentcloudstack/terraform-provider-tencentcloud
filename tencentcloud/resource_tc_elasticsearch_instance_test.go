package tencentcloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccTencentCloudElasticsearchInstance_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckElasticsearchInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccElasticsearchInstance,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckElasticsearchInstanceExists("tencentcloud_elasticsearch_instance.foo"),
					resource.TestCheckResourceAttr("tencentcloud_elasticsearch_instance.foo", "instance_name", "tf-ci-test"),
					resource.TestCheckResourceAttr("tencentcloud_elasticsearch_instance.foo", "availability_zone", defaultAZone),
					resource.TestCheckResourceAttr("tencentcloud_elasticsearch_instance.foo", "version", "7.5.1"),
					resource.TestCheckResourceAttr("tencentcloud_elasticsearch_instance.foo", "vpc_id", defaultVpcId),
					resource.TestCheckResourceAttr("tencentcloud_elasticsearch_instance.foo", "subnet_id", defaultSubnetId),
					resource.TestCheckResourceAttr("tencentcloud_elasticsearch_instance.foo", "license_type", "oss"),
					resource.TestCheckResourceAttr("tencentcloud_elasticsearch_instance.foo", "node_info_list.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_elasticsearch_instance.foo", "node_info_list.0.node_num", "2"),
					resource.TestCheckResourceAttr("tencentcloud_elasticsearch_instance.foo", "node_info_list.0.node_type", "ES.S1.SMALL2"),
					resource.TestCheckResourceAttr("tencentcloud_elasticsearch_instance.foo", "node_info_list.0.type", "hotData"),
					resource.TestCheckResourceAttr("tencentcloud_elasticsearch_instance.foo", "node_info_list.0.encrypt", "false"),
					resource.TestCheckResourceAttr("tencentcloud_elasticsearch_instance.foo", "tags.test", "terraform"),
				),
			},
			{
				Config: testAccElasticsearchInstanceUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckElasticsearchInstanceExists("tencentcloud_elasticsearch_instance.foo"),
					resource.TestCheckResourceAttr("tencentcloud_elasticsearch_instance.foo", "instance_name", "tf-ci-test-update"),
					resource.TestCheckResourceAttr("tencentcloud_elasticsearch_instance.foo", "license_type", "basic"),
					resource.TestCheckResourceAttr("tencentcloud_elasticsearch_instance.foo", "basic_security_type", "2"),
					resource.TestCheckResourceAttr("tencentcloud_elasticsearch_instance.foo", "tags.test", "test"),
				),
			},
			{
				ResourceName:            "tencentcloud_elasticsearch_instance.foo",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
		},
	})
}

func testAccCheckElasticsearchInstanceDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	elasticsearchService := ElasticsearchService{
		client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_elasticsearch_instance" {
			continue
		}

		instance, err := elasticsearchService.DescribeInstanceById(ctx, rs.Primary.ID)
		if err != nil {
			err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
				instance, err = elasticsearchService.DescribeInstanceById(ctx, rs.Primary.ID)
				if err != nil {
					return retryError(err)
				}
				return nil
			})
		}
		if err != nil {
			return err
		}
		if instance != nil {
			return fmt.Errorf("elasticsearch instance still exists: %s", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckElasticsearchInstanceExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("elasticsearch instance %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("elasticsearch instance id is not set")
		}
		elasticsearchService := ElasticsearchService{
			client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
		}
		instance, err := elasticsearchService.DescribeInstanceById(ctx, rs.Primary.ID)
		if err != nil {
			err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
				instance, err = elasticsearchService.DescribeInstanceById(ctx, rs.Primary.ID)
				if err != nil {
					return retryError(err)
				}
				return nil
			})
		}
		if err != nil {
			return err
		}
		if instance == nil {
			return fmt.Errorf("elasticsearch instance is not found")
		}
		return nil
	}
}

const testAccElasticsearchInstance = defaultVpcVariable + `
resource "tencentcloud_elasticsearch_instance" "foo" {
	instance_name     = "tf-ci-test"
	availability_zone = var.availability_zone
	version           = "7.5.1"
	vpc_id            = var.vpc_id
	subnet_id         = var.subnet_id
	password          = "Test1234"
	license_type      = "oss"
  
	node_info_list {
	  node_num  = 2
	  node_type = "ES.S1.SMALL2"
	}
  
	tags = {
	  test = "terraform"
	}
  }
`

const testAccElasticsearchInstanceUpdate = defaultVpcVariable + `
resource "tencentcloud_elasticsearch_instance" "foo" {
	instance_name       = "tf-ci-test-update"
	availability_zone   = var.availability_zone
	version             = "7.5.1"
	vpc_id              = var.vpc_id
	subnet_id           = var.subnet_id
	password            = "Test12345"
	license_type        = "basic"
	basic_security_type = 2
  
	node_info_list {
	  node_num  = 2
	  node_type = "ES.S1.SMALL2"
	}
  
	tags = {
	  test = "test"
	}
  }
`
