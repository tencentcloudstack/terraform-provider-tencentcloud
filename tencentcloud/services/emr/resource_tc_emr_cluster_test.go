package emr_test

import (
	"context"
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	emr "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/emr/v20190103"

	svccdb "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cdb"
	svcemr "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/emr"
)

func init() {
	// go test -v ./tencentcloud -sweep=ap-guangzhou -sweep-run=tencentcloud_emr
	resource.AddTestSweepers("tencentcloud_emr", &resource.Sweeper{
		Name: "tencentcloud_emr",
		F: func(r string) error {
			logId := tccommon.GetLogId(tccommon.ContextNil)
			ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
			sharedClient, err := tcacctest.SharedClientForRegion(r)
			if err != nil {
				return fmt.Errorf("getting tencentcloud client error: %s", err.Error())
			}
			client := sharedClient.(tccommon.ProviderMeta).GetAPIV3Conn()

			emrService := svcemr.NewEMRService(client)
			filters := make(map[string]interface{})
			filters["display_strategy"] = svcemr.DisplayStrategyIsclusterList
			clusters, err := emrService.DescribeInstances(ctx, filters)
			if err != nil {
				return nil
			}

			// add scanning resources
			var resources, nonKeepResources []*tccommon.ResourceInstance
			for _, v := range clusters {
				if !tccommon.CheckResourcePersist(*v.ClusterId, *v.AddTime) {
					nonKeepResources = append(nonKeepResources, &tccommon.ResourceInstance{
						Id:   *v.ClusterId,
						Name: *v.ClusterName,
					})
				}
				resources = append(resources, &tccommon.ResourceInstance{
					Id:         *v.ClusterId,
					Name:       *v.ClusterName,
					CreateTime: *v.AddTime,
				})
			}
			tccommon.ProcessScanCloudResources(client, resources, nonKeepResources, "CreateInstance")

			for _, cluster := range clusters {
				clusterName := *cluster.ClusterName
				if strings.HasPrefix(clusterName, tcacctest.KeepResource) || strings.HasPrefix(clusterName, tcacctest.DefaultResource) {
					continue
				}
				now := time.Now()
				createTime := tccommon.StringToTime(*cluster.AddTime)
				interval := now.Sub(createTime).Minutes()
				// less than 30 minute, not delete
				if tccommon.NeedProtect == 1 && int64(interval) < 30 {
					continue
				}
				metaDB := cluster.MetaDb
				instanceId := *cluster.ClusterId
				request := emr.NewTerminateInstanceRequest()
				request.InstanceId = &instanceId
				if _, err = client.UseEmrClient().TerminateInstance(request); err != nil {
					return nil
				}
				err = resource.Retry(10*tccommon.ReadRetryTimeout, func() *resource.RetryError {
					clusters, err := emrService.DescribeInstancesById(ctx, instanceId, svcemr.DisplayStrategyIsclusterList)

					if e, ok := err.(*errors.TencentCloudSDKError); ok {
						if e.GetCode() == "InternalError.ClusterNotFound" {
							return nil
						}
						if e.GetCode() == "UnauthorizedOperation" {
							return nil
						}
					}

					if len(clusters) > 0 {
						status := *(clusters[0].Status)
						if status != svcemr.EmrInternetStatusDeleted {
							return resource.RetryableError(
								fmt.Errorf("%v create cluster endpoint status still is %v", instanceId, status))
						}
					}

					if err != nil {
						return resource.RetryableError(err)
					}
					return nil
				})
				if err != nil {
					return nil
				}

				if metaDB != nil && *metaDB != "" {
					// remove metadb
					mysqlService := svccdb.NewMysqlService(client)

					err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
						err := mysqlService.OfflineIsolatedInstances(ctx, *metaDB)
						if err != nil {
							return tccommon.RetryError(err, tccommon.InternalError)
						}
						return nil
					})

					if err != nil {
						return nil
					}
				}
			}
			return nil
		},
	})
}

var testEmrClusterResourceKey = "tencentcloud_emr_cluster.emrrrr"

