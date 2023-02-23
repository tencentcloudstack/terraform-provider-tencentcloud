package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var testTkeClusterName = "tencentcloud_kubernetes_cluster"
var testTkeClusterResourceKey = testTkeClusterName + ".managed_cluster"

func init() {
	// go test -v ./tencentcloud -sweep=ap-guangzhou -sweep-run=tencentcloud_kubernetes_cluster
	resource.AddTestSweepers("tencentcloud_kubernetes_cluster", &resource.Sweeper{
		Name: "tencentcloud_kubernetes_cluster",
		F: func(r string) error {
			logId := getLogId(contextNil)
			ctx := context.WithValue(context.TODO(), logIdKey, logId)
			cli, _ := sharedClientForRegion(r)
			client := cli.(*TencentCloudClient).apiV3Conn
			service := TkeService{client: client}
			clusters, err := service.DescribeClusters(ctx, "", "")
			if err != nil {
				return err
			}

			for _, v := range clusters {
				id := v.ClusterId
				name := v.ClusterName
				createdTime, _ := time.Parse(time.RFC3339, v.CreatedTime)
				if isResourcePersist(name, &createdTime) {
					continue
				}
				if err := service.DeleteCluster(ctx, id); err != nil {
					return err
				}
			}

			return nil
		},
	})
}

func TestAccTencentCloudKubernetesClusterResourceBasic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTkeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTkeCluster,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTkeExists(testTkeClusterResourceKey),
					resource.TestCheckResourceAttr(testTkeClusterResourceKey, "cluster_cidr", "10.31.0.0/23"),
					resource.TestCheckResourceAttr(testTkeClusterResourceKey, "cluster_max_pod_num", "32"),
					resource.TestCheckResourceAttr(testTkeClusterResourceKey, "cluster_name", "test"),
					resource.TestCheckResourceAttr(testTkeClusterResourceKey, "cluster_desc", "test cluster desc"),
					resource.TestCheckResourceAttr(testTkeClusterResourceKey, "cluster_node_num", "1"),
					resource.TestCheckResourceAttr(testTkeClusterResourceKey, "worker_instances_list.#", "1"),
					resource.TestCheckResourceAttrSet(testTkeClusterResourceKey, "worker_instances_list.0.instance_id"),
					resource.TestCheckResourceAttrSet(testTkeClusterResourceKey, "certification_authority"),
					resource.TestCheckResourceAttrSet(testTkeClusterResourceKey, "user_name"),
					resource.TestCheckResourceAttrSet(testTkeClusterResourceKey, "password"),
					resource.TestCheckResourceAttr(testTkeClusterResourceKey, "tags.test", "test"),
					//resource.TestCheckResourceAttr(testTkeClusterResourceKey, "security_policy.#", "2"),
					//resource.TestCheckResourceAttrSet(testTkeClusterResourceKey, "cluster_external_endpoint"),
					resource.TestCheckResourceAttr(testTkeClusterResourceKey, "cluster_level", "L5"),
					resource.TestCheckResourceAttr(testTkeClusterResourceKey, "auto_upgrade_cluster_level", "true"),
					resource.TestCheckResourceAttr(testTkeClusterResourceKey, "labels.test1", "test1"),
					resource.TestCheckResourceAttr(testTkeClusterResourceKey, "labels.test2", "test2"),
					resource.TestCheckResourceAttr(testTkeClusterResourceKey, "cluster_internet_domain", "tf.cluster-internet.com"),
					resource.TestCheckResourceAttr(testTkeClusterResourceKey, "cluster_intranet_domain", "tf.cluster-intranet.com"),
				),
			},
			{
				Config: testAccTkeClusterUpdateAccess,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTkeExists(testTkeClusterResourceKey),
					resource.TestCheckResourceAttr(testTkeClusterResourceKey, "cluster_name", "test2"),
					resource.TestCheckResourceAttr(testTkeClusterResourceKey, "cluster_desc", "test cluster desc 2"),
					resource.TestCheckResourceAttr(testTkeClusterResourceKey, "cluster_level", "L5"),
					resource.TestCheckResourceAttr(testTkeClusterResourceKey, "cluster_internet_domain", "tf2.cluster-internet.com"),
				),
			},
			{
				Config: testAccTkeClusterUpdateLevel,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTkeExists(testTkeClusterResourceKey),
					resource.TestCheckResourceAttr(testTkeClusterResourceKey, "cluster_desc", "test cluster desc 3"),
					resource.TestCheckResourceAttr(testTkeClusterResourceKey, "cluster_level", "L20"),
					resource.TestCheckResourceAttr(testTkeClusterResourceKey, "auto_upgrade_cluster_level", "false"),
				),
			},
		},
	})
}

