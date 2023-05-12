package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	lighthouse "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/lighthouse/v20200324"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

type LightHouseService struct {
	client *connectivity.TencentCloudClient
}

// lighthouse instance

func (me *LightHouseService) DescribeLighthouseInstanceById(ctx context.Context, instanceId string) (instance *lighthouse.Instance, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = lighthouse.NewDescribeInstancesRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()

	request.InstanceIds = append(request.InstanceIds, helper.String(instanceId))
	ratelimit.Check(request.GetAction())

	var offset int64 = 0
	var pageSize int64 = 100
	instances := make([]*lighthouse.Instance, 0)

	for {
		request.Offset = &offset
		request.Limit = &pageSize
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseLighthouseClient().DescribeInstances(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.InstanceSet) < 1 {
			break
		}
		instances = append(instances, response.Response.InstanceSet...)
		if len(response.Response.InstanceSet) < int(pageSize) {
			break
		}
		offset += pageSize
	}

	if len(instances) < 1 {
		return
	}
	instance = instances[0]

	return
}

func (me *LightHouseService) DeleteLighthouseInstanceById(ctx context.Context, id string) (errRet error) {
	logId := getLogId(ctx)

	request := lighthouse.NewTerminateInstancesRequest()
	request.InstanceIds = append(request.InstanceIds, &id)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "delete object", request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseLighthouseClient().TerminateInstances(request)
	if err != nil {
		errRet = err
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *LightHouseService) IsolateLighthouseInstanceById(ctx context.Context, id string) (errRet error) {
	logId := getLogId(ctx)

	request := lighthouse.NewIsolateInstancesRequest()
	request.InstanceIds = append(request.InstanceIds, &id)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "delete object", request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseLighthouseClient().IsolateInstances(request)
	if err != nil {
		errRet = err
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *LightHouseService) LighthouseBlueprintStateRefreshFunc(blueprintId string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		ctx := contextNil

		object, err := me.DescribeLighthouseBlueprintById(ctx, blueprintId)

		if err != nil {
			return nil, "", err
		}

		return object, helper.PString(object.BlueprintState), nil
	}
}

func (me *LightHouseService) DeleteLighthouseBlueprintById(ctx context.Context, blueprintId string) (errRet error) {
	logId := getLogId(ctx)

	request := lighthouse.NewDeleteBlueprintsRequest()
	request.BlueprintIds = []*string{&blueprintId}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseLighthouseClient().DeleteBlueprints(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *LightHouseService) DescribeLighthouseBlueprintById(ctx context.Context, blueprintId string) (blueprint *lighthouse.Blueprint, errRet error) {
	logId := getLogId(ctx)

	request := lighthouse.NewDescribeBlueprintsRequest()
	request.BlueprintIds = []*string{&blueprintId}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseLighthouseClient().DescribeBlueprints(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.BlueprintSet) < 1 {
		return
	}

	blueprint = response.Response.BlueprintSet[0]
	return
}

func (me *LightHouseService) DescribeLighthouseFirewallRuleById(ctx context.Context, instance_id string) (firewallRules []*lighthouse.FirewallRuleInfo, errRet error) {
	logId := getLogId(ctx)

	request := lighthouse.NewDescribeFirewallRulesRequest()

	request.InstanceId = helper.String(instance_id)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	limit := 50
	offset := 0
	firewallRules = make([]*lighthouse.FirewallRuleInfo, 0)
	for {
		ratelimit.Check(request.GetAction())
		request.Limit = helper.IntInt64(limit)
		request.Offset = helper.IntInt64(offset)
		response, err := me.client.UseLighthouseClient().DescribeFirewallRules(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
		firewallRules = append(firewallRules, response.Response.FirewallRuleSet...)

		if len(response.Response.FirewallRuleSet) < limit {
			break
		}
		offset += limit
	}

	return
}

func (me *LightHouseService) DescribeLighthouseFirewallRulesTemplateByFilter(ctx context.Context) (firewallRulesTemplate []*lighthouse.FirewallRuleInfo, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = lighthouse.NewDescribeFirewallRulesTemplateRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseLighthouseClient().DescribeFirewallRulesTemplate(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response != nil && response.Response != nil && len(response.Response.FirewallRuleSet) != 0 {
		firewallRulesTemplate = append(firewallRulesTemplate, response.Response.FirewallRuleSet...)
	}
	return
}

func (me *LightHouseService) DescribeLighthouseDiskBackupById(ctx context.Context, diskBackupId string) (diskBackup *lighthouse.DiskBackup, errRet error) {
	logId := getLogId(ctx)

	request := lighthouse.NewDescribeDiskBackupsRequest()
	request.DiskBackupIds = []*string{&diskBackupId}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseLighthouseClient().DescribeDiskBackups(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.DiskBackupSet) < 1 {
		return
	}

	diskBackup = response.Response.DiskBackupSet[0]
	return
}

func (me *LightHouseService) DeleteLighthouseDiskBackupById(ctx context.Context, diskBackupId string) (errRet error) {
	logId := getLogId(ctx)

	request := lighthouse.NewDeleteDiskBackupsRequest()
	request.DiskBackupIds = []*string{&diskBackupId}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseLighthouseClient().DeleteDiskBackups(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *LightHouseService) LighthouseDiskBackupStateRefreshFunc(diskBackupId string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		ctx := contextNil

		object, err := me.DescribeLighthouseDiskBackupById(ctx, diskBackupId)

		if err != nil {
			return nil, "", err
		}

		return object, helper.PString(object.DiskBackupState), nil
	}
}

func (me *LightHouseService) LighthouseApplyDiskBackupStateRefreshFunc(diskBackupId string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		ctx := contextNil

		object, err := me.DescribeLighthouseDiskBackupById(ctx, diskBackupId)

		if err != nil {
			return nil, "", err
		}

		return object, helper.PString(object.LatestOperationState), nil
	}
}

func (me *LightHouseService) DescribeLighthouseDiskAttachmentById(ctx context.Context, diskId string) (diskAttachment *lighthouse.Disk, errRet error) {
	logId := getLogId(ctx)

	request := lighthouse.NewDescribeDisksRequest()
	request.DiskIds = []*string{&diskId}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseLighthouseClient().DescribeDisks(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.DiskSet) < 1 {
		return
	}

	diskAttachment = response.Response.DiskSet[0]
	return
}

func (me *LightHouseService) LighthouseDiskAttachmentStateRefreshFunc(diskId string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		ctx := contextNil

		object, err := me.DescribeLighthouseDiskAttachmentById(ctx, diskId)

		if err != nil {
			return nil, "", err
		}

		return object, helper.PString(object.DiskState), nil
	}
}

func (me *LightHouseService) DeleteLighthouseDiskAttachmentById(ctx context.Context, diskId string) (errRet error) {
	logId := getLogId(ctx)

	request := lighthouse.NewDetachDisksRequest()
	request.DiskIds = []*string{&diskId}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseLighthouseClient().DetachDisks(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *LightHouseService) DescribeLighthouseKeyPairById(ctx context.Context, keyId string) (keyPair *lighthouse.KeyPair, errRet error) {
	logId := getLogId(ctx)

	request := lighthouse.NewDescribeKeyPairsRequest()
	request.KeyIds = []*string{&keyId}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseLighthouseClient().DescribeKeyPairs(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.KeyPairSet) < 1 {
		return
	}

	keyPair = response.Response.KeyPairSet[0]
	return
}

func (me *LightHouseService) DeleteLighthouseKeyPairById(ctx context.Context, keyId string) (errRet error) {
	logId := getLogId(ctx)

	request := lighthouse.NewDeleteKeyPairsRequest()
	request.KeyIds = []*string{&keyId}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseLighthouseClient().DeleteKeyPairs(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *LightHouseService) DescribeLighthouseSnapshotById(ctx context.Context, snapshotId string) (snapshot *lighthouse.Snapshot, errRet error) {
	logId := getLogId(ctx)

	request := lighthouse.NewDescribeSnapshotsRequest()
	request.SnapshotIds = []*string{&snapshotId}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseLighthouseClient().DescribeSnapshots(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.SnapshotSet) < 1 {
		return
	}

	snapshot = response.Response.SnapshotSet[0]
	return
}

func (me *LightHouseService) DeleteLighthouseSnapshotById(ctx context.Context, snapshotId string) (errRet error) {
	logId := getLogId(ctx)

	request := lighthouse.NewDeleteSnapshotsRequest()
	request.SnapshotIds = []*string{&snapshotId}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseLighthouseClient().DeleteSnapshots(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *LightHouseService) LighthouseSnapshotStateRefreshFunc(snapshotId string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		ctx := contextNil

		object, err := me.DescribeLighthouseSnapshotById(ctx, snapshotId)

		if err != nil {
			return nil, "", err
		}

		return object, helper.PString(object.SnapshotState), nil
	}
}

func (me *LightHouseService) LighthouseApplySnapshotStateRefreshFunc(snapshotId string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		ctx := contextNil

		object, err := me.DescribeLighthouseSnapshotById(ctx, snapshotId)

		if err != nil {
			return nil, "", err
		}

		return object, helper.PString(object.LatestOperationState), nil
	}
}

func (me *LightHouseService) DescribeLighthouseBundleByFilter(ctx context.Context, param map[string]interface{}) (bundle []*lighthouse.Bundle, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = lighthouse.NewDescribeBundlesRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	var (
		offset = 0
		limit  = 20
	)

	for k, v := range param {
		if k == "bundle_ids" {
			request.BundleIds = helper.Strings(v.([]string))
		}
		if k == "offset" {
			offset = v.(int)
		}
		if k == "limit" {
			limit = v.(int)
		}
		if k == "zones" {
			request.Zones = helper.Strings(v.([]string))
		}
		if k == "filters" {
			request.Filters = v.([]*lighthouse.Filter)
		}
	}

	ratelimit.Check(request.GetAction())
	request.Offset = helper.IntInt64(offset)
	request.Limit = helper.IntInt64(limit)
	response, err := me.client.UseLighthouseClient().DescribeBundles(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || response.Response == nil || len(response.Response.BundleSet) < 1 {
		return
	}
	bundle = append(bundle, response.Response.BundleSet...)

	return
}

func (me *LightHouseService) DescribeLighthouseZoneByFilter(ctx context.Context, param map[string]interface{}) (zone []*lighthouse.ZoneInfo, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = lighthouse.NewDescribeZonesRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "order_field" {
			request.OrderField = helper.String(v.(string))
		}
		if k == "order" {
			request.Order = helper.String(v.(string))
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseLighthouseClient().DescribeZones(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || len(response.Response.ZoneInfoSet) < 1 {
		errRet = fmt.Errorf("Response body is null")
		return
	}
	zone = append(zone, response.Response.ZoneInfoSet...)

	return
}

func (me *LightHouseService) DescribeLighthouseSceneByFilter(ctx context.Context, param map[string]interface{}) (scene []*lighthouse.Scene, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = lighthouse.NewDescribeScenesRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "scene_ids" {
			request.SceneIds = helper.Strings(v.([]string))
		}
		if k == "offset" {
			request.Offset = helper.IntInt64(v.(int))
		}
		if k == "limit" {
			request.Limit = helper.IntInt64(v.(int))
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
		response, err := me.client.UseLighthouseClient().DescribeScenes(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.SceneSet) < 1 {
			break
		}
		scene = append(scene, response.Response.SceneSet...)
		if len(response.Response.SceneSet) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *LightHouseService) DescribeLighthouseResetInstanceBlueprintByFilter(ctx context.Context, param map[string]interface{}) (resetInstanceBlueprint []*lighthouse.ResetInstanceBlueprint, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = lighthouse.NewDescribeResetInstanceBlueprintsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	var (
		offset int = 0
		limit  int = 20
	)
	for k, v := range param {
		if k == "instance_id" {
			request.InstanceId = helper.String(v.(string))
		}
		if k == "offset" {
			offset = v.(int)
		}
		if k == "limit" {
			limit = v.(int)
		}
		if k == "filters" {
			request.Filters = v.([]*lighthouse.Filter)
		}
	}

	ratelimit.Check(request.GetAction())

	request.Offset = helper.IntInt64(offset)
	request.Limit = helper.IntInt64(limit)
	response, err := me.client.UseLighthouseClient().DescribeResetInstanceBlueprints(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || len(response.Response.ResetInstanceBlueprintSet) < 1 {
		errRet = fmt.Errorf("Response body is null")
		return
	}
	resetInstanceBlueprint = append(resetInstanceBlueprint, response.Response.ResetInstanceBlueprintSet...)

	return
}

func (me *LightHouseService) DescribeLighthouseRegionByFilter(ctx context.Context, param map[string]interface{}) (region []*lighthouse.RegionInfo, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = lighthouse.NewDescribeRegionsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseLighthouseClient().DescribeRegions(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || len(response.Response.RegionSet) < 1 {
		errRet = fmt.Errorf("Response is null")
	}
	region = append(region, response.Response.RegionSet...)

	return
}

func (me *LightHouseService) DescribeLighthouseInstanceVncUrlByFilter(ctx context.Context, instanceId string) (instanceVncUrl string, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = lighthouse.NewDescribeInstanceVncUrlRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	request.InstanceId = helper.String(instanceId)

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseLighthouseClient().DescribeInstanceVncUrl(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || response.Response == nil || response.Response.InstanceVncUrl == nil {
		errRet = fmt.Errorf("Response is null")
	}
	instanceVncUrl = *response.Response.InstanceVncUrl

	return
}

func (me *LightHouseService) DescribeLighthouseInstanceTrafficPackageByFilter(ctx context.Context, param map[string]interface{}) (instanceTrafficPackage []*lighthouse.InstanceTrafficPackage, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = lighthouse.NewDescribeInstancesTrafficPackagesRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	var (
		offset = 0
		limit  = 20
	)

	for k, v := range param {
		if k == "instance_ids" {
			request.InstanceIds = helper.Strings(v.([]string))
		}
		if k == "offset" {
			offset = v.(int)
		}
		if k == "limit" {
			limit = v.(int)
		}
		if k == "instance_ids" {
			request.InstanceIds = helper.Strings(v.([]string))
		}
	}

	ratelimit.Check(request.GetAction())

	for {
		request.Offset = helper.IntInt64(offset)
		request.Limit = helper.IntInt64(limit)
		response, err := me.client.UseLighthouseClient().DescribeInstancesTrafficPackages(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.InstanceTrafficPackageSet) < 1 {
			break
		}
		instanceTrafficPackage = append(instanceTrafficPackage, response.Response.InstanceTrafficPackageSet...)
		if len(response.Response.InstanceTrafficPackageSet) < limit {
			break
		}

		offset += limit
	}

	return
}

func (me *LightHouseService) DescribeLighthouseInstanceDiskNum(ctx context.Context, instanceIds []string) (instanceDiskNum []*lighthouse.AttachDetail, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = lighthouse.NewDescribeInstancesDiskNumRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	request.InstanceIds = helper.Strings(instanceIds)

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseLighthouseClient().DescribeInstancesDiskNum(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || len(response.Response.AttachDetailSet) < 1 {
		errRet = fmt.Errorf("Response is null")
		return
	}
	instanceDiskNum = append(instanceDiskNum, response.Response.AttachDetailSet...)

	return
}

func (me *LightHouseService) DescribeLighthouseInstanceBlueprintByFilter(ctx context.Context, instanceIds []string) (instanceBlueprint []*lighthouse.BlueprintInstance, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = lighthouse.NewDescribeBlueprintInstancesRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	request.InstanceIds = helper.Strings(instanceIds)

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseLighthouseClient().DescribeBlueprintInstances(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || len(response.Response.BlueprintInstanceSet) < 1 {
		errRet = fmt.Errorf("Response is null")
		return
	}
	instanceBlueprint = append(instanceBlueprint, response.Response.BlueprintInstanceSet...)

	return
}

func (me *LightHouseService) DescribeLighthouseDiskConfigByFilter(ctx context.Context, param map[string]interface{}) (diskConfig []*lighthouse.DiskConfig, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = lighthouse.NewDescribeDiskConfigsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "filters" {
			request.Filters = v.([]*lighthouse.Filter)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseLighthouseClient().DescribeDiskConfigs(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || len(response.Response.DiskConfigSet) < 1 {
		errRet = fmt.Errorf("Response is null")
		return
	}
	diskConfig = append(diskConfig, response.Response.DiskConfigSet...)

	return
}
