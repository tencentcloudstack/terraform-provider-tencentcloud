package cls

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cls "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cls/v20201016"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudClsRemoteWriteTask() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudClsRemoteWriteTaskCreate,
		Read:   resourceTencentCloudClsRemoteWriteTaskRead,
		Update: resourceTencentCloudClsRemoteWriteTaskUpdate,
		Delete: resourceTencentCloudClsRemoteWriteTaskDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"topic_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Log topic ID.",
			},

			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "RemoteWrite task name.",
			},

			"target": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Target service name.",
			},

			"remote_write_url": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Target address for RemoteWrite.",
			},

			"auth_type": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Authentication type. 0: no auth, 1: basic_auth, 2: token.",
			},

			"net_type": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Network type. 1: intranet, 2: internet.",
			},

			"vpc_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Private network ID.",
			},

			"virtual_gateway_type": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Backend service type. 0: CVM, 1025: CLB.",
			},

			"enable": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Task status. 0: disabled, 1: enabled.",
			},

			"auth_info": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Authentication information block.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"username": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Basic auth username.",
						},
						"password": {
							Type:        schema.TypeString,
							Optional:    true,
							Sensitive:   true,
							Description: "Basic auth password.",
						},
						"token": {
							Type:        schema.TypeString,
							Optional:    true,
							Sensitive:   true,
							Description: "Basic auth token.",
						},
					},
				},
			},

			// computed
			"task_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "RemoteWrite task ID.",
			},

			"status": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Task running status. 1: running, 2: paused, 3: failed.",
			},

			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Task creation time.",
			},

			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Task update time.",
			},

			"logset_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Logset ID.",
			},
		},
	}
}

