/*
Provides a resource to create a monitor tmpScrapeJob

Example Usage

```hcl
resource "tencentcloud_monitor_tmp_scrape_job" "tmpScrapeJob" {
  instance_id = "prom-dko9d0nu"
  agent_id = "agent-6a7g40k2"
  config = <<-EOT
job_name: demo-config
honor_timestamps: true
metrics_path: /metrics
scheme: https
EOT
}

```
Import

monitor tmpScrapeJob can be imported using the id, e.g.
```
$ terraform import tencentcloud_monitor_tmp_scrape_job.tmpScrapeJob tmpScrapeJob_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	monitor "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/monitor/v20180724"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudMonitorTmpScrapeJob() *schema.Resource {
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
	defer logElapsed("resource.tencentcloud_monitor_tmp_scrape_job.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

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

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseMonitorClient().CreatePrometheusScrapeJob(request)
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
		log.Printf("[CRITAL]%s create monitor tmpScrapeJob failed, reason:%+v", logId, err)
		return err
	}

	tmpScrapeJobId := *response.Response.JobId

	d.SetId(strings.Join([]string{tmpScrapeJobId, instanceId, agentId}, FILED_SP))

	return resourceTencentCloudMonitorTmpScrapeJobRead(d, meta)
}

func resourceTencentCloudMonitorTmpScrapeJobRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_monitor_tmpScrapeJob.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := MonitorService{client: meta.(*TencentCloudClient).apiV3Conn}

	tmpScrapeJobId := d.Id()

	tmpScrapeJob, err := service.DescribeMonitorTmpScrapeJob(ctx, tmpScrapeJobId)

	if err != nil {
		return err
	}

	if tmpScrapeJob == nil {
		d.SetId("")
		return fmt.Errorf("resource `tmpScrapeJob` %s does not exist", tmpScrapeJobId)
	}

	_ = d.Set("instance_id", strings.Split(tmpScrapeJobId, FILED_SP)[1])
	if tmpScrapeJob.AgentId != nil {
		_ = d.Set("agent_id", tmpScrapeJob.AgentId)
	}

	if tmpScrapeJob.Config != nil {
		_ = d.Set("config", tmpScrapeJob.Config)
	}

	return nil
}

func resourceTencentCloudMonitorTmpScrapeJobUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_monitor_tmp_scrape_job.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := monitor.NewUpdatePrometheusScrapeJobRequest()

	ids := strings.Split(d.Id(), FILED_SP)

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

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseMonitorClient().UpdatePrometheusScrapeJob(request)
		if e != nil {
			return retryError(e)
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
	defer logElapsed("resource.tencentcloud_monitor_tmp_scrape_job.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := MonitorService{client: meta.(*TencentCloudClient).apiV3Conn}
	tmpScrapeJobId := d.Id()

	if err := service.DeleteMonitorTmpScrapeJobById(ctx, tmpScrapeJobId); err != nil {
		return err
	}

	return nil
}
