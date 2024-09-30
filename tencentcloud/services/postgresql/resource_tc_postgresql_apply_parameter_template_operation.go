package postgresql

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	postgres "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/postgres/v20170312"
)

func ResourceTencentCloudPostgresqlApplyParameterTemplateOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudPostgresqlApplyParameterTemplateOperationCreate,
		Read:   resourceTencentCloudPostgresqlApplyParameterTemplateOperationRead,
		Delete: resourceTencentCloudPostgresqlApplyParameterTemplateOperationDelete,
		Schema: map[string]*schema.Schema{
			"db_instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Instance ID.",
			},
			"template_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Template ID.",
			},
		},
	}
}

func resourceTencentCloudPostgresqlApplyParameterTemplateOperationCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_postgresql_apply_parameter_template_operation.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := PostgresqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var (
		err, innerErr        error
		dbParmas, diffParmas map[string]string
		templateAttributes   *postgres.DescribeParameterTemplateAttributesResponseParams
	)
	dbInstanceId := d.Get("db_instance_id").(string)
	err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		dbParmas, innerErr = service.DescribePgParams(ctx, dbInstanceId)
		if innerErr != nil {
			return tccommon.RetryError(innerErr)
		}

		return nil
	})

	if err != nil {
		return err
	}

	templateId := d.Get("template_id").(string)
	err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		templateAttributes, innerErr = service.DescribePostgresqlParameterTemplateById(ctx, templateId)
		if innerErr != nil {
			return tccommon.RetryError(innerErr)
		}

		return nil
	})

	if err != nil {
		return err
	}

	diffParmas = make(map[string]string)
	if templateAttributes != nil {
		for _, param := range templateAttributes.ParamInfoSet {
			name := param.Name
			value := param.CurrentValue
			if name != nil && value != nil && dbParmas[*name] != *value {
				diffParmas[*name] = *value
			}
		}
	}

	err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		innerErr = service.ModifyPgParams(ctx, dbInstanceId, diffParmas)
		if innerErr != nil {
			return tccommon.RetryError(innerErr)
		}

		return nil
	})

	if err != nil {
		return err
	}

	err = service.CheckDBInstanceStatus(ctx, dbInstanceId)
	if err != nil {
		return err
	}

	d.SetId(dbInstanceId + tccommon.FILED_SP + templateId)

	return resourceTencentCloudPostgresqlApplyParameterTemplateOperationRead(d, meta)
}

func resourceTencentCloudPostgresqlApplyParameterTemplateOperationRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_postgresql_apply_parameter_template_operation.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudPostgresqlApplyParameterTemplateOperationDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_postgresql_apply_parameter_template_operation.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
