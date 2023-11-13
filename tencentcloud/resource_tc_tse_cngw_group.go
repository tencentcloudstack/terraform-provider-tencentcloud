/*
Provides a resource to create a tse cngw_group

Example Usage

```hcl
resource "tencentcloud_tse_cngw_group" "cngw_group" {
  gateway_id = ""
  name = ""
  node_config {
		specification = ""
		number =

  }
  subnet_id = ""
  description = ""
  internet_max_bandwidth_out =
  internet_config {
		internet_address_version = ""
		internet_pay_mode = ""
		internet_max_bandwidth_out =
		description = ""
		sla_type = ""
		multi_zone_flag =
		master_zone_id = ""
		slave_zone_id = ""

  }
}
```

Import

tse cngw_group can be imported using the id, e.g.

```
terraform import tencentcloud_tse_cngw_group.cngw_group cngw_group_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tse "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tse/v20201207"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
	"strings"
	"time"
)

func resourceTencentCloudTseCngwGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTseCngwGroupCreate,
		Read:   resourceTencentCloudTseCngwGroupRead,
		Update: resourceTencentCloudTseCngwGroupUpdate,
		Delete: resourceTencentCloudTseCngwGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"gateway_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Gateway IDonly postpaid gateway supported.",
			},

			"name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Gateway group name.",
			},

			"node_config": {
				Required:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Group node configration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"specification": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Group specification, 1c2g|2c4g|4c8g|8c16g.",
						},
						"number": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Group node number, 2-50.",
						},
					},
				},
			},

			"subnet_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Subnet ID. Assign an IP address to the engine in the VPC subnet. Reference value:- subnet-ahde9me9.",
			},

			"description": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Description information of group.",
			},

			"internet_max_bandwidth_out": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Public network outbound traffic bandwidth,[1,2048]Mbps.",
			},

			"internet_config": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Internet configration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"internet_address_version": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Internet type. Reference value:- IPV4 (default value)- IPV6.",
						},
						"internet_pay_mode": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Trade type of internet. Reference value:- BANDWIDTH- TRAFFIC (default value).",
						},
						"internet_max_bandwidth_out": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Public network bandwidth.",
						},
						"description": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Description of clb.",
						},
						"sla_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Specification type of clb. Default shared type when this parameter is empty. Reference value:- SLA LCU-supported.",
						},
						"multi_zone_flag": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether load balancing has multiple availability zones.",
						},
						"master_zone_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Primary availability zone.",
						},
						"slave_zone_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Alternate availability zone.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudTseCngwGroupCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tse_cngw_group.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request   = tse.NewCreateNativeGatewayServerGroupRequest()
		response  = tse.NewCreateNativeGatewayServerGroupResponse()
		gatewayId string
		groupId   string
	)
	if v, ok := d.GetOk("gateway_id"); ok {
		gatewayId = v.(string)
		request.GatewayId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "node_config"); ok {
		cloudNativeAPIGatewayNodeConfig := tse.CloudNativeAPIGatewayNodeConfig{}
		if v, ok := dMap["specification"]; ok {
			cloudNativeAPIGatewayNodeConfig.Specification = helper.String(v.(string))
		}
		if v, ok := dMap["number"]; ok {
			cloudNativeAPIGatewayNodeConfig.Number = helper.IntInt64(v.(int))
		}
		request.NodeConfig = &cloudNativeAPIGatewayNodeConfig
	}

	if v, ok := d.GetOk("subnet_id"); ok {
		request.SubnetId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("description"); ok {
		request.Description = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("internet_max_bandwidth_out"); ok {
		request.InternetMaxBandwidthOut = helper.IntUint64(v.(int))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "internet_config"); ok {
		internetConfig := tse.InternetConfig{}
		if v, ok := dMap["internet_address_version"]; ok {
			internetConfig.InternetAddressVersion = helper.String(v.(string))
		}
		if v, ok := dMap["internet_pay_mode"]; ok {
			internetConfig.InternetPayMode = helper.String(v.(string))
		}
		if v, ok := dMap["internet_max_bandwidth_out"]; ok {
			internetConfig.InternetMaxBandwidthOut = helper.IntUint64(v.(int))
		}
		if v, ok := dMap["description"]; ok {
			internetConfig.Description = helper.String(v.(string))
		}
		if v, ok := dMap["sla_type"]; ok {
			internetConfig.SlaType = helper.String(v.(string))
		}
		if v, ok := dMap["multi_zone_flag"]; ok {
			internetConfig.MultiZoneFlag = helper.Bool(v.(bool))
		}
		if v, ok := dMap["master_zone_id"]; ok {
			internetConfig.MasterZoneId = helper.String(v.(string))
		}
		if v, ok := dMap["slave_zone_id"]; ok {
			internetConfig.SlaveZoneId = helper.String(v.(string))
		}
		request.InternetConfig = &internetConfig
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTseClient().CreateNativeGatewayServerGroup(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create tse cngwGroup failed, reason:%+v", logId, err)
		return err
	}

	gatewayId = *response.Response.GatewayId
	d.SetId(strings.Join([]string{gatewayId, groupId}, FILED_SP))

	service := TseService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{"Running"}, 5*readRetryTimeout, time.Second, service.TseCngwGroupStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudTseCngwGroupRead(d, meta)
}

func resourceTencentCloudTseCngwGroupRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tse_cngw_group.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TseService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	gatewayId := idSplit[0]
	groupId := idSplit[1]

	cngwGroup, err := service.DescribeTseCngwGroupById(ctx, gatewayId, groupId)
	if err != nil {
		return err
	}

	if cngwGroup == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `TseCngwGroup` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if cngwGroup.GatewayId != nil {
		_ = d.Set("gateway_id", cngwGroup.GatewayId)
	}

	if cngwGroup.Name != nil {
		_ = d.Set("name", cngwGroup.Name)
	}

	if cngwGroup.NodeConfig != nil {
		nodeConfigMap := map[string]interface{}{}

		if cngwGroup.NodeConfig.Specification != nil {
			nodeConfigMap["specification"] = cngwGroup.NodeConfig.Specification
		}

		if cngwGroup.NodeConfig.Number != nil {
			nodeConfigMap["number"] = cngwGroup.NodeConfig.Number
		}

		_ = d.Set("node_config", []interface{}{nodeConfigMap})
	}

	if cngwGroup.SubnetId != nil {
		_ = d.Set("subnet_id", cngwGroup.SubnetId)
	}

	if cngwGroup.Description != nil {
		_ = d.Set("description", cngwGroup.Description)
	}

	if cngwGroup.InternetMaxBandwidthOut != nil {
		_ = d.Set("internet_max_bandwidth_out", cngwGroup.InternetMaxBandwidthOut)
	}

	if cngwGroup.InternetConfig != nil {
		internetConfigMap := map[string]interface{}{}

		if cngwGroup.InternetConfig.InternetAddressVersion != nil {
			internetConfigMap["internet_address_version"] = cngwGroup.InternetConfig.InternetAddressVersion
		}

		if cngwGroup.InternetConfig.InternetPayMode != nil {
			internetConfigMap["internet_pay_mode"] = cngwGroup.InternetConfig.InternetPayMode
		}

		if cngwGroup.InternetConfig.InternetMaxBandwidthOut != nil {
			internetConfigMap["internet_max_bandwidth_out"] = cngwGroup.InternetConfig.InternetMaxBandwidthOut
		}

		if cngwGroup.InternetConfig.Description != nil {
			internetConfigMap["description"] = cngwGroup.InternetConfig.Description
		}

		if cngwGroup.InternetConfig.SlaType != nil {
			internetConfigMap["sla_type"] = cngwGroup.InternetConfig.SlaType
		}

		if cngwGroup.InternetConfig.MultiZoneFlag != nil {
			internetConfigMap["multi_zone_flag"] = cngwGroup.InternetConfig.MultiZoneFlag
		}

		if cngwGroup.InternetConfig.MasterZoneId != nil {
			internetConfigMap["master_zone_id"] = cngwGroup.InternetConfig.MasterZoneId
		}

		if cngwGroup.InternetConfig.SlaveZoneId != nil {
			internetConfigMap["slave_zone_id"] = cngwGroup.InternetConfig.SlaveZoneId
		}

		_ = d.Set("internet_config", []interface{}{internetConfigMap})
	}

	return nil
}

func resourceTencentCloudTseCngwGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tse_cngw_group.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := tse.NewModifyNativeGatewayServerGroupRequest()

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	gatewayId := idSplit[0]
	groupId := idSplit[1]

	request.GatewayId = &gatewayId
	request.GroupId = &groupId

	immutableArgs := []string{"gateway_id", "name", "node_config", "subnet_id", "description", "internet_max_bandwidth_out", "internet_config"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("gateway_id") {
		if v, ok := d.GetOk("gateway_id"); ok {
			request.GatewayId = helper.String(v.(string))
		}
	}

	if d.HasChange("name") {
		if v, ok := d.GetOk("name"); ok {
			request.Name = helper.String(v.(string))
		}
	}

	if d.HasChange("description") {
		if v, ok := d.GetOk("description"); ok {
			request.Description = helper.String(v.(string))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTseClient().ModifyNativeGatewayServerGroup(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update tse cngwGroup failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudTseCngwGroupRead(d, meta)
}

func resourceTencentCloudTseCngwGroupDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tse_cngw_group.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TseService{client: meta.(*TencentCloudClient).apiV3Conn}
	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	gatewayId := idSplit[0]
	groupId := idSplit[1]

	if err := service.DeleteTseCngwGroupById(ctx, gatewayId, groupId); err != nil {
		return err
	}

	service := TseService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{"Deleted"}, 5*readRetryTimeout, time.Second, service.TseCngwGroupStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return nil
}