func TestAccTencentCloudEmrClusterResource_Basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testEmrBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEmrExists(testEmrClusterResourceKey),
					resource.TestCheckResourceAttr(testEmrClusterResourceKey, "product_id", "38"),
					resource.TestCheckResourceAttr(testEmrClusterResourceKey, "vpc_settings.vpc_id", tcacctest.DefaultEMRVpcId),
					resource.TestCheckResourceAttr(testEmrClusterResourceKey, "vpc_settings.subnet_id", tcacctest.DefaultEMRSubnetId),
					resource.TestCheckResourceAttr(testEmrClusterResourceKey, "softwares.#", "5"),
					resource.TestCheckResourceAttr(testEmrClusterResourceKey, "support_ha", "0"),
					resource.TestCheckResourceAttr(testEmrClusterResourceKey, "instance_name", "emr-test-demo"),
					resource.TestCheckResourceAttr(testEmrClusterResourceKey, "resource_spec.#", "1"),
					resource.TestCheckResourceAttr(testEmrClusterResourceKey, "login_settings.password", "Tencent@cloud123"),
					resource.TestCheckResourceAttr(testEmrClusterResourceKey, "time_span", "3600"),
					resource.TestCheckResourceAttr(testEmrClusterResourceKey, "time_unit", "s"),
					resource.TestCheckResourceAttr(testEmrClusterResourceKey, "pay_mode", "0"),
					resource.TestCheckResourceAttr(testEmrClusterResourceKey, "placement_info.0.zone", "ap-guangzhou-3"),
					resource.TestCheckResourceAttr(testEmrClusterResourceKey, "placement_info.0.project_id", "0"),
					resource.TestCheckResourceAttrSet(testEmrClusterResourceKey, "instance_id"),
					resource.TestCheckResourceAttr(testEmrClusterResourceKey, "sg_id", tcacctest.DefaultEMRSgId),
					resource.TestCheckResourceAttr(testEmrClusterResourceKey, "tags.emr-key", "emr-value"),
					resource.TestCheckResourceAttr(testEmrClusterResourceKey, "resource_spec.0.core_count", "2"),
					resource.TestCheckResourceAttr(testEmrClusterResourceKey, "resource_spec.0.master_resource_spec.0.multi_disks.#", "1"),
					resource.TestCheckResourceAttr(testEmrClusterResourceKey, "resource_spec.0.core_resource_spec.0.multi_disks.#", "1"),
				),
			},
			{
				Config: testEmrBasic_AddCoreNode,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEmrExists(testEmrClusterResourceKey),
					resource.TestCheckResourceAttr(testEmrClusterResourceKey, "resource_spec.0.core_count", "3"),
				),
			},
			{
				ResourceName:            testEmrClusterResourceKey,
				ImportState:             true,
				ImportStateVerifyIgnore: []string{"display_strategy", "placement", "time_span", "time_unit", "login_settings", "terminate_node_info"},
			},
		},
	})
}

func TestAccTencentCloudEmrClusterResource_Zookeeper(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testEmrZookeeper,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEmrExists("tencentcloud_emr_cluster.emr_zookeeper"),
					resource.TestCheckResourceAttr("tencentcloud_emr_cluster.emr_zookeeper", "scene_name", "Hadoop-Zookeeper"),
					resource.TestCheckResourceAttr("tencentcloud_emr_cluster.emr_zookeeper", "resource_spec.0.common_count", "3"),
				),
			},
			{
				ResourceName:            "tencentcloud_emr_cluster.emr_zookeeper",
				ImportState:             true,
				ImportStateVerifyIgnore: []string{"display_strategy", "placement", "time_span", "time_unit", "login_settings", "terminate_node_info"},
			},
		},
	})
}

