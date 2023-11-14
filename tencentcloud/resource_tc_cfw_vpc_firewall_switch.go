/*
Provides a resource to create a cfw vpc_firewall_switch

Example Usage

```hcl
resource "tencentcloud_cfw_vpc_firewall_switch" "vpc_firewall_switch" {
  enable = 1
  all_switch = 0
  switch_list {
		switch_mode = 1
		switch_id = "cfws-id"

  }
}
```

Import

cfw vpc_firewall_switch can be imported using the id, e.g.

```
terraform import tencentcloud_cfw_vpc_firewall_switch.vpc_firewall_switch vpc_firewall_switch_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cfw "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cfw/v20190904"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudCfwVpcFirewallSwitch() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCfwVpcFirewallSwitchCreate,
		Read:   resourceTencentCloudCfwVpcFirewallSwitchRead,
		Update: resourceTencentCloudCfwVpcFirewallSwitchUpdate,
		Delete: resourceTencentCloudCfwVpcFirewallSwitchDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"enable": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Turn the switch on or off0: turn off the switch1: Turn on the switch.",
			},

			"all_switch": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Whether to operate all switches 0 does not operate all switches, 1 operates all switches.",
			},

			"switch_list": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "Switch list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"switch_mode": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Switching mode of firewall instance 1: Single point interworking 2: Multipoint interworking 3: Full interworking 4: Custom routing.",
						},
						"switch_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Firewall switch ID Support three types 1. Edge switch (single point interconnection) 2. Point switch (multi-point interoperability) 3. Full switch (full interoperability).",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudCfwVpcFirewallSwitchCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cfw_vpc_firewall_switch.create")()
	defer inconsistentCheck(d, meta)()

	var switchId string
	if v, ok := d.GetOk("switch_id"); ok {
		switchId = v.(string)
	}

	d.SetId(switchId)

	return resourceTencentCloudCfwVpcFirewallSwitchUpdate(d, meta)
}

func resourceTencentCloudCfwVpcFirewallSwitchRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cfw_vpc_firewall_switch.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CfwService{client: meta.(*TencentCloudClient).apiV3Conn}

	vpcFirewallSwitchId := d.Id()

	vpcFirewallSwitch, err := service.DescribeCfwVpcFirewallSwitchById(ctx, switchId)
	if err != nil {
		return err
	}

	if vpcFirewallSwitch == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CfwVpcFirewallSwitch` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if vpcFirewallSwitch.Enable != nil {
		_ = d.Set("enable", vpcFirewallSwitch.Enable)
	}

	if vpcFirewallSwitch.AllSwitch != nil {
		_ = d.Set("all_switch", vpcFirewallSwitch.AllSwitch)
	}

	if vpcFirewallSwitch.SwitchList != nil {
		switchListList := []interface{}{}
		for _, switchList := range vpcFirewallSwitch.SwitchList {
			switchListMap := map[string]interface{}{}

			if vpcFirewallSwitch.SwitchList.SwitchMode != nil {
				switchListMap["switch_mode"] = vpcFirewallSwitch.SwitchList.SwitchMode
			}

			if vpcFirewallSwitch.SwitchList.SwitchId != nil {
				switchListMap["switch_id"] = vpcFirewallSwitch.SwitchList.SwitchId
			}

			switchListList = append(switchListList, switchListMap)
		}

		_ = d.Set("switch_list", switchListList)

	}

	return nil
}

func resourceTencentCloudCfwVpcFirewallSwitchUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cfw_vpc_firewall_switch.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := cfw.NewModifyFwGroupSwitchRequest()

	vpcFirewallSwitchId := d.Id()

	request.SwitchId = &switchId

	immutableArgs := []string{"enable", "all_switch", "switch_list"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("enable") {
		if v, ok := d.GetOkExists("enable"); ok {
			request.Enable = helper.IntInt64(v.(int))
		}
	}

	if d.HasChange("all_switch") {
		if v, ok := d.GetOkExists("all_switch"); ok {
			request.AllSwitch = helper.IntInt64(v.(int))
		}
	}

	if d.HasChange("switch_list") {
		if v, ok := d.GetOk("switch_list"); ok {
			for _, item := range v.([]interface{}) {
				fwGroupSwitch := cfw.FwGroupSwitch{}
				if v, ok := dMap["switch_mode"]; ok {
					fwGroupSwitch.SwitchMode = helper.IntInt64(v.(int))
				}
				if v, ok := dMap["switch_id"]; ok {
					fwGroupSwitch.SwitchId = helper.String(v.(string))
				}
				request.SwitchList = append(request.SwitchList, &fwGroupSwitch)
			}
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCfwClient().ModifyFwGroupSwitch(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update cfw vpcFirewallSwitch failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudCfwVpcFirewallSwitchRead(d, meta)
}

func resourceTencentCloudCfwVpcFirewallSwitchDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cfw_vpc_firewall_switch.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
