/*
Provides a resource to create a cfw edge_firewall_switch

Example Usage

```hcl
resource "tencentcloud_cfw_edge_firewall_switch" "edge_firewall_switch" {
  enable = 1
  edge_ip_switch_lst {
		public_ip = "1.1.1.1"
		subnet_id = "subnet-id"
		endpoint_ip = ""
		switch_mode = 0

  }
}
```

Import

cfw edge_firewall_switch can be imported using the id, e.g.

```
terraform import tencentcloud_cfw_edge_firewall_switch.edge_firewall_switch edge_firewall_switch_id
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

func resourceTencentCloudCfwEdgeFirewallSwitch() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCfwEdgeFirewallSwitchCreate,
		Read:   resourceTencentCloudCfwEdgeFirewallSwitchRead,
		Update: resourceTencentCloudCfwEdgeFirewallSwitchUpdate,
		Delete: resourceTencentCloudCfwEdgeFirewallSwitchDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"enable": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "0 off switch1 Turn on the switch.",
			},

			"edge_ip_switch_lst": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "Operation switch details.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"public_ip": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Public Ip.",
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The first EIP switch in the vpc is turned on, and you need to specify a subnet to create a private connection.",
						},
						"endpoint_ip": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Create a private connection to specify an IP.",
						},
						"switch_mode": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "0: bypass; 1: serial.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudCfwEdgeFirewallSwitchCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cfw_edge_firewall_switch.create")()
	defer inconsistentCheck(d, meta)()

	var publicIp string
	if v, ok := d.GetOk("public_ip"); ok {
		publicIp = v.(string)
	}

	d.SetId(publicIp)

	return resourceTencentCloudCfwEdgeFirewallSwitchUpdate(d, meta)
}

func resourceTencentCloudCfwEdgeFirewallSwitchRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cfw_edge_firewall_switch.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CfwService{client: meta.(*TencentCloudClient).apiV3Conn}

	edgeFirewallSwitchId := d.Id()

	edgeFirewallSwitch, err := service.DescribeCfwEdgeFirewallSwitchById(ctx, publicIp)
	if err != nil {
		return err
	}

	if edgeFirewallSwitch == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CfwEdgeFirewallSwitch` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if edgeFirewallSwitch.Enable != nil {
		_ = d.Set("enable", edgeFirewallSwitch.Enable)
	}

	if edgeFirewallSwitch.EdgeIpSwitchLst != nil {
		edgeIpSwitchLstList := []interface{}{}
		for _, edgeIpSwitchLst := range edgeFirewallSwitch.EdgeIpSwitchLst {
			edgeIpSwitchLstMap := map[string]interface{}{}

			if edgeFirewallSwitch.EdgeIpSwitchLst.PublicIp != nil {
				edgeIpSwitchLstMap["public_ip"] = edgeFirewallSwitch.EdgeIpSwitchLst.PublicIp
			}

			if edgeFirewallSwitch.EdgeIpSwitchLst.SubnetId != nil {
				edgeIpSwitchLstMap["subnet_id"] = edgeFirewallSwitch.EdgeIpSwitchLst.SubnetId
			}

			if edgeFirewallSwitch.EdgeIpSwitchLst.EndpointIp != nil {
				edgeIpSwitchLstMap["endpoint_ip"] = edgeFirewallSwitch.EdgeIpSwitchLst.EndpointIp
			}

			if edgeFirewallSwitch.EdgeIpSwitchLst.SwitchMode != nil {
				edgeIpSwitchLstMap["switch_mode"] = edgeFirewallSwitch.EdgeIpSwitchLst.SwitchMode
			}

			edgeIpSwitchLstList = append(edgeIpSwitchLstList, edgeIpSwitchLstMap)
		}

		_ = d.Set("edge_ip_switch_lst", edgeIpSwitchLstList)

	}

	return nil
}

func resourceTencentCloudCfwEdgeFirewallSwitchUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cfw_edge_firewall_switch.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := cfw.NewModifyEdgeIpSwitchRequest()

	edgeFirewallSwitchId := d.Id()

	request.PublicIp = &publicIp

	immutableArgs := []string{"enable", "edge_ip_switch_lst"}

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

	if d.HasChange("edge_ip_switch_lst") {
		if v, ok := d.GetOk("edge_ip_switch_lst"); ok {
			for _, item := range v.([]interface{}) {
				edgeIpSwitch := cfw.EdgeIpSwitch{}
				if v, ok := dMap["public_ip"]; ok {
					edgeIpSwitch.PublicIp = helper.String(v.(string))
				}
				if v, ok := dMap["subnet_id"]; ok {
					edgeIpSwitch.SubnetId = helper.String(v.(string))
				}
				if v, ok := dMap["endpoint_ip"]; ok {
					edgeIpSwitch.EndpointIp = helper.String(v.(string))
				}
				if v, ok := dMap["switch_mode"]; ok {
					edgeIpSwitch.SwitchMode = helper.IntInt64(v.(int))
				}
				request.EdgeIpSwitchLst = append(request.EdgeIpSwitchLst, &edgeIpSwitch)
			}
		}
	}

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

	return resourceTencentCloudCfwEdgeFirewallSwitchRead(d, meta)
}

func resourceTencentCloudCfwEdgeFirewallSwitchDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cfw_edge_firewall_switch.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
