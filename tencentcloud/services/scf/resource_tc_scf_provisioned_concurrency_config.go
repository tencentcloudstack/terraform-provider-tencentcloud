package scf

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	scf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/scf/v20180416"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudScfProvisionedConcurrencyConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudScfProvisionedConcurrencyConfigCreate,
		Read:   resourceTencentCloudScfProvisionedConcurrencyConfigRead,
		Delete: resourceTencentCloudScfProvisionedConcurrencyConfigDelete,
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
				Description: "Provisioned concurrency amount. Note: there is an upper limit for the sum of provisioned concurrency amounts of all versions, which currently is the function&amp;#39;s maximum concurrency quota minus 100.",
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
	defer tccommon.LogElapsed("resource.tencentcloud_scf_provisioned_concurrency_config.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request      = scf.NewPutProvisionedConcurrencyConfigRequest()
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

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseScfClient().PutProvisionedConcurrencyConfig(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create scf ProvisionedConcurrencyConfig failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(functionName + tccommon.FILED_SP + qualifier + tccommon.FILED_SP + namespace)

	// dirty code
	time.Sleep(5 * time.Second)

	return resourceTencentCloudScfProvisionedConcurrencyConfigRead(d, meta)
}

func resourceTencentCloudScfProvisionedConcurrencyConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_scf_provisioned_concurrency_config.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := ScfService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	functionName := idSplit[0]
	qualifier := idSplit[1]
	namespace := idSplit[2]

	provisionedConcurrencyConfig, err := service.DescribeScfProvisionedConcurrencyConfigById(ctx, functionName, qualifier, namespace)
	if err != nil {
		return err
	}

	if provisionedConcurrencyConfig == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `ScfProvisionedConcurrencyConfig` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("function_name", functionName)
	_ = d.Set("qualifier", qualifier)
	_ = d.Set("namespace", namespace)

	if provisionedConcurrencyConfig.AvailableProvisionedConcurrencyNum != nil {
		_ = d.Set("version_provisioned_concurrency_num", provisionedConcurrencyConfig.AvailableProvisionedConcurrencyNum)
	}

	if provisionedConcurrencyConfig.TriggerActions != nil {
		triggerActionsList := []interface{}{}
		for _, action := range provisionedConcurrencyConfig.TriggerActions {
			triggerActionsMap := map[string]interface{}{}

			if action.TriggerName != nil {
				triggerActionsMap["trigger_name"] = action.TriggerName
			}

			if action.TriggerProvisionedConcurrencyNum != nil {
				triggerActionsMap["trigger_provisioned_concurrency_num"] = action.TriggerProvisionedConcurrencyNum
			}

			if action.TriggerCronConfig != nil {
				triggerActionsMap["trigger_cron_config"] = action.TriggerCronConfig
			}

			if action.ProvisionedType != nil {
				triggerActionsMap["provisioned_type"] = action.ProvisionedType
			}

			triggerActionsList = append(triggerActionsList, triggerActionsMap)
		}

		_ = d.Set("trigger_actions", triggerActionsList)

	}

	return nil
}

func resourceTencentCloudScfProvisionedConcurrencyConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_scf_provisioned_concurrency_config.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := ScfService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
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
