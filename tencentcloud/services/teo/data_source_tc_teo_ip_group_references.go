package teo

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	teo "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

func DataSourceTencentCloudTeoIPGroupReferences() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTeoIPGroupReferencesRead,
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Site ID.",
			},

			"group_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "IP group ID.",
			},

			"references": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of references to the IP group.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"zone_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Site ID.",
						},
						"entity_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Entity type.",
						},
						"entity_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Entity ID.",
						},
						"entity_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Entity name.",
						},
						"sub_entity_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Sub entity type.",
						},
						"sub_entity_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Sub entity ID.",
						},
						"sub_entity_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Sub entity name.",
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

func dataSourceTencentCloudTeoIPGroupReferencesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_teo_ip_group_references.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	client := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoClient()

	request := teo.NewDescribeIPGroupReferencesRequest()

	if v, ok := d.GetOk("zone_id"); ok {
		request.ZoneId = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("group_id"); ok {
		request.GroupId = helper.IntInt64(v.(int))
	}

	var references []*teo.IPGroupReference
	var offset int64 = 0
	var limit int64 = 200

	for {
		request.Offset = &offset
		request.Limit = &limit
		ratelimit.Check(request.GetAction())

		var response *teo.DescribeIPGroupReferencesResponse
		err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			resp, e := client.DescribeIPGroupReferences(request)
			if e != nil {
				return tccommon.RetryError(e)
			}

			if resp == nil || resp.Response == nil || len(resp.Response.References) == 0 {
				return resource.NonRetryableError(fmt.Errorf("teo_ip_group_references DescribeIPGroupReferences response is empty, zone_id=%s, group_id=%d", d.Get("zone_id"), d.Get("group_id")))
			}

			response = resp
			return nil
		})
		if err != nil {
			log.Printf("[DATASOURCE] read empty, skip SetId")
			return err
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response.Response.References != nil {
			references = append(references, response.Response.References...)
		}

		if len(response.Response.References) < int(limit) {
			break
		}
		offset += limit
	}

	referencesList := make([]map[string]interface{}, 0, len(references))
	if references != nil {
		for _, reference := range references {
			referenceMap := map[string]interface{}{}
			if reference.ZoneId != nil {
				referenceMap["zone_id"] = reference.ZoneId
			}

			if reference.EntityType != nil {
				referenceMap["entity_type"] = reference.EntityType
			}

			if reference.EntityId != nil {
				referenceMap["entity_id"] = reference.EntityId
			}

			if reference.EntityName != nil {
				referenceMap["entity_name"] = reference.EntityName
			}

			if reference.SubEntityType != nil {
				referenceMap["sub_entity_type"] = reference.SubEntityType
			}

			if reference.SubEntityId != nil {
				referenceMap["sub_entity_id"] = reference.SubEntityId
			}

			if reference.SubEntityName != nil {
				referenceMap["sub_entity_name"] = reference.SubEntityName
			}

			referencesList = append(referencesList, referenceMap)
		}
	}

	_ = d.Set("references", referencesList)

	d.SetId(helper.BuildToken())

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), d); e != nil {
			return e
		}
	}

	return nil
}
