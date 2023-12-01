/*
Provides a resource to create a tke tmpPrometheusConfig

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

resource "tencentcloud_monitor_tmp_tke_config" "foo" {
  instance_id  = tencentcloud_monitor_tmp_instance.foo.id
  cluster_type = var.cluster_type
  cluster_id   = tencentcloud_kubernetes_cluster.example.id

  raw_jobs {
    name   = "raw_jobs_001"
    config = "your config for raw_jobs_001\n"
  }

  service_monitors {
    name   = "kube-system/service-monitor-001" # name with default namespace kube-system
    config = "apiVersion: monitoring.coreos.com/v1\nkind: ServiceMonitor\nmetadata:\n  name: service-monitor-001\n  namespace: kube-system\n"
  }

  pod_monitors {
    name   = "mynamespace/pod-monitor-001" # name with the specified namespace
    config = "apiVersion: monitoring.coreos.com/v1\nkind: PodMonitor\nmetadata:\n  name: pod-monitor-001\n  namespace: mynamespace\n"
  }
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
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudMonitorTmpTkeConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTkeTmpConfigCreate,
		Read:   resourceTencentCloudTkeTmpConfigRead,
		Update: resourceTencentCloudTkeTmpConfigUpdate,
		Delete: resourceTencentCloudTkeTmpConfigDelete,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of instance.",
			},
			"cluster_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Type of cluster.",
			},
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of cluster.",
			},
			"config": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Global configuration.",
			},
			"service_monitors": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Configuration of the service monitors.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Name. The naming rule is: namespace/name. If you don't have any namespace, use the default namespace: kube-system, otherwise use the specified one.",
						},
						"config": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Config.",
						},
						"template_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Used for output parameters, if the configuration comes from a template, it is the template id.",
						},
					},
				},
			},
			"pod_monitors": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Configuration of the pod monitors.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Name. The naming rule is: namespace/name. If you don't have any namespace, use the default namespace: kube-system, otherwise use the specified one.",
						},
						"config": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Config.",
						},
						"template_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Used for output parameters, if the configuration comes from a template, it is the template id.",
						},
					},
				},
			},
			"raw_jobs": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Configuration of the native prometheus job.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Name.",
						},
						"config": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Config.",
						},
						"template_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Used for output parameters, if the configuration comes from a template, it is the template id.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudTkeTmpConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tke_tmp_config.read")()
	defer inconsistentCheck(d, meta)()

	var (
		logId    = getLogId(contextNil)
		ctx      = context.WithValue(context.TODO(), logIdKey, logId)
		service  = MonitorService{client: meta.(*TencentCloudClient).apiV3Conn}
		configId = d.Id()
	)

	params, err := service.DescribeTkeTmpConfigById(ctx, configId)

	if err != nil {
		return err
	}

	if params == nil {
		d.SetId("")
		return fmt.Errorf("resource `prometheus_config` %s does not exist", configId)
	}

	if e := d.Set("config", params.Config); e != nil {
		log.Printf("[CRITAL]%s provider set config fail, reason:%s\n", logId, e.Error())
		return e
	}

	return nil
}

func resourceTencentCloudTkeTmpConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tke_tmp_config.create")()
	defer inconsistentCheck(d, meta)()

	var (
		logId   = getLogId(contextNil)
		request = monitor.NewCreatePrometheusConfigRequest()
	)

	if v, ok := d.GetOk("instance_id"); ok {
		request.InstanceId = helper.String(v.(string))
	}
	if v, ok := d.GetOk("cluster_id"); ok {
		request.ClusterId = helper.String(v.(string))
	}
	if v, ok := d.GetOk("cluster_type"); ok {
		request.ClusterType = helper.String(v.(string))
	}
	if v, ok := d.GetOk("service_monitors"); ok {
		request.ServiceMonitors = serializePromConfigItems(v)
	}
	if v, ok := d.GetOk("pod_monitors"); ok {
		request.PodMonitors = serializePromConfigItems(v)
	}
	if v, ok := d.GetOk("raw_jobs"); ok {
		request.RawJobs = serializePromConfigItems(v)
	}

	ids := strings.Join([]string{*request.InstanceId, *request.ClusterType, *request.ClusterId}, FILED_SP)

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		response, e := meta.(*TencentCloudClient).apiV3Conn.UseMonitorClient().CreatePrometheusConfig(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, ids [%s], request body [%s], response body [%s]\n",
				logId, request.GetAction(), ids, request.ToJsonString(), response.ToJsonString())
		}
		return nil
	})

	if err != nil {
		return err
	}

	d.SetId(ids)

	return resourceTencentCloudTkeTmpConfigRead(d, meta)
}

func resourceTencentCloudTkeTmpConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tke_tmp_config.update, Id: %s", d.Id())()
	defer inconsistentCheck(d, meta)()

	var (
		logId   = getLogId(contextNil)
		request = monitor.NewModifyPrometheusConfigRequest()
		client  = meta.(*TencentCloudClient).apiV3Conn
		service = MonitorService{client: client}
	)

	ids, err := service.parseConfigId(d.Id())
	if err != nil {
		return err
	}
	request.ClusterId = &ids.ClusterId
	request.ClusterType = &ids.ClusterType
	request.InstanceId = &ids.InstanceId

	if d.HasChange("service_monitors") {
		if v, ok := d.GetOk("service_monitors"); ok {
			request.ServiceMonitors = serializePromConfigItems(v)
		}
	}

	if d.HasChange("pod_monitors") {
		if v, ok := d.GetOk("pod_monitors"); ok {
			request.PodMonitors = serializePromConfigItems(v)
		}
	}

	if d.HasChange("raw_jobs") {
		if v, ok := d.GetOk("raw_jobs"); ok {
			request.RawJobs = serializePromConfigItems(v)
		}
	}

	err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		response, e := client.UseMonitorClient().ModifyPrometheusConfig(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, ids [%s], request body [%s], response body [%s]\n",
				logId, request.GetAction(), d.Id(), request.ToJsonString(), response.ToJsonString())
		}
		return nil
	})

	if err != nil {
		return err
	}

	return resourceTencentCloudTkeTmpConfigRead(d, meta)
}

func resourceTencentCloudTkeTmpConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tke_tmp_config.delete, Id: %s", d.Id())()
	defer inconsistentCheck(d, meta)()

	var (
		logId           = getLogId(contextNil)
		ctx             = context.WithValue(context.TODO(), logIdKey, logId)
		service         = MonitorService{client: meta.(*TencentCloudClient).apiV3Conn}
		ServiceMonitors = []*string{}
		PodMonitors     = []*string{}
		RawJobs         = []*string{}
	)

	if v, ok := d.GetOk("service_monitors"); ok {
		ServiceMonitors = serializePromConfigItemNames(v)
	}

	if v, ok := d.GetOk("pod_monitors"); ok {
		PodMonitors = serializePromConfigItemNames(v)
	}

	if v, ok := d.GetOk("raw_jobs"); ok {
		RawJobs = serializePromConfigItemNames(v)
	}

	if err := service.DeleteTkeTmpConfigByName(ctx, d.Id(), ServiceMonitors, PodMonitors, RawJobs); err != nil {
		return err
	}

	return nil
}

func serializePromConfigItems(v interface{}) []*monitor.PrometheusConfigItem {
	resList := v.([]interface{})
	items := make([]*monitor.PrometheusConfigItem, 0, len(resList))
	for _, res := range resList {
		vv := res.(map[string]interface{})
		var item monitor.PrometheusConfigItem
		if v, ok := vv["name"]; ok {
			item.Name = helper.String(v.(string))
		}
		if v, ok := vv["config"]; ok {
			item.Config = helper.String(v.(string))
		}
		if v, ok := vv["template_id"]; ok {
			item.TemplateId = helper.String(v.(string))
		}
		items = append(items, &item)
	}
	return items
}

func serializePromConfigItemNames(v interface{}) []*string {
	resList := v.([]interface{})
	names := make([]*string, 0, len(resList))
	for _, res := range resList {
		vv := res.(map[string]interface{})
		if v, ok := vv["name"]; ok {
			names = append(names, helper.String(v.(string)))
		}
	}
	return names
}