func TestAccTencentCloudEmrClusterResource_PreExecutedFileSettings(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testEmrBasicPreExecutedFileSettings,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEmrExists(testEmrClusterResourceKey),
					resource.TestCheckResourceAttr(testEmrClusterResourceKey, "pre_executed_file_settings.#", "1"),
					resource.TestCheckResourceAttr(testEmrClusterResourceKey, "pre_executed_file_settings.0.cos_file_name", "test"),
					resource.TestCheckResourceAttr(testEmrClusterResourceKey, "pre_executed_file_settings.0.when_run", "resourceAfter"),
					resource.TestCheckResourceAttrSet(testEmrClusterResourceKey, "pre_executed_file_settings.0.cos_file_uri"),
				),
			},
		},
	})
}
func TestAccTencentCloudEmrClusterResource_Prepay(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testEmrBasicPrepay,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEmrExists(testEmrClusterResourceKey),
					resource.TestCheckResourceAttr(testEmrClusterResourceKey, "product_id", "38"),
					resource.TestCheckResourceAttr(testEmrClusterResourceKey, "vpc_settings.vpc_id", tcacctest.DefaultEMRVpcId),
					resource.TestCheckResourceAttr(testEmrClusterResourceKey, "vpc_settings.subnet_id", tcacctest.DefaultEMRSubnetId),
					resource.TestCheckResourceAttr(testEmrClusterResourceKey, "softwares.#", "5"),
					resource.TestCheckResourceAttr(testEmrClusterResourceKey, "support_ha", "0"),
					resource.TestCheckResourceAttr(testEmrClusterResourceKey, "instance_name", "emr-test-demo"),
					resource.TestCheckResourceAttr(testEmrClusterResourceKey, "resource_spec.#", "1"),
					resource.TestCheckResourceAttr(testEmrClusterResourceKey, "login_settings.password", "Tencent@cloud123"),
					resource.TestCheckResourceAttr(testEmrClusterResourceKey, "time_span", "1"),
					resource.TestCheckResourceAttr(testEmrClusterResourceKey, "time_unit", "m"),
					resource.TestCheckResourceAttr(testEmrClusterResourceKey, "pay_mode", "1"),
					resource.TestCheckResourceAttr(testEmrClusterResourceKey, "placement_info.0.zone", "ap-guangzhou-3"),
					resource.TestCheckResourceAttr(testEmrClusterResourceKey, "placement_info.0.project_id", "0"),
					resource.TestCheckResourceAttrSet(testEmrClusterResourceKey, "instance_id"),
					resource.TestCheckResourceAttr(testEmrClusterResourceKey, "sg_id", tcacctest.DefaultEMRSgId),
					resource.TestCheckResourceAttr(testEmrClusterResourceKey, "tags.emr-key", "emr-value"),
					resource.TestCheckResourceAttr(testEmrClusterResourceKey, "auto_renew", "1"),
				),
			},
		},
	})
}

func TestAccTencentCloudEmrClusterResource_NotNeedMasterWan(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testEmrNotNeedMasterWan,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEmrExists(testEmrClusterResourceKey),
					resource.TestCheckResourceAttr(testEmrClusterResourceKey, "need_master_wan", "NOT_NEED_MASTER_WAN"),
				),
			},
		},
	})
}

func TestAccTencentCloudEmrClusterResource_multiZone(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testEmrMultiZone,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEmrExists(testEmrClusterResourceKey),
					resource.TestCheckResourceAttr(testEmrClusterResourceKey, "multi_zone", "true"),
					resource.TestCheckResourceAttr(testEmrClusterResourceKey, "multi_zone_setting.#", "2"),
				),
			},
		},
	})
}

func testAccCheckEmrExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("emr cluster %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("emr cluster id is not set")
		}

		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		service := svcemr.NewEMRService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())

		instanceId := rs.Primary.ID
		clusters, err := service.DescribeInstancesById(ctx, instanceId, svcemr.DisplayStrategyIsclusterList)
		if err != nil {
			err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
				clusters, err = service.DescribeInstancesById(ctx, instanceId, svcemr.DisplayStrategyIsclusterList)
				if err != nil {
					return tccommon.RetryError(err)
				}
				return nil
			})
		}

		if err != nil {
			return nil
		}
		if len(clusters) <= 0 {
			return fmt.Errorf("emr cluster create fail")
		} else {
			log.Printf("[DEBUG]emr cluster  %s create  ok", rs.Primary.ID)
			return nil
		}

	}
}

