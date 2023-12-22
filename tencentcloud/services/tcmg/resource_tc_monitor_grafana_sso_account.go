package tcmg

import (
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svcmonitor "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/monitor"

	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	monitor "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/monitor/v20180724"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudMonitorGrafanaSsoAccount() *schema.Resource {
	return &schema.Resource{
		Read:   resourceTencentCloudMonitorGrafanaSsoAccountRead,
		Create: resourceTencentCloudMonitorGrafanaSsoAccountCreate,
		Update: resourceTencentCloudMonitorGrafanaSsoAccountUpdate,
		Delete: resourceTencentCloudMonitorGrafanaSsoAccountDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "grafana instance id.",
			},

			"user_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "sub account uin of specific user.",
			},

			"notes": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "account related description.",
			},

			"role": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Description: "grafana role.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"organization": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Grafana organization id string.",
						},
						"role": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Grafana role, one of {Admin,Editor,Viewer}.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudMonitorGrafanaSsoAccountCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_monitor_grafana_sso_account.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request    = monitor.NewCreateSSOAccountRequest()
		response   *monitor.CreateSSOAccountResponse
		instanceId string
		userId     string
	)

	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("user_id"); ok {
		request.UserId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("notes"); ok {
		request.Notes = helper.String(v.(string))
	}

	if v, ok := d.GetOk("role"); ok {
		roleList := v.([]interface{})
		for _, r := range roleList {
			rr := r.(map[string]interface{})
			var role monitor.GrafanaAccountRole
			if vv, ok := rr["role"]; ok {
				role.Role = helper.String(vv.(string))
			}
			if vv, ok := rr["organization"]; ok {
				role.Organization = helper.String(vv.(string))
			}
			request.Role = append(request.Role, &role)
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMonitorClient().CreateSSOAccount(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create monitor ssoAccount failed, reason:%+v", logId, err)
		return err
	}

	userId = *response.Response.UserId

	d.SetId(strings.Join([]string{instanceId, userId}, tccommon.FILED_SP))
	return resourceTencentCloudMonitorGrafanaSsoAccountRead(d, meta)
}

func resourceTencentCloudMonitorGrafanaSsoAccountRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_monitor_grafana_sso_account.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := svcmonitor.NewMonitorService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	instanceId := idSplit[0]
	userId := idSplit[1]

	ssoAccount, err := service.DescribeMonitorSsoAccount(ctx, instanceId, userId)

	if err != nil {
		return err
	}

	if ssoAccount == nil {
		d.SetId("")
		return fmt.Errorf("resource `ssoAccount` %s does not exist", userId)
	}

	_ = d.Set("instance_id", instanceId)

	if ssoAccount.UserId != nil {
		_ = d.Set("user_id", ssoAccount.UserId)
	}

	if ssoAccount.Notes != nil {
		_ = d.Set("notes", ssoAccount.Notes)
	}

	if ssoAccount.Role != nil {
		roleList := make([]map[string]interface{}, 0, len(ssoAccount.Role))
		for _, role := range ssoAccount.Role {
			roleList = append(roleList, map[string]interface{}{
				"role":         role.Role,
				"organization": role.Organization,
			})
		}
		_ = d.Set("role", roleList)
	}

	return nil
}

func resourceTencentCloudMonitorGrafanaSsoAccountUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_monitor_grafana_sso_account.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	request := monitor.NewUpdateSSOAccountRequest()

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	instanceId := idSplit[0]
	userId := idSplit[1]

	request.InstanceId = &instanceId
	request.UserId = &userId

	if d.HasChange("instance_id") {
		return fmt.Errorf("`instance_id` do not support change now.")
	}

	if d.HasChange("user_id") {
		return fmt.Errorf("`user_id` do not support change now.")
	}

	if d.HasChange("notes") {
		if v, ok := d.GetOk("notes"); ok {
			request.Notes = helper.String(v.(string))
		}
	}

	if d.HasChange("role") {
		if v, ok := d.GetOk("role"); ok {
			roleList := v.([]interface{})
			for _, r := range roleList {
				rr := r.(map[string]interface{})
				var role monitor.GrafanaAccountRole
				if vv, ok := rr["role"]; ok {
					role.Role = helper.String(vv.(string))
				}
				if vv, ok := rr["organization"]; ok {
					role.Organization = helper.String(vv.(string))
				}
				request.Role = append(request.Role, &role)
			}
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMonitorClient().UpdateSSOAccount(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})

	if err != nil {
		return err
	}

	return resourceTencentCloudMonitorGrafanaSsoAccountRead(d, meta)
}

func resourceTencentCloudMonitorGrafanaSsoAccountDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_monitor_grafana_sso_account.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := svcmonitor.NewMonitorService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	instanceId := idSplit[0]
	userId := idSplit[1]

	if err := service.DeleteMonitorSsoAccountById(ctx, instanceId, userId); err != nil {
		return err
	}

	return nil
}
