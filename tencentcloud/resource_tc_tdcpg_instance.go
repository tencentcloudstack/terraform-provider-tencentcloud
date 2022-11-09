/*
Provides a resource to create a tdcpg instance

Example Usage

```hcl
resource "tencentcloud_tdcpg_instance" "instance" {
  cluster_id = ""
  c_p_u = ""
  memory = ""
  instance_name = ""
  instance_count = ""
  operation_timing = ""
}

```
Import

tdcpg instance can be imported using the id, e.g.
```
$ terraform import tencentcloud_tdcpg_instance.instance instance_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	tdcpg "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tdcpg/v20211118"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudTdcpgInstance() *schema.Resource {
	return &schema.Resource{
		Read:   resourceTencentCloudTdcpgInstanceRead,
		Create: resourceTencentCloudTdcpgInstanceCreate,
		Update: resourceTencentCloudTdcpgInstanceUpdate,
		Delete: resourceTencentCloudTdcpgInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "cluster id.",
			},

			"c_p_u": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "cpu cores.",
			},

			"memory": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "memory size.",
			},

			"instance_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "instance name.",
			},

			"instance_count": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "instance count.",
			},

			"operation_timing": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "operation timing, optional value is IMMEDIATE or MAINTAIN_PERIOD.",
			},
		},
	}
}

func resourceTencentCloudTdcpgInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tdcpg_instance.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = tdcpg.NewCreateClusterInstancesRequest()
		response   *tdcpg.CreateClusterInstancesResponse
		instanceId string
		clusterId  string
	)

	if v, ok := d.GetOk("cluster_id"); ok {
		clusterId = v.(string)
		request.ClusterId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("CPU"); ok {
		request.CPU = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("memory"); ok {
		request.Memory = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("instance_name"); ok {

		request.InstanceName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("instance_count"); ok {
		request.InstanceCount = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("operation_timing"); ok {

		request.OperationTiming = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTdcpgClient().CreateClusterInstances(request)
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
		log.Printf("[CRITAL]%s create tdcpg instance failed, reason:%+v", logId, err)
		return err
	}

	instanceId := *response.Response.InstanceId

	service := TdcpgService{client: meta.(*TencentCloudClient).apiV3Conn}
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	err = resource.Retry(5*readRetryTimeout, func() *resource.RetryError {
		instance, errRet := service.DescribeTdcpgInstance(ctx, instanceId)
		if errRet != nil {
			return retryError(errRet, InternalError)
		}
		if *instance.resourceInfo == ready {
			return nil
		}
		if *instance.resourceInfo == failed {
			return resource.NonRetryableError(fmt.Errorf("instance status is %v, operate failed.", *instance.resourceInfo))
		}
		return resource.RetryableError(fmt.Errorf("instance status is %v, retry...", *instance.resourceInfo))
	})
	if err != nil {
		return err
	}

	d.SetId(instanceId + FILED_SP + clusterId)
	return resourceTencentCloudTdcpgInstanceRead(d, meta)
}

func resourceTencentCloudTdcpgInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tdcpg_instance.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TdcpgService{client: meta.(*TencentCloudClient).apiV3Conn}

	instanceId := d.Id()

	instance, err := service.DescribeTdcpgInstance(ctx, instanceId)

	if err != nil {
		return err
	}

	if instance == nil {
		d.SetId("")
		return fmt.Errorf("resource `instance` %s does not exist", instanceId)
	}

	if instance.ClusterId != nil {
		_ = d.Set("cluster_id", instance.ClusterId)
	}

	if instance.CPU != nil {
		_ = d.Set("c_p_u", instance.CPU)
	}

	if instance.Memory != nil {
		_ = d.Set("memory", instance.Memory)
	}

	if instance.InstanceName != nil {
		_ = d.Set("instance_name", instance.InstanceName)
	}

	if instance.InstanceCount != nil {
		_ = d.Set("instance_count", instance.InstanceCount)
	}

	if instance.OperationTiming != nil {
		_ = d.Set("operation_timing", instance.OperationTiming)
	}

	return nil
}

func resourceTencentCloudTdcpgInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tdcpg_instance.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	request := tdcpg.NewModifyClusterInstancesSpecRequest()

	instanceId := d.Id()

	request.InstanceId = &instanceId

	if d.HasChange("cluster_id") {

		return fmt.Errorf("`cluster_id` do not support change now.")

	}

	if d.HasChange("c_p_u") {
		if v, ok := d.GetOk("c_p_u"); ok {
			request.CPU = helper.IntInt64(v.(int))
		}

	}

	if d.HasChange("memory") {
		if v, ok := d.GetOk("memory"); ok {
			request.Memory = helper.IntInt64(v.(int))
		}

	}

	if d.HasChange("instance_name") {

		return fmt.Errorf("`instance_name` do not support change now.")

	}

	if d.HasChange("instance_count") {

		return fmt.Errorf("`instance_count` do not support change now.")

	}

	if d.HasChange("operation_timing") {

		return fmt.Errorf("`operation_timing` do not support change now.")

	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTdcpgClient().ModifyClusterInstancesSpec(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create tdcpg instance failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudTdcpgInstanceRead(d, meta)
}

func resourceTencentCloudTdcpgInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tdcpg_instance.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TdcpgService{client: meta.(*TencentCloudClient).apiV3Conn}

	instanceId := d.Id()

	if err := service.DeleteTdcpgInstanceById(ctx, instanceId); err != nil {
		return err
	}

	return nil
}
