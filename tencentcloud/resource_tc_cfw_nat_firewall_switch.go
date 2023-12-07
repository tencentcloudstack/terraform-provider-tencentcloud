package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cfw "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cfw/v20190904"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudCfwNatFirewallSwitch() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCfwNatFirewallSwitchCreate,
		Read:   resourceTencentCloudCfwNatFirewallSwitchRead,
		Update: resourceTencentCloudCfwNatFirewallSwitchUpdate,
		Delete: resourceTencentCloudCfwNatFirewallSwitchDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"nat_ins_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Firewall instance id.",
			},
			"subnet_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "subnet id.",
			},
			"enable": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Switch, 0: off, 1: on.",
			},
		},
	}
}

func resourceTencentCloudCfwNatFirewallSwitchCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cfw_nat_firewall_switch.create")()
	defer inconsistentCheck(d, meta)()

	natInsId := d.Get("nat_ins_id").(string)
	subnetId := d.Get("subnet_id").(string)

	d.SetId(strings.Join([]string{natInsId, subnetId}, FILED_SP))

	return resourceTencentCloudCfwNatFirewallSwitchUpdate(d, meta)
}

func resourceTencentCloudCfwNatFirewallSwitchRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cfw_nat_firewall_switch.read")()
	defer inconsistentCheck(d, meta)()

	var (
		logId   = getLogId(contextNil)
		ctx     = context.WithValue(context.TODO(), logIdKey, logId)
		service = CfwService{client: meta.(*TencentCloudClient).apiV3Conn}
	)

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", idSplit)
	}
	natInsId := idSplit[0]
	subnetId := idSplit[1]

	natFirewallSwitch, err := service.DescribeCfwNatFirewallSwitchById(ctx, natInsId, subnetId)
	if err != nil {
		return err
	}

	if natFirewallSwitch == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CfwNatFirewallSwitch` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if natFirewallSwitch.Enable != nil {
		_ = d.Set("enable", natFirewallSwitch.Enable)
	}

	return nil
}

func resourceTencentCloudCfwNatFirewallSwitchUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cfw_nat_firewall_switch.update")()
	defer inconsistentCheck(d, meta)()

	var (
		logId   = getLogId(contextNil)
		ctx     = context.WithValue(context.TODO(), logIdKey, logId)
		service = CfwService{client: meta.(*TencentCloudClient).apiV3Conn}
		request = cfw.NewModifyNatFwSwitchRequest()
	)

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", idSplit)
	}
	natInsId := idSplit[0]
	subnetId := idSplit[1]

	if v, ok := d.GetOkExists("enable"); ok {
		request.Enable = helper.IntInt64(v.(int))
	}

	request.SubnetIdList = common.StringPtrs([]string{subnetId})

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCfwClient().ModifyNatFwSwitch(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s update cfw natFirewallSwitch failed, reason:%+v", logId, err)
		return err
	}

	// wait
	err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		switchDetail, e := service.DescribeCfwNatFirewallSwitchById(ctx, natInsId, subnetId)
		if e != nil {
			return retryError(e)
		}

		if *switchDetail.Status == 0 {
			return nil
		}

		return resource.RetryableError(fmt.Errorf("update cfw natFirewallSwitch status is %d", *switchDetail.Status))
	})

	if err != nil {
		return err
	}

	return resourceTencentCloudCfwNatFirewallSwitchRead(d, meta)
}

func resourceTencentCloudCfwNatFirewallSwitchDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cfw_nat_firewall_switch.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
