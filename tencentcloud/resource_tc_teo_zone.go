/*
Provides a resource to create a teo zone

Example Usage

```hcl
resource "tencentcloud_teo_zone" "zone" {
  name           = "sfurnace.work"
  plan_type      = "ent_cm_with_bot"
  type           = "full"
  paused         = false
  cname_speed_up = "enabled"

  #  vanity_name_servers {
  #    switch  = "on"
  #    servers = ["2.2.2.2"]
  #  }
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
	"errors"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	teo "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220106"
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
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Site name.",
			},

			"plan_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Plan type of the zone. See details in data source `zone_available_plans`.",
			},

			"original_name_servers": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed:    true,
				Description: "List of name servers used.",
			},

			"name_servers": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed:    true,
				Description: "List of name servers assigned to users by Tencent Cloud.",
			},

			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Site status.",
			},

			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies how the site is connected to EdgeOne.",
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
				Computed:    true,
				Description: "User-defined name server information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"switch": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Whether to enable the custom name server.",
						},
						"servers": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional:    true,
							Computed:    true,
							Description: "List of custom name servers.",
						},
					},
				},
			},

			"vanity_name_servers_ips": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "User-defined name server IP information.",
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
				Description: "Specifies whether to enable CNAME acceleration, enabled: Enable; disabled: Disable.",
			},

			"cname_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Ownership verification status of the site when it accesses via CNAME.",
			},

			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Tag description list.",
			},
			"area": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Acceleration area of the zone. Valid values: `mainland`, `overseas`.",
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
	)

	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOk("type"); ok {
		request.Type = helper.String(v.(string))
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
	}

	result, e := meta.(*TencentCloudClient).apiV3Conn.UseTeoClient().CreateZone(request)
	if e != nil {
		return e
	} else {
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
	}
	response = result

	if response == nil || response.Response == nil || response.Response.Id == nil {
		return errors.New("CreateZone create teo zone failed")
	}

	zoneId := *response.Response.Id

	var planRequest = teo.NewCreatePlanForZoneRequest()
	planRequest.ZoneId = &zoneId
	if v, ok := d.GetOk("plan_type"); ok {
		planRequest.PlanType = helper.String(v.(string))
	}

	resultPlan, err := meta.(*TencentCloudClient).apiV3Conn.UseTeoClient().CreatePlanForZone(planRequest)
	if e != nil {
		log.Printf("[CRITAL]%s create teo zone plan failed, reason:%+v", logId, e)
		return err
	} else {
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), resultPlan.ToJsonString())
	}

	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		tagService := TagService{client: meta.(*TencentCloudClient).apiV3Conn}
		region := meta.(*TencentCloudClient).apiV3Conn.Region
		resourceName := fmt.Sprintf("qcs::teo:%s:uin/:zone/%s", region, zoneId)
		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
	}
	d.SetId(zoneId)
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

	if zone.Id != nil {
		_ = d.Set("id", zone.Id)
	}

	if zone.Name != nil {
		_ = d.Set("name", zone.Name)
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

	tcClient := meta.(*TencentCloudClient).apiV3Conn
	tagService := &TagService{client: tcClient}
	tags, err := tagService.DescribeResourceTags(ctx, "teo", "zone", tcClient.Region, d.Id())
	if err != nil {
		return err
	}
	_ = d.Set("tags", tags)

	if zone.Area != nil {
		_ = d.Set("area", zone.Area)
	}

	return nil
}

func resourceTencentCloudTeoZoneUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_teo_zone.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	request := teo.NewModifyZoneRequest()

	request.Id = helper.String(d.Id())

	if d.HasChange("name") {
		return fmt.Errorf("`name` do not support change now.")
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
		return err
	}

	if d.HasChange("tags") {
		tcClient := meta.(*TencentCloudClient).apiV3Conn
		tagService := &TagService{client: tcClient}
		oldTags, newTags := d.GetChange("tags")
		replaceTags, deleteTags := diffTags(oldTags.(map[string]interface{}), newTags.(map[string]interface{}))
		resourceName := BuildTagResourceName("teo", "zone", tcClient.Region, d.Id())
		if err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags); err != nil {
			return err
		}
	}

	if d.HasChange("cname_speed_up") {
		requestCnameSpeedUp := teo.NewModifyZoneCnameSpeedUpRequest()
		requestCnameSpeedUp.Id = helper.String(d.Id())
		if v, ok := d.GetOk("cname_speed_up"); ok {
			requestCnameSpeedUp.Status = helper.String(v.(string))
		}
		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			result, e := meta.(*TencentCloudClient).apiV3Conn.UseTeoClient().ModifyZoneCnameSpeedUp(requestCnameSpeedUp)
			if e != nil {
				return retryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
					logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})

		if err != nil {
			return err
		}
	}

	if d.HasChange("paused") {
		requestPaused := teo.NewModifyZoneStatusRequest()
		requestPaused.Id = helper.String(d.Id())
		v, _ := d.GetOk("paused")
		requestPaused.Paused = helper.Bool(v.(bool))
		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			result, e := meta.(*TencentCloudClient).apiV3Conn.UseTeoClient().ModifyZoneStatus(requestPaused)
			if e != nil {
				return retryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
					logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})

		if err != nil {
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
