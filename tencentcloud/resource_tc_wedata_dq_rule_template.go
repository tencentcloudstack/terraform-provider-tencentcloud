/*
Provides a resource to create a wedata dq_rule_template

Example Usage

```hcl
resource "tencentcloud_wedata_dq_rule_template" "dq_rule_template" {
  type =
  name = ""
  quality_dim =
  source_object_type =
  description = ""
  source_engine_types =
  multi_source_flag =
  sql_expression = ""
  project_id = ""
  where_flag =
}
```

Import

wedata dq_rule_template can be imported using the id, e.g.

```
terraform import tencentcloud_wedata_dq_rule_template.dq_rule_template dq_rule_template_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	wedata "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/wedata/v20210820"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudWedataDq_ruleTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudWedataDq_ruleTemplateCreate,
		Read:   resourceTencentCloudWedataDq_ruleTemplateRead,
		Update: resourceTencentCloudWedataDq_ruleTemplateUpdate,
		Delete: resourceTencentCloudWedataDq_ruleTemplateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"type": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Template Type 1. System Template 2. User-defined template.",
			},

			"name": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Template name.",
			},

			"quality_dim": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Quality detection dimension 1. Accuracy 2. Uniqueness 3. Completeness 4. Consistency 5. Timeliness 6. effectiveness.",
			},

			"source_object_type": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Source end data object type 1. Constant 2. Offline table level 2. Offline field level.",
			},

			"description": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Template description.",
			},

			"source_engine_types": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Description: "Type of the engine on the source end.",
			},

			"multi_source_flag": {
				Optional:    true,
				Type:        schema.TypeBool,
				Description: "Whether to associate other tables.",
			},

			"sql_expression": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "SQL expression.",
			},

			"project_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Project id.",
			},

			"where_flag": {
				Optional:    true,
				Type:        schema.TypeBool,
				Description: "Add &amp;amp;#39;where&amp;amp;#39; parameter or not.",
			},
		},
	}
}

func resourceTencentCloudWedataDq_ruleTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_wedata_dq_rule_template.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = wedata.NewCreateRuleTemplateRequest()
		response   = wedata.NewCreateRuleTemplateResponse()
		templateId int
	)
	if v, ok := d.GetOkExists("type"); ok {
		request.Type = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("quality_dim"); ok {
		request.QualityDim = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("source_object_type"); ok {
		request.SourceObjectType = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("description"); ok {
		request.Description = helper.String(v.(string))
	}

	if v, ok := d.GetOk("source_engine_types"); ok {
		sourceEngineTypesSet := v.(*schema.Set).List()
		for i := range sourceEngineTypesSet {
			sourceEngineTypes := sourceEngineTypesSet[i].(int)
			request.SourceEngineTypes = append(request.SourceEngineTypes, helper.IntUint64(sourceEngineTypes))
		}
	}

	if v, ok := d.GetOkExists("multi_source_flag"); ok {
		request.MultiSourceFlag = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOk("sql_expression"); ok {
		request.SqlExpression = helper.String(v.(string))
	}

	if v, ok := d.GetOk("project_id"); ok {
		request.ProjectId = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("where_flag"); ok {
		request.WhereFlag = helper.Bool(v.(bool))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseWedataClient().CreateRuleTemplate(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create wedata dq_ruleTemplate failed, reason:%+v", logId, err)
		return err
	}

	templateId = *response.Response.TemplateId
	d.SetId(helper.Int64ToStr(int64(templateId)))

	return resourceTencentCloudWedataDq_ruleTemplateRead(d, meta)
}

func resourceTencentCloudWedataDq_ruleTemplateRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_wedata_dq_rule_template.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := WedataService{client: meta.(*TencentCloudClient).apiV3Conn}

	dq_ruleTemplateId := d.Id()

	dq_ruleTemplate, err := service.DescribeWedataDq_ruleTemplateById(ctx, templateId)
	if err != nil {
		return err
	}

	if dq_ruleTemplate == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `WedataDq_ruleTemplate` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if dq_ruleTemplate.Type != nil {
		_ = d.Set("type", dq_ruleTemplate.Type)
	}

	if dq_ruleTemplate.Name != nil {
		_ = d.Set("name", dq_ruleTemplate.Name)
	}

	if dq_ruleTemplate.QualityDim != nil {
		_ = d.Set("quality_dim", dq_ruleTemplate.QualityDim)
	}

	if dq_ruleTemplate.SourceObjectType != nil {
		_ = d.Set("source_object_type", dq_ruleTemplate.SourceObjectType)
	}

	if dq_ruleTemplate.Description != nil {
		_ = d.Set("description", dq_ruleTemplate.Description)
	}

	if dq_ruleTemplate.SourceEngineTypes != nil {
		_ = d.Set("source_engine_types", dq_ruleTemplate.SourceEngineTypes)
	}

	if dq_ruleTemplate.MultiSourceFlag != nil {
		_ = d.Set("multi_source_flag", dq_ruleTemplate.MultiSourceFlag)
	}

	if dq_ruleTemplate.SqlExpression != nil {
		_ = d.Set("sql_expression", dq_ruleTemplate.SqlExpression)
	}

	if dq_ruleTemplate.ProjectId != nil {
		_ = d.Set("project_id", dq_ruleTemplate.ProjectId)
	}

	if dq_ruleTemplate.WhereFlag != nil {
		_ = d.Set("where_flag", dq_ruleTemplate.WhereFlag)
	}

	return nil
}

func resourceTencentCloudWedataDq_ruleTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_wedata_dq_rule_template.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := wedata.NewModifyRuleTemplateRequest()

	dq_ruleTemplateId := d.Id()

	request.TemplateId = &templateId

	immutableArgs := []string{"type", "name", "quality_dim", "source_object_type", "description", "source_engine_types", "multi_source_flag", "sql_expression", "project_id", "where_flag"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("type") {
		if v, ok := d.GetOkExists("type"); ok {
			request.Type = helper.IntUint64(v.(int))
		}
	}

	if d.HasChange("name") {
		if v, ok := d.GetOk("name"); ok {
			request.Name = helper.String(v.(string))
		}
	}

	if d.HasChange("quality_dim") {
		if v, ok := d.GetOkExists("quality_dim"); ok {
			request.QualityDim = helper.IntUint64(v.(int))
		}
	}

	if d.HasChange("source_object_type") {
		if v, ok := d.GetOkExists("source_object_type"); ok {
			request.SourceObjectType = helper.IntUint64(v.(int))
		}
	}

	if d.HasChange("description") {
		if v, ok := d.GetOk("description"); ok {
			request.Description = helper.String(v.(string))
		}
	}

	if d.HasChange("source_engine_types") {
		if v, ok := d.GetOk("source_engine_types"); ok {
			sourceEngineTypesSet := v.(*schema.Set).List()
			for i := range sourceEngineTypesSet {
				sourceEngineTypes := sourceEngineTypesSet[i].(int)
				request.SourceEngineTypes = append(request.SourceEngineTypes, helper.IntUint64(sourceEngineTypes))
			}
		}
	}

	if d.HasChange("multi_source_flag") {
		if v, ok := d.GetOkExists("multi_source_flag"); ok {
			request.MultiSourceFlag = helper.Bool(v.(bool))
		}
	}

	if d.HasChange("sql_expression") {
		if v, ok := d.GetOk("sql_expression"); ok {
			request.SqlExpression = helper.String(v.(string))
		}
	}

	if d.HasChange("project_id") {
		if v, ok := d.GetOk("project_id"); ok {
			request.ProjectId = helper.String(v.(string))
		}
	}

	if d.HasChange("where_flag") {
		if v, ok := d.GetOkExists("where_flag"); ok {
			request.WhereFlag = helper.Bool(v.(bool))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseWedataClient().ModifyRuleTemplate(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update wedata dq_ruleTemplate failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudWedataDq_ruleTemplateRead(d, meta)
}

func resourceTencentCloudWedataDq_ruleTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_wedata_dq_rule_template.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := WedataService{client: meta.(*TencentCloudClient).apiV3Conn}
	dq_ruleTemplateId := d.Id()

	if err := service.DeleteWedataDq_ruleTemplateById(ctx, templateId); err != nil {
		return err
	}

	return nil
}
