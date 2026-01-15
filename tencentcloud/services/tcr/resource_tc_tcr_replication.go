package tcr

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tcrv20190924 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tcr/v20190924"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudTcrReplication() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTcrReplicationCreate,
		Read:   resourceTencentCloudTcrReplicationRead,
		Delete: resourceTencentCloudTcrReplicationDelete,
		Schema: map[string]*schema.Schema{
			"source_registry_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Source instance ID.",
			},

			"destination_registry_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Destination instance ID.",
			},

			"rule": {
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				MaxItems:    1,
				Description: "Synchronization rule.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "Name of synchronization rule.",
						},
						"dest_namespace": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "Destination namespace.",
						},
						"override": {
							Type:        schema.TypeBool,
							Required:    true,
							ForceNew:    true,
							Description: "Whether to override.",
						},
						"filters": {
							Type:        schema.TypeList,
							Required:    true,
							ForceNew:    true,
							Description: "Synchronization filters.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:        schema.TypeString,
										Required:    true,
										ForceNew:    true,
										Description: "Type (`name`, `tag` and `resource`).",
									},
									"value": {
										Type:        schema.TypeString,
										Optional:    true,
										ForceNew:    true,
										Description: "It is left blank by default. If the type is `resource` it supports `image`, `chart`, and an empty string. If the type is `name` it supports Namespace name/**, Namespace name/Repository name.",
									},
								},
							},
						},
						"deletion": {
							Type:        schema.TypeBool,
							Optional:    true,
							ForceNew:    true,
							Description: "Whether synchronous deletion event.",
						},
					},
				},
			},

			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Rule description.",
			},

			"destination_region_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: "Region ID of the destination instance. For example, `1` represents Guangzhou.",
			},

			"peer_replication_option": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				MaxItems:    1,
				Description: "Configuration of the synchronization rule.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"peer_registry_uin": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "UIN of the destination instance.",
						},
						"peer_registry_token": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "Permanent access Token for the destination instance.",
						},
						"enable_peer_replication": {
							Type:        schema.TypeBool,
							Required:    true,
							ForceNew:    true,
							Description: "Whether to enable cross-account synchronization.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudTcrReplicationCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tcr_replication.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId            = tccommon.GetLogId(tccommon.ContextNil)
		ctx              = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request          = tcrv20190924.NewManageReplicationRequest()
		sourceRegistryId string
		ruleName         string
	)

	if v, ok := d.GetOk("source_registry_id"); ok {
		request.SourceRegistryId = helper.String(v.(string))
		sourceRegistryId = v.(string)
	}

	if v, ok := d.GetOk("destination_registry_id"); ok {
		request.DestinationRegistryId = helper.String(v.(string))
	}

	if ruleMap, ok := helper.InterfacesHeadMap(d, "rule"); ok {
		replicationRule := tcrv20190924.ReplicationRule{}
		if v, ok := ruleMap["name"].(string); ok && v != "" {
			replicationRule.Name = helper.String(v)
			ruleName = v
		}

		if v, ok := ruleMap["dest_namespace"].(string); ok {
			replicationRule.DestNamespace = helper.String(v)
		}

		if v, ok := ruleMap["override"].(bool); ok {
			replicationRule.Override = helper.Bool(v)
		}

		if v, ok := ruleMap["filters"]; ok {
			for _, item := range v.([]interface{}) {
				filtersMap := item.(map[string]interface{})
				replicationFilter := tcrv20190924.ReplicationFilter{}
				if v, ok := filtersMap["type"].(string); ok && v != "" {
					replicationFilter.Type = helper.String(v)
				}

				if v, ok := filtersMap["value"].(string); ok {
					replicationFilter.Value = helper.String(v)
				}

				replicationRule.Filters = append(replicationRule.Filters, &replicationFilter)
			}
		}

		if v, ok := ruleMap["deletion"].(bool); ok {
			replicationRule.Deletion = helper.Bool(v)
		}

		request.Rule = &replicationRule
	}

	if v, ok := d.GetOk("description"); ok {
		request.Description = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("destination_region_id"); ok {
		request.DestinationRegionId = helper.IntUint64(v.(int))
	}

	if peerReplicationOptionMap, ok := helper.InterfacesHeadMap(d, "peer_replication_option"); ok {
		peerReplicationOption := tcrv20190924.PeerReplicationOption{}
		if v, ok := peerReplicationOptionMap["peer_registry_uin"].(string); ok && v != "" {
			peerReplicationOption.PeerRegistryUin = helper.String(v)
		}

		if v, ok := peerReplicationOptionMap["peer_registry_token"].(string); ok && v != "" {
			peerReplicationOption.PeerRegistryToken = helper.String(v)
		}

		if v, ok := peerReplicationOptionMap["enable_peer_replication"].(bool); ok {
			peerReplicationOption.EnablePeerReplication = helper.Bool(v)
		}

		request.PeerReplicationOption = &peerReplicationOption
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTCRClient().ManageReplicationWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create tcr replication failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	d.SetId(strings.Join([]string{sourceRegistryId, ruleName}, tccommon.FILED_SP))
	return resourceTencentCloudTcrReplicationRead(d, meta)
}

func resourceTencentCloudTcrReplicationRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tcr_replication.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = TCRService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	sourceRegistryId := idSplit[0]
	ruleName := idSplit[1]

	respData, err := service.DescribeTcrReplicationById(ctx, sourceRegistryId, ruleName)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_tcr_replication` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	return nil
}

func resourceTencentCloudTcrReplicationDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tcr_replication.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = tcrv20190924.NewDeleteReplicationRuleRequest()
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	sourceRegistryId := idSplit[0]
	ruleName := idSplit[1]

	request.SourceRegistryId = &sourceRegistryId
	request.RuleName = &ruleName
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTCRClient().DeleteReplicationRuleWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete tcr replication failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return nil
}
