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

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
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
	vpcService := VpcService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}

	eipId := ""
	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		id, errRet := vpcService.CreateEip(ctx)
		if errRet != nil {
			return retryError(errRet)
		}
		eipId = id
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

	return resourceTencentCloudEipRead(d, meta)
}

func resourceTencentCloudEipRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_eip.read")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)
	vpcService := VpcService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}

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

	d.Set("name", eip.AddressName)
	d.Set("public_ip", eip.AddressIp)
	d.Set("status", eip.AddressStatus)
	return nil
}

func resourceTencentCloudEipUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_eip.update")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)
	vpcService := VpcService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}

	eipId := d.Id()
	if d.HasChange("name") {
		name := d.Get("name").(string)
		err := vpcService.ModifyEipName(ctx, eipId, name)
		if err != nil {
			return err
		}
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
	var eip *vpc.Address
	var errRet error
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		eip, errRet = vpcService.DescribeEipById(ctx, eipId)
		if errRet != nil {
			return retryError(errRet, "InternalError")
		}
		return nil
	})
	if err != nil {
		return err
	}
	if eip == nil {
		return nil
	}

	if eip.InstanceId != nil {
		err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			errRet := vpcService.UnattachEip(ctx, eipId)
			if errRet != nil {
				return retryError(errRet)
			}
			return nil
		})
		if err != nil {
			return err
		}
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
