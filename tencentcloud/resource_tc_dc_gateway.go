/*
Provides a resource to creating direct connect gateway instance.

Example Usage

```hcl
resource "tencentcloud_vpc" "main" {
    name="ci-vpc-instance-test"
    cidr_block="10.0.0.0/16"
}

resource "tencentcloud_dc_gateway" "vpc_main" {
  name                = "ci-cdg-vpc-test"
  network_instance_id = tencentcloud_vpc.main.id
  network_type        = "VPC"
  gateway_type        = "NAT"
}
```

Import

Direct connect gateway instance can be imported, e.g.

```
$ terraform import tencentcloud_dc_gateway.instance dcg-id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceTencentCloudDcGatewayInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDcGatewayCreate,
		Read:   resourceTencentCloudDcGatewayRead,
		Update: resourceTencentCloudDcGatewayUpdate,
		Delete: resourceTencentCloudDcGatewayDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateStringLengthInRange(1, 60),
				Description:  "Name of the DCG.",
			},
			"network_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateAllowedStringValue(DCG_NETWORK_TYPES),
				Description:  "Type of associated network, the available value include 'VPC' and 'CCN'.",
			},
			"network_instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "If the 'network_type' value is 'VPC', the available value is VPC ID. But when the 'network_type' value is 'CCN', the available value is CCN instance ID.",
			},
			"gateway_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      DCG_GATEWAY_TYPE_NORMAL,
				ValidateFunc: validateAllowedStringValue(DCG_GATEWAY_TYPES),
				Description:  "Type of the gateway, the available value include 'NORMAL' and 'NAT'. Default is 'NORMAL'. NOTES: CCN only supports 'NORMAL' and a vpc can create two DCGs, the one is NAT type and the other is non-NAT type.",
			},

			//compute
			"cnn_route_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Type of CCN route, the available value include 'BGP' and 'STATIC'. The property is available when the DCG type is CCN gateway and BGP enabled.",
			},
			"enable_bgp": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether the BGP is enabled.",
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation time of resource.",
			},
		},
	}
}

func resourceTencentCloudDcGatewayCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dc_gateway.create")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	service := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}

	var (
		name              = d.Get("name").(string)
		networkType       = d.Get("network_type").(string)
		networkInstanceId = d.Get("network_instance_id").(string)
		gatewayType       = d.Get("gateway_type").(string)
	)

	if networkType == DCG_NETWORK_TYPE_VPC &&
		!strings.HasPrefix(networkInstanceId, "vpc") {

		return fmt.Errorf("if `network_type` is '%s', the field `network_instance_id` must be a VPC resource",
			DCG_NETWORK_TYPE_VPC)
	}

	if networkType == DCG_NETWORK_TYPE_CCN &&
		!strings.HasPrefix(networkInstanceId, "ccn") {

		return fmt.Errorf("if `network_type` is '%s', the field `network_instance_id` must be a CCN resource",
			DCG_NETWORK_TYPE_CCN)
	}

	if networkType == DCG_NETWORK_TYPE_CCN && gatewayType != DCG_GATEWAY_TYPE_NORMAL {

		return fmt.Errorf("if `network_type` is '%s', the field `gateway_type` must be '%s'",
			DCG_NETWORK_TYPE_CCN,
			DCG_GATEWAY_TYPE_NORMAL)
	}

	dcgId, err := service.CreateDirectConnectGateway(ctx, name, networkType, networkInstanceId, gatewayType)
	if err != nil {
		return err
	}

	d.SetId(dcgId)

	return resourceTencentCloudDcGatewayRead(d, meta)
}

func resourceTencentCloudDcGatewayRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dc_gateway.read")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	service := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		info, has, e := service.DescribeDirectConnectGateway(ctx, d.Id())
		if e != nil {
			return retryError(e)
		}

		if has == 0 {
			d.SetId("")
			return nil
		}

		_ = d.Set("name", info.name)
		_ = d.Set("network_type", info.networkType)
		_ = d.Set("network_instance_id", info.networkInstanceId)
		_ = d.Set("gateway_type", info.gatewayType)
		_ = d.Set("cnn_route_type", info.cnnRouteType)
		_ = d.Set("enable_bgp", info.enableBGP)
		_ = d.Set("create_time", info.createTime)
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func resourceTencentCloudDcGatewayUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dc_gateway.update")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	service := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}
	if d.HasChange("name") {
		var name = d.Get("name").(string)
		return service.ModifyDirectConnectGatewayAttribute(ctx, d.Id(), name)
	}

	return resourceTencentCloudDcGatewayRead(d, meta)
}

func resourceTencentCloudDcGatewayDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dc_gateway.delete")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	service := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		_, has, e := service.DescribeDirectConnectGateway(ctx, d.Id())
		if e != nil {
			return retryError(e)
		}

		if has == 0 {
			return nil
		}
		return nil
	})
	if err != nil {
		return err
	}
	return service.DeleteDirectConnectGateway(ctx, d.Id())
}
