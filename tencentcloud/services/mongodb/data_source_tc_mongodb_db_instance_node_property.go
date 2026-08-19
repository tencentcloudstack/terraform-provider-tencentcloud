package mongodb

import (
	"fmt"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mongodb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mongodb/v20190725"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

func DataSourceTencentCloudMongodbDbInstanceNodeProperty() *schema.Resource {
	nodePropertySchema := map[string]*schema.Schema{
		"zone": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The availability zone where the node is located.",
		},
		"node_name": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Node name.",
		},
		"address": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Node access address.",
		},
		"wan_service_address": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Node public network access address (IP or domain name).",
		},
		"role": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Node role. Valid values: PRIMARY, SECONDARY, READONLY, ARBITER.",
		},
		"hidden": {
			Type:        schema.TypeBool,
			Computed:    true,
			Description: "Whether the node is a Hidden node.",
		},
		"status": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Node status. Valid values: NORMAL, STARTUP, STARTUP2, RECOVERING, DOWN, UNKNOWN, ROLLBACK, REMOVED.",
		},
		"slave_delay": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Primary-secondary sync delay in seconds.",
		},
		"priority": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Node priority. Value range: [0, 100].",
		},
		"votes": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Node votes. 1: has votes; 0: no votes.",
		},
		"tags": {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "Node tags.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"tag_key": {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "Node tag key.",
					},
					"tag_value": {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "Node tag value.",
					},
				},
			},
		},
		"replicate_set_id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Replica set ID.",
		},
	}

	return &schema.Resource{
		Read: dataSourceTencentCloudMongodbDbInstanceNodePropertyRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance ID.",
			},

			"node_ids": {
				Optional:    true,
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Node ID list.",
			},

			"roles": {
				Optional:    true,
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Node role list. Valid values: PRIMARY, SECONDARY, READONLY, ARBITER.",
			},

			"only_hidden": {
				Optional:    true,
				Type:        schema.TypeBool,
				Description: "Whether to query only Hidden nodes. Default is false.",
			},

			"priority": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Node priority. Value range: [0, 100].",
			},

			"votes": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Node votes. 1: has votes; 0: no votes.",
			},

			"tags": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "Node tags for filtering.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tag_key": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Node tag key.",
						},
						"tag_value": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Node tag value.",
						},
					},
				},
			},

			"mongos": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Mongos node property list.",
				Elem: &schema.Resource{
					Schema: nodePropertySchema,
				},
			},

			"replicate_sets": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Replica set node info list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"nodes": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Node property list in the replica set.",
							Elem: &schema.Resource{
								Schema: nodePropertySchema,
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

func nodePropertyToMap(node *mongodb.NodeProperty) map[string]interface{} {
	nodeMap := map[string]interface{}{}
	if node.Zone != nil {
		nodeMap["zone"] = node.Zone
	}
	if node.NodeName != nil {
		nodeMap["node_name"] = node.NodeName
	}
	if node.Address != nil {
		nodeMap["address"] = node.Address
	}
	if node.WanServiceAddress != nil {
		nodeMap["wan_service_address"] = node.WanServiceAddress
	}
	if node.Role != nil {
		nodeMap["role"] = node.Role
	}
	if node.Hidden != nil {
		nodeMap["hidden"] = node.Hidden
	}
	if node.Status != nil {
		nodeMap["status"] = node.Status
	}
	if node.SlaveDelay != nil {
		nodeMap["slave_delay"] = node.SlaveDelay
	}
	if node.Priority != nil {
		nodeMap["priority"] = node.Priority
	}
	if node.Votes != nil {
		nodeMap["votes"] = node.Votes
	}
	if node.Tags != nil {
		tagList := make([]map[string]interface{}, 0, len(node.Tags))
		for _, tag := range node.Tags {
			tagMap := map[string]interface{}{}
			if tag.TagKey != nil {
				tagMap["tag_key"] = tag.TagKey
			}
			if tag.TagValue != nil {
				tagMap["tag_value"] = tag.TagValue
			}
			tagList = append(tagList, tagMap)
		}
		nodeMap["tags"] = tagList
	}
	if node.ReplicateSetId != nil {
		nodeMap["replicate_set_id"] = node.ReplicateSetId
	}
	return nodeMap
}

func dataSourceTencentCloudMongodbDbInstanceNodePropertyRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_mongodb_db_instance_node_property.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	instanceId := d.Get("instance_id").(string)

	request := mongodb.NewDescribeDBInstanceNodePropertyRequest()
	request.InstanceId = helper.String(instanceId)

	if v, ok := d.GetOk("node_ids"); ok {
		nodeIds := v.([]interface{})
		request.NodeIds = make([]*string, 0, len(nodeIds))
		for _, id := range nodeIds {
			request.NodeIds = append(request.NodeIds, helper.String(id.(string)))
		}
	}

	if v, ok := d.GetOk("roles"); ok {
		roles := v.([]interface{})
		request.Roles = make([]*string, 0, len(roles))
		for _, role := range roles {
			request.Roles = append(request.Roles, helper.String(role.(string)))
		}
	}

	if v, ok := d.GetOkExists("only_hidden"); ok {
		request.OnlyHidden = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOkExists("priority"); ok {
		request.Priority = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("votes"); ok {
		request.Votes = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("tags"); ok {
		tagList := v.([]interface{})
		request.Tags = make([]*mongodb.NodeTag, 0, len(tagList))
		for _, item := range tagList {
			tagMap := item.(map[string]interface{})
			nodeTag := &mongodb.NodeTag{}
			if key, ok := tagMap["tag_key"].(string); ok && key != "" {
				nodeTag.TagKey = helper.String(key)
			}
			if val, ok := tagMap["tag_value"].(string); ok && val != "" {
				nodeTag.TagValue = helper.String(val)
			}
			request.Tags = append(request.Tags, nodeTag)
		}
	}

	client := meta.(tccommon.ProviderMeta).GetAPIV3Conn()

	var response *mongodb.DescribeDBInstanceNodePropertyResponse
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := client.UseMongodbClient().DescribeDBInstanceNodeProperty(request)
		if e != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]",
				logId, request.GetAction(), request.ToJsonString(), e.Error())
			return tccommon.RetryError(e)
		}
		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("db_instance_node_property: DescribeDBInstanceNodeProperty returned empty response"))
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]",
			logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		response = result
		return nil
	})
	if err != nil {
		return err
	}

	var mongosResult []map[string]interface{}
	if response.Response.Mongos != nil {
		mongosResult = make([]map[string]interface{}, 0, len(response.Response.Mongos))
		for _, node := range response.Response.Mongos {
			mongosResult = append(mongosResult, nodePropertyToMap(node))
		}
	} else {
		mongosResult = []map[string]interface{}{}
	}
	_ = d.Set("mongos", mongosResult)

	var replicateSetsResult []map[string]interface{}
	if response.Response.ReplicateSets != nil {
		replicateSetsResult = make([]map[string]interface{}, 0, len(response.Response.ReplicateSets))
		for _, rs := range response.Response.ReplicateSets {
			rsMap := map[string]interface{}{}
			if rs.Nodes != nil {
				nodes := make([]map[string]interface{}, 0, len(rs.Nodes))
				for _, node := range rs.Nodes {
					nodes = append(nodes, nodePropertyToMap(node))
				}
				rsMap["nodes"] = nodes
			} else {
				rsMap["nodes"] = []map[string]interface{}{}
			}
			replicateSetsResult = append(replicateSetsResult, rsMap)
		}
	} else {
		replicateSetsResult = []map[string]interface{}{}
	}
	_ = d.Set("replicate_sets", replicateSetsResult)

	d.SetId(instanceId)

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		resultMap := map[string]interface{}{
			"mongos":         mongosResult,
			"replicate_sets": replicateSetsResult,
		}
		if e := tccommon.WriteToFile(output.(string), resultMap); e != nil {
			return e
		}
	}
	return nil
}
