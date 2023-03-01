package tencentcloud

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func init() {
	// go test -v ./tencentcloud -sweep=ap-guangzhou -sweep-run=tencentcloud_elasticsearch_instance
	resource.AddTestSweepers("tencentcloud_elasticsearch_instance", &resource.Sweeper{
		Name: "tencentcloud_elasticsearch_instance",
		F: func(r string) error {
			logId := getLogId(contextNil)
			ctx := context.WithValue(context.TODO(), logIdKey, logId)
			cli, _ := sharedClientForRegion(r)
			client := cli.(*TencentCloudClient).apiV3Conn

			service := ElasticsearchService{client: client}

			es, err := service.DescribeInstancesByFilter(ctx, "", "tf-ci-test", nil)
			if err != nil {
				return err
			}

			for _, v := range es {
				id := *v.InstanceId
				name := *v.InstanceName
				if !strings.Contains(name, "tf-ci-test") {
					continue
				}
				if err := service.DeleteInstance(ctx, id); err != nil {
					continue
				}
			}

			return nil
		},
	})
}

func TestAccTencentCloudElasticsearchInstanceResource_basic(t *testing.T) {
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
					resource.TestCheckResourceAttr("tencentcloud_elasticsearch_instance.foo", "version", "7.10.1"),
					resource.TestCheckResourceAttr("tencentcloud_elasticsearch_instance.foo", "vpc_id", defaultVpcId),
					resource.TestCheckResourceAttr("tencentcloud_elasticsearch_instance.foo", "subnet_id", defaultSubnetId),
					resource.TestCheckResourceAttr("tencentcloud_elasticsearch_instance.foo", "license_type", "basic"),
					resource.TestCheckResourceAttr("tencentcloud_elasticsearch_instance.foo", "basic_security_type", "1"),
					resource.TestCheckResourceAttr("tencentcloud_elasticsearch_instance.foo", "web_node_type_info.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_elasticsearch_instance.foo", "web_node_type_info.0.node_num", "1"),
					resource.TestCheckResourceAttr("tencentcloud_elasticsearch_instance.foo", "web_node_type_info.0.node_type", "ES.S1.MEDIUM4"),
					resource.TestCheckResourceAttr("tencentcloud_elasticsearch_instance.foo", "node_info_list.0.node_num", "2"),
					resource.TestCheckResourceAttr("tencentcloud_elasticsearch_instance.foo", "node_info_list.0.node_type", "ES.S1.MEDIUM4"),
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
					resource.TestCheckResourceAttr("tencentcloud_elasticsearch_instance.foo", "license_type", "platinum"),
					resource.TestCheckResourceAttr("tencentcloud_elasticsearch_instance.foo", "tags.test", "test"),
					resource.TestCheckResourceAttr("tencentcloud_elasticsearch_instance.foo", "web_node_type_info.0.node_num", "1"),
					resource.TestCheckResourceAttr("tencentcloud_elasticsearch_instance.foo", "web_node_type_info.0.node_type", "ES.S1.MEDIUM8"),
					resource.TestCheckResourceAttr("tencentcloud_elasticsearch_instance.foo", "node_info_list.0.node_type", "ES.S1.MEDIUM8"),
					resource.TestCheckResourceAttr("tencentcloud_elasticsearch_instance.foo", "node_info_list.0.disk_size", "200"),
					resource.TestCheckResourceAttr("tencentcloud_elasticsearch_instance.foo", "es_acl.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_elasticsearch_instance.foo", "es_acl.0.white_list.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_elasticsearch_instance.foo", "es_acl.0.black_list.#", "1"),
				),
			},
			{
				ResourceName:            "tencentcloud_elasticsearch_instance.foo",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password", "basic_security_type"},
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
	instance_name       = "tf-ci-test"
	availability_zone   = var.availability_zone
	version             = "7.10.1"
	vpc_id              = var.vpc_id
	subnet_id           = var.subnet_id
	password            = "Test1234"
	license_type        = "basic"
	basic_security_type = 1

    web_node_type_info {
      node_num = 1
      node_type = "ES.S1.MEDIUM4"
    }

	node_info_list {
	  node_num          = 2
	  node_type         = "ES.S1.MEDIUM4"
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
	version             = "7.10.1"
	vpc_id              = var.vpc_id
	subnet_id           = var.subnet_id
	password            = "Test12345"
	license_type        = "platinum"

	web_node_type_info {
      node_num = 1
      node_type = "ES.S1.MEDIUM8"
  	}

	node_info_list {
	  node_num          = 2
	  node_type         = "ES.S1.MEDIUM8"
	  disk_size         = 200
	}

	es_acl {
	  white_list = [
		"0.0.0.0"
	  ]
	  black_list = [
		"1.1.1.1"
	  ]
	}
  
	tags = {
	  test = "test"
	}
  }
`
