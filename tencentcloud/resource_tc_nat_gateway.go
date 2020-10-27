/*
Provides a resource to create a NAT gateway.

Example Usage

```hcl
resource "tencentcloud_nat_gateway" "foo" {
  name             = "test_nat_gateway"
  vpc_id           = "vpc-4xxr2cy7"
  bandwidth        = 100
  max_concurrent   = 1000000
  assigned_eip_set = ["1.1.1.1"]

  tags = {
    test = "tf"
  }
}
```

Import

NAT gateway can be imported using the id, e.g.

```
$ terraform import tencentcloud_nat_gateway.foo nat-1asg3t63
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudNatGateway() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudNatGatewayCreate,
		Read:   resourceTencentCloudNatGatewayRead,
		Update: resourceTencentCloudNatGatewayUpdate,
		Delete: resourceTencentCloudNatGatewayDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the vpc.",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateStringLengthInRange(1, 60),
				Description:  "Name of the NAT gateway.",
			},
			"max_concurrent": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      1000000,
				ValidateFunc: validateAllowedIntValue([]int{1000000, 3000000, 10000000}),
				Description:  "The upper limit of concurrent connection of NAT gateway, the available values include: 1000000,3000000,10000000. Default is 1000000.",
			},
			"bandwidth": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     100,
				Description: "The maximum public network output bandwidth of NAT gateway (unit: Mbps), the available values include: 20,50,100,200,500,1000,2000,5000. Default is 100.",
			},
			"assigned_eip_set": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validateIp,
				},
				MinItems:    1,
				MaxItems:    10,
				Description: "EIP IP address set bound to the gateway. The value of at least 1 and at most 10.",
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "The available tags within this NAT gateway.",
			},
			//computed
			"created_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Create time of the NAT gateway.",
			},
		},
	}
}

func resourceTencentCloudNatGatewayCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_nat_gateway.create")()

	logId := getLogId(contextNil)
	request := vpc.NewCreateNatGatewayRequest()
	vpcId := d.Get("vpc_id").(string)
	natGatewayName := d.Get("name").(string)
	request.VpcId = &vpcId
	request.NatGatewayName = &natGatewayName
	//test default value
	bandwidth := uint64(d.Get("bandwidth").(int))
	request.InternetMaxBandwidthOut = &bandwidth
	maxConcurrent := uint64(d.Get("max_concurrent").(int))
	request.MaxConcurrentConnection = &maxConcurrent
	if v, ok := d.GetOk("assigned_eip_set"); ok {
		eipSet := v.(*schema.Set).List()
		//set request public ips
		for i := range eipSet {
			publicIp := eipSet[i].(string)
			request.PublicIpAddresses = append(request.PublicIpAddresses, &publicIp)
		}
	}

	var response *vpc.CreateNatGatewayResponse
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseVpcClient().CreateNatGateway(request)
		if e != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), e.Error())
			return retryError(e)
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create NAT gateway failed, reason:%s\n", logId, err.Error())
		return err
	}

	if len(response.Response.NatGatewaySet) < 1 {
		return fmt.Errorf("NAT gateway ID is nil")
	}
	d.SetId(*response.Response.NatGatewaySet[0].NatGatewayId)
	// must wait for finishing creating NAT
	statRequest := vpc.NewDescribeNatGatewaysRequest()
	statRequest.NatGatewayIds = []*string{response.Response.NatGatewaySet[0].NatGatewayId}
	err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseVpcClient().DescribeNatGateways(statRequest)
		if e != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), e.Error())
			return retryError(e)
		} else {
			//if not, quit
			if len(result.Response.NatGatewaySet) != 1 {
				return resource.NonRetryableError(fmt.Errorf("creating error"))
			}
			//else get stat
			nat := result.Response.NatGatewaySet[0]
			stat := *nat.State

			if stat == "AVAILABLE" {
				return nil
			}
			return resource.RetryableError(fmt.Errorf("creating not ready retry"))
		}
	})
	if err != nil {
		log.Printf("[CRITAL]%s create NAT gateway failed, reason:%s\n", logId, err.Error())
		return err
	}

	//cs::vpc:ap-guangzhou:uin/12345:nat/nat-nxxx
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		tcClient := meta.(*TencentCloudClient).apiV3Conn
		tagService := &TagService{client: tcClient}
		resourceName := BuildTagResourceName("vpc", "nat", tcClient.Region, d.Id())
		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
	}
	return resourceTencentCloudNatGatewayRead(d, meta)
}

func resourceTencentCloudNatGatewayRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_nat_gateway.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	natGatewayId := d.Id()
	request := vpc.NewDescribeNatGatewaysRequest()
	request.NatGatewayIds = []*string{&natGatewayId}
	var response *vpc.DescribeNatGatewaysResponse
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseVpcClient().DescribeNatGateways(request)
		if e != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), e.Error())
			return retryError(e)
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read NAT gateway failed, reason:%s\n", logId, err.Error())
		return err
	}
	if len(response.Response.NatGatewaySet) < 1 {
		d.SetId("")
		return nil
	}

	nat := response.Response.NatGatewaySet[0]

	_ = d.Set("vpc_id", *nat.VpcId)
	_ = d.Set("name", *nat.NatGatewayName)
	_ = d.Set("max_concurrent", *nat.MaxConcurrentConnection)
	_ = d.Set("bandwidth", *nat.InternetMaxBandwidthOut)
	_ = d.Set("create_time", *nat.CreatedTime)
	_ = d.Set("assigned_eip_set", flattenAddressList((*nat).PublicIpAddressSet))

	tcClient := meta.(*TencentCloudClient).apiV3Conn
	tagService := &TagService{client: tcClient}
	tags, err := tagService.DescribeResourceTags(ctx, "vpc", "nat", tcClient.Region, d.Id())
	if err != nil {
		return err
	}
	_ = d.Set("tags", tags)

	return nil
}

func resourceTencentCloudNatGatewayUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_nat_gateway.update")()

	logId := getLogId(contextNil)

	d.Partial(true)
	natGatewayId := d.Id()
	request := vpc.NewModifyNatGatewayAttributeRequest()
	request.NatGatewayId = &natGatewayId
	changed := false
	if d.HasChange("name") {
		request.NatGatewayName = helper.String(d.Get("name").(string))
		changed = true
	}
	if d.HasChange("bandwidth") {
		bandwidth := d.Get("bandwidth").(int)
		bandwidth64 := uint64(bandwidth)
		request.InternetMaxBandwidthOut = &bandwidth64
		changed = true
	}
	if changed {
		err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
			_, e := meta.(*TencentCloudClient).apiV3Conn.UseVpcClient().ModifyNatGatewayAttribute(request)
			if e != nil {
				log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
					logId, request.GetAction(), request.ToJsonString(), e.Error())
				return retryError(e)
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s modify NAT gateway failed, reason:%s\n", logId, err.Error())
			return err
		}
	}
	if d.HasChange("name") {
		d.SetPartial("name")
	}
	if d.HasChange("bandwidth") {
		d.SetPartial("bandwidth")
	}
	//max concurrent
	if d.HasChange("max_concurrent") {
		concurrentReq := vpc.NewResetNatGatewayConnectionRequest()
		concurrentReq.NatGatewayId = &natGatewayId
		concurrent := d.Get("max_concurrent").(int)
		concurrent64 := uint64(concurrent)
		concurrentReq.MaxConcurrentConnection = &concurrent64
		err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
			_, e := meta.(*TencentCloudClient).apiV3Conn.UseVpcClient().ResetNatGatewayConnection(concurrentReq)
			if e != nil {
				log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
					logId, concurrentReq.GetAction(), concurrentReq.ToJsonString(), e.Error())
				return retryError(e, InternalError)
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s modify NAT gateway concurrent failed, reason:%s\n", logId, err.Error())
			return err
		}
		d.SetPartial("max_concurrent")
	}

	//eip

	if d.HasChange("assigned_eip_set") {
		eipSetLength := 0
		if v, ok := d.GetOk("assigned_eip_set"); ok {
			eipSet := v.(*schema.Set).List()
			eipSetLength = len(eipSet)
		}
		if d.HasChange("assigned_eip_set") {
			o, n := d.GetChange("assigned_eip_set")
			os := o.(*schema.Set)
			ns := n.(*schema.Set)
			oldEipSet := os.List()
			newEipSet := ns.List()

			//in case of no union set
			backUpOldIp := ""
			backUpNewIp := ""
			//Unassign eips
			if len(oldEipSet) > 0 {
				unassignedRequest := vpc.NewDisassociateNatGatewayAddressRequest()
				unassignedRequest.PublicIpAddresses = make([]*string, 0, len(oldEipSet))
				unassignedRequest.NatGatewayId = &natGatewayId
				//set request public ips
				for i := range oldEipSet {
					publicIp := oldEipSet[i].(string)
					isIn := false
					for j := range newEipSet {
						if publicIp == newEipSet[j] {
							isIn = true
						}
					}
					if !isIn {
						if len(unassignedRequest.PublicIpAddresses)+1 == len(oldEipSet) {
							backUpOldIp = publicIp
						} else {
							unassignedRequest.PublicIpAddresses = append(unassignedRequest.PublicIpAddresses, &publicIp)
						}
					}
				}

				if len(unassignedRequest.PublicIpAddresses) > 0 {
					err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
						_, e := meta.(*TencentCloudClient).apiV3Conn.UseVpcClient().DisassociateNatGatewayAddress(unassignedRequest)
						if e != nil {
							log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
								logId, unassignedRequest.GetAction(), unassignedRequest.ToJsonString(), e.Error())
							return retryError(e)
						}
						return nil
					})
					if err != nil {
						log.Printf("[CRITAL]%s modify NAT gateway EIP failed, reason:%s\n", logId, err.Error())
						return err
					}
				}
			}
			time.Sleep(3 * time.Minute)
			//Assign new EIP
			if len(newEipSet) > 0 {
				assignedRequest := vpc.NewAssociateNatGatewayAddressRequest()
				assignedRequest.PublicIpAddresses = make([]*string, 0, len(newEipSet))
				assignedRequest.NatGatewayId = &natGatewayId
				//set request public ips
				for i := range newEipSet {
					publicIp := newEipSet[i].(string)
					isIn := false
					for j := range oldEipSet {
						if publicIp == oldEipSet[j] {
							isIn = true
						}
					}
					if !isIn {
						if len(assignedRequest.PublicIpAddresses)+eipSetLength+1 == NAT_EIP_MAX_LIMIT {
							backUpNewIp = publicIp
						} else {
							assignedRequest.PublicIpAddresses = append(assignedRequest.PublicIpAddresses, &publicIp)
						}
					}
				}
				if len(assignedRequest.PublicIpAddresses) > 0 {
					err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
						_, e := meta.(*TencentCloudClient).apiV3Conn.UseVpcClient().AssociateNatGatewayAddress(assignedRequest)
						if e != nil {
							log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
								logId, assignedRequest.GetAction(), assignedRequest.ToJsonString(), e.Error())
							return retryError(e)
						}
						return nil
					})
					if err != nil {
						log.Printf("[CRITAL]%s modify NAT gateway EIP failed, reason:%s\n", logId, err.Error())
						return err
					}
				}
			}
			time.Sleep(3 * time.Minute)
			if backUpOldIp != "" {
				//disassociate one old ip
				unassignedRequest := vpc.NewDisassociateNatGatewayAddressRequest()
				unassignedRequest.NatGatewayId = &natGatewayId
				unassignedRequest.PublicIpAddresses = []*string{&backUpOldIp}
				err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
					_, e := meta.(*TencentCloudClient).apiV3Conn.UseVpcClient().DisassociateNatGatewayAddress(unassignedRequest)
					if e != nil {
						log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
							logId, unassignedRequest.GetAction(), unassignedRequest.ToJsonString(), e.Error())
						return retryError(e)
					}
					return nil
				})
				if err != nil {
					log.Printf("[CRITAL]%s modify NAT gateway EIP failed, reason:%s\n", logId, err.Error())
					return err
				}
			}
			if backUpNewIp != "" {
				//associate one new ip
				assignedRequest := vpc.NewAssociateNatGatewayAddressRequest()
				assignedRequest.NatGatewayId = &natGatewayId
				assignedRequest.PublicIpAddresses = []*string{&backUpNewIp}
				err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
					_, e := meta.(*TencentCloudClient).apiV3Conn.UseVpcClient().AssociateNatGatewayAddress(assignedRequest)
					if e != nil {
						log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
							logId, assignedRequest.GetAction(), assignedRequest.ToJsonString(), e.Error())
						return retryError(e)
					}
					return nil
				})
				if err != nil {
					log.Printf("[CRITAL]%s modify NAT gateway EIP failed, reason:%s\n", logId, err.Error())
					return err
				}
			}
			d.SetPartial("assigned_eip_set")
		}

	}

	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	if d.HasChange("tags") {

		oldValue, newValue := d.GetChange("tags")
		replaceTags, deleteTags := diffTags(oldValue.(map[string]interface{}), newValue.(map[string]interface{}))

		tcClient := meta.(*TencentCloudClient).apiV3Conn
		tagService := &TagService{client: tcClient}
		resourceName := BuildTagResourceName("vpc", "nat", tcClient.Region, d.Id())
		err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags)
		if err != nil {
			return err
		}
		d.SetPartial("tags")
	}

	d.Partial(false)

	return nil
}

func resourceTencentCloudNatGatewayDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_nat_gateway.delete")()

	logId := getLogId(contextNil)

	natGatewayId := d.Id()
	request := vpc.NewDeleteNatGatewayRequest()
	request.NatGatewayId = &natGatewayId
	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		_, e := meta.(*TencentCloudClient).apiV3Conn.UseVpcClient().DeleteNatGateway(request)
		if e != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), e.Error())
			return retryError(e)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s delete NAT gateway failed, reason:%s\n", logId, err.Error())
		return err
	}
	// must wait for finishing deleting NAT
	time.Sleep(10 * time.Second)
	//to get the status of NAT

	statRequest := vpc.NewDescribeNatGatewaysRequest()
	statRequest.NatGatewayIds = []*string{&natGatewayId}
	err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseVpcClient().DescribeNatGateways(statRequest)
		if e != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), e.Error())
			return retryError(e)
		} else {
			//if not, quit
			if len(result.Response.NatGatewaySet) == 0 {
				log.Printf("deleting done")
				return nil
			}
			//else get stat
			nat := result.Response.NatGatewaySet[0]
			stat := *nat.State
			if stat == NAT_FAILED_STATE {
				return resource.NonRetryableError(fmt.Errorf("delete NAT failed"))
			}
			time.Sleep(3 * time.Second)

			return resource.RetryableError(fmt.Errorf("deleting retry"))
		}
	})
	if err != nil {
		log.Printf("[CRITAL]%s delete NAT gateway failed, reason:%s\n", logId, err.Error())
		return err
	}
	return nil
}

func flattenAddressList(addresses []*vpc.NatGatewayAddress) (eips []*string) {
	for _, address := range addresses {
		eips = append(eips, address.PublicIpAddress)
	}
	return
}