func TestAccTencentCloudKubernetesClusterResourceLogsAddons(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTkeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTkeClusterLogsAddons,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTkeExists(testTkeClusterResourceKey),
					resource.TestCheckResourceAttr(testTkeClusterResourceKey, "cluster_cidr", "192.168.0.0/18"),
					resource.TestCheckResourceAttr(testTkeClusterResourceKey, "cluster_name", "test"),
					resource.TestCheckResourceAttr(testTkeClusterResourceKey, "cluster_desc", "test cluster desc"),
					resource.TestCheckResourceAttr(testTkeClusterResourceKey, "log_agent.0.enabled", "true"),
					resource.TestCheckResourceAttr(testTkeClusterResourceKey, "event_persistence.0.enabled", "true"),
					resource.TestCheckResourceAttr(testTkeClusterResourceKey, "cluster_audit.0.enabled", "false"),
				),
			},
			{
				PreConfig: func() {
					// do not update so fast
					time.Sleep(10 * time.Second)
				},
				Config: testAccTkeClusterLogsAddonsUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTkeExists(testTkeClusterResourceKey),
					resource.TestCheckResourceAttr(testTkeClusterResourceKey, "cluster_cidr", "192.168.0.0/18"),
					resource.TestCheckResourceAttr(testTkeClusterResourceKey, "cluster_name", "test"),
					resource.TestCheckResourceAttr(testTkeClusterResourceKey, "cluster_desc", "test cluster desc"),
					resource.TestCheckResourceAttr(testTkeClusterResourceKey, "log_agent.0.enabled", "true"),
					resource.TestCheckResourceAttr(testTkeClusterResourceKey, "event_persistence.0.enabled", "false"),
					resource.TestCheckResourceAttr(testTkeClusterResourceKey, "event_persistence.0.delete_event_log_and_topic",
						"true"),
					resource.TestCheckResourceAttr(testTkeClusterResourceKey, "cluster_audit.0.enabled", "true"),
					resource.TestCheckResourceAttr(testTkeClusterResourceKey, "cluster_audit.0.delete_audit_log_and_topic",
						"true"),
				),
			},
		},
	})
}

func TestUnitTkeAddonDiff(t *testing.T) {
	t.Parallel()
	addons1 := []interface{}{
		map[string]interface{}{
			"name":  "tcr",
			"param": `{ "version": "1.0" }`,
		},
		map[string]interface{}{
			"name":  "cos",
			"param": `{ "version": "1.2" }`,
		},
		map[string]interface{}{
			"name":  "oom_guard",
			"param": `{ "version": "1.2" }`,
		},
		map[string]interface{}{
			"name":  "npdplus",
			"param": `{ "version": "1.1" }`,
		},
	}

	addons2 := []interface{}{
		map[string]interface{}{
			"name":  "tcr",
			"param": `{ "version": "1.0" }`,
		},
		map[string]interface{}{
			"name":  "oom_guard",
			"param": `{ "version": "2.0" }`,
		},
		map[string]interface{}{
			"name":  "prom",
			"param": `{ "version": "1.1" }`,
		},
		map[string]interface{}{
			"name":  "npdplus",
			"param": `{ "version": "1.3" }`,
		},
	}

	adds, removes, changes := resourceTkeGetAddonsDiffs(addons1, addons2)

	assert.Len(t, adds, 1)
	assert.Len(t, removes, 1)
	assert.Len(t, changes, 2)

	assert.Contains(t, adds, map[string]interface{}{
		"name":  "prom",
		"param": `{ "version": "1.1" }`,
	})

	assert.Contains(t, removes, map[string]interface{}{
		"name":  "cos",
		"param": `{ "version": "1.2" }`,
	})

	assert.Contains(t, changes, map[string]interface{}{
		"name":  "oom_guard",
		"param": `{ "version": "2.0" }`,
	})

	assert.Contains(t, changes, map[string]interface{}{
		"name":  "npdplus",
		"param": `{ "version": "1.3" }`,
	})
}

func testAccCheckTkeDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TkeService{
		client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != testTkeClusterName {
			continue
		}
		_, has, err := service.DescribeCluster(ctx, rs.Primary.ID)
		if err != nil {
			err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
				_, has, err = service.DescribeCluster(ctx, rs.Primary.ID)
				if err != nil {
					code := err.(*sdkErrors.TencentCloudSDKError).Code
					if code == "ResourceUnavailable.ClusterState" {
						return nil
					}
					return retryError(err)
				}
				return nil
			})
		}

		if err != nil {
			return nil
		}

		if !has {
			log.Printf("[DEBUG]tke cluster  %s delete  ok", rs.Primary.ID)
			return nil
		} else {
			return fmt.Errorf("tke cluster delete fail,%s", rs.Primary.ID)
		}

	}
	return nil
}

func testAccCheckTkeExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("tke cluster %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("tke cluster id is not set")
		}

		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		service := TkeService{
			client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
		}

		_, has, err := service.DescribeCluster(ctx, rs.Primary.ID)
		if err != nil {
			err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
				_, has, err = service.DescribeCluster(ctx, rs.Primary.ID)
				if err != nil {
					return retryError(err)
				}
				return nil
			})
		}

		if err != nil {
			return nil
		}
		if !has {
			return fmt.Errorf("tke cluster create fail")
		} else {
			log.Printf("[DEBUG]tke cluster  %s create  ok", rs.Primary.ID)
			return nil
		}

	}
}

const testAccTkeExtensionAddons = `
variable "addons" {
  default = [{
    name  = "CFS",
    param = {
      "kind" : "App", "spec" : {
        "chart" : { "chartName" : "cfs", "chartVersion" : "1.0.7" },
        "values" : { "values" : [], "rawValues" : "e30=", "rawValuesType" : "json" }
      }
    }
  },
    {
      name  = "OOMGuard",
      param = {
        "kind" : "App", "spec" : { "chart" : { "chartName" : "oomguard", "chartVersion" : "1.0.1" } }
      }
    }]
}

variable "addons_update" {
  default = [{
    name  = "CFS",
    param = {
      "kind" : "App", "spec" : {
        "chart" : { "chartName" : "cfs", "chartVersion" : "1.0.8" },
        "values" : { "values" : [], "rawValues" : "e30=", "rawValuesType" : "json" }
      }
    }
  },
    {
      name  = "OOMGuard",
      param = {
        "kind" : "App", "spec" : { "chart" : { "chartName" : "oomguard", "chartVersion" : "1.0.1" } }
      }
    },
    {
      name  = "cos",
      param = {
        "kind" : "App", "spec" : { "chart" : { "chartName" : "cos", "chartVersion" : "1.0.1" } }
      }
    }]
}
`

const TkeDeps = TkeExclusiveNetwork + TkeInstanceType + TkeCIDRs + defaultImages + defaultSecurityGroupData

const testAccTkeCluster = TkeDeps + `
variable "availability_zone" {
  default = "ap-guangzhou-3"
}

resource "tencentcloud_kubernetes_cluster" "managed_cluster" {
  vpc_id                                     = local.vpc_id
  cluster_cidr                               = var.tke_cidr_a.0
  cluster_max_pod_num                        = 32
  cluster_name                               = "test"
  cluster_desc                               = "test cluster desc"
  cluster_max_service_num                    = 32
  cluster_internet                           = true
  cluster_internet_domain                    = "tf.cluster-internet.com"
  cluster_intranet                           = true
  cluster_intranet_domain                    = "tf.cluster-intranet.com"
  cluster_version                            = "1.18.4"
  cluster_os                                 = "tlinux2.2(tkernel3)x86_64"
  cluster_level								 = "L5"
  auto_upgrade_cluster_level				 = true
  cluster_intranet_subnet_id                 = local.subnet_id
  cluster_internet_security_group               = local.sg_id
  managed_cluster_internet_security_policies = ["3.3.3.3", "1.1.1.1"]
  worker_config {
    count                      = 1
    availability_zone          = var.availability_zone
    instance_type              = local.final_type
    system_disk_type           = "CLOUD_SSD"
    system_disk_size           = 60
    internet_charge_type       = "TRAFFIC_POSTPAID_BY_HOUR"
    internet_max_bandwidth_out = 100
    public_ip_assigned         = true
    subnet_id                  = local.subnet_id
    img_id                     = var.default_img_id
    security_group_ids         = [local.sg_id]

    data_disk {
      disk_type = "CLOUD_PREMIUM"
      disk_size = 50
      file_system = "ext3"
	  auto_format_and_mount = "true"
	  mount_target = "/var/lib/docker"
      disk_partition = "/dev/sdb1"
    }

    enhanced_security_service = false
    enhanced_monitor_service  = false
    user_data                 = "dGVzdA=="
    password                  = "ZZXXccvv1212"
  }

  cluster_deploy_type = "MANAGED_CLUSTER"

  tags = {
    "test" = "test"
  }

  unschedulable = 0

  labels = {
    "test1" = "test1",
    "test2" = "test2",
  }
  extra_args = [
 	"root-dir=/var/lib/kubelet"
  ]
}
`

