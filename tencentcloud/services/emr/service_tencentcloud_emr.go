package emr

import (
	"context"
	"fmt"
	"log"
	"strconv"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svctag "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	emr "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/emr/v20190103"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

func NewEMRService(client *connectivity.TencentCloudClient) EMRService {
	return EMRService{client: client}
}

type EMRService struct {
	client *connectivity.TencentCloudClient
}

func (me *EMRService) UpdateInstance(ctx context.Context, request *emr.ScaleOutInstanceRequest) (id string, err error) {
	logId := tccommon.GetLogId(ctx)
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseEmrClient().ScaleOutInstance(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return
	}
	id = *response.Response.InstanceId
	return
}

func (me *EMRService) DeleteInstance(ctx context.Context, d *schema.ResourceData) error {
	logId := tccommon.GetLogId(ctx)
	request := emr.NewTerminateInstanceRequest()
	if v, ok := d.GetOk("instance_id"); ok {
		request.InstanceId = common.StringPtr(v.(string))
	}
	ratelimit.Check(request.GetAction())
	//API: https://cloud.tencent.com/document/api/589/34261
	_, err := me.client.UseEmrClient().TerminateInstance(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
	}
	return nil
}

func (me *EMRService) TerminateInstance(ctx context.Context, instanceId string) error {
	logId := tccommon.GetLogId(ctx)
	request := emr.NewTerminateInstanceRequest()
	request.InstanceId = helper.String(instanceId)
	ratelimit.Check(request.GetAction())
	_, err := me.client.UseEmrClient().TerminateInstance(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
	}
	return nil
}

func (me *EMRService) CreateInstance(ctx context.Context, d *schema.ResourceData) (id string, err error) {
	logId := tccommon.GetLogId(ctx)
	request := emr.NewCreateInstanceRequest()

	if v, ok := d.GetOk("scene_name"); ok {
		request.SceneName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("auto_renew"); ok {
		request.AutoRenew = common.Uint64Ptr((uint64)(v.(int)))
	}

	if v, ok := d.GetOk("product_id"); ok {
		request.ProductId = common.Uint64Ptr((uint64)(v.(int)))
	}

	if v, ok := d.GetOk("vpc_settings"); ok {
		value := v.(map[string]interface{})
		var vpcId string
		var subnetId string

		if subV, ok := value["vpc_id"]; ok {
			vpcId = subV.(string)
		}
		if subV, ok := value["subnet_id"]; ok {
			subnetId = subV.(string)
		}
		vpcSettings := &emr.VPCSettings{VpcId: &vpcId, SubnetId: &subnetId}
		request.VPCSettings = vpcSettings
	}

	if v, ok := d.GetOk("softwares"); ok {
		softwares := v.(*schema.Set).List()
		request.Software = make([]*string, 0)
		for _, software := range softwares {
			request.Software = append(request.Software, common.StringPtr(software.(string)))
		}
	}

	if v, ok := d.GetOk("cos_bucket"); ok {
		request.CosBucket = helper.String(v.(string))
	}

	if v, ok := d.GetOk("resource_spec"); ok {
		tmpResourceSpec := v.([]interface{})
		resourceSpec := tmpResourceSpec[0].(map[string]interface{})
		request.ResourceSpec = &emr.NewResourceSpec{}
		for k, v := range resourceSpec {
			if k == "master_resource_spec" {
				if len(v.([]interface{})) > 0 {
					spec := v.([]interface{})[0].(map[string]interface{})
					err = validateMultiDisks(spec)
					if err != nil {
						return
					}
					request.ResourceSpec.MasterResourceSpec = ParseResource(spec)
				}
			} else if k == "core_resource_spec" {
				if len(v.([]interface{})) > 0 {
					spec := v.([]interface{})[0].(map[string]interface{})
					err = validateMultiDisks(spec)
					if err != nil {
						return
					}
					request.ResourceSpec.CoreResourceSpec = ParseResource(spec)
				}
			} else if k == "task_resource_spec" {
				if len(v.([]interface{})) > 0 {
					spec := v.([]interface{})[0].(map[string]interface{})
					err = validateMultiDisks(spec)
					if err != nil {
						return
					}
					request.ResourceSpec.TaskResourceSpec = ParseResource(spec)
				}
			} else if k == "master_count" {
				request.ResourceSpec.MasterCount = common.Int64Ptr((int64)(v.(int)))
			} else if k == "core_count" {
				request.ResourceSpec.CoreCount = common.Int64Ptr((int64)(v.(int)))
			} else if k == "task_count" {
				request.ResourceSpec.TaskCount = common.Int64Ptr((int64)(v.(int)))
			} else if k == "common_resource_spec" {
				if len(v.([]interface{})) > 0 {
					spec := v.([]interface{})[0].(map[string]interface{})
					err = validateMultiDisks(spec)
					if err != nil {
						return
					}
					request.ResourceSpec.CommonResourceSpec = ParseResource(spec)
				}
			} else if k == "common_count" {
				request.ResourceSpec.CommonCount = common.Int64Ptr((int64)(v.(int)))
			}
		}
	}

	if v, ok := d.GetOkExists("multi_zone"); ok {
		request.MultiZone = helper.Bool(v.(bool))
		request.VersionID = helper.IntInt64(1)
	}
	if v, ok := d.GetOk("multi_zone_setting"); ok {
		multiZoneSettings := v.([]interface{})
		request.MultiZoneSettings = make([]*emr.MultiZoneSetting, 0)
		for idx, zone := range multiZoneSettings {
			if zone == nil {
				err = fmt.Errorf("multi_zone_setting element with index %d is nil", idx+1)
				return
			}
			zoneMap := zone.(map[string]interface{})
			tmpZone := &emr.MultiZoneSetting{}
			if v, ok := zoneMap["zone_tag"].(string); ok && v != "" {
				tmpZone.ZoneTag = helper.String(v)
			}
			if v, ok := zoneMap["vpc_settings"]; ok {
				value := v.(map[string]interface{})
				var vpcId string
				var subnetId string

				if subV, ok := value["vpc_id"]; ok {
					vpcId = subV.(string)
				}
				if subV, ok := value["subnet_id"]; ok {
					subnetId = subV.(string)
				}
				vpcSettings := &emr.VPCSettings{VpcId: &vpcId, SubnetId: &subnetId}
				tmpZone.VPCSettings = vpcSettings
			}
			if v, ok := zoneMap["placement"]; ok {
				tmpZone.Placement = &emr.Placement{}
				placementList := v.([]interface{})
				if len(placementList) == 0 {
					err = fmt.Errorf("placement in multi_zone_setting is empty")
					return
				}
				placement := placementList[0].(map[string]interface{})

				if projectId, ok := placement["project_id"]; ok {
					projectIdInt64, _ := strconv.ParseInt(projectId.(string), 10, 64)
					tmpZone.Placement.ProjectId = common.Int64Ptr(projectIdInt64)
				} else {
					tmpZone.Placement.ProjectId = common.Int64Ptr(0)
				}
				if z, ok := placement["zone"]; ok {
					tmpZone.Placement.Zone = common.StringPtr(z.(string))
				}
			}
			if v, ok := zoneMap["resource_spec"]; ok {
				tmpResourceSpec := v.([]interface{})
				resourceSpec := tmpResourceSpec[0].(map[string]interface{})
				tmpZone.ResourceSpec = &emr.NewResourceSpec{}
				for k, v := range resourceSpec {
					if k == "master_resource_spec" {
						if len(v.([]interface{})) > 0 {
							spec := v.([]interface{})[0].(map[string]interface{})
							err = validateMultiDisks(spec)
							if err != nil {
								return
							}
							tmpZone.ResourceSpec.MasterResourceSpec = ParseResource(spec)
						}
					} else if k == "core_resource_spec" {
						if len(v.([]interface{})) > 0 {
							spec := v.([]interface{})[0].(map[string]interface{})
							err = validateMultiDisks(spec)
							if err != nil {
								return
							}
							tmpZone.ResourceSpec.CoreResourceSpec = ParseResource(spec)
						}
					} else if k == "task_resource_spec" {
						if len(v.([]interface{})) > 0 {
							spec := v.([]interface{})[0].(map[string]interface{})
							err = validateMultiDisks(spec)
							if err != nil {
								return
							}
							tmpZone.ResourceSpec.TaskResourceSpec = ParseResource(spec)
						}
					} else if k == "master_count" {
						tmpZone.ResourceSpec.MasterCount = common.Int64Ptr((int64)(v.(int)))
					} else if k == "core_count" {
						tmpZone.ResourceSpec.CoreCount = common.Int64Ptr((int64)(v.(int)))
					} else if k == "task_count" {
						tmpZone.ResourceSpec.TaskCount = common.Int64Ptr((int64)(v.(int)))
					} else if k == "common_resource_spec" {
						if len(v.([]interface{})) > 0 {
							spec := v.([]interface{})[0].(map[string]interface{})
							err = validateMultiDisks(spec)
							if err != nil {
								return
							}
							tmpZone.ResourceSpec.CommonResourceSpec = ParseResource(spec)
						}
					} else if k == "common_count" {
						tmpZone.ResourceSpec.CommonCount = common.Int64Ptr((int64)(v.(int)))
					}
				}
			}
			request.MultiZoneSettings = append(request.MultiZoneSettings, tmpZone)
		}
	}
	if v, ok := d.GetOk("support_ha"); ok {
		request.SupportHA = common.Uint64Ptr((uint64)(v.(int)))
	} else {
		request.SupportHA = common.Uint64Ptr(0)
	}

	if v, ok := d.GetOk("instance_name"); ok {
		request.InstanceName = common.StringPtr(v.(string))
	}

	needMasterWan := d.Get("need_master_wan").(string)
	request.NeedMasterWan = common.StringPtr(needMasterWan)
	payMode := d.Get("pay_mode")
	request.PayMode = common.Uint64Ptr((uint64)(payMode.(int)))
	if v, ok := d.GetOk("placement"); ok {
		request.Placement = &emr.Placement{}
		placement := v.(map[string]interface{})

		if projectId, ok := placement["project_id"]; ok {
			projectIdInt64, _ := strconv.ParseInt(projectId.(string), 10, 64)
			request.Placement.ProjectId = common.Int64Ptr(projectIdInt64)
		} else {
			request.Placement.ProjectId = common.Int64Ptr(0)
		}
		if zone, ok := placement["zone"]; ok {
			request.Placement.Zone = common.StringPtr(zone.(string))
		}
	}

	if v, ok := d.GetOk("placement_info"); ok {
		request.Placement = &emr.Placement{}
		placementList := v.([]interface{})
		placement := placementList[0].(map[string]interface{})

		if v, ok := placement["project_id"]; ok {
			projectId := v.(int)
			request.Placement.ProjectId = helper.IntInt64(projectId)
		} else {
			request.Placement.ProjectId = helper.IntInt64(0)
		}
		if zone, ok := placement["zone"]; ok {
			request.Placement.Zone = common.StringPtr(zone.(string))
		}
	}

	if v, ok := d.GetOk("time_span"); ok {
		request.TimeSpan = common.Uint64Ptr((uint64)(v.(int)))
	}
	if v, ok := d.GetOk("time_unit"); ok {
		request.TimeUnit = common.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("login_settings"); ok {
		request.LoginSettings = &emr.LoginSettings{}
		loginSettings := v.(map[string]interface{})
		if password, ok := loginSettings["password"]; ok {
			request.LoginSettings.Password = common.StringPtr(password.(string))
		}
		if publicKeyId, ok := loginSettings["public_key_id"]; ok {
			request.LoginSettings.PublicKeyId = common.StringPtr(publicKeyId.(string))
		}
	}
	if v, ok := d.GetOk("sg_id"); ok {
		request.SgId = common.StringPtr(v.(string))
	}

	if v, ok := d.GetOk("extend_fs_field"); ok {
		request.ExtendFsField = common.StringPtr(v.(string))
	}
	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		emrTags := make([]*emr.Tag, 0)
		for k, v := range tags {
			tagKey := k
			tagValue := v
			emrTags = append(emrTags, &emr.Tag{
				TagKey:   helper.String(tagKey),
				TagValue: helper.String(tagValue),
			})
		}
		request.Tags = emrTags
	}

	if v, ok := d.GetOk("pre_executed_file_settings"); ok {
		preExecutedFileSettings := v.([]interface{})
		for idx, preExecutedFileSetting := range preExecutedFileSettings {
			if preExecutedFileSetting == nil {
				err = fmt.Errorf("pre_executed_file_settings element with index %d is nil", idx+1)
				return
			}
			preExecutedFileSettingMap := preExecutedFileSetting.(map[string]interface{})
			tmpPreExecutedFileSetting := &emr.PreExecuteFileSettings{}
			if v, ok := preExecutedFileSettingMap["args"]; ok {
				tmpPreExecutedFileSetting.Args = helper.InterfacesStringsPoint(v.([]interface{}))
			}
			if v, ok := preExecutedFileSettingMap["run_order"]; ok {
				tmpPreExecutedFileSetting.RunOrder = helper.IntInt64(v.(int))
			}
			if v, ok := preExecutedFileSettingMap["when_run"]; ok {
				tmpPreExecutedFileSetting.WhenRun = helper.String(v.(string))
			}
			if v, ok := preExecutedFileSettingMap["cos_file_name"]; ok {
				tmpPreExecutedFileSetting.CosFileName = helper.String(v.(string))
			}
			if v, ok := preExecutedFileSettingMap["cos_file_uri"]; ok {
				tmpPreExecutedFileSetting.CosFileURI = helper.String(v.(string))
			}
			if v, ok := preExecutedFileSettingMap["cos_secret_id"]; ok {
				tmpPreExecutedFileSetting.CosSecretId = helper.String(v.(string))
			}
			if v, ok := preExecutedFileSettingMap["cos_secret_key"]; ok {
				tmpPreExecutedFileSetting.CosSecretKey = helper.String(v.(string))
			}
			if v, ok := preExecutedFileSettingMap["remark"]; ok {
				tmpPreExecutedFileSetting.Remark = helper.String(v.(string))
			}
			request.PreExecutedFileSettings = append(request.PreExecutedFileSettings, tmpPreExecutedFileSetting)
		}
	}
	ratelimit.Check(request.GetAction())
	//API: https://cloud.tencent.com/document/api/589/34261
	response, err := me.client.UseEmrClient().CreateInstance(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return
	}
	id = *response.Response.InstanceId
	return
}

func (me *EMRService) DescribeInstances(ctx context.Context, filters map[string]interface{}) (clusters []*emr.ClusterInstancesInfo, errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := emr.NewDescribeInstancesRequest()

	ratelimit.Check(request.GetAction())
	// API: https://cloud.tencent.com/document/api/589/41707
	if v, ok := filters["instance_ids"]; ok {
		instances := v.([]interface{})
		request.InstanceIds = make([]*string, 0)
		for _, instance := range instances {
			request.InstanceIds = append(request.InstanceIds, common.StringPtr(instance.(string)))
		}
	}
	if v, ok := filters["display_strategy"]; ok {
		request.DisplayStrategy = common.StringPtr(v.(string))
	}
	if v, ok := filters["project_id"]; ok {
		request.ProjectId = common.Int64Ptr(v.(int64))
	}
	response, err := me.client.UseEmrClient().DescribeInstances(request)

	if err != nil {
		if sdkError, ok := err.(*sdkErrors.TencentCloudSDKError); ok {
			if sdkError.Code == "ResourceNotFound.ClusterNotFound" {
				return
			}
		}
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	clusters = response.Response.ClusterList
	return
}

func (me *EMRService) DescribeInstancesById(ctx context.Context, instanceId string, displayStrategy string) (clusters []*emr.ClusterInstancesInfo, errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := emr.NewDescribeInstancesRequest()

	ratelimit.Check(request.GetAction())
	// API: https://cloud.tencent.com/document/api/589/41707
	request.ProjectId = helper.IntInt64(-1)
	request.InstanceIds = make([]*string, 0)
	request.InstanceIds = append(request.InstanceIds, common.StringPtr(instanceId))
	request.DisplayStrategy = common.StringPtr(displayStrategy)

	response, err := me.client.UseEmrClient().DescribeInstances(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	clusters = response.Response.ClusterList
	return
}

func (me *EMRService) DescribeClusterNodes(ctx context.Context, instanceId, nodeFlag, hardwareResourceType string, offset, limit int) (nodes []*emr.NodeHardwareInfo, errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := emr.NewDescribeClusterNodesRequest()

	ratelimit.Check(request.GetAction())
	// API: https://cloud.tencent.com/document/api/589/41707
	request.InstanceId = &instanceId
	request.NodeFlag = &nodeFlag
	request.HardwareResourceType = &hardwareResourceType
	response, err := me.client.UseEmrClient().DescribeClusterNodes(request)

	if err != nil {
		if sdkError, ok := err.(*sdkErrors.TencentCloudSDKError); ok {
			if sdkError.Code == "ResourceNotFound.ClusterNotFound" {
				return
			}
		}
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	nodes = response.Response.NodeList
	return
}

func (me *EMRService) ModifyResourcesTags(ctx context.Context, region string, instanceId string, oldTags, newTags map[string]interface{}) error {
	resourceName := tccommon.BuildTagResourceName("emr", "emr-instance", region, instanceId)
	rTags, dTags := svctag.DiffTags(oldTags, newTags)
	tagService := svctag.NewTagService(me.client)
	if err := tagService.ModifyTags(ctx, resourceName, rTags, dTags); err != nil {
		return err
	}

	addTags := make([]*emr.Tag, 0)
	modifyTags := make([]*emr.Tag, 0)
	deleteTags := make([]*emr.Tag, 0)
	for k, v := range newTags {
		tagKey := k
		tageValue := v.(string)
		_, ok := oldTags[tagKey]
		if !ok {
			addTags = append(addTags, &emr.Tag{
				TagKey:   &tagKey,
				TagValue: &tageValue,
			})
		} else if oldTags[tagKey].(string) != tageValue {
			modifyTags = append(modifyTags, &emr.Tag{
				TagKey:   &tagKey,
				TagValue: &tageValue,
			})
		}
	}
	for k, v := range oldTags {
		tagKey := k
		tageValue := v.(string)
		_, ok := newTags[tagKey]
		if !ok {
			deleteTags = append(deleteTags, &emr.Tag{
				TagKey:   &tagKey,
				TagValue: &tageValue,
			})
		}
	}
	modifyResourceTags := &emr.ModifyResourceTags{
		Resource:       helper.String(resourceName),
		ResourceId:     helper.String(instanceId),
		ResourceRegion: helper.String(region),
	}
	if len(addTags) > 0 {
		modifyResourceTags.AddTags = addTags
	}
	if len(modifyTags) > 0 {
		modifyResourceTags.ModifyTags = modifyTags
	}
	if len(deleteTags) > 0 {
		modifyResourceTags.DeleteTags = deleteTags
	}

	request := emr.NewModifyResourcesTagsRequest()
	ratelimit.Check(request.GetAction())
	request.ModifyType = helper.String("Cluster")
	request.ModifyResourceTagsInfoList = []*emr.ModifyResourceTags{modifyResourceTags}

	response, err := me.client.UseEmrClient().ModifyResourcesTags(request)
	if err != nil {
		return err
	}
	if response != nil && response.Response != nil && len(response.Response.FailList) > 0 {
		return fmt.Errorf("file resource list: %v", response.Response.FailList)
	}
	return nil
}

func (me *EMRService) DescribeEmrUserManagerById(ctx context.Context, instanceId string, userName string) (userManager *emr.DescribeUsersForUserManagerResponseParams, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := emr.NewDescribeUsersForUserManagerRequest()
	request.InstanceId = &instanceId
	request.UserManagerFilter = &emr.UserManagerFilter{
		UserName: &userName,
	}
	request.PageNo = helper.IntInt64(0)
	request.PageSize = helper.IntInt64(100)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseEmrClient().DescribeUsersForUserManager(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	userManager = response.Response
	return
}

func (me *EMRService) DeleteEmrUserManagerById(ctx context.Context, instanceId string, userName string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := emr.NewDeleteUserManagerUserListRequest()
	request.InstanceId = &instanceId
	request.UserNameList = []*string{&userName}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseEmrClient().DeleteUserManagerUserList(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *EMRService) DescribeEmrCvmQuotaByFilter(ctx context.Context, param map[string]interface{}) (cvmQuota *emr.DescribeCvmQuotaResponseParams, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = emr.NewDescribeCvmQuotaRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "ClusterId" {
			request.ClusterId = v.(*string)
		}
		if k == "ZoneId" {
			request.ZoneId = v.(*int64)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseEmrClient().DescribeCvmQuota(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	cvmQuota = response.Response
	return
}

func (me *EMRService) DescribeEmrAutoScaleRecordsByFilter(ctx context.Context, param map[string]interface{}) (autoScaleRecords []*emr.AutoScaleRecord, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = emr.NewDescribeAutoScaleRecordsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "InstanceId" {
			request.InstanceId = v.(*string)
		}
		if k == "Filters" {
			request.Filters = v.([]*emr.KeyValue)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset int64 = 0
		limit  int64 = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseEmrClient().DescribeAutoScaleRecords(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.RecordList) < 1 {
			break
		}
		autoScaleRecords = append(autoScaleRecords, response.Response.RecordList...)
		if len(response.Response.RecordList) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *EMRService) GetSLInstanceStatus(ctx context.Context, instanceId string) (instanceInfo *emr.SLInstanceInfo, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := emr.NewDescribeSLInstanceListRequest()
	filter := &emr.Filters{
		Name:   helper.String("ClusterId"),
		Values: []*string{&instanceId},
	}
	request.Filters = []*emr.Filters{filter}
	request.DisplayStrategy = helper.String("clusterList")

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, e := me.client.UseEmrClient().DescribeSLInstanceList(request)
	if e != nil {
		errRet = e
		return
	}

	if response.Response != nil && len(response.Response.InstancesList) > 0 {
		instanceInfo = response.Response.InstancesList[0]
	}
	return

}

func (me *EMRService) SLInstanceStateRefreshFunc(instanceId string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		ctx := tccommon.ContextNil

		object, err := me.GetSLInstanceStatus(ctx, instanceId)

		if err != nil {
			if e, ok := err.(*sdkErrors.TencentCloudSDKError); ok {
				if e.GetCode() == "UnauthorizedOperation" {
					return &emr.SLInstanceInfo{}, "-2", nil
				}
			}
			return nil, "", err
		}

		if object == nil {
			return &emr.SLInstanceInfo{}, "-2", nil
		}

		return object, helper.UInt64ToStr(*object.Status), nil
	}
}

func (me *EMRService) DescribeLiteHbaseInstancesByFilter(ctx context.Context, param map[string]interface{}) (instances []*emr.SLInstanceInfo, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = emr.NewDescribeSLInstanceListRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "DisplayStrategy" {
			request.DisplayStrategy = v.(*string)
		}
		if k == "OrderField" {
			request.OrderField = v.(*string)
		}
		if k == "Asc" {
			request.Asc = v.(*int64)
		}
		if k == "Filters" {
			request.Filters = v.([]*emr.Filters)
		}
	}

	var (
		offset   int64 = 0
		limit    int64 = 100
		response *emr.DescribeSLInstanceListResponse
		innerErr error
	)

	for {
		request.Offset = helper.Int64(offset)
		request.Limit = helper.Int64(limit)

		err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())
			response, innerErr = me.client.UseEmrV20190103Client().DescribeSLInstanceList(request)
			if innerErr != nil {
				return tccommon.RetryError(innerErr)
			}
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

			if response != nil && response.Response != nil && len(response.Response.InstancesList) > 0 {
				instances = append(instances, response.Response.InstancesList...)
			}
			return nil
		})
		if err != nil {
			errRet = err
			return
		}

		if len(response.Response.InstancesList) < int(limit) {
			break
		}
		offset += limit
	}

	return
}

func (me *EMRService) TerminateClusterNodes(ctx context.Context, instanceIds []string, instanceId, nodeFlag string) (flowId int64, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := emr.NewTerminateClusterNodesRequest()
	request.CvmInstanceIds = helper.Strings(instanceIds)
	request.InstanceId = helper.String(instanceId)
	request.NodeFlag = helper.String(nodeFlag)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	var (
		response *emr.TerminateClusterNodesResponse
		innerErr error
	)
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		response, innerErr = me.client.UseEmrClient().TerminateClusterNodes(request)
		if innerErr != nil {
			return tccommon.RetryError(innerErr)
		}
		return nil
	})
	if err != nil {
		errRet = err
		return
	}

	if response.Response != nil && response.Response.FlowId != nil {
		flowId = *response.Response.FlowId
		return
	}
	return

}

func (me *EMRService) FlowStatusRefreshFunc(instanceId, flowId, flowType string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {

		request := emr.NewDescribeClusterFlowStatusDetailRequest()
		request.InstanceId = helper.String(instanceId)
		request.FlowParam = &emr.FlowParam{
			FKey:   helper.String(flowType),
			FValue: helper.String(flowId),
		}

		var (
			response *emr.DescribeClusterFlowStatusDetailResponse
			innerErr error
		)
		err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())
			response, innerErr = me.client.UseEmrClient().DescribeClusterFlowStatusDetail(request)
			if innerErr != nil {
				return tccommon.RetryError(innerErr)
			}
			return nil
		})
		if err != nil {
			return nil, "", err
		}

		if response.Response == nil || response.Response.FlowTotalStatus == nil {
			return nil, "", fmt.Errorf("Not found flow.")
		}
		return response.Response.FlowTotalStatus, helper.Int64ToStr(*response.Response.FlowTotalStatus), nil
	}
}

