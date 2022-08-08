/*
Provides a resource to create a teo originGroup

Example Usage

```hcl
resource "tencentcloud_teo_origin_group" "originGroup" {
  record {
		private_parameter {}
  }
  tags = {
    "createdBy" = "terraform"
  }
}

```
Import

teo originGroup can be imported using the id, e.g.
```
$ terraform import tencentcloud_teo_origin_group.originGroup originGroup_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	teo "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220106"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudTeoOriginGroup() *schema.Resource {
	return &schema.Resource{
		Read:   resourceTencentCloudTeoOriginGroupRead,
		Create: resourceTencentCloudTeoOriginGroupCreate,
		Update: resourceTencentCloudTeoOriginGroupUpdate,
		Delete: resourceTencentCloudTeoOriginGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"origin_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: ".",
			},

			"origin_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: ".",
			},

			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "area, weight.",
			},

			"record": {
				Type:        schema.TypeList,
				Required:    true,
				Description: ".",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"record_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: ".",
						},
						"record": {
							Type:        schema.TypeString,
							Required:    true,
							Description: ".",
						},
						"area": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Required:    true,
							Description: ".",
						},
						"weight": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "1-100.",
						},
						"port": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: ".",
						},
						"private": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: ".",
						},
						"private_parameter": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: ".",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Required:    true,
										Description: ".",
									},
									"value": {
										Type:        schema.TypeString,
										Required:    true,
										Description: ".",
									},
								},
							},
						},
					},
				},
			},

			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: ".",
			},

			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: ".",
			},

			"zone_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: ".",
			},

			"origin_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: ".",
			},

			"application_proxy_used": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: ".",
			},

			"load_balancing_used": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: ".",
			},

			"load_balancing_used_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "none, dns_only, proxied, both.",
			},

			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Tag description list.",
			},
		},
	}
}

func resourceTencentCloudTeoOriginGroupCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_teo_origin_group.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request  = teo.NewCreateOriginGroupRequest()
		response *teo.CreateOriginGroupResponse
		zoneId   string
	)

	if v, ok := d.GetOk("origin_name"); ok {
		request.OriginName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("type"); ok {
		request.Type = helper.String(v.(string))
	}

	if v, ok := d.GetOk("record"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			originRecord := teo.OriginRecord{}
			if v, ok := dMap["record"]; ok {
				originRecord.Record = helper.String(v.(string))
			}
			if v, ok := dMap["area"]; ok {
				areaSet := v.(*schema.Set).List()
				for i := range areaSet {
					area := areaSet[i].(string)
					originRecord.Area = append(originRecord.Area, &area)
				}
			}
			if v, ok := dMap["weight"]; ok {
				originRecord.Weight = helper.IntUint64(v.(int))
			}
			if v, ok := dMap["port"]; ok {
				originRecord.Port = helper.IntUint64(v.(int))
			}
			if v, ok := dMap["private"]; ok {
				originRecord.Private = helper.Bool(v.(bool))
			}
			if v, ok := dMap["private_parameter"]; ok {
				for _, item := range v.([]interface{}) {
					PrivateParameterMap := item.(map[string]interface{})
					originRecordPrivateParameter := teo.OriginRecordPrivateParameter{}
					if v, ok := PrivateParameterMap["name"]; ok {
						originRecordPrivateParameter.Name = helper.String(v.(string))
					}
					if v, ok := PrivateParameterMap["value"]; ok {
						originRecordPrivateParameter.Value = helper.String(v.(string))
					}
					originRecord.PrivateParameter = append(originRecord.PrivateParameter, &originRecordPrivateParameter)
				}
			}

		}
	}

	if v, ok := d.GetOk("zone_id"); ok {
		zoneId = v.(string)
		request.ZoneId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("origin_type"); ok {
		request.OriginType = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTeoClient().CreateOriginGroup(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create teo originGroup failed, reason:%+v", logId, err)
		return err
	}

	originGroupId := *response.Response.OriginId

	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		tagService := TagService{client: meta.(*TencentCloudClient).apiV3Conn}
		region := meta.(*TencentCloudClient).apiV3Conn.Region
		resourceName := fmt.Sprintf("qcs::teo:%s:uin/:zone/%s", region, originGroupId)
		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
	}
	d.SetId(zoneId + "#" + originGroupId)
	return resourceTencentCloudTeoOriginGroupRead(d, meta)
}

func resourceTencentCloudTeoOriginGroupRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_teo_originGroup.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TeoService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	zoneId := idSplit[0]
	originGroupId := idSplit[1]

	originGroup, err := service.DescribeTeoOriginGroup(ctx, zoneId, originGroupId)

	if err != nil {
		return err
	}

	if originGroup == nil {
		d.SetId("")
		return fmt.Errorf("resource `originGroup` %s does not exist", originGroupId)
	}

	if originGroup.OriginId != nil {
		_ = d.Set("origin_id", originGroup.OriginId)
	}

	if originGroup.OriginName != nil {
		_ = d.Set("origin_name", originGroup.OriginName)
	}

	if originGroup.Type != nil {
		_ = d.Set("type", originGroup.Type)
	}

	if originGroup.Record != nil {
		recordList := []interface{}{}
		for _, record := range originGroup.Record {
			recordMap := map[string]interface{}{}
			if record.RecordId != nil {
				recordMap["record_id"] = record.RecordId
			}
			if record.Record != nil {
				recordMap["record"] = record.Record
			}
			if record.Area != nil {
				recordMap["area"] = record.Area
			}
			if record.Weight != nil {
				recordMap["weight"] = record.Weight
			}
			if record.Port != nil {
				recordMap["port"] = record.Port
			}
			if record.Private != nil {
				recordMap["private"] = record.Private
			}
			if record.PrivateParameter != nil {
				privateParameterList := []interface{}{}
				for _, privateParameter := range record.PrivateParameter {
					privateParameterMap := map[string]interface{}{}
					if privateParameter.Name != nil {
						privateParameterMap["name"] = privateParameter.Name
					}
					if privateParameter.Value != nil {
						privateParameterMap["value"] = privateParameter.Value
					}

					privateParameterList = append(privateParameterList, privateParameterMap)
				}
				recordMap["private_parameter"] = privateParameterList
			}

			recordList = append(recordList, recordMap)
		}
		_ = d.Set("record", recordList)
	}

	if originGroup.UpdateTime != nil {
		_ = d.Set("update_time", originGroup.UpdateTime)
	}

	if originGroup.ZoneId != nil {
		_ = d.Set("zone_id", originGroup.ZoneId)
	}

	if originGroup.ZoneName != nil {
		_ = d.Set("zone_name", originGroup.ZoneName)
	}

	if originGroup.OriginType != nil {
		_ = d.Set("origin_type", originGroup.OriginType)
	}

	if originGroup.ApplicationProxyUsed != nil {
		_ = d.Set("application_proxy_used", originGroup.ApplicationProxyUsed)
	}

	if originGroup.LoadBalancingUsed != nil {
		_ = d.Set("load_balancing_used", originGroup.LoadBalancingUsed)
	}

	if originGroup.LoadBalancingUsedType != nil {
		_ = d.Set("load_balancing_used_type", originGroup.LoadBalancingUsedType)
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

func resourceTencentCloudTeoOriginGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_teo_origin_group.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	request := teo.NewModifyOriginGroupRequest()

	request.ZoneId = helper.String(d.Id())

	if d.HasChange("origin_id") {
		return fmt.Errorf("`origin_id` do not support change now.")
	}

	if d.HasChange("origin_name") {
		if v, ok := d.GetOk("origin_name"); ok {
			request.OriginName = helper.String(v.(string))
		}
	}

	if d.HasChange("type") {
		if v, ok := d.GetOk("type"); ok {
			request.Type = helper.String(v.(string))
		}
	}

	if d.HasChange("record") {
		if v, ok := d.GetOk("record"); ok {
			for _, item := range v.([]interface{}) {
				dMap := item.(map[string]interface{})
				originRecord := teo.OriginRecord{}
				if v, ok := dMap["record"]; ok {
					originRecord.Record = helper.String(v.(string))
				}
				if v, ok := dMap["area"]; ok {
					areaSet := v.(*schema.Set).List()
					for i := range areaSet {
						area := areaSet[i].(string)
						originRecord.Area = append(originRecord.Area, &area)
					}
				}
				if v, ok := dMap["weight"]; ok {
					originRecord.Weight = helper.IntUint64(v.(int))
				}
				if v, ok := dMap["port"]; ok {
					originRecord.Port = helper.IntUint64(v.(int))
				}
				if v, ok := dMap["private"]; ok {
					originRecord.Private = helper.Bool(v.(bool))
				}
				if v, ok := dMap["private_parameter"]; ok {
					for _, item := range v.([]interface{}) {
						PrivateParameterMap := item.(map[string]interface{})
						originRecordPrivateParameter := teo.OriginRecordPrivateParameter{}
						if v, ok := PrivateParameterMap["name"]; ok {
							originRecordPrivateParameter.Name = helper.String(v.(string))
						}
						if v, ok := PrivateParameterMap["value"]; ok {
							originRecordPrivateParameter.Value = helper.String(v.(string))
						}
						originRecord.PrivateParameter = append(originRecord.PrivateParameter, &originRecordPrivateParameter)
					}
				}

			}
		}
	}

	if d.HasChange("zone_id") {
		if v, ok := d.GetOk("zone_id"); ok {
			request.ZoneId = helper.String(v.(string))
		}
	}

	if d.HasChange("origin_type") {
		if v, ok := d.GetOk("origin_type"); ok {
			request.OriginType = helper.String(v.(string))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTeoClient().ModifyOriginGroup(request)
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

	return resourceTencentCloudTeoOriginGroupRead(d, meta)
}

func resourceTencentCloudTeoOriginGroupDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_teo_origin_group.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TeoService{client: meta.(*TencentCloudClient).apiV3Conn}
	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	zoneId := idSplit[0]
	originGroupId := idSplit[1]

	if err := service.DeleteTeoOriginGroupById(ctx, zoneId, originGroupId); err != nil {
		return err
	}

	return nil
}
