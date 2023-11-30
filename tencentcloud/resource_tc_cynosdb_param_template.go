package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cynosdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cynosdb/v20190107"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudCynosdbParamTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCynosdbParamTemplateCreate,
		Read:   resourceTencentCloudCynosdbParamTemplateRead,
		Update: resourceTencentCloudCynosdbParamTemplateUpdate,
		Delete: resourceTencentCloudCynosdbParamTemplateDelete,

		Schema: map[string]*schema.Schema{
			"template_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Template Name.",
			},

			"engine_version": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "MySQL version number.",
			},

			"template_description": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Template Description.",
			},

			"template_id": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Optional parameter, template ID to be copied.",
			},

			"db_mode": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Database type, optional values: NORMAL (default), SERVERLESS.",
			},

			"param_list": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeSet,
				Description: "parameter list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"param_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Parameter Name.",
						},
						"current_value": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Current value.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudCynosdbParamTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_param_template.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request  = cynosdb.NewCreateParamTemplateRequest()
		response = cynosdb.NewCreateParamTemplateResponse()
	)
	if v, ok := d.GetOk("template_name"); ok {
		request.TemplateName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("engine_version"); ok {
		request.EngineVersion = helper.String(v.(string))
	}

	if v, ok := d.GetOk("template_description"); ok {
		request.TemplateDescription = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("template_id"); ok {
		request.TemplateId = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("db_mode"); ok {
		request.DbMode = helper.String(v.(string))
	}

	if v, ok := d.GetOk("param_list"); ok {
		for _, item := range v.(*schema.Set).List() {
			dMap := item.(map[string]interface{})
			paramItem := cynosdb.ParamItem{}
			if v, ok := dMap["param_name"]; ok {
				paramItem.ParamName = helper.String(v.(string))
			}
			if v, ok := dMap["current_value"]; ok {
				paramItem.CurrentValue = helper.String(v.(string))
			}
			request.ParamList = append(request.ParamList, &paramItem)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCynosdbClient().CreateParamTemplate(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create cynosdb paramTemplate failed, reason:%+v", logId, err)
		return err
	}

	templateId := *response.Response.TemplateId
	d.SetId(strconv.FormatInt(templateId, 10))

	return resourceTencentCloudCynosdbParamTemplateRead(d, meta)
}

func resourceTencentCloudCynosdbParamTemplateRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_param_template.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CynosdbService{client: meta.(*TencentCloudClient).apiV3Conn}

	templateId, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return err
	}

	paramTemplate, err := service.DescribeCynosdbParamTemplateById(ctx, templateId)
	if err != nil {
		return err
	}

	if paramTemplate == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CynosdbParamTemplate` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if paramTemplate.TemplateName != nil {
		_ = d.Set("template_name", paramTemplate.TemplateName)
	}

	if paramTemplate.EngineVersion != nil {
		_ = d.Set("engine_version", paramTemplate.EngineVersion)
	}

	if paramTemplate.TemplateDescription != nil {
		_ = d.Set("template_description", paramTemplate.TemplateDescription)
	}

	if paramTemplate.TemplateId != nil {
		_ = d.Set("template_id", paramTemplate.TemplateId)
	}

	if paramTemplate.DbMode != nil {
		_ = d.Set("db_mode", paramTemplate.DbMode)
	}

	params := make([]string, 0)
	if v, ok := d.GetOk("param_list"); ok {
		for _, item := range v.(*schema.Set).List() {
			if item != nil {
				dMap := item.(map[string]interface{})
				if v, ok := dMap["param_name"]; ok {
					params = append(params, v.(string))
				}
			}
		}
	}

	if paramTemplate.Items != nil {
		if len(params) > 0 {
			paramInfoSetList := make([]map[string]interface{}, 0, len(params))
			for _, param := range params {
				for _, paramList := range paramTemplate.Items {
					if *paramList.ParamName == param {
						paramListMap := map[string]interface{}{}
						if paramList.ParamName != nil {
							paramListMap["param_name"] = paramList.ParamName
						}
						if paramList.CurrentValue != nil {
							paramListMap["current_value"] = paramList.CurrentValue
						}
						paramInfoSetList = append(paramInfoSetList, paramListMap)
						break
					}
				}
			}
			_ = d.Set("param_list", paramInfoSetList)
		}
	}

	return nil
}

func resourceTencentCloudCynosdbParamTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_param_template.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	request := cynosdb.NewModifyParamTemplateRequest()

	templateId, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return err
	}

	request.TemplateId = &templateId

	immutableArgs := []string{"engine_version", "template_id", "db_mode"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("template_name") {
		if v, ok := d.GetOk("template_name"); ok {
			request.TemplateName = helper.String(v.(string))
		}
	}

	if d.HasChange("template_description") {
		if v, ok := d.GetOk("template_description"); ok {
			request.TemplateDescription = helper.String(v.(string))
		}
	}

	if d.HasChange("param_list") {
		oldParam, _ := d.GetChange("param_list")
		oldItem := oldParam.(*schema.Set).List()
		oldParamItem := make(map[string]string)
		for _, v := range oldItem {
			dMap := v.(map[string]interface{})
			key := dMap["param_name"].(string)
			value := dMap["current_value"].(string)
			oldParamItem[key] = value
		}

		if v, ok := d.GetOk("param_list"); ok {
			for _, item := range v.(*schema.Set).List() {
				dMap := item.(map[string]interface{})
				paramItem := cynosdb.ModifyParamItem{}
				if v, ok := dMap["param_name"]; ok {
					paramItem.ParamName = helper.String(v.(string))
				}
				if v, ok := dMap["current_value"]; ok {
					paramItem.CurrentValue = helper.String(v.(string))
				}
				if oldParamItem[*paramItem.ParamName] != "" {
					paramItem.OldValue = helper.String(oldParamItem[*paramItem.ParamName])
				}
				request.ParamList = append(request.ParamList, &paramItem)
			}
		}
	}

	err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCynosdbClient().ModifyParamTemplate(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update cynosdb paramTemplate failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudCynosdbParamTemplateRead(d, meta)
}

func resourceTencentCloudCynosdbParamTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_param_template.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CynosdbService{client: meta.(*TencentCloudClient).apiV3Conn}
	templateId, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return err
	}

	if err := service.DeleteCynosdbParamTemplateById(ctx, templateId); err != nil {
		return err
	}

	return nil
}
