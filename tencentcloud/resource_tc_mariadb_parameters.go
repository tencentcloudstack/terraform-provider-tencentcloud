/*
Provides a resource to create a mariadb parameters

Example Usage

```hcl
resource "tencentcloud_mariadb_parameters" "parameters" {
  instance_id = &lt;nil&gt;
  params {
		param = &lt;nil&gt;
		value = &lt;nil&gt;

  }
}
```

Import

mariadb parameters can be imported using the id, e.g.

```
terraform import tencentcloud_mariadb_parameters.parameters parameters_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mariadb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mariadb/v20170312"
	"log"
)

func resourceTencentCloudMariadbParameters() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMariadbParametersCreate,
		Read:   resourceTencentCloudMariadbParametersRead,
		Update: resourceTencentCloudMariadbParametersUpdate,
		Delete: resourceTencentCloudMariadbParametersDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance id.",
			},

			"params": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "Parameters list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"param": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Parameter name.",
						},
						"value": {
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

	parametersId := d.Id()

	parameters, err := service.DescribeMariadbParametersById(ctx, instanceId)
	if err != nil {
		return err
	}

	if parameters == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `MariadbParameters` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if parameters.InstanceId != nil {
		_ = d.Set("instance_id", parameters.InstanceId)
	}

	if parameters.Params != nil {
		paramsList := []interface{}{}
		for _, params := range parameters.Params {
			paramsMap := map[string]interface{}{}

			if parameters.Params.Param != nil {
				paramsMap["param"] = parameters.Params.Param
			}

			if parameters.Params.Value != nil {
				paramsMap["value"] = parameters.Params.Value
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

	parametersId := d.Id()

	request.InstanceId = &instanceId

	immutableArgs := []string{"instance_id", "params"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseMariadbClient().ModifyDBParameters(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update mariadb parameters failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudMariadbParametersRead(d, meta)
}

func resourceTencentCloudMariadbParametersDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mariadb_parameters.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