func resourceTencentCloudClsRemoteWriteTaskCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cls_remote_write_task.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId    = tccommon.GetLogId(tccommon.ContextNil)
		ctx      = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request  = cls.NewCreateRemoteWriteTaskRequest()
		response = cls.NewCreateRemoteWriteTaskResponse()
		topicId  string
	)

	if v, ok := d.GetOk("topic_id"); ok {
		request.TopicId = helper.String(v.(string))
		topicId = v.(string)
	}

	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOk("target"); ok {
		request.Target = helper.String(v.(string))
	}

	if v, ok := d.GetOk("remote_write_url"); ok {
		request.RemoteWriteURL = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("auth_type"); ok {
		request.AuthType = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("net_type"); ok {
		request.NetType = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("vpc_id"); ok {
		request.VpcId = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("virtual_gateway_type"); ok {
		request.VirtualGatewayType = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("auth_info"); ok {
		authInfoList := v.([]interface{})
		if len(authInfoList) > 0 {
			authInfoMap := authInfoList[0].(map[string]interface{})
			authInfo := cls.RemoteWriteAuthInfo{}
			if v, ok := authInfoMap["username"].(string); ok && v != "" {
				authInfo.Username = helper.String(v)
			}
			if v, ok := authInfoMap["password"].(string); ok && v != "" {
				authInfo.Password = helper.String(v)
			}
			if v, ok := authInfoMap["token"].(string); ok && v != "" {
				authInfo.Token = helper.String(v)
			}
			request.AuthInfo = &authInfo
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClsClient().CreateRemoteWriteTaskWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create cls remote_write_task failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create cls remote_write_task failed, reason:%+v", logId, err)
		return err
	}

	log.Printf("[CRUD] cls remote_write_task logId=%s, topicId=%s, taskId=%s", logId, topicId, *response.Response.TaskId)

	if response.Response.TaskId == nil || *response.Response.TaskId == "" {
		return fmt.Errorf("Create cls remote_write_task failed, TaskId is nil or empty.")
	}

	taskId := *response.Response.TaskId
	d.SetId(strings.Join([]string{topicId, taskId}, tccommon.FILED_SP))
	return resourceTencentCloudClsRemoteWriteTaskRead(d, meta)
}

func resourceTencentCloudClsRemoteWriteTaskRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cls_remote_write_task.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
		ctx   = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	taskId := idSplit[1]

	request := cls.NewDescribeRemoteWriteTasksRequest()
	request.Offset = helper.Uint64(uint64(0))
	request.Limit = helper.Uint64(uint64(100))
	request.Filters = []*cls.Filter{
		{
			Key:    helper.String("taskId"),
			Values: []*string{helper.String(taskId)},
		},
	}

	var infos []*cls.RemoteWriteInfo
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClsClient().DescribeRemoteWriteTasksWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil || len(result.Response.Infos) == 0 {
			return nil
		}

		infos = result.Response.Infos
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s read cls remote_write_task failed, reason:%+v", logId, err)
		return err
	}

	if len(infos) == 0 {
		log.Printf("[CRUD] tencentcloud_cls_remote_write_task id=%s", d.Id())
		d.SetId("")
		return nil
	}

	info := infos[0]

	if info.TopicId != nil {
		_ = d.Set("topic_id", info.TopicId)
	}

	if info.Name != nil {
		_ = d.Set("name", info.Name)
	}

	if info.Target != nil {
		_ = d.Set("target", info.Target)
	}

	if info.RemoteWriteURL != nil {
		_ = d.Set("remote_write_url", info.RemoteWriteURL)
	}

	if info.AuthType != nil {
		_ = d.Set("auth_type", info.AuthType)
	}

	if info.NetType != nil {
		_ = d.Set("net_type", info.NetType)
	}

	if info.VpcId != nil {
		_ = d.Set("vpc_id", info.VpcId)
	}

	if info.VirtualGatewayType != nil {
		_ = d.Set("virtual_gateway_type", info.VirtualGatewayType)
	}

	if info.Enable != nil {
		_ = d.Set("enable", info.Enable)
	}

	if info.AuthInfo != nil {
		authInfoMap := map[string]interface{}{}
		if info.AuthInfo.Username != nil {
			authInfoMap["username"] = info.AuthInfo.Username
		}
		if info.AuthInfo.Password != nil {
			authInfoMap["password"] = info.AuthInfo.Password
		}
		if info.AuthInfo.Token != nil {
			authInfoMap["token"] = info.AuthInfo.Token
		}
		_ = d.Set("auth_info", []interface{}{authInfoMap})
	}

	if info.TaskId != nil {
		_ = d.Set("task_id", info.TaskId)
	}

	if info.Status != nil {
		_ = d.Set("status", info.Status)
	}

	if info.CreateTime != nil {
		_ = d.Set("create_time", info.CreateTime)
	}

	if info.UpdateTime != nil {
		_ = d.Set("update_time", info.UpdateTime)
	}

	if info.LogsetId != nil {
		_ = d.Set("logset_id", info.LogsetId)
	}

	return nil
}

func resourceTencentCloudClsRemoteWriteTaskUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cls_remote_write_task.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
		ctx   = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	topicId := idSplit[0]
	taskId := idSplit[1]

	needChange := false
	mutableArgs := []string{"name", "net_type", "vpc_id", "target", "remote_write_url", "auth_type", "enable", "virtual_gateway_type", "auth_info"}
	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {
		request := cls.NewModifyRemoteWriteTaskRequest()
		request.TaskId = helper.String(taskId)
		request.TopicId = helper.String(topicId)

		if v, ok := d.GetOk("name"); ok {
			request.Name = helper.String(v.(string))
		}

		if v, ok := d.GetOk("net_type"); ok {
			request.NetType = helper.IntUint64(v.(int))
		}

		if v, ok := d.GetOk("vpc_id"); ok {
			request.VpcId = helper.String(v.(string))
		}

		if v, ok := d.GetOk("target"); ok {
			request.Target = helper.String(v.(string))
		}

		if v, ok := d.GetOk("remote_write_url"); ok {
			request.RemoteWriteURL = helper.String(v.(string))
		}

		if v, ok := d.GetOk("auth_type"); ok {
			request.AuthType = helper.IntUint64(v.(int))
		}

		if v, ok := d.GetOkExists("enable"); ok {
			request.Enable = helper.IntUint64(v.(int))
		}

		if v, ok := d.GetOkExists("virtual_gateway_type"); ok {
			request.VirtualGatewayType = helper.IntInt64(v.(int))
		}

		if v, ok := d.GetOk("auth_info"); ok {
			authInfoList := v.([]interface{})
			if len(authInfoList) > 0 {
				authInfoMap := authInfoList[0].(map[string]interface{})
				authInfo := cls.RemoteWriteAuthInfo{}
				if v, ok := authInfoMap["username"].(string); ok && v != "" {
					authInfo.Username = helper.String(v)
				}
				if v, ok := authInfoMap["password"].(string); ok && v != "" {
					authInfo.Password = helper.String(v)
				}
				if v, ok := authInfoMap["token"].(string); ok && v != "" {
					authInfo.Token = helper.String(v)
				}
				request.AuthInfo = &authInfo
			}
		}

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClsClient().ModifyRemoteWriteTaskWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			if result == nil || result.Response == nil {
				return resource.NonRetryableError(fmt.Errorf("Modify cls remote_write_task failed, Response is nil."))
			}

			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s update cls remote_write_task failed, reason:%+v", logId, err)
			return err
		}
	}

	return resourceTencentCloudClsRemoteWriteTaskRead(d, meta)
}

func resourceTencentCloudClsRemoteWriteTaskDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cls_remote_write_task.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
		ctx   = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	topicId := idSplit[0]
	taskId := idSplit[1]

	request := cls.NewDeleteRemoteWriteTaskRequest()
	request.TaskId = helper.String(taskId)
	request.TopicId = helper.String(topicId)

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClsClient().DeleteRemoteWriteTaskWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Delete cls remote_write_task failed, Response is nil."))
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s delete cls remote_write_task failed, reason:%+v", logId, err)
		return err
	}

	return nil
}
