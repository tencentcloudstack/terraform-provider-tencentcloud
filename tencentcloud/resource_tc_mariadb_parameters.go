package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mariadb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mariadb/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudMariadbParameters() *schema.Resource {
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
	defer logElapsed("resource.tencentcloud_mariadb_parameters.create")()
	defer inconsistentCheck(d, meta)()

	var instanceId string
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	d.SetId(instanceId)
	return resourceTencentCloudMariadbParametersUpdate(d, meta)
}

func resourceTencentCloudMariadbParametersRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mariadb_parameters.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := MariadbService{client: meta.(*TencentCloudClient).apiV3Conn}

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
	defer logElapsed("resource.tencentcloud_mariadb_parameters.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

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

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseMariadbClient().ModifyDBParameters(request)
		if e != nil {
			return retryError(e)
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
	defer logElapsed("resource.tencentcloud_mariadb_parameters.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
