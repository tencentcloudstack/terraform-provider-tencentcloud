/*
Provides a resource to create a monitor tmpInstance

Example Usage

```hcl
resource "tencentcloud_monitor_tmp_instance" "tmpInstance" {
  instance_name = "logset-hello"
  vpc_id = "vpc-2hfyray3"
  subnet_id = "subnet-rdkj0agk"
  data_retention_time = 30
  zone = "ap-guangzhou-3"
}

```
Import

monitor tmp instance can be imported using the id, e.g.
```
$ terraform import tencentcloud_monitor_tmp_instance.tmpInstance tmpInstance_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	monitor "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/monitor/v20180724"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudMonitorTmpInstance() *schema.Resource {
	return &schema.Resource{
		Read:   resourceTencentCloudMonitorTmpInstanceRead,
		Create: resourceTencentCloudMonitorTmpInstanceCreate,
		Update: resourceTencentCloudMonitorTmpInstanceUpdate,
		Delete: resourceTencentCloudMonitorTmpInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Instance name.",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Vpc Id.",
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Subnet Id.",
			},
			"data_retention_time": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Data retention time.",
			},
			"zone": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Available zone.",
			},
			"grafana_instance_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Associated grafana instance id.",
			},
		},
	}
}

func resourceTencentCloudMonitorTmpInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_monitor_tmp_instance.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request  = monitor.NewCreatePrometheusMultiTenantInstancePostPayModeRequest()
		response *monitor.CreatePrometheusMultiTenantInstancePostPayModeResponse
	)

	if v, ok := d.GetOk("instance_name"); ok {
		request.InstanceName = helper.String(v.(string))
	}
	if v, ok := d.GetOk("vpc_id"); ok {
		request.VpcId = helper.String(v.(string))
	}
	if v, ok := d.GetOk("subnet_id"); ok {
		request.SubnetId = helper.String(v.(string))
	}
	if v, ok := d.GetOk("data_retention_time"); ok {
		request.DataRetentionTime = helper.IntInt64(v.(int))
	}
	if v, ok := d.GetOk("zone"); ok {
		request.Zone = helper.String(v.(string))
	}
	if v, ok := d.GetOk("grafana_instance_id"); ok {
		request.GrafanaInstanceId = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseMonitorClient().CreatePrometheusMultiTenantInstancePostPayMode(request)
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
		log.Printf("[CRITAL]%s create monitor tmpInstance failed, reason:%+v", logId, err)
		return err
	}

	tmpInstanceId := *response.Response.InstanceId
	d.SetId(tmpInstanceId)

	return resourceTencentCloudMonitorTmpInstanceRead(d, meta)
}

func resourceTencentCloudMonitorTmpInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_monitor_tmp_instance.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := MonitorService{client: meta.(*TencentCloudClient).apiV3Conn}

	tmpInstanceId := d.Id()

	tmpInstance, err := service.DescribeMonitorTmpInstanceById(ctx, tmpInstanceId)

	if err != nil {
		return err
	}

	if tmpInstance == nil {
		d.SetId("")
		return fmt.Errorf("resource `tmpInstance` %s does not exist", tmpInstanceId)
	}

	if tmpInstance.InstanceName != nil {
		_ = d.Set("instance_name", tmpInstance.InstanceName)
	}
	if tmpInstance.VpcId != nil {
		_ = d.Set("vpc_id", tmpInstance.VpcId)
	}
	if tmpInstance.SubnetId != nil {
		_ = d.Set("subnet_id", tmpInstance.SubnetId)
	}
	if tmpInstance.DataRetentionTime != nil {
		_ = d.Set("data_retention_time", tmpInstance.DataRetentionTime)
	}
	if tmpInstance.Zone != nil {
		_ = d.Set("zone", tmpInstance.Zone)
	}
	if tmpInstance.GrafanaInstanceId != nil {
		_ = d.Set("grafana_instance_id", tmpInstance.GrafanaInstanceId)
	}

	return nil
}

func resourceTencentCloudMonitorTmpInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_monitor_tmp_instance.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := monitor.NewModifyPrometheusInstanceAttributesRequest()

	request.InstanceId = helper.String(d.Id())

	if d.HasChange("instance_name") {
		if v, ok := d.GetOk("instance_name"); ok {
			request.InstanceName = helper.String(v.(string))
		}
	}

	if d.HasChange("data_retention_time") {
		if v, ok := d.GetOk("data_retention_time"); ok {
			request.DataRetentionTime = helper.IntInt64(v.(int))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseMonitorClient().ModifyPrometheusInstanceAttributes(request)
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

	return resourceTencentCloudMonitorTmpInstanceRead(d, meta)
}

func resourceTencentCloudMonitorTmpInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_monitor_tmp_instance.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := MonitorService{client: meta.(*TencentCloudClient).apiV3Conn}
	id := d.Id()

	if err := service.DeleteMonitorTmpInstance(ctx, id); err != nil {
		return err
	}

	return nil
}
