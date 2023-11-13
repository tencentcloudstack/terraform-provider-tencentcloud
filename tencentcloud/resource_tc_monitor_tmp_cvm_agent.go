/*
Provides a resource to create a monitor tmp_cvm_agent

Example Usage

```hcl
resource "tencentcloud_monitor_tmp_cvm_agent" "tmp_cvm_agent" {
  instance_id = "prom-dko9d0nu"
  name = "agent"
}
```

Import

monitor tmp_cvm_agent can be imported using the id, e.g.

```
terraform import tencentcloud_monitor_tmp_cvm_agent.tmp_cvm_agent tmp_cvm_agent_id
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

func resourceTencentCloudMonitorTmpCvmAgent() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMonitorTmpCvmAgentCreate,
		Read:   resourceTencentCloudMonitorTmpCvmAgentRead,
		Update: resourceTencentCloudMonitorTmpCvmAgentUpdate,
		Delete: resourceTencentCloudMonitorTmpCvmAgentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance id.",
			},

			"name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Agent name.",
			},
		},
	}
}

func resourceTencentCloudMonitorTmpCvmAgentCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_monitor_tmp_cvm_agent.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request  = monitor.NewCreatePrometheusAgentRequest()
		response = monitor.NewCreatePrometheusAgentResponse()
		agentId  string
	)
	if v, ok := d.GetOk("instance_id"); ok {
		request.InstanceId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseMonitorClient().CreatePrometheusAgent(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create monitor tmpCvmAgent failed, reason:%+v", logId, err)
		return err
	}

	agentId = *response.Response.AgentId
	d.SetId(agentId)

	return resourceTencentCloudMonitorTmpCvmAgentRead(d, meta)
}

func resourceTencentCloudMonitorTmpCvmAgentRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_monitor_tmp_cvm_agent.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := MonitorService{client: meta.(*TencentCloudClient).apiV3Conn}

	tmpCvmAgentId := d.Id()

	tmpCvmAgent, err := service.DescribeMonitorTmpCvmAgentById(ctx, agentId)
	if err != nil {
		return err
	}

	if tmpCvmAgent == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `MonitorTmpCvmAgent` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if tmpCvmAgent.InstanceId != nil {
		_ = d.Set("instance_id", tmpCvmAgent.InstanceId)
	}

	if tmpCvmAgent.Name != nil {
		_ = d.Set("name", tmpCvmAgent.Name)
	}

	return nil
}

func resourceTencentCloudMonitorTmpCvmAgentUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_monitor_tmp_cvm_agent.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := monitor.NewRequest()

	tmpCvmAgentId := d.Id()

	request.AgentId = &agentId

	immutableArgs := []string{"instance_id", "name"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseMonitorClient().(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update monitor tmpCvmAgent failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudMonitorTmpCvmAgentRead(d, meta)
}

func resourceTencentCloudMonitorTmpCvmAgentDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_monitor_tmp_cvm_agent.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := MonitorService{client: meta.(*TencentCloudClient).apiV3Conn}
	tmpCvmAgentId := d.Id()

	if err := service.DeleteMonitorTmpCvmAgentById(ctx, agentId); err != nil {
		return err
	}

	return nil
}
