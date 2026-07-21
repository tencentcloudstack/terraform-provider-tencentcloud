package es_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
	es "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/es/v20180416"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
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

func TestAccTencentCloudElasticsearchInstanceResource_kibanaPrivateAccess(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckElasticsearchInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccElasticsearchInstanceKibanaPrivateAccessUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckElasticsearchInstanceExists("tencentcloud_elasticsearch_instance.es_kibana"),
					resource.TestCheckResourceAttr("tencentcloud_elasticsearch_instance.es_kibana", "kibana_public_access", "CLOSE"),
					resource.TestCheckResourceAttr("tencentcloud_elasticsearch_instance.es_kibana", "kibana_private_access", "OPEN"),
					resource.TestCheckResourceAttrSet("tencentcloud_elasticsearch_instance.es_kibana", "kibana_private_url"),
				),
			},
			{
				Config: testAccElasticsearchInstanceKibanaPrivateAccessDefault,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckElasticsearchInstanceExists("tencentcloud_elasticsearch_instance.es_kibana"),
					resource.TestCheckResourceAttr("tencentcloud_elasticsearch_instance.es_kibana", "kibana_public_access", "OPEN"),
					resource.TestCheckResourceAttr("tencentcloud_elasticsearch_instance.es_kibana", "kibana_private_access", "CLOSE"),
				),
			},
			{
				Config: testAccElasticsearchInstanceKibanaPrivateAccessUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckElasticsearchInstanceExists("tencentcloud_elasticsearch_instance.es_kibana"),
					resource.TestCheckResourceAttr("tencentcloud_elasticsearch_instance.es_kibana", "kibana_public_access", "CLOSE"),
					resource.TestCheckResourceAttr("tencentcloud_elasticsearch_instance.es_kibana", "kibana_private_access", "OPEN"),
					resource.TestCheckResourceAttrSet("tencentcloud_elasticsearch_instance.es_kibana", "kibana_private_url"),
				),
			},
		},
	})
}

func TestAccTencentCloudElasticsearchInstanceResource_publicAccess(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckElasticsearchInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccElasticsearchInstancePublicAccessDefault,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckElasticsearchInstanceExists("tencentcloud_elasticsearch_instance.es_kibana"),
					resource.TestCheckResourceAttr("tencentcloud_elasticsearch_instance.es_kibana", "public_access", "OPEN"),
					resource.TestCheckResourceAttr("tencentcloud_elasticsearch_instance.es_kibana", "es_public_acl.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_elasticsearch_instance.es_kibana", "es_public_acl.0.white_ip_list.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_elasticsearch_instance.es_kibana", "es_public_acl.0.white_ip_list.0", "127.0.0.1"),
				),
			},
			{
				Config: testAccElasticsearchInstanceKibanaPublicAccessUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckElasticsearchInstanceExists("tencentcloud_elasticsearch_instance.es_kibana"),
					resource.TestCheckResourceAttr("tencentcloud_elasticsearch_instance.es_kibana", "public_access", "OPEN"),
					resource.TestCheckResourceAttr("tencentcloud_elasticsearch_instance.es_kibana", "es_public_acl.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_elasticsearch_instance.es_kibana", "es_public_acl.0.white_ip_list.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_elasticsearch_instance.es_kibana", "es_public_acl.0.white_ip_list.0", "127.0.0.2"),
				),
			},
		},
	})
}

func TestAccTencentCloudElasticsearchInstanceResource_https(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckElasticsearchInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccElasticsearchInstanceKibanaPublicAccessHttps,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckElasticsearchInstanceExists("tencentcloud_elasticsearch_instance.es_kibana"),
					resource.TestCheckResourceAttr("tencentcloud_elasticsearch_instance.es_kibana", "protocol", "https"),
				),
			},
		},
	})
}
func TestAccTencentCloudElasticsearchInstanceResource_httpTohttps(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckElasticsearchInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccElasticsearchInstanceKibanaPublicAccessHttp,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckElasticsearchInstanceExists("tencentcloud_elasticsearch_instance.es_kibana"),
				),
			},
			{
				Config: testAccElasticsearchInstanceKibanaPublicAccessHttps,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckElasticsearchInstanceExists("tencentcloud_elasticsearch_instance.es_kibana"),
					resource.TestCheckResourceAttr("tencentcloud_elasticsearch_instance.es_kibana", "protocol", "https"),
				),
			},
		},
	})
}

