/*
Provides a resource to create a tmp tke template

Example Usage

```hcl
resource "tencentcloud_monitor_tmp_tke_template" "template" {
  template {
    name = "test"
    level = "cluster"
    describe = "template"
    service_monitors {
      name = "test"
      config = "xxxxx"
    }
  }
}
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	tke "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tke/v20180525"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudMonitorTmpTkeTemplate() *schema.Resource {
	return &schema.Resource{
		Read:   resourceTencentCloudMonitorTmpTkeTemplateRead,
		Create: resourceTencentCloudMonitorTmpTkeTemplateCreate,
		Update: resourceTencentCloudMonitorTmpTkeTemplateUpdate,
		Delete: resourceTencentCloudMonitorTmpTkeTemplateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"template": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Required:    true,
				Description: "Template settings.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Template name.",
						},
						"level": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Template dimensions, the following types are supported `instance` instance level, `cluster` cluster level.",
						},
						"describe": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Template description.",
						},
						"record_rules": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Effective when Level is instance, A list of aggregation rules in the template.",
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
										Description: "Used for the argument, if the configuration comes to the template, the template id.",
									},
								},
							},
						},
						"service_monitors": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Effective when Level is a cluster, A list of ServiceMonitor rules in the template.",
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
										Description: "Used for the argument, if the configuration comes to the template, the template id.",
									},
								},
							},
						},
						"pod_monitors": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Effective when Level is a cluster, A list of PodMonitors rules in the template.",
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
										Description: "Used for the argument, if the configuration comes to the template, the template id.",
									},
								},
							},
						},
						"raw_jobs": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Effective when Level is a cluster, A list of RawJobs rules in the template.",
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
										Description: "Used for the argument, if the configuration comes to the template, the template id.",
									},
								},
							},
						},
						"template_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The ID of the template, which is used for the outgoing reference.",
						},
						"update_time": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Last updated, for outgoing references.",
						},
						"version": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Whether the system-supplied default template is used for outgoing references.",
						},
						"is_default": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether the system-supplied default template is used for outgoing references.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudMonitorTmpTkeTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_monitor_tmp_tke_template.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request  = tke.NewCreatePrometheusTempRequest()
		response *tke.CreatePrometheusTempResponse
	)

	if dMap, ok := helper.InterfacesHeadMap(d, "template"); ok {
		prometheusTemp := tke.PrometheusTemp{}
		if v, ok := dMap["name"]; ok {
			prometheusTemp.Name = helper.String(v.(string))
		}
		if v, ok := dMap["level"]; ok {
			prometheusTemp.Level = helper.String(v.(string))
		}
		if v, ok := dMap["describe"]; ok {
			prometheusTemp.Describe = helper.String(v.(string))
		}
		if v, ok := d.GetOk("record_rules"); ok {
			resList := v.([]interface{})
			prometheusConfigItem := make([]*tke.PrometheusConfigItem, 0, len(resList))
			for _, res := range resList {
				vv := res.(map[string]interface{})
				var item tke.PrometheusConfigItem
				if v, ok := vv["name"]; ok {
					item.Name = helper.String(v.(string))
				}
				if v, ok := vv["config"]; ok {
					item.Config = helper.String(v.(string))
				}
				if v, ok := vv["template_id"]; ok {
					item.TemplateId = helper.String(v.(string))
				}
				prometheusConfigItem = append(prometheusConfigItem, &item)
			}

			prometheusTemp.RecordRules = prometheusConfigItem
		}
		if v, ok := d.GetOk("service_monitors"); ok {
			resList := v.([]interface{})
			prometheusConfigItem := make([]*tke.PrometheusConfigItem, 0, len(resList))
			for _, res := range resList {
				vv := res.(map[string]interface{})
				var item tke.PrometheusConfigItem
				if v, ok := vv["name"]; ok {
					item.Name = helper.String(v.(string))
				}
				if v, ok := vv["config"]; ok {
					item.Config = helper.String(v.(string))
				}
				if v, ok := vv["template_id"]; ok {
					item.TemplateId = helper.String(v.(string))
				}
				prometheusConfigItem = append(prometheusConfigItem, &item)
			}
			prometheusTemp.ServiceMonitors = prometheusConfigItem
		}
		if v, ok := d.GetOk("pod_monitors"); ok {
			resList := v.([]interface{})
			prometheusConfigItem := make([]*tke.PrometheusConfigItem, 0, len(resList))
			for _, res := range resList {
				vv := res.(map[string]interface{})
				var item tke.PrometheusConfigItem
				if v, ok := vv["name"]; ok {
					item.Name = helper.String(v.(string))
				}
				if v, ok := vv["config"]; ok {
					item.Config = helper.String(v.(string))
				}
				if v, ok := vv["template_id"]; ok {
					item.TemplateId = helper.String(v.(string))
				}
				prometheusConfigItem = append(prometheusConfigItem, &item)
			}
			prometheusTemp.PodMonitors = prometheusConfigItem
		}
		if v, ok := d.GetOk("raw_jobs"); ok {
			resList := v.([]interface{})
			prometheusConfigItem := make([]*tke.PrometheusConfigItem, 0, len(resList))
			for _, res := range resList {
				vv := res.(map[string]interface{})
				var item tke.PrometheusConfigItem
				if v, ok := vv["name"]; ok {
					item.Name = helper.String(v.(string))
				}
				if v, ok := vv["config"]; ok {
					item.Config = helper.String(v.(string))
				}
				if v, ok := vv["template_id"]; ok {
					item.TemplateId = helper.String(v.(string))
				}
				prometheusConfigItem = append(prometheusConfigItem, &item)
			}
			prometheusTemp.RawJobs = prometheusConfigItem
		}
		if v, ok := dMap["template_id"]; ok {
			prometheusTemp.TemplateId = helper.String(v.(string))
		}
		if v, ok := dMap["update_time"]; ok {
			prometheusTemp.UpdateTime = helper.String(v.(string))
		}
		if v, ok := dMap["version"]; ok {
			prometheusTemp.Version = helper.String(v.(string))
		}
		if v, ok := dMap["is_default"]; ok {
			prometheusTemp.IsDefault = helper.Bool(v.(bool))
		}
		request.Template = &prometheusTemp

	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTkeClient().CreatePrometheusTemp(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create tke template failed, reason:%+v", logId, err)
		return err
	}

	templateId := *response.Response.TemplateId
	d.SetId(templateId)

	return resourceTencentCloudMonitorTmpTkeTemplateRead(d, meta)
}

func resourceTencentCloudMonitorTmpTkeTemplateRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_monitor_tmp_tke_template.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TkeService{client: meta.(*TencentCloudClient).apiV3Conn}

	templateId := d.Id()

	template, err := service.DescribeTmpTkeTemplateById(ctx, templateId)

	if err != nil {
		return err
	}

	if template == nil {
		d.SetId("")
		return fmt.Errorf("resource `template` %s does not exist", templateId)
	}

	templates := make([]map[string]interface{}, 0)
	templates = append(templates, map[string]interface{}{
		"name":       template.Name,
		"level":      template.Level,
		"describe":   template.Describe,
		"is_default": template.IsDefault,
	})
	_ = d.Set("template", templates)
	return nil
}

func resourceTencentCloudMonitorTmpTkeTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_monitor_tmp_tke_template.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := tke.NewModifyPrometheusTempRequest()

	request.TemplateId = helper.String(d.Id())

	if d.HasChange("template") {
		if dMap, ok := helper.InterfacesHeadMap(d, "template"); ok {
			prometheusTemp := tke.PrometheusTempModify{}
			if v, ok := dMap["name"]; ok {
				prometheusTemp.Name = helper.String(v.(string))
			}
			if v, ok := dMap["describe"]; ok {
				prometheusTemp.Describe = helper.String(v.(string))
			}
			if v, ok := d.GetOk("record_rules"); ok {
				resList := v.([]interface{})
				prometheusConfigItem := make([]*tke.PrometheusConfigItem, 0, len(resList))
				for _, res := range resList {
					vv := res.(map[string]interface{})
					var item tke.PrometheusConfigItem
					if v, ok := vv["name"]; ok {
						item.Name = helper.String(v.(string))
					}
					if v, ok := vv["config"]; ok {
						item.Config = helper.String(v.(string))
					}
					if v, ok := vv["template_id"]; ok {
						item.TemplateId = helper.String(v.(string))
					}
					prometheusConfigItem = append(prometheusConfigItem, &item)
				}

				prometheusTemp.RecordRules = prometheusConfigItem
			}
			if v, ok := d.GetOk("service_monitors"); ok {
				resList := v.([]interface{})
				prometheusConfigItem := make([]*tke.PrometheusConfigItem, 0, len(resList))
				for _, res := range resList {
					vv := res.(map[string]interface{})
					var item tke.PrometheusConfigItem
					if v, ok := vv["name"]; ok {
						item.Name = helper.String(v.(string))
					}
					if v, ok := vv["config"]; ok {
						item.Config = helper.String(v.(string))
					}
					if v, ok := vv["template_id"]; ok {
						item.TemplateId = helper.String(v.(string))
					}
					prometheusConfigItem = append(prometheusConfigItem, &item)
				}
				prometheusTemp.ServiceMonitors = prometheusConfigItem
			}
			if v, ok := d.GetOk("pod_monitors"); ok {
				resList := v.([]interface{})
				prometheusConfigItem := make([]*tke.PrometheusConfigItem, 0, len(resList))
				for _, res := range resList {
					vv := res.(map[string]interface{})
					var item tke.PrometheusConfigItem
					if v, ok := vv["name"]; ok {
						item.Name = helper.String(v.(string))
					}
					if v, ok := vv["config"]; ok {
						item.Config = helper.String(v.(string))
					}
					if v, ok := vv["template_id"]; ok {
						item.TemplateId = helper.String(v.(string))
					}
					prometheusConfigItem = append(prometheusConfigItem, &item)
				}
				prometheusTemp.PodMonitors = prometheusConfigItem
			}
			if v, ok := d.GetOk("raw_jobs"); ok {
				resList := v.([]interface{})
				prometheusConfigItem := make([]*tke.PrometheusConfigItem, 0, len(resList))
				for _, res := range resList {
					vv := res.(map[string]interface{})
					var item tke.PrometheusConfigItem
					if v, ok := vv["name"]; ok {
						item.Name = helper.String(v.(string))
					}
					if v, ok := vv["config"]; ok {
						item.Config = helper.String(v.(string))
					}
					if v, ok := vv["template_id"]; ok {
						item.TemplateId = helper.String(v.(string))
					}
					prometheusConfigItem = append(prometheusConfigItem, &item)
				}
				prometheusTemp.RawJobs = prometheusConfigItem
			}
			request.Template = &prometheusTemp
		}
	}
	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTkeClient().ModifyPrometheusTemp(request)
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

	return resourceTencentCloudMonitorTmpTkeTemplateRead(d, meta)
}

func resourceTencentCloudMonitorTmpTkeTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_monitor_tmp_tke_template.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TkeService{client: meta.(*TencentCloudClient).apiV3Conn}
	id := d.Id()

	if err := service.DeleteTmpTkeTemplate(ctx, id); err != nil {
		return err
	}

	return nil
}
