package cbs

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	cbs "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cbs/v20170312"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

func NewCbsService(client *connectivity.TencentCloudClient) CbsService {
	return CbsService{client: client}
}

type CbsService struct {
	client *connectivity.TencentCloudClient
}

func (me *CbsService) DescribeDiskSetByIds(ctx context.Context, diskSetIds string) (disks []*cbs.Disk, errRet error) {

	diskSet, err := helper.StrToStrList(diskSetIds)
	if err != nil {
		return
	}

	disks, err = me.DescribeDiskList(ctx, helper.StringsStringsPoint(diskSet))
	if err != nil {
		errRet = err
		return
	}

	if len(disks) < 1 {
		return
	}

	return
}

func (me *CbsService) DescribeDiskById(ctx context.Context, diskId string) (disk *cbs.Disk, errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := cbs.NewDescribeDisksRequest()
	request.DiskIds = common.StringPtrs([]string{diskId})
	request.Limit = helper.IntUint64(100)
	ratelimit.Check(request.GetAction())

	var iacExtInfo connectivity.IacExtInfo
	iacExtInfo.InstanceId = diskId
	response, err := me.client.UseCbsClient(iacExtInfo).DescribeDisks(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if err != nil {
		errRet = err
		return
	}

	if len(response.Response.DiskSet) > 0 {
		disk = response.Response.DiskSet[0]
	}

	return
}

func (me *CbsService) DescribeDiskList(ctx context.Context, diskIds []*string) (disk []*cbs.Disk, errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := cbs.NewDescribeDisksRequest()
	request.DiskIds = diskIds
	request.Limit = helper.IntUint64(100)
	ratelimit.Check(request.GetAction())

	var iacExtInfo connectivity.IacExtInfo
	tmpList := make([]string, len(diskIds))
	for k, v := range diskIds {
		tmpList[k] = *v
	}
	iacExtInfo.InstanceId = strings.Join(tmpList, tccommon.FILED_SP)
	response, err := me.client.UseCbsClient(iacExtInfo).DescribeDisks(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.DiskSet) > 0 {
		disk = response.Response.DiskSet
	}
	return
}

func (me *CbsService) DescribeDisksByFilter(ctx context.Context, params map[string]interface{}) (disks []*cbs.Disk, errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := cbs.NewDescribeDisksRequest()
	request.Filters = make([]*cbs.Filter, 0, len(params))
	for k, v := range params {
		filter := &cbs.Filter{
			Name: helper.String(k),
		}
		switch v := v.(type) {
		case string:
			filter.Values = []*string{helper.String(v)}
		case []*string:
			filter.Values = v
		}
		request.Filters = append(request.Filters, filter)
	}

	offset := 0
	pageSize := 100
	disks = make([]*cbs.Disk, 0)
	for {
		request.Offset = helper.IntUint64(offset)
		request.Limit = helper.IntUint64(pageSize)
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseCbsClient().DescribeDisks(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.DiskSet) < 1 {
			break
		}

		disks = append(disks, response.Response.DiskSet...)

		if len(response.Response.DiskSet) < pageSize {
			break
		}
		offset += pageSize
	}
	return
}

func (me *CbsService) DescribeDisksInParallelByFilter(ctx context.Context, params map[string]interface{}) (disks []*cbs.Disk, errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := cbs.NewDescribeDisksRequest()

	request.Filters = make([]*cbs.Filter, 0, len(params))
	for k, v := range params {
		filter := &cbs.Filter{
			Name: helper.String(k),
		}
		switch v := v.(type) {
		case string:
			filter.Values = []*string{helper.String(v)}
		case []*string:
			filter.Values = v
		}
		request.Filters = append(request.Filters, filter)
	}
	response, err := me.client.UseCbsClient().DescribeDisks(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}

	if response == nil || len(response.Response.DiskSet) < 1 {
		return
	}

	total := response.Response.TotalCount

	var limit = 100
	num := int(*total) / limit
	g := tccommon.NewGoRoutine(num + 1)
	wg := sync.WaitGroup{}

	var diskSetList = make([]interface{}, num+1)

	for i := 0; i <= num; i++ {
		wg.Add(1)
		value := i
		goFunc := func() {
			offset := value * limit
			request := cbs.NewDescribeDisksRequest()
			request.Filters = make([]*cbs.Filter, 0, len(params))
			for k, v := range params {
				filter := &cbs.Filter{
					Name: helper.String(k),
				}
				switch v := v.(type) {
				case string:
					filter.Values = []*string{helper.String(v)}
				case []*string:
					filter.Values = v
				}
				request.Filters = append(request.Filters, filter)
			}

			request.Offset = helper.IntUint64(offset)
			request.Limit = helper.IntUint64(limit)

			ratelimit.Check(request.GetAction())
			response, err := me.client.UseCbsClient().DescribeDisks(request)
			if err != nil {
				log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
					logId, request.GetAction(), request.ToJsonString(), err.Error())
				errRet = err
				return
			}
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

			diskSetList[value] = response.Response.DiskSet

			wg.Done()
		}
		g.Run(goFunc)
	}
	wg.Wait()

	for _, v := range diskSetList {
		disks = append(disks, v.([]*cbs.Disk)...)
	}

	return
}

