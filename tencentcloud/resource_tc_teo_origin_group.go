/*
Provides a resource to create a teo originGroup

Example Usage

```hcl
locals {
  group0 = [
    {
      "record" = "1.1.1.1"
      "port"   = 80
      "weight" = 30
    }, {
      "record" = "2.2.2.2"
      "port"   = 443
      "weight" = 70
    }
  ]
}

resource "tencentcloud_teo_origin_group" "origin_group" {
  zone_id     = tencentcloud_teo_zone.zone.id
  origin_name = "group0"
  origin_type = "self"
  type        = "weight"

  dynamic "record" {
    for_each = local.group0
    content {
      record = record.value["record"]
      port   = record.value["port"]
      weight = record.value["weight"]
      area   = []
    }
  }
}

```
Import

teo origin_group can be imported using the id, e.g.
```
$ terraform import tencentcloud_teo_origin_group.origin_group zoneId#originId
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
			"origin_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "OriginGroup Name.",
			},

			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Type of the origin group, this field is required only when `OriginType` is `self`. Valid values:- area: select an origin by using Geo info of the client IP and `Area` field in Records.- weight: weighted select an origin by using `Weight` field in Records.",
			},

			"record": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Origin website records.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"record": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Record Value.",
						},
						"area": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Required:    true,
							Description: "Indicating origin website&#39;s area when `Type` field is `area`. An empty List indicate the default area.",
						},
						"weight": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Indicating origin website&#39;s weight when `Type` field is `weight`. Valid value range: 1-100. Sum of all weights should be 100.",
						},
						"port": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Port of the origin website.",
						},
						"private": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether origin website is using private authentication. Only valid when `OriginType` is `third_party`.",
						},
						"private_parameter": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Parameters for private authentication. Only valid when `Private` is `true`.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Parameter Name. Valid values:- AccessKeyId: Access Key ID.- SecretAccessKey: Secret Access Key.",
									},
									"value": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Parameter value.",
									},
								},
							},
						},
						"record_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Record Id.",
						},
					},
				},
			},

			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Site ID.",
			},

			"origin_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Type of the origin website. Valid values:- self: self-build website.- cos: tencent cos.- third_party: third party cos.",
			},

			"zone_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Site Name.",
			},

			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Last modification date.",
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
		originId string
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

			request.Record = append(request.Record, &originRecord)
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

	originId = *response.Response.OriginId

	d.SetId(zoneId + FILED_SP + originId)
	return resourceTencentCloudTeoOriginGroupRead(d, meta)
}

func resourceTencentCloudTeoOriginGroupRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_teo_origin_group.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TeoService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	zoneId := idSplit[0]
	originId := idSplit[1]

	originGroup, err := service.DescribeTeoOriginGroup(ctx, zoneId, originId)

	if err != nil {
		return err
	}

	if originGroup == nil {
		d.SetId("")
		return fmt.Errorf("resource `originGroup` %s does not exist", originId)
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
			if record.RecordId != nil {
				recordMap["record_id"] = record.RecordId
			}

			recordList = append(recordList, recordMap)
		}
		_ = d.Set("record", recordList)
	}

	if originGroup.ZoneId != nil {
		_ = d.Set("zone_id", originGroup.ZoneId)
	}

	if originGroup.OriginType != nil {
		_ = d.Set("origin_type", originGroup.OriginType)
	}

	if originGroup.ZoneName != nil {
		_ = d.Set("zone_name", originGroup.ZoneName)
	}

	if originGroup.UpdateTime != nil {
		_ = d.Set("update_time", originGroup.UpdateTime)
	}

	return nil
}

func resourceTencentCloudTeoOriginGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_teo_origin_group.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := teo.NewModifyOriginGroupRequest()

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	zoneId := idSplit[0]
	originId := idSplit[1]

	request.ZoneId = &zoneId
	request.OriginId = &originId

	if v, ok := d.GetOk("origin_name"); ok {
		request.OriginName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("type"); ok {
		request.Type = helper.String(v.(string))
	}

	if v, ok := d.GetOk("origin_type"); ok {
		request.OriginType = helper.String(v.(string))
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

				request.Record = append(request.Record, &originRecord)
			}
		}
	}

	if d.HasChange("zone_id") {
		if v, ok := d.GetOk("zone_id"); ok {
			request.ZoneId = helper.String(v.(string))
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
	originId := idSplit[1]

	if err := service.DeleteTeoOriginGroupById(ctx, zoneId, originId); err != nil {
		return err
	}

	return nil
}
