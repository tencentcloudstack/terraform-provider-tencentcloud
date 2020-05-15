/*
Provides a resource to create a VPN customer gateway.

Example Usage

```hcl
resource "tencentcloud_vpn_customer_gateway" "foo" {
  name              = "test_vpn_customer_gateway"
  public_ip_address = "1.1.1.1"

  tags = {
	  tag = "test"
  }
}
```

Import

VPN customer gateway can be imported using the id, e.g.

```
$ terraform import tencentcloud_vpn_customer_gateway.foo cgw-xfqag
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
	"github.com/terraform-providers/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudVpnCustomerGateway() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudVpnCustomerGatewayCreate,
		Read:   resourceTencentCloudVpnCustomerGatewayRead,
		Update: resourceTencentCloudVpnCustomerGatewayUpdate,
		Delete: resourceTencentCloudVpnCustomerGatewayDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateStringLengthInRange(1, 60),
				Description:  "Name of the customer gateway. The length of character is limited to 1-60.",
			},
			"public_ip_address": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateIp,
				Description:  "Public ip of the customer gateway.",
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "A list of tags used to associate different resources.",
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Create time of the customer gateway.",
			},
		},
	}
}

func resourceTencentCloudVpnCustomerGatewayCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vpn_customer_gateway.create")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	request := vpc.NewCreateCustomerGatewayRequest()
	request.CustomerGatewayName = helper.String(d.Get("name").(string))
	request.IpAddress = helper.String(d.Get("public_ip_address").(string))
	var response *vpc.CreateCustomerGatewayResponse
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseVpcClient().CreateCustomerGateway(request)
		if e != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), e.Error())
			return retryError(e)
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create VPN customer gateway failed, reason:%s\n", logId, err.Error())
		return err
	}

	if response.Response.CustomerGateway == nil {
		return fmt.Errorf("VPN customer gateway id is nil")
	}
	customerGatewayId := *response.Response.CustomerGateway.CustomerGatewayId
	d.SetId(customerGatewayId)
	// must wait for finishing creating customer gateway
	statRequest := vpc.NewDescribeCustomerGatewaysRequest()
	statRequest.CustomerGatewayIds = []*string{response.Response.CustomerGateway.CustomerGatewayId}
	err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseVpcClient().DescribeCustomerGateways(statRequest)
		if e != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), e.Error())
			return retryError(e)
		} else {
			//if not, quit
			if len(result.Response.CustomerGatewaySet) != 1 {
				return resource.NonRetryableError(fmt.Errorf("creating error"))
			}
			//else consider created, cos there is no status of gateway
			return nil
		}
	})
	if err != nil {
		log.Printf("[CRITAL]%s create VPN customer gateway failed, reason:%s\n", logId, err.Error())
		return err
	}

	//modify tags
	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		tagService := TagService{client: meta.(*TencentCloudClient).apiV3Conn}

		region := meta.(*TencentCloudClient).apiV3Conn.Region
		resourceName := BuildTagResourceName("vpc", "cgw", region, customerGatewayId)

		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
	}

	return resourceTencentCloudVpnCustomerGatewayRead(d, meta)
}

func resourceTencentCloudVpnCustomerGatewayRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vpn_customer_gateway.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	customerGatewayId := d.Id()
	request := vpc.NewDescribeCustomerGatewaysRequest()
	request.CustomerGatewayIds = []*string{&customerGatewayId}
	var response *vpc.DescribeCustomerGatewaysResponse
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseVpcClient().DescribeCustomerGateways(request)
		if e != nil {
			ee, ok := e.(*errors.TencentCloudSDKError)
			if !ok {
				return retryError(e)
			}
			if ee.Code == VPCNotFound {
				log.Printf("[CRITAL]%s api[%s] success, request body [%s], reason[%s]\n",
					logId, request.GetAction(), request.ToJsonString(), e.Error())
				return nil
			} else {
				log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
					logId, request.GetAction(), request.ToJsonString(), e.Error())
				return retryError(e)
			}
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read VPN customer gateway failed, reason:%s\n", logId, err.Error())
		return err
	}
	if len(response.Response.CustomerGatewaySet) < 1 {
		d.SetId("")
		return nil
	}

	gateway := response.Response.CustomerGatewaySet[0]

	_ = d.Set("name", *gateway.CustomerGatewayName)
	_ = d.Set("public_ip_address", *gateway.IpAddress)
	_ = d.Set("create_time", *gateway.CreatedTime)

	//tags
	tagService := TagService{client: meta.(*TencentCloudClient).apiV3Conn}
	region := meta.(*TencentCloudClient).apiV3Conn.Region
	tags, err := tagService.DescribeResourceTags(ctx, "vpc", "cgw", region, customerGatewayId)
	if err != nil {
		return err
	}
	_ = d.Set("tags", tags)

	return nil
}

func resourceTencentCloudVpnCustomerGatewayUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vpn_customer_gateway.update")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	d.Partial(true)
	customerGatewayId := d.Id()
	request := vpc.NewModifyCustomerGatewayAttributeRequest()
	request.CustomerGatewayId = &customerGatewayId
	if d.HasChange("name") {
		request.CustomerGatewayName = helper.String(d.Get("name").(string))
		err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
			_, e := meta.(*TencentCloudClient).apiV3Conn.UseVpcClient().ModifyCustomerGatewayAttribute(request)
			if e != nil {
				log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
					logId, request.GetAction(), request.ToJsonString(), e.Error())
				return retryError(e)
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s modify VPN customer gateway failed, reason:%s\n", logId, err.Error())
			return err
		}
		d.SetPartial("name")
	}

	//tag
	if d.HasChange("tags") {
		oldInterface, newInterface := d.GetChange("tags")
		replaceTags, deleteTags := diffTags(oldInterface.(map[string]interface{}), newInterface.(map[string]interface{}))
		tagService := TagService{
			client: meta.(*TencentCloudClient).apiV3Conn,
		}
		region := meta.(*TencentCloudClient).apiV3Conn.Region
		resourceName := BuildTagResourceName("vpc", "cgw", region, customerGatewayId)
		err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags)
		if err != nil {
			return err
		}
		d.SetPartial("tags")
	}
	d.Partial(false)

	return resourceTencentCloudVpnCustomerGatewayRead(d, meta)
}

func resourceTencentCloudVpnCustomerGatewayDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vpn_customer_gateway.delete")()

	logId := getLogId(contextNil)

	customerGatewayId := d.Id()

	//check the customer gateway is not related with any tunnel
	tRequest := vpc.NewDescribeVpnConnectionsRequest()
	tRequest.Filters = make([]*vpc.Filter, 0, 1)
	params := make(map[string]string)
	params["customer-gateway-id"] = customerGatewayId

	for k, v := range params {
		filter := &vpc.Filter{
			Name:   helper.String(k),
			Values: []*string{helper.String(v)},
		}
		tRequest.Filters = append(tRequest.Filters, filter)
	}
	offset := uint64(0)
	tRequest.Offset = &offset

	tErr := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseVpcClient().DescribeVpnConnections(tRequest)

		if e != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, tRequest.GetAction(), tRequest.ToJsonString(), e.Error())
			return retryError(e)
		} else {
			if len(result.Response.VpnConnectionSet) == 0 {
				return nil
			} else {
				return resource.NonRetryableError(fmt.Errorf("There is associated tunnel exists, please delete associated tunnels first."))
			}
		}
	})
	if tErr != nil {
		log.Printf("[CRITAL]%s describe VPN connection failed, reason:%s\n", logId, tErr.Error())
		return tErr
	}

	request := vpc.NewDeleteCustomerGatewayRequest()
	request.CustomerGatewayId = &customerGatewayId
	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		_, e := meta.(*TencentCloudClient).apiV3Conn.UseVpcClient().DeleteCustomerGateway(request)
		if e != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), e.Error())
			return retryError(e)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s delete VPN customer gateway failed, reason:%s\n", logId, err.Error())
		return err
	}
	//to get the status of customer gateway
	statRequest := vpc.NewDescribeCustomerGatewaysRequest()
	statRequest.CustomerGatewayIds = []*string{&customerGatewayId}
	err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseVpcClient().DescribeCustomerGateways(statRequest)
		if e != nil {
			ee, ok := e.(*errors.TencentCloudSDKError)
			if !ok {
				return retryError(e)
			}
			if ee.Code == VPCNotFound {
				log.Printf("[CRITAL]%s api[%s] success, request body [%s], reason[%s]\n",
					logId, request.GetAction(), request.ToJsonString(), e.Error())
				return nil
			} else {
				log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
					logId, request.GetAction(), request.ToJsonString(), e.Error())
				return retryError(e)
			}
		} else {
			//if not, quit
			if len(result.Response.CustomerGatewaySet) == 0 {
				return nil
			}
			//else consider delete fail
			return resource.RetryableError(fmt.Errorf("deleting retry"))
		}
	})
	if err != nil {
		log.Printf("[CRITAL]%s delete VPN customer gateway failed, reason:%s\n", logId, err.Error())
		return err
	}
	return nil
}