func (me *CbsService) ModifyDiskAttributes(ctx context.Context, diskId, diskName string, projectId int) error {
	logId := tccommon.GetLogId(ctx)
	request := cbs.NewModifyDiskAttributesRequest()
	request.DiskIds = []*string{&diskId}
	if diskName != "" {
		request.DiskName = &diskName
	}
	if projectId >= 0 {
		request.ProjectId = helper.IntUint64(projectId)
	}
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseCbsClient().ModifyDiskAttributes(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return nil
}

func (me *CbsService) DeleteDiskSetByIds(ctx context.Context, diskSetIds string) error {
	logId := tccommon.GetLogId(ctx)
	request := cbs.NewTerminateDisksRequest()

	diskSet, err := helper.StrToStrList(diskSetIds)
	if err != nil {
		return err
	}

	request.DiskIds = helper.StringsStringsPoint(diskSet)
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseCbsClient().TerminateDisks(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	return nil
}

func (me *CbsService) DeleteDiskById(ctx context.Context, diskId string) error {
	logId := tccommon.GetLogId(ctx)
	request := cbs.NewTerminateDisksRequest()
	request.DiskIds = []*string{&diskId}
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseCbsClient().TerminateDisks(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	return nil
}

func (me *CbsService) ResizeDisk(ctx context.Context, diskId string, diskSize int) error {
	logId := tccommon.GetLogId(ctx)
	request := cbs.NewResizeDiskRequest()
	request.DiskId = &diskId
	request.DiskSize = helper.IntUint64(diskSize)
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseCbsClient().ResizeDisk(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		storage, e := me.DescribeDiskById(ctx, diskId)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if *storage.DiskState == CBS_STORAGE_STATUS_EXPANDING {
			return resource.RetryableError(fmt.Errorf("cbs storage status is %s", *storage.DiskState))
		}

		if *storage.DiskSize != uint64(diskSize) {
			return resource.RetryableError(fmt.Errorf("waiting for cbs size changed to %d, now %d", diskSize, *storage.DiskSize))
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s resize cbs failed, reason:%s\n ", logId, err.Error())
		return err
	}

	return nil
}

func (me *CbsService) ModifyThroughputPerformance(ctx context.Context, diskId string, throughputPerformance int) error {
	logId := tccommon.GetLogId(ctx)
	request := cbs.NewModifyDiskExtraPerformanceRequest()
	request.DiskId = &diskId
	request.ThroughputPerformance = helper.IntUint64(throughputPerformance)
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseCbsClient().ModifyDiskExtraPerformance(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	return nil
}

func (me *CbsService) ApplySnapshot(ctx context.Context, diskId, snapshotId string) error {
	logId := tccommon.GetLogId(ctx)
	request := cbs.NewApplySnapshotRequest()
	request.DiskId = &diskId
	request.SnapshotId = &snapshotId
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseCbsClient().ApplySnapshot(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	return nil
}

func (me *CbsService) AttachDisk(ctx context.Context, diskId, instanceId string) error {
	logId := tccommon.GetLogId(ctx)
	request := cbs.NewAttachDisksRequest()
	request.DiskIds = []*string{&diskId}
	request.InstanceId = &instanceId
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseCbsClient().AttachDisks(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	return nil
}

func (me *CbsService) DetachDisk(ctx context.Context, diskId, instanceId string) error {
	logId := tccommon.GetLogId(ctx)
	request := cbs.NewDetachDisksRequest()
	request.DiskIds = []*string{&diskId}
	request.InstanceId = &instanceId
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseCbsClient().DetachDisks(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	return nil
}

func (me *CbsService) CreateSnapshot(ctx context.Context, diskId, snapshotName string) (snapshotId string, errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := cbs.NewCreateSnapshotRequest()
	request.DiskId = &diskId
	request.SnapshotName = &snapshotName
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseCbsClient().CreateSnapshot(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || response.Response == nil || response.Response.SnapshotId == nil {
		errRet = fmt.Errorf("CreateSnapshot response is nil.")
		return
	}

	snapshotId = *response.Response.SnapshotId
	return
}

func (me *CbsService) DescribeSnapshotById(ctx context.Context, snapshotId string) (snapshot *cbs.Snapshot, errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := cbs.NewDescribeSnapshotsRequest()
	request.SnapshotIds = []*string{&snapshotId}
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseCbsClient().DescribeSnapshots(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.SnapshotSet) > 0 {
		snapshot = response.Response.SnapshotSet[0]
	}
	return
}

func (me *CbsService) DescribeSnapshotByIds(ctx context.Context, snapshotIdsParam []*string) (snapshots []*cbs.Snapshot, errRet error) {
	if len(snapshotIdsParam) == 0 {
		return
	}

	var (
		logId            = tccommon.GetLogId(ctx)
		request          = cbs.NewDescribeSnapshotsRequest()
		err              error
		response         *cbs.DescribeSnapshotsResponse
		offset, pageSize uint64 = 0, 100
	)
	request.SnapshotIds = snapshotIdsParam

	for {
		request.Offset = &offset
		request.Limit = &pageSize

		err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())
			response, err = me.client.UseCbsClient().DescribeSnapshots(request)
			if err != nil {
				return tccommon.RetryError(err, tccommon.InternalError)
			}
			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			errRet = err
			return
		}

		snapshots = append(snapshots, response.Response.SnapshotSet...)
		if len(response.Response.SnapshotSet) < int(pageSize) {
			break
		}
		offset += pageSize
	}
	return
}

func (me *CbsService) DescribeSnapshotsByFilter(ctx context.Context, params map[string]string) (snapshots []*cbs.Snapshot, errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := cbs.NewDescribeSnapshotsRequest()
	request.Filters = make([]*cbs.Filter, 0, len(params))
	for k, v := range params {
		filter := &cbs.Filter{
			Name:   helper.String(k),
			Values: []*string{helper.String(v)},
		}
		request.Filters = append(request.Filters, filter)
	}

	offset := 0
	pageSize := 100
	for {
		request.Offset = helper.IntUint64(offset)
		request.Limit = helper.IntUint64(pageSize)
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseCbsClient().DescribeSnapshots(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.SnapshotSet) < 1 {
			break
		}

		snapshots = append(snapshots, response.Response.SnapshotSet...)

		if len(response.Response.SnapshotSet) < pageSize {
			break
		}
		offset += pageSize
	}
	return
}

func (me *CbsService) ModifySnapshotName(ctx context.Context, snapshotId, snapshotName string) error {
	logId := tccommon.GetLogId(ctx)
	request := cbs.NewModifySnapshotAttributeRequest()
	request.SnapshotId = &snapshotId
	request.SnapshotName = &snapshotName
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseCbsClient().ModifySnapshotAttribute(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	return nil
}

func (me *CbsService) DeleteSnapshot(ctx context.Context, snapshotId string) error {
	logId := tccommon.GetLogId(ctx)
	request := cbs.NewDeleteSnapshotsRequest()
	request.SnapshotIds = []*string{&snapshotId}
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseCbsClient().DeleteSnapshots(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	return nil
}

func (me *CbsService) DescribeSnapshotPolicyById(ctx context.Context, policyId string) (policy *cbs.AutoSnapshotPolicy, errRet error) {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	request := cbs.NewDescribeAutoSnapshotPoliciesRequest()
	request.AutoSnapshotPolicyIds = []*string{&policyId}
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseCbsClient().DescribeAutoSnapshotPolicies(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.AutoSnapshotPolicySet) > 0 {
		policy = response.Response.AutoSnapshotPolicySet[0]
	}
	return
}

func (me *CbsService) DescribeSnapshotPolicy(ctx context.Context, policyId, policyName string) (policies []*cbs.AutoSnapshotPolicy, errRet error) {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	request := cbs.NewDescribeAutoSnapshotPoliciesRequest()
	request.Filters = make([]*cbs.Filter, 0)
	if policyId != "" {
		filter := cbs.Filter{
			Name:   helper.String("auto-snapshot-policy-id"),
			Values: []*string{&policyId},
		}
		request.Filters = append(request.Filters, &filter)
	}
	if policyName != "" {
		filter := cbs.Filter{
			Name:   helper.String("auto-snapshot-policy-name"),
			Values: []*string{&policyName},
		}
		request.Filters = append(request.Filters, &filter)
	}
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseCbsClient().DescribeAutoSnapshotPolicies(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	policies = response.Response.AutoSnapshotPolicySet
	return
}

func (me *CbsService) DeleteSnapshotPolicy(ctx context.Context, policyId string) error {
	logId := tccommon.GetLogId(ctx)
	request := cbs.NewDeleteAutoSnapshotPoliciesRequest()
	request.AutoSnapshotPolicyIds = []*string{&policyId}
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseCbsClient().DeleteAutoSnapshotPolicies(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	return nil
}

func (me *CbsService) AttachSnapshotPolicy(ctx context.Context, diskId, policyId string) error {
	logId := tccommon.GetLogId(ctx)
	request := cbs.NewBindAutoSnapshotPolicyRequest()
	request.AutoSnapshotPolicyId = &policyId
	request.DiskIds = []*string{&diskId}
	ratelimit.Check(request.GetAction())
	_, err := me.client.UseCbsClient().BindAutoSnapshotPolicy(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
	}
	return nil
}

func (me *CbsService) DescribeAttachedSnapshotPolicy(ctx context.Context, diskId, policyId string) (policy *cbs.AutoSnapshotPolicy, errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := cbs.NewDescribeDiskAssociatedAutoSnapshotPolicyRequest()
	request.DiskId = &diskId
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseCbsClient().DescribeDiskAssociatedAutoSnapshotPolicy(request)
	if err != nil {
		errRet = err
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return
	}
	for i, item := range response.Response.AutoSnapshotPolicySet {
		if *item.AutoSnapshotPolicyId == policyId {
			policy = response.Response.AutoSnapshotPolicySet[i]
			break
		}
	}
	return
}

func (me *CbsService) UnattachSnapshotPolicy(ctx context.Context, diskId, policyId string) error {
	logId := tccommon.GetLogId(ctx)
	request := cbs.NewUnbindAutoSnapshotPolicyRequest()
	request.AutoSnapshotPolicyId = &policyId
	request.DiskIds = []*string{&diskId}
	ratelimit.Check(request.GetAction())
	_, err := me.client.UseCbsClient().UnbindAutoSnapshotPolicy(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
	}
	return nil
}

func (me *CbsService) ModifyDiskChargeType(ctx context.Context, storageId string, chargeType string, renewFlag string, period int) error {
	logId := tccommon.GetLogId(ctx)
	request := cbs.NewModifyDisksChargeTypeRequest()
	request.DiskIds = []*string{&storageId}
	request.DiskChargePrepaid = &cbs.DiskChargePrepaid{Period: helper.IntUint64(period), RenewFlag: &renewFlag}
	ratelimit.Check(request.GetAction())
	_, err := me.client.UseCbsClient().ModifyDisksChargeType(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
	}
	return nil
}

func (me *CbsService) ModifyDisksRenewFlag(ctx context.Context, storageId string, renewFlag string) error {
	logId := tccommon.GetLogId(ctx)
	request := cbs.NewModifyDisksRenewFlagRequest()
	request.DiskIds = []*string{&storageId}
	request.RenewFlag = &renewFlag

	ratelimit.Check(request.GetAction())
	_, err := me.client.UseCbsClient().ModifyDisksRenewFlag(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
	}
	return nil
}

func (me *CbsService) DescribeCbsDiskBackupById(ctx context.Context, diskBackupId string) (DiskBackup *cbs.DiskBackup, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := cbs.NewDescribeDiskBackupsRequest()
	request.DiskBackupIds = []*string{&diskBackupId}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCbsClient().DescribeDiskBackups(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.DiskBackupSet) < 1 {
		return
	}

	DiskBackup = response.Response.DiskBackupSet[0]
	return
}

func (me *CbsService) DeleteCbsDiskBackupById(ctx context.Context, diskBackupId string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := cbs.NewDeleteDiskBackupsRequest()
	request.DiskBackupIds = []*string{&diskBackupId}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCbsClient().DeleteDiskBackups(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *CbsService) ModifyDiskBackupQuota(ctx context.Context, diskId string, diskBackupQuota int) (errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := cbs.NewModifyDiskBackupQuotaRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	request.DiskId = helper.String(diskId)
	request.DiskBackupQuota = helper.IntUint64(diskBackupQuota)

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCbsClient().ModifyDiskBackupQuota(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *CbsService) CreateDiskBackup(ctx context.Context, diskId, diskBackupName string) (diskBackupId string, errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := cbs.NewCreateDiskBackupRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()
	request.DiskId = helper.String(diskId)
	request.DiskBackupName = helper.String(diskBackupName)

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseCbsClient().CreateDiskBackup(request)
		if e != nil {
			if sdkError, ok := e.(*errors.TencentCloudSDKError); ok {
				if sdkError.Code == "ResourceUnavailable.NotSupported" {
					return resource.NonRetryableError(e)
				}
			}
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		diskBackupId = *result.Response.DiskBackupId
		return nil
	})
	if err != nil {
		errRet = err
		log.Printf("[CRITAL]%s create cbs DiskBackup failed, reason:%+v", logId, err)
		return
	}
	return
}

func (me *CbsService) DescribeCbsSnapshotSharePermissionById(ctx context.Context, snapshotId string) (snapshotSharePermissions []*cbs.SharePermission, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := cbs.NewDescribeSnapshotSharePermissionRequest()
	request.SnapshotId = &snapshotId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCbsClient().DescribeSnapshotSharePermission(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	snapshotSharePermissions = response.Response.SharePermissionSet
	return
}

func (me *CbsService) ModifySnapshotsSharePermission(ctx context.Context, snapshotId, permission string, accountIds []string) (errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := cbs.NewModifySnapshotsSharePermissionRequest()

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()
	request.SnapshotIds = []*string{&snapshotId}
	request.Permission = helper.String(permission)
	request.AccountIds = helper.StringsStringsPoint(accountIds)
	ratelimit.Check(request.GetAction())

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := me.client.UseCbsClient().ModifySnapshotsSharePermission(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create cbs SnapshotSharePermission failed, reason:%+v", logId, err)
		return err
	}
	return
}

func (me *CbsService) ApplyDiskBackup(ctx context.Context, diskBackupId, diskId string) (errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := cbs.NewApplyDiskBackupRequest()

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()
	request.DiskBackupId = helper.String(diskBackupId)
	request.DiskId = helper.String(diskId)
	ratelimit.Check(request.GetAction())

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := me.client.UseCbsClient().ApplyDiskBackup(request)
		if e != nil {
			if sdkError, ok := e.(*errors.TencentCloudSDKError); ok {
				if sdkError.Code == "ResourceUnavailable.NotSupported" {
					return resource.NonRetryableError(e)
				}
			}
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s ApplyDiskBackup failed, reason:%+v", logId, err)
		return err
	}
	return
}

func (me *CbsService) DescribeDiskConfigQuota(ctx context.Context, cvmInfo map[string]interface{}) (diskConfigSet []*cbs.DiskConfig, errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := cbs.NewDescribeDiskConfigQuotaRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	request.InquiryType = helper.String("INQUIRY_CVM_CONFIG")
	request.Zones = helper.Strings([]string{cvmInfo["availability_zone"].(string)})
	request.CPU = helper.Int64Uint64(cvmInfo["cpu_core_count"].(int64))
	request.Memory = helper.Int64Uint64(cvmInfo["memory_size"].(int64))
	request.InstanceFamilies = helper.Strings([]string{cvmInfo["family"].(string)})
	request.DiskTypes = helper.Strings(cvmInfo["disk_types"].([]string))
	request.DiskChargeType = helper.String(cvmInfo["disk_charge_type"].(string))
	request.DiskUsage = helper.String(cvmInfo["disk_usage"].(string))
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseCbsClient().DescribeDiskConfigQuota(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		diskConfigSet = result.Response.DiskConfigSet
		return nil
	})
	if err != nil {
		errRet = err
		log.Printf("[CRITAL]%s create cbs DiskBackup failed, reason:%+v", logId, err)
		return
	}
	return
}
