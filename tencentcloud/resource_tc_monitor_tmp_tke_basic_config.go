/*
Provides a resource to create a monitor tmp_tke_basic_config

Example Usage

```hcl
resource "tencentcloud_monitor_tmp_tke_basic_config" "tmp_tke_basic_config" {
  instance_id = ""
  cluster_type = ""
  cluster_id = ""
  service_monitors {
		name = ""
		config = ""
		template_id = ""

  }
  pod_monitors {
		name = ""
		config = ""
		template_id = ""

  }
  raw_jobs {
		name = ""
		config = ""
		template_id = ""

  }
}
```

Import

monitor tmp_tke_basic_config can be imported using the id, e.g.

```
terraform import tencentcloud_monitor_tmp_tke_basic_config.tmp_tke_basic_config tmp_tke_basic_config_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	monitor "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/monitor/v20180724"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
	"strings"
)

func resourceTencentCloudMonitorTmpTkeBasicConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMonitorTmpTkeBasicConfigCreate,
		Read:   resourceTencentCloudMonitorTmpTkeBasicConfigRead,
		Update: resourceTencentCloudMonitorTmpTkeBasicConfigUpdate,
		Delete: resourceTencentCloudMonitorTmpTkeBasicConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
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

			"service_monitors": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "Configuration of the service monitors.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Name. The naming rule is: namespace/name. If you don&amp;#39;t have any namespace, use the default namespace: kube-system, otherwise use the specified one.",
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
				Optional:    true,
				Type:        schema.TypeList,
				Description: "Configuration of the pod monitors.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Name. The naming rule is: namespace/name. If you don&amp;#39;t have any namespace, use the default namespace: kube-system, otherwise use the specified one.",
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
				Optional:    true,
				Type:        schema.TypeList,
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

func resourceTencentCloudMonitorTmpTkeBasicConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_monitor_tmp_tke_basic_config.create")()
	defer inconsistentCheck(d, meta)()

	var instanceId string
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	var clusterType string
	if v, ok := d.GetOk("cluster_type"); ok {
		clusterType = v.(string)
	}

	var clusterId string
	if v, ok := d.GetOk("cluster_id"); ok {
		clusterId = v.(string)
	}

	d.SetId(strings.Join([]string{instanceId, clusterType, clusterId}, FILED_SP))

	return resourceTencentCloudMonitorTmpTkeBasicConfigUpdate(d, meta)
}

func resourceTencentCloudMonitorTmpTkeBasicConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_monitor_tmp_tke_basic_config.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := MonitorService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	instanceId := idSplit[0]
	clusterType := idSplit[1]
	clusterId := idSplit[2]

	tmpTkeBasicConfig, err := service.DescribeMonitorTmpTkeBasicConfigById(ctx, instanceId, clusterType, clusterId)
	if err != nil {
		return err
	}

	if tmpTkeBasicConfig == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `MonitorTmpTkeBasicConfig` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if tmpTkeBasicConfig.InstanceId != nil {
		_ = d.Set("instance_id", tmpTkeBasicConfig.InstanceId)
	}

	if tmpTkeBasicConfig.ClusterType != nil {
		_ = d.Set("cluster_type", tmpTkeBasicConfig.ClusterType)
	}

	if tmpTkeBasicConfig.ClusterId != nil {
		_ = d.Set("cluster_id", tmpTkeBasicConfig.ClusterId)
	}

	if tmpTkeBasicConfig.ServiceMonitors != nil {
		serviceMonitorsList := []interface{}{}
		for _, serviceMonitors := range tmpTkeBasicConfig.ServiceMonitors {
			serviceMonitorsMap := map[string]interface{}{}

			if tmpTkeBasicConfig.ServiceMonitors.Name != nil {
				serviceMonitorsMap["name"] = tmpTkeBasicConfig.ServiceMonitors.Name
			}

			if tmpTkeBasicConfig.ServiceMonitors.Config != nil {
				serviceMonitorsMap["config"] = tmpTkeBasicConfig.ServiceMonitors.Config
			}

			if tmpTkeBasicConfig.ServiceMonitors.TemplateId != nil {
				serviceMonitorsMap["template_id"] = tmpTkeBasicConfig.ServiceMonitors.TemplateId
			}

			serviceMonitorsList = append(serviceMonitorsList, serviceMonitorsMap)
		}

		_ = d.Set("service_monitors", serviceMonitorsList)

	}

	if tmpTkeBasicConfig.PodMonitors != nil {
		podMonitorsList := []interface{}{}
		for _, podMonitors := range tmpTkeBasicConfig.PodMonitors {
			podMonitorsMap := map[string]interface{}{}

			if tmpTkeBasicConfig.PodMonitors.Name != nil {
				podMonitorsMap["name"] = tmpTkeBasicConfig.PodMonitors.Name
			}

			if tmpTkeBasicConfig.PodMonitors.Config != nil {
				podMonitorsMap["config"] = tmpTkeBasicConfig.PodMonitors.Config
			}

			if tmpTkeBasicConfig.PodMonitors.TemplateId != nil {
				podMonitorsMap["template_id"] = tmpTkeBasicConfig.PodMonitors.TemplateId
			}

			podMonitorsList = append(podMonitorsList, podMonitorsMap)
		}

		_ = d.Set("pod_monitors", podMonitorsList)

	}

	if tmpTkeBasicConfig.RawJobs != nil {
		rawJobsList := []interface{}{}
		for _, rawJobs := range tmpTkeBasicConfig.RawJobs {
			rawJobsMap := map[string]interface{}{}

			if tmpTkeBasicConfig.RawJobs.Name != nil {
				rawJobsMap["name"] = tmpTkeBasicConfig.RawJobs.Name
			}

			if tmpTkeBasicConfig.RawJobs.Config != nil {
				rawJobsMap["config"] = tmpTkeBasicConfig.RawJobs.Config
			}

			if tmpTkeBasicConfig.RawJobs.TemplateId != nil {
				rawJobsMap["template_id"] = tmpTkeBasicConfig.RawJobs.TemplateId
			}

			rawJobsList = append(rawJobsList, rawJobsMap)
		}

		_ = d.Set("raw_jobs", rawJobsList)

	}

	return nil
}

func resourceTencentCloudMonitorTmpTkeBasicConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_monitor_tmp_tke_basic_config.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := monitor.NewModifyPrometheusConfigRequest()

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	instanceId := idSplit[0]
	clusterType := idSplit[1]
	clusterId := idSplit[2]

	request.InstanceId = &instanceId
	request.ClusterType = &clusterType
	request.ClusterId = &clusterId

	immutableArgs := []string{"instance_id", "cluster_type", "cluster_id", "service_monitors", "pod_monitors", "raw_jobs"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("instance_id") {
		if v, ok := d.GetOk("instance_id"); ok {
			request.InstanceId = helper.String(v.(string))
		}
	}

	if d.HasChange("cluster_type") {
		if v, ok := d.GetOk("cluster_type"); ok {
			request.ClusterType = helper.String(v.(string))
		}
	}

	if d.HasChange("cluster_id") {
		if v, ok := d.GetOk("cluster_id"); ok {
			request.ClusterId = helper.String(v.(string))
		}
	}

	if d.HasChange("service_monitors") {
		if v, ok := d.GetOk("service_monitors"); ok {
			for _, item := range v.([]interface{}) {
				prometheusConfigItem := monitor.PrometheusConfigItem{}
				if v, ok := dMap["name"]; ok {
					prometheusConfigItem.Name = helper.String(v.(string))
				}
				if v, ok := dMap["config"]; ok {
					prometheusConfigItem.Config = helper.String(v.(string))
				}
				if v, ok := dMap["template_id"]; ok {
					prometheusConfigItem.TemplateId = helper.String(v.(string))
				}
				request.ServiceMonitors = append(request.ServiceMonitors, &prometheusConfigItem)
			}
		}
	}

	if d.HasChange("pod_monitors") {
		if v, ok := d.GetOk("pod_monitors"); ok {
			for _, item := range v.([]interface{}) {
				prometheusConfigItem := monitor.PrometheusConfigItem{}
				if v, ok := dMap["name"]; ok {
					prometheusConfigItem.Name = helper.String(v.(string))
				}
				if v, ok := dMap["config"]; ok {
					prometheusConfigItem.Config = helper.String(v.(string))
				}
				if v, ok := dMap["template_id"]; ok {
					prometheusConfigItem.TemplateId = helper.String(v.(string))
				}
				request.PodMonitors = append(request.PodMonitors, &prometheusConfigItem)
			}
		}
	}

	if d.HasChange("raw_jobs") {
		if v, ok := d.GetOk("raw_jobs"); ok {
			for _, item := range v.([]interface{}) {
				prometheusConfigItem := monitor.PrometheusConfigItem{}
				if v, ok := dMap["name"]; ok {
					prometheusConfigItem.Name = helper.String(v.(string))
				}
				if v, ok := dMap["config"]; ok {
					prometheusConfigItem.Config = helper.String(v.(string))
				}
				if v, ok := dMap["template_id"]; ok {
					prometheusConfigItem.TemplateId = helper.String(v.(string))
				}
				request.RawJobs = append(request.RawJobs, &prometheusConfigItem)
			}
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
