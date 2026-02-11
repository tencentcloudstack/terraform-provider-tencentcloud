package postgresql

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	postgresql "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/postgres/v20170312"
	postgresv20170312 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/postgres/v20170312"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudPostgresqlParameterTemplateConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudPostgresqlParameterTemplateConfigCreate,
		Read:   resourceTencentCloudPostgresqlParameterTemplateConfigRead,
		Update: resourceTencentCloudPostgresqlParameterTemplateConfigUpdate,
		Delete: resourceTencentCloudPostgresqlParameterTemplateConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"template_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Specifies the parameter template ID, which uniquely identifies the parameter template and cannot be modified. it can be obtained through the api [DescribeParameterTemplates](https://www.tencentcloud.comom/document/api/409/84067?from_cn_redirect=1).",
			},

			"modify_param_entry_set": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "The set of parameters to be modified or added.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Parameter name.",
						},
						"expected_value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The new value to which the parameter will be modified. When this parameter is used as an input parameter, its value must be a string, such as `0.1` (decimal), `1000` (integer), and `replica` (enum).",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudPostgresqlParameterTemplateConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_postgresql_parameter_template_config.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var templateId string

	if v, ok := d.GetOk("template_id"); ok {
		templateId = v.(string)
	}

	d.SetId(templateId)
	return resourceTencentCloudPostgresqlParameterTemplateConfigUpdate(d, meta)
}

func resourceTencentCloudPostgresqlParameterTemplateConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_postgresql_parameter_template_config.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service    = PostgresqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		templateId = d.Id()
	)

	respData, err := service.DescribePostgresqlParameterTemplateConfigById(ctx, templateId)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_postgresql_parameter_template_config` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	if respData.TemplateId != nil {
		_ = d.Set("template_id", respData.TemplateId)
	}

	if respData.ParamInfoSet != nil {
		paramInfoSetList := make([]map[string]interface{}, 0, len(respData.ParamInfoSet))
		for _, paramInfoSet := range respData.ParamInfoSet {
			if paramInfoSet != nil && paramInfoSet.Name != nil && paramInfoSet.CurrentValue != nil {
				paramInfoSetList = append(paramInfoSetList, map[string]interface{}{
					"name":           *paramInfoSet.Name,
					"expected_value": *paramInfoSet.CurrentValue,
				})
			}
		}

		_ = d.Set("modify_param_entry_set", paramInfoSetList)
	}

	return nil
}

func resourceTencentCloudPostgresqlParameterTemplateConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_postgresql_parameter_template_config.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		templateId = d.Id()
	)

	if d.HasChange("modify_param_entry_set") {
		oldInterface, newInterface := d.GetChange("modify_param_entry_set")
		olds := oldInterface.(*schema.Set)
		news := newInterface.(*schema.Set)
		remove := olds.Difference(news).List()
		add := news.Difference(olds).List()
		if len(remove) > 0 {
			request := postgresv20170312.NewModifyParameterTemplateRequest()
			for _, item := range remove {
				dMap := item.(map[string]interface{})
				var name string
				if v, ok := dMap["name"]; ok {
					name = v.(string)
				}

				request.DeleteParamSet = append(request.DeleteParamSet, &name)
			}

			request.TemplateId = &templateId
			reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UsePostgresqlClient().ModifyParameterTemplate(request)
				if e != nil {
					return tccommon.RetryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
				}

				return nil
			})

			if reqErr != nil {
				log.Printf("[CRITAL]%s update postgresql parameter template failed, reason:%+v", logId, reqErr)
				return reqErr
			}
		}

		if len(add) > 0 {
			request := postgresv20170312.NewModifyParameterTemplateRequest()
			for _, item := range add {
				dMap := item.(map[string]interface{})
				paramEntry := postgresql.ParamEntry{}
				if v, ok := dMap["name"]; ok {
					paramEntry.Name = helper.String(v.(string))
				}

				if v, ok := dMap["expected_value"]; ok {
					paramEntry.ExpectedValue = helper.String(v.(string))
				}

				request.ModifyParamEntrySet = append(request.ModifyParamEntrySet, &paramEntry)
			}

			request.TemplateId = &templateId
			reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UsePostgresqlClient().ModifyParameterTemplate(request)
				if e != nil {
					return tccommon.RetryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
				}

				return nil
			})

			if reqErr != nil {
				log.Printf("[CRITAL]%s update postgresql parameter template failed, reason:%+v", logId, reqErr)
				return reqErr
			}
		}
	}

	return resourceTencentCloudPostgresqlParameterTemplateConfigRead(d, meta)
}

func resourceTencentCloudPostgresqlParameterTemplateConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_postgresql_parameter_template_config.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
