/*
Provides a resource to create a tsf instances_attachment

Example Usage

```hcl
resource "tencentcloud_tsf_instances_attachment" "instances_attachment" {
  cluster_id = "cluster-123456"
  instance_id_list = [""]
  os_name = "Ubuntu 20.04"
  image_id = "img-123456"
  password = "MyP@ssw0rd"
  key_id = "key-123456"
  sg_id = "sg-123456"
  instance_import_mode = "R"
  os_customize_type = "my_customize"
  feature_id_list =
  instance_advanced_settings {
	mount_target = "/mnt/data"
	docker_graph_path = "/var/lib/docker"
  }
  security_group_ids = [""]
}
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tsf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tsf/v20180326"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudTsfInstancesAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTsfInstancesAttachmentCreate,
		Read:   resourceTencentCloudTsfInstancesAttachmentRead,
		Delete: resourceTencentCloudTsfInstancesAttachmentDelete,
		// Importer: &schema.ResourceImporter{
		// 	State: schema.ImportStatePassthrough,
		// },
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Cluster ID.",
			},

			"instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Cloud server ID.",
			},

			"os_name": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Operating system name.",
			},

			"image_id": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Operating system image ID.",
			},

			"password": {
				Optional:    true,
				ForceNew:    true,
				Sensitive:   true,
				Type:        schema.TypeString,
				Description: "Reset system password.",
			},

			"key_id": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Associated key for system reinstallation.",
			},

			"sg_id": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Security group setting.",
			},

			"instance_import_mode": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Cloud server import mode, required for virtual machine clusters, not required for container clusters. R: Reinstall TSF system image, M: Manual installation of agent.",
			},

			"os_customize_type": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Image customization type.",
			},

			"feature_id_list": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Image feature ID list.",
			},

			"instance_advanced_settings": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Additional instance parameter information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"mount_target": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Data disk mount point, data disks are not mounted by default. Data disks with formatted ext3, ext4, xfs file systems will be mounted directly, other file systems or unformatted data disks will be automatically formatted as ext4 and mounted. Please back up your data! This setting does not take effect for cloud servers with no data disks or multiple data disks. Note: This field may return null, indicating that no valid values can be obtained.",
						},
						"docker_graph_path": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Dockerd --graph specifies the value, default is /var/lib/docker Note: This field may return null, indicating that no valid values can be obtained.",
						},
					},
				},
			},

			"security_group_ids": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Security group.",
			},
		},
	}
}

func resourceTencentCloudTsfInstancesAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tsf_instances_attachment.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = tsf.NewAddClusterInstancesRequest()
		clusterId  string
		instanceId string
	)
	if v, ok := d.GetOk("cluster_id"); ok {
		clusterId = v.(string)
		request.ClusterId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceIdList = helper.Strings([]string{v.(string)})
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

	if v, ok := d.GetOk("os_customize_type"); ok {
		request.OsCustomizeType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("feature_id_list"); ok {
		featureIdListSet := v.(*schema.Set).List()
		for i := range featureIdListSet {
			featureIdList := featureIdListSet[i].(string)
			request.FeatureIdList = append(request.FeatureIdList, &featureIdList)
		}
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "instance_advanced_settings"); ok {
		instanceAdvancedSettings := tsf.InstanceAdvancedSettings{}
		if v, ok := dMap["mount_target"]; ok {
			instanceAdvancedSettings.MountTarget = helper.String(v.(string))
		}
		if v, ok := dMap["docker_graph_path"]; ok {
			instanceAdvancedSettings.DockerGraphPath = helper.String(v.(string))
		}
		request.InstanceAdvancedSettings = &instanceAdvancedSettings
	}

	if v, ok := d.GetOk("security_group_ids"); ok {
		securityGroupIdsSet := v.(*schema.Set).List()
		for i := range securityGroupIdsSet {
			securityGroupIds := securityGroupIdsSet[i].(string)
			request.SecurityGroupIds = append(request.SecurityGroupIds, &securityGroupIds)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTsfClient().AddClusterInstances(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate tsf addClusterInstancesOperation failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(clusterId + FILED_SP + instanceId)

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

	if instancesAttachment.ClusterId != nil {
		_ = d.Set("cluster_id", instancesAttachment.ClusterId)
	}

	if instancesAttachment.InstanceId != nil {
		_ = d.Set("instance_id", instancesAttachment.InstanceId)
	}

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

	// if instancesAttachment.OsCustomizeType != nil {
	// 	_ = d.Set("os_customize_type", instancesAttachment.OsCustomizeType)
	// }

	// if instancesAttachment.FeatureIdList != nil {
	// 	_ = d.Set("feature_id_list", instancesAttachment.FeatureIdList)
	// }

	// if instancesAttachment.InstanceAdvancedSettings != nil {
	// 	instanceAdvancedSettingsMap := map[string]interface{}{}

	// 	if instancesAttachment.InstanceAdvancedSettings.MountTarget != nil {
	// 		instanceAdvancedSettingsMap["mount_target"] = instancesAttachment.InstanceAdvancedSettings.MountTarget
	// 	}

	// 	if instancesAttachment.InstanceAdvancedSettings.DockerGraphPath != nil {
	// 		instanceAdvancedSettingsMap["docker_graph_path"] = instancesAttachment.InstanceAdvancedSettings.DockerGraphPath
	// 	}

	// 	_ = d.Set("instance_advanced_settings", []interface{}{instanceAdvancedSettingsMap})
	// }

	// if instancesAttachment.SecurityGroupIds != nil {
	// 	_ = d.Set("security_group_ids", instancesAttachment.SecurityGroupIds)
	// }

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