const testAccTkeClusterUpdateAccess = TkeDeps + `
variable "availability_zone" {
  default = "ap-guangzhou-3"
}

resource "tencentcloud_kubernetes_cluster" "managed_cluster" {
  vpc_id                                     = local.vpc_id
  cluster_cidr                               = var.tke_cidr_a.0
  cluster_max_pod_num                        = 32
  cluster_name                               = "test2"
  cluster_desc                               = "test cluster desc 2"
  cluster_max_service_num                    = 32
  cluster_internet                           = true
  cluster_internet_domain                    = "tf2.cluster-internet.com"
  cluster_intranet                           = false
  cluster_version                            = "1.18.4"
  cluster_os                                 = "tlinux2.2(tkernel3)x86_64"
  cluster_level								 = "L5"
  cluster_internet_security_group               = local.sg_id2
  auto_upgrade_cluster_level				 = true
  managed_cluster_internet_security_policies = ["3.3.3.3", "1.1.1.1"]
  worker_config {
    count                      = 1
    availability_zone          = var.availability_zone
    instance_type              = local.final_type
    system_disk_type           = "CLOUD_SSD"
    system_disk_size           = 60
    internet_charge_type       = "TRAFFIC_POSTPAID_BY_HOUR"
    internet_max_bandwidth_out = 100
    public_ip_assigned         = true
    subnet_id                  = local.subnet_id
    img_id                     = var.default_img_id
    security_group_ids         = [local.sg_id]

    data_disk {
      disk_type = "CLOUD_PREMIUM"
      disk_size = 50
      file_system = "ext3"
	  auto_format_and_mount = "true"
	  mount_target = "/var/lib/docker"
      disk_partition = "/dev/sdb1"
    }

    enhanced_security_service = false
    enhanced_monitor_service  = false
    user_data                 = "dGVzdA=="
    password                  = "ZZXXccvv1212"
  }

  cluster_deploy_type = "MANAGED_CLUSTER"

  tags = {
    "test" = "test"
  }

  unschedulable = 0

  labels = {
    "test1" = "test1",
    "test2" = "test2",
  }
  extra_args = [
 	"root-dir=/var/lib/kubelet"
  ]
}
`
const testAccTkeClusterUpdateLevel = TkeDeps + `
variable "availability_zone" {
  default = "ap-guangzhou-3"
}

resource "tencentcloud_kubernetes_cluster" "managed_cluster" {
  vpc_id                                     = local.vpc_id
  cluster_cidr                               = var.tke_cidr_a.0
  cluster_max_pod_num                        = 32
  cluster_name                               = "test2"
  cluster_desc                               = "test cluster desc 3"
  cluster_max_service_num                    = 32
  cluster_internet                           = false
  cluster_version                            = "1.18.4"
  cluster_os                                 = "tlinux2.2(tkernel3)x86_64"
  cluster_level								 = "L20"
  auto_upgrade_cluster_level				 = false
  managed_cluster_internet_security_policies = ["3.3.3.3", "1.1.1.1"]
  worker_config {
    count                      = 1
    availability_zone          = var.availability_zone
    instance_type              = local.final_type
    system_disk_type           = "CLOUD_SSD"
    system_disk_size           = 60
    internet_charge_type       = "TRAFFIC_POSTPAID_BY_HOUR"
    internet_max_bandwidth_out = 100
    public_ip_assigned         = true
    subnet_id                  = local.subnet_id
    img_id                     = var.default_img_id
    security_group_ids         = [local.sg_id]

    data_disk {
      disk_type = "CLOUD_PREMIUM"
      disk_size = 50
      file_system = "ext3"
	  auto_format_and_mount = "true"
	  mount_target = "/var/lib/docker"
      disk_partition = "/dev/sdb1"
    }

    enhanced_security_service = false
    enhanced_monitor_service  = false
    user_data                 = "dGVzdA=="
    password                  = "ZZXXccvv1212"
  }

  cluster_deploy_type = "MANAGED_CLUSTER"

  tags = {
    "abc" = "abc"
  }

  unschedulable = 0

  labels = {
    "test1" = "test1",
    "test2" = "test2",
  }
  extra_args = [
 	"root-dir=/var/lib/kubelet"
  ]
}
`

