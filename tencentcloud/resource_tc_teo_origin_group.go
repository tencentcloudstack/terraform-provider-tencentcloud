/*
Provides a resource to create a teo origin_group

Example Usage

```hcl
resource "tencentcloud_teo_origin_group" "origin_group" {
  zone_id = &lt;nil&gt;
  origin_group_id = &lt;nil&gt;
  origin_group_name = &lt;nil&gt;
  origin_type = &lt;nil&gt;
  configuration_type = &lt;nil&gt;
  origin_records {
		record = &lt;nil&gt;
		port = &lt;nil&gt;
		weight = &lt;nil&gt;
		area = &lt;nil&gt;
		private = &lt;nil&gt;
		private_parameter {
			name = &lt;nil&gt;
			value = &lt;nil&gt;
		}

  }
  }
```

Import

teo origin_group can be imported using the id, e.g.

```
terraform import tencentcloud_teo_origin_group.origin_group origin_group_id
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
	"strings"
)

func resourceTencentCloudTeoOriginGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTeoOriginGroupCreate,
		Read:   resourceTencentCloudTeoOriginGroupRead,
		Update: resourceTencentCloudTeoOriginGroupUpdate,
		Delete: resourceTencentCloudTeoOriginGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Site ID.",
			},

			"origin_group_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "OriginGroup ID.",
			},

			"origin_group_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "OriginGroup Name.",
			},

			"origin_type": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Type of the origin site. Valid values:- `self`: self-build website.- `cos`: tencent cos.- `third_party`: third party cos.",
			},

			"configuration_type": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Type of the origin group, this field should be set when `OriginType` is self, otherwise leave it empty. Valid values:- `area`: select an origin by using Geo info of the client IP and `Area` field in Records.- `weight`: weighted select an origin by using `Weight` field in Records.- `proto`: config by HTTP protocol.",
			},

			"origin_records": {
				Required:    true,
				Type:        schema.TypeList,
				Description: "Origin site records.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"record_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Record Id.",
						},
						"record": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Record value, which could be an IPv4/IPv6 address or a domain.",
						},
						"port": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Port of the origin site. Valid value range: 1-65535.",
						},
						"weight": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Indicating origin site&amp;#39;s weight when `Type` field is `weight`. Valid value range: 1-100. Sum of all weights should be 100.",
						},
						"area": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional:    true,
							Description: "Indicating origin site&amp;#39;s area when `Type` field is `area`. An empty List indicate the default area. Valid value:- Asia, Americas, Europe, Africa or Oceania.- 2 characters ISO 3166 area code.",
						},
						"private": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether origin site is using private authentication. Only valid when `OriginType` is `third_party`.",
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
										Description: "Parameter Name. Valid values:- AccessKeyId：Access Key ID.- SecretAccessKey：Secret Access Key.",
									},
									"value": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Parameter value.",
									},
								},
							},
						},
					},
				},
			},

			"update_time": {
				Computed:    true,
				Type:        schema.TypeString,
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
		request       = teo.NewCreateOriginGroupRequest()
		response      = teo.NewCreateOriginGroupResponse()
		zoneId        string
		originGroupId string
	)
	if v, ok := d.GetOk("zone_id"); ok {
		zoneId = v.(string)
		request.ZoneId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("origin_group_id"); ok {
		originGroupId = v.(string)
		request.OriginGroupId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("origin_group_name"); ok {
		request.OriginGroupName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("origin_type"); ok {
		request.OriginType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("configuration_type"); ok {
		request.ConfigurationType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("origin_records"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			originRecord := teo.OriginRecord{}
			if v, ok := dMap["record"]; ok {
				originRecord.Record = helper.String(v.(string))
			}
			if v, ok := dMap["port"]; ok {
				originRecord.Port = helper.IntUint64(v.(int))
			}
			if v, ok := dMap["weight"]; ok {
				originRecord.Weight = helper.IntUint64(v.(int))
			}
			if v, ok := dMap["area"]; ok {
				areaSet := v.(*schema.Set).List()
				for i := range areaSet {
					area := areaSet[i].(string)
					originRecord.Area = append(originRecord.Area, &area)
				}
			}
			if v, ok := dMap["private"]; ok {
				originRecord.Private = helper.Bool(v.(bool))
			}
			if v, ok := dMap["private_parameter"]; ok {
				for _, item := range v.([]interface{}) {
					privateParameterMap := item.(map[string]interface{})
					originRecordPrivateParameter := teo.OriginRecordPrivateParameter{}
					if v, ok := privateParameterMap["name"]; ok {
						originRecordPrivateParameter.Name = helper.String(v.(string))
					}
					if v, ok := privateParameterMap["value"]; ok {
						originRecordPrivateParameter.Value = helper.String(v.(string))
					}
					originRecord.PrivateParameter = append(originRecord.PrivateParameter, &originRecordPrivateParameter)
				}
			}
			request.OriginRecords = append(request.OriginRecords, &originRecord)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTeoClient().CreateOriginGroup(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create teo originGroup failed, reason:%+v", logId, err)
		return err
	}

	zoneId = *response.Response.ZoneId
	d.SetId(strings.Join([]string{zoneId, originGroupId}, FILED_SP))

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
	originGroupId := idSplit[1]

	originGroup, err := service.DescribeTeoOriginGroupById(ctx, zoneId, originGroupId)
	if err != nil {
		return err
	}

	if originGroup == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `TeoOriginGroup` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if originGroup.ZoneId != nil {
		_ = d.Set("zone_id", originGroup.ZoneId)
	}

	if originGroup.OriginGroupId != nil {
		_ = d.Set("origin_group_id", originGroup.OriginGroupId)
	}

	if originGroup.OriginGroupName != nil {
		_ = d.Set("origin_group_name", originGroup.OriginGroupName)
	}

	if originGroup.OriginType != nil {
		_ = d.Set("origin_type", originGroup.OriginType)
	}

	if originGroup.ConfigurationType != nil {
		_ = d.Set("configuration_type", originGroup.ConfigurationType)
	}

	if originGroup.OriginRecords != nil {
		originRecordsList := []interface{}{}
		for _, originRecords := range originGroup.OriginRecords {
			originRecordsMap := map[string]interface{}{}

			if originGroup.OriginRecords.RecordId != nil {
				originRecordsMap["record_id"] = originGroup.OriginRecords.RecordId
			}

			if originGroup.OriginRecords.Record != nil {
				originRecordsMap["record"] = originGroup.OriginRecords.Record
			}

			if originGroup.OriginRecords.Port != nil {
				originRecordsMap["port"] = originGroup.OriginRecords.Port
			}

			if originGroup.OriginRecords.Weight != nil {
				originRecordsMap["weight"] = originGroup.OriginRecords.Weight
			}

			if originGroup.OriginRecords.Area != nil {
				originRecordsMap["area"] = originGroup.OriginRecords.Area
			}

			if originGroup.OriginRecords.Private != nil {
				originRecordsMap["private"] = originGroup.OriginRecords.Private
			}

			if originGroup.OriginRecords.PrivateParameter != nil {
				privateParameterList := []interface{}{}
				for _, privateParameter := range originGroup.OriginRecords.PrivateParameter {
					privateParameterMap := map[string]interface{}{}

					if privateParameter.Name != nil {
						privateParameterMap["name"] = privateParameter.Name
					}

					if privateParameter.Value != nil {
						privateParameterMap["value"] = privateParameter.Value
					}

					privateParameterList = append(privateParameterList, privateParameterMap)
				}

				originRecordsMap["private_parameter"] = []interface{}{privateParameterList}
			}

			originRecordsList = append(originRecordsList, originRecordsMap)
		}

		_ = d.Set("origin_records", originRecordsList)

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
	originGroupId := idSplit[1]

	request.ZoneId = &zoneId
	request.OriginGroupId = &originGroupId

	immutableArgs := []string{"zone_id", "origin_group_id", "origin_group_name", "origin_type", "configuration_type", "origin_records", "update_time"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("zone_id") {
		if v, ok := d.GetOk("zone_id"); ok {
			request.ZoneId = helper.String(v.(string))
		}
	}

	if d.HasChange("origin_group_name") {
		if v, ok := d.GetOk("origin_group_name"); ok {
			request.OriginGroupName = helper.String(v.(string))
		}
	}

	if d.HasChange("origin_type") {
		if v, ok := d.GetOk("origin_type"); ok {
			request.OriginType = helper.String(v.(string))
		}
	}

	if d.HasChange("configuration_type") {
		if v, ok := d.GetOk("configuration_type"); ok {
			request.ConfigurationType = helper.String(v.(string))
		}
	}

	if d.HasChange("origin_records") {
		if v, ok := d.GetOk("origin_records"); ok {
			for _, item := range v.([]interface{}) {
				originRecord := teo.OriginRecord{}
				if v, ok := dMap["record"]; ok {
					originRecord.Record = helper.String(v.(string))
				}
				if v, ok := dMap["port"]; ok {
					originRecord.Port = helper.IntUint64(v.(int))
				}
				if v, ok := dMap["weight"]; ok {
					originRecord.Weight = helper.IntUint64(v.(int))
				}
				if v, ok := dMap["area"]; ok {
					areaSet := v.(*schema.Set).List()
					for i := range areaSet {
						area := areaSet[i].(string)
						originRecord.Area = append(originRecord.Area, &area)
					}
				}
				if v, ok := dMap["private"]; ok {
					originRecord.Private = helper.Bool(v.(bool))
				}
				if v, ok := dMap["private_parameter"]; ok {
					for _, item := range v.([]interface{}) {
						privateParameterMap := item.(map[string]interface{})
						originRecordPrivateParameter := teo.OriginRecordPrivateParameter{}
						if v, ok := privateParameterMap["name"]; ok {
							originRecordPrivateParameter.Name = helper.String(v.(string))
						}
						if v, ok := privateParameterMap["value"]; ok {
							originRecordPrivateParameter.Value = helper.String(v.(string))
						}
						originRecord.PrivateParameter = append(originRecord.PrivateParameter, &originRecordPrivateParameter)
					}
				}
				request.OriginRecords = append(request.OriginRecords, &originRecord)
			}
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTeoClient().ModifyOriginGroup(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update teo originGroup failed, reason:%+v", logId, err)
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
	originGroupId := idSplit[1]

	if err := service.DeleteTeoOriginGroupById(ctx, zoneId, originGroupId); err != nil {
		return err
	}

	return nil
}
