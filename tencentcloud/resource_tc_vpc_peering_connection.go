/*
Provides a resource to create a vpc peering_connection

Example Usage

```hcl
resource "tencentcloud_vpc_peering_connection" "peering_connection" {
  source_vpc_id = "vpc-abcdef"
  peering_connection_name = "name"
  destination_vpc_id = "vpc-abc1234"
  destination_uin = "12345678"
  destination_region = "ap-beijing"
  bandwidth = 100
  type = "VPC_PEER"
  charge_type = "POSTPAID_BY_DAY_MAX"
  qos_level = "AU"
}
```

Import

vpc peering_connection can be imported using the id, e.g.

```
terraform import tencentcloud_vpc_peering_connection.peering_connection peering_connection_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudVpcPeeringConnection() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudVpcPeeringConnectionCreate,
		Read:   resourceTencentCloudVpcPeeringConnectionRead,
		Update: resourceTencentCloudVpcPeeringConnectionUpdate,
		Delete: resourceTencentCloudVpcPeeringConnectionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"source_vpc_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The unique ID of the local VPC.",
			},

			"peering_connection_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Peer connection name.",
			},

			"destination_vpc_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The unique ID of the peer VPC.",
			},

			"destination_uin": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Peer user UIN.",
			},

			"destination_region": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Peer region.",
			},

			"bandwidth": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Bandwidth upper limit, unit Mbps.",
			},

			"type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Interworking type, `VPC_PEER`: interworking between VPCs, `VPC_BM_PEER`: interworking between VPC and BM Network.",
			},

			"charge_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Billing mode, daily peak value: `POSTPAID_BY_DAY_MAX`, monthly 95 value: `POSTPAID_BY_MONTH_95`.",
			},

			"qos_level": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Service classification PT, AU, AG.",
			},
		},
	}
}

func resourceTencentCloudVpcPeeringConnectionCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vpc_peering_connection.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request             = vpc.NewCreateVpcPeeringConnectionRequest()
		response            = vpc.NewCreateVpcPeeringConnectionResponse()
		peeringConnectionId string
	)
	if v, ok := d.GetOk("source_vpc_id"); ok {
		request.SourceVpcId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("peering_connection_name"); ok {
		request.PeeringConnectionName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("destination_vpc_id"); ok {
		request.DestinationVpcId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("destination_uin"); ok {
		request.DestinationUin = helper.String(v.(string))
	}

	if v, ok := d.GetOk("destination_region"); ok {
		request.DestinationRegion = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("bandwidth"); ok {
		request.Bandwidth = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("type"); ok {
		request.Type = helper.String(v.(string))
	}

	if v, ok := d.GetOk("charge_type"); ok {
		request.ChargeType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("qos_level"); ok {
		request.QosLevel = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseVpcClient().CreateVpcPeeringConnection(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create vpc PeeringConnection failed, reason:%+v", logId, err)
		return err
	}

	peeringConnectionId = *response.Response.PeeringConnectionId
	d.SetId(peeringConnectionId)

	return resourceTencentCloudVpcPeeringConnectionRead(d, meta)
}

func resourceTencentCloudVpcPeeringConnectionRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vpc_peering_connection.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}

	peeringConnectionId := d.Id()

	peeringConnection, err := service.DescribeVpcPeeringConnectionById(ctx, peeringConnectionId)
	if err != nil {
		return err
	}

	if peeringConnection == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `VpcPeeringConnection` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if peeringConnection.SourceVpcId != nil {
		_ = d.Set("source_vpc_id", peeringConnection.SourceVpcId)
	}

	if peeringConnection.PeeringConnectionName != nil {
		_ = d.Set("peering_connection_name", peeringConnection.PeeringConnectionName)
	}

	if peeringConnection.PeerVpcId != nil {
		_ = d.Set("destination_vpc_id", peeringConnection.PeerVpcId)
	}

	if peeringConnection.DestinationUin != nil {
		_ = d.Set("destination_uin", helper.Int64ToStr(*peeringConnection.DestinationUin))
	}

	if peeringConnection.DestinationRegion != nil {
		_ = d.Set("destination_region", peeringConnection.DestinationRegion)
	}

	if peeringConnection.Bandwidth != nil {
		_ = d.Set("bandwidth", peeringConnection.Bandwidth)
	}

	if peeringConnection.Type != nil {
		_ = d.Set("type", peeringConnection.Type)
	}

	if peeringConnection.ChargeType != nil {
		_ = d.Set("charge_type", peeringConnection.ChargeType)
	}

	if peeringConnection.QosLevel != nil {
		_ = d.Set("qos_level", peeringConnection.QosLevel)
	}

	return nil
}

func resourceTencentCloudVpcPeeringConnectionUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vpc_peering_connection.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := vpc.NewModifyVpcPeeringConnectionRequest()

	peeringConnectionId := d.Id()

	request.PeeringConnectionId = &peeringConnectionId

	immutableArgs := []string{"source_vpc_id", "peering_connection_name", "destination_vpc_id", "destination_uin", "destination_region", "bandwidth", "type", "charge_type", "qos_level"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("peering_connection_name") {
		if v, ok := d.GetOk("peering_connection_name"); ok {
			request.PeeringConnectionName = helper.String(v.(string))
		}
	}

	if d.HasChange("bandwidth") {
		if v, ok := d.GetOkExists("bandwidth"); ok {
			request.Bandwidth = helper.IntInt64(v.(int))
		}
	}

	if d.HasChange("charge_type") {
		if v, ok := d.GetOk("charge_type"); ok {
			request.ChargeType = helper.String(v.(string))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseVpcClient().ModifyVpcPeeringConnection(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update vpc PeeringConnection failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudVpcPeeringConnectionRead(d, meta)
}

func resourceTencentCloudVpcPeeringConnectionDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vpc_peering_connection.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}
	peeringConnectionId := d.Id()

	if err := service.DeleteVpcPeeringConnectionById(ctx, peeringConnectionId); err != nil {
		return err
	}

	return nil
}
