package sqlserver

import (
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sqlserver "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sqlserver/v20180328"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudSqlserverStartXevent() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudSqlserverStartXeventCreate,
		Read:   resourceTencentCloudSqlserverStartXeventRead,
		Delete: resourceTencentCloudSqlserverStartXeventDelete,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Instance ID.",
			},
			"event_config": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				Description: "Whether to start or stop an extended event.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"event_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Event type. Valid values: slow (set threshold for slow SQL ), blocked (set threshold for the blocking and deadlock).",
						},
						"threshold": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Threshold in milliseconds. Valid values: 0(disable), non-zero (enable).",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudSqlserverStartXeventCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_sqlserver_start_xevent.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		request    = sqlserver.NewStartInstanceXEventRequest()
		instanceId string
	)
	if v, ok := d.GetOk("instance_id"); ok {
		request.InstanceId = helper.String(v.(string))
		instanceId = v.(string)
	}

	if v, ok := d.GetOk("event_config"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			eventConfig := sqlserver.EventConfig{}
			if v, ok := dMap["event_type"]; ok {
				eventConfig.EventType = helper.String(v.(string))
			}
			if v, ok := dMap["threshold"]; ok {
				eventConfig.Threshold = helper.IntInt64(v.(int))
			}
			request.EventConfig = append(request.EventConfig, &eventConfig)
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseSqlserverClient().StartInstanceXEvent(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s operate sqlserver startXevent failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(instanceId)

	return resourceTencentCloudSqlserverStartXeventRead(d, meta)
}

func resourceTencentCloudSqlserverStartXeventRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_sqlserver_start_xevent.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudSqlserverStartXeventDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_sqlserver_start_xevent.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
