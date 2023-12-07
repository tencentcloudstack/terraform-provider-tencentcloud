package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tat "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tat/v20201028"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudTatInvoker() *schema.Resource {
	return &schema.Resource{
		Read:   resourceTencentCloudTatInvokerRead,
		Create: resourceTencentCloudTatInvokerCreate,
		Update: resourceTencentCloudTatInvokerUpdate,
		Delete: resourceTencentCloudTatInvokerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Invoker name.",
			},

			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Invoker type. It can only be `SCHEDULE` (recurring invokers).",
			},

			"command_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Remote command ID.",
			},

			"instance_ids": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Required:    true,
				Description: "ID of the instance bound to the trigger. Up to 100 IDs are allowed.",
			},

			"username": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The user who executes the command.",
			},

			"parameters": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Custom parameters of the command.",
			},

			"schedule_settings": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "Settings required for a recurring invoker.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"policy": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Execution policy: `ONCE`: Execute once; `RECURRENCE`: Execute repeatedly.",
						},
						"recurrence": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Trigger the crontab expression. This field is required if `Policy` is `RECURRENCE`. The crontab expression is parsed in UTC+8.",
						},
						"invoke_time": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The next execution time of the invoker. This field is required if Policy is ONCE.",
						},
					},
				},
			},

			"invoker_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Invoker ID.",
			},

			"enable": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether to enable the invoker.",
			},

			"created_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation time.",
			},

			"updated_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Modification time.",
			},
		},
	}
}

func resourceTencentCloudTatInvokerCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tat_invoker.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request   = tat.NewCreateInvokerRequest()
		response  *tat.CreateInvokerResponse
		invokerId string
	)

	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOk("type"); ok {
		request.Type = helper.String(v.(string))
	}

	if v, ok := d.GetOk("command_id"); ok {
		request.CommandId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("instance_ids"); ok {
		instanceIdsSet := v.(*schema.Set).List()
		for i := range instanceIdsSet {
			instanceIds := instanceIdsSet[i].(string)
			request.InstanceIds = append(request.InstanceIds, &instanceIds)
		}
	}

	if v, ok := d.GetOk("username"); ok {
		request.Username = helper.String(v.(string))
	}

	if v, ok := d.GetOk("parameters"); ok {
		request.Parameters = helper.String(v.(string))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "schedule_settings"); ok {
		scheduleSettings := tat.ScheduleSettings{}
		if v, ok := dMap["policy"]; ok {
			scheduleSettings.Policy = helper.String(v.(string))
		}
		if v, ok := dMap["recurrence"]; ok {
			scheduleSettings.Recurrence = helper.String(v.(string))
		}
		if v, ok := dMap["invoke_time"]; ok {
			scheduleSettings.InvokeTime = helper.String(v.(string))
		}

		request.ScheduleSettings = &scheduleSettings
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTatClient().CreateInvoker(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create tat invoker failed, reason:%+v", logId, err)
		return err
	}

	invokerId = *response.Response.InvokerId

	d.SetId(invokerId)
	return resourceTencentCloudTatInvokerRead(d, meta)
}

func resourceTencentCloudTatInvokerRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tat_invoker.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TatService{client: meta.(*TencentCloudClient).apiV3Conn}

	invokerId := d.Id()

	invoker, err := service.DescribeTatInvoker(ctx, invokerId)

	if err != nil {
		return err
	}

	if invoker == nil {
		d.SetId("")
		return fmt.Errorf("resource `invoker` %s does not exist", invokerId)
	}

	if invoker.Name != nil {
		_ = d.Set("name", invoker.Name)
	}

	if invoker.Type != nil {
		_ = d.Set("type", invoker.Type)
	}

	if invoker.CommandId != nil {
		_ = d.Set("command_id", invoker.CommandId)
	}

	if invoker.InstanceIds != nil {
		_ = d.Set("instance_ids", invoker.InstanceIds)
	}

	if invoker.Username != nil {
		_ = d.Set("username", invoker.Username)
	}

	if invoker.Parameters != nil {
		_ = d.Set("parameters", invoker.Parameters)
	}

	if invoker.ScheduleSettings != nil {
		scheduleSettingsMap := map[string]interface{}{}
		if invoker.ScheduleSettings.Policy != nil {
			scheduleSettingsMap["policy"] = invoker.ScheduleSettings.Policy
		}
		if invoker.ScheduleSettings.Recurrence != nil {
			scheduleSettingsMap["recurrence"] = invoker.ScheduleSettings.Recurrence
		}
		if invoker.ScheduleSettings.InvokeTime != nil {
			scheduleSettingsMap["invoke_time"] = invoker.ScheduleSettings.InvokeTime
		}

		_ = d.Set("schedule_settings", []interface{}{scheduleSettingsMap})
	}

	if invoker.InvokerId != nil {
		_ = d.Set("invoker_id", invoker.InvokerId)
	}

	if invoker.Enable != nil {
		_ = d.Set("enable", invoker.Enable)
	}

	if invoker.CreatedTime != nil {
		_ = d.Set("created_time", invoker.CreatedTime)
	}

	if invoker.UpdatedTime != nil {
		_ = d.Set("updated_time", invoker.UpdatedTime)
	}

	return nil
}

func resourceTencentCloudTatInvokerUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tat_invoker.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := tat.NewModifyInvokerRequest()

	invokerId := d.Id()

	request.InvokerId = &invokerId

	if d.HasChange("name") {
		if v, ok := d.GetOk("name"); ok {
			request.Name = helper.String(v.(string))
		}
	}

	if d.HasChange("type") {
		if v, ok := d.GetOk("type"); ok {
			request.Type = helper.String(v.(string))
		}
	}

	if d.HasChange("command_id") {
		if v, ok := d.GetOk("command_id"); ok {
			request.CommandId = helper.String(v.(string))
		}
	}

	if d.HasChange("instance_ids") {
		if v, ok := d.GetOk("instance_ids"); ok {
			instanceIdsSet := v.(*schema.Set).List()
			for i := range instanceIdsSet {
				instanceIds := instanceIdsSet[i].(string)
				request.InstanceIds = append(request.InstanceIds, &instanceIds)
			}
		}
	}

	if d.HasChange("username") {
		if v, ok := d.GetOk("username"); ok {
			request.Username = helper.String(v.(string))
		}
	}

	if d.HasChange("parameters") {
		if v, ok := d.GetOk("parameters"); ok {
			request.Parameters = helper.String(v.(string))
		}
	}

	if d.HasChange("schedule_settings") {
		if dMap, ok := helper.InterfacesHeadMap(d, "schedule_settings"); ok {
			scheduleSettings := tat.ScheduleSettings{}
			if v, ok := dMap["policy"]; ok {
				scheduleSettings.Policy = helper.String(v.(string))
			}
			if v, ok := dMap["recurrence"]; ok {
				scheduleSettings.Recurrence = helper.String(v.(string))
			}
			if v, ok := dMap["invoke_time"]; ok {
				scheduleSettings.InvokeTime = helper.String(v.(string))
			}

			request.ScheduleSettings = &scheduleSettings
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTatClient().ModifyInvoker(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create tat invoker failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudTatInvokerRead(d, meta)
}

func resourceTencentCloudTatInvokerDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tat_invoker.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TatService{client: meta.(*TencentCloudClient).apiV3Conn}

	invokerId := d.Id()

	if err := service.DeleteTatInvokerById(ctx, invokerId); err != nil {
		return err
	}

	return nil
}
