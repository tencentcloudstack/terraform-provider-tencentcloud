/*
Provides an EIP resource.

Example Usage

```hcl
resource "tencentcloud_eip" "foo" {
  name = "awesome_gateway_ip"
}
```

Import

EIP can be imported using the id, e.g.

```
$ terraform import tencentcloud_eip.foo eip-nyvf60va
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
	"github.com/terraform-providers/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

func resourceTencentCloudEip() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudEipCreate,
		Read:   resourceTencentCloudEipRead,
		Update: resourceTencentCloudEipUpdate,
		Delete: resourceTencentCloudEipDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateStringLengthInRange(1, 20),
				Description:  "The name of eip.",
			},
			"type": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      EIP_TYPE_EIP,
				ForceNew:     true,
				ValidateFunc: validateAllowedStringValue(EIP_TYPE),
				Description:  "The type of eip, and available values include `EIP` and `AnycastEIP`. Default is `EIP`.",
			},
			"anycast_zone": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateAllowedStringValue(EIP_ANYCAST_ZONE),
				Description:  "The zone of anycast, and available values include `ANYCAST_ZONE_GLOBAL` and `ANYCAST_ZONE_OVERSEAS`.",
			},
			"applicable_for_clb": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Description: "Indicates whether the anycast eip can be associated to a CLB.",
			},
			"internet_service_provider": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateAllowedStringValue(EIP_INTERNET_PROVIDER),
				Description:  "Internet service provider of eip, and available values include `BGP`, `CMCC`, `CTCC` and `CUCC`.",
			},
			"internet_charge_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateAllowedStringValue(CVM_INTERNET_CHARGE_TYPE),
				Description:  "The charge type of eip, and available values include `BANDWIDTH_PACKAGE`, `BANDWIDTH_POSTPAID_BY_HOUR` and `TRAFFIC_POSTPAID_BY_HOUR`.",
			},
			"internet_max_bandwidth_out": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateIntegerInRange(1, 1000),
				Description:  "The bandwidth limit of eip, unit is Mbps, and the range is 1-1000.",
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "The tags of eip.",
			},

			// computed
			"public_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The elastic ip address.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The eip current status.",
			},
		},
	}
}

func resourceTencentCloudEipCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_eip.create")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	client := meta.(*TencentCloudClient).apiV3Conn
	vpcService := VpcService{client: client}
	tagService := TagService{client: client}
	region := client.Region

	request := vpc.NewAllocateAddressesRequest()
	if v, ok := d.GetOk("type"); ok {
		request.AddressType = stringToPointer(v.(string))
	}
	if v, ok := d.GetOk("anycast_zone"); ok {
		request.AnycastZone = stringToPointer(v.(string))
	}
	if v, ok := d.GetOkExists("applicable_for_clb"); ok {
		applicable := v.(bool)
		request.ApplicableForCLB = &applicable
	}
	if v, ok := d.GetOk("internet_service_provider"); ok {
		request.InternetServiceProvider = stringToPointer(v.(string))
	}
	if v, ok := d.GetOk("internet_charge_type"); ok {
		request.InternetChargeType = stringToPointer(v.(string))
	}
	if v, ok := d.GetOk("internet_max_bandwidth_out"); ok {
		request.InternetMaxBandwidthOut = int64ToPointer(v.(int))
	}

	eipId := ""
	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		response, err := client.UseVpcClient().AllocateAddresses(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			return retryError(err)
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if len(response.Response.AddressSet) < 1 {
			return resource.NonRetryableError(fmt.Errorf("eip id is nil"))
		}
		eipId = *response.Response.AddressSet[0]
		return nil
	})
	if err != nil {
		return err
	}
	d.SetId(eipId)

	// wait for status
	err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		eip, errRet := vpcService.DescribeEipById(ctx, eipId)
		if errRet != nil {
			return retryError(errRet)
		}
		if *eip.AddressStatus == EIP_STATUS_CREATING {
			return resource.RetryableError(fmt.Errorf("eip is still creating"))
		}
		return nil
	})
	if err != nil {
		return err
	}

	if v, ok := d.GetOk("name"); ok {
		name := v.(string)
		err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			errRet := vpcService.ModifyEipName(ctx, eipId, name)
			if errRet != nil {
				return retryError(errRet)
			}
			return nil
		})
		if err != nil {
			return err
		}
	}

	if tags := getTags(d, "tags"); len(tags) > 0 {
		resourceName := BuildTagResourceName(VPC_SERVICE_TYPE, EIP_RESOURCE_TYPE, region, eipId)
		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			log.Printf("[CRITAL]%s set eip tags failed: %+v", logId, err)
			return err
		}
	}

	return resourceTencentCloudEipRead(d, meta)
}

func resourceTencentCloudEipRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_eip.read")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	client := meta.(*TencentCloudClient).apiV3Conn
	vpcService := VpcService{client: client}
	tagService := TagService{client: client}
	region := client.Region

	eipId := d.Id()
	var eip *vpc.Address
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		instance, errRet := vpcService.DescribeEipById(ctx, eipId)
		if errRet != nil {
			return retryError(errRet)
		}
		eip = instance
		return nil
	})
	if err != nil {
		return err
	}
	if eip == nil {
		d.SetId("")
		return nil
	}

	tags, err := tagService.DescribeResourceTags(ctx, VPC_SERVICE_TYPE, EIP_RESOURCE_TYPE, region, eipId)
	if err != nil {
		log.Printf("[CRITAL]%s describe eip tags failed: %+v", logId, err)
		return err
	}

	d.Set("name", eip.AddressName)
	d.Set("type", eip.AddressType)
	d.Set("public_ip", eip.AddressIp)
	d.Set("status", eip.AddressStatus)
	d.Set("tags", tags)
	return nil
}

func resourceTencentCloudEipUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_eip.update")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	client := meta.(*TencentCloudClient).apiV3Conn
	vpcService := VpcService{client: client}
	tagService := TagService{client: client}
	region := client.Region

	eipId := d.Id()

	d.Partial(true)

	if d.HasChange("name") {
		name := d.Get("name").(string)
		err := vpcService.ModifyEipName(ctx, eipId, name)
		if err != nil {
			return err
		}

		d.SetPartial("name")
	}

	if d.HasChange("tags") {
		oldTags, newTags := d.GetChange("tags")
		replaceTags, deleteTags := diffTags(oldTags.(map[string]interface{}), newTags.(map[string]interface{}))
		resourceName := BuildTagResourceName(VPC_SERVICE_TYPE, EIP_RESOURCE_TYPE, region, eipId)

		if err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags); err != nil {
			log.Printf("[CRITAL]%s update eip tags failed: %+v", logId, err)
			return err
		}
		d.SetPartial("tags")
	}

	return resourceTencentCloudEipRead(d, meta)
}

func resourceTencentCloudEipDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_eip.delete")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)
	vpcService := VpcService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	eipId := d.Id()
	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		errRet := vpcService.UnattachEip(ctx, eipId)
		if errRet != nil {
			return retryError(errRet, "DesOperation.MutexTaskRunning")
		}
		return nil
	})
	if err != nil {
		return err
	}

	err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		errRet := vpcService.DeleteEip(ctx, eipId)
		if err != nil {
			return retryError(errRet)
		}
		return nil
	})
	if err != nil {
		return err
	}

	err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		eip, errRet := vpcService.DescribeEipById(ctx, eipId)
		if errRet != nil {
			return retryError(errRet)
		}
		if eip != nil {
			return resource.RetryableError(fmt.Errorf("eip is still deleting"))
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
