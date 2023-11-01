/*
Provides a resource to create a wedata dq_rule_template

Example Usage

```hcl
resource "tencentcloud_wedata_dq_rule_template" "example" {
  type                = 2
  name                = "tf_example"
  quality_dim         = 1
  source_object_type  = 2
  description         = "description."
  source_engine_types = [2]
  multi_source_flag   = true
  sql_expression      = "c2VsZWN0"
  project_id          = "1948767646355341312"
  where_flag          = true
}
```

Import

wedata dq_rule_template can be imported using the id, e.g.

```
terraform import tencentcloud_wedata_dq_rule_template.example 1948767646355341312#9480
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	wedata "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/wedata/v20210820"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudWedataDqRuleTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudWedataDqRuleTemplateCreate,
		Read:   resourceTencentCloudWedataDqRuleTemplateRead,
		Update: resourceTencentCloudWedataDqRuleTemplateUpdate,
		Delete: resourceTencentCloudWedataDqRuleTemplateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"type": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Template Type 1. System Template 2. User-defined template.",
			},
			"name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Template name.",
			},
			"quality_dim": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Quality detection dimension 1. Accuracy 2. Uniqueness 3. Completeness 4. Consistency 5. Timeliness 6. effectiveness.",
			},
			"source_object_type": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Source end data object type 1. Constant 2. Offline table level 2. Offline field level.",
			},
			"description": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Template description.",
			},
			"source_engine_types": {
				Required:    true,
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeInt},
				Description: "Type of the engine on the source end.",
			},
			"multi_source_flag": {
				Required:    true,
				Type:        schema.TypeBool,
				Description: "Whether to associate other tables.",
			},
			"sql_expression": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "SQL expression.",
			},
			"project_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Project id.",
			},
			"where_flag": {
				Required:    true,
				Type:        schema.TypeBool,
				Description: "Add where parameter or not.",
			},
			"template_id": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Template ID.",
			},
		},
	}
}

func resourceTencentCloudWedataDqRuleTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_wedata_dq_rule_template.create")()
	defer inconsistentCheck(d, meta)()

	var (
		logId      = getLogId(contextNil)
		request    = wedata.NewCreateRuleTemplateRequest()
		response   = wedata.NewCreateRuleTemplateResponse()
		templateId string
		projectId  string
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
		projectId = v.(string)
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

		if result == nil {
			e = fmt.Errorf("wedata dqRuleTemplate not exists")
			return resource.NonRetryableError(e)
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create wedata dqRuleTemplate failed, reason:%+v", logId, err)
		return err
	}

	data := *response.Response.Data
	templateId = helper.UInt64ToStr(data)
	d.SetId(strings.Join([]string{projectId, templateId}, FILED_SP))

	return resourceTencentCloudWedataDqRuleTemplateRead(d, meta)
}

func resourceTencentCloudWedataDqRuleTemplateRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_wedata_dq_rule_template.read")()
	defer inconsistentCheck(d, meta)()

	var (
		logId   = getLogId(contextNil)
		ctx     = context.WithValue(context.TODO(), logIdKey, logId)
		service = WedataService{client: meta.(*TencentCloudClient).apiV3Conn}
	)

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", idSplit)
	}
	projectId := idSplit[0]
	templateId := idSplit[1]

	dqRuleTemplate, err := service.DescribeWedataDqRuleTemplateById(ctx, projectId, templateId)
	if err != nil {
		return err
	}

	if dqRuleTemplate == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `WedataDqRuleTemplate` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("project_id", projectId)
	_ = d.Set("template_id", templateId)

	if dqRuleTemplate.Type != nil {
		_ = d.Set("type", dqRuleTemplate.Type)
	}

	if dqRuleTemplate.Name != nil {
		_ = d.Set("name", dqRuleTemplate.Name)
	}

	if dqRuleTemplate.QualityDim != nil {
		_ = d.Set("quality_dim", dqRuleTemplate.QualityDim)
	}

	if dqRuleTemplate.SourceObjectType != nil {
		_ = d.Set("source_object_type", dqRuleTemplate.SourceObjectType)
	}

	if dqRuleTemplate.Description != nil {
		_ = d.Set("description", dqRuleTemplate.Description)
	}

	if dqRuleTemplate.SourceEngineTypes != nil {
		_ = d.Set("source_engine_types", dqRuleTemplate.SourceEngineTypes)
	}

	if dqRuleTemplate.MultiSourceFlag != nil {
		_ = d.Set("multi_source_flag", dqRuleTemplate.MultiSourceFlag)
	}

	if dqRuleTemplate.SqlExpression != nil {
		_ = d.Set("sql_expression", dqRuleTemplate.SqlExpression)
	}

	if dqRuleTemplate.WhereFlag != nil {
		_ = d.Set("where_flag", dqRuleTemplate.WhereFlag)
	}

	return nil
}

func resourceTencentCloudWedataDqRuleTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_wedata_dq_rule_template.update")()
	defer inconsistentCheck(d, meta)()

	var (
		logId   = getLogId(contextNil)
		request = wedata.NewModifyRuleTemplateRequest()
	)

	immutableArgs := []string{"project_id", "multi_source_flag", "where_flag"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", idSplit)
	}
	projectId := idSplit[0]
	templateId := idSplit[1]

	request.ProjectId = &projectId
	templateIdInt, _ := strconv.ParseUint(templateId, 10, 64)
	request.TemplateId = &templateIdInt

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

	if v, ok := d.GetOkExists("where_flag"); ok {
		request.WhereFlag = helper.Bool(v.(bool))
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
		log.Printf("[CRITAL]%s update wedata dqRuleTemplate failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudWedataDqRuleTemplateRead(d, meta)
}

func resourceTencentCloudWedataDqRuleTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_wedata_dq_rule_template.delete")()
	defer inconsistentCheck(d, meta)()

	var (
		logId   = getLogId(contextNil)
		ctx     = context.WithValue(context.TODO(), logIdKey, logId)
		service = WedataService{client: meta.(*TencentCloudClient).apiV3Conn}
	)

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", idSplit)
	}
	projectId := idSplit[0]
	templateId := idSplit[1]

	if err := service.DeleteWedataDqRuleTemplateById(ctx, projectId, templateId); err != nil {
		return err
	}

	return nil
}
