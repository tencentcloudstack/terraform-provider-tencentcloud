package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudDcGatewayAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDcGatewayAttachmentCreate,
		Read:   resourceTencentCloudDcGatewayAttachmentRead,
		Delete: resourceTencentCloudDcGatewayAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"vpc_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "vpc id.",
			},

			"nat_gateway_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "NatGatewayId.",
			},

			"direct_connect_gateway_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "DirectConnectGatewayId.",
			},
		},
	}
}

func resourceTencentCloudDcGatewayAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dc_gateway_attachment.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request                = vpc.NewAssociateDirectConnectGatewayNatGatewayRequest()
		vpcId                  string
		directConnectGatewayId string
		natGatewayId           string
	)
	if v, ok := d.GetOk("vpc_id"); ok {
		vpcId = v.(string)
		request.VpcId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("nat_gateway_id"); ok {
		natGatewayId = v.(string)
		request.NatGatewayId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("direct_connect_gateway_id"); ok {
		directConnectGatewayId = v.(string)
		request.DirectConnectGatewayId = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseVpcClient().AssociateDirectConnectGatewayNatGateway(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create vpc dcGatewayAttachment failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(vpcId + FILED_SP + directConnectGatewayId + FILED_SP + natGatewayId)

	return resourceTencentCloudDcGatewayAttachmentRead(d, meta)
}

func resourceTencentCloudDcGatewayAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vpc_dc_gateway_attachment.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	vpcId := idSplit[0]
	directConnectGatewayId := idSplit[1]
	natGatewayId := idSplit[2]

	dcGatewayAttachment, err := service.DescribeDcGatewayAttachmentById(ctx, vpcId, directConnectGatewayId, natGatewayId)
	if err != nil {
		return err
	}

	if dcGatewayAttachment == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `VpcDcGatewayAttachment` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if dcGatewayAttachment.VpcId != nil {
		_ = d.Set("vpc_id", dcGatewayAttachment.VpcId)
	}

	if dcGatewayAttachment.NatGatewayId != nil {
		_ = d.Set("nat_gateway_id", dcGatewayAttachment.NatGatewayId)
	}

	if dcGatewayAttachment.DirectConnectGatewayId != nil {
		_ = d.Set("direct_connect_gateway_id", dcGatewayAttachment.DirectConnectGatewayId)
	}

	return nil
}

func resourceTencentCloudDcGatewayAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vpc_dc_gateway_attachment.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}
	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	vpcId := idSplit[0]
	directConnectGatewayId := idSplit[1]
	natGatewayId := idSplit[2]

	if err := service.DeleteDcGatewayAttachmentById(ctx, vpcId, directConnectGatewayId, natGatewayId); err != nil {
		return err
	}

	return nil
}
