/*
Provides a resource to create a tmp tke template attachment

Example Usage

```hcl
# create tke
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

# create monitor template
resource "tencentcloud_monitor_tmp_tke_template" "foo" {
  template {
    name     = "tf-template"
    level    = "cluster"
    describe = "template"
    service_monitors {
      name   = "tf-ServiceMonitor"
      config = <<-EOT
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  name: example-service-monitor
  namespace: monitoring
  labels:
    k8s-app: example-service
spec:
  selector:
    matchLabels:
      k8s-app: example-service
  namespaceSelector:
    matchNames:
      - default
  endpoints:
  - port: http-metrics
    interval: 30s
    path: /metrics
    scheme: http
    bearerTokenFile: /var/run/secrets/kubernetes.io/serviceaccount/token
    tlsConfig:
      insecureSkipVerify: true
      EOT
    }

    pod_monitors {
      name   = "tf-PodMonitors"
      config = <<-EOT
apiVersion: monitoring.coreos.com/v1
kind: PodMonitor
metadata:
  name: example-pod-monitor
  namespace: monitoring
  labels:
    k8s-app: example-pod
spec:
  selector:
    matchLabels:
      k8s-app: example-pod
  namespaceSelector:
    matchNames:
      - default
  podMetricsEndpoints:
  - port: http-metrics
    interval: 30s
    path: /metrics
    scheme: http
    bearerTokenFile: /var/run/secrets/kubernetes.io/serviceaccount/token
    tlsConfig:
      insecureSkipVerify: true
EOT
    }

    pod_monitors {
      name   = "tf-RawJobs"
      config = <<-EOT
scrape_configs:
  - job_name: 'example-job'
    scrape_interval: 30s
    static_configs:
      - targets: ['example-service.default.svc.cluster.local:8080']
    metrics_path: /metrics
    scheme: http
    bearer_token_file: /var/run/secrets/kubernetes.io/serviceaccount/token
    tls_config:
      insecure_skip_verify: true
EOT
    }
  }
}

resource "tencentcloud_monitor_tmp_tke_template_attachment" "temp_attachment" {
  template_id = tencentcloud_monitor_tmp_tke_template.foo.id

  targets {
    cluster_type = var.cluster_type
    cluster_id   = tencentcloud_kubernetes_cluster.example.id
    region       = var.zone
    instance_id  = tencentcloud_monitor_tmp_instance.foo.id
  }

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
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudMonitorTmpTkeTemplateAttachment() *schema.Resource {
	return &schema.Resource{
		Read:   resourceTencentCloudMonitorTmpTkeTemplateAttachmentRead,
		Create: resourceTencentCloudMonitorTmpTkeTemplateAttachmentCreate,
		Delete: resourceTencentCloudMonitorTmpTkeTemplateAttachmentDelete,
		Schema: map[string]*schema.Schema{
			"template_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the template, which is used for the outgoing reference.",
			},

			"targets": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Required:    true,
				ForceNew:    true,
				Description: "Sync target details.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"region": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "target area.",
						},
						"instance_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "instance id.",
						},
						"cluster_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "ID of the cluster.",
						},
						"sync_time": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Last sync template time.",
						},
						"version": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Template version currently in use.",
						},
						"cluster_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Cluster type.",
						},
						"instance_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Name of the prometheus instance.",
						},
						"cluster_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Name the cluster.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudMonitorTmpTkeTemplateAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_monitor_tmp_tke_template_attachment.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := monitor.NewSyncPrometheusTempRequest()

	if v, ok := d.GetOk("template_id"); ok {
		request.TemplateId = helper.String(v.(string))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "targets"); ok {
		var prometheusTarget monitor.PrometheusTemplateSyncTarget
		if v, ok := dMap["region"]; ok {
			prometheusTarget.Region = helper.String(v.(string))
		}

		if v, ok := dMap["instance_id"]; ok {
			prometheusTarget.InstanceId = helper.String(v.(string))
		}

		if v, ok := dMap["cluster_id"]; ok {
			prometheusTarget.ClusterId = helper.String(v.(string))
		}

		if v, ok := dMap["sync_time"]; ok {
			prometheusTarget.SyncTime = helper.String(v.(string))
		}

		if v, ok := dMap["version"]; ok {
			prometheusTarget.Version = helper.String(v.(string))
		}

		if v, ok := dMap["cluster_type"]; ok {
			prometheusTarget.ClusterType = helper.String(v.(string))
		}

		if v, ok := dMap["instance_name"]; ok {
			prometheusTarget.InstanceName = helper.String(v.(string))
		}

		if v, ok := dMap["cluster_name"]; ok {
			prometheusTarget.ClusterName = helper.String(v.(string))
		}

		prometheusTargets := make([]*monitor.PrometheusTemplateSyncTarget, 0)
		prometheusTargets = append(prometheusTargets, &prometheusTarget)
		request.Targets = prometheusTargets

	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseMonitorClient().SyncPrometheusTemp(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s sync tke template failed, reason:%+v", logId, err)
		return err
	}

	templateId := *request.TemplateId
	instanceId := *request.Targets[0].InstanceId
	region := *request.Targets[0].Region
	d.SetId(strings.Join([]string{templateId, instanceId, region}, FILED_SP))

	return resourceTencentCloudMonitorTmpTkeTemplateAttachmentRead(d, meta)
}

func resourceTencentCloudMonitorTmpTkeTemplateAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_monitor_tmp_tke_template_attachment.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := MonitorService{client: meta.(*TencentCloudClient).apiV3Conn}

	ids := strings.Split(d.Id(), FILED_SP)
	if len(ids) != 3 {
		return fmt.Errorf("id is broken, id is %s", d.Id())
	}

	templateId := ids[0]
	instanceId := ids[1]
	region := ids[2]

	targets, err := service.DescribePrometheusTempSync(ctx, templateId)

	if err != nil {
		return err
	}

	if targets == nil || len(targets) < 1 {
		d.SetId("")
		return fmt.Errorf("resource `targets` %s does not exist", templateId)
	}

	tempTargets := make([]map[string]interface{}, 0)
	for _, v := range targets {
		if *v.InstanceId == instanceId && *v.Region == region {
			tempTargets = append(tempTargets, map[string]interface{}{
				"region":      v.Region,
				"instance_id": v.InstanceId,
				//"cluster_id":    v.ClusterId,
				//"sync_time":     v.SyncTime,
				//"version":       v.Version,
				//"cluster_type":  v.ClusterType,
				//"instance_name": v.InstanceName,
				//"cluster_name":  v.ClusterName,
			})
		}
	}
	_ = d.Set("targets", tempTargets)

	return nil
}

func resourceTencentCloudMonitorTmpTkeTemplateAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_monitor_tmp_tke_template_attachment.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	request := monitor.NewDeletePrometheusTempSyncRequest()

	ids := strings.Split(d.Id(), FILED_SP)
	if len(ids) != 3 {
		return fmt.Errorf("id is broken, id is %s", d.Id())
	}

	templateId := ids[0]
	instanceId := ids[1]
	region := ids[2]

	request.TemplateId = &templateId
	var targets []*monitor.PrometheusTemplateSyncTarget
	target := monitor.PrometheusTemplateSyncTarget{
		Region:     &region,
		InstanceId: &instanceId,
	}
	targets = append(targets, &target)
	request.Targets = targets

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseMonitorClient().DeletePrometheusTempSync(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