const testAccTkeClusterLogsAddons = TkeDeps + testAccTkeExtensionAddons + `
variable "availability_zone" {
  default = "ap-guangzhou-3"
}

resource "tencentcloud_kubernetes_cluster" "managed_cluster" {
  vpc_id                                     = local.vpc_id
  cluster_cidr                               = var.tke_cidr_c.0
  cluster_max_pod_num                        = 32
  cluster_name                               = "test"
  cluster_desc                               = "test cluster desc"
  cluster_max_service_num                    = 32
  cluster_version                            = "1.20.6"
  cluster_os                                 = "tlinux2.2(tkernel3)x86_64"
  cluster_level								 = "L5"
  auto_upgrade_cluster_level				 = true
  cluster_deploy_type 						 = "MANAGED_CLUSTER"

  worker_config {
    count                      = 1
    availability_zone          = var.availability_zone
    instance_type              = local.final_type
    system_disk_type           = "CLOUD_SSD"
    system_disk_size           = 60
    internet_charge_type       = "TRAFFIC_POSTPAID_BY_HOUR"
    internet_max_bandwidth_out = 10
    public_ip_assigned         = true
    subnet_id                  = local.subnet_id
    img_id                     = var.default_img_id
    security_group_ids         = [local.sg_id]
    enhanced_security_service = false
    enhanced_monitor_service  = false
    user_data                 = "dGVzdA=="
    password                  = "ZZXXccvv1212"
  }

  dynamic "extension_addon" {
    for_each = var.addons
    content {
      name = extension_addon.value.name
      param = jsonencode(extension_addon.value.param)
    }
  }

  log_agent {
    enabled = true
  }

  event_persistence {
    enabled = true
  }

  cluster_audit {
    enabled = false
  }
}`

const testAccTkeClusterLogsAddonsUpdate = TkeDeps + testAccTkeExtensionAddons + `
variable "availability_zone" {
  default = "ap-guangzhou-3"
}

resource "tencentcloud_kubernetes_cluster" "managed_cluster" {
  vpc_id                                     = local.vpc_id
  cluster_cidr                               = var.tke_cidr_c.0
  cluster_max_pod_num                        = 32
  cluster_name                               = "test"
  cluster_desc                               = "test cluster desc"
  cluster_max_service_num                    = 32
  cluster_version                            = "1.20.6"
  cluster_os                                 = "tlinux2.2(tkernel3)x86_64"
  cluster_level								 = "L5"
  auto_upgrade_cluster_level				 = true
  cluster_deploy_type 						 = "MANAGED_CLUSTER"

  dynamic "extension_addon" {
    for_each = var.addons_update
    content {
      name = extension_addon.value.name
      param = jsonencode(extension_addon.value.param)
    }
  }

  worker_config {
    count                      = 1
    availability_zone          = var.availability_zone
    instance_type              = local.final_type
    system_disk_type           = "CLOUD_SSD"
    system_disk_size           = 60
    internet_charge_type       = "TRAFFIC_POSTPAID_BY_HOUR"
    internet_max_bandwidth_out = 10
    public_ip_assigned         = true
    subnet_id                  = local.subnet_id
    img_id                     = var.default_img_id
    security_group_ids         = [local.sg_id]
    enhanced_security_service = false
    enhanced_monitor_service  = false
    user_data                 = "dGVzdA=="
    password                  = "ZZXXccvv1212"
  }

  log_agent {
    enabled = true
  }

  event_persistence {
    enabled = false
    delete_event_log_and_topic = true
  }

  cluster_audit {
    enabled = true
    delete_audit_log_and_topic = true
  }
}`