const testEmrBasic = tcacctest.DefaultEMRVariable + `
data "tencentcloud_instance_types" "cvm4c8m" {
	exclude_sold_out=true
	cpu_core_count=4
	memory_size=8
    filter {
      name   = "instance-charge-type"
      values = ["POSTPAID_BY_HOUR"]
    }
    filter {
    name   = "zone"
    values = ["ap-guangzhou-3"]
  }
}

resource "tencentcloud_emr_cluster" "emrrrr" {
	product_id=38
	vpc_settings={
	  vpc_id=var.vpc_id
	  subnet_id=var.subnet_id
	}
	softwares = [
	  "hdfs-2.8.5",
	  "knox-1.6.1",
	  "openldap-2.4.44",
	  "yarn-2.8.5",
	  "zookeeper-3.6.3",
	]
	support_ha=0
	instance_name="emr-test-demo"
	resource_spec {
	  master_resource_spec {
		mem_size=8192
		cpu=4
		disk_size=100
		disk_type="CLOUD_PREMIUM"
		spec="CVM.${data.tencentcloud_instance_types.cvm4c8m.instance_types.0.family}"
		storage_type=5
		root_size=50
		multi_disks {
			disk_type = "CLOUD_PREMIUM"
			volume = 200
			count = 1
		}
	  }
	  core_resource_spec {
		mem_size=8192
		cpu=4
		disk_size=100
		disk_type="CLOUD_PREMIUM"
		spec="CVM.${data.tencentcloud_instance_types.cvm4c8m.instance_types.0.family}"
		storage_type=5
		root_size=50
		multi_disks {
			disk_type = "CLOUD_PREMIUM"
			volume = 100
			count = 2
		}
	  }
	  master_count=1
	  core_count=2
	}
	login_settings={
	  password="Tencent@cloud123"
	}
	time_span=3600
	time_unit="s"
	pay_mode=0
	placement_info {
	  zone="ap-guangzhou-3"
	  project_id=0
	}
	sg_id=var.sg_id
	tags = {
        emr-key = "emr-value"
    }
  }
`

const testEmrBasicPreExecutedFileSettings = tcacctest.DefaultEMRVariable + `
data "tencentcloud_instance_types" "cvm4c8m" {
	exclude_sold_out=true
	cpu_core_count=4
	memory_size=8
    filter {
      name   = "instance-charge-type"
      values = ["POSTPAID_BY_HOUR"]
    }
    filter {
    name   = "zone"
    values = ["ap-guangzhou-3"]
  }
}

resource "tencentcloud_emr_cluster" "emrrrr" {
	product_id=38
	vpc_settings={
	  vpc_id=var.vpc_id
	  subnet_id=var.subnet_id
	}
	softwares = [
	  "hdfs-2.8.5",
	  "knox-1.6.1",
	  "openldap-2.4.44",
	  "yarn-2.8.5",
	  "zookeeper-3.6.3",
	]
	support_ha=0
	instance_name="emr-test-demo"
	resource_spec {
	  master_resource_spec {
		mem_size=8192
		cpu=4
		disk_size=100
		disk_type="CLOUD_PREMIUM"
		spec="CVM.${data.tencentcloud_instance_types.cvm4c8m.instance_types.0.family}"
		storage_type=5
		root_size=50
		multi_disks {
			disk_type = "CLOUD_PREMIUM"
			volume = 200
			count = 1
		}
	  }
	  core_resource_spec {
		mem_size=8192
		cpu=4
		disk_size=100
		disk_type="CLOUD_PREMIUM"
		spec="CVM.${data.tencentcloud_instance_types.cvm4c8m.instance_types.0.family}"
		storage_type=5
		root_size=50
		multi_disks {
			disk_type = "CLOUD_PREMIUM"
			volume = 100
			count = 2
		}
	  }
	  master_count=1
	  core_count=2
	}
	login_settings={
	  password="Tencent@cloud123"
	}
	time_span=3600
	time_unit="s"
	pay_mode=0
	placement_info {
	  zone="ap-guangzhou-3"
	  project_id=0
	}
	sg_id=var.sg_id
	tags = {
        emr-key = "emr-value"
    }
	pre_executed_file_settings {
		cos_file_name = "test"
		cos_file_uri = "https://keep-tf-test-1308726196.cos.ap-guangzhou.myqcloud.com/test/tmp.sh"
		when_run = "resourceAfter"
	}
  }
`

