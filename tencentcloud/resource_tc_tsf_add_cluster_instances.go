/*
Provides a resource to create a tsf add_cluster_instances

Example Usage

```hcl
resource "tencentcloud_tsf_add_cluster_instances" "add_cluster_instances" {
  cluster_id = "cluster-123456"
  instance_id_list =
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
  security_group_ids =
}
```

Import

tsf add_cluster_instances can be imported using the id, e.g.

```
terraform import tencentcloud_tsf_add_cluster_instances.add_cluster_instances add_cluster_instances_id
```
*/
package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tsf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tsf/v20180326"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudTsfAddClusterInstances() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTsfAddClusterInstancesCreate,
		Read:   resourceTencentCloudTsfAddClusterInstancesRead,
		Delete: resourceTencentCloudTsfAddClusterInstancesDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Cluster ID.",
			},

			"instance_id_list": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Cloud server ID list.",
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
				Description: "Cloud server import mode, required for virtual machine clusters, not required for container clusters. R : Reinstall TSF system image, M: Manual installation of agent.",
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
							Description: " dockerd --graph specifies the value, default is /var/lib/docker Note: This field may return null, indicating that no valid values can be obtained.",
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

func resourceTencentCloudTsfAddClusterInstancesCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tsf_add_cluster_instances.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request   = tsf.NewAddClusterInstancesRequest()
		response  = tsf.NewAddClusterInstancesResponse()
		clusterId string
	)
	if v, ok := d.GetOk("cluster_id"); ok {
		clusterId = v.(string)
		request.ClusterId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("instance_id_list"); ok {
		instanceIdListSet := v.(*schema.Set).List()
		for i := range instanceIdListSet {
			instanceIdList := instanceIdListSet[i].(string)
			request.InstanceIdList = append(request.InstanceIdList, &instanceIdList)
		}
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
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate tsf addClusterInstances failed, reason:%+v", logId, err)
		return err
	}

	clusterId = *response.Response.ClusterId
	d.SetId(clusterId)

	return resourceTencentCloudTsfAddClusterInstancesRead(d, meta)
}

func resourceTencentCloudTsfAddClusterInstancesRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tsf_add_cluster_instances.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudTsfAddClusterInstancesDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tsf_add_cluster_instances.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
