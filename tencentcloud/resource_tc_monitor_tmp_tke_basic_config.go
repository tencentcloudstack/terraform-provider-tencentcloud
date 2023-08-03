/*
Provides a resource to create a monitor tmp_tke_basic_config

Example Usage

```hcl
variable "default_instance_type" {
  default = "SA1.MEDIUM2"
}

variable "availability_zone_first" {
  default = "ap-guangzhou-3"
}

variable "availability_zone_second" {
  default = "ap-guangzhou-4"
}

variable "example_cluster_cidr" {
  default = "10.31.0.0/16"
}

locals {
  first_vpc_id     = data.tencentcloud_vpc_subnets.vpc_one.instance_list.0.vpc_id
  first_subnet_id  = data.tencentcloud_vpc_subnets.vpc_one.instance_list.0.subnet_id
  second_vpc_id    = data.tencentcloud_vpc_subnets.vpc_two.instance_list.0.vpc_id
  second_subnet_id = data.tencentcloud_vpc_subnets.vpc_two.instance_list.0.subnet_id
  sg_id            = tencentcloud_security_group.sg.id
  image_id         = data.tencentcloud_images.default.image_id
}

data "tencentcloud_vpc_subnets" "vpc_one" {
  is_default        = true
  availability_zone = var.availability_zone_first
}

data "tencentcloud_vpc_subnets" "vpc_two" {
  is_default        = true
  availability_zone = var.availability_zone_second
}

resource "tencentcloud_security_group" "sg" {
  name = "tf-example-sg"
}

resource "tencentcloud_security_group_lite_rule" "sg_rule" {
  security_group_id = tencentcloud_security_group.sg.id

  ingress = [
    "ACCEPT#10.0.0.0/16#ALL#ALL",
    "ACCEPT#172.16.0.0/22#ALL#ALL",
    "DROP#0.0.0.0/0#ALL#ALL",
  ]

  egress = [
    "ACCEPT#172.16.0.0/22#ALL#ALL",
  ]
}

data "tencentcloud_images" "default" {
  image_type       = ["PUBLIC_IMAGE"]
  image_name_regex = "Final"
}

resource "tencentcloud_kubernetes_cluster" "example" {
  vpc_id                          = local.first_vpc_id
  cluster_cidr                    = var.example_cluster_cidr
  cluster_max_pod_num             = 32
  cluster_name                    = "tf_example_cluster"
  cluster_desc                    = "example for tke cluster"
  cluster_max_service_num         = 32
  cluster_internet                = false
  cluster_internet_security_group = local.sg_id
  cluster_version                 = "1.22.5"
  cluster_deploy_type             = "MANAGED_CLUSTER"

  worker_config {
    count                      = 1
    availability_zone          = var.availability_zone_first
    instance_type              = var.default_instance_type
    system_disk_type           = "CLOUD_SSD"
    system_disk_size           = 60
    internet_charge_type       = "TRAFFIC_POSTPAID_BY_HOUR"
    internet_max_bandwidth_out = 100
    public_ip_assigned         = true
    subnet_id                  = local.first_subnet_id
    img_id                     = local.image_id

    data_disk {
      disk_type = "CLOUD_PREMIUM"
      disk_size = 50
    }

    enhanced_security_service = false
    enhanced_monitor_service  = false
    user_data                 = "dGVzdA=="
    # key_ids                   = ["skey-11112222"]
    password = "ZZXXccvv1212" // Optional, should be set if key_ids not set.
  }

  worker_config {
    count                      = 1
    availability_zone          = var.availability_zone_second
    instance_type              = var.default_instance_type
    system_disk_type           = "CLOUD_SSD"
    system_disk_size           = 60
    internet_charge_type       = "TRAFFIC_POSTPAID_BY_HOUR"
    internet_max_bandwidth_out = 100
    public_ip_assigned         = true
    subnet_id                  = local.second_subnet_id

    data_disk {
      disk_type = "CLOUD_PREMIUM"
      disk_size = 50
    }

    enhanced_security_service = false
    enhanced_monitor_service  = false
    user_data                 = "dGVzdA=="
    # key_ids                   = ["skey-11112222"]
    cam_role_name = "CVM_QcsRole"
    password      = "ZZXXccvv1212" // Optional, should be set if key_ids not set.
  }

  labels = {
    "test1" = "test1",
    "test2" = "test2",
  }
}


# create monitor
variable "zone" {
  default = "ap-guangzhou"
}

variable "cluster_type" {
  default = "tke"
}

resource "tencentcloud_monitor_tmp_instance" "foo" {
  instance_name       = "tf-tmp-instance"
  vpc_id              = local.first_vpc_id
  subnet_id           = local.first_subnet_id
  data_retention_time = 30
  zone                = var.availability_zone_second
  tags = {
    "createdBy" = "terraform"
  }
}

# tmp tke bind
resource "tencentcloud_monitor_tmp_tke_cluster_agent" "foo" {
  instance_id = tencentcloud_monitor_tmp_instance.foo.id

  agents {
    region          = var.zone
    cluster_type    = var.cluster_type
    cluster_id      = tencentcloud_kubernetes_cluster.example.id
    enable_external = false
  }
}

resource "tencentcloud_monitor_tmp_tke_basic_config" "tmp_tke_basic_config" {
  instance_id  = tencentcloud_monitor_tmp_instance.foo.id
  cluster_type = var.cluster_type
  cluster_id   = tencentcloud_kubernetes_cluster.example.id
  name = "kube-system/kube-state-metrics"
  metrics_name = ["kube_job_status_succeeded"]

  depends_on = [tencentcloud_monitor_tmp_tke_cluster_agent.foo]
}

```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	monitor "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/monitor/v20180724"
)

const (
	SERVICE_MONITORS string = "service_monitors"
	POD_MONITORS     string = "pod_monitors"
	RAW_JOBS         string = "raw_jobs"
)

func resourceTencentCloudMonitorTmpTkeBasicConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMonitorTmpTkeBasicConfigCreate,
		Read:   resourceTencentCloudMonitorTmpTkeBasicConfigRead,
		Update: resourceTencentCloudMonitorTmpTkeBasicConfigUpdate,
		Delete: resourceTencentCloudMonitorTmpTkeBasicConfigDelete,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "ID of instance.",
			},

			"cluster_type": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Type of cluster.",
			},

			"cluster_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "ID of cluster.",
			},

			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name. The naming rule is: namespace/name. If you don&#39;t have any namespace, use the default namespace: kube-system, otherwise use the specified one.",
			},

			"metrics_name": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validateNotEmpty,
				},
				Description: "Configure the name of the metric to keep on.",
			},

			"config_type": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "config type, `service_monitors`, `pod_monitors`, `raw_jobs`.",
			},

			"config": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Full configuration in yaml format.",
			},
		},
	}
}

func resourceTencentCloudMonitorTmpTkeBasicConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_monitor_tmp_tke_basic_config.create")()
	defer inconsistentCheck(d, meta)()

	var (
		instanceId  string
		clusterType string
		clusterId   string
		name        string
	)
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}
	if v, ok := d.GetOk("cluster_id"); ok {
		clusterId = v.(string)
	}
	if v, ok := d.GetOk("cluster_type"); ok {
		clusterType = v.(string)
	}
	if v, ok := d.GetOk("name"); ok {
		name = v.(string)
	}

	d.SetId(strings.Join([]string{instanceId, clusterType, clusterId, name}, FILED_SP))

	return resourceTencentCloudMonitorTmpTkeBasicConfigUpdate(d, meta)
}

func resourceTencentCloudMonitorTmpTkeBasicConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_monitor_tmp_tke_basic_config.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 4 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	instanceId := idSplit[0]
	clusterType := idSplit[1]
	clusterId := idSplit[2]
	name := idSplit[3]

	service := MonitorService{client: meta.(*TencentCloudClient).apiV3Conn}
	tmpTkeBasicConfig, err := service.DescribeTkeTmpBasicConfigById(ctx, clusterId, clusterType, instanceId)
	if err != nil {
		return err
	}

	if tmpTkeBasicConfig == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `MonitorTmpTkeBasicConfig` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("instance_id", instanceId)
	_ = d.Set("cluster_type", clusterType)
	_ = d.Set("cluster_id", clusterId)
	_ = d.Set("name", name)

	configType, config, err := service.GetConfigType(name, tmpTkeBasicConfig)
	if err != nil {
		return err
	}
	_ = d.Set("config_type", configType)
	_ = d.Set("config", config.Config)

	return nil
}

func resourceTencentCloudMonitorTmpTkeBasicConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_monitor_tmp_tke_basic_config.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	request := monitor.NewModifyPrometheusConfigRequest()

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 4 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	instanceId := idSplit[0]
	clusterType := idSplit[1]
	clusterId := idSplit[2]
	name := idSplit[3]

	request.InstanceId = &instanceId
	request.ClusterType = &clusterType
	request.ClusterId = &clusterId

	if v, ok := d.GetOk("metrics_name"); ok {
		regexs := []string{}
		regexSet := v.(*schema.Set).List()
		for i := range regexSet {
			regex := regexSet[i].(string)
			regexs = append(regexs, regex)
		}

		service := MonitorService{client: meta.(*TencentCloudClient).apiV3Conn}
		tmpTkeBasicConfig, err := service.DescribeTkeTmpBasicConfigById(ctx, clusterId, clusterType, instanceId)
		if err != nil {
			return err
		}
		configType, config, err := service.GetConfigType(name, tmpTkeBasicConfig)
		if err != nil {
			return err
		}

		serviceMonitors, podMonitors, rawMobs, err := configInit(configType, config, regexs)
		if err != nil {
			return err
		}

		if serviceMonitors != "" {
			prometheusConfig := []*monitor.PrometheusConfigItem{}
			prometheusConfig = append(prometheusConfig, &monitor.PrometheusConfigItem{
				Name:   &name,
				Config: &serviceMonitors,
			})
			request.ServiceMonitors = prometheusConfig
		}
		if podMonitors != "" {
			prometheusConfig := []*monitor.PrometheusConfigItem{}
			prometheusConfig = append(prometheusConfig, &monitor.PrometheusConfigItem{
				Name:   &name,
				Config: &podMonitors,
			})
			request.PodMonitors = prometheusConfig
		}
		if rawMobs != "" {
			prometheusConfig := []*monitor.PrometheusConfigItem{}
			prometheusConfig = append(prometheusConfig, &monitor.PrometheusConfigItem{
				Name:   &name,
				Config: &rawMobs,
			})
			request.RawJobs = prometheusConfig
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseMonitorClient().ModifyPrometheusConfig(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update monitor tmpTkeBasicConfig failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudMonitorTmpTkeBasicConfigRead(d, meta)
}

func resourceTencentCloudMonitorTmpTkeBasicConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_monitor_tmp_tke_basic_config.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func configInit(configType string, respParams *monitor.PrometheusConfigItem, regexs []string) (serviceMonitorConfig, podMonitorConfig, rawMobConfig string, errRet error) {
	config := PrometheusConfig{
		Config: respParams.Config,
		Regex:  regexs,
	}
	switch configType {
	case SERVICE_MONITORS:
		serviceMonitor, err := config.UnmarshalToMap()
		if err != nil {
			errRet = err
			return
		}
		spec := serviceMonitor["spec"].(map[interface{}]interface{})["endpoints"].([]interface{})
		serviceMonitors, err := config.SetRegex(spec)
		serviceMonitor["spec"].(map[interface{}]interface{})["endpoints"] = serviceMonitors
		if err != nil {
			errRet = err
			return
		}
		serviceMonitorConfig, errRet = config.MarshalToYaml(&serviceMonitor)
		return
	case POD_MONITORS:
		serviceMonitor, err := config.UnmarshalToMap()
		if err != nil {
			errRet = err
			return
		}
		spec := serviceMonitor["spec"].(map[interface{}]interface{})["podMetricsEndpoints"].([]interface{})
		serviceMonitors, err := config.SetRegex(spec)
		serviceMonitor["spec"].(map[interface{}]interface{})["podMetricsEndpoints"] = serviceMonitors
		if err != nil {
			errRet = err
			return
		}
		podMonitorConfig, errRet = config.MarshalToYaml(&serviceMonitor)
		return
	case RAW_JOBS:
		rawMob, err := config.UnmarshalToMap()
		if err != nil {
			errRet = err
			return
		}
		configs := rawMob["scrape_configs"].([]interface{})
		rawMobConfigs, err := config.SetRegex(configs)
		rawMob["scrape_configs"] = rawMobConfigs
		if err != nil {
			errRet = err
			return
		}
		rawMobConfig, errRet = config.MarshalToYaml(&rawMob)
		return
	}
	return
}
