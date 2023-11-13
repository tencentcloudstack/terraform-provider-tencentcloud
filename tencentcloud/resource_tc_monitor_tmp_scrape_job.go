/*
Provides a resource to create a monitor tmp_scrape_job

Example Usage

```hcl
resource "tencentcloud_monitor_tmp_scrape_job" "tmp_scrape_job" {
  instance_id = "prom-dko9d0nu"
  agent_id = "agent-6a7g40k2"
  config = "job_name: demo-config"
}
```

Import

monitor tmp_scrape_job can be imported using the id, e.g.

```
terraform import tencentcloud_monitor_tmp_scrape_job.tmp_scrape_job tmp_scrape_job_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	monitor "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/monitor/v20180724"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudMonitorTmpScrapeJob() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMonitorTmpScrapeJobCreate,
		Read:   resourceTencentCloudMonitorTmpScrapeJobRead,
		Update: resourceTencentCloudMonitorTmpScrapeJobUpdate,
		Delete: resourceTencentCloudMonitorTmpScrapeJobDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance id.",
			},

			"agent_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Agent id.",
			},

			"config": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Job content.",
			},
		},
	}
}

func resourceTencentCloudMonitorTmpScrapeJobCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_monitor_tmp_scrape_job.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request  = monitor.NewCreatePrometheusScrapeJobRequest()
		response = monitor.NewCreatePrometheusScrapeJobResponse()
		jobId    string
	)
	if v, ok := d.GetOk("instance_id"); ok {
		request.InstanceId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("agent_id"); ok {
		request.AgentId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("config"); ok {
		request.Config = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseMonitorClient().CreatePrometheusScrapeJob(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create monitor tmpScrapeJob failed, reason:%+v", logId, err)
		return err
	}

	jobId = *response.Response.JobId
	d.SetId(jobId)

	return resourceTencentCloudMonitorTmpScrapeJobRead(d, meta)
}

func resourceTencentCloudMonitorTmpScrapeJobRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_monitor_tmp_scrape_job.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := MonitorService{client: meta.(*TencentCloudClient).apiV3Conn}

	tmpScrapeJobId := d.Id()

	tmpScrapeJob, err := service.DescribeMonitorTmpScrapeJobById(ctx, jobId)
	if err != nil {
		return err
	}

	if tmpScrapeJob == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `MonitorTmpScrapeJob` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if tmpScrapeJob.InstanceId != nil {
		_ = d.Set("instance_id", tmpScrapeJob.InstanceId)
	}

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

	tmpScrapeJobId := d.Id()

	request.JobId = &jobId

	immutableArgs := []string{"instance_id", "agent_id", "config"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
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
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update monitor tmpScrapeJob failed, reason:%+v", logId, err)
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

	if err := service.DeleteMonitorTmpScrapeJobById(ctx, jobId); err != nil {
		return err
	}

	return nil
}
