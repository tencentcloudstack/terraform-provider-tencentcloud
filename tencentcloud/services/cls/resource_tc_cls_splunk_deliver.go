package cls

import (
	"context"
	"fmt"
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cls "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cls/v20201016"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudClsSplunkDeliver() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudClsSplunkDeliverCreate,
		Read:   resourceTencentCloudClsSplunkDeliverRead,
		Update: resourceTencentCloudClsSplunkDeliverUpdate,
		Delete: resourceTencentCloudClsSplunkDeliverDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"topic_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Log topic ID.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Splunk deliver task name.",
			},
			"net_info": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: "Splunk deliver task target network information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"host": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Network address.",
						},
						"port": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Port.",
						},
						"token": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Authentication token.",
						},
						"net_type": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Network type. 1: internal network, 2: external network.",
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "VPC ID. Required when net_type is internal network.",
						},
						"virtual_gateway_type": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Network service type. Required when net_type is internal network. 0: CVM, 3: Direct Connect Gateway, 11: CCN, 1025: CLB.",
						},
						"is_ssl": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether to use SSL. Default is false.",
						},
					},
				},
			},
			"metadata_info": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: "Splunk deliver task metadata information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"format": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Data format. Valid values: rawlog, json.",
						},
						"meta_fields": {
							Type:        schema.TypeSet,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "Delivery fields, including __SOURCE__, __FILENAME__, __TIMESTAMP__, __HOSTNAME__, __PKG_ID__.",
						},
						"enable_tag": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether to deliver __TAG__ field.",
						},
						"tag_json_tiled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether to flatten JSON. Required when enable_tag is true.",
						},
					},
				},
			},
			"has_service_log": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Whether to enable service log. 1: disable, 2: enable. Default: enable.",
			},
			"index_ack": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Whether to enable indexer. 1: disable, 2: enable. Default: 1.",
			},
			"source": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Advanced configuration - data source. No more than 64 characters.",
			},
			"source_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Advanced configuration - data source type. No more than 64 characters.",
			},
			"index": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Advanced configuration - Splunk index. No more than 64 characters.",
			},
			"channel": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Advanced configuration - channel. Required if indexer is enabled.",
			},
			"dsl_filter": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Log pre-filtering DSL statement for raw data written to Splunk.",
			},
			"external_role": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Advanced configuration - cross-account delivery role authorization information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"role_arn": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Cross-account delivery role RoleArn.",
						},
						"external_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Cross-account delivery role name.",
						},
					},
				},
			},
			"enable": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Delivery task enable status. 0: disable, 1: enable.",
			},
			// computed
			"task_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Splunk deliver task ID.",
			},
		},
	}
}

func resourceTencentCloudClsSplunkDeliverCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cls_splunk_deliver.create")()

	var (
		logId    = tccommon.GetLogId(tccommon.ContextNil)
		request  = cls.NewCreateSplunkDeliverRequest()
		response *cls.CreateSplunkDeliverResponse
		topicId  string
	)

	if v, ok := d.GetOk("topic_id"); ok {
		topicId = v.(string)
		request.TopicId = helper.String(topicId)
	}

	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOk("net_info"); ok {
		netInfoList := v.([]interface{})
		if len(netInfoList) == 1 {
			netInfoMap := netInfoList[0].(map[string]interface{})
			netInfo := cls.NetInfo{}
			if v, ok := netInfoMap["host"]; ok {
				netInfo.Host = helper.String(v.(string))
			}
			if v, ok := netInfoMap["port"]; ok {
				netInfo.Port = helper.IntUint64(v.(int))
			}
			if v, ok := netInfoMap["token"]; ok {
				netInfo.Token = helper.String(v.(string))
			}
			if v, ok := netInfoMap["net_type"]; ok {
				netInfo.NetType = helper.IntUint64(v.(int))
			}
			if v, ok := netInfoMap["vpc_id"]; ok {
				netInfo.VpcId = helper.String(v.(string))
			}
			if v, ok := netInfoMap["virtual_gateway_type"]; ok {
				netInfo.VirtualGatewayType = helper.IntUint64(v.(int))
			}
			if v, ok := netInfoMap["is_ssl"]; ok {
				netInfo.IsSSL = helper.Bool(v.(bool))
			}
			request.NetInfo = &netInfo
		}
	}

	if v, ok := d.GetOk("metadata_info"); ok {
		metadataInfoList := v.([]interface{})
		if len(metadataInfoList) == 1 {
			metadataInfoMap := metadataInfoList[0].(map[string]interface{})
			metadataInfo := cls.MetadataInfo{}
			if v, ok := metadataInfoMap["format"]; ok {
				metadataInfo.Format = helper.String(v.(string))
			}
			if v, ok := metadataInfoMap["meta_fields"]; ok {
				metaFields := v.(*schema.Set).List()
				for _, mf := range metaFields {
					metadataInfo.MetaFields = append(metadataInfo.MetaFields, helper.String(mf.(string)))
				}
			}
			if v, ok := metadataInfoMap["enable_tag"]; ok {
				metadataInfo.EnableTag = helper.Bool(v.(bool))
			}
			if v, ok := metadataInfoMap["tag_json_tiled"]; ok {
				metadataInfo.TagJsonTiled = helper.Bool(v.(bool))
			}
			request.MetadataInfo = &metadataInfo
		}
	}

	if v, ok := d.GetOkExists("has_service_log"); ok {
		request.HasServiceLog = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("index_ack"); ok {
		request.IndexAck = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("source"); ok {
		request.Source = helper.String(v.(string))
	}

	if v, ok := d.GetOk("source_type"); ok {
		request.SourceType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("index"); ok {
		request.Index = helper.String(v.(string))
	}

	if v, ok := d.GetOk("channel"); ok {
		request.Channel = helper.String(v.(string))
	}

	if v, ok := d.GetOk("dsl_filter"); ok {
		request.DSLFilter = helper.String(v.(string))
	}

	if v, ok := d.GetOk("external_role"); ok {
		externalRoleList := v.([]interface{})
		if len(externalRoleList) == 1 {
			externalRoleMap := externalRoleList[0].(map[string]interface{})
			externalRole := cls.ExternalRole{}
			if v, ok := externalRoleMap["role_arn"]; ok {
				externalRole.RoleArn = helper.String(v.(string))
			}
			if v, ok := externalRoleMap["external_id"]; ok {
				externalRole.ExternalId = helper.String(v.(string))
			}
			request.ExternalRole = &externalRole
		}
	}

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClsClient().CreateSplunkDeliver(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create cls splunk_deliver failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create cls splunk_deliver failed, reason:%+v", logId, err)
		return err
	}

	log.Printf("[CRUD] cls splunk_deliver create, logId=%s, topicId=%s", logId, topicId)

	if response.Response.TaskId == nil || *response.Response.TaskId == "" {
		return fmt.Errorf("splunk_deliver TaskId is nil or empty.")
	}

	taskId := *response.Response.TaskId
	d.SetId(strings.Join([]string{topicId, taskId}, tccommon.FILED_SP))
	return resourceTencentCloudClsSplunkDeliverRead(d, meta)
}

func resourceTencentCloudClsSplunkDeliverRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cls_splunk_deliver.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	id := d.Id()
	idSplit := strings.Split(id, tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("splunk_deliver id is invalid, id format should be topic_id#task_id, got: %s", id)
	}

	topicId := idSplit[0]
	taskId := idSplit[1]

	request := cls.NewDescribeSplunkDeliversRequest()
	request.TopicId = helper.String(topicId)
	request.Filters = []*cls.Filter{
		{
			Key:    helper.String("taskId"),
			Values: []*string{helper.String(taskId)},
		},
	}
	request.Limit = helper.IntUint64(100)

	var response *cls.DescribeSplunkDeliversResponse
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClsClient().DescribeSplunkDeliversWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe cls splunk_deliver failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s read cls splunk_deliver failed, reason:%+v", logId, err)
		return err
	}

	if response.Response.Infos == nil || len(response.Response.Infos) == 0 {
		log.Printf("[CRUD] cls splunk_deliver id=%s, read empty, skip SetId", d.Id())
		d.SetId("")
		return nil
	}

	info := response.Response.Infos[0]

	_ = d.Set("topic_id", info.TopicId)
	_ = d.Set("task_id", info.TaskId)

	if info.Name != nil {
		_ = d.Set("name", info.Name)
	}

	if info.Enable != nil {
		_ = d.Set("enable", info.Enable)
	}

	if info.NetInfo != nil {
		netInfoMap := map[string]interface{}{
			"host":  info.NetInfo.Host,
			"token": info.NetInfo.Token,
		}
		if info.NetInfo.Port != nil {
			netInfoMap["port"] = info.NetInfo.Port
		}
		if info.NetInfo.NetType != nil {
			netInfoMap["net_type"] = info.NetInfo.NetType
		}
		if info.NetInfo.VpcId != nil {
			netInfoMap["vpc_id"] = info.NetInfo.VpcId
		}
		if info.NetInfo.VirtualGatewayType != nil {
			netInfoMap["virtual_gateway_type"] = info.NetInfo.VirtualGatewayType
		}
		if info.NetInfo.IsSSL != nil {
			netInfoMap["is_ssl"] = info.NetInfo.IsSSL
		}
		_ = d.Set("net_info", []interface{}{netInfoMap})
	}

	if info.Metadata != nil {
		metadataInfoMap := map[string]interface{}{
			"format":      info.Metadata.Format,
			"meta_fields": info.Metadata.MetaFields,
		}
		if info.Metadata.EnableTag != nil {
			metadataInfoMap["enable_tag"] = info.Metadata.EnableTag
		}
		if info.Metadata.TagJsonTiled != nil {
			metadataInfoMap["tag_json_tiled"] = info.Metadata.TagJsonTiled
		}
		_ = d.Set("metadata_info", []interface{}{metadataInfoMap})
	}

	if info.HasServiceLog != nil {
		_ = d.Set("has_service_log", info.HasServiceLog)
	}

	if info.IndexAck != nil {
		_ = d.Set("index_ack", info.IndexAck)
	}

	if info.Source != nil {
		_ = d.Set("source", info.Source)
	}

	if info.SourceType != nil {
		_ = d.Set("source_type", info.SourceType)
	}

	if info.Index != nil {
		_ = d.Set("index", info.Index)
	}

	if info.Channel != nil {
		_ = d.Set("channel", info.Channel)
	}

	if info.DSLFilter != nil {
		_ = d.Set("dsl_filter", info.DSLFilter)
	}

	if info.ExternalRole != nil {
		externalRoleMap := map[string]interface{}{
			"role_arn":    info.ExternalRole.RoleArn,
			"external_id": info.ExternalRole.ExternalId,
		}
		_ = d.Set("external_role", []interface{}{externalRoleMap})
	}

	return nil
}

func resourceTencentCloudClsSplunkDeliverUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cls_splunk_deliver.update")()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	id := d.Id()
	idSplit := strings.Split(id, tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("splunk_deliver id is invalid, id format should be topic_id#task_id, got: %s", id)
	}

	topicId := idSplit[0]
	taskId := idSplit[1]

	request := cls.NewModifySplunkDeliverRequest()
	request.TaskId = helper.String(taskId)
	request.TopicId = helper.String(topicId)

	if d.HasChange("name") {
		if v, ok := d.GetOk("name"); ok {
			request.Name = helper.String(v.(string))
		}
	}

	if d.HasChange("enable") {
		if v, ok := d.GetOkExists("enable"); ok {
			request.Enable = helper.IntInt64(v.(int))
		}
	}

	if d.HasChange("net_info") {
		if v, ok := d.GetOk("net_info"); ok {
			netInfoList := v.([]interface{})
			if len(netInfoList) == 1 {
				netInfoMap := netInfoList[0].(map[string]interface{})
				netInfo := cls.NetInfo{}
				if v, ok := netInfoMap["host"]; ok {
					netInfo.Host = helper.String(v.(string))
				}
				if v, ok := netInfoMap["port"]; ok {
					netInfo.Port = helper.IntUint64(v.(int))
				}
				if v, ok := netInfoMap["token"]; ok {
					netInfo.Token = helper.String(v.(string))
				}
				if v, ok := netInfoMap["net_type"]; ok {
					netInfo.NetType = helper.IntUint64(v.(int))
				}
				if v, ok := netInfoMap["vpc_id"]; ok {
					netInfo.VpcId = helper.String(v.(string))
				}
				if v, ok := netInfoMap["virtual_gateway_type"]; ok {
					netInfo.VirtualGatewayType = helper.IntUint64(v.(int))
				}
				if v, ok := netInfoMap["is_ssl"]; ok {
					netInfo.IsSSL = helper.Bool(v.(bool))
				}
				request.NetInfo = &netInfo
			}
		}
	}

	if d.HasChange("metadata_info") {
		if v, ok := d.GetOk("metadata_info"); ok {
			metadataInfoList := v.([]interface{})
			if len(metadataInfoList) == 1 {
				metadataInfoMap := metadataInfoList[0].(map[string]interface{})
				metadataInfo := cls.MetadataInfo{}
				if v, ok := metadataInfoMap["format"]; ok {
					metadataInfo.Format = helper.String(v.(string))
				}
				if v, ok := metadataInfoMap["meta_fields"]; ok {
					metaFields := v.(*schema.Set).List()
					for _, mf := range metaFields {
						metadataInfo.MetaFields = append(metadataInfo.MetaFields, helper.String(mf.(string)))
					}
				}
				if v, ok := metadataInfoMap["enable_tag"]; ok {
					metadataInfo.EnableTag = helper.Bool(v.(bool))
				}
				if v, ok := metadataInfoMap["tag_json_tiled"]; ok {
					metadataInfo.TagJsonTiled = helper.Bool(v.(bool))
				}
				request.MetadataInfo = &metadataInfo
			}
		}
	}

	if d.HasChange("has_service_log") {
		if v, ok := d.GetOkExists("has_service_log"); ok {
			request.HasServiceLog = helper.IntInt64(v.(int))
		}
	}

	if d.HasChange("index_ack") {
		if v, ok := d.GetOkExists("index_ack"); ok {
			request.IndexAck = helper.IntInt64(v.(int))
		}
	}

	if d.HasChange("source") {
		if v, ok := d.GetOk("source"); ok {
			request.Source = helper.String(v.(string))
		}
	}

	if d.HasChange("source_type") {
		if v, ok := d.GetOk("source_type"); ok {
			request.SourceType = helper.String(v.(string))
		}
	}

	if d.HasChange("index") {
		if v, ok := d.GetOk("index"); ok {
			request.Index = helper.String(v.(string))
		}
	}

	if d.HasChange("channel") {
		if v, ok := d.GetOk("channel"); ok {
			request.Channel = helper.String(v.(string))
		}
	}

	if d.HasChange("dsl_filter") {
		if v, ok := d.GetOk("dsl_filter"); ok {
			request.DSLFilter = helper.String(v.(string))
		}
	}

	if d.HasChange("external_role") {
		if v, ok := d.GetOk("external_role"); ok {
			externalRoleList := v.([]interface{})
			if len(externalRoleList) == 1 {
				externalRoleMap := externalRoleList[0].(map[string]interface{})
				externalRole := cls.ExternalRole{}
				if v, ok := externalRoleMap["role_arn"]; ok {
					externalRole.RoleArn = helper.String(v.(string))
				}
				if v, ok := externalRoleMap["external_id"]; ok {
					externalRole.ExternalId = helper.String(v.(string))
				}
				request.ExternalRole = &externalRole
			}
		}
	}

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClsClient().ModifySplunkDeliver(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s update cls splunk_deliver failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudClsSplunkDeliverRead(d, meta)
}

func resourceTencentCloudClsSplunkDeliverDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cls_splunk_deliver.delete")()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	id := d.Id()
	idSplit := strings.Split(id, tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("splunk_deliver id is invalid, id format should be topic_id#task_id, got: %s", id)
	}

	topicId := idSplit[1]
	taskId := idSplit[0]

	request := cls.NewDeleteSplunkDeliverRequest()
	request.TaskId = helper.String(taskId)
	request.TopicId = helper.String(topicId)

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClsClient().DeleteSplunkDeliver(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s delete cls splunk_deliver failed, reason:%+v", logId, err)
		return err
	}

	return nil
}
