package tencentcloud

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceTencentCloudCnnBandwidthLimit() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCnnBandwidthLimitCreate,
		Read:   resourceTencentCloudCnnBandwidthLimitRead,
		Update: resourceTencentCloudCnnBandwidthLimitUpdate,
		Delete: resourceTencentCloudCnnBandwidthLimitDelete,

		Schema: map[string]*schema.Schema{
			"cnn_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"region": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"bandwidth_limit": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
		},
	}
}
func resourceTencentCloudCnnBandwidthLimitCreate(d *schema.ResourceData, meta interface{}) error {

	logId := GetLogId(nil)
	defer LogElapsed(logId + "resource.tencentcloud_cnn_bandwidth_limit.create")()

	ctx := context.WithValue(context.TODO(), "logId", logId)
	service := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}

	var (
		cnnId  = d.Get("cnn_id").(string)
		region = d.Get("region").(string)
	)

	_, has, err := service.DescribeCcn(ctx, cnnId)
	if err != nil {
		return err
	}
	if has == 0 {
		return fmt.Errorf("cnn[%s] doesn't exist", cnnId)
	}

	id := fmt.Sprintf("%s#%s", cnnId, region)

	if limitTemp, ok := d.GetOk("bandwidth_limit"); ok {
		if err := service.SetCcnRegionBandwidthLimits(ctx, cnnId, region, int64(limitTemp.(int))); err != nil {
			return err
		}
	}
	d.SetId(id)

	return resourceTencentCloudCnnBandwidthLimitRead(d, meta)
}

func resourceTencentCloudCnnBandwidthLimitUpdate(d *schema.ResourceData, meta interface{}) error {

	logId := GetLogId(nil)
	defer LogElapsed(logId + "resource.tencentcloud_cnn_bandwidth_limit.update")()

	ctx := context.WithValue(context.TODO(), "logId", logId)
	service := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}

	var (
		cnnId  = d.Get("cnn_id").(string)
		region = d.Get("region").(string)
	)
	_, has, err := service.DescribeCcn(ctx, cnnId)
	if err != nil {
		return err
	}
	if has == 0 {
		return fmt.Errorf("cnn[%s] doesn't exist", cnnId)
	}

	if limitTemp, ok := d.GetOk("bandwidth_limit"); ok {
		if err := service.SetCcnRegionBandwidthLimits(ctx, cnnId, region, int64(limitTemp.(int))); err != nil {
			return err
		}
	}

	return resourceTencentCloudCnnBandwidthLimitRead(d, meta)
}

func resourceTencentCloudCnnBandwidthLimitRead(d *schema.ResourceData, meta interface{}) error {

	logId := GetLogId(nil)
	defer LogElapsed(logId + "resource.tencentcloud_cnn_bandwidth_limit.read")()

	ctx := context.WithValue(context.TODO(), "logId", logId)
	service := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}

	var (
		cnnId  = d.Get("cnn_id").(string)
		region = d.Get("region").(string)
	)

	_, has, err := service.DescribeCcn(ctx, cnnId)
	if err != nil {
		return err
	}

	if has == 0 {
		d.SetId("")
		return nil
	}

	bandwidth, err := service.DescribeCcnRegionBandwidthLimit(ctx, cnnId, region)
	if err != nil {
		return err
	}
	d.Set("bandwidth_limit", bandwidth)
	return nil
}

func resourceTencentCloudCnnBandwidthLimitDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
