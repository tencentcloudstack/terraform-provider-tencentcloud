package tmp

import (
	"context"
	"fmt"
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svcmonitor "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/monitor"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	monitor "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/monitor/v20180724"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudMonitorTmpScrapeJob() *schema.Resource {
	return &schema.Resource{
		Read:   resourceTencentCloudMonitorTmpScrapeJobRead,
		Create: resourceTencentCloudMonitorTmpScrapeJobCreate,
		Update: resourceTencentCloudMonitorTmpScrapeJobUpdate,
		Delete: resourceTencentCloudMonitorTmpScrapeJobDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Instance id.",
			},

			"agent_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Agent id.",
			},

			"config": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Job content.",
			},
		},
	}
}

func resourceTencentCloudMonitorTmpScrapeJobCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_monitor_tmp_scrape_job.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var instanceId string
	var agentId string

	var (
		request  = monitor.NewCreatePrometheusScrapeJobRequest()
		response *monitor.CreatePrometheusScrapeJobResponse
	)

	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceId = helper.String(instanceId)
	}

	if v, ok := d.GetOk("agent_id"); ok {
		agentId = v.(string)
		request.AgentId = helper.String(agentId)
	}

	if v, ok := d.GetOk("config"); ok {
		request.Config = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMonitorClient().CreatePrometheusScrapeJob(request)
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
		log.Printf("[CRITAL]%s create monitor tmpScrapeJob failed, reason:%+v", logId, err)
		return err
	}

	tmpScrapeJobId := *response.Response.JobId

	d.SetId(strings.Join([]string{tmpScrapeJobId, instanceId, agentId}, tccommon.FILED_SP))

	return resourceTencentCloudMonitorTmpScrapeJobRead(d, meta)
}

func resourceTencentCloudMonitorTmpScrapeJobRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_monitor_tmpScrapeJob.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := svcmonitor.NewMonitorService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())

	tmpScrapeJobId := d.Id()

	tmpScrapeJob, err := service.DescribeMonitorTmpScrapeJob(ctx, tmpScrapeJobId)

	if err != nil {
		return err
	}

	if tmpScrapeJob == nil {
		d.SetId("")
		return fmt.Errorf("resource `tmpScrapeJob` %s does not exist", tmpScrapeJobId)
	}

	_ = d.Set("instance_id", strings.Split(tmpScrapeJobId, tccommon.FILED_SP)[1])
	if tmpScrapeJob.AgentId != nil {
		_ = d.Set("agent_id", tmpScrapeJob.AgentId)
	}

	if tmpScrapeJob.Config != nil {
		_ = d.Set("config", tmpScrapeJob.Config)
	}

	return nil
}

func resourceTencentCloudMonitorTmpScrapeJobUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_monitor_tmp_scrape_job.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	request := monitor.NewUpdatePrometheusScrapeJobRequest()

	ids := strings.Split(d.Id(), tccommon.FILED_SP)

	request.JobId = &ids[0]
	request.InstanceId = &ids[1]
	request.AgentId = &ids[2]

	if d.HasChange("instance_id") {
		return fmt.Errorf("`instance_id` do not support change now.")
	}

	if d.HasChange("agent_id") {
		return fmt.Errorf("`agent_id` do not support change now.")
	}

	if d.HasChange("config") {
		if v, ok := d.GetOk("config"); ok {
			request.Config = helper.String(v.(string))
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMonitorClient().UpdatePrometheusScrapeJob(request)
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

	return resourceTencentCloudMonitorTmpScrapeJobRead(d, meta)
}

func resourceTencentCloudMonitorTmpScrapeJobDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_monitor_tmp_scrape_job.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := svcmonitor.NewMonitorService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
	tmpScrapeJobId := d.Id()

	if err := service.DeleteMonitorTmpScrapeJobById(ctx, tmpScrapeJobId); err != nil {
		return err
	}

	return nil
}
