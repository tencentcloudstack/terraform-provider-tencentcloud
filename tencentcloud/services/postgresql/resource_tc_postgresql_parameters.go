// Code generated by iacg; DO NOT EDIT.
package postgresql

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	postgresv20170312 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/postgres/v20170312"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudPostgresqlParameters() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudPostgresqlParametersCreate,
		Read:   resourceTencentCloudPostgresqlParametersRead,
		Update: resourceTencentCloudPostgresqlParametersUpdate,
		Delete: resourceTencentCloudPostgresqlParametersDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"db_instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Instance ID.",
			},

			"param_list": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Parameters to be modified and expected values.",
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
						"default_value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The default value of the parameter. Returned as a string.",
						},
						"param_description_ch": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Parameter Chinese Description.",
						},
						"param_description_en": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Parameter English Description.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudPostgresqlParametersCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_postgresql_parameters.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		dBInstanceId string
	)
	if v, ok := d.GetOk("db_instance_id"); ok {
		dBInstanceId = v.(string)
	}

	d.SetId(dBInstanceId)

	return resourceTencentCloudPostgresqlParametersUpdate(d, meta)
}

func resourceTencentCloudPostgresqlParametersRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_postgresql_parameters.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	service := PostgresqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	dBInstanceId := d.Id()

	_ = d.Set("db_instance_id", dBInstanceId)

	respData, err := service.DescribePostgresqlParametersById(ctx, dBInstanceId)
	if err != nil {
		return err
	}

	if respData == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `postgresql_parameters` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}
	if err := resourceTencentCloudPostgresqlParametersReadPreHandleResponse0(ctx, respData); err != nil {
		return err
	}

	detailList := make([]map[string]interface{}, 0, len(respData.Detail))
	if respData.Detail != nil {
		for _, detail := range respData.Detail {
			detailMap := map[string]interface{}{}

			if detail.Name != nil {
				detailMap["name"] = detail.Name
			}

			if detail.DefaultValue != nil {
				detailMap["default_value"] = detail.DefaultValue
			}

			if detail.CurrentValue != nil {
				detailMap["expected_value"] = detail.CurrentValue
			}

			if detail.ParamDescriptionCH != nil {
				detailMap["param_description_ch"] = detail.ParamDescriptionCH
			}

			if detail.ParamDescriptionEN != nil {
				detailMap["param_description_en"] = detail.ParamDescriptionEN
			}

			detailList = append(detailList, detailMap)
		}

		_ = d.Set("param_list", detailList)
	}

	return nil
}

func resourceTencentCloudPostgresqlParametersUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_postgresql_parameters.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	dBInstanceId := d.Id()

	needChange := false
	mutableArgs := []string{"param_list"}
	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {
		request := postgresv20170312.NewModifyDBInstanceParametersRequest()
		response := postgresv20170312.NewModifyDBInstanceParametersResponse()

		if v, ok := d.GetOk("db_instance_id"); ok {
			request.DBInstanceId = helper.String(v.(string))
		}

		if v, ok := d.GetOk("param_list"); ok {
			for _, item := range v.([]interface{}) {
				paramListMap := item.(map[string]interface{})
				paramEntry := postgresv20170312.ParamEntry{}
				if v, ok := paramListMap["name"].(string); ok && v != "" {
					paramEntry.Name = helper.String(v)
				}
				if v, ok := paramListMap["expected_value"].(string); ok && v != "" {
					paramEntry.ExpectedValue = helper.String(v)
				}
				request.ParamList = append(request.ParamList, &paramEntry)
			}
		}

		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UsePostgresV20170312Client().ModifyDBInstanceParametersWithContext(ctx, request)
			if e != nil {
				return resourceTencentCloudPostgresqlParametersUpdateRequestOnError0(ctx, request, e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			response = result
			return nil
		})
		if reqErr != nil {
			log.Printf("[CRITAL]%s update postgresql parameters failed, reason:%+v", logId, reqErr)
			return reqErr
		}
		if err := resourceTencentCloudPostgresqlParametersUpdatePreHandleResponse0(ctx, response); err != nil {
			return err
		}
	}

	_ = dBInstanceId
	return resourceTencentCloudPostgresqlParametersRead(d, meta)
}

func resourceTencentCloudPostgresqlParametersDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_postgresql_parameters.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
