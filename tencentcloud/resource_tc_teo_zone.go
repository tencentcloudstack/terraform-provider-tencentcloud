/*
Provides a resource to create a teo zone

Example Usage

```hcl
resource "tencentcloud_teo_zone" "zone" {
    zone_name = &lt;nil&gt;
  plan_type = &lt;nil&gt;
          type = &lt;nil&gt;
  paused = &lt;nil&gt;
      vanity_name_servers {
		switch = &lt;nil&gt;
		servers = &lt;nil&gt;

  }
    cname_speed_up = &lt;nil&gt;
    tags {
		tag_key = &lt;nil&gt;
		tag_value = &lt;nil&gt;

  }
      tags = {
    "createdBy" = "terraform"
  }
}
```

Import

teo zone can be imported using the id, e.g.

```
terraform import tencentcloud_teo_zone.zone zone_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	teo "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudTeoZone() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTeoZoneCreate,
		Read:   resourceTencentCloudTeoZoneRead,
		Update: resourceTencentCloudTeoZoneUpdate,
		Delete: resourceTencentCloudTeoZoneDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Site ID.",
			},

			"zone_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Site name.",
			},

			"plan_type": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Plan type of the zone. See details in data source `zone_available_plans`.",
			},

			"area": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Acceleration area of the zone. Valid values: `mainland`, `overseas`.",
			},

			"original_name_servers": {
				Computed: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Name server used by the site.",
			},

			"name_servers": {
				Computed: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "List of name servers assigned by Tencent Cloud.",
			},

			"status": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Site status. Valid values:- `active`: NS is switched.- `pending`: NS is not switched.- `moved`: NS is moved.- `deactivated`: this site is blocked.",
			},

			"type": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Specifies how the site is connected to EdgeOne.- `full`: The site is connected via NS.- `partial`: The site is connected via CNAME.",
			},

			"paused": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeBool,
				Description: "Indicates whether the site is disabled.",
			},

			"created_on": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Site creation date.",
			},

			"modified_on": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Site modification date.",
			},

			"vanity_name_servers": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
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
				Computed:    true,
				Type:        schema.TypeList,
				Description: "User-defined name server IP information. Note: This field may return null, indicating that no valid value can be obtained.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the custom name server.",
						},
						"i_pv4": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "IPv4 address of the custom name server.",
						},
					},
				},
			},

			"cname_speed_up": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Specifies whether CNAME acceleration is enabled. Valid values: `enabled`, `disabled`.",
			},

			"cname_status": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Ownership verification status of the site when it accesses via CNAME.- `finished`: The site is verified.- `pending`: The site is waiting for verification.",
			},

			"tags": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "Tag list of the site.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tag_key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Tag key. Note: This field may return null, indicating that no valid value can be obtained.",
						},
						"tag_value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Tag value. Note: This field may return null, indicating that no valid value can be obtained.",
						},
					},
				},
			},

			"resources": {
				Computed:    true,
				Type:        schema.TypeList,
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

			"area": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Site access area. Valid values: `global`, `mainland`, `overseas`.",
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
		createZoneRequest  = teo.NewCreateZoneRequest()
		createZoneResponse = teo.NewCreateZoneResponse()
	)
	if v, ok := d.GetOk("zone_name"); ok {
		request.ZoneName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("plan_type"); ok {
		request.PlanType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("type"); ok {
		request.Type = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("paused"); ok {
		request.Paused = helper.Bool(v.(bool))
	}

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

	if v, ok := d.GetOk("cname_speed_up"); ok {
		request.CnameSpeedUp = helper.String(v.(string))
	}

	if v, ok := d.GetOk("tags"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			tag := teo.Tag{}
			if v, ok := dMap["tag_key"]; ok {
				tag.TagKey = helper.String(v.(string))
			}
			if v, ok := dMap["tag_value"]; ok {
				tag.TagValue = helper.String(v.(string))
			}
			request.Tags = append(request.Tags, &tag)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTeoClient().CreateZone(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
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

	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		tagService := TagService{client: meta.(*TencentCloudClient).apiV3Conn}
		region := meta.(*TencentCloudClient).apiV3Conn.Region
		resourceName := fmt.Sprintf("qcs::teo:%s:uin/:zone/%s", region, d.Id())
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

	zone, err := service.DescribeTeoZoneById(ctx, zoneId)
	if err != nil {
		return err
	}

	if zone == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `TeoZone` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if zone.ZoneId != nil {
		_ = d.Set("zone_id", zone.ZoneId)
	}

	if zone.ZoneName != nil {
		_ = d.Set("zone_name", zone.ZoneName)
	}

	if zone.PlanType != nil {
		_ = d.Set("plan_type", zone.PlanType)
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

			if zone.VanityNameServersIps.Name != nil {
				vanityNameServersIpsMap["name"] = zone.VanityNameServersIps.Name
			}

			if zone.VanityNameServersIps.IPv4 != nil {
				vanityNameServersIpsMap["i_pv4"] = zone.VanityNameServersIps.IPv4
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

	if zone.Tags != nil {
		tagsList := []interface{}{}
		for _, tags := range zone.Tags {
			tagsMap := map[string]interface{}{}

			if zone.Tags.TagKey != nil {
				tagsMap["tag_key"] = zone.Tags.TagKey
			}

			if zone.Tags.TagValue != nil {
				tagsMap["tag_value"] = zone.Tags.TagValue
			}

			tagsList = append(tagsList, tagsMap)
		}

		_ = d.Set("tags", tagsList)

	}

	if zone.Resources != nil {
		resourcesList := []interface{}{}
		for _, resources := range zone.Resources {
			resourcesMap := map[string]interface{}{}

			if zone.Resources.Id != nil {
				resourcesMap["id"] = zone.Resources.Id
			}

			if zone.Resources.CreateTime != nil {
				resourcesMap["create_time"] = zone.Resources.CreateTime
			}

			if zone.Resources.PayMode != nil {
				resourcesMap["pay_mode"] = zone.Resources.PayMode
			}

			if zone.Resources.EnableTime != nil {
				resourcesMap["enable_time"] = zone.Resources.EnableTime
			}

			if zone.Resources.ExpireTime != nil {
				resourcesMap["expire_time"] = zone.Resources.ExpireTime
			}

			if zone.Resources.Status != nil {
				resourcesMap["status"] = zone.Resources.Status
			}

			if zone.Resources.Sv != nil {
				svList := []interface{}{}
				for _, sv := range zone.Resources.Sv {
					svMap := map[string]interface{}{}

					if sv.Key != nil {
						svMap["key"] = sv.Key
					}

					if sv.Value != nil {
						svMap["value"] = sv.Value
					}

					svList = append(svList, svMap)
				}

				resourcesMap["sv"] = []interface{}{svList}
			}

			if zone.Resources.AutoRenewFlag != nil {
				resourcesMap["auto_renew_flag"] = zone.Resources.AutoRenewFlag
			}

			if zone.Resources.PlanId != nil {
				resourcesMap["plan_id"] = zone.Resources.PlanId
			}

			if zone.Resources.Area != nil {
				resourcesMap["area"] = zone.Resources.Area
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
	tags, err := tagService.DescribeResourceTags(ctx, "teo", "zone", tcClient.Region, d.Id())
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

	var (
		modifyZoneRequest  = teo.NewModifyZoneRequest()
		modifyZoneResponse = teo.NewModifyZoneResponse()
	)

	zoneId := d.Id()

	request.ZoneId = &zoneId

	immutableArgs := []string{"zone_id", "zone_name", "plan_type", "area", "original_name_servers", "name_servers", "status", "type", "paused", "created_on", "modified_on", "vanity_name_servers", "vanity_name_servers_ips", "cname_speed_up", "cname_status", "tags", "resources", "area"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("plan_type") {
		if v, ok := d.GetOk("plan_type"); ok {
			request.PlanType = helper.String(v.(string))
		}
	}

	if d.HasChange("type") {
		if v, ok := d.GetOk("type"); ok {
			request.Type = helper.String(v.(string))
		}
	}

	if d.HasChange("paused") {
		if v, ok := d.GetOkExists("paused"); ok {
			request.Paused = helper.Bool(v.(bool))
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

	if d.HasChange("cname_speed_up") {
		if v, ok := d.GetOk("cname_speed_up"); ok {
			request.CnameSpeedUp = helper.String(v.(string))
		}
	}

	if d.HasChange("tags") {
		if v, ok := d.GetOk("tags"); ok {
			for _, item := range v.([]interface{}) {
				tag := teo.Tag{}
				if v, ok := dMap["tag_key"]; ok {
					tag.TagKey = helper.String(v.(string))
				}
				if v, ok := dMap["tag_value"]; ok {
					tag.TagValue = helper.String(v.(string))
				}
				request.Tags = append(request.Tags, &tag)
			}
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTeoClient().ModifyZone(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update teo zone failed, reason:%+v", logId, err)
		return err
	}

	if d.HasChange("tags") {
		ctx := context.WithValue(context.TODO(), logIdKey, logId)
		tcClient := meta.(*TencentCloudClient).apiV3Conn
		tagService := &TagService{client: tcClient}
		oldTags, newTags := d.GetChange("tags")
		replaceTags, deleteTags := diffTags(oldTags.(map[string]interface{}), newTags.(map[string]interface{}))
		resourceName := BuildTagResourceName("teo", "zone", tcClient.Region, d.Id())
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
