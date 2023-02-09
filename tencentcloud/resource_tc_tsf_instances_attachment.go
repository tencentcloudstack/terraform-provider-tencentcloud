/*
Provides a resource to create a tsf instances_attachment

Example Usage

```hcl
resource "tencentcloud_tsf_instances_attachment" "instances_attachment" {
  cluster_id = ""
  instance_id = ""
  os_name = ""
  image_id = ""
  password = ""
  key_id = ""
  sg_id = ""
  instance_import_mode = ""
  tags = {
    "createdBy" = "terraform"
  }
}
```

Import

tsf instances_attachment can be imported using the id, e.g.

```
terraform import tencentcloud_tsf_instances_attachment.instances_attachment instances_attachment_id
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
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	tsf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tsf/v20180326"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudTsfInstancesAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTsfInstancesAttachmentCreate,
		Read:   resourceTencentCloudTsfInstancesAttachmentRead,
		Delete: resourceTencentCloudTsfInstancesAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Cluster id.",
			},

			"instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Instance id.",
			},

			"os_name": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "operating system name.",
			},

			"image_id": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "OS image ID.",
			},

			"password": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Reinstall system password settings.",
			},

			"key_id": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Reinstall the system, associate key settings.",
			},

			"sg_id": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Security Group Settings.",
			},

			"instance_import_mode": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Cloud host import method, virtual machine clusters are required, container clusters do not fill in this field, `R`: reinstall the TSF system image, `M`: manually install the agent.",
			},

			"instance_name": {
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Instance name.",
			},

			"lan_ip": {
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Machine intranet address IP.",
			},

			"wan_ip": {
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Machine external network address IP.",
			},

			"instance_desc": {
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Instance desc.",
			},

			"cluster_name": {
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Cluster name.",
			},

			"instance_status": {
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "VM status Virtual machine: the status of the virtual machine, container: the status of the virtual machine where the Pod resides.",
			},

			"instance_available_status": {
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Available state of VM Virtual machine: whether the virtual machine can be used as a resource, container: whether the virtual machine can be used as a resource to deploy POD.",
			},

			"service_instance_status": {
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "The status of the service instance under the service Virtual machine: Whether the application is available + Agent status, Container: Pod status.",
			},

			"count_in_tsf": {
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "Identifies whether this instance has been added in tsf.",
			},

			"group_id": {
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "The ID of the deployment group to which the machine belongs.",
			},

			"application_id": {
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "The ID of the deployment group to which the machine belongs.",
			},

			"application_name": {
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "The name of the application to which the machine belongs.",
			},

			"instance_created_time": {
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "The creation time of the machine instance in CVM.",
			},

			"instance_expired_time": {
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "The expiration time of the machine instance in CVM.",
			},

			"instance_charge_type": {
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "The billing mode of the machine instance in CVM.",
			},

			"instance_total_cpu": {
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeFloat,
				Description: "Total CPU information of the machine instance.",
			},

			"instance_total_mem": {
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeFloat,
				Description: "Total memory information of the machine instance.",
			},

			"instance_used_cpu": {
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeFloat,
				Description: "CPU information used by the machine instance.",
			},

			"instance_used_mem": {
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeFloat,
				Description: "Memory information used by the machine instance.",
			},

			"instance_limit_cpu": {
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeFloat,
				Description: "Machine instance Limit CPU information.",
			},

			"instance_limit_mem": {
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeFloat,
				Description: "Machine instance Limit memory information.",
			},

			"instance_pkg_version": {
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "package version.",
			},

			"cluster_type": {
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Cluster type.",
			},

			"restrict_state": {
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Machine instance business status.",
			},

			"update_time": {
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Update time.",
			},

			"operation_state": {
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "instance execution state.",
			},

			"namespace_id": {
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "NamespaceId Ns ID.",
			},

			"instance_zone_id": {
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "InstanceZoneId availability zone ID.",
			},

			"application_type": {
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "application type.",
			},

			"application_resource_type": {
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Resource Type.",
			},

			"service_sidecar_status": {
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "sidecar status.",
			},

			"group_name": {
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Group name.",
			},

			"namespace_name": {
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Namespace name.",
			},

			"reason": {
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "health check reason.",
			},

			"agent_version": {
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "agent version.",
			},

			"node_instance_id": {
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Container host instance ID.",
			},
		},
	}
}

func resourceTencentCloudTsfInstancesAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tsf_instances_attachment.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = tsf.NewAddInstancesRequest()
		response   = tsf.NewAddInstancesResponse()
		clusterId  string
		instanceId string
	)
	if v, ok := d.GetOk("cluster_id"); ok {
		clusterId = v.(string)
		request.ClusterId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceIdList = []*string{helper.String(v.(string))}
	}

	if v, ok := d.GetOk("os_name"); ok {
		request.OsName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("image_id"); ok {
		request.ImageId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("password"); ok {
		request.Password = helper.String(v.(string))
	}

	if v, ok := d.GetOk("key_id"); ok {
		request.KeyId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("sg_id"); ok {
		request.SgId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("instance_import_mode"); ok {
		request.InstanceImportMode = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTsfClient().AddInstances(request)
		if e != nil {
			ee, ok := e.(*sdkErrors.TencentCloudSDKError)
			if !ok {
				return resource.NonRetryableError(e)
			}
			if ee.Code == "ResourceInUse.InstanceHasBeenUsed" {
				return resource.NonRetryableError(fmt.Errorf("The machine instance [%v] is already in use", instanceId))
			}
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create tsf instancesAttachment failed, reason:%+v", logId, err)
		return err
	}

	if *response.Response.Result {
		d.SetId(clusterId + FILED_SP + instanceId)
	} else {
		return fmt.Errorf("Failed to import cloud host %v", instanceId)
	}

	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		tagService := TagService{client: meta.(*TencentCloudClient).apiV3Conn}
		region := meta.(*TencentCloudClient).apiV3Conn.Region
		resourceName := fmt.Sprintf("qcs::tsf:%s:uin/:instance/%s", region, d.Id())
		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
	}

	return resourceTencentCloudTsfInstancesAttachmentRead(d, meta)
}

func resourceTencentCloudTsfInstancesAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tsf_instances_attachment.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TsfService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	clusterId := idSplit[0]
	instanceId := idSplit[1]

	instancesAttachment, err := service.DescribeTsfInstancesAttachmentById(ctx, clusterId, instanceId)
	if err != nil {
		return err
	}

	if instancesAttachment == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `TsfInstancesAttachment` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("cluster_id", clusterId)
	_ = d.Set("instance_id", instanceId)

	// if instancesAttachment.OsName != nil {
	// 	_ = d.Set("os_name", instancesAttachment.OsName)
	// }

	// if instancesAttachment.ImageId != nil {
	// 	_ = d.Set("image_id", instancesAttachment.ImageId)
	// }

	// if instancesAttachment.Password != nil {
	// 	_ = d.Set("password", instancesAttachment.Password)
	// }

	// if instancesAttachment.KeyId != nil {
	// 	_ = d.Set("key_id", instancesAttachment.KeyId)
	// }

	// if instancesAttachment.SgId != nil {
	// 	_ = d.Set("sg_id", instancesAttachment.SgId)
	// }

	if instancesAttachment.InstanceImportMode != nil {
		_ = d.Set("instance_import_mode", instancesAttachment.InstanceImportMode)
	}

	if instancesAttachment.InstanceName != nil {
		_ = d.Set("instance_name", instancesAttachment.InstanceName)
	}

	if instancesAttachment.LanIp != nil {
		_ = d.Set("lan_ip", instancesAttachment.LanIp)
	}

	if instancesAttachment.WanIp != nil {
		_ = d.Set("wan_ip", instancesAttachment.WanIp)
	}

	if instancesAttachment.InstanceDesc != nil {
		_ = d.Set("instance_desc", instancesAttachment.InstanceDesc)
	}

	if instancesAttachment.ClusterName != nil {
		_ = d.Set("cluster_name", instancesAttachment.ClusterName)
	}

	if instancesAttachment.InstanceStatus != nil {
		_ = d.Set("instance_status", instancesAttachment.InstanceStatus)
	}

	if instancesAttachment.InstanceAvailableStatus != nil {
		_ = d.Set("instance_available_status", instancesAttachment.InstanceAvailableStatus)
	}

	if instancesAttachment.ServiceInstanceStatus != nil {
		_ = d.Set("service_instance_status", instancesAttachment.ServiceInstanceStatus)
	}

	if instancesAttachment.CountInTsf != nil {
		_ = d.Set("count_in_tsf", instancesAttachment.CountInTsf)
	}

	if instancesAttachment.GroupId != nil {
		_ = d.Set("group_id", instancesAttachment.GroupId)
	}

	if instancesAttachment.ApplicationId != nil {
		_ = d.Set("application_id", instancesAttachment.ApplicationId)
	}

	if instancesAttachment.ApplicationName != nil {
		_ = d.Set("application_name", instancesAttachment.ApplicationName)
	}

	if instancesAttachment.InstanceCreatedTime != nil {
		_ = d.Set("instance_created_time", instancesAttachment.InstanceCreatedTime)
	}

	if instancesAttachment.InstanceExpiredTime != nil {
		_ = d.Set("instance_expired_time", instancesAttachment.InstanceExpiredTime)
	}

	if instancesAttachment.InstanceChargeType != nil {
		_ = d.Set("instance_charge_type", instancesAttachment.InstanceChargeType)
	}

	if instancesAttachment.InstanceTotalCpu != nil {
		_ = d.Set("instance_total_cpu", instancesAttachment.InstanceTotalCpu)
	}

	if instancesAttachment.InstanceTotalMem != nil {
		_ = d.Set("instance_total_mem", instancesAttachment.InstanceTotalMem)
	}

	if instancesAttachment.InstanceUsedCpu != nil {
		_ = d.Set("instance_used_cpu", instancesAttachment.InstanceUsedCpu)
	}

	if instancesAttachment.InstanceUsedMem != nil {
		_ = d.Set("instance_used_mem", instancesAttachment.InstanceUsedMem)
	}

	if instancesAttachment.InstanceLimitCpu != nil {
		_ = d.Set("instance_limit_cpu", instancesAttachment.InstanceLimitCpu)
	}

	if instancesAttachment.InstanceLimitMem != nil {
		_ = d.Set("instance_limit_mem", instancesAttachment.InstanceLimitMem)
	}

	if instancesAttachment.InstancePkgVersion != nil {
		_ = d.Set("instance_pkg_version", instancesAttachment.InstancePkgVersion)
	}

	if instancesAttachment.ClusterType != nil {
		_ = d.Set("cluster_type", instancesAttachment.ClusterType)
	}

	if instancesAttachment.RestrictState != nil {
		_ = d.Set("restrict_state", instancesAttachment.RestrictState)
	}

	if instancesAttachment.UpdateTime != nil {
		_ = d.Set("update_time", instancesAttachment.UpdateTime)
	}

	if instancesAttachment.OperationState != nil {
		_ = d.Set("operation_state", instancesAttachment.OperationState)
	}

	if instancesAttachment.NamespaceId != nil {
		_ = d.Set("namespace_id", instancesAttachment.NamespaceId)
	}

	if instancesAttachment.InstanceZoneId != nil {
		_ = d.Set("instance_zone_id", instancesAttachment.InstanceZoneId)
	}

	if instancesAttachment.ApplicationType != nil {
		_ = d.Set("application_type", instancesAttachment.ApplicationType)
	}

	if instancesAttachment.ApplicationResourceType != nil {
		_ = d.Set("application_resource_type", instancesAttachment.ApplicationResourceType)
	}

	if instancesAttachment.ServiceSidecarStatus != nil {
		_ = d.Set("service_sidecar_status", instancesAttachment.ServiceSidecarStatus)
	}

	if instancesAttachment.GroupName != nil {
		_ = d.Set("group_name", instancesAttachment.GroupName)
	}

	if instancesAttachment.NamespaceName != nil {
		_ = d.Set("namespace_name", instancesAttachment.NamespaceName)
	}

	if instancesAttachment.Reason != nil {
		_ = d.Set("reason", instancesAttachment.Reason)
	}

	if instancesAttachment.AgentVersion != nil {
		_ = d.Set("agent_version", instancesAttachment.AgentVersion)
	}

	if instancesAttachment.NodeInstanceId != nil {
		_ = d.Set("node_instance_id", instancesAttachment.NodeInstanceId)
	}

	tcClient := meta.(*TencentCloudClient).apiV3Conn
	tagService := &TagService{client: tcClient}
	tags, err := tagService.DescribeResourceTags(ctx, "tsf", "instance", tcClient.Region, d.Id())
	if err != nil {
		return err
	}
	_ = d.Set("tags", tags)

	return nil
}

func resourceTencentCloudTsfInstancesAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tsf_instances_attachment.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TsfService{client: meta.(*TencentCloudClient).apiV3Conn}
	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	clusterId := idSplit[0]
	instanceId := idSplit[1]

	if err := service.DeleteTsfInstancesAttachmentById(ctx, clusterId, instanceId); err != nil {
		return err
	}

	return nil
}