const testEmrBasic_AddCoreNode = tcacctest.DefaultEMRVariable + `
data "tencentcloud_instance_types" "cvm4c8m" {
	exclude_sold_out=true
	cpu_core_count=4
	memory_size=8
    filter {
      name   = "instance-charge-type"
      values = ["POSTPAID_BY_HOUR"]
    }
    filter {
    name   = "zone"
    values = ["ap-guangzhou-3"]
  }
}

resource "tencentcloud_emr_cluster" "emrrrr" {
	product_id=38
	vpc_settings={
	  vpc_id=var.vpc_id
	  subnet_id=var.subnet_id
	}
	softwares = [
	  "hdfs-2.8.5",
	  "knox-1.6.1",
	  "openldap-2.4.44",
	  "yarn-2.8.5",
	  "zookeeper-3.6.3",
	]
	support_ha=0
	instance_name="emr-test-demo"
	resource_spec {
	  master_resource_spec {
		mem_size=8192
		cpu=4
		disk_size=100
		disk_type="CLOUD_PREMIUM"
		spec="CVM.${data.tencentcloud_instance_types.cvm4c8m.instance_types.0.family}"
		storage_type=5
		root_size=50
		multi_disks {
			disk_type = "CLOUD_PREMIUM"
			volume = 200
			count = 1
		}
	  }
	  core_resource_spec {
		mem_size=8192
		cpu=4
		disk_size=100
		disk_type="CLOUD_PREMIUM"
		spec="CVM.${data.tencentcloud_instance_types.cvm4c8m.instance_types.0.family}"
		storage_type=5
		root_size=50
		multi_disks {
			disk_type = "CLOUD_PREMIUM"
			volume = 100
			count = 2
		}
	  }
	  master_count=1
	  core_count=3
	}
	login_settings={
	  password="Tencent@cloud123"
	}
	time_span=3600
	time_unit="s"
	pay_mode=0
	placement_info {
	  zone="ap-guangzhou-3"
	  project_id=0
	}
	sg_id=var.sg_id
	tags = {
        emr-key = "emr-value"
    }
  }
`

