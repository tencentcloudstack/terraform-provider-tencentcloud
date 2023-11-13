/*
Provides a resource to create a sqlserver config_instance_param

Example Usage

```hcl
resource "tencentcloud_sqlserver_config_instance_param" "config_instance_param" {
  instance_ids =
  param_list {
		name = "fill factor(%)"
		current_value = "90"

  }
  wait_switch = 0
}
```

Import

sqlserver config_instance_param can be imported using the id, e.g.

```
terraform import tencentcloud_sqlserver_config_instance_param.config_instance_param config_instance_param_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sqlserver "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sqlserver/v20180328"
	"log"
)

func resourceTencentCloudSqlserverConfigInstanceParam() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudSqlserverConfigInstanceParamCreate,
		Read:   resourceTencentCloudSqlserverConfigInstanceParamRead,
		Update: resourceTencentCloudSqlserverConfigInstanceParamUpdate,
		Delete: resourceTencentCloudSqlserverConfigInstanceParamDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_ids": {
				Required: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Instance ID list.",
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

			"wait_switch": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "When to execute the parameter modification task. Valid values: 0 (execute immediately), 1 (execute during maintenance time). Default value: 0.",
			},
		},
	}
}

func resourceTencentCloudSqlserverConfigInstanceParamCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_config_instance_param.create")()
	defer inconsistentCheck(d, meta)()

	var instanceId string
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	d.SetId(instanceId)

	return resourceTencentCloudSqlserverConfigInstanceParamUpdate(d, meta)
}

func resourceTencentCloudSqlserverConfigInstanceParamRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_config_instance_param.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := SqlserverService{client: meta.(*TencentCloudClient).apiV3Conn}

	configInstanceParamId := d.Id()

	configInstanceParam, err := service.DescribeSqlserverConfigInstanceParamById(ctx, instanceId)
	if err != nil {
		return err
	}

	if configInstanceParam == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `SqlserverConfigInstanceParam` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if configInstanceParam.InstanceIds != nil {
		_ = d.Set("instance_ids", configInstanceParam.InstanceIds)
	}

	if configInstanceParam.ParamList != nil {
		paramListList := []interface{}{}
		for _, paramList := range configInstanceParam.ParamList {
			paramListMap := map[string]interface{}{}

			if configInstanceParam.ParamList.Name != nil {
				paramListMap["name"] = configInstanceParam.ParamList.Name
			}

			if configInstanceParam.ParamList.CurrentValue != nil {
				paramListMap["current_value"] = configInstanceParam.ParamList.CurrentValue
			}

			paramListList = append(paramListList, paramListMap)
		}

		_ = d.Set("param_list", paramListList)

	}

	if configInstanceParam.WaitSwitch != nil {
		_ = d.Set("wait_switch", configInstanceParam.WaitSwitch)
	}

	return nil
}

func resourceTencentCloudSqlserverConfigInstanceParamUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_config_instance_param.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := sqlserver.NewModifyInstanceParamRequest()

	configInstanceParamId := d.Id()

	request.InstanceId = &instanceId

	immutableArgs := []string{"instance_ids", "param_list", "wait_switch"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseSqlserverClient().ModifyInstanceParam(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
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
	defer logElapsed("resource.tencentcloud_sqlserver_config_instance_param.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
