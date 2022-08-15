/*
Provides a resource to create a tke tmpPrometheusConfig

Example Usage

```hcl

resource "tencentcloud_monitor_tmp_tke_config" "foo" {
  instance_id  = "xxx"
  cluster_type = "xxx"
  cluster_id   = "xxx"

  raw_jobs {
    name   = "rawjob_001"
    config = "your_config_for_raw_jobs\n"
  }

  service_monitors {
    name   = "servicemonitors_001"
    config = "your_config_for_service_monitors\n"
  }

  pod_monitors {
    name   = "pod_monitors_001"
    config = "your_config_for_pod_monitors\n"
  }
}
*/
package tencentcloud

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	tke "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tke/v20180525"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudMonitorTmpTkeConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTkeTmpConfigCreate,
		Read:   resourceTencentCloudTkeTmpConfigRead,
		Update: resourceTencentCloudTkeTmpConfigUpdate,
		Delete: resourceTencentCloudTkeTmpConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of instance.",
			},
			"cluster_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Type of cluster.",
			},
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
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
			"pod_monitors": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Configuration of the pod monitors.",
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
		service  = TkeService{client: meta.(*TencentCloudClient).apiV3Conn}
		configId = d.Id()
	)

	params, err := service.DescribeTkeTmpConfigById(logId, configId)

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
	if e := d.Set("service_monitors", flattenPrometheusConfigItems(params.ServiceMonitors)); e != nil {
		log.Printf("[CRITAL]%s provider set service_monitors fail, reason:%s\n", logId, e.Error())
		return e
	}
	if e := d.Set("pod_monitors", flattenPrometheusConfigItems(params.PodMonitors)); e != nil {
		log.Printf("[CRITAL]%s provider set pod_monitors fail, reason:%s\n", logId, e.Error())
		return e
	}
	if e := d.Set("raw_jobs", flattenPrometheusConfigItems(params.RawJobs)); e != nil {
		log.Printf("[CRITAL]%s provider set raw_jobs fail, reason:%s\n", logId, e.Error())
		return e
	}
	return nil
}

func resourceTencentCloudTkeTmpConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tke_tmp_config.create")()
	defer inconsistentCheck(d, meta)()

	var (
		logId   = getLogId(contextNil)
		request = tke.NewCreatePrometheusConfigRequest()
		client  = meta.(*TencentCloudClient).apiV3Conn.UseTkeClient()
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
		response, e := client.CreatePrometheusConfig(request)
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
		request = tke.NewModifyPrometheusConfigRequest()
		client  = meta.(*TencentCloudClient).apiV3Conn.UseTkeClient()
	)

	if d.HasChange("instance_id") {
		return fmt.Errorf("`instance_id` do not support change now.")
	}
	if d.HasChange("cluster_id") {
		return fmt.Errorf("`cluster_id` do not support change now.")
	}
	if d.HasChange("cluster_type") {
		return fmt.Errorf("`cluster_type` do not support change now.")
	}

	ids, err := parseId(d.Id())
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
		response, e := client.ModifyPrometheusConfig(request)
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

	service := TkeService{client: meta.(*TencentCloudClient).apiV3Conn}

	if err := service.DeleteTkeTmpConfigById(d); err != nil {
		return err
	}

	return nil
}

func flattenPrometheusConfigItems(objList []*tke.PrometheusConfigItem) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(objList))
	for i := range objList {
		v := objList[i]
		item := map[string]interface{}{
			"config": v.Config,
			"name":   v.Name,
		}
		if v.TemplateId != nil {
			item["template_id"] = v.TemplateId
		}
		result = append(result, item)
	}
	return result
}
