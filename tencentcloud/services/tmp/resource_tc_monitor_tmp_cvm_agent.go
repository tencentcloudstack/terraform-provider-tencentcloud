package tmp

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

func ResourceTencentCloudMonitorTmpCvmAgent() *schema.Resource {
	return &schema.Resource{
		Read:   resourceTencentCloudMonitorTmpCvmAgentRead,
		Create: resourceTencentCloudMonitorTmpCvmAgentCreate,
		//Update: resourceTencentCloudMonitorTmpCvmAgentUpdate,
		Delete: resourceTencentCloudMonitorTmpCvmAgentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Instance id.",
			},

			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Agent name.",
			},

			"agent_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Agent id.",
			},
		},
	}
}

func resourceTencentCloudMonitorTmpCvmAgentCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_monitor_tmp_cvm_agent.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request  = monitor.NewCreatePrometheusAgentRequest()
		response *monitor.CreatePrometheusAgentResponse
	)

	var instanceId string

	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceId = helper.String(instanceId)
	}

	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMonitorClient().CreatePrometheusAgent(request)
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
		log.Printf("[CRITAL]%s create monitor tmpCvmAgent failed, reason:%+v", logId, err)
		return err
	}

	tmpCvmAgentId := *response.Response.AgentId

	d.SetId(strings.Join([]string{instanceId, tmpCvmAgentId}, tccommon.FILED_SP))
	return resourceTencentCloudMonitorTmpCvmAgentRead(d, meta)
}

func resourceTencentCloudMonitorTmpCvmAgentRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_monitor_tmpCvmAgent.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := svcmonitor.NewMonitorService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())

	ids := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(ids) != 2 {
		return fmt.Errorf("id is broken, id is %s", d.Id())
	}

	tmpCvmAgent, err := service.DescribeMonitorTmpCvmAgent(ctx, ids[0], ids[1])

	if err != nil {
		return err
	}

	if tmpCvmAgent == nil {
		d.SetId("")
		return fmt.Errorf("resource `tmpCvmAgent` %s does not exist", ids[1])
	}

	if tmpCvmAgent.InstanceId != nil {
		_ = d.Set("instance_id", tmpCvmAgent.InstanceId)
	}

	if tmpCvmAgent.Name != nil {
		_ = d.Set("name", tmpCvmAgent.Name)
	}

	if tmpCvmAgent.AgentId != nil {
		_ = d.Set("agent_id", tmpCvmAgent.AgentId)
	}

	return nil
}

func resourceTencentCloudMonitorTmpCvmAgentDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_monitor_tmp_cvm_agent.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
