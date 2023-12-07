package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cfw "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cfw/v20190904"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudCfwEdgeFirewallSwitch() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCfwEdgeFirewallSwitchCreate,
		Read:   resourceTencentCloudCfwEdgeFirewallSwitchRead,
		Update: resourceTencentCloudCfwEdgeFirewallSwitchUpdate,
		Delete: resourceTencentCloudCfwEdgeFirewallSwitchDelete,

		Schema: map[string]*schema.Schema{
			"public_ip": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Public Ip.",
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The first EIP switch in the vpc is turned on, and you need to specify a subnet to create a private connection. If `switch_mode` is 1 and `enable` is 1, this field is required.",
			},
			"switch_mode": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "0: bypass; 1: serial.",
			},
			"enable": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Switch, 0: off, 1: on.",
			},
		},
	}
}

func resourceTencentCloudCfwEdgeFirewallSwitchCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cfw_edge_firewall_switch.create")()
	defer inconsistentCheck(d, meta)()

	publicIp := d.Get("public_ip").(string)
	d.SetId(publicIp)

	return resourceTencentCloudCfwEdgeFirewallSwitchUpdate(d, meta)
}

func resourceTencentCloudCfwEdgeFirewallSwitchRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cfw_edge_firewall_switch.read")()
	defer inconsistentCheck(d, meta)()

	var (
		logId    = getLogId(contextNil)
		ctx      = context.WithValue(context.TODO(), logIdKey, logId)
		service  = CfwService{client: meta.(*TencentCloudClient).apiV3Conn}
		publicIp = d.Id()
	)

	edgeFirewallSwitch, err := service.DescribeCfwEdgeFirewallSwitchById(ctx, publicIp)
	if err != nil {
		return err
	}

	if edgeFirewallSwitch == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CfwEdgeFirewallSwitch` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if edgeFirewallSwitch.PublicIp != nil {
		_ = d.Set("public_ip", edgeFirewallSwitch.PublicIp)
	}

	if edgeFirewallSwitch.SwitchMode != nil {
		_ = d.Set("switch_mode", edgeFirewallSwitch.SwitchMode)
	}

	if edgeFirewallSwitch.Status != nil {
		_ = d.Set("enable", edgeFirewallSwitch.Status)
	}

	return nil
}

func resourceTencentCloudCfwEdgeFirewallSwitchUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cfw_edge_firewall_switch.update")()
	defer inconsistentCheck(d, meta)()

	var (
		logId        = getLogId(contextNil)
		ctx          = context.WithValue(context.TODO(), logIdKey, logId)
		service      = CfwService{client: meta.(*TencentCloudClient).apiV3Conn}
		request      = cfw.NewModifyEdgeIpSwitchRequest()
		edgeIpSwitch = cfw.EdgeIpSwitch{}
		publicIp     = d.Id()
	)

	edgeIpSwitch.PublicIp = &publicIp

	if v, ok := d.GetOk("subnet_id"); ok {
		edgeIpSwitch.SubnetId = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("switch_mode"); ok {
		edgeIpSwitch.SwitchMode = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("enable"); ok {
		request.Enable = helper.IntInt64(v.(int))
	}

	request.EdgeIpSwitchLst = append(request.EdgeIpSwitchLst, &edgeIpSwitch)

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCfwClient().ModifyEdgeIpSwitch(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s update cfw edgeFirewallSwitch failed, reason:%+v", logId, err)
		return err
	}

	// wait
	err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		switchDetail, e := service.DescribeCfwEdgeFirewallSwitchById(ctx, publicIp)
		if e != nil {
			return retryError(e)
		}

		if *switchDetail.Status == 0 || *switchDetail.Status == 1 {
			return nil
		}

		return resource.RetryableError(fmt.Errorf("update cfw edgeFirewallSwitch status is %d", *switchDetail.Status))
	})

	if err != nil {
		return err
	}

	return resourceTencentCloudCfwEdgeFirewallSwitchRead(d, meta)
}

func resourceTencentCloudCfwEdgeFirewallSwitchDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cfw_edge_firewall_switch.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
