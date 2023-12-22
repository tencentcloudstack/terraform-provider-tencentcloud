package mariadb

import (
	"context"
	"fmt"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mariadb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mariadb/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudMariadbParameters() *schema.Resource {
	return &schema.Resource{
		Read:   resourceTencentCloudMariadbParametersRead,
		Create: resourceTencentCloudMariadbParametersCreate,
		Update: resourceTencentCloudMariadbParametersUpdate,
		Delete: resourceTencentCloudMariadbParametersDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "instance id.",
			},

			"params": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Number of days to keep, no more than 30.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"param": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "parameter name.",
						},
						"value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "parameter value.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudMariadbParametersCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mariadb_parameters.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var instanceId string
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	d.SetId(instanceId)
	return resourceTencentCloudMariadbParametersUpdate(d, meta)
}

func resourceTencentCloudMariadbParametersRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mariadb_parameters.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := MariadbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	instanceId := d.Id()

	parameters, err := service.DescribeMariadbParameters(ctx, instanceId)

	if err != nil {
		return err
	}

	if parameters == nil {
		d.SetId("")
		return fmt.Errorf("resource `parameters` %s does not exist", instanceId)
	}

	if parameters.InstanceId != nil {
		_ = d.Set("instance_id", parameters.InstanceId)
	}

	if parameters.Params != nil {
		paramsList := []interface{}{}
		for _, params := range parameters.Params {
			paramsMap := map[string]interface{}{}
			if params.Param != nil {
				paramsMap["param"] = params.Param
			}
			if params.Value != nil {
				paramsMap["value"] = params.Value
			}

			paramsList = append(paramsList, paramsMap)
		}
		_ = d.Set("params", paramsList)
	}

	return nil
}

func resourceTencentCloudMariadbParametersUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mariadb_parameters.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	request := mariadb.NewModifyDBParametersRequest()

	instanceId := d.Id()

	request.InstanceId = &instanceId

	if v, ok := d.GetOk("params"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			dBParamValue := mariadb.DBParamValue{}
			if v, ok := dMap["param"]; ok {
				dBParamValue.Param = helper.String(v.(string))
			}
			if v, ok := dMap["value"]; ok {
				dBParamValue.Value = helper.String(v.(string))
			}

			request.Params = append(request.Params, &dBParamValue)
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMariadbClient().ModifyDBParameters(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create mariadb parameters failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudMariadbParametersRead(d, meta)
}

func resourceTencentCloudMariadbParametersDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mariadb_parameters.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
