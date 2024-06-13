package cvm

import (
	"context"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"

	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
)

func dataSourceTencentCloudInstancesSetReadPostFillRequest0(ctx context.Context, req map[string]interface{}) error {
	d := tccommon.ResourceDataFromContext(ctx)
	if v, ok := d.GetOk("tags"); ok {
		for key, value := range v.(map[string]interface{}) {
			req["tag:"+key] = value.(string)
		}
	}

	return nil
}

func dataSourceTencentCloudInstancesSetReadPostHandleResponse0(ctx context.Context, req map[string]interface{}, resp *[]*cvm.Instance) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	d := tccommon.ResourceDataFromContext(ctx)
	instances := *resp
	instanceList := make([]map[string]interface{}, 0, len(instances))
	ids := make([]string, 0, len(instances))
	for _, instance := range instances {
		mapping := map[string]interface{}{
			"instance_id":                instance.InstanceId,
			"instance_name":              instance.InstanceName,
			"instance_type":              instance.InstanceType,
			"cpu":                        instance.CPU,
			"memory":                     instance.Memory,
			"availability_zone":          instance.Placement.Zone,
			"project_id":                 instance.Placement.ProjectId,
			"image_id":                   instance.ImageId,
			"instance_charge_type":       instance.InstanceChargeType,
			"system_disk_type":           instance.SystemDisk.DiskType,
			"system_disk_size":           instance.SystemDisk.DiskSize,
			"system_disk_id":             instance.SystemDisk.DiskId,
			"vpc_id":                     instance.VirtualPrivateCloud.VpcId,
			"subnet_id":                  instance.VirtualPrivateCloud.SubnetId,
			"internet_charge_type":       instance.InternetAccessible.InternetChargeType,
			"internet_max_bandwidth_out": instance.InternetAccessible.InternetMaxBandwidthOut,
			"allocate_public_ip":         instance.InternetAccessible.PublicIpAssigned,
			"status":                     instance.InstanceState,
			"security_groups":            helper.StringsInterfaces(instance.SecurityGroupIds),
			"tags":                       flattenCvmTagsMapping(instance.Tags),
			"create_time":                instance.CreatedTime,
			"expired_time":               instance.ExpiredTime,
			"instance_charge_type_prepaid_renew_flag": instance.RenewFlag,
			"cam_role_name": instance.CamRoleName,
		}
		if len(instance.PublicIpAddresses) > 0 {
			mapping["public_ip"] = *instance.PublicIpAddresses[0]
		}
		if len(instance.PrivateIpAddresses) > 0 {
			mapping["private_ip"] = *instance.PrivateIpAddresses[0]
		}
		dataDisks := make([]map[string]interface{}, 0, len(instance.DataDisks))
		for _, v := range instance.DataDisks {
			dataDisk := map[string]interface{}{
				"data_disk_type":       v.DiskType,
				"data_disk_size":       v.DiskSize,
				"data_disk_id":         v.DiskId,
				"delete_with_instance": v.DeleteWithInstance,
			}
			dataDisks = append(dataDisks, dataDisk)
		}
		mapping["data_disks"] = dataDisks
		instanceList = append(instanceList, mapping)
		ids = append(ids, *instance.InstanceId)
	}
	log.Printf("[DEBUG]%s set instance attribute finished", logId)
	d.SetId(helper.DataResourceIdsHash(ids))
	err := d.Set("instance_list", instanceList)
	if err != nil {
		log.Printf("[CRITAL]%s provider set instance list fail, reason:%s\n ", logId, err.Error())
		return err
	}
	context.WithValue(ctx, "instanceList", instanceList)
	return nil
}

func dataSourceTencentCloudInstancesSetReadOutputContent(ctx context.Context) interface{} {
	instanceList := ctx.Value("instanceList")
	return instanceList
}
