package cfs

import (
	"context"
	"fmt"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cfs "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cfs/v20190719"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

func ResourceTencentCloudCfsFileSystem() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCfsFileSystemCreate,
		Read:   resourceTencentCloudCfsFileSystemRead,
		Update: resourceTencentCloudCfsFileSystemUpdate,
		Delete: resourceTencentCloudCfsFileSystemDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Name of a file system.",
			},
			"availability_zone": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The available zone that the file system locates at.",
			},
			"access_group_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of a access group.",
			},
			"net_interface": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      CFS_NET_VPC,
				ValidateFunc: tccommon.ValidateAllowedStringValue(CFS_NET),
				Description:  "Network type, Default `VPC`. Valid values: `VPC` and `CCN`. Select `VPC` for a Standard or High-Performance file system, and `CCN` for a Standard Turbo or High-Performance Turbo one.",
			},
			"protocol": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      CFS_PROTOCOL_NFS,
				ValidateFunc: tccommon.ValidateAllowedStringValue(CFS_PROTOCOL),
				ForceNew:     true,
				Description:  "File system protocol. Valid values: `NFS`, `CIFS`, `TURBO`. If this parameter is left empty, `NFS` is used by default. For the Turbo series, you must set this parameter to `TURBO`.",
			},
			"storage_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      CFS_STORAGETYPE_SD,
				ValidateFunc: tccommon.ValidateAllowedStringValue(CFS_STORAGETYPE),
				ForceNew:     true,
				Description:  "Storage type of the file system. Valid values: `SD` (Standard), `HP` (High-Performance), `TB` (Standard Turbo), and `TP` (High-Performance Turbo). Default value: `SD`.",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "ID of a VPC network.",
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "ID of a subnet.",
			},
			"mount_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: "IP of mount point.",
			},
			"ccn_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "CCN instance ID (required if the network type is CCN).",
			},
			"cidr_block": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "CCN IP range used by the CFS (required if the network type is CCN), which cannot conflict with other IP ranges bound in CCN.",
			},
			"capacity": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "File system capacity, in GiB (required for the Turbo series). For Standard Turbo, the minimum purchase required is 40,960 GiB (40 TiB) and the expansion increment is 20,480 GiB (20 TiB). For High-Performance Turbo, the minimum purchase required is 20,480 GiB (20 TiB) and the expansion increment is 10,240 GiB (10 TiB).",
			},
			// computed
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Create time of the file system.",
			},
			"fs_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Mount root-directory.",
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Instance tags.",
			},
		},
	}
}

func resourceTencentCloudCfsFileSystemCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cfs_file_system.create")()
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	cfsService := CfsService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}

	request := cfs.NewCreateCfsFileSystemRequest()
	request.Zone = helper.String(d.Get("availability_zone").(string))
	request.NetInterface = helper.String(d.Get("net_interface").(string))
	request.PGroupId = helper.String(d.Get("access_group_id").(string))
	request.Protocol = helper.String(d.Get("protocol").(string))
	request.VpcId = helper.String(d.Get("vpc_id").(string))
	request.SubnetId = helper.String(d.Get("subnet_id").(string))
	request.StorageType = helper.String(d.Get("storage_type").(string))
	if v, ok := d.GetOk("name"); ok {
		request.FsName = helper.String(v.(string))
	}
	if v, ok := d.GetOk("mount_ip"); ok {
		request.MountIP = helper.String(v.(string))
	}
	if v, ok := d.GetOk("ccn_id"); ok {
		request.CcnId = helper.String(v.(string))
	}
	if v, ok := d.GetOk("cidr_block"); ok {
		request.CidrBlock = helper.String(v.(string))
	}
	if v, ok := d.GetOkExists("capacity"); ok {
		request.Capacity = helper.IntUint64(v.(int))
	}

	if v := helper.GetTags(d, "tags"); len(v) > 0 {
		for tagKey, tagValue := range v {
			tag := cfs.TagInfo{
				TagKey:   helper.String(tagKey),
				TagValue: helper.String(tagValue),
			}
			request.ResourceTags = append(request.ResourceTags, &tag)
		}
	}

	fsId := ""
	err := resource.Retry(3*tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		response, err := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCfsClient().CreateCfsFileSystem(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			return tccommon.RetryError(err)
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response.Response.FileSystemId == nil {
			err = fmt.Errorf("file system id is nil")
			return resource.NonRetryableError(err)
		}
		fsId = *response.Response.FileSystemId
		return nil
	})
	if err != nil {
		return err
	}
	d.SetId(fsId)

	// wait for success status
	err = resource.Retry(3*tccommon.ReadRetryTimeout, func() *resource.RetryError {
		fileSystems, errRet := cfsService.DescribeFileSystem(ctx, fsId, "", "")
		if errRet != nil {
			return tccommon.RetryError(errRet, tccommon.InternalError)
		}
		if len(fileSystems) < 1 {
			return resource.RetryableError(fmt.Errorf("file system %s not exist", fsId))
		}
		if *fileSystems[0].LifeCycleState == CFS_FILE_SYSTEM_STATUS_CREATING {
			return resource.RetryableError(fmt.Errorf("file system status is %s, retry...", *fileSystems[0].LifeCycleState))
		}
		return nil
	})
	if err != nil {
		return err
	}
	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		tcClient := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
		tagService := &TagService{client: tcClient}
		resourceName := tccommon.BuildTagResourceName("cfs", "filesystem", tcClient.Region, d.Id())
		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
	}

	return resourceTencentCloudCfsFileSystemRead(d, meta)
}

func resourceTencentCloudCfsFileSystemRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cfs_file_system.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	fsId := d.Id()
	cfsService := CfsService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}
	var fileSystem *cfs.FileSystemInfo
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		fileSystems, errRet := cfsService.DescribeFileSystem(ctx, fsId, "", "")
		if errRet != nil {
			return tccommon.RetryError(errRet, tccommon.InternalError)
		}
		if len(fileSystems) > 0 {
			fileSystem = fileSystems[0]
		}
		return nil
	})
	if err != nil {
		return err
	}
	if fileSystem == nil {
		d.SetId("")
		return nil
	}

	_ = d.Set("name", fileSystem.FsName)
	_ = d.Set("availability_zone", fileSystem.Zone)
	_ = d.Set("access_group_id", fileSystem.PGroup.PGroupId)
	_ = d.Set("protocol", fileSystem.Protocol)
	_ = d.Set("create_time", fileSystem.CreationTime)
	_ = d.Set("storage_type", fileSystem.StorageType)
	_ = d.Set("capacity", fileSystem.SizeLimit)

	var mountTarget *cfs.MountInfo
	err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		targets, errRet := cfsService.DescribeMountTargets(ctx, fsId)
		if errRet != nil {
			return tccommon.RetryError(errRet, tccommon.InternalError)
		}
		if len(targets) > 0 {
			mountTarget = targets[0]
		}
		return nil
	})
	if err != nil {
		return err
	}
	tcClient := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
	tagService := &TagService{client: tcClient}
	tags, err := tagService.DescribeResourceTags(ctx, "cfs", "filesystem", tcClient.Region, d.Id())
	if err != nil {
		return err
	}

	_ = d.Set("tags", tags)

	if mountTarget != nil {
		_ = d.Set("vpc_id", mountTarget.VpcId)
		_ = d.Set("subnet_id", mountTarget.SubnetId)
		_ = d.Set("mount_ip", mountTarget.IpAddress)
		_ = d.Set("ccn_id", mountTarget.CcnID)
		_ = d.Set("cidr_block", mountTarget.CidrBlock)
		_ = d.Set("net_interface", mountTarget.NetworkInterface)
	}

	return nil
}

func resourceTencentCloudCfsFileSystemUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cfs_file_system.update")()
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	fsId := d.Id()
	cfsService := CfsService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}

	immutableArgs := []string{"ccn_id", "cidr_block", "net_interface", "capacity"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	d.Partial(true)

	if d.HasChange("name") {
		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			errRet := cfsService.ModifyFileSystemName(ctx, fsId, d.Get("name").(string))
			if errRet != nil {
				return tccommon.RetryError(errRet)
			}
			return nil
		})
		if err != nil {
			return err
		}

	}

	if d.HasChange("access_group_id") {
		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			errRet := cfsService.ModifyFileSystemAccessGroup(ctx, fsId, d.Get("access_group_id").(string))
			if errRet != nil {
				return tccommon.RetryError(errRet)
			}
			return nil
		})
		if err != nil {
			return err
		}

	}

	if d.HasChange("tags") {

		oldValue, newValue := d.GetChange("tags")
		replaceTags, deleteTags := diffTags(oldValue.(map[string]interface{}), newValue.(map[string]interface{}))

		tcClient := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
		tagService := &TagService{client: tcClient}
		resourceName := tccommon.BuildTagResourceName("cfs", "filesystem", tcClient.Region, d.Id())
		err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags)
		if err != nil {
			return err
		}

	}

	d.Partial(false)

	return resourceTencentCloudCfsFileSystemRead(d, meta)
}

func resourceTencentCloudCfsFileSystemDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cfs_file_system.delete")()
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	fsId := d.Id()
	cfsService := CfsService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		errRet := cfsService.DeleteFileSystem(ctx, fsId)
		if errRet != nil {
			return tccommon.RetryError(errRet)
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
