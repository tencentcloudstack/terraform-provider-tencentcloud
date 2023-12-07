package tencentcloud

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
)

func resourceTencentCloudReservedInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudReservedInstanceCreate,
		Read:   resourceTencentCloudReservedInstanceRead,
		Update: resourceTencentCloudReservedInstanceUpdate,
		Delete: resourceTencentCloudReservedInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"config_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Configuration ID of the reserved instance.",
			},
			"instance_count": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validateIntegerMin(1),
				Description:  "Number of reserved instances to be purchased.",
			},
			"reserved_instance_name": {
				Type:     schema.TypeString,
				Optional: true,
				Description: `Reserved Instance display name.
				- If you do not specify an instance display name, 'Unnamed' is displayed by default.
				- Up to 60 characters (including pattern strings) are supported.`,
			},

			// computed
			"start_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Start time of the RI.",
			},
			"end_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Expiry time of the RI.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status of the RI at the time of purchase.",
			},
		},
	}
}

func resourceTencentCloudReservedInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_reserved_instance.create")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	configId := d.Get("config_id").(string)
	count := d.Get("instance_count").(int)
	cvmService := CvmService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	extendParams := make(map[string]interface{})
	if v, ok := d.GetOk("reserved_instance_name"); ok {
		extendParams["reserved_instance_name"] = v.(string)
	}
	var instanceId string
	var errRet error
	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		instanceId, errRet = cvmService.CreateReservedInstance(ctx, configId, int64(count), extendParams)
		if errRet != nil {
			return retryError(errRet)
		}
		return nil
	})
	if err != nil {
		return err
	}
	d.SetId(instanceId)

	return resourceTencentCloudReservedInstanceRead(d, meta)
}

func resourceTencentCloudReservedInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_reserved_instance.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	id := d.Id()
	cvmService := CvmService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	filter := map[string]string{
		"reserved-instances-id": id,
	}
	var instances []*cvm.ReservedInstances
	var errRet error
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		instances, errRet = cvmService.DescribeReservedInstanceByFilter(ctx, filter)
		if errRet != nil {
			return retryError(errRet, InternalError)
		}
		return nil
	})
	if err != nil {
		return err
	}
	if len(instances) < 1 {
		d.SetId("")
		return nil
	}
	instance := instances[0]

	_ = d.Set("instance_count", instance.InstanceCount)
	_ = d.Set("start_time", instance.StartTime)
	_ = d.Set("end_time", instance.EndTime)
	_ = d.Set("status", instance.State)
	_ = d.Set("reserved_instance_name", instance.ReservedInstanceName)

	return nil
}

func resourceTencentCloudReservedInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	return fmt.Errorf("reserved instance not allowed to modify")
}

func resourceTencentCloudReservedInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
