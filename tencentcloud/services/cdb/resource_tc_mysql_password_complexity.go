package cdb

import (
	"context"
	"fmt"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mysql "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudMysqlPasswordComplexity() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMysqlPasswordComplexityCreate,
		Read:   resourceTencentCloudMysqlPasswordComplexityRead,
		Update: resourceTencentCloudMysqlPasswordComplexityUpdate,
		Delete: resourceTencentCloudMysqlPasswordComplexityDelete,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance ID.",
			},

			"param_list": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "List of parameters to be modified. Every element is a combination of `Name` (parameter name) and `CurrentValue` (new value). Valid values for `Name` of version 8.0: `validate_password.policy`, `validate_password.length`, `validate_password.mixed_case_count`, `validate_password.number_count`, `validate_password.special_char_count`. Valid values for `Name` of version 5.6 and 5.7: `validate_password_policy`, `validate_password_length` `validate_password_mixed_case_count`, `validate_password_number_count`, `validate_password_special_char_count`.",
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

func resourceTencentCloudMysqlPasswordComplexityCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mysql_password_complexity.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var instanceId string
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	d.SetId(instanceId)

	return resourceTencentCloudMysqlPasswordComplexityUpdate(d, meta)
}

func resourceTencentCloudMysqlPasswordComplexityRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mysql_password_complexity.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := MysqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	instanceId := d.Id()

	passwordComplexitys, err := service.DescribeMysqlPasswordComplexityById(ctx, instanceId)
	if err != nil {
		return err
	}

	if passwordComplexitys == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `MysqlPasswordComplexity` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("instance_id", instanceId)

	return nil
}

func resourceTencentCloudMysqlPasswordComplexityUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mysql_password_complexity.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	request := mysql.NewModifyInstancePasswordComplexityRequest()
	response := mysql.NewModifyInstancePasswordComplexityResponse()

	instanceId := d.Id()

	request.InstanceIds = []*string{&instanceId}

	if v, ok := d.GetOk("param_list"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			parameter := mysql.Parameter{}
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
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMysqlClient().ModifyInstancePasswordComplexity(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update mysql passwordComplexity failed, reason:%+v", logId, err)
		return err
	}

	asyncRequestId := *response.Response.AsyncRequestId
	service := MysqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		taskStatus, message, err := service.DescribeAsyncRequestInfo(ctx, asyncRequestId)
		if err != nil {
			return resource.NonRetryableError(err)
		}
		if taskStatus == MYSQL_TASK_STATUS_SUCCESS {
			return nil
		}
		if taskStatus == MYSQL_TASK_STATUS_INITIAL || taskStatus == MYSQL_TASK_STATUS_RUNNING {
			return resource.RetryableError(fmt.Errorf("%s update mysql passwordComplexity status is %s", instanceId, taskStatus))
		}
		err = fmt.Errorf("%s update mysql passwordComplexity status is %s,we won't wait for it finish ,it show message:%s", instanceId, taskStatus, message)
		return resource.NonRetryableError(err)
	})

	if err != nil {
		log.Printf("[CRITAL]%s update mysql passwordComplexity fail, reason:%s\n ", logId, err.Error())
		return err
	}

	return resourceTencentCloudMysqlPasswordComplexityRead(d, meta)
}

func resourceTencentCloudMysqlPasswordComplexityDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mysql_password_complexity.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
