/*
Use this data source to query detailed information of teo zoneDDoSPolicy

Example Usage

```hcl
data "tencentcloud_teo_zone_ddos_policy" "zoneDDoSPolicy" {
  zone_id = ""
}
```
*/
package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	teo "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"
)

func dataSourceTencentCloudTeoZoneDDoSPolicy() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTeoZoneDDoSPolicyRead,
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Site ID.",
			},

			//"app_id": {
			//	Type:        schema.TypeInt,
			//	Computed:    true,
			//	Description: "AppID of the account.",
			//},

			"shield_areas": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Shielded areas of the zone.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"zone_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Site ID.",
						},
						"policy_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Policy ID.",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Valid values: `domain`, `application`.",
						},
						"entity": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "When `Type` is `domain`, this field is `ZoneId`. When `Type` is `application`, this field is `ProxyId`. Note: This field may return null, indicating that no valid value can be obtained.",
						},
						"entity_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "When `Type` is `domain`, this field is `ZoneName`. When `Type` is `application`, this field is `ProxyName`. Note: This field may return null, indicating that no valid value can be obtained.",
						},
						"tcp_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "TCP forwarding rule number of layer 4 application.",
						},
						"udp_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "UDP forwarding rule number of layer 4 application.",
						},
						"application": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "DDoS layer 7 application.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"host": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Subdomain.",
									},
									"status": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Status of the subdomain. Valid values:- `init`: waiting to config NS.- `offline`: need to enable site accelerating.- `process`: processing the config deployment.- `online`: normal status. Note: This field may return null, indicating that no valid value can be obtained.",
									},
									"accelerate_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Acceleration function switch. Valid values:- `on`: Enable.- `off`: Disable.",
									},
									"security_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Security function switch. Valid values:- `on`: Enable.- `off`: Disable.",
									},
								},
							},
						},
					},
				},
			},

			"domains": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "All subdomain info. Note: This field may return null, indicating that no valid value can be obtained.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"host": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Subdomain.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Status of the subdomain. Valid values:- `init`: waiting to config NS.- `offline`: need to enable site accelerating.- `process`: processing the config deployment.- `online`: normal status. Note: This field may return null, indicating that no valid value can be obtained.",
						},
						"accelerate_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Acceleration function switch. Valid values:- `on`: Enable.- `off`: Disable.",
						},
						"security_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Security function switch. Valid values:- `on`: Enable.- `off`: Disable.",
						},
					},
				},
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudTeoZoneDDoSPolicyRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_teo_zone_ddos_policy.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	var zoneId string

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("zone_id"); ok {
		zoneId = v.(string)
		paramMap["zone_id"] = v
	}

	teoService := TeoService{client: meta.(*TencentCloudClient).apiV3Conn}

	var ddosPolicy *teo.DescribeZoneDDoSPolicyResponseParams
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		results, e := teoService.DescribeTeoZoneDDoSPolicyByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		ddosPolicy = results
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read Teo planInfo failed, reason:%+v", logId, err)
		return err
	}

	shieldAreasList := []interface{}{}
	if ddosPolicy != nil {
		for _, shieldAreas := range ddosPolicy.ShieldAreas {
			shieldAreasMap := map[string]interface{}{}
			if shieldAreas.ZoneId != nil {
				shieldAreasMap["zone_id"] = shieldAreas.ZoneId
			}
			if shieldAreas.PolicyId != nil {
				shieldAreasMap["policy_id"] = shieldAreas.PolicyId
			}
			if shieldAreas.Type != nil {
				shieldAreasMap["type"] = shieldAreas.Type
			}
			if shieldAreas.Entity != nil {
				shieldAreasMap["entity"] = shieldAreas.Entity
			}
			if shieldAreas.EntityName != nil {
				shieldAreasMap["entity_name"] = shieldAreas.EntityName
			}
			if shieldAreas.TcpNum != nil {
				shieldAreasMap["tcp_num"] = shieldAreas.TcpNum
			}
			if shieldAreas.UdpNum != nil {
				shieldAreasMap["udp_num"] = shieldAreas.UdpNum
			}
			if shieldAreas.DDoSHosts != nil {
				applicationList := []interface{}{}
				for _, ddosHost := range shieldAreas.DDoSHosts {
					applicationMap := map[string]interface{}{}
					if ddosHost.Host != nil {
						applicationMap["host"] = ddosHost.Host
					}
					if ddosHost.Status != nil {
						applicationMap["status"] = ddosHost.Status
					}
					if ddosHost.AccelerateType != nil {
						applicationMap["accelerate_type"] = ddosHost.AccelerateType
					}
					if ddosHost.SecurityType != nil {
						applicationMap["security_type"] = ddosHost.SecurityType
					}

					applicationList = append(applicationList, applicationMap)
				}
				shieldAreasMap["application"] = applicationList
			}

			shieldAreasList = append(shieldAreasList, shieldAreasMap)
		}
		_ = d.Set("shield_areas", shieldAreasList)
	}

	ddosHostList := []interface{}{}
	if ddosPolicy != nil {
		for _, planInfo := range ddosPolicy.DDoSHosts {
			ddosHostMap := map[string]interface{}{}
			if planInfo.Host != nil {
				ddosHostMap["host"] = planInfo.Host
			}
			if planInfo.Status != nil {
				ddosHostMap["status"] = planInfo.Status
			}
			if planInfo.AccelerateType != nil {
				ddosHostMap["accelerate_type"] = planInfo.AccelerateType
			}
			if planInfo.SecurityType != nil {
				ddosHostMap["security_type"] = planInfo.SecurityType
			}

			ddosHostList = append(ddosHostList, ddosHostMap)
		}
		_ = d.Set("domains", ddosHostList)
	}

	d.SetId(zoneId)

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), map[string]interface{}{
			"Shield_areas": shieldAreasList,
			"domains":      ddosHostList,
		}); e != nil {
			return e
		}
	}

	return nil
}
