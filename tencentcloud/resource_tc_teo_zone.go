/*
Provides a resource to create a teo zone

Example Usage

```hcl
resource "tencentcloud_teo_zone" "zone" {
  zone_name = "toutiao2.com"
  plan_type = "sta"
  type      = "full"
  paused    = false
#  vanity_name_servers {
#    switch = ""
#    servers = ""
#
#  }
  cname_speed_up = "enabled"
#  tags {
#    tag_key = ""
#    tag_value = ""
#
#  }
  tags = {
    "createdBy" = "terraform"
  }
}

```
Import

teo zone can be imported using the id, e.g.
```
$ terraform import tencentcloud_teo_zone.zone zone_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	teo "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudTeoZone() *schema.Resource {
	return &schema.Resource{
		Read:   resourceTencentCloudTeoZoneRead,
		Create: resourceTencentCloudTeoZoneCreate,
		Update: resourceTencentCloudTeoZoneUpdate,
		Delete: resourceTencentCloudTeoZoneDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Site ID.",
			},

			"zone_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Site name.",
			},

			"plan_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Plan type of the zone. See details in data source `zone_available_plans`.",
			},

			"area": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Acceleration area of the zone. Valid values: `mainland`, `overseas`.",
			},

			"original_name_servers": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed:    true,
				Description: "Name server used by the site.",
			},

			"name_servers": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed:    true,
				Description: "List of name servers assigned by Tencent Cloud.",
			},

			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Site status. Valid values:- `active`: NS is switched.- `pending`: NS is not switched.- `moved`: NS is moved.- `deactivated`: this site is blocked.",
			},

			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies how the site is connected to EdgeOne.- `full`: The site is connected via NS.- `partial`: The site is connected via CNAME.",
			},

			"paused": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Indicates whether the site is disabled.",
			},

			"created_on": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Site creation date.",
			},

			"modified_on": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Site modification date.",
			},

			"vanity_name_servers": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "User-defined name server information. Note: This field may return null, indicating that no valid value can be obtained.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"switch": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Whether to enable the custom name server.- `on`: Enable.- `off`: Disable.",
						},
						"servers": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional:    true,
							Description: "List of custom name servers.",
						},
					},
				},
			},

			"vanity_name_servers_ips": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "User-defined name server IP information. Note: This field may return null, indicating that no valid value can be obtained.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the custom name server.",
						},
						"ipv4": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "IPv4 address of the custom name server.",
						},
					},
				},
			},

			"cname_speed_up": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies whether CNAME acceleration is enabled. Valid values: `enabled`, `disabled`.",
			},

			"cname_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Ownership verification status of the site when it accesses via CNAME.- `finished`: The site is verified.- `pending`: The site is waiting for verification.",
			},

			"resources": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Billing resources of the zone.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Resource ID.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Resource creation date.",
						},
						"pay_mode": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Resource pay mode. Valid values:- `0`: post pay mode.",
						},
						"enable_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Enable time of the resource.",
						},
						"expire_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Expire time of the resource.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Status of the resource. Valid values: `normal`, `isolated`, `destroyed`.",
						},
						"sv": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Price inquiry parameters.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Parameter Key.",
									},
									"value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Parameter Value.",
									},
								},
							},
						},
						"auto_renew_flag": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Whether to automatically renew. Valid values:- `0`: Default.- `1`: Enable automatic renewal.- `2`: Disable automatic renewal.",
						},
						"plan_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Associated plan ID.",
						},
						"area": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Valid values: `mainland`, `overseas`.",
						},
					},
				},
			},

			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Tag description list.",
			},
		},
	}
}

func resourceTencentCloudTeoZoneCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_teo_zone.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request  = teo.NewCreateZoneRequest()
		response *teo.CreateZoneResponse
		zoneId   string
		zoneName string
		planType string
	)

	if v, ok := d.GetOk("zone_name"); ok {
		zoneName = v.(string)
		request.ZoneName = &zoneName
	}

	if v, ok := d.GetOk("plan_type"); ok {
		planType = v.(string)
	}

	if v, _ := d.GetOk("type"); v != nil {
		request.Type = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTeoClient().CreateZone(request)
		if e != nil {
			if isExpectError(e, []string{"ResourceInUse", "ResourceInUse.Others"}) {
				return resource.NonRetryableError(e)
			}
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create teo zone failed, reason:%+v", logId, err)
		return err
	}

	zoneId = *response.Response.ZoneId
	d.SetId(zoneId)

	if zoneId != "" {
		var planRequest = teo.NewCreatePlanForZoneRequest()
		planRequest.ZoneId = &zoneId
		planRequest.PlanType = &planType
		planErr := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			result, e := meta.(*TencentCloudClient).apiV3Conn.UseTeoClient().CreatePlanForZone(planRequest)
			if e != nil {
				return retryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
					logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if planErr != nil {
			log.Printf("[CRITAL]%s create teo zone failed, reason:%+v", logId, planErr)
			return planErr
		}
	}

	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		tagService := TagService{client: meta.(*TencentCloudClient).apiV3Conn}
		resourceName := fmt.Sprintf("qcs::teo::uin/:zone/%s", zoneId)
		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
	}
	return resourceTencentCloudTeoZoneRead(d, meta)
}

func resourceTencentCloudTeoZoneRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_teo_zone.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TeoService{client: meta.(*TencentCloudClient).apiV3Conn}

	zoneId := d.Id()

	zone, err := service.DescribeTeoZone(ctx, zoneId)

	if err != nil {
		return err
	}

	if zone == nil {
		d.SetId("")
		return fmt.Errorf("resource `zone` %s does not exist", zoneId)
	}

	if zone.ZoneId != nil {
		_ = d.Set("zone_id", zone.ZoneId)
	}

	if zone.ZoneName != nil {
		_ = d.Set("zone_name", zone.ZoneName)
	}

	if zone.Area != nil {
		_ = d.Set("area", zone.Area)
	}

	if zone.OriginalNameServers != nil {
		_ = d.Set("original_name_servers", zone.OriginalNameServers)
	}

	if zone.NameServers != nil {
		_ = d.Set("name_servers", zone.NameServers)
	}

	if zone.Status != nil {
		_ = d.Set("status", zone.Status)
	}

	if zone.Type != nil {
		_ = d.Set("type", zone.Type)
	}

	if zone.Paused != nil {
		_ = d.Set("paused", zone.Paused)
	}

	if zone.CreatedOn != nil {
		_ = d.Set("created_on", zone.CreatedOn)
	}

	if zone.ModifiedOn != nil {
		_ = d.Set("modified_on", zone.ModifiedOn)
	}

	if zone.VanityNameServers != nil {
		vanityNameServersMap := map[string]interface{}{}
		if zone.VanityNameServers.Switch != nil {
			vanityNameServersMap["switch"] = zone.VanityNameServers.Switch
		}
		if zone.VanityNameServers.Servers != nil {
			vanityNameServersMap["servers"] = zone.VanityNameServers.Servers
		}

		_ = d.Set("vanity_name_servers", []interface{}{vanityNameServersMap})
	}

	if zone.VanityNameServersIps != nil {
		vanityNameServersIpsList := []interface{}{}
		for _, vanityNameServersIps := range zone.VanityNameServersIps {
			vanityNameServersIpsMap := map[string]interface{}{}
			if vanityNameServersIps.Name != nil {
				vanityNameServersIpsMap["name"] = vanityNameServersIps.Name
			}
			if vanityNameServersIps.IPv4 != nil {
				vanityNameServersIpsMap["ipv4"] = vanityNameServersIps.IPv4
			}

			vanityNameServersIpsList = append(vanityNameServersIpsList, vanityNameServersIpsMap)
		}
		_ = d.Set("vanity_name_servers_ips", vanityNameServersIpsList)
	}

	if zone.CnameSpeedUp != nil {
		_ = d.Set("cname_speed_up", zone.CnameSpeedUp)
	}

	if zone.CnameStatus != nil {
		_ = d.Set("cname_status", zone.CnameStatus)
	}

	if zone.Resources != nil {
		resourcesList := []interface{}{}
		for _, resources := range zone.Resources {
			resourcesMap := map[string]interface{}{}
			if resources.Id != nil {
				resourcesMap["id"] = resources.Id
			}
			if resources.CreateTime != nil {
				resourcesMap["create_time"] = resources.CreateTime
			}
			if resources.PayMode != nil {
				resourcesMap["pay_mode"] = resources.PayMode
			}
			if resources.EnableTime != nil {
				resourcesMap["enable_time"] = resources.EnableTime
			}
			if resources.ExpireTime != nil {
				resourcesMap["expire_time"] = resources.ExpireTime
			}
			if resources.Status != nil {
				resourcesMap["status"] = resources.Status
			}
			if resources.Sv != nil {
				svList := []interface{}{}
				for _, sv := range resources.Sv {
					svMap := map[string]interface{}{}
					if sv.Key != nil {
						svMap["key"] = sv.Key
					}
					if sv.Value != nil {
						svMap["value"] = sv.Value
					}

					svList = append(svList, svMap)
				}
				resourcesMap["sv"] = svList
			}
			if resources.AutoRenewFlag != nil {
				resourcesMap["auto_renew_flag"] = resources.AutoRenewFlag
			}
			if resources.PlanId != nil {
				resourcesMap["plan_id"] = resources.PlanId
			}
			if resources.Area != nil {
				resourcesMap["area"] = resources.Area
			}

			resourcesList = append(resourcesList, resourcesMap)
		}
		_ = d.Set("resources", resourcesList)
	}

	if zone.Area != nil {
		_ = d.Set("area", zone.Area)
	}

	tcClient := meta.(*TencentCloudClient).apiV3Conn
	tagService := &TagService{client: tcClient}
	tags, err := tagService.DescribeResourceTags(ctx, "teo", "zone", "", zoneId)
	if err != nil {
		return err
	}
	_ = d.Set("tags", tags)

	return nil
}

func resourceTencentCloudTeoZoneUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_teo_zone.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	request := teo.NewModifyZoneRequest()

	zoneId := d.Id()
	request.ZoneId = &zoneId

	if d.HasChange("zone_name") {
		return fmt.Errorf("`zone_name` do not support change now.")
	}

	if d.HasChange("plan_type") {
		log.Printf("[WARN] change `plan_type` is not supported now.")
		_ = d.Set("plan_type", d.Get("plan_type"))
	}

	if d.HasChange("type") {
		if v, ok := d.GetOk("type"); ok {
			request.Type = helper.String(v.(string))
		}
	}

	if d.HasChange("vanity_name_servers") {
		if dMap, ok := helper.InterfacesHeadMap(d, "vanity_name_servers"); ok {
			vanityNameServers := teo.VanityNameServers{}
			if v, ok := dMap["switch"]; ok {
				vanityNameServers.Switch = helper.String(v.(string))
			}
			if v, ok := dMap["servers"]; ok {
				serversSet := v.(*schema.Set).List()
				for i := range serversSet {
					servers := serversSet[i].(string)
					vanityNameServers.Servers = append(vanityNameServers.Servers, &servers)
				}
			}
			request.VanityNameServers = &vanityNameServers
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTeoClient().ModifyZone(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create teo zone failed, reason:%+v", logId, err)
		return err
	}

	if d.HasChange("paused") {
		if v := d.Get("paused"); v != nil {
			req := teo.NewModifyZoneStatusRequest()
			req.ZoneId, req.Paused = &zoneId, helper.Bool(v.(bool))
			_, e := meta.(*TencentCloudClient).apiV3Conn.UseTeoClient().ModifyZoneStatus(req)
			if e != nil {
				log.Printf("[CRITAL]%s modify zone status failed, reason:%+v", logId, err)
				return err
			}
		}
	}

	if d.HasChange("cname_speed_up") {
		if v, ok := d.GetOk("cname_speed_up"); ok {
			req := teo.NewModifyZoneCnameSpeedUpRequest()
			req.ZoneId, req.Status = &zoneId, helper.String(v.(string))
			_, e := meta.(*TencentCloudClient).apiV3Conn.UseTeoClient().ModifyZoneCnameSpeedUp(req)
			if e != nil {
				log.Printf("[CRITAL]%s modify zone cname_speed_up failed, reason:%+v", logId, err)
				return err
			}
		}
	}

	if d.HasChange("tags") {
		tcClient := meta.(*TencentCloudClient).apiV3Conn
		tagService := &TagService{client: tcClient}
		oldTags, newTags := d.GetChange("tags")
		replaceTags, deleteTags := diffTags(oldTags.(map[string]interface{}), newTags.(map[string]interface{}))
		resourceName := BuildTagResourceName("teo", "zone", "", zoneId)
		if err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags); err != nil {
			return err
		}
	}

	return resourceTencentCloudTeoZoneRead(d, meta)
}

func resourceTencentCloudTeoZoneDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_teo_zone.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TeoService{client: meta.(*TencentCloudClient).apiV3Conn}
	zoneId := d.Id()

	if err := service.DeleteTeoZoneById(ctx, zoneId); err != nil {
		return err
	}

	return nil
}
