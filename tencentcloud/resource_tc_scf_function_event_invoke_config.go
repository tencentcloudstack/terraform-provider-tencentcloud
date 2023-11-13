/*
Provides a resource to create a scf function_event_invoke_config

Example Usage

```hcl
resource "tencentcloud_scf_function_event_invoke_config" "function_event_invoke_config" {
  async_trigger_config {
		retry_config {
			retry_num = 2
		}
		msg_t_t_l = 24

  }
  function_name = "test_function"
  namespace = "test_namespace"
}
```

Import

scf function_event_invoke_config can be imported using the id, e.g.

```
terraform import tencentcloud_scf_function_event_invoke_config.function_event_invoke_config function_event_invoke_config_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	scf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/scf/v20180416"
	"log"
	"strings"
)

func resourceTencentCloudScfFunctionEventInvokeConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudScfFunctionEventInvokeConfigCreate,
		Read:   resourceTencentCloudScfFunctionEventInvokeConfigRead,
		Update: resourceTencentCloudScfFunctionEventInvokeConfigUpdate,
		Delete: resourceTencentCloudScfFunctionEventInvokeConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"async_trigger_config": {
				Required:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Async retry configuration information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"retry_config": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "Async retry configuration of function upon user error.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"retry_num": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Number of retry attempts.",
									},
								},
							},
						},
						"msg_t_t_l": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Message retention period.",
						},
					},
				},
			},

			"function_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Function name.",
			},

			"namespace": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Function namespace. Default value: default.",
			},
		},
	}
}

func resourceTencentCloudScfFunctionEventInvokeConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_scf_function_event_invoke_config.create")()
	defer inconsistentCheck(d, meta)()

	var namespace string
	if v, ok := d.GetOk("namespace"); ok {
		namespace = v.(string)
	}

	var functionName string
	if v, ok := d.GetOk("function_name"); ok {
		functionName = v.(string)
	}

	d.SetId(strings.Join([]string{namespace, functionName}, FILED_SP))

	return resourceTencentCloudScfFunctionEventInvokeConfigUpdate(d, meta)
}

func resourceTencentCloudScfFunctionEventInvokeConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_scf_function_event_invoke_config.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := ScfService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	namespace := idSplit[0]
	functionName := idSplit[1]

	FunctionEventInvokeConfig, err := service.DescribeScfFunctionEventInvokeConfigById(ctx, namespace, functionName)
	if err != nil {
		return err
	}

	if FunctionEventInvokeConfig == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `ScfFunctionEventInvokeConfig` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if FunctionEventInvokeConfig.AsyncTriggerConfig != nil {
		asyncTriggerConfigMap := map[string]interface{}{}

		if FunctionEventInvokeConfig.AsyncTriggerConfig.RetryConfig != nil {
			retryConfigList := []interface{}{}
			for _, retryConfig := range FunctionEventInvokeConfig.AsyncTriggerConfig.RetryConfig {
				retryConfigMap := map[string]interface{}{}

				if retryConfig.RetryNum != nil {
					retryConfigMap["retry_num"] = retryConfig.RetryNum
				}

				retryConfigList = append(retryConfigList, retryConfigMap)
			}

			asyncTriggerConfigMap["retry_config"] = []interface{}{retryConfigList}
		}

		if FunctionEventInvokeConfig.AsyncTriggerConfig.MsgTTL != nil {
			asyncTriggerConfigMap["msg_t_t_l"] = FunctionEventInvokeConfig.AsyncTriggerConfig.MsgTTL
		}

		_ = d.Set("async_trigger_config", []interface{}{asyncTriggerConfigMap})
	}

	if FunctionEventInvokeConfig.FunctionName != nil {
		_ = d.Set("function_name", FunctionEventInvokeConfig.FunctionName)
	}

	if FunctionEventInvokeConfig.Namespace != nil {
		_ = d.Set("namespace", FunctionEventInvokeConfig.Namespace)
	}

	return nil
}

func resourceTencentCloudScfFunctionEventInvokeConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_scf_function_event_invoke_config.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := scf.NewUpdateFunctionEventInvokeConfigRequest()

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	namespace := idSplit[0]
	functionName := idSplit[1]

	request.Namespace = &namespace
	request.FunctionName = &functionName

	immutableArgs := []string{"async_trigger_config", "function_name", "namespace"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseScfClient().UpdateFunctionEventInvokeConfig(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update scf FunctionEventInvokeConfig failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudScfFunctionEventInvokeConfigRead(d, meta)
}

func resourceTencentCloudScfFunctionEventInvokeConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_scf_function_event_invoke_config.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
