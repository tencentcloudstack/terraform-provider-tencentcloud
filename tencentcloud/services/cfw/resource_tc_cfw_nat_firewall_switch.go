package cfw

import (
	"context"
	"fmt"
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cfw "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cfw/v20190904"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudCfwNatFirewallSwitch() *schema.Resource {
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
	defer tccommon.LogElapsed("resource.tencentcloud_cfw_nat_firewall_switch.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	natInsId := d.Get("nat_ins_id").(string)
	subnetId := d.Get("subnet_id").(string)

	d.SetId(strings.Join([]string{natInsId, subnetId}, tccommon.FILED_SP))

	return resourceTencentCloudCfwNatFirewallSwitchUpdate(d, meta)
}

func resourceTencentCloudCfwNatFirewallSwitchRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cfw_nat_firewall_switch.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = CfwService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
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
	defer tccommon.LogElapsed("resource.tencentcloud_cfw_nat_firewall_switch.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = CfwService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		request = cfw.NewModifyNatFwSwitchRequest()
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", idSplit)
	}
	natInsId := idSplit[0]
	subnetId := idSplit[1]

	if v, ok := d.GetOkExists("enable"); ok {
		request.Enable = helper.IntInt64(v.(int))
	}

	request.SubnetIdList = common.StringPtrs([]string{subnetId})

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCfwClient().ModifyNatFwSwitch(request)
		if e != nil {
			return tccommon.RetryError(e)
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
	err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		switchDetail, e := service.DescribeCfwNatFirewallSwitchById(ctx, natInsId, subnetId)
		if e != nil {
			return tccommon.RetryError(e)
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
	defer tccommon.LogElapsed("resource.tencentcloud_cfw_nat_firewall_switch.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