const testEmrBasicPrepay = tcacctest.DefaultEMRVariable + `
data "tencentcloud_instance_types" "cvm4c8m" {
	exclude_sold_out=true
	cpu_core_count=4
	memory_size=8
    filter {
      name   = "instance-charge-type"
      values = ["POSTPAID_BY_HOUR"]
    }
    filter {
    name   = "zone"
    values = ["ap-guangzhou-3"]
  }
}

resource "tencentcloud_emr_cluster" "emrrrr" {
	product_id=38
	vpc_settings={
	  vpc_id=var.vpc_id
	  subnet_id=var.subnet_id
	}
	softwares = [
	  "hdfs-2.8.5",
	  "knox-1.6.1",
	  "openldap-2.4.44",
	  "yarn-2.8.5",
	  "zookeeper-3.6.3",
	]
	support_ha=0
	instance_name="emr-test-demo"
	resource_spec {
	  master_resource_spec {
		mem_size=8192
		cpu=4
		disk_size=100
		disk_type="CLOUD_PREMIUM"
		spec="CVM.${data.tencentcloud_instance_types.cvm4c8m.instance_types.0.family}"
		storage_type=5
		root_size=50
	  }
	  core_resource_spec {
		mem_size=8192
		cpu=4
		disk_size=100
		disk_type="CLOUD_PREMIUM"
		spec="CVM.${data.tencentcloud_instance_types.cvm4c8m.instance_types.0.family}"
		storage_type=5
		root_size=50
	  }
	  master_count=1
	  core_count=2
	}
	login_settings={
	  password="Tencent@cloud123"
	}
	time_span=1
	time_unit="m"
	pay_mode=1
	placement_info {
	  zone="ap-guangzhou-3"
	  project_id=0
	}
	sg_id=var.sg_id
	tags = {
        emr-key = "emr-value"
    }
	auto_renew=1
  }
`
const testEmrNotNeedMasterWan = tcacctest.DefaultEMRVariable + `
data "tencentcloud_instance_types" "cvm4c8m" {
	exclude_sold_out=true
	cpu_core_count=4
	memory_size=8
    filter {
      name   = "instance-charge-type"
      values = ["POSTPAID_BY_HOUR"]
    }
    filter {
    name   = "zone"
    values = ["ap-guangzhou-3"]
  }
}

resource "tencentcloud_emr_cluster" "emrrrr" {
	product_id=38
	vpc_settings={
	  vpc_id=var.vpc_id
	  subnet_id=var.subnet_id
	}
	softwares = [
	  "hdfs-2.8.5",
	  "knox-1.6.1",
	  "openldap-2.4.44",
	  "yarn-2.8.5",
	  "zookeeper-3.6.3",
	]
	support_ha=0
	instance_name="emr-test-demo"
	resource_spec {
	  master_resource_spec {
		mem_size=8192
		cpu=4
		disk_size=100
		disk_type="CLOUD_PREMIUM"
		spec="CVM.${data.tencentcloud_instance_types.cvm4c8m.instance_types.0.family}"
		storage_type=5
		root_size=50
	  }
	  core_resource_spec {
		mem_size=8192
		cpu=4
		disk_size=100
		disk_type="CLOUD_PREMIUM"
		spec="CVM.${data.tencentcloud_instance_types.cvm4c8m.instance_types.0.family}"
		storage_type=5
		root_size=50
	  }
	  master_count=1
	  core_count=2
	}
	login_settings={
	  password="Tencent@cloud123"
	}
	time_span=3600
	time_unit="s"
	pay_mode=0
	placement_info {
	  zone="ap-guangzhou-3"
	  project_id=0
	}
	sg_id=var.sg_id
	tags = {
        emr-key = "emr-value"
    }
	need_master_wan = "NOT_NEED_MASTER_WAN"
  }
`

const testEmrZookeeper = tcacctest.DefaultEMRVariable + `
data "tencentcloud_instance_types" "cvm2c4m" {
	exclude_sold_out=true
	cpu_core_count=2
	memory_size=4
    filter {
      name   = "instance-charge-type"
      values = ["POSTPAID_BY_HOUR"]
    }
    filter {
    name   = "zone"
    values = ["ap-guangzhou-3"]
  }
}

resource "tencentcloud_emr_cluster" "emr_zookeeper" {
  product_id = 37
  vpc_settings = {
    vpc_id    = var.vpc_id
    subnet_id = var.subnet_id
  }
  softwares = [
    "zookeeper-3.6.3",
  ]
  support_ha    = 1
  instance_name = "emr-test-demo"
  resource_spec {
    common_resource_spec {
      mem_size     = 4096
      cpu          = 2
      disk_size    = 100
      disk_type    = "CLOUD_SSD"
      spec         = "CVM.${data.tencentcloud_instance_types.cvm2c4m.instance_types.0.family}"
      storage_type = 4
      root_size    = 50
    }
    common_count = 3
  }
  login_settings = {
    password = "Tencent@cloud123"
  }
  time_span = 3600
  time_unit = "s"
  pay_mode  = 0
  placement_info {
    zone       = "ap-guangzhou-3"
    project_id = 0
  }
  sg_id = var.sg_id
  need_master_wan = "NOT_NEED_MASTER_WAN"
  scene_name = "Hadoop-Zookeeper"
}
`

