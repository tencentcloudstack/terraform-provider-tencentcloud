package es

import (
	"context"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	elasticsearch "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/es/v20180416"
)

func ResourceTencentCloudElasticsearchSecurityGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudElasticsearchSecurityGroupCreate,
		Read:   resourceTencentCloudElasticsearchSecurityGroupRead,
		Update: resourceTencentCloudElasticsearchSecurityGroupUpdate,
		Delete: resourceTencentCloudElasticsearchSecurityGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Instance Id.",
			},

			"security_group_ids": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Security group id list.",
			},
		},
	}
}

func resourceTencentCloudElasticsearchSecurityGroupCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_elasticsearch_security_group.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var instanceId string
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	d.SetId(instanceId)

	return resourceTencentCloudElasticsearchSecurityGroupUpdate(d, meta)
}

func resourceTencentCloudElasticsearchSecurityGroupRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_elasticsearch_security_group.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := ElasticsearchService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	instanceId := d.Id()

	securityGroup, err := service.DescribeInstanceById(ctx, instanceId)
	if err != nil {
		return err
	}

	if securityGroup == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `ElasticsearchSecurityGroup` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if securityGroup.InstanceId != nil {
		_ = d.Set("instance_id", securityGroup.InstanceId)
	}

	if securityGroup.SecurityGroups != nil {
		_ = d.Set("security_group_ids", securityGroup.SecurityGroups)
	}

	return nil
}

func resourceTencentCloudElasticsearchSecurityGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_elasticsearch_security_group.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	request := elasticsearch.NewModifyEsVipSecurityGroupRequest()

	instanceId := d.Id()
	request.InstanceId = &instanceId

	if v, ok := d.GetOk("security_group_ids"); ok {
		securityGroupIdsSet := v.(*schema.Set).List()
		for i := range securityGroupIdsSet {
			securityGroupIds := securityGroupIdsSet[i].(string)
			request.SecurityGroupIds = append(request.SecurityGroupIds, &securityGroupIds)
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseEsClient().ModifyEsVipSecurityGroup(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update elasticsearch securityGroup failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudElasticsearchSecurityGroupRead(d, meta)
}

func resourceTencentCloudElasticsearchSecurityGroupDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_elasticsearch_security_group.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
