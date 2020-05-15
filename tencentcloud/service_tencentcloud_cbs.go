package tencentcloud

import (
	"context"
	"log"

	cbs "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cbs/v20170312"
	"github.com/terraform-providers/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/terraform-providers/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/terraform-providers/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

type CbsService struct {
	client *connectivity.TencentCloudClient
}

func (me *CbsService) DescribeDiskById(ctx context.Context, diskId string) (disk *cbs.Disk, errRet error) {
	logId := getLogId(ctx)
	request := cbs.NewDescribeDisksRequest()
	request.DiskIds = []*string{&diskId}
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

	if len(response.Response.DiskSet) > 0 {
		disk = response.Response.DiskSet[0]
	}
	return
}

func (me *CbsService) DescribeDisksByFilter(ctx context.Context, params map[string]string) (disks []*cbs.Disk, errRet error) {
	logId := getLogId(ctx)
	request := cbs.NewDescribeDisksRequest()
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

func (me *CbsService) ModifyDiskAttributes(ctx context.Context, diskId, diskName string, projectId int) error {
	logId := getLogId(ctx)
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

func (me *CbsService) DeleteDiskById(ctx context.Context, diskId string) error {
	logId := getLogId(ctx)
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
	logId := getLogId(ctx)
	request := cbs.NewResizeDiskRequest()
	request.DiskId = &diskId
	request.DiskSize = helper.IntUint64(diskSize)
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseCbsClient().ResizeDisk(request)
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
	logId := getLogId(ctx)
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
	logId := getLogId(ctx)
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
	logId := getLogId(ctx)
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
	logId := getLogId(ctx)
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

	snapshotId = *response.Response.SnapshotId
	return
}

func (me *CbsService) DescribeSnapshotById(ctx context.Context, snapshotId string) (snapshot *cbs.Snapshot, errRet error) {
	logId := getLogId(ctx)
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

func (me *CbsService) DescribeSnapshotsByFilter(ctx context.Context, params map[string]string) (snapshots []*cbs.Snapshot, errRet error) {
	logId := getLogId(ctx)
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
	logId := getLogId(ctx)
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
	logId := getLogId(ctx)
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
	logId := getLogId(contextNil)
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
	logId := getLogId(contextNil)
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
	logId := getLogId(ctx)
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
	logId := getLogId(ctx)
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
	logId := getLogId(ctx)
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
	logId := getLogId(ctx)
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
	logId := getLogId(ctx)
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
	logId := getLogId(ctx)
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

func flattenCbsTagsMapping(tags []*cbs.Tag) (mapping map[string]string) {
	mapping = make(map[string]string)
	for _, tag := range tags {
		mapping[*tag.Key] = *tag.Value
	}
	return
}