func (me *EMRService) ScaleOutInstance(ctx context.Context, request *emr.ScaleOutInstanceRequest) (traceId string, err error) {
	logId := tccommon.GetLogId(ctx)
	err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		ratelimit.Check(request.GetAction())
		response, e := me.client.UseEmrClient().ScaleOutInstance(request)
		if e != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), e.Error())
			return tccommon.RetryError(e)
		}
		traceId = *response.Response.TraceId
		return nil
	})

	if err != nil {
		return
	}
	return
}

func (me *EMRService) DescribeEmrYarnById(ctx context.Context, instanceId string) (ret *emr.DescribeGlobalConfigResponseParams, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := emr.NewDescribeGlobalConfigRequest()
	request.InstanceId = helper.String(instanceId)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseEmrClient().DescribeGlobalConfig(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	ret = response.Response
	return
}

func (me *EMRService) DescribeEmrJobStatusDetailByFilter(ctx context.Context, param map[string]interface{}) (ret *emr.DescribeClusterFlowStatusDetailResponseParams, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = emr.NewDescribeClusterFlowStatusDetailRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "InstanceId" {
			request.InstanceId = v.(*string)
		}
		if k == "FlowParam" {
			request.FlowParam = v.(*emr.FlowParam)
		}
		if k == "NeedExtraDetail" {
			request.NeedExtraDetail = v.(*bool)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseEmrClient().DescribeClusterFlowStatusDetail(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || response.Response == nil {
		return
	}

	ret = response.Response
	return
}

func (me *EMRService) DescribeEmrServiceNodeInfosByFilter(ctx context.Context, param map[string]interface{}) (ret *emr.DescribeServiceNodeInfosResponseParams, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = emr.NewDescribeServiceNodeInfosRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "InstanceId" {
			request.InstanceId = v.(*string)
		}
		if k == "Offset" {
			request.Offset = v.(*int64)
		}
		if k == "Limit" {
			request.Limit = v.(*int64)
		}
		if k == "SearchText" {
			request.SearchText = v.(*string)
		}
		if k == "ConfStatus" {
			request.ConfStatus = v.(*int64)
		}
		if k == "MaintainStateId" {
			request.MaintainStateId = v.(*int64)
		}
		if k == "OperatorStateId" {
			request.OperatorStateId = v.(*int64)
		}
		if k == "HealthStateId" {
			request.HealthStateId = v.(*string)
		}
		if k == "ServiceName" {
			request.ServiceName = v.(*string)
		}
		if k == "NodeTypeName" {
			request.NodeTypeName = v.(*string)
		}
		if k == "DataNodeMaintenanceId" {
			request.DataNodeMaintenanceId = v.(*int64)
		}
		if k == "SearchFields" {
			searchFields := v.([]interface{})
			for _, searchField := range searchFields {
				request.SearchFields = append(request.SearchFields, searchField.(*emr.SearchItem))
			}
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseEmrClient().DescribeServiceNodeInfos(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || response.Response == nil {
		return
	}

	ret = response.Response
	return
}

func (me *EMRService) DescribeEmrAutoScaleStrategyById(ctx context.Context, instanceId string) (ret *emr.DescribeAutoScaleStrategiesResponseParams, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := emr.NewDescribeAutoScaleStrategiesRequest()
	request.InstanceId = helper.String(instanceId)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseEmrClient().DescribeAutoScaleStrategies(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	ret = response.Response
	return
}

func (me *EMRService) DescribeEmrAutoScaleStrategy(ctx context.Context, instanceId, name string) (ret *emr.DescribeAutoScaleStrategiesResponseParams, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := emr.NewDescribeAutoScaleStrategiesRequest()
	request.InstanceId = helper.String(instanceId)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseEmrClient().DescribeAutoScaleStrategies(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	ret = response.Response
	return
}

func (me *EMRService) DeleteAutoScaleStrategy(ctx context.Context, instanceId, strategyType string, strategyId int64) error {
	logId := tccommon.GetLogId(ctx)
	var (
		request  = emr.NewDeleteAutoScaleStrategyRequest()
		response = emr.NewDeleteAutoScaleStrategyResponse()
	)
	request.InstanceId = helper.String(instanceId)

	request.StrategyType = helper.StrToInt64Point(strategyType)
	request.StrategyId = helper.Int64(strategyId)

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := me.client.UseEmrClient().DeleteAutoScaleStrategyWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s delete emr auto scale strategy failed, reason:%+v", logId, err)
		return err
	}

	_ = response
	return nil

}

func (me *EMRService) DescribeEmrClusterNewById(ctx context.Context, instanceId string) (ret *emr.ClusterInstancesInfo, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := emr.NewDescribeInstancesRequest()
	request.DisplayStrategy = helper.String("clusterList")

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseEmrClient().DescribeInstances(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.ClusterList) < 1 {
		return
	}

	ret = response.Response.ClusterList[0]
	return
}
