/*
Provides an eip resource associated with other resource like CVM or ENI.

~> **NOTE:** Please DO NOT define `allocate_public_ip` in `tencentcloud_instance` resource when using `tencentcloud_eip_association`.

Example Usage

```hcl
resource "tencentcloud_eip_association" "foo" {
  eip_id      = "eip-xxxxxx"
  instance_id = "ins-xxxxxx"
}
```

or

```hcl
resource "tencentcloud_eip_association" "bar" {
  eip_id               = "eip-xxxxxx"
  network_interface_id = "eni-xxxxxx"
  private_ip           = "10.0.1.22"
}
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/terraform-providers/terraform-provider-tencentcloud/tencentcloud/ratelimit"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
)

func resourceTencentCloudEipAssociation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudEipAssociationCreate,
		Read:   resourceTencentCloudEipAssociationRead,
		Delete: resourceTencentCloudEipAssociationDelete,

		Schema: map[string]*schema.Schema{
			"eip_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateStringLengthInRange(1, 25),
				Description:  "The id of eip.",
			},
			"instance_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
				Computed: true,
				ConflictsWith: []string{
					"network_interface_id",
					"private_ip",
				},
				ValidateFunc: validateStringLengthInRange(1, 25),
				Description:  "The instance id going to bind with the EIP. This field is conflict with `network_interface_id` and `private_ip fields`.",
			},
			"network_interface_id": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateStringLengthInRange(1, 25),
				ConflictsWith: []string{
					"instance_id",
				},
				Description: "Indicates the network interface id like `eni-xxxxxx`. This field is conflict with `instance_id`.",
			},
			"private_ip": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateStringLengthInRange(7, 25),
				ConflictsWith: []string{
					"instance_id",
				},
				Description: "Indicates an IP belongs to the `network_interface_id`. This field is conflict with `instance_id`.",
			},
		},
	}
}

func resourceTencentCloudEipAssociationCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_eip_association.create")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)
	vpcService := VpcService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}

	eipId := d.Get("eip_id").(string)
	var eip *vpc.Address
	var errRet error
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		eip, errRet = vpcService.DescribeEipById(ctx, eipId)
		if errRet != nil {
			return retryError(errRet, "InternalError")
		}
		if eip == nil {
			return resource.NonRetryableError(fmt.Errorf("eip is not found"))
		}
		return nil
	})
	if err != nil {
		return err
	}
	if *eip.AddressStatus != EIP_STATUS_UNBIND {
		return fmt.Errorf("eip status is illegal %s", *eip.AddressStatus)
	}

	if v, ok := d.GetOk("instance_id"); ok {
		instanceId := v.(string)
		err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			err := vpcService.AttachEip(ctx, eipId, instanceId)
			if err != nil {
				return retryError(err)
			}
			return nil
		})
		if err != nil {
			return err
		}
		associationId := fmt.Sprintf("%v::%v", eipId, instanceId)
		d.SetId(associationId)

		err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
			eip, errRet = vpcService.DescribeEipById(ctx, eipId)
			if errRet != nil {
				return retryError(errRet)
			}
			if eip == nil {
				return resource.NonRetryableError(fmt.Errorf("eip is not found"))
			}
			if *eip.AddressStatus == EIP_STATUS_BIND {
				return nil
			}
			return resource.RetryableError(fmt.Errorf("wait for binding success: %s", *eip.AddressStatus))
		})
		if err != nil {
			return err
		}
		return resourceTencentCloudEipAssociationRead(d, meta)
	}

	needRequest := false
	request := vpc.NewAssociateAddressRequest()
	request.AddressId = &eipId
	var networkId string
	var privateIp string
	if v, ok := d.GetOk("network_interface_id"); ok {
		needRequest = true
		networkId = v.(string)
		request.NetworkInterfaceId = &networkId
	}
	if v, ok := d.GetOk("private_ip"); ok {
		needRequest = true
		privateIp = v.(string)
		request.PrivateIpAddress = &privateIp
	}
	if needRequest {
		err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())
			response, err := meta.(*TencentCloudClient).apiV3Conn.UseVpcClient().AssociateAddress(request)
			if err != nil {
				log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
					logId, request.GetAction(), request.ToJsonString(), err.Error())
				return retryError(err)
			}
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
			return nil
		})
		if err != nil {
			return err
		}
		id := fmt.Sprintf("%v::%v::%v", eipId, networkId, privateIp)
		d.SetId(id)

		err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
			eip, errRet = vpcService.DescribeEipById(ctx, eipId)
			if errRet != nil {
				return retryError(errRet)
			}
			if eip == nil {
				return resource.NonRetryableError(fmt.Errorf("eip is not found"))
			}
			if *eip.AddressStatus == EIP_STATUS_BIND {
				return nil
			}
			return resource.RetryableError(fmt.Errorf("wait for binding success: %s", *eip.AddressStatus))
		})
		if err != nil {
			return err
		}

		return resourceTencentCloudEipAssociationRead(d, meta)
	}

	return nil
}

func resourceTencentCloudEipAssociationRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_eip_association.read")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)
	vpcService := VpcService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}

	id := d.Id()
	association, err := parseEipAssociationId(id)
	if err != nil {
		return err
	}

	err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		eip, errRet := vpcService.DescribeEipById(ctx, association.EipId)
		if errRet != nil {
			return retryError(errRet)
		}
		if eip == nil {
			d.SetId("")
		}
		return nil
	})
	if err != nil {
		return err
	}

	d.Set("eip_id", association.EipId)
	// associate with instance
	if len(association.InstanceId) > 0 {
		d.Set("instance_id", association.InstanceId)
		return nil
	}

	d.Set("network_interface_id", association.NetworkInterfaceId)
	d.Set("private_ip", association.PrivateIp)
	return nil
}

func resourceTencentCloudEipAssociationDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_eip_association.delete")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)
	vpcService := VpcService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}

	id := d.Id()
	association, err := parseEipAssociationId(id)
	if err != nil {
		return err
	}

	var eip *vpc.Address
	var errRet error
	err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		eip, errRet = vpcService.DescribeEipById(ctx, association.EipId)
		if errRet != nil {
			return retryError(errRet)
		}
		return nil
	})
	if err != nil {
		return err
	}
	if eip == nil {
		return nil
	}

	err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		err := vpcService.UnattachEip(ctx, association.EipId)
		if err != nil {
			return retryError(err)
		}
		return nil
	})
	if err != nil {
		return err
	}

	err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		eip, errRet = vpcService.DescribeEipById(ctx, association.EipId)
		if errRet != nil {
			return retryError(errRet)
		}
		if *eip.AddressStatus == EIP_STATUS_UNBIND {
			return nil
		}
		return resource.RetryableError(fmt.Errorf("wait for unbind success: %s", *eip.AddressStatus))
	})
	if err != nil {
		return err
	}

	return nil
}

type EipAssociationId struct {
	EipId              string
	InstanceId         string
	NetworkInterfaceId string
	PrivateIp          string
}

func parseEipAssociationId(associationId string) (association EipAssociationId, errRet error) {
	ids := strings.Split(associationId, "::")
	if len(ids) < 2 || len(ids) > 3 {
		errRet = fmt.Errorf("Invalid eip association ID: %v", associationId)
		return
	}
	association.EipId = ids[0]

	// associate with instance
	if len(ids) == 2 {
		association.InstanceId = ids[1]
		return
	}

	// associate with network interface
	association.NetworkInterfaceId = ids[1]
	association.PrivateIp = ids[2]
	return
}
