package scf

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	scf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/scf/v20180416"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudScfTrigger() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudScfTriggerCreate,
		Read:   resourceTencentCloudScfTriggerRead,
		Update: resourceTencentCloudScfTriggerUpdate,
		Delete: resourceTencentCloudScfTriggerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"function_name": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Name of the SCF function that the trigger binds to.",
			},

			"trigger_name": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Name of the trigger.",
			},

			"type": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Trigger type. Valid values: `cos`, `cls`, `timer`, `ckafka`, `http`.",
			},

			"trigger_desc": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Trigger description parameter, see the trigger description documentation for details.",
			},

			"namespace": {
				Optional:    true,
				ForceNew:    true,
				Default:     "default",
				Type:        schema.TypeString,
				Description: "Function namespace. Defaults to `default`.",
			},

			"qualifier": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Function version or alias that the trigger points to. Defaults to `$LATEST`.",
			},

			"enable": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Trigger enable status. Valid values: `OPEN` (enabled), `CLOSE` (disabled).",
			},

			"custom_argument": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "User custom parameter, only supported by timer trigger.",
			},

			"description": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Trigger description.",
			},

			"available_status": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Trigger available status.",
			},

			"add_time": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Trigger creation time.",
			},

			"mod_time": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Trigger last modified time.",
			},
		},
	}
}

func resourceTencentCloudScfTriggerCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_scf_trigger.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId        = tccommon.GetLogId(tccommon.ContextNil)
		ctx          = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request      = scf.NewCreateTriggerRequest()
		functionName string
		triggerName  string
		namespace    string
	)

	if v, ok := d.GetOk("function_name"); ok {
		request.FunctionName = helper.String(v.(string))
		functionName = v.(string)
	}

	if v, ok := d.GetOk("trigger_name"); ok {
		request.TriggerName = helper.String(v.(string))
		triggerName = v.(string)
	}

	if v, ok := d.GetOk("type"); ok {
		request.Type = helper.String(v.(string))
	}

	if v, ok := d.GetOk("trigger_desc"); ok {
		request.TriggerDesc = helper.String(v.(string))
	}

	if v, ok := d.GetOk("namespace"); ok {
		request.Namespace = helper.String(v.(string))
		namespace = v.(string)
	}

	if v, ok := d.GetOk("qualifier"); ok {
		request.Qualifier = helper.String(v.(string))
	}

	if v, ok := d.GetOk("enable"); ok {
		request.Enable = helper.String(v.(string))
	}

	if v, ok := d.GetOk("custom_argument"); ok {
		request.CustomArgument = helper.String(v.(string))
	}

	if v, ok := d.GetOk("description"); ok {
		request.Description = helper.String(v.(string))
	}

	var response *scf.CreateTriggerResponse
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseScfClient().CreateTriggerWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create scf_trigger failed, Response is nil."))
		}

		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create scf_trigger failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	if response.Response.TriggerInfo == nil {
		log.Printf("[CRITAL]%s create scf_trigger failed, TriggerInfo is nil, d.Id()=%s", logId, d.Id())
		return fmt.Errorf("create scf_trigger failed, TriggerInfo is nil.")
	}

	d.SetId(functionName + tccommon.FILED_SP + namespace + tccommon.FILED_SP + triggerName)

	return resourceTencentCloudScfTriggerRead(d, meta)
}

func resourceTencentCloudScfTriggerRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_scf_trigger.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = ScfService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	functionName := idSplit[0]
	namespace := idSplit[1]
	triggerName := idSplit[2]

	var triggerInfo *scf.TriggerInfo
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeScfTriggerById(ctx, functionName, namespace, triggerName)
		if e != nil {
			return tccommon.RetryError(e)
		}
		triggerInfo = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read scf_trigger failed, reason:%+v", logId, err)
		return err
	}

	if triggerInfo == nil {
		log.Printf("[CRUD] scf_trigger id=%s", d.Id())
		d.SetId("")
		return nil
	}

	_ = d.Set("function_name", functionName)
	_ = d.Set("namespace", namespace)

	if triggerInfo.TriggerName != nil {
		_ = d.Set("trigger_name", triggerInfo.TriggerName)
	}

	if triggerInfo.Type != nil {
		_ = d.Set("type", triggerInfo.Type)
	}

	if triggerInfo.TriggerDesc != nil {
		_ = d.Set("trigger_desc", triggerInfo.TriggerDesc)
	}

	if triggerInfo.Qualifier != nil {
		_ = d.Set("qualifier", triggerInfo.Qualifier)
	}

	if triggerInfo.CustomArgument != nil {
		_ = d.Set("custom_argument", triggerInfo.CustomArgument)
	}

	if triggerInfo.Description != nil {
		_ = d.Set("description", triggerInfo.Description)
	}

	if triggerInfo.Enable != nil {
		if *triggerInfo.Enable == 1 {
			_ = d.Set("enable", "OPEN")
		} else {
			_ = d.Set("enable", "CLOSE")
		}
	}

	if triggerInfo.AvailableStatus != nil {
		_ = d.Set("available_status", triggerInfo.AvailableStatus)
	}

	if triggerInfo.AddTime != nil {
		_ = d.Set("add_time", triggerInfo.AddTime)
	}

	if triggerInfo.ModTime != nil {
		_ = d.Set("mod_time", triggerInfo.ModTime)
	}

	return nil
}

func resourceTencentCloudScfTriggerUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_scf_trigger.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
		ctx   = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	functionName := idSplit[0]
	namespace := idSplit[1]
	triggerName := idSplit[2]

	needChange := false
	mutableArgs := []string{"enable", "qualifier", "trigger_desc", "description", "custom_argument"}
	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {
		request := scf.NewUpdateTriggerRequest()
		request.FunctionName = &functionName
		request.Namespace = &namespace
		request.TriggerName = &triggerName

		if v, ok := d.GetOk("type"); ok {
			request.Type = helper.String(v.(string))
		}

		if v, ok := d.GetOk("enable"); ok {
			request.Enable = helper.String(v.(string))
		}

		if v, ok := d.GetOk("qualifier"); ok {
			request.Qualifier = helper.String(v.(string))
		}

		if v, ok := d.GetOk("trigger_desc"); ok {
			request.TriggerDesc = helper.String(v.(string))
		}

		if v, ok := d.GetOk("description"); ok {
			request.Description = helper.String(v.(string))
		}

		if v, ok := d.GetOk("custom_argument"); ok {
			request.CustomArgument = helper.String(v.(string))
		}

		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseScfClient().UpdateTriggerWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s update scf_trigger failed, reason:%+v", logId, reqErr)
			return reqErr
		}
	}

	return resourceTencentCloudScfTriggerRead(d, meta)
}

func resourceTencentCloudScfTriggerDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_scf_trigger.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
		ctx   = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	functionName := idSplit[0]
	namespace := idSplit[1]
	triggerName := idSplit[2]

	request := scf.NewDeleteTriggerRequest()
	request.FunctionName = &functionName
	request.TriggerName = &triggerName
	request.Namespace = &namespace

	if v, ok := d.GetOk("type"); ok {
		request.Type = helper.String(v.(string))
	}

	if v, ok := d.GetOk("trigger_desc"); ok {
		request.TriggerDesc = helper.String(v.(string))
	}

	if v, ok := d.GetOk("qualifier"); ok {
		request.Qualifier = helper.String(v.(string))
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseScfClient().DeleteTriggerWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete scf_trigger failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return nil
}
