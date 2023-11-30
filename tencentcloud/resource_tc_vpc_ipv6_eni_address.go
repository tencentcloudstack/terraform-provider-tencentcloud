package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudVpcIpv6EniAddress() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudVpcIpv6EniAddressCreate,
		Read:   resourceTencentCloudVpcIpv6EniAddressRead,
		Update: resourceTencentCloudVpcIpv6EniAddressUpdate,
		Delete: resourceTencentCloudVpcIpv6EniAddressDelete,
		Schema: map[string]*schema.Schema{
			"vpc_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "VPC `ID`, in the form of `vpc-m6dyj72l`.",
			},

			"network_interface_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "ENI instance `ID`, in the form of `eni-m6dyj72l`.",
			},

			"ipv6_addresses": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "The specified `IPv6` address list, up to 10 can be specified at a time. Combined with the input parameter `Ipv6AddressCount` to calculate the quota. Mandatory one with Ipv6AddressCount.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"address": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "`IPv6` address, in the form of: `3402:4e00:20:100:0:8cd9:2a67:71f3`.",
						},
						"primary": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether to master `IP`.",
						},
						"address_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "`EIP` instance `ID`, such as:`eip-hxlqja90`.",
						},
						"description": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Description.",
						},
						"is_wan_ip_blocked": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether the public network IP is blocked.",
						},
						"state": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "`IPv6` address status: `PENDING`: pending, `MIGRATING`: migrating, `DELETING`: deleting, `AVAILABLE`: available.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudVpcIpv6EniAddressCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vpc_ipv6_eni_address.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request            = vpc.NewAssignIpv6AddressesRequest()
		response           = vpc.NewAssignIpv6AddressesResponse()
		vpcId              string
		networkInterfaceId string
	)
	if v, ok := d.GetOk("vpc_id"); ok {
		vpcId = v.(string)
	}

	if v, ok := d.GetOk("network_interface_id"); ok {
		networkInterfaceId = v.(string)
		request.NetworkInterfaceId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("ipv6_addresses"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			ipv6Address := vpc.Ipv6Address{}
			if v, ok := dMap["address"]; ok {
				ipv6Address.Address = helper.String(v.(string))
			}
			if v, ok := dMap["primary"]; ok {
				ipv6Address.Primary = helper.Bool(v.(bool))
			}
			if v, ok := dMap["address_id"]; ok {
				ipv6Address.AddressId = helper.String(v.(string))
			}
			if v, ok := dMap["description"]; ok {
				ipv6Address.Description = helper.String(v.(string))
			}
			if v, ok := dMap["is_wan_ip_blocked"]; ok {
				ipv6Address.IsWanIpBlocked = helper.Bool(v.(bool))
			}
			if v, ok := dMap["state"]; ok {
				ipv6Address.State = helper.String(v.(string))
			}
			request.Ipv6Addresses = append(request.Ipv6Addresses, &ipv6Address)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseVpcClient().AssignIpv6Addresses(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create vpc ipv6EniAddress failed, reason:%+v", logId, err)
		return err
	}

	addressSet := response.Response.Ipv6AddressSet
	if len(addressSet) < 1 {
		return fmt.Errorf("assign ipv6 addresses failed.")
	}

	time.Sleep(5 * time.Second)

	d.SetId(vpcId + FILED_SP + networkInterfaceId + FILED_SP + *addressSet[0].Address)

	return resourceTencentCloudVpcIpv6EniAddressRead(d, meta)
}

func resourceTencentCloudVpcIpv6EniAddressRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vpc_ipv6_eni_address.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	vpcId := idSplit[0]
	networkInterfaceId := idSplit[1]
	address := idSplit[2]

	ipv6EniAddress, err := service.DescribeVpcIpv6EniAddressById(ctx, vpcId, address)
	if err != nil {
		return err
	}

	if ipv6EniAddress == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `VpcIpv6EniAddress` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("vpc_id", vpcId)
	_ = d.Set("network_interface_id", networkInterfaceId)

	return nil
}

func resourceTencentCloudVpcIpv6EniAddressUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vpc_ipv6_eni_address.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := vpc.NewModifyIpv6AddressesAttributeRequest()

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	networkInterfaceId := idSplit[1]

	request.NetworkInterfaceId = &networkInterfaceId

	if d.HasChange("ipv6_addresses") {
		if v, ok := d.GetOk("ipv6_addresses"); ok {
			for _, item := range v.([]interface{}) {
				ipv6Address := vpc.Ipv6Address{}
				dMap := item.(map[string]interface{})
				if v, ok := dMap["address"]; ok {
					ipv6Address.Address = helper.String(v.(string))
				}
				if v, ok := dMap["primary"]; ok {
					ipv6Address.Primary = helper.Bool(v.(bool))
				}
				if v, ok := dMap["address_id"]; ok {
					ipv6Address.AddressId = helper.String(v.(string))
				}
				if v, ok := dMap["description"]; ok {
					ipv6Address.Description = helper.String(v.(string))
				}
				if v, ok := dMap["is_wan_ip_blocked"]; ok {
					ipv6Address.IsWanIpBlocked = helper.Bool(v.(bool))
				}
				if v, ok := dMap["state"]; ok {
					ipv6Address.State = helper.String(v.(string))
				}
				request.Ipv6Addresses = append(request.Ipv6Addresses, &ipv6Address)
			}
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseVpcClient().ModifyIpv6AddressesAttribute(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update vpc ipv6EniAddress failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudVpcIpv6EniAddressRead(d, meta)
}

func resourceTencentCloudVpcIpv6EniAddressDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vpc_ipv6_eni_address.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}
	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	networkInterfaceId := idSplit[1]
	address := idSplit[2]

	if err := service.DeleteVpcIpv6EniAddressById(ctx, networkInterfaceId, address); err != nil {
		return err
	}

	return nil
}
