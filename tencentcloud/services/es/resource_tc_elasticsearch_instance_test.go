package es_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	svces "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/es"
)

func init() {
	// go test -v ./tencentcloud -sweep=ap-guangzhou -sweep-run=tencentcloud_elasticsearch_instance
	resource.AddTestSweepers("tencentcloud_elasticsearch_instance", &resource.Sweeper{
		Name: "tencentcloud_elasticsearch_instance",
		F: func(r string) error {
			logId := tccommon.GetLogId(tccommon.ContextNil)
			ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
			cli, _ := tcacctest.SharedClientForRegion(r)
			client := cli.(tccommon.ProviderMeta).GetAPIV3Conn()

			service := svces.NewElasticsearchService(client)

			es, err := service.DescribeInstancesByFilter(ctx, "", "tf-ci-test", nil)
			if err != nil {
				return err
			}

			// add scanning resources
			var resources, nonKeepResources []*tccommon.ResourceInstance
			for _, v := range es {
				if !tccommon.CheckResourcePersist(*v.InstanceName, *v.CreateTime) {
					nonKeepResources = append(nonKeepResources, &tccommon.ResourceInstance{
						Id:   *v.InstanceId,
						Name: *v.InstanceName,
					})
				}
				resources = append(resources, &tccommon.ResourceInstance{
					Id:         *v.InstanceId,
					Name:       *v.InstanceName,
					CreateTime: *v.CreateTime,
				})
			}
			tccommon.ProcessScanCloudResources(client, resources, nonKeepResources, "CreateInstance")

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
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckElasticsearchInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccElasticsearchInstance,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckElasticsearchInstanceExists("tencentcloud_elasticsearch_instance.foo"),
					resource.TestCheckResourceAttr("tencentcloud_elasticsearch_instance.foo", "instance_name", "tf-ci-test"),
					resource.TestCheckResourceAttr("tencentcloud_elasticsearch_instance.foo", "availability_zone", tcacctest.DefaultAZone),
					resource.TestCheckResourceAttr("tencentcloud_elasticsearch_instance.foo", "version", "7.10.1"),
					resource.TestCheckResourceAttr("tencentcloud_elasticsearch_instance.foo", "vpc_id", tcacctest.DefaultEsVpcId),
					resource.TestCheckResourceAttr("tencentcloud_elasticsearch_instance.foo", "subnet_id", tcacctest.DefaultEsSubnetId),
					resource.TestCheckResourceAttr("tencentcloud_elasticsearch_instance.foo", "license_type", "basic"),
					resource.TestCheckResourceAttr("tencentcloud_elasticsearch_instance.foo", "basic_security_type", "2"),
					resource.TestCheckResourceAttr("tencentcloud_elasticsearch_instance.foo", "web_node_type_info.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_elasticsearch_instance.foo", "web_node_type_info.0.node_num", "1"),
					resource.TestCheckResourceAttr("tencentcloud_elasticsearch_instance.foo", "web_node_type_info.0.node_type", "ES.S1.MEDIUM4"),
					resource.TestCheckResourceAttr("tencentcloud_elasticsearch_instance.foo", "node_info_list.0.node_num", "2"),
					resource.TestCheckResourceAttr("tencentcloud_elasticsearch_instance.foo", "node_info_list.0.node_type", "ES.S1.MEDIUM4"),
					resource.TestCheckResourceAttr("tencentcloud_elasticsearch_instance.foo", "node_info_list.0.type", "hotData"),
					resource.TestCheckResourceAttr("tencentcloud_elasticsearch_instance.foo", "node_info_list.0.encrypt", "false"),
					resource.TestCheckResourceAttr("tencentcloud_elasticsearch_instance.foo", "tags.test", "terraform"),
					resource.TestCheckResourceAttr("tencentcloud_elasticsearch_instance.foo", "kibana_public_access", "OPEN"),
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
func TestAccTencentCloudElasticsearchInstanceResource_kibanaPublicAccess(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckElasticsearchInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccElasticsearchInstanceKibanaPublicAccessOpen,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckElasticsearchInstanceExists("tencentcloud_elasticsearch_instance.es_kibana"),
					resource.TestCheckResourceAttr("tencentcloud_elasticsearch_instance.es_kibana", "kibana_public_access", "OPEN"),
					resource.TestCheckResourceAttr("tencentcloud_elasticsearch_instance.es_kibana", "public_access", "OPEN"),
					resource.TestCheckResourceAttr("tencentcloud_elasticsearch_instance.es_kibana", "es_public_acl.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_elasticsearch_instance.es_kibana", "es_public_acl.0.white_ip_list.#", "1"),
					resource.TestCheckResourceAttrSet("tencentcloud_elasticsearch_instance.es_kibana", "es_public_url"),
				),
			},
			{
				Config: testAccElasticsearchInstanceKibanaPublicAccessClose,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckElasticsearchInstanceExists("tencentcloud_elasticsearch_instance.es_kibana"),
					resource.TestCheckResourceAttr("tencentcloud_elasticsearch_instance.es_kibana", "kibana_public_access", "CLOSE"),
					resource.TestCheckResourceAttr("tencentcloud_elasticsearch_instance.es_kibana", "public_access", "CLOSE"),
					resource.TestCheckResourceAttr("tencentcloud_elasticsearch_instance.es_kibana", "es_public_acl.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_elasticsearch_instance.es_kibana", "es_public_acl.0.white_ip_list.#", "1"),
				),
			},
			{
				Config: testAccElasticsearchInstanceKibanaPublicAccessOpen,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckElasticsearchInstanceExists("tencentcloud_elasticsearch_instance.es_kibana"),
					resource.TestCheckResourceAttr("tencentcloud_elasticsearch_instance.es_kibana", "kibana_public_access", "OPEN"),
					resource.TestCheckResourceAttr("tencentcloud_elasticsearch_instance.es_kibana", "public_access", "OPEN"),
					resource.TestCheckResourceAttr("tencentcloud_elasticsearch_instance.es_kibana", "es_public_acl.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_elasticsearch_instance.es_kibana", "es_public_acl.0.white_ip_list.#", "1"),
				),
			},
		},
	})
}

func testAccCheckElasticsearchInstanceDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	elasticsearchService := svces.NewElasticsearchService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_elasticsearch_instance" {
			continue
		}

		instance, err := elasticsearchService.DescribeInstanceById(ctx, rs.Primary.ID)
		if err != nil {
			err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
				instance, err = elasticsearchService.DescribeInstanceById(ctx, rs.Primary.ID)
				if err != nil {
					return tccommon.RetryError(err)
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
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("elasticsearch instance %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("elasticsearch instance id is not set")
		}
		elasticsearchService := svces.NewElasticsearchService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
		instance, err := elasticsearchService.DescribeInstanceById(ctx, rs.Primary.ID)
		if err != nil {
			err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
				instance, err = elasticsearchService.DescribeInstanceById(ctx, rs.Primary.ID)
				if err != nil {
					return tccommon.RetryError(err)
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

const testAccElasticsearchInstance = tcacctest.DefaultEsVariables + `
resource "tencentcloud_elasticsearch_instance" "foo" {
	instance_name       = "tf-ci-test"
	availability_zone   = var.availability_zone
	version             = "7.10.1"
	vpc_id              = var.vpc_id
	subnet_id           = var.subnet_id
	password            = "Test1234"
	license_type        = "basic"
	basic_security_type = 2

    web_node_type_info {
      node_num = 1
      node_type = "ES.S1.MEDIUM4"
    }

	node_info_list {
	  node_num          = 2
	  node_type         = "ES.S1.MEDIUM4"
	}

	es_acl {
	  white_list = [
		"127.0.0.2"
	  ]
	  black_list = [
		"1.1.1.1"
	  ]
	}
  
	tags = {
	  test = "terraform"
	}
  }
`

const testAccElasticsearchInstanceUpdate = tcacctest.DefaultEsVariables + `
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
		"127.0.0.1"
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

const testAccElasticsearchInstanceKibanaPublicAccessClose = tcacctest.DefaultEsVariables + `
resource "tencentcloud_elasticsearch_instance" "es_kibana" {
	instance_name        = "tf-ci-test-kibana"
	availability_zone    = var.availability_zone
	version              = "7.10.1"
	vpc_id               = var.vpc_id
	subnet_id            = var.subnet_id
	password             = "Test1234"
	license_type         = "basic"
	basic_security_type  = 2
	kibana_public_access = "CLOSE"
	public_access = "CLOSE"
	es_public_acl {
	  white_ip_list = [
		"127.0.0.1"
	  ]
	}
  
	node_info_list {
	  node_num  = 2
	  node_type = "ES.S1.MEDIUM4"
	}
  }
`

const testAccElasticsearchInstanceKibanaPublicAccessOpen = tcacctest.DefaultEsVariables + `
resource "tencentcloud_elasticsearch_instance" "es_kibana" {
	instance_name        = "tf-ci-test-kibana"
	availability_zone    = var.availability_zone
	version              = "7.10.1"
	vpc_id               = var.vpc_id
	subnet_id            = var.subnet_id
	password             = "Test1234"
	license_type         = "basic"
	basic_security_type  = 2
	kibana_public_access = "OPEN"
	public_access = "OPEN"
	es_public_acl {
	  white_ip_list = [
		"127.0.0.1"
	  ]
	}
  
	node_info_list {
	  node_num  = 2
	  node_type = "ES.S1.MEDIUM4"
	}
  }
`
