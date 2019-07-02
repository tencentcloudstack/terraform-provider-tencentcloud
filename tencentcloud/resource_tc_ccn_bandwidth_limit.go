/*
Provides a resource to limit CCN bandwidth.

Example Usage

```hcl
variable "other_region1" {
    default = "ap-shanghai"
}
resource tencentcloud_ccn main{
	name ="ci-temp-test-ccn"
	description="ci-temp-test-ccn-des"
	qos ="AG"
}

resource tencentcloud_ccn_bandwidth_limit limit1 {
	ccn_id ="${tencentcloud_ccn.main.id}"
	region ="${var.other_region1}"
	bandwidth_limit = 500
}
```
*/
package tencentcloud

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceTencentCloudCcnBandwidthLimit() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCcnBandwidthLimitCreate,
		Read:   resourceTencentCloudCcnBandwidthLimitRead,
		Update: resourceTencentCloudCcnBandwidthLimitUpdate,
		Delete: resourceTencentCloudCcnBandwidthLimitDelete,

		Schema: map[string]*schema.Schema{
			"ccn_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the CCN",
			},
			"region": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Limitation of region.",
			},
			"bandwidth_limit": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Limitation of bandwidth.",
			},
		},
	}
}
func resourceTencentCloudCcnBandwidthLimitCreate(d *schema.ResourceData, meta interface{}) error {

	logId := GetLogId(nil)
	defer LogElapsed(logId + "resource.tencentcloud_ccn_bandwidth_limit.create")()

	ctx := context.WithValue(context.TODO(), "logId", logId)
	service := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}

	var (
		ccnId  = d.Get("ccn_id").(string)
		region = d.Get("region").(string)
	)

	_, has, err := service.DescribeCcn(ctx, ccnId)
	if err != nil {
		return err
	}
	if has == 0 {
		return fmt.Errorf("ccn[%s] doesn't exist", ccnId)
	}

	id := fmt.Sprintf("%s#%s", ccnId, region)

	if limitTemp, ok := d.GetOk("bandwidth_limit"); ok {
		if err := service.SetCcnRegionBandwidthLimits(ctx, ccnId, region, int64(limitTemp.(int))); err != nil {
			return err
		}
	}
	d.SetId(id)

	return resourceTencentCloudCcnBandwidthLimitRead(d, meta)
}

func resourceTencentCloudCcnBandwidthLimitUpdate(d *schema.ResourceData, meta interface{}) error {

	logId := GetLogId(nil)
	defer LogElapsed(logId + "resource.tencentcloud_ccn_bandwidth_limit.update")()

	ctx := context.WithValue(context.TODO(), "logId", logId)
	service := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}

	var (
		ccnId  = d.Get("ccn_id").(string)
		region = d.Get("region").(string)
	)
	_, has, err := service.DescribeCcn(ctx, ccnId)
	if err != nil {
		return err
	}
	if has == 0 {
		return fmt.Errorf("ccn[%s] doesn't exist", ccnId)
	}

	if limitTemp, ok := d.GetOk("bandwidth_limit"); ok {
		if err := service.SetCcnRegionBandwidthLimits(ctx, ccnId, region, int64(limitTemp.(int))); err != nil {
			return err
		}
	}

	return resourceTencentCloudCcnBandwidthLimitRead(d, meta)
}

func resourceTencentCloudCcnBandwidthLimitRead(d *schema.ResourceData, meta interface{}) error {

	logId := GetLogId(nil)
	defer LogElapsed(logId + "resource.tencentcloud_ccn_bandwidth_limit.read")()

	ctx := context.WithValue(context.TODO(), "logId", logId)
	service := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}

	var (
		ccnId  = d.Get("ccn_id").(string)
		region = d.Get("region").(string)
	)

	_, has, err := service.DescribeCcn(ctx, ccnId)
	if err != nil {
		return err
	}

	if has == 0 {
		d.SetId("")
		return nil
	}

	bandwidth, err := service.DescribeCcnRegionBandwidthLimit(ctx, ccnId, region)
	if err != nil {
		return err
	}
	d.Set("bandwidth_limit", bandwidth)
	return nil
}

func resourceTencentCloudCcnBandwidthLimitDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
