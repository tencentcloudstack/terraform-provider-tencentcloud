package tmp

import (
	"context"
	"fmt"
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svcmonitor "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/monitor"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	monitor "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/monitor/v20180724"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudMonitorTmpTkeConfig() *schema.Resource {
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
	defer tccommon.LogElapsed("resource.tencentcloud_tke_tmp_config.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId    = tccommon.GetLogId(tccommon.ContextNil)
		ctx      = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service  = svcmonitor.NewMonitorService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
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
	defer tccommon.LogElapsed("resource.tencentcloud_tke_tmp_config.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
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

	ids := strings.Join([]string{*request.InstanceId, *request.ClusterType, *request.ClusterId}, tccommon.FILED_SP)

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		response, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMonitorClient().CreatePrometheusConfig(request)
		if e != nil {
			return tccommon.RetryError(e)
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
	defer tccommon.LogElapsed("resource.tencentcloud_tke_tmp_config.update, Id: %s", d.Id())()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		request = monitor.NewModifyPrometheusConfigRequest()
		client  = meta.(tccommon.ProviderMeta).GetAPIV3Conn()
		service = svcmonitor.NewMonitorService(client)
	)

	ids, err := service.ParseConfigId(d.Id())
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

	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		response, e := client.UseMonitorClient().ModifyPrometheusConfig(request)
		if e != nil {
			return tccommon.RetryError(e)
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
	defer tccommon.LogElapsed("resource.tencentcloud_tke_tmp_config.delete, Id: %s", d.Id())()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId           = tccommon.GetLogId(tccommon.ContextNil)
		ctx             = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service         = svcmonitor.NewMonitorService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
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
