/*
Provides a resource to create a cfw natfw

Example Usage

```hcl
resource "tencentcloud_cfw_natfw" "natfw" {
  name = "test natfw"
  width = 20
  mode = 0
  new_mode_items {
		vpc_list =
		eips =
		add_count = 1

  }
  nat_gw_list =
  zone = "ap-guangzhou-1"
  zone_bak = "ap-guangzhou-2"
  cross_a_zone = 1
  is_create_domain = 0
  domain = ""
  fw_cidr_info {
		fw_cidr_type = "VpcSelf"
		fw_cidr_lst {
			vpc_id = "vpc-id"
			fw_cidr = "10.96.0.0/16"
		}
		com_fw_cidr = ""

  }
}
```

Import

cfw natfw can be imported using the id, e.g.

```
terraform import tencentcloud_cfw_natfw.natfw natfw_id
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
	"time"
)

func resourceTencentCloudCfwNatfw() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCfwNatfwCreate,
		Read:   resourceTencentCloudCfwNatfwRead,
		Update: resourceTencentCloudCfwNatfwUpdate,
		Delete: resourceTencentCloudCfwNatfwDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Firewall instance name.",
			},

			"width": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Bandwidth.",
			},

			"mode": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Mode 1: access mode; 0: new mode.",
			},

			"new_mode_items": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "New mode passing parameters are added, at least one of NewModeItems and NatgwList is passed.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vpc_list": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Required:    true,
							Description: "List of vpcs connected in new mode.",
						},
						"eips": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional:    true,
							Description: "List of egress elastic public network IPs bound in the new mode, at least one of which is passed in Eips and AddCount.",
						},
						"add_count": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "In the new mode, the number of newly bound egress elastic public network IPs is specified, at least one of which is passed in Eips and AddCount.",
						},
					},
				},
			},

			"nat_gw_list": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "A list of nat gateways connected to the access mode, at least one of NewModeItems and NatgwList is passed.",
			},

			"zone": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Main Availability Zone, if empty, select the default Availability Zone.",
			},

			"zone_bak": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Standby Availability Zone, if empty, select the default Availability Zone.",
			},

			"cross_a_zone": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Off-site disaster recovery 1: use off-site disaster recovery; 0: do not use off-site disaster recovery; if empty, the default is not to use off-site disaster recovery.",
			},

			"is_create_domain": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "0 does not create a domain name, 1 creates a domain name.",
			},

			"domain": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Required if you want to create a domain name.",
			},

			"fw_cidr_info": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Specify the network segment information used by the firewall.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"fw_cidr_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The type of network segment used by the firewall. The values VpcSelf/Assis/Custom respectively represent own network segment priority/extended network segment priority/custom.",
						},
						"fw_cidr_lst": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Specify the network segment of the firewall for each vpc.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"vpc_id": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Vpc id.",
									},
									"fw_cidr": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Firewall network segment, at least /24 network segment.",
									},
								},
							},
						},
						"com_fw_cidr": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Other firewalls occupy the network segment, which is usually the network segment specified when the firewall needs to exclusively occupy the vpc.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudCfwNatfwCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cfw_natfw.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request  = cfw.NewCreateNatFwInstanceWithDomainRequest()
		response = cfw.NewCreateNatFwInstanceWithDomainResponse()
		natinsId string
	)
	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("width"); ok {
		request.Width = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("mode"); ok {
		request.Mode = helper.IntInt64(v.(int))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "new_mode_items"); ok {
		newModeItems := cfw.NewModeItems{}
		if v, ok := dMap["vpc_list"]; ok {
			vpcListSet := v.(*schema.Set).List()
			for i := range vpcListSet {
				vpcList := vpcListSet[i].(string)
				newModeItems.VpcList = append(newModeItems.VpcList, &vpcList)
			}
		}
		if v, ok := dMap["eips"]; ok {
			eipsSet := v.(*schema.Set).List()
			for i := range eipsSet {
				eips := eipsSet[i].(string)
				newModeItems.Eips = append(newModeItems.Eips, &eips)
			}
		}
		if v, ok := dMap["add_count"]; ok {
			newModeItems.AddCount = helper.IntInt64(v.(int))
		}
		request.NewModeItems = &newModeItems
	}

	if v, ok := d.GetOk("nat_gw_list"); ok {
		natGwListSet := v.(*schema.Set).List()
		for i := range natGwListSet {
			natGwList := natGwListSet[i].(string)
			request.NatGwList = append(request.NatGwList, &natGwList)
		}
	}

	if v, ok := d.GetOk("zone"); ok {
		request.Zone = helper.String(v.(string))
	}

	if v, ok := d.GetOk("zone_bak"); ok {
		request.ZoneBak = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("cross_a_zone"); ok {
		request.CrossAZone = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("is_create_domain"); ok {
		request.IsCreateDomain = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("domain"); ok {
		request.Domain = helper.String(v.(string))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "fw_cidr_info"); ok {
		fwCidrInfo := cfw.FwCidrInfo{}
		if v, ok := dMap["fw_cidr_type"]; ok {
			fwCidrInfo.FwCidrType = helper.String(v.(string))
		}
		if v, ok := dMap["fw_cidr_lst"]; ok {
			for _, item := range v.([]interface{}) {
				fwCidrLstMap := item.(map[string]interface{})
				fwVpcCidr := cfw.FwVpcCidr{}
				if v, ok := fwCidrLstMap["vpc_id"]; ok {
					fwVpcCidr.VpcId = helper.String(v.(string))
				}
				if v, ok := fwCidrLstMap["fw_cidr"]; ok {
					fwVpcCidr.FwCidr = helper.String(v.(string))
				}
				fwCidrInfo.FwCidrLst = append(fwCidrInfo.FwCidrLst, &fwVpcCidr)
			}
		}
		if v, ok := dMap["com_fw_cidr"]; ok {
			fwCidrInfo.ComFwCidr = helper.String(v.(string))
		}
		request.FwCidrInfo = &fwCidrInfo
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCfwClient().CreateNatFwInstanceWithDomain(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create cfw natfw failed, reason:%+v", logId, err)
		return err
	}

	natinsId = *response.Response.NatinsId
	d.SetId(natinsId)

	service := CfwService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{0}, 30*readRetryTimeout, time.Second, service.CfwNatfwStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudCfwNatfwRead(d, meta)
}

func resourceTencentCloudCfwNatfwRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cfw_natfw.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CfwService{client: meta.(*TencentCloudClient).apiV3Conn}

	natfwId := d.Id()

	natfw, err := service.DescribeCfwNatfwById(ctx, natinsId)
	if err != nil {
		return err
	}

	if natfw == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CfwNatfw` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if natfw.Name != nil {
		_ = d.Set("name", natfw.Name)
	}

	if natfw.Width != nil {
		_ = d.Set("width", natfw.Width)
	}

	if natfw.Mode != nil {
		_ = d.Set("mode", natfw.Mode)
	}

	if natfw.NewModeItems != nil {
		newModeItemsMap := map[string]interface{}{}

		if natfw.NewModeItems.VpcList != nil {
			newModeItemsMap["vpc_list"] = natfw.NewModeItems.VpcList
		}

		if natfw.NewModeItems.Eips != nil {
			newModeItemsMap["eips"] = natfw.NewModeItems.Eips
		}

		if natfw.NewModeItems.AddCount != nil {
			newModeItemsMap["add_count"] = natfw.NewModeItems.AddCount
		}

		_ = d.Set("new_mode_items", []interface{}{newModeItemsMap})
	}

	if natfw.NatGwList != nil {
		_ = d.Set("nat_gw_list", natfw.NatGwList)
	}

	if natfw.Zone != nil {
		_ = d.Set("zone", natfw.Zone)
	}

	if natfw.ZoneBak != nil {
		_ = d.Set("zone_bak", natfw.ZoneBak)
	}

	if natfw.CrossAZone != nil {
		_ = d.Set("cross_a_zone", natfw.CrossAZone)
	}

	if natfw.IsCreateDomain != nil {
		_ = d.Set("is_create_domain", natfw.IsCreateDomain)
	}

	if natfw.Domain != nil {
		_ = d.Set("domain", natfw.Domain)
	}

	if natfw.FwCidrInfo != nil {
		fwCidrInfoMap := map[string]interface{}{}

		if natfw.FwCidrInfo.FwCidrType != nil {
			fwCidrInfoMap["fw_cidr_type"] = natfw.FwCidrInfo.FwCidrType
		}

		if natfw.FwCidrInfo.FwCidrLst != nil {
			fwCidrLstList := []interface{}{}
			for _, fwCidrLst := range natfw.FwCidrInfo.FwCidrLst {
				fwCidrLstMap := map[string]interface{}{}

				if fwCidrLst.VpcId != nil {
					fwCidrLstMap["vpc_id"] = fwCidrLst.VpcId
				}

				if fwCidrLst.FwCidr != nil {
					fwCidrLstMap["fw_cidr"] = fwCidrLst.FwCidr
				}

				fwCidrLstList = append(fwCidrLstList, fwCidrLstMap)
			}

			fwCidrInfoMap["fw_cidr_lst"] = []interface{}{fwCidrLstList}
		}

		if natfw.FwCidrInfo.ComFwCidr != nil {
			fwCidrInfoMap["com_fw_cidr"] = natfw.FwCidrInfo.ComFwCidr
		}

		_ = d.Set("fw_cidr_info", []interface{}{fwCidrInfoMap})
	}

	return nil
}

func resourceTencentCloudCfwNatfwUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cfw_natfw.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := cfw.NewModifyNatInstanceRequest()

	natfwId := d.Id()

	request.NatinsId = &natinsId

	immutableArgs := []string{"name", "width", "mode", "new_mode_items", "nat_gw_list", "zone", "zone_bak", "cross_a_zone", "is_create_domain", "domain", "fw_cidr_info"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCfwClient().ModifyNatInstance(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update cfw natfw failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudCfwNatfwRead(d, meta)
}

func resourceTencentCloudCfwNatfwDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cfw_natfw.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CfwService{client: meta.(*TencentCloudClient).apiV3Conn}
	natfwId := d.Id()

	if err := service.DeleteCfwNatfwById(ctx, natinsId); err != nil {
		return err
	}

	return nil
}
