/*
Provide a resource to create a CynosDB readonly instance.

Example Usage

```hcl
resource "tencentcloud_cynosdb_readonly_instance" "foo" {
  cluster_id           = cynosdbmysql-dzj5l8gz
  instance_name        = "tf-cynosdb-readonly-instance"
  force_delete         = true
  instance_cpu_core    = 2
  instance_memory_size = 4

  instance_maintain_duration   = 7200
  instance_maintain_start_time = 21600
  instance_maintain_weekdays   = [
    "Fri",
    "Mon",
    "Sat",
    "Sun",
    "Thu",
    "Wed",
    "Tue",
  ]
}
```

Import

CynosDB readonly instance can be imported using the id, e.g.

```
$ terraform import tencentcloud_cynosdb_readonly_instance.foo cynosdbmysql-ins-dhwynib6
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cynosdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cynosdb/v20190107"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

func resourceTencentCloudCynosdbReadonlyInstance() *schema.Resource {
	instanceInfo := map[string]*schema.Schema{
		"cluster_id": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "Cluster ID which the readonly instance belongs to.",
		},
		"instance_name": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "Name of instance.",
		},
		"force_delete": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Indicate whether to delete readonly instance directly or not. Default is false. If set true, instance will be deleted instead of staying recycle bin. Note: works for both `PREPAID` and `POSTPAID_BY_HOUR` cluster.",
		},
	}
	basic := TencentCynosdbInstanceBaseInfo()
	delete(basic, "instance_id")
	delete(basic, "instance_name")
	for k, v := range basic {
		instanceInfo[k] = v
	}

	return &schema.Resource{
		Create: resourceTencentCloudCynosdbReadonlyInstanceCreate,
		Read:   resourceTencentCloudCynosdbReadonlyInstanceRead,
		Update: resourceTencentCloudCynosdbReadonlyInstanceUpdate,
		Delete: resourceTencentCloudCynosdbReadonlyInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: instanceInfo,
	}
}

func resourceTencentCloudCynosdbReadonlyInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_readonly_instance.create")()

	var (
		logId = getLogId(contextNil)
		ctx   = context.WithValue(context.TODO(), logIdKey, logId)

		client         = meta.(*TencentCloudClient).apiV3Conn
		cynosdbService = CynosdbService{client: client}

		request = cynosdb.NewAddInstancesRequest()
	)

	// instance info
	request.ClusterId = helper.String(d.Get("cluster_id").(string))
	request.InstanceName = helper.String(d.Get("instance_name").(string))
	request.Cpu = helper.IntInt64(d.Get("instance_cpu_core").(int))
	request.Memory = helper.IntInt64(d.Get("instance_memory_size").(int))
	request.ReadOnlyCount = helper.Int64(1)

	var response *cynosdb.AddInstancesResponse
	var err error
	err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		response, err = meta.(*TencentCloudClient).apiV3Conn.UseCynosdbClient().AddInstances(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, reason:%s", logId, request.GetAction(), err.Error())
			return retryError(err)
		}
		return nil
	})
	if err != nil {
		return err
	}
	if response != nil && response.Response != nil && len(response.Response.ResourceIds) != 1 {
		return fmt.Errorf("cynosdb readonly instance id count isn't 1")
	}
	d.SetId(*response.Response.ResourceIds[0])
	id := d.Id()

	// set maintenance info
	var weekdays []interface{}
	if v, ok := d.GetOk("instance_maintain_weekdays"); ok {
		weekdays = v.(*schema.Set).List()
	} else {
		weekdays = []interface{}{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"}
	}
	reqWeekdays := make([]*string, 0, len(weekdays))
	for _, v := range weekdays {
		reqWeekdays = append(reqWeekdays, helper.String(v.(string)))
	}
	startTime := int64(d.Get("instance_maintain_start_time").(int))
	duration := int64(d.Get("instance_maintain_duration").(int))
	err = cynosdbService.ModifyMaintainPeriodConfig(ctx, id, startTime, duration, reqWeekdays)
	if err != nil {
		return err
	}

	return resourceTencentCloudCynosdbReadonlyInstanceRead(d, meta)
}

func resourceTencentCloudCynosdbReadonlyInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_readonly_instance.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	id := d.Id()

	client := meta.(*TencentCloudClient).apiV3Conn
	cynosdbService := CynosdbService{client: client}
	clusterId, instance, has, err := cynosdbService.DescribeInstanceById(ctx, id)
	if err != nil {
		return err
	}
	if !has {
		d.SetId("")
		return nil
	}

	_ = d.Set("cluster_id", clusterId)
	_ = d.Set("instance_cpu_core", instance.Cpu)
	_ = d.Set("instance_memory_size", instance.Memory)
	_ = d.Set("instance_name", instance.InstanceName)
	_ = d.Set("instance_status", instance.Status)
	_ = d.Set("instance_storage_size", instance.Storage)

	maintain, err := cynosdbService.DescribeMaintainPeriod(ctx, id)
	if err != nil {
		return err
	}
	_ = d.Set("instance_maintain_weekdays", maintain.Response.MaintainWeekDays)
	_ = d.Set("instance_maintain_start_time", maintain.Response.MaintainStartTime)
	_ = d.Set("instance_maintain_duration", maintain.Response.MaintainDuration)

	return nil
}

func resourceTencentCloudCynosdbReadonlyInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_readonly_instance.update")()

	var (
		logId          = getLogId(contextNil)
		ctx            = context.WithValue(context.TODO(), logIdKey, logId)
		instanceId     = d.Id()
		client         = meta.(*TencentCloudClient).apiV3Conn
		cynosdbService = CynosdbService{client: client}
	)

	d.Partial(true)

	if d.HasChange("instance_cpu_core") || d.HasChange("instance_memory_size") {
		cpu := int64(d.Get("instance_cpu_core").(int))
		memory := int64(d.Get("instance_memory_size").(int))
		err := cynosdbService.UpgradeInstance(ctx, instanceId, cpu, memory)
		if err != nil {
			return err
		}

		errUpdate := resource.Retry(readRetryTimeout, func() *resource.RetryError {
			_, infos, has, e := cynosdbService.DescribeInstanceById(ctx, instanceId)
			if e != nil {
				return resource.NonRetryableError(e)
			}
			if !has {
				return resource.NonRetryableError(fmt.Errorf("[CRITAL]%s updating cynosdb cluster instance failed, instance doesn't exist", logId))
			}

			cpuReal := *infos.Cpu
			memReal := *infos.Memory
			if cpu != cpuReal || memory != memReal {
				return resource.RetryableError(fmt.Errorf("[CRITAL] updating cynosdb instance, current cpu and memory values: %d, %d, waiting for them becoming new value: %d, %d", cpuReal, memReal, cpu, memory))
			}
			return nil
		})
		if errUpdate != nil {
			return errUpdate
		}

	}

	if d.HasChange("instance_maintain_weekdays") || d.HasChange("instance_maintain_start_time") || d.HasChange("instance_maintain_duration") {
		weekdays := d.Get("instance_maintain_weekdays").(*schema.Set).List()
		reqWeekdays := make([]*string, 0, len(weekdays))
		for _, v := range weekdays {
			reqWeekdays = append(reqWeekdays, helper.String(v.(string)))
		}
		startTime := int64(d.Get("instance_maintain_start_time").(int))
		duration := int64(d.Get("instance_maintain_duration").(int))
		err := cynosdbService.ModifyMaintainPeriodConfig(ctx, instanceId, startTime, duration, reqWeekdays)
		if err != nil {
			return err
		}

	}

	d.Partial(false)

	return resourceTencentCloudCynosdbReadonlyInstanceRead(d, meta)
}

func resourceTencentCloudCynosdbReadonlyInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_readonly_instance.delete")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	instanceId := d.Id()
	clusterId := d.Get("cluster_id").(string)
	cynosdbService := CynosdbService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}

	forceDelete := d.Get("force_delete").(bool)

	var err error
	if err = cynosdbService.IsolateInstance(ctx, clusterId, instanceId); err != nil {
		return err
	}

	if forceDelete {
		if err = cynosdbService.OfflineInstance(ctx, clusterId, instanceId); err != nil {
			return err
		}
	}

	return nil
}
