package tencentcloud

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceTencentCloudDayuEip() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDayuEipCreate,
		Read:   resourceTencentCloudDayuEipRead,
		Delete: resourceTencentCloudDayuEipDelete,

		Schema: map[string]*schema.Schema{
			"resource_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the resource.",
			},
			"eip": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Eip of the resource.",
			},
			"bind_resource_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Resource id to bind.",
			},
			"bind_resource_region": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Resource region to bind.",
			},
			"bind_resource_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateAllowedStringValue(DDOS_EIP_BIND_RESOURCE_TYPE),
				Description:  "Resource type to bind, value range [`clb`, `cvm`].",
			},
			"resource_region": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Region of the resource instance.",
			},
			"eip_bound_rsc_ins": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Eip bound rsc ins of the resource instance.",
			},
			"eip_bound_rsc_eni": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Eip bound rsc eni of the resource instance.",
			},
			"eip_bound_rsc_vip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Eip bound rsc vip of the resource instance.",
			},
			"eip_address_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Eip address status of the resource instance.",
			},
			"protection_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Protection status of the resource instance.",
			},
			"created_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Created time of the resource instance.",
			},
			"expired_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Expired time of the resource instance.",
			},
			"modify_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Modify time of the resource instance.",
			},
		},
	}
}

func resourceTencentCloudDayuEipCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dayu_l4_rule.create")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	antiddosService := AntiddosService{client: meta.(*TencentCloudClient).apiV3Conn}
	bindResourceType := d.Get("bind_resource_type").(string)
	resourceId := d.Get("resource_id").(string)
	eip := d.Get("eip").(string)
	bindResourceId := d.Get("bind_resource_id").(string)
	bindResourceRegion := d.Get("bind_resource_region").(string)
	if bindResourceType == DDOS_EIP_BIND_RESOURCE_TYPE_CLB {
		err := antiddosService.AssociateDDoSEipLoadBalancer(ctx, resourceId, eip, bindResourceId, bindResourceRegion)
		if err != nil {
			return err
		}
	}

	if bindResourceType == DDOS_EIP_BIND_RESOURCE_TYPE_CVM {
		err := antiddosService.AssociateDDoSEipAddress(ctx, resourceId, eip, bindResourceId, bindResourceRegion)
		if err != nil {
			return err
		}
	}

	for {
		bgpIPInstances, err := antiddosService.DescribeListBGPIPInstances(ctx, resourceId, DDOS_EIP_BIND_STATUS, 0, 10)
		if err != nil {
			return err
		}
		if len(bgpIPInstances) != 0 {
			break
		}
	}
	d.SetId(resourceId + FILED_SP + eip)
	return resourceTencentCloudDayuEipRead(d, meta)
}

func resourceTencentCloudDayuEipRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dayu_l4_rule.read")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	items := strings.Split(d.Id(), FILED_SP)
	if len(items) < 2 {
		return fmt.Errorf("broken ID of dayu eip.")
	}
	resourceId := items[0]
	antiddosService := AntiddosService{client: meta.(*TencentCloudClient).apiV3Conn}
	bgpIPInstances, err := antiddosService.DescribeListBGPIPInstances(ctx, resourceId, DDOS_EIP_BIND_STATUS, 0, 10)
	if err != nil {
		return err
	}
	if len(bgpIPInstances) != 0 {
		posBGPIPInstance := bgpIPInstances[0]
		_ = d.Set("resource_region", *posBGPIPInstance.Region.Region)
		_ = d.Set("eip_bound_rsc_ins", *posBGPIPInstance.EipAddressInfo.EipBoundRscIns)
		_ = d.Set("eip_bound_rsc_eni", *posBGPIPInstance.EipAddressInfo.EipBoundRscEni)
		_ = d.Set("eip_bound_rsc_vip", *posBGPIPInstance.EipAddressInfo.EipBoundRscVip)
		_ = d.Set("eip_address_status", *posBGPIPInstance.EipAddressStatus)
		_ = d.Set("protection_status", *posBGPIPInstance.Status)
		_ = d.Set("created_time", *posBGPIPInstance.CreatedTime)
		_ = d.Set("expired_time", *posBGPIPInstance.ExpiredTime)
		_ = d.Set("modify_time", *posBGPIPInstance.EipAddressInfo.ModifyTime)
	}

	return nil
}

func resourceTencentCloudDayuEipDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dayu_l4_rule.delete")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	items := strings.Split(d.Id(), FILED_SP)
	if len(items) < 2 {
		return fmt.Errorf("broken ID of dayu eip.")
	}
	resourceId := items[0]
	eip := items[1]
	antiddosService := AntiddosService{client: meta.(*TencentCloudClient).apiV3Conn}
	err := antiddosService.DisassociateDDoSEipAddress(ctx, resourceId, eip)
	if err != nil {
		return err
	}
	return nil
}
