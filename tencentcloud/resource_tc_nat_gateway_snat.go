package tencentcloud

import (
	"context"
	"errors"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudNatGatewaySnat() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudNatGatewaySnatCreate,
		Read:   resourceTencentCloudNatGatewaySnatRead,
		Update: resourceTencentCloudNatGatewaySnatUpdate,
		Delete: resourceTencentCloudNatGatewaySnatDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: NatGatewaySnatPara(),
	}
}

func NatGatewaySnatPara() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"nat_gateway_id": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "NAT gateway ID.",
		},
		"resource_type": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "Resource type. Valid values: " + NAT_GATEWAY_TYPE_SUBNET + ", " + NAT_GATEWAY_TYPE_NETWORK_INTERFACE + ".",
		},
		"subnet_id": {
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
			Description: "Subnet instance ID, required when `resource_type` is " + NAT_GATEWAY_TYPE_SUBNET + ".",
		},
		"subnet_cidr_block": {
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
			Description: "The IPv4 CIDR of the subnet, required when `resource_type` is " + NAT_GATEWAY_TYPE_SUBNET + ".",
		},
		"instance_id": {
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
			Description: "Instance ID, required when `resource_type` is " + NAT_GATEWAY_TYPE_NETWORK_INTERFACE + ".",
		},
		"instance_private_ip_addr": {
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
			Description: "Private IPs of the instance's primary ENI, required when `resource_type` is " + NAT_GATEWAY_TYPE_NETWORK_INTERFACE + ".",
		},
		"public_ip_addr": {
			Type:        schema.TypeList,
			Required:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
			Description: "Elastic IP address pool.",
		},
		"description": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Description.",
		},
		"snat_id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "SNAT rule ID.",
		},
		"create_time": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Create time.",
		},
	}
}

func resourceTencentCloudNatGatewaySnatCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_nat_gateway_snat.create")()

	var (
		logId        = getLogId(contextNil)
		ctx          = context.WithValue(context.TODO(), logIdKey, logId)
		vpcService   = VpcService{client: meta.(*TencentCloudClient).apiV3Conn}
		natGatewayId = d.Get("nat_gateway_id").(string)
	)

	// param valid
	err := paramValid(d)
	if err != nil {
		return err
	}

	snat := getSnat(d)

	err = vpcService.CreateNatGatewaySnat(ctx, natGatewayId, snat)
	if err != nil {
		log.Printf("[CRITAL]%s create nat gateway snat failed, reason:%s\n", logId, err.Error())
		return err
	}

	err, result := vpcService.DescribeNatGatewaySnats(ctx, natGatewayId, []*vpc.Filter{{
		Name:   helper.String("resource-id"),
		Values: []*string{snat.ResourceId},
	}, {
		Name:   helper.String("public-ip-address"),
		Values: snat.PublicIpAddresses,
	}, {
		Name:   helper.String("description"),
		Values: []*string{snat.Description},
	}})
	if err != nil {
		return err
	}
	if len(result) == 0 {
		return errors.New("[CRITAL] create nat gateway snat failed: read result is empty")
	}
	rule := result[len(result)-1]
	d.SetId(helper.IdFormat(*rule.NatGatewayId, *rule.ResourceId))

	return resourceTencentCloudNatGatewaySnatRead(d, meta)
}

func resourceTencentCloudNatGatewaySnatRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_nat_gateway_snat.read")()
	defer inconsistentCheck(d, meta)()

	var (
		logId   = getLogId(contextNil)
		ctx     = context.WithValue(context.TODO(), logIdKey, logId)
		service = VpcService{client: meta.(*TencentCloudClient).apiV3Conn}
		id      = d.Id()
	)
	compositeId := helper.IdParse(id)
	if len(compositeId) != 2 {
		return errors.New("the id format must be '{nat_gateway_id}#{resource_id}'")
	}

	err, snatList := service.DescribeNatGatewaySnats(ctx, compositeId[0], nil)
	if err != nil {
		log.Printf("[CRITAL]%s read nat gateway snat failed, reason:%s\n", logId, err.Error())
		return err
	}
	var snat *vpc.SourceIpTranslationNatRule
	for _, s := range snatList {
		if compositeId[1] == *s.ResourceId {
			snat = s
		}
	}
	if snat == nil {
		d.SetId("")
		return nil
	}
	_ = d.Set("nat_gateway_id", snat.NatGatewayId)
	_ = d.Set("resource_type", snat.ResourceType)
	_ = d.Set("public_ip_addr", sortSnatPublicIpAddr(d, snat.PublicIpAddresses))
	_ = d.Set("description", snat.Description)
	_ = d.Set("snat_id", snat.NatGatewaySnatId)
	_ = d.Set("create_time", snat.CreatedTime)
	resourceType := snat.ResourceType
	if *resourceType == NAT_GATEWAY_TYPE_SUBNET {
		_ = d.Set("subnet_id", snat.ResourceId)
		_ = d.Set("subnet_cidr_block", snat.PrivateIpAddress)
	} else if *resourceType == NAT_GATEWAY_TYPE_NETWORK_INTERFACE {
		_ = d.Set("instance_id", snat.ResourceId)
		_ = d.Set("instance_private_ip_addr", snat.PrivateIpAddress)
	}
	return nil
}

func resourceTencentCloudNatGatewaySnatUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vpn_gateway_route.update")()

	var (
		logId   = getLogId(contextNil)
		ctx     = context.WithValue(context.TODO(), logIdKey, logId)
		service = VpcService{client: meta.(*TencentCloudClient).apiV3Conn}
		id      = d.Id()
	)
	compositeId := helper.IdParse(id)
	if len(compositeId) != 2 {
		return errors.New("the id format must be '{nat_gateway_id}#{snat_id}'")
	}

	// param valid
	err := paramValid(d)
	if err != nil {
		return err
	}

	// update
	snat := getSnat(d)
	err = service.ModifyNatGatewaySnat(ctx, compositeId[0], snat)
	if err != nil {
		log.Printf("[CRITAL]%s modify nat gateway snat failed, reason:%s\n", logId, err.Error())
		return err
	}

	return resourceTencentCloudNatGatewaySnatRead(d, meta)
}

func resourceTencentCloudNatGatewaySnatDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_nat_gateway_snat.delete")()

	var (
		logId   = getLogId(contextNil)
		ctx     = context.WithValue(context.TODO(), logIdKey, logId)
		service = VpcService{client: meta.(*TencentCloudClient).apiV3Conn}
	)

	natGatewayId := d.Get("nat_gateway_id").(string)
	snatId := d.Get("snat_id").(string)
	err := service.DeleteNatGatewaySnat(ctx, natGatewayId, snatId)
	if err != nil {
		log.Printf("[CRITAL]%s delete nat gateway snat failed, reason:%s\n", logId, err.Error())
		return err
	}
	return nil
}

func paramValid(d *schema.ResourceData) error {
	logId := getLogId(contextNil)
	resourceType := d.Get("resource_type")
	_, hasSubnetId := d.GetOk("subnet_id")
	_, hasCidrBlock := d.GetOk("subnet_cidr_block")
	_, hasInstanceId := d.GetOk("instance_id")
	_, hasPrivateIpAddr := d.GetOk("instance_private_ip_addr")
	if resourceType == NAT_GATEWAY_TYPE_SUBNET && !(hasSubnetId && hasCidrBlock) {
		log.Printf("[CRITAL]%s `resource_type` is %s, but hasSubnetId = %v, hasCidrBlock = %v",
			logId, NAT_GATEWAY_TYPE_SUBNET, hasSubnetId, hasCidrBlock)
		return errors.New("`subnet_id` and `subnet_cidr_block` required when `resource_type` is " + NAT_GATEWAY_TYPE_SUBNET)
	} else if resourceType == NAT_GATEWAY_TYPE_NETWORK_INTERFACE && !(hasInstanceId && hasPrivateIpAddr) {
		log.Printf("[CRITAL]%s `resource_type` is %s, but hasInstanceId = %v, hasPrivateIpAddr = %v",
			logId, NAT_GATEWAY_TYPE_NETWORK_INTERFACE, hasInstanceId, hasPrivateIpAddr)
		return errors.New("`instance_id` and `instance_private_ip_addr` required when `resource_type` is " + NAT_GATEWAY_TYPE_NETWORK_INTERFACE)
	}
	return nil
}

func getSnat(d *schema.ResourceData) *vpc.SourceIpTranslationNatRule {
	var (
		resourceType  = d.Get("resource_type").(string)
		publicIpAddrs = helper.InterfacesStringsPoint(d.Get("public_ip_addr").([]interface{}))
		description   = helper.String(d.Get("description").(string))
	)
	var resourceId, privateIpAddr string
	if resourceType == NAT_GATEWAY_TYPE_SUBNET {
		resourceId = d.Get("subnet_id").(string)
		privateIpAddr = d.Get("subnet_cidr_block").(string)
	} else if resourceType == NAT_GATEWAY_TYPE_NETWORK_INTERFACE {
		resourceId = d.Get("instance_id").(string)
		privateIpAddr = d.Get("instance_private_ip_addr").(string)
	}
	sNat := &vpc.SourceIpTranslationNatRule{
		ResourceId:        &resourceId,
		ResourceType:      &resourceType,
		PrivateIpAddress:  &privateIpAddr,
		PublicIpAddresses: publicIpAddrs,
		Description:       description,
	}
	if v, ok := d.GetOk("snat_id"); ok {
		sNat.NatGatewaySnatId = helper.String(v.(string))
	}
	return sNat
}

func sortSnatPublicIpAddr(d *schema.ResourceData, publicIpAddresses []*string) []*string {
	if v, ok := d.GetOk("public_ip_addr"); ok {
		result := make([]*string, 0)
		paramObjs := v.([]interface{})
		for _, paramObj := range paramObjs {
			for _, obj := range publicIpAddresses {
				if paramObj.(string) == *obj {
					result = append(result, obj)
					break
				}
			}
		}
		for _, obj := range publicIpAddresses {
			contain := false
			for _, exist := range result {
				if exist == obj {
					contain = true
					break
				}
			}
			if !contain {
				result = append(result, obj)
			}
		}
		return result
	}
	return publicIpAddresses
}
