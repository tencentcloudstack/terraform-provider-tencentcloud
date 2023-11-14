/*
Provides a resource to create a scf provisioned_concurrency_config

Example Usage

```hcl
resource "tencentcloud_scf_provisioned_concurrency_config" "provisioned_concurrency_config" {
  function_name = "test_function"
  qualifier = "1"
  version_provisioned_concurrency_num = 2
  namespace = "test_namespace"
  trigger_actions {
		trigger_name = ""
		trigger_provisioned_concurrency_num =
		trigger_cron_config = ""
		provisioned_type = ""

  }
  provisioned_type = ""
  tracking_target =
  min_capacity =
  max_capacity =
}
```

Import

scf provisioned_concurrency_config can be imported using the id, e.g.

```
terraform import tencentcloud_scf_provisioned_concurrency_config.provisioned_concurrency_config provisioned_concurrency_config_id
```
*/
package tencentcloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	scf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/scf/v20180416"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
	"strings"
)

func resourceTencentCloudScfProvisionedConcurrencyConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudScfProvisionedConcurrencyConfigCreate,
		Read:   resourceTencentCloudScfProvisionedConcurrencyConfigRead,
		Delete: resourceTencentCloudScfProvisionedConcurrencyConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"function_name": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Name of the function for which to set the provisioned concurrency.",
			},

			"qualifier": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Function version number. Note: the $LATEST version does not support provisioned concurrency.",
			},

			"version_provisioned_concurrency_num": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "Provisioned concurrency amount. Note: there is an upper limit for the sum of provisioned concurrency amounts of all versions, which currently is the function&amp;amp;#39;s maximum concurrency quota minus 100.",
			},

			"namespace": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Function namespace. Default value: default.",
			},

			"trigger_actions": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				Description: "Scheduled provisioned concurrency scaling action.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"trigger_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Scheduled action name Note: this field may return null, indicating that no valid values can be obtained.",
						},
						"trigger_provisioned_concurrency_num": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Target provisioned concurrency of the scheduled scaling action Note: this field may return null, indicating that no valid values can be obtained.",
						},
						"trigger_cron_config": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Trigger time of the scheduled action in Cron expression. Seven fields are required and should be separated with a space. Note: this field may return null, indicating that no valid values can be obtained.",
						},
						"provisioned_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The provision type. Value: Default Note: This field may return null, indicating that no valid value can be found.",
						},
					},
				},
			},

			"provisioned_type": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Specifies the provisioned concurrency type. Default: Static provisioned concurrency. ConcurrencyUtilizationTracking: Scales the concurrency automatically according to the concurrency utilization. If ConcurrencyUtilizationTracking is passed in, TrackingTarget, MinCapacity and MaxCapacity are required, and VersionProvisionedConcurrencyNum must be 0.",
			},

			"tracking_target": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeFloat,
				Description: "The target concurrency utilization. Range: (0,1) (two decimal places).",
			},

			"min_capacity": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "The minimum number of instances. It can not be smaller than 1.",
			},

			"max_capacity": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "The maximum number of instances.",
			},
		},
	}
}

