package eb

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	eb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/eb/v20210416"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudEbBus() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudEbBusRead,
		Schema: map[string]*schema.Schema{
			"order_by": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "According to which field to sort the returned results, the following fields are supported: AddTime (creation time), ModTime (modification time).",
			},

			"order": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Return results in ascending or descending order, optional values ASC (ascending) and DESC (descending).",
			},

			"filters": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "Filter conditions. The upper limit of Filters per request is 10, and the upper limit of Filter.Values 5.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"values": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Required:    true,
							Description: "One or more filter values.",
						},
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The name of the filter key.",
						},
					},
				},
			},

			"event_buses": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "event set information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"mod_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "update time.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Event set description, unlimited character type, description within 200 characters.",
						},
						"add_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "create time.",
						},
						"event_bus_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Event set name, which can only contain letters, numbers, underscores, hyphens, starts with a letter and ends with a number or letter, 2~60 characters.",
						},
						"event_bus_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "event bus Id.",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "event bus type.",
						},
						"pay_mode": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Billing mode, note: this field may return null, indicating that no valid value can be obtained.",
						},
						"connection_briefs": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Connector basic information, note: this field may return null, indicating that no valid value can be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Connector type, note: this field may return null, indicating that no valid value can be obtained.",
									},
									"status": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Connector status, note: this field may return null, indicating that no valid value can be obtained.",
									},
								},
							},
						},
						"target_briefs": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Target brief information, note: this field may return null, indicating that no valid value can be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"target_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Target ID.",
									},
									"type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Target type.",
									},
								},
							},
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

func dataSourceTencentCloudEbBusRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_eb_bus.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("order_by"); ok {
		paramMap["OrderBy"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("order"); ok {
		paramMap["Order"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("filters"); ok {
		filtersSet := v.([]interface{})
		tmpSet := make([]*eb.Filter, 0, len(filtersSet))

		for _, item := range filtersSet {
			filter := eb.Filter{}
			filterMap := item.(map[string]interface{})

			if v, ok := filterMap["values"]; ok {
				valuesSet := v.(*schema.Set).List()
				filter.Values = helper.InterfacesStringsPoint(valuesSet)
			}
			if v, ok := filterMap["name"]; ok {
				filter.Name = helper.String(v.(string))
			}
			tmpSet = append(tmpSet, &filter)
		}
		paramMap["filters"] = tmpSet
	}

	service := EbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var eventBuses []*eb.EventBus

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeEbBusByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		eventBuses = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(eventBuses))
	tmpList := make([]map[string]interface{}, 0, len(eventBuses))

	if eventBuses != nil {
		for _, eventBus := range eventBuses {
			eventBusMap := map[string]interface{}{}

			if eventBus.ModTime != nil {
				eventBusMap["mod_time"] = eventBus.ModTime
			}

			if eventBus.Description != nil {
				eventBusMap["description"] = eventBus.Description
			}

			if eventBus.AddTime != nil {
				eventBusMap["add_time"] = eventBus.AddTime
			}

			if eventBus.EventBusName != nil {
				eventBusMap["event_bus_name"] = eventBus.EventBusName
			}

			if eventBus.EventBusId != nil {
				eventBusMap["event_bus_id"] = eventBus.EventBusId
			}

			if eventBus.Type != nil {
				eventBusMap["type"] = eventBus.Type
			}

			if eventBus.PayMode != nil {
				eventBusMap["pay_mode"] = eventBus.PayMode
			}

			if eventBus.ConnectionBriefs != nil {
				connectionBriefsList := []interface{}{}
				for _, connectionBriefs := range eventBus.ConnectionBriefs {
					connectionBriefsMap := map[string]interface{}{}

					if connectionBriefs.Type != nil {
						connectionBriefsMap["type"] = connectionBriefs.Type
					}

					if connectionBriefs.Status != nil {
						connectionBriefsMap["status"] = connectionBriefs.Status
					}

					connectionBriefsList = append(connectionBriefsList, connectionBriefsMap)
				}

				eventBusMap["connection_briefs"] = []interface{}{connectionBriefsList}
			}

			if eventBus.TargetBriefs != nil {
				targetBriefsList := []interface{}{}
				for _, targetBriefs := range eventBus.TargetBriefs {
					targetBriefsMap := map[string]interface{}{}

					if targetBriefs.TargetId != nil {
						targetBriefsMap["target_id"] = targetBriefs.TargetId
					}

					if targetBriefs.Type != nil {
						targetBriefsMap["type"] = targetBriefs.Type
					}

					targetBriefsList = append(targetBriefsList, targetBriefsMap)
				}

				eventBusMap["target_briefs"] = []interface{}{targetBriefsList}
			}

			ids = append(ids, *eventBus.EventBusId)
			tmpList = append(tmpList, eventBusMap)
		}

		_ = d.Set("event_buses", tmpList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
