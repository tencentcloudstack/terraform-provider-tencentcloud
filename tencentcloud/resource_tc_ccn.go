/*
resource tencentcloud_ccn main{
	name ="ci-temp-test-ccn"
	description="ci-temp-test-ccn-des"
	qos ="AG"
}
*/
package tencentcloud

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceTencentCloudCcn() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCcnCreate,
		Read:   resourceTencentCloudCcnRead,
		Update: resourceTencentCloudCcnUpdate,
		Delete: resourceTencentCloudCcnDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateStringLengthInRange(1, 60),
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateStringLengthInRange(0, 100),
			},
			"qos": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      CNN_QOS_AU,
				ValidateFunc: validateAllowedStringValue([]string{CNN_QOS_PT, CNN_QOS_AU, CNN_QOS_AG}),
			},
			// Computed values
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"instance_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}
func resourceTencentCloudCcnCreate(d *schema.ResourceData, meta interface{}) error {

	logId := GetLogId(nil)
	defer LogElapsed(logId + "resource.tencentcloud_ccn.create")()

	ctx := context.WithValue(context.TODO(), "logId", logId)
	service := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}

	var (
		name        string = d.Get("name").(string)
		description string = ""
		qos         string = d.Get("qos").(string)
	)
	if temp, ok := d.GetOk("description"); ok {
		description = temp.(string)
	}
	info, err := service.CreateCcn(ctx, name, description, qos)
	if err != nil {
		return err
	}
	d.SetId(info.ccnId)

	return resourceTencentCloudCcnRead(d, meta)
}
func resourceTencentCloudCcnRead(d *schema.ResourceData, meta interface{}) error {

	logId := GetLogId(nil)
	defer LogElapsed(logId + "resource.tencentcloud_ccn.read")()

	ctx := context.WithValue(context.TODO(), "logId", logId)
	service := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}

	info, has, err := service.DescribeCcn(ctx, d.Id())
	if err != nil {
		return err
	}

	if has == 0 {
		d.SetId("")
		return nil
	}

	d.Set("name", info.name)
	d.Set("description", info.description)
	d.Set("qos", strings.ToUpper(info.qos))
	d.Set("state", strings.ToUpper(info.state))
	d.Set("instance_count", info.instanceCount)
	d.Set("create_time", info.createTime)

	return nil
}
func resourceTencentCloudCcnUpdate(d *schema.ResourceData, meta interface{}) error {

	logId := GetLogId(nil)
	defer LogElapsed(logId + "resource.tencentcloud_ccn.update")()

	ctx := context.WithValue(context.TODO(), "logId", logId)
	service := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}

	var (
		name        string = ""
		description string = ""
		change      bool   = false
	)
	if d.HasChange("name") {
		name = d.Get("name").(string)
		change = true
	}

	if d.HasChange("description") {
		if temp, ok := d.GetOk("description"); ok {
			description = temp.(string)
		}
		if description == "" {
			return fmt.Errorf("can not set description='' ")
		}
		change = true
	}

	if change {
		if err := service.ModifyCcnAttribute(ctx, d.Id(), name, description); err != nil {
			return err
		}
	}
	return resourceTencentCloudCcnRead(d, meta)
}

func resourceTencentCloudCcnDelete(d *schema.ResourceData, meta interface{}) error {
	logId := GetLogId(nil)
	defer LogElapsed(logId + "resource.tencentcloud_ccn.delete")()

	ctx := context.WithValue(context.TODO(), "logId", logId)
	service := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}

	_, has, err := service.DescribeCcn(ctx, d.Id())
	if err != nil {
		return err
	}
	if has == 0 {
		return nil
	}
	if err = service.DeleteCcn(ctx, d.Id()); err != nil {
		return err
	}

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, has, err := service.DescribeCcn(ctx, d.Id())
		if err != nil {
			return resource.RetryableError(err)
		}
		if has == 0 {
			return nil
		}
		return resource.RetryableError(fmt.Errorf("delete fail"))
	})
}
