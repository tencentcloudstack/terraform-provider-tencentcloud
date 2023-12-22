package sqlserver

import (
	"context"
	"fmt"
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sqlserver "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sqlserver/v20180328"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudSqlserverConfigInstanceParam() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudSqlserverConfigInstanceParamCreate,
		Read:   resourceTencentCloudSqlserverConfigInstanceParamRead,
		Update: resourceTencentCloudSqlserverConfigInstanceParamUpdate,
		Delete: resourceTencentCloudSqlserverConfigInstanceParamDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance ID.",
			},
			"param_list": {
				Required:    true,
				Type:        schema.TypeList,
				Description: "List of modified parameters. Each list element has two fields: Name and CurrentValue. Set Name to the parameter name and CurrentValue to the new value after modification. Note: if the instance needs to be restarted for the modified parameter to take effect, it will be restarted immediately or during the maintenance time. Before you modify a parameter, you can use the DescribeInstanceParams API to query whether the instance needs to be restarted.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Parameter name.",
						},
						"current_value": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Parameter value.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudSqlserverConfigInstanceParamCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_sqlserver_config_instance_param.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var instanceId string

	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	params := make([]string, 0)
	params = append(params, instanceId)
	if v, ok := d.GetOk("param_list"); ok {
		for _, item := range v.([]interface{}) {
			if item != nil {
				dMap := item.(map[string]interface{})
				if v, ok := dMap["name"]; ok {
					params = append(params, v.(string))
				}
			}
		}
	}

	d.SetId(strings.Join(params, tccommon.FILED_SP))

	return resourceTencentCloudSqlserverConfigInstanceParamUpdate(d, meta)
}

func resourceTencentCloudSqlserverConfigInstanceParamRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_sqlserver_config_instance_param.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = SqlserverService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) == 1 {
		return fmt.Errorf("id is broken, id is %s", d.Id())
	}
	instanceId := idSplit[0]
	params := idSplit[1:]

	configInstanceParam, err := service.DescribeSqlserverConfigInstanceParamById(ctx, instanceId)
	if err != nil {
		return err
	}

	if configInstanceParam == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `SqlserverConfigInstanceParam` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	paramItemsList := []interface{}{}
	if len(params) > 0 {
		for _, param := range params {
			for _, paramList := range configInstanceParam {
				if *paramList.Name == param {
					paramListMap := map[string]interface{}{}
					if paramList.Name != nil {
						paramListMap["name"] = paramList.Name
					}
					if paramList.CurrentValue != nil {
						paramListMap["current_value"] = paramList.CurrentValue
					}
					paramItemsList = append(paramItemsList, paramListMap)
				}
			}
		}
		_ = d.Set("param_list", paramItemsList)
	}

	_ = d.Set("instance_id", instanceId)

	return nil
}

func resourceTencentCloudSqlserverConfigInstanceParamUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_sqlserver_config_instance_param.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		request = sqlserver.NewModifyInstanceParamRequest()
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) == 1 {
		return fmt.Errorf("id is broken, id is %s", d.Id())
	}
	instanceId := idSplit[0]

	request.InstanceIds = []*string{&instanceId}

	if v, ok := d.GetOk("param_list"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			parameter := sqlserver.Parameter{}
			if v, ok := dMap["name"]; ok {
				parameter.Name = helper.String(v.(string))
			}
			if v, ok := dMap["current_value"]; ok {
				parameter.CurrentValue = helper.String(v.(string))
			}
			request.ParamList = append(request.ParamList, &parameter)
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseSqlserverClient().ModifyInstanceParam(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil {
			e = fmt.Errorf("sqlserver configInstanceParam not exists")
			return resource.NonRetryableError(e)
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s update sqlserver configInstanceParam failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudSqlserverConfigInstanceParamRead(d, meta)
}

func resourceTencentCloudSqlserverConfigInstanceParamDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_sqlserver_config_instance_param.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
