package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudVpcSnapshotPolicyAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudVpcSnapshotPolicyAttachmentCreate,
		Read:   resourceTencentCloudVpcSnapshotPolicyAttachmentRead,
		Delete: resourceTencentCloudVpcSnapshotPolicyAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"snapshot_policy_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Snapshot policy Id.",
			},

			"instances": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeSet,
				Description: "Associated instance information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_id": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "InstanceId.",
						},
						"instance_type": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "Instance type, currently supports set: `securitygroup`.",
						},
						"instance_region": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "The region where the instance is located.",
						},
						"instance_name": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Computed:    true,
							Description: "Instance name.",
						},
						"snapshot_policy_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Snapshot policy Id.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudVpcSnapshotPolicyAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vpc_snapshot_policy_attachment.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request          = vpc.NewAttachSnapshotInstancesRequest()
		snapshotPolicyId string
	)
	if v, ok := d.GetOk("snapshot_policy_id"); ok {
		snapshotPolicyId = v.(string)
		request.SnapshotPolicyId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("instances"); ok {
		for _, item := range v.(*schema.Set).List() {
			dMap := item.(map[string]interface{})
			snapshotInstance := vpc.SnapshotInstance{}
			if v, ok := dMap["instance_id"]; ok {
				snapshotInstance.InstanceId = helper.String(v.(string))
			}
			if v, ok := dMap["instance_type"]; ok {
				snapshotInstance.InstanceType = helper.String(v.(string))
			}
			if v, ok := dMap["instance_region"]; ok {
				snapshotInstance.InstanceRegion = helper.String(v.(string))
			}
			if v, ok := dMap["snapshot_policy_id"]; ok {
				snapshotInstance.SnapshotPolicyId = helper.String(v.(string))
			}
			if v, ok := dMap["instance_name"]; ok {
				snapshotInstance.InstanceName = helper.String(v.(string))
			}
			request.Instances = append(request.Instances, &snapshotInstance)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseVpcClient().AttachSnapshotInstances(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create vpc snapshotPolicyAttachment failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(snapshotPolicyId)

	return resourceTencentCloudVpcSnapshotPolicyAttachmentRead(d, meta)
}

func resourceTencentCloudVpcSnapshotPolicyAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vpc_snapshot_policy_attachment.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}

	snapshotPolicyId := d.Id()

	snapshotPolicyAttachment, err := service.DescribeVpcSnapshotPolicyAttachmentById(ctx, snapshotPolicyId)
	if err != nil {
		return err
	}

	if snapshotPolicyAttachment == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `VpcSnapshotPolicyAttachment` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("snapshot_policy_id", snapshotPolicyId)

	if snapshotPolicyAttachment != nil {
		instancesList := []interface{}{}
		for _, instance := range snapshotPolicyAttachment {
			instancesMap := map[string]interface{}{}

			if instance.InstanceId != nil {
				instancesMap["instance_id"] = instance.InstanceId
			}

			if instance.InstanceType != nil {
				instancesMap["instance_type"] = instance.InstanceType
			}

			if instance.InstanceRegion != nil {
				instancesMap["instance_region"] = instance.InstanceRegion
			}

			if instance.SnapshotPolicyId != nil {
				instancesMap["snapshot_policy_id"] = instance.SnapshotPolicyId
			}

			if instance.InstanceName != nil {
				instancesMap["instance_name"] = instance.InstanceName
			}

			instancesList = append(instancesList, instancesMap)
		}

		_ = d.Set("instances", instancesList)

	}

	return nil
}

func resourceTencentCloudVpcSnapshotPolicyAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vpc_snapshot_policy_attachment.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}
	snapshotPolicyId := d.Id()

	if err := service.DeleteVpcSnapshotPolicyAttachmentById(ctx, snapshotPolicyId); err != nil {
		return err
	}

	return nil
}
