/*
Provides a resource to create a monitor tmpInstance

Example Usage

```hcl
resource "tencentcloud_monitor_tmp_instance" "tmpInstance" {
  instance_name = "demo"
  vpc_id = "vpc-2hfyray3"
  subnet_id = "subnet-rdkj0agk"
  data_retention_time = 30
  zone = "ap-guangzhou-3"
  tags = {
    "createdBy" = "terraform"
  }
}

```
Import

monitor tmpInstance can be imported using the id, e.g.
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

			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Tag description list.",
			},

			"ipv4_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Instance IPv4 address.",
			},

			"remote_write": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Prometheus remote write address.",
			},

			"api_root_path": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Prometheus HTTP API root address.",
			},

			"proxy_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Proxy address.",
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

	service := MonitorService{client: meta.(*TencentCloudClient).apiV3Conn}
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	err = resource.Retry(3*readRetryTimeout, func() *resource.RetryError {
		instance, errRet := service.DescribeMonitorTmpInstance(ctx, tmpInstanceId)
		if errRet != nil {
			return retryError(errRet, InternalError)
		}
		if *instance.InstanceStatus == 2 {
			return nil
		}
		if *instance.InstanceStatus == 3 {
			return resource.NonRetryableError(fmt.Errorf("tmpInstance status is %v, operate failed.", *instance.InstanceStatus))
		}
		return resource.RetryableError(fmt.Errorf("tmpInstance status is %v, retry...", *instance.InstanceStatus))
	})
	if err != nil {
		return err
	}

	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		tagService := TagService{client: meta.(*TencentCloudClient).apiV3Conn}
		region := meta.(*TencentCloudClient).apiV3Conn.Region
		resourceName := fmt.Sprintf("qcs::monitor:%s:uin/:prom-instance/%s", region, tmpInstanceId)
		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
	}
	d.SetId(tmpInstanceId)
	return resourceTencentCloudMonitorTmpInstanceRead(d, meta)
}

func resourceTencentCloudMonitorTmpInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_monitor_tmpInstance.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := MonitorService{client: meta.(*TencentCloudClient).apiV3Conn}

	tmpInstanceId := d.Id()

	tmpInstance, err := service.DescribeMonitorTmpInstance(ctx, tmpInstanceId)

	if err != nil {
		return err
	}

	if tmpInstance == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `tmpInstance` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
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

	if tmpInstance.IPv4Address != nil {
		_ = d.Set("ipv4_address", tmpInstance.IPv4Address)
	}

	if tmpInstance.RemoteWrite != nil {
		_ = d.Set("remote_write", tmpInstance.RemoteWrite)
	}

	if tmpInstance.ApiRootPath != nil {
		_ = d.Set("api_root_path", tmpInstance.ApiRootPath)
	}

	if tmpInstance.ProxyAddress != nil {
		_ = d.Set("proxy_address", tmpInstance.ProxyAddress)
	}

	tcClient := meta.(*TencentCloudClient).apiV3Conn
	tagService := &TagService{client: tcClient}
	tags, err := tagService.DescribeResourceTags(ctx, "monitor", "prom-instance", tcClient.Region, d.Id())
	if err != nil {
		return err
	}
	_ = d.Set("tags", tags)

	return nil
}

func resourceTencentCloudMonitorTmpInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_monitor_tmp_instance.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	request := monitor.NewModifyPrometheusInstanceAttributesRequest()

	request.InstanceId = helper.String(d.Id())

	if v, ok := d.GetOk("instance_name"); ok {
		request.InstanceName = helper.String(v.(string))
	}

	if d.HasChange("vpc_id") {
		return fmt.Errorf("`vpc_id` do not support change now.")
	}

	if d.HasChange("subnet_id") {
		return fmt.Errorf("`subnet_id` do not support change now.")
	}

	if d.HasChange("data_retention_time") {
		if v, ok := d.GetOk("data_retention_time"); ok {
			request.DataRetentionTime = helper.IntInt64(v.(int))
		}
	}

	if d.HasChange("zone") {
		return fmt.Errorf("`zone` do not support change now.")
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

	if d.HasChange("tags") {
		tcClient := meta.(*TencentCloudClient).apiV3Conn
		tagService := &TagService{client: tcClient}
		oldTags, newTags := d.GetChange("tags")
		replaceTags, deleteTags := diffTags(oldTags.(map[string]interface{}), newTags.(map[string]interface{}))
		resourceName := BuildTagResourceName("monitor", "prom-instance", tcClient.Region, d.Id())
		if err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags); err != nil {
			return err
		}
	}

	return resourceTencentCloudMonitorTmpInstanceRead(d, meta)
}

func resourceTencentCloudMonitorTmpInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_monitor_tmp_instance.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := MonitorService{client: meta.(*TencentCloudClient).apiV3Conn}
	tmpInstanceId := d.Id()

	if err := service.IsolateMonitorTmpInstanceById(ctx, tmpInstanceId); err != nil {
		return err
	}

	err := resource.Retry(1*readRetryTimeout, func() *resource.RetryError {
		instance, errRet := service.DescribeMonitorTmpInstance(ctx, tmpInstanceId)
		if errRet != nil {
			return retryError(errRet, InternalError)
		}
		if *instance.InstanceStatus == 6 {
			return nil
		}
		if *instance.InstanceStatus == 3 {
			return resource.NonRetryableError(fmt.Errorf("tmpInstance status is %v, operate failed.", *instance.InstanceStatus))
		}
		return resource.RetryableError(fmt.Errorf("tmpInstance status is %v, retry...", *instance.InstanceStatus))
	})
	if err != nil {
		return err
	}

	if err := service.DeleteMonitorTmpInstanceById(ctx, tmpInstanceId); err != nil {
		return err
	}
	return nil
}