func TestAccTencentCloudElasticsearchInstanceResource_nodeInfoList(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckElasticsearchInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccElasticsearchInstanceNodeInfoList,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckElasticsearchInstanceExists("tencentcloud_elasticsearch_instance.es_node_info_list"),
					resource.TestCheckResourceAttr("tencentcloud_elasticsearch_instance.es_node_info_list", "node_info_list.#", "1"),
				),
			},
			{
				Config: testAccElasticsearchInstanceNodeInfoListUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckElasticsearchInstanceExists("tencentcloud_elasticsearch_instance.es_node_info_list"),
					resource.TestCheckResourceAttr("tencentcloud_elasticsearch_instance.es_node_info_list", "node_info_list.#", "2"),
				),
			},
		},
	})
}

func TestAccTencentCloudElasticsearchInstanceResource_nodeInfoListIO(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckElasticsearchInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccElasticsearchInstanceNodeInfoListIO,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckElasticsearchInstanceExists("tencentcloud_elasticsearch_instance.es_node_info_list"),
					resource.TestCheckResourceAttr("tencentcloud_elasticsearch_instance.es_node_info_list", "node_info_list.#", "2"),
				),
			},
		},
	})
}

func TestAccTencentCloudElasticsearchInstanceResource_MultiZoneInfo(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckElasticsearchInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccElasticsearchInstanceMultiZoneInfo,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckElasticsearchInstanceExists("tencentcloud_elasticsearch_instance.foo"),
					resource.TestCheckResourceAttr("tencentcloud_elasticsearch_instance.foo", "node_info_list.#", "2"),
					resource.TestCheckResourceAttr("tencentcloud_elasticsearch_instance.foo", "deploy_mode", "0"),
				),
			},
			{
				Config: testAccElasticsearchInstanceMultiZoneInfoUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckElasticsearchInstanceExists("tencentcloud_elasticsearch_instance.foo"),
					resource.TestCheckResourceAttr("tencentcloud_elasticsearch_instance.foo", "node_info_list.#", "3"),
					resource.TestCheckResourceAttr("tencentcloud_elasticsearch_instance.foo", "deploy_mode", "1"),
					resource.TestCheckResourceAttr("tencentcloud_elasticsearch_instance.foo", "multi_zone_infos.#", "2"),
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

const testAccElasticsearchInstanceKibanaPrivateAccessDefault = tcacctest.DefaultEsVariables + `
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
	kibana_private_access = "CLOSE"
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

const testAccElasticsearchInstanceKibanaPrivateAccessUpdate = tcacctest.DefaultEsVariables + `
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
	kibana_private_access = "OPEN"
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

const testAccElasticsearchInstancePublicAccessDefault = tcacctest.DefaultEsVariables + `
resource "tencentcloud_elasticsearch_instance" "es_kibana" {
	instance_name        = "tf-ci-test-kibana"
	availability_zone    = var.availability_zone
	version              = "7.10.1"
	vpc_id               = var.vpc_id
	subnet_id            = var.subnet_id
	password             = "Test1234"
	license_type         = "basic"
	basic_security_type  = 2
	public_access = "OPEN"
	es_acl {
	  white_list = [
		"127.0.0.1"
	  ]
	}
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

const testAccElasticsearchInstanceKibanaPublicAccessUpdate = tcacctest.DefaultEsVariables + `
resource "tencentcloud_elasticsearch_instance" "es_kibana" {
	instance_name        = "tf-ci-test-kibana"
	availability_zone    = var.availability_zone
	version              = "7.10.1"
	vpc_id               = var.vpc_id
	subnet_id            = var.subnet_id
	password             = "Test1234"
	license_type         = "basic"
	basic_security_type  = 2
	public_access = "OPEN"
	es_acl {
	  white_list = [
		"127.0.0.2"
	  ]
	}
	es_public_acl {
	  white_ip_list = [
		"127.0.0.2"
	  ]
	}
  
	node_info_list {
	  node_num  = 2
	  node_type = "ES.S1.MEDIUM4"
	}
  }
`
const testAccElasticsearchInstanceKibanaPublicAccessHttp = tcacctest.DefaultEsVariables + `
resource "tencentcloud_elasticsearch_instance" "es_kibana" {
	instance_name        = "tf-ci-test-kibana"
	availability_zone    = var.availability_zone
	version              = "7.10.1"
	vpc_id               = var.vpc_id
	subnet_id            = var.subnet_id
	password             = "Test1234"
	license_type         = "basic"
	basic_security_type  = 2
	public_access = "OPEN"
	es_acl {
	  white_list = [
		"127.0.0.2"
	  ]
	}
	es_public_acl {
	  white_ip_list = [
		"127.0.0.2"
	  ]
	}
  
	node_info_list {
	  node_num  = 2
	  node_type = "ES.S1.MEDIUM4"
	}
  }
`

const testAccElasticsearchInstanceKibanaPublicAccessHttps = tcacctest.DefaultEsVariables + `
resource "tencentcloud_elasticsearch_instance" "es_kibana" {
	instance_name        = "tf-ci-test-kibana"
	availability_zone    = var.availability_zone
	version              = "7.10.1"
	vpc_id               = var.vpc_id
	subnet_id            = var.subnet_id
	password             = "Test1234"
	license_type         = "basic"
	basic_security_type  = 2
	public_access = "OPEN"
	protocol = "https"
	es_acl {
	  white_list = [
		"127.0.0.2"
	  ]
	}
	es_public_acl {
	  white_ip_list = [
		"127.0.0.2"
	  ]
	}
  
	node_info_list {
	  node_num  = 2
	  node_type = "ES.S1.MEDIUM4"
	}
  }
`

const testAccElasticsearchInstanceNodeInfoListIO = tcacctest.DefaultEsVariables + `
resource "tencentcloud_elasticsearch_instance" "es_node_info_list" {
	instance_name        = "tf-ci-test-node"
	availability_zone    = var.availability_zone
	version              = "7.10.1"
	vpc_id               = var.vpc_id
	subnet_id            = var.subnet_id
	password             = "Test1234"
	license_type         = "basic"
	basic_security_type  = 2
	public_access = "OPEN"
	protocol = "https"
	es_acl {
	  white_list = [
		"127.0.0.2"
	  ]
	}
	es_public_acl {
	  white_ip_list = [
		"127.0.0.2"
	  ]
	}
  
	node_info_list {
	  node_num  = 2
	  node_type = "ES.I1.4XLARGE64"
	  type      = "hotData"
	}

	node_info_list {
	  node_num  = 3
	  node_type = "ES.S1.MEDIUM4"
	  disk_size = 50
	  type      = "dedicatedMaster"
	  disk_type = "CLOUD_SSD"
	}
  }
`

const testAccElasticsearchInstanceNodeInfoList = tcacctest.DefaultEsVariables + `
resource "tencentcloud_elasticsearch_instance" "es_node_info_list" {
	instance_name        = "tf-ci-test-node"
	availability_zone    = var.availability_zone
	version              = "7.10.1"
	vpc_id               = var.vpc_id
	subnet_id            = var.subnet_id
	password             = "Test1234"
	license_type         = "basic"
	basic_security_type  = 2
	public_access = "OPEN"
	protocol = "https"
	es_acl {
	  white_list = [
		"127.0.0.2"
	  ]
	}
	es_public_acl {
	  white_ip_list = [
		"127.0.0.2"
	  ]
	}
  
	node_info_list {
	  node_num  = 2
	  node_type = "ES.I1.4XLARGE64"
	  type      = "hotData"
	}

	node_info_list {
	  node_num  = 3
	  node_type = "ES.S1.MEDIUM4"
	  disk_size = 50
	  type      = "dedicatedMaster"
	  disk_type = "CLOUD_SSD"
	}
  }
`

const testAccElasticsearchInstanceNodeInfoListUpdate = tcacctest.DefaultEsVariables + `
resource "tencentcloud_elasticsearch_instance" "es_node_info_list" {
	instance_name        = "tf-ci-test-node"
	availability_zone    = var.availability_zone
	version              = "7.10.1"
	vpc_id               = var.vpc_id
	subnet_id            = var.subnet_id
	password             = "Test1234"
	license_type         = "basic"
	basic_security_type  = 2
	public_access = "OPEN"
	protocol = "https"
	es_acl {
	  white_list = [
		"127.0.0.2"
	  ]
	}
	es_public_acl {
	  white_ip_list = [
		"127.0.0.2"
	  ]
	}
  
	node_info_list {
	  node_num  = 3
	  node_type = "ES.S1.MEDIUM8"
	  disk_size = 100
	  type      = "hotData"
	  disk_type = "CLOUD_SSD"
	}
	node_info_list {
	  node_num  = 3
	  node_type = "ES.S1.MEDIUM8"
	  disk_type = "CLOUD_SSD"
	  type      = "dedicatedMaster"
	  disk_size = 50
  	}
  }
`
const testAccElasticsearchInstanceMultiZoneInfo = `
resource "tencentcloud_elasticsearch_instance" "foo" {
  instance_name       = "tf-ci-test"
  availability_zone   = "ap-guangzhou-3"
  version             = "7.10.1"
  vpc_id              = "vpc-axrsmmrv"
  subnet_id           = "subnet-j5vja918"
  password            = "Test1234"
  license_type        = "basic"
  basic_security_type = 2

  node_info_list {
    node_num  = 3
    node_type = "ES.S1.MEDIUM4"
    disk_size = 50
    type      = "hotData"
    disk_type = "CLOUD_SSD"
  }
  node_info_list {
    node_num  = 3
    node_type = "ES.S1.MEDIUM8"
    # disk_type = "CLOUD_SSD"
    type      = "dedicatedMaster"
    disk_size = 50
  }
  es_acl {
    white_list = [
      "127.0.0.2"
    ]
    black_list = [
      "1.1.1.1"
    ]
  }
}
`

const testAccElasticsearchInstanceMultiZoneInfoUpdate = `
resource "tencentcloud_elasticsearch_instance" "foo" {
  instance_name       = "tf-ci-test"
  availability_zone   = "ap-guangzhou-3"
  version             = "7.10.1"
  vpc_id              = "vpc-axrsmmrv"
  subnet_id           = "subnet-j5vja918"
  password            = "Test1234"
  license_type        = "basic"
  basic_security_type = 2

  node_info_list {
    node_num  = 2
    node_type = "ES.S1.MEDIUM4"
    disk_size = 50
    type      = "warmData"
    disk_type = "CLOUD_PREMIUM"
  }
  node_info_list {
    node_num  = 6
    node_type = "ES.S1.MEDIUM4"
    disk_size = 50
    type      = "hotData"
    disk_type = "CLOUD_SSD"
  }
  node_info_list {
    node_num  = 3
    node_type = "ES.S1.MEDIUM8"
    # disk_type = "CLOUD_SSD"
    type      = "dedicatedMaster"
    disk_size = 50
  }
  es_acl {
    white_list = [
      "127.0.0.2"
    ]
    black_list = [
      "1.1.1.1"
    ]
  }
  deploy_mode = 1
  multi_zone_infos {
    availability_zone = "ap-guangzhou-3"
    subnet_id         = "subnet-j5vja918"
  }
  multi_zone_infos {
    availability_zone = "ap-guangzhou-4"
    subnet_id         = "subnet-oi7ya2j6"
  }
}
`

type mockMetaEs struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaEs) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaEs{}

func newMockMetaEs() *mockMetaEs {
	return &mockMetaEs{client: &connectivity.TencentCloudClient{}}
}

func ptrStringEs(s string) *string {
	return &s
}

func ptrInt64Es(v int64) *int64 {
	return &v
}

// go test ./tencentcloud/services/es/ -run "TestEsInstanceDestroyProtection" -v -count=1 -gcflags="all=-l"

// TestEsInstanceDestroyProtection_Schema validates the enable_destroy_protection schema field definition
func TestEsInstanceDestroyProtection_Schema(t *testing.T) {
	res := svces.ResourceTencentCloudElasticsearchInstance()
	assert.NotNil(t, res)
	assert.Contains(t, res.Schema, "enable_destroy_protection")

	field := res.Schema["enable_destroy_protection"]
	assert.Equal(t, schema.TypeString, field.Type)
	assert.True(t, field.Optional)
	assert.True(t, field.Computed)
	assert.NotNil(t, field.ValidateFunc)
}

// TestEsInstanceDestroyProtection_Read_NonNil verifies the read flow sets enable_destroy_protection
// from InstanceInfo.EnableDestroyProtection when the API returns a non-nil value.
func TestEsInstanceDestroyProtection_Read_NonNil(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	esClient := &es.Client{}
	patches.ApplyMethodReturn(newMockMetaEs().client, "UseEsClient", esClient)

	patches.ApplyMethodFunc(esClient, "DescribeInstances", func(request *es.DescribeInstancesRequest) (*es.DescribeInstancesResponse, error) {
		resp := es.NewDescribeInstancesResponse()
		resp.Response = &es.DescribeInstancesResponseParams{
			InstanceList: []*es.InstanceInfo{
				{
					InstanceId:              ptrStringEs("es-destroy-protection-test"),
					InstanceName:            ptrStringEs("tf-test-instance"),
					Zone:                    ptrStringEs("ap-guangzhou-3"),
					EsVersion:               ptrStringEs("7.10.1"),
					VpcUid:                  ptrStringEs("vpc-test"),
					SubnetUid:               ptrStringEs("subnet-test"),
					ChargeType:              ptrStringEs("POSTPAID_BY_HOUR"),
					Status:                  ptrInt64Es(1),
					LicenseType:             ptrStringEs("platinum"),
					EnableDestroyProtection: ptrStringEs("OPEN"),
				},
			},
		}
		return resp, nil
	})

	meta := newMockMetaEs()
	res := svces.ResourceTencentCloudElasticsearchInstance()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"version":  "7.10.1",
		"vpc_id":   "vpc-test",
		"password": "Test1234",
	})
	d.SetId("es-destroy-protection-test")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "OPEN", d.Get("enable_destroy_protection").(string))
}

// TestEsInstanceDestroyProtection_Read_Nil verifies the read flow does not overwrite
// enable_destroy_protection state when the API returns a nil value.
func TestEsInstanceDestroyProtection_Read_Nil(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	esClient := &es.Client{}
	patches.ApplyMethodReturn(newMockMetaEs().client, "UseEsClient", esClient)

	patches.ApplyMethodFunc(esClient, "DescribeInstances", func(request *es.DescribeInstancesRequest) (*es.DescribeInstancesResponse, error) {
		resp := es.NewDescribeInstancesResponse()
		resp.Response = &es.DescribeInstancesResponseParams{
			InstanceList: []*es.InstanceInfo{
				{
					InstanceId:   ptrStringEs("es-destroy-protection-test"),
					InstanceName: ptrStringEs("tf-test-instance"),
					Zone:         ptrStringEs("ap-guangzhou-3"),
					EsVersion:    ptrStringEs("7.10.1"),
					VpcUid:       ptrStringEs("vpc-test"),
					SubnetUid:    ptrStringEs("subnet-test"),
					ChargeType:   ptrStringEs("POSTPAID_BY_HOUR"),
					Status:       ptrInt64Es(1),
					LicenseType:  ptrStringEs("platinum"),
				},
			},
		}
		return resp, nil
	})

	meta := newMockMetaEs()
	res := svces.ResourceTencentCloudElasticsearchInstance()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"version":                   "7.10.1",
		"vpc_id":                    "vpc-test",
		"password":                  "Test1234",
		"enable_destroy_protection": "OPEN",
	})
	d.SetId("es-destroy-protection-test")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	// state should be preserved (nil-safe read does not overwrite)
	assert.Equal(t, "OPEN", d.Get("enable_destroy_protection").(string))
}

// TestEsInstanceDestroyProtection_Create verifies the create flow invokes UpdateInstance with
// EnableDestroyProtection set to the configured value after instance creation.
func TestEsInstanceDestroyProtection_Create(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	esClient := &es.Client{}
	patches.ApplyMethodReturn(newMockMetaEs().client, "UseEsClient", esClient)

	// mock CreateInstance
	patches.ApplyMethodFunc(esClient, "CreateInstance", func(request *es.CreateInstanceRequest) (*es.CreateInstanceResponse, error) {
		resp := es.NewCreateInstanceResponse()
		resp.Response = &es.CreateInstanceResponseParams{
			InstanceId: ptrStringEs("es-create-destroy-protection"),
		}
		return resp, nil
	})

	// track whether UpdateInstance was called with EnableDestroyProtection = OPEN
	var capturedEnableDestroyProtection *string
	patches.ApplyMethodFunc(esClient, "UpdateInstance", func(request *es.UpdateInstanceRequest) (*es.UpdateInstanceResponse, error) {
		capturedEnableDestroyProtection = request.EnableDestroyProtection
		resp := es.NewUpdateInstanceResponse()
		resp.Response = &es.UpdateInstanceResponseParams{}
		return resp, nil
	})

	// mock DescribeInstances for the post-create status waiting + final read
	callCount := 0
	patches.ApplyMethodFunc(esClient, "DescribeInstances", func(request *es.DescribeInstancesRequest) (*es.DescribeInstancesResponse, error) {
		callCount++
		resp := es.NewDescribeInstancesResponse()
		resp.Response = &es.DescribeInstancesResponseParams{
			InstanceList: []*es.InstanceInfo{
				{
					InstanceId:              ptrStringEs("es-create-destroy-protection"),
					InstanceName:            ptrStringEs("tf-test-instance"),
					Zone:                    ptrStringEs("ap-guangzhou-3"),
					EsVersion:               ptrStringEs("7.10.1"),
					VpcUid:                  ptrStringEs("vpc-test"),
					SubnetUid:               ptrStringEs("subnet-test"),
					ChargeType:              ptrStringEs("POSTPAID_BY_HOUR"),
					Status:                  ptrInt64Es(1),
					LicenseType:             ptrStringEs("platinum"),
					KibanaPublicAccess:      ptrStringEs("OPEN"),
					EnableDestroyProtection: ptrStringEs("OPEN"),
				},
			},
		}
		return resp, nil
	})

	meta := newMockMetaEs()
	res := svces.ResourceTencentCloudElasticsearchInstance()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_name":             "tf-test-instance",
		"availability_zone":         "ap-guangzhou-3",
		"version":                   "7.10.1",
		"vpc_id":                    "vpc-test",
		"subnet_id":                 "subnet-test",
		"password":                  "Test1234",
		"charge_type":               "POSTPAID_BY_HOUR",
		"license_type":              "platinum",
		"enable_destroy_protection": "OPEN",
		"node_info_list": []interface{}{
			map[string]interface{}{
				"node_num":  2,
				"node_type": "ES.S1.MEDIUM8",
				"type":      "hotData",
				"disk_type": "CLOUD_SSD",
				"disk_size": 100,
				"encrypt":   false,
			},
		},
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "es-create-destroy-protection", d.Id())
	// the UpdateInstance call for destroy protection must have carried EnableDestroyProtection = OPEN
	assert.NotNil(t, capturedEnableDestroyProtection)
	assert.Equal(t, "OPEN", *capturedEnableDestroyProtection)
	// final read should have populated state from DescribeInstances
	assert.Equal(t, "OPEN", d.Get("enable_destroy_protection").(string))
}

// TestEsInstanceDestroyProtection_Update verifies the update flow invokes UpdateInstance with
// EnableDestroyProtection when d.HasChange("enable_destroy_protection") is true.
func TestEsInstanceDestroyProtection_Update(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	esClient := &es.Client{}
	patches.ApplyMethodReturn(newMockMetaEs().client, "UseEsClient", esClient)

	var capturedEnableDestroyProtection *string
	patches.ApplyMethodFunc(esClient, "UpdateInstance", func(request *es.UpdateInstanceRequest) (*es.UpdateInstanceResponse, error) {
		if request.EnableDestroyProtection != nil {
			capturedEnableDestroyProtection = request.EnableDestroyProtection
		}
		resp := es.NewUpdateInstanceResponse()
		resp.Response = &es.UpdateInstanceResponseParams{}
		return resp, nil
	})

	// mock DescribeInstances for the upgrade-wait helper + final read
	patches.ApplyMethodFunc(esClient, "DescribeInstances", func(request *es.DescribeInstancesRequest) (*es.DescribeInstancesResponse, error) {
		resp := es.NewDescribeInstancesResponse()
		resp.Response = &es.DescribeInstancesResponseParams{
			InstanceList: []*es.InstanceInfo{
				{
					InstanceId:              ptrStringEs("es-update-destroy-protection"),
					InstanceName:            ptrStringEs("tf-test-instance"),
					Zone:                    ptrStringEs("ap-guangzhou-3"),
					EsVersion:               ptrStringEs("7.10.1"),
					VpcUid:                  ptrStringEs("vpc-test"),
					SubnetUid:               ptrStringEs("subnet-test"),
					ChargeType:              ptrStringEs("POSTPAID_BY_HOUR"),
					Status:                  ptrInt64Es(1),
					LicenseType:             ptrStringEs("platinum"),
					EnableDestroyProtection: ptrStringEs("OPEN"),
				},
			},
		}
		return resp, nil
	})

	meta := newMockMetaEs()
	res := svces.ResourceTencentCloudElasticsearchInstance()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_name":             "tf-test-instance",
		"availability_zone":         "ap-guangzhou-3",
		"version":                   "7.10.1",
		"vpc_id":                    "vpc-test",
		"subnet_id":                 "subnet-test",
		"password":                  "Test1234",
		"charge_type":               "POSTPAID_BY_HOUR",
		"license_type":              "platinum",
		"enable_destroy_protection": "OPEN",
	})
	d.SetId("es-update-destroy-protection")

	// force only enable_destroy_protection to be detected as changed
	patches.ApplyMethodFunc(d, "HasChange", func(key string) bool {
		return key == "enable_destroy_protection"
	})

	err := res.Update(d, meta)
	assert.NoError(t, err)
	assert.NotNil(t, capturedEnableDestroyProtection)
	assert.Equal(t, "OPEN", *capturedEnableDestroyProtection)
	assert.Equal(t, "OPEN", d.Get("enable_destroy_protection").(string))
}

// TestEsInstanceDestroyProtection_Update_NoChange verifies the update flow does NOT call
// UpdateInstance for destroy protection when d.HasChange("enable_destroy_protection") is false.
func TestEsInstanceDestroyProtection_Update_NoChange(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	esClient := &es.Client{}
	patches.ApplyMethodReturn(newMockMetaEs().client, "UseEsClient", esClient)

	updateCalled := false
	patches.ApplyMethodFunc(esClient, "UpdateInstance", func(request *es.UpdateInstanceRequest) (*es.UpdateInstanceResponse, error) {
		if request.EnableDestroyProtection != nil {
			updateCalled = true
		}
		resp := es.NewUpdateInstanceResponse()
		resp.Response = &es.UpdateInstanceResponseParams{}
		return resp, nil
	})

	// mock DescribeInstances for the final read
	patches.ApplyMethodFunc(esClient, "DescribeInstances", func(request *es.DescribeInstancesRequest) (*es.DescribeInstancesResponse, error) {
		resp := es.NewDescribeInstancesResponse()
		resp.Response = &es.DescribeInstancesResponseParams{
			InstanceList: []*es.InstanceInfo{
				{
					InstanceId:              ptrStringEs("es-update-destroy-protection"),
					InstanceName:            ptrStringEs("tf-test-instance"),
					Zone:                    ptrStringEs("ap-guangzhou-3"),
					EsVersion:               ptrStringEs("7.10.1"),
					VpcUid:                  ptrStringEs("vpc-test"),
					SubnetUid:               ptrStringEs("subnet-test"),
					ChargeType:              ptrStringEs("POSTPAID_BY_HOUR"),
					Status:                  ptrInt64Es(1),
					LicenseType:             ptrStringEs("platinum"),
					EnableDestroyProtection: ptrStringEs("CLOSE"),
				},
			},
		}
		return resp, nil
	})

	meta := newMockMetaEs()
	res := svces.ResourceTencentCloudElasticsearchInstance()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_name":             "tf-test-instance",
		"availability_zone":         "ap-guangzhou-3",
		"version":                   "7.10.1",
		"vpc_id":                    "vpc-test",
		"subnet_id":                 "subnet-test",
		"password":                  "Test1234",
		"charge_type":               "POSTPAID_BY_HOUR",
		"license_type":              "platinum",
		"enable_destroy_protection": "CLOSE",
	})
	d.SetId("es-update-destroy-protection")

	// no changes detected
	patches.ApplyMethodFunc(d, "HasChange", func(key string) bool {
		return false
	})

	err := res.Update(d, meta)
	assert.NoError(t, err)
	assert.False(t, updateCalled)
}
