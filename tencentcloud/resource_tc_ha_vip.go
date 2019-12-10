/*
Provides a resource to create a HA VIP.

Example Usage

```hcl
resource "tencentcloud_ha_vip" "foo" {
  name      = "terraform_test"
  vpc_id    = "vpc-gzea3dd7"
  subnet_id = "subnet-4d4m4cd4s"
  vip       = "10.0.4.16"
}
```

Import

HA VIP can be imported using the id, e.g.

```
$ terraform import tencentcloud_ha_vip.foo havip-kjqwe4ba
```
*/
package tencentcloud

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/pkg/errors"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
)

func resourceTencentCloudHaVip() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudHaVipCreate,
		Read:   resourceTencentCloudHaVipRead,
		Update: resourceTencentCloudHaVipUpdate,
		Delete: resourceTencentCloudHaVipDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateStringLengthInRange(1, 60),
				Description:  "Name of the HA VIP. The length of character is limited to 1-60.",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "VPC id.",
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Subnet id.",
			},
			"vip": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validateIp,
				Description:  "Virtual IP address, it must not be occupied and in this VPC network segment. If not set, it will be assigned after resource created automatically.",
			},
			"state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "State of the HA VIP, values are `AVAILABLE`, `UNBIND`.",
			},
			"instance_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Instance id that is associated.",
			},
			"network_interface_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Network interface id that is associated.",
			},
			"address_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "EIP that is associated.",
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Create time of the HA VIP.",
			},
		},
	}
}

func resourceTencentCloudHaVipCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ha_vip.create")()

	logId := getLogId(contextNil)

	request := vpc.NewCreateHaVipRequest()
	request.VpcId = stringToPointer(d.Get("vpc_id").(string))
	request.SubnetId = stringToPointer(d.Get("subnet_id").(string))
	request.HaVipName = stringToPointer(d.Get("name").(string))
	//optional
	if v, ok := d.GetOk("vip"); ok {
		request.Vip = stringToPointer(v.(string))
	}
	var response *vpc.CreateHaVipResponse
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseVpcClient().CreateHaVip(request)
		if e != nil {
			return retryError(errors.WithStack(e))
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create HA VIP failed, reason:%+v", logId, err)
		return err
	}

	if response.Response.HaVip == nil {
		return fmt.Errorf("HA VIP id is nil")
	}
	haVipId := *response.Response.HaVip.HaVipId
	d.SetId(haVipId)

	return resourceTencentCloudHaVipRead(d, meta)
}

func resourceTencentCloudHaVipRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ha_vip.read")()

	logId := getLogId(contextNil)

	haVipId := d.Id()
	request := vpc.NewDescribeHaVipsRequest()
	request.HaVipIds = []*string{&haVipId}
	var response *vpc.DescribeHaVipsResponse
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseVpcClient().DescribeHaVips(request)
		if e != nil {
			return retryError(errors.WithStack(e))
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read HA VIP failed, reason:%+v", logId, err)
		return err
	}
	if len(response.Response.HaVipSet) < 1 {
		d.SetId("")
		return nil
	}

	haVip := response.Response.HaVipSet[0]

	d.Set("name", *haVip.HaVipName)
	d.Set("create_time", *haVip.CreatedTime)
	d.Set("vip", *haVip.Vip)
	d.Set("vpc_id", *haVip.VpcId)
	d.Set("subnet_id", *haVip.SubnetId)
	d.Set("address_id", *haVip.AddressIp)
	d.Set("state", *haVip.State)
	d.Set("network_interface_id", *haVip.NetworkInterfaceId)
	d.Set("instance_id", *haVip.InstanceId)

	return nil
}

func resourceTencentCloudHaVipUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ha_vip.update")()

	logId := getLogId(contextNil)

	haVipId := d.Id()
	request := vpc.NewModifyHaVipAttributeRequest()
	request.HaVipId = &haVipId
	if d.HasChange("name") {
		request.HaVipName = stringToPointer(d.Get("name").(string))
		err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
			_, e := meta.(*TencentCloudClient).apiV3Conn.UseVpcClient().ModifyHaVipAttribute(request)
			if e != nil {
				return retryError(errors.WithStack(e))
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s modify HA VIP failed, reason:%+v", logId, err)
			return err
		}
	}

	return resourceTencentCloudHaVipRead(d, meta)
}

func resourceTencentCloudHaVipDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ha_vip.delete")()

	logId := getLogId(contextNil)

	haVipId := d.Id()

	request := vpc.NewDeleteHaVipRequest()
	request.HaVipId = &haVipId
	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		_, e := meta.(*TencentCloudClient).apiV3Conn.UseVpcClient().DeleteHaVip(request)
		if e != nil {
			return retryError(errors.WithStack(e), VPCUnsupportedOperation)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s delete HA VIP failed, reason:%+v", logId, err)
		return err
	}
	//to get the status of haVip
	statRequest := vpc.NewDescribeHaVipsRequest()
	statRequest.HaVipIds = []*string{&haVipId}
	err = resource.Retry(3*time.Minute, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseVpcClient().DescribeHaVips(statRequest)
		if e != nil {
			ee, ok := e.(*sdkErrors.TencentCloudSDKError)
			if !ok {
				return retryError(errors.WithStack(ee))
			}
			if ee.Code == VPCNotFound {
				log.Printf("[CRITAL]%s api[%s] success, request body [%s], reason[%s]\n",
					logId, statRequest.GetAction(), statRequest.ToJsonString(), e)
				return nil
			} else {
				//when associated eip is in deleting process, delete ha vip may return unsupported operation error
				return retryError(errors.WithStack(e), VPCUnsupportedOperation)
			}
		} else {
			//if not, quit
			if len(result.Response.HaVipSet) == 0 {
				return nil
			}
			//else consider delete fail
			return resource.RetryableError(fmt.Errorf("deleting retry"))
		}
	})
	if err != nil {
		log.Printf("[CRITAL]%s delete HA VIP failed, reason:%+v", logId, err)
		return err
	}
	return nil
}