const testEmrMultiZone = `
data "tencentcloud_instance_types" "cvm4c8m" {
  exclude_sold_out = true
  cpu_core_count   = 4
  memory_size      = 8
  filter {
    name   = "instance-charge-type"
    values = ["POSTPAID_BY_HOUR"]
  }
  filter {
    name   = "zone"
    values = ["ap-guangzhou-6", "ap-guangzhou-7"]
  }
  filter {
    name   = "instance-family"
    values = ["S5"]
  }
}

resource "tencentcloud_emr_cluster" "emrrrr" {
  need_master_wan = "NOT_NEED_MASTER_WAN"
  product_id      = 38
  vpc_settings = {
    vpc_id    = "vpc-axrsmmrv"
    subnet_id = "subnet-kxaxknmg"
  }
  softwares = [
    "hdfs-2.8.5",
    "knox-1.6.1",
    "openldap-2.4.44",
    "yarn-2.8.5",
    "zookeeper-3.6.3",
  ]
  support_ha    = 1
  instance_name = "emr-test-demo"
  multi_zone    = true
  multi_zone_setting {
    vpc_settings = {
      vpc_id    = "vpc-axrsmmrv"
      subnet_id = "subnet-kxaxknmg"
    }
    placement {
      zone = "ap-guangzhou-6"
    }

    resource_spec {
      master_resource_spec {
        mem_size     = 8192
        cpu          = 4
        disk_size    = 100
        disk_type    = "CLOUD_PREMIUM"
        spec         = "CVM.${data.tencentcloud_instance_types.cvm4c8m.instance_types.0.family}"
        storage_type = 5
        root_size    = 50
        multi_disks {
          disk_type = "CLOUD_PREMIUM"
          volume    = 200
          count     = 1
        }
      }
      core_resource_spec {
        mem_size     = 8192
        cpu          = 4
        disk_size    = 100
        disk_type    = "CLOUD_PREMIUM"
        spec         = "CVM.${data.tencentcloud_instance_types.cvm4c8m.instance_types.0.family}"
        storage_type = 5
        root_size    = 50
        multi_disks {
          disk_type = "CLOUD_PREMIUM"
          volume    = 100
          count     = 2
        }
      }
      common_resource_spec {
        mem_size     = 4096
        cpu          = 2
        disk_size    = 100
        disk_type    = "CLOUD_SSD"
        spec         = "CVM.${data.tencentcloud_instance_types.cvm4c8m.instance_types.0.family}"
        storage_type = 4
        root_size    = 50
      }
      common_count = 2
      master_count = 1
      core_count   = 3
    }
  }
  multi_zone_setting {
    vpc_settings = {
      vpc_id    = "vpc-axrsmmrv"
      subnet_id = "subnet-861wd75e"
    }
    placement {
      zone = "ap-guangzhou-7"
    }

    resource_spec {
      master_resource_spec {
        mem_size     = 8192
        cpu          = 4
        disk_size    = 100
        disk_type    = "CLOUD_PREMIUM"
        spec         = "CVM.${data.tencentcloud_instance_types.cvm4c8m.instance_types.0.family}"
        storage_type = 5
        root_size    = 50
        multi_disks {
          disk_type = "CLOUD_PREMIUM"
          volume    = 200
          count     = 1
        }
      }
      core_resource_spec {
        mem_size     = 8192
        cpu          = 4
        disk_size    = 100
        disk_type    = "CLOUD_PREMIUM"
        spec         = "CVM.${data.tencentcloud_instance_types.cvm4c8m.instance_types.0.family}"
        storage_type = 5
        root_size    = 50
        multi_disks {
          disk_type = "CLOUD_PREMIUM"
          volume    = 100
          count     = 2
        }
      }
      common_resource_spec {
        mem_size     = 4096
        cpu          = 2
        disk_size    = 100
        disk_type    = "CLOUD_SSD"
        spec         = "CVM.${data.tencentcloud_instance_types.cvm4c8m.instance_types.0.family}"
        storage_type = 4
        root_size    = 50
      }
      common_count = 1
      master_count = 1
      core_count   = 4
    }
  }
  login_settings = {
    password = "Tencent@cloud123"
  }
  time_span = 3600
  time_unit = "s"
  pay_mode  = 0
  sg_id = "sg-bzbu5ezt"
  placement_info {
    zone       = "ap-guangzhou-6"
    project_id = 0
  }
}
`
