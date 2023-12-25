package tcr

import (
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tcr "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tcr/v20190924"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudTcrManageReplicationOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTcrManageReplicationOperationCreate,
		Read:   resourceTencentCloudTcrManageReplicationOperationRead,
		Delete: resourceTencentCloudTcrManageReplicationOperationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"source_registry_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "copy source instance Id.",
			},

			"destination_registry_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "copy destination instance Id.",
			},

			"rule": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "synchronization rules.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "synchronization rule names.",
						},
						"dest_namespace": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "target namespace.",
						},
						"override": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "whether to cover.",
						},
						"filters": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "sync filters.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "type (name, tag, and resource).",
									},
									"value": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "empty by default.",
									},
								},
							},
						},
					},
				},
			},

			"description": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "rule description.",
			},

			"destination_region_id": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "the region ID of the target instance, such as Guangzhou is 1.",
			},

			"peer_replication_option": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "enable synchronization of configuration items across master account instances.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"peer_registry_uin": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "uin of the instance to be synchronized.",
						},
						"peer_registry_token": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "access permanent token of the instance to be synchronized.",
						},
						"enable_peer_replication": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "whether to enable cross-master account instance synchronization.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudTcrManageReplicationOperationCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tcr_manage_replication_operation.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request               = tcr.NewManageReplicationRequest()
		sourceRegistryId      string
		destinationRegistryId string
		ruleName              string
	)
	if v, ok := d.GetOk("source_registry_id"); ok {
		sourceRegistryId = v.(string)
		request.SourceRegistryId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("destination_registry_id"); ok {
		destinationRegistryId = v.(string)
		request.DestinationRegistryId = helper.String(v.(string))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "rule"); ok {
		replicationRule := tcr.ReplicationRule{}
		if v, ok := dMap["name"]; ok {
			ruleName = v.(string)
			replicationRule.Name = helper.String(v.(string))
		}
		if v, ok := dMap["dest_namespace"]; ok {
			replicationRule.DestNamespace = helper.String(v.(string))
		}
		if v, ok := dMap["override"]; ok {
			replicationRule.Override = helper.Bool(v.(bool))
		}
		if v, ok := dMap["filters"]; ok {
			for _, item := range v.([]interface{}) {
				filtersMap := item.(map[string]interface{})
				replicationFilter := tcr.ReplicationFilter{}
				if v, ok := filtersMap["type"]; ok {
					replicationFilter.Type = helper.String(v.(string))
				}
				if v, ok := filtersMap["value"]; ok {
					replicationFilter.Value = helper.String(v.(string))
				}
				replicationRule.Filters = append(replicationRule.Filters, &replicationFilter)
			}
		}
		request.Rule = &replicationRule
	}

	if v, ok := d.GetOk("description"); ok {
		request.Description = helper.String(v.(string))
	}

	if v, _ := d.GetOk("destination_region_id"); v != nil {
		request.DestinationRegionId = helper.IntUint64(v.(int))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "peer_replication_option"); ok {
		peerReplicationOption := tcr.PeerReplicationOption{}
		if v, ok := dMap["peer_registry_uin"]; ok {
			peerReplicationOption.PeerRegistryUin = helper.String(v.(string))
		}
		if v, ok := dMap["peer_registry_token"]; ok {
			peerReplicationOption.PeerRegistryToken = helper.String(v.(string))
		}
		if v, ok := dMap["enable_peer_replication"]; ok {
			peerReplicationOption.EnablePeerReplication = helper.Bool(v.(bool))
		}
		request.PeerReplicationOption = &peerReplicationOption
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTCRClient().ManageReplication(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate tcr ManageReplicationOperation failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(strings.Join([]string{sourceRegistryId, destinationRegistryId, ruleName}, tccommon.FILED_SP))

	return resourceTencentCloudTcrManageReplicationOperationRead(d, meta)
}

func resourceTencentCloudTcrManageReplicationOperationRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tcr_manage_replication_operation.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudTcrManageReplicationOperationDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tcr_manage_replication_operation.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
