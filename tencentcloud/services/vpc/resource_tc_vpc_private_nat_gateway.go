package vpc

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudVpcPrivateNatGateway() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudVpcPrivateNatGatewayCreate,
		Read:   resourceTencentCloudVpcPrivateNatGatewayRead,
		Update: resourceTencentCloudVpcPrivateNatGatewayUpdate,
		Delete: resourceTencentCloudVpcPrivateNatGatewayDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"nat_gateway_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Private network gateway name.",
			},

			"vpc_id": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Private Cloud instance ID. This parameter is required when creating a VPC type private network NAT gateway or a private network NAT gateway of private network gateway.",
			},

			"cross_domain": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeBool,
				Description: "Cross-domain parameters. Cross-domain binding of VPCs is supported only when the value is True.",
			},

			"vpc_type": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeBool,
				Description: "VPC type private network NAT gateway. Only when the value is True will a VPC type private network NAT gateway be created.",
			},

			"ccn_id": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Cloud Connect Network type The Cloud Connect Network instance ID required to be bound to the private network NAT gateway.",
			},
		},
	}
}

func resourceTencentCloudVpcPrivateNatGatewayCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vpc_private_nat_gateway.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	var (
		request    = vpc.NewCreatePrivateNatGatewayRequest()
		response   = vpc.NewCreatePrivateNatGatewayResponse()
		instanceId string
	)
	if v, ok := d.GetOk("nat_gateway_name"); ok {
		request.NatGatewayName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("vpc_id"); ok {
		request.VpcId = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("cross_domain"); ok {
		request.CrossDomain = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOkExists("vpc_type"); ok {
		request.VpcType = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOk("ccn_id"); ok {
		request.CcnId = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().CreatePrivateNatGateway(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create vpc privateNatGateway failed, reason:%+v", logId, err)
		return err
	}

	if response.Response != nil && len(response.Response.PrivateNatGatewaySet) > 0 {
		privateNatGateway := response.Response.PrivateNatGatewaySet[0]
		if privateNatGateway.NatGatewayId != nil {
			instanceId = *response.Response.PrivateNatGatewaySet[0].NatGatewayId
		}
	}
	if instanceId == "" {
		d.SetId("")
		return fmt.Errorf("instanceId is nil")
	}

	service := VpcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	err = resource.Retry(5*tccommon.ReadRetryTimeout, func() *resource.RetryError {
		privateNatGateway, errRet := service.DescribeVpcPrivateNatGatewayById(ctx, instanceId)
		if errRet != nil {
			return tccommon.RetryError(errRet)
		}
		if privateNatGateway.Status == nil {
			return resource.RetryableError(fmt.Errorf("waiting for instance create"))
		}
		if *privateNatGateway.Status != "AVAILABLE" {
			return resource.RetryableError(fmt.Errorf("waiting for instance create"))
		}
		return nil
	})
	if err != nil {
		return err
	}
	d.SetId(instanceId)

	return resourceTencentCloudVpcPrivateNatGatewayRead(d, meta)
}

func resourceTencentCloudVpcPrivateNatGatewayRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vpc_private_nat_gateway.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := VpcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	instanceId := d.Id()

	privateNatGateway, err := service.DescribeVpcPrivateNatGatewayById(ctx, instanceId)
	if err != nil {
		return err
	}

	if privateNatGateway == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `VpcPrivateNatGateway` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if privateNatGateway.NatGatewayName != nil {
		_ = d.Set("nat_gateway_name", privateNatGateway.NatGatewayName)
	}

	if privateNatGateway.VpcId != nil {
		_ = d.Set("vpc_id", privateNatGateway.VpcId)
	}

	if privateNatGateway.CrossDomain != nil {
		_ = d.Set("cross_domain", privateNatGateway.CrossDomain)
	}

	if privateNatGateway.VpcType != nil {
		_ = d.Set("vpc_type", privateNatGateway.VpcType)
	}

	if privateNatGateway.CcnId != nil {
		_ = d.Set("ccn_id", privateNatGateway.CcnId)
	}

	return nil
}

func resourceTencentCloudVpcPrivateNatGatewayUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vpc_private_nat_gateway.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	request := vpc.NewModifyPrivateNatGatewayAttributeRequest()

	instanceId := d.Id()

	request.NatGatewayId = &instanceId

	immutableArgs := []string{"vpc_id", "cross_domain", "vpc_type", "ccn_id"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("nat_gateway_name") {
		if v, ok := d.GetOk("nat_gateway_name"); ok {
			request.NatGatewayName = helper.String(v.(string))
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().ModifyPrivateNatGatewayAttribute(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update vpc privateNatGateway failed, reason:%+v", logId, err)
		return err
	}

	service := VpcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	err = resource.Retry(5*tccommon.ReadRetryTimeout, func() *resource.RetryError {
		privateNatGateway, errRet := service.DescribeVpcPrivateNatGatewayById(ctx, instanceId)
		if errRet != nil {
			return tccommon.RetryError(errRet)
		}
		if privateNatGateway.Status == nil {
			return resource.RetryableError(fmt.Errorf("waiting for instance update"))
		}
		if *privateNatGateway.Status != "AVAILABLE" {
			return resource.RetryableError(fmt.Errorf("waiting for instance update"))
		}
		return nil
	})
	if err != nil {
		return err
	}
	return resourceTencentCloudVpcPrivateNatGatewayRead(d, meta)
}

func resourceTencentCloudVpcPrivateNatGatewayDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vpc_private_nat_gateway.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := VpcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	instanceId := d.Id()

	if err := service.DeleteVpcPrivateNatGatewayById(ctx, instanceId); err != nil {
		return err
	}

	err := resource.Retry(5*tccommon.ReadRetryTimeout, func() *resource.RetryError {
		privateNatGateway, errRet := service.DescribeVpcPrivateNatGatewayById(ctx, instanceId)
		if errRet != nil {
			return tccommon.RetryError(errRet)
		}
		if privateNatGateway != nil {
			return resource.RetryableError(fmt.Errorf("waiting for instance delete"))
		}
		return nil
	})

	if err != nil {
		return err
	}
	return nil
}
