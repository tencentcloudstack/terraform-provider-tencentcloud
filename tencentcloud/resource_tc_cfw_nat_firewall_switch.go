/*
Provides a resource to create a cfw nat_firewall_switch

Example Usage

```hcl
resource "tencentcloud_cfw_nat_firewall_switch" "nat_firewall_switch" {
  enable = 1
  cfw_ins_id_list =
  subnet_id_list =
  route_table_id_list =
}
```

Import

cfw nat_firewall_switch can be imported using the id, e.g.

```
terraform import tencentcloud_cfw_nat_firewall_switch.nat_firewall_switch nat_firewall_switch_id
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
			"enable": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Switch, 0: off, 1: on.",
			},

			"cfw_ins_id_list": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Firewall instance id list, of which only one type can be passed: CfwInsIdList, SubnetIdList and RouteTableIdList.",
			},

			"subnet_id_list": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "A list of subnet ids, of which only one type can be passed: CfwInsIdList, SubnetIdList and RouteTableIdList.",
			},

			"route_table_id_list": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "A list of routing table ids, of which only one type can be passed: CfwInsIdList, SubnetIdList and RouteTableIdList.",
			},
		},
	}
}

func resourceTencentCloudCfwNatFirewallSwitchCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cfw_nat_firewall_switch.create")()
	defer inconsistentCheck(d, meta)()

	var subnetId string
	if v, ok := d.GetOk("subnet_id"); ok {
		subnetId = v.(string)
	}

	d.SetId(subnetId)

	return resourceTencentCloudCfwNatFirewallSwitchUpdate(d, meta)
}

func resourceTencentCloudCfwNatFirewallSwitchRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cfw_nat_firewall_switch.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CfwService{client: meta.(*TencentCloudClient).apiV3Conn}

	natFirewallSwitchId := d.Id()

	natFirewallSwitch, err := service.DescribeCfwNatFirewallSwitchById(ctx, subnetId)
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

	if natFirewallSwitch.CfwInsIdList != nil {
		_ = d.Set("cfw_ins_id_list", natFirewallSwitch.CfwInsIdList)
	}

	if natFirewallSwitch.SubnetIdList != nil {
		_ = d.Set("subnet_id_list", natFirewallSwitch.SubnetIdList)
	}

	if natFirewallSwitch.RouteTableIdList != nil {
		_ = d.Set("route_table_id_list", natFirewallSwitch.RouteTableIdList)
	}

	return nil
}

func resourceTencentCloudCfwNatFirewallSwitchUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cfw_nat_firewall_switch.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := cfw.NewModifyNatFwSwitchRequest()

	natFirewallSwitchId := d.Id()

	request.SubnetId = &subnetId

	immutableArgs := []string{"enable", "cfw_ins_id_list", "subnet_id_list", "route_table_id_list"}

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

	if d.HasChange("cfw_ins_id_list") {
		if v, ok := d.GetOk("cfw_ins_id_list"); ok {
			cfwInsIdListSet := v.(*schema.Set).List()
			for i := range cfwInsIdListSet {
				cfwInsIdList := cfwInsIdListSet[i].(string)
				request.CfwInsIdList = append(request.CfwInsIdList, &cfwInsIdList)
			}
		}
	}

	if d.HasChange("subnet_id_list") {
		if v, ok := d.GetOk("subnet_id_list"); ok {
			subnetIdListSet := v.(*schema.Set).List()
			for i := range subnetIdListSet {
				subnetIdList := subnetIdListSet[i].(string)
				request.SubnetIdList = append(request.SubnetIdList, &subnetIdList)
			}
		}
	}

	if d.HasChange("route_table_id_list") {
		if v, ok := d.GetOk("route_table_id_list"); ok {
			routeTableIdListSet := v.(*schema.Set).List()
			for i := range routeTableIdListSet {
				routeTableIdList := routeTableIdListSet[i].(string)
				request.RouteTableIdList = append(request.RouteTableIdList, &routeTableIdList)
			}
		}
	}

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

	return resourceTencentCloudCfwNatFirewallSwitchRead(d, meta)
}

func resourceTencentCloudCfwNatFirewallSwitchDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cfw_nat_firewall_switch.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