func resourceTencentCloudScfProvisionedConcurrencyConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_scf_provisioned_concurrency_config.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request      = scf.NewPutProvisionedConcurrencyConfigRequest()
		response     = scf.NewPutProvisionedConcurrencyConfigResponse()
		functionName string
		qualifier    string
		namespace    string
	)
	if v, ok := d.GetOk("function_name"); ok {
		functionName = v.(string)
		request.FunctionName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("qualifier"); ok {
		qualifier = v.(string)
		request.Qualifier = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("version_provisioned_concurrency_num"); ok {
		request.VersionProvisionedConcurrencyNum = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("namespace"); ok {
		namespace = v.(string)
		request.Namespace = helper.String(v.(string))
	}

	if v, ok := d.GetOk("trigger_actions"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			triggerAction := scf.TriggerAction{}
			if v, ok := dMap["trigger_name"]; ok {
				triggerAction.TriggerName = helper.String(v.(string))
			}
			if v, ok := dMap["trigger_provisioned_concurrency_num"]; ok {
				triggerAction.TriggerProvisionedConcurrencyNum = helper.IntUint64(v.(int))
			}
			if v, ok := dMap["trigger_cron_config"]; ok {
				triggerAction.TriggerCronConfig = helper.String(v.(string))
			}
			if v, ok := dMap["provisioned_type"]; ok {
				triggerAction.ProvisionedType = helper.String(v.(string))
			}
			request.TriggerActions = append(request.TriggerActions, &triggerAction)
		}
	}

	if v, ok := d.GetOk("provisioned_type"); ok {
		request.ProvisionedType = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("tracking_target"); ok {
		request.TrackingTarget = helper.Float64(v.(float64))
	}

	if v, ok := d.GetOkExists("min_capacity"); ok {
		request.MinCapacity = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("max_capacity"); ok {
		request.MaxCapacity = helper.IntUint64(v.(int))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseScfClient().PutProvisionedConcurrencyConfig(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create scf ProvisionedConcurrencyConfig failed, reason:%+v", logId, err)
		return err
	}

	functionName = *response.Response.FunctionName
	d.SetId(strings.Join([]string{functionName, qualifier, namespace}, FILED_SP))

	return resourceTencentCloudScfProvisionedConcurrencyConfigRead(d, meta)
}

func resourceTencentCloudScfProvisionedConcurrencyConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_scf_provisioned_concurrency_config.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := ScfService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	functionName := idSplit[0]
	qualifier := idSplit[1]
	namespace := idSplit[2]

	ProvisionedConcurrencyConfig, err := service.DescribeScfProvisionedConcurrencyConfigById(ctx, functionName, qualifier, namespace)
	if err != nil {
		return err
	}

	if ProvisionedConcurrencyConfig == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `ScfProvisionedConcurrencyConfig` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if ProvisionedConcurrencyConfig.FunctionName != nil {
		_ = d.Set("function_name", ProvisionedConcurrencyConfig.FunctionName)
	}

	if ProvisionedConcurrencyConfig.Qualifier != nil {
		_ = d.Set("qualifier", ProvisionedConcurrencyConfig.Qualifier)
	}

	if ProvisionedConcurrencyConfig.VersionProvisionedConcurrencyNum != nil {
		_ = d.Set("version_provisioned_concurrency_num", ProvisionedConcurrencyConfig.VersionProvisionedConcurrencyNum)
	}

	if ProvisionedConcurrencyConfig.Namespace != nil {
		_ = d.Set("namespace", ProvisionedConcurrencyConfig.Namespace)
	}

	if ProvisionedConcurrencyConfig.TriggerActions != nil {
		triggerActionsList := []interface{}{}
		for _, triggerActions := range ProvisionedConcurrencyConfig.TriggerActions {
			triggerActionsMap := map[string]interface{}{}

			if ProvisionedConcurrencyConfig.TriggerActions.TriggerName != nil {
				triggerActionsMap["trigger_name"] = ProvisionedConcurrencyConfig.TriggerActions.TriggerName
			}

			if ProvisionedConcurrencyConfig.TriggerActions.TriggerProvisionedConcurrencyNum != nil {
				triggerActionsMap["trigger_provisioned_concurrency_num"] = ProvisionedConcurrencyConfig.TriggerActions.TriggerProvisionedConcurrencyNum
			}

			if ProvisionedConcurrencyConfig.TriggerActions.TriggerCronConfig != nil {
				triggerActionsMap["trigger_cron_config"] = ProvisionedConcurrencyConfig.TriggerActions.TriggerCronConfig
			}

			if ProvisionedConcurrencyConfig.TriggerActions.ProvisionedType != nil {
				triggerActionsMap["provisioned_type"] = ProvisionedConcurrencyConfig.TriggerActions.ProvisionedType
			}

			triggerActionsList = append(triggerActionsList, triggerActionsMap)
		}

		_ = d.Set("trigger_actions", triggerActionsList)

	}

	if ProvisionedConcurrencyConfig.ProvisionedType != nil {
		_ = d.Set("provisioned_type", ProvisionedConcurrencyConfig.ProvisionedType)
	}

	if ProvisionedConcurrencyConfig.TrackingTarget != nil {
		_ = d.Set("tracking_target", ProvisionedConcurrencyConfig.TrackingTarget)
	}

	if ProvisionedConcurrencyConfig.MinCapacity != nil {
		_ = d.Set("min_capacity", ProvisionedConcurrencyConfig.MinCapacity)
	}

	if ProvisionedConcurrencyConfig.MaxCapacity != nil {
		_ = d.Set("max_capacity", ProvisionedConcurrencyConfig.MaxCapacity)
	}

	return nil
}

func resourceTencentCloudScfProvisionedConcurrencyConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_scf_provisioned_concurrency_config.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := ScfService{client: meta.(*TencentCloudClient).apiV3Conn}
	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	functionName := idSplit[0]
	qualifier := idSplit[1]
	namespace := idSplit[2]

	if err := service.DeleteScfProvisionedConcurrencyConfigById(ctx, functionName, qualifier, namespace); err != nil {
		return err
	}

	return nil
}
