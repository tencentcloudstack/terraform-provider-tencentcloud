// Copyright (c) 2017-2025 Tencent. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package v20170312

const (
	// error codes for specific actions

	// Your account failed the qualification verification.
	ACCOUNTQUALIFICATIONRESTRICTIONS = "AccountQualificationRestrictions"

	// Role name authentication failed.
	AUTHFAILURE_CAMROLENAMEAUTHENTICATEFAILED = "AuthFailure.CamRoleNameAuthenticateFailed"

	// ENIs do not support changing subnets.
	ENINOTALLOWEDCHANGESUBNET = "EniNotAllowedChangeSubnet"

	// The account already exists.
	FAILEDOPERATION_ACCOUNTALREADYEXISTS = "FailedOperation.AccountAlreadyExists"

	// You cannot share images with yourself.
	FAILEDOPERATION_ACCOUNTISYOURSELF = "FailedOperation.AccountIsYourSelf"

	// The instance `ins-xxxxxxx` is already in the instance placement group `dgroup-xxxxxx`.
	FAILEDOPERATION_ALREADYINDISASTERRECOVERGROUP = "FailedOperation.AlreadyInDisasterRecoverGroup"

	// BYOL images cannot be shared.
	FAILEDOPERATION_BYOLIMAGESHAREFAILED = "FailedOperation.BYOLImageShareFailed"

	// The specified spread placement group does not exist.
	FAILEDOPERATION_DISASTERRECOVERGROUPNOTFOUND = "FailedOperation.DisasterRecoverGroupNotFound"

	// Failed to obtain the status of TencentCloud Automation Tools for the instance.
	FAILEDOPERATION_GETINSTANCETATAGENTSTATUSFAILED = "FailedOperation.GetInstanceTATAgentStatusFailed"

	// The tag key contains invalid characters.
	FAILEDOPERATION_ILLEGALTAGKEY = "FailedOperation.IllegalTagKey"

	// The tag value contains invalid characters.
	FAILEDOPERATION_ILLEGALTAGVALUE = "FailedOperation.IllegalTagValue"

	// Price query failed.
	FAILEDOPERATION_INQUIRYPRICEFAILED = "FailedOperation.InquiryPriceFailed"

	// Failed to query the refund: the payment order is not found. Check whether the instance `ins-xxxxxxx` has expired.
	FAILEDOPERATION_INQUIRYREFUNDPRICEFAILED = "FailedOperation.InquiryRefundPriceFailed"

	// The image is busy. Please try again later.
	FAILEDOPERATION_INVALIDIMAGESTATE = "FailedOperation.InvalidImageState"

	// The applicationRole instance does not support the operation.
	FAILEDOPERATION_INVALIDINSTANCEAPPLICATIONROLE = "FailedOperation.InvalidInstanceApplicationRole"

	// The EMR instance `ins-xxxxxxxx` does not support this operation.
	FAILEDOPERATION_INVALIDINSTANCEAPPLICATIONROLEEMR = "FailedOperation.InvalidInstanceApplicationRoleEmr"

	// No available IPs in the subnet.
	FAILEDOPERATION_NOAVAILABLEIPADDRESSCOUNTINSUBNET = "FailedOperation.NoAvailableIpAddressCountInSubnet"

	// This instance does not bind an EIP.
	FAILEDOPERATION_NOTFOUNDEIP = "FailedOperation.NotFoundEIP"

	// You’re using a collaborator account. Please enter a root account.
	FAILEDOPERATION_NOTMASTERACCOUNT = "FailedOperation.NotMasterAccount"

	// The specified placement group is not empty.
	FAILEDOPERATION_PLACEMENTSETNOTEMPTY = "FailedOperation.PlacementSetNotEmpty"

	// The configuration or billing mode of the CVM instances purchased during the promotion period cannot be modified.
	FAILEDOPERATION_PROMOTIONALPERIORESTRICTION = "FailedOperation.PromotionalPerioRestriction"

	// The service is not available in this country/region.
	FAILEDOPERATION_PROMOTIONALREGIONRESTRICTION = "FailedOperation.PromotionalRegionRestriction"

	// Image sharing failed.
	FAILEDOPERATION_QIMAGESHAREFAILED = "FailedOperation.QImageShareFailed"

	// Image sharing failed.
	FAILEDOPERATION_RIMAGESHAREFAILED = "FailedOperation.RImageShareFailed"

	// Security group operation failed.
	FAILEDOPERATION_SECURITYGROUPACTIONFAILED = "FailedOperation.SecurityGroupActionFailed"

	// The snapshot size is larger than the disk capacity. You need a larger disk space.
	FAILEDOPERATION_SNAPSHOTSIZELARGERTHANDATASIZE = "FailedOperation.SnapshotSizeLargerThanDataSize"

	// The snapshot size should be larger than the cloud disk capacity.
	FAILEDOPERATION_SNAPSHOTSIZELESSTHANDATASIZE = "FailedOperation.SnapshotSizeLessThanDataSize"

	// The tag key specified in the request is reserved for the system.
	FAILEDOPERATION_TAGKEYRESERVED = "FailedOperation.TagKeyReserved"

	// This image is not a Linux&x86_64 image.
	FAILEDOPERATION_TATAGENTNOTSUPPORT = "FailedOperation.TatAgentNotSupport"

	// The instance is unreturnable.
	FAILEDOPERATION_UNRETURNABLE = "FailedOperation.Unreturnable"

	// The image quota has been exceeded.
	IMAGEQUOTALIMITEXCEEDED = "ImageQuotaLimitExceeded"

	// You are trying to create more instances than your remaining quota allows.
	INSTANCESQUOTALIMITEXCEEDED = "InstancesQuotaLimitExceeded"

	// Internal error.
	INTERNALERROR = "InternalError"

	// Internal error.
	INTERNALERROR_TRADEUNKNOWNERROR = "InternalError.TradeUnknownError"

	// Internal error.
	INTERNALSERVERERROR = "InternalServerError"

	// Insufficient account balance.
	INVALIDACCOUNT_INSUFFICIENTBALANCE = "InvalidAccount.InsufficientBalance"

	// The account has unpaid orders.
	INVALIDACCOUNT_UNPAIDORDER = "InvalidAccount.UnpaidOrder"

	// Invalid account ID.
	INVALIDACCOUNTID_NOTFOUND = "InvalidAccountId.NotFound"

	// You cannot share images with yourself.
	INVALIDACCOUNTIS_YOURSELF = "InvalidAccountIs.YourSelf"

	// The specified ClientToken exceeds the maximum length of 64 bytes.
	INVALIDCLIENTTOKEN_TOOLONG = "InvalidClientToken.TooLong"

	// Invalid filter.
	INVALIDFILTER = "InvalidFilter"

	// The value of [`Filter`](https://www.tencentcloud.com/document/api/213/15753?from_cn_redirect=1#Filter) exceeds the limit.
	INVALIDFILTERVALUE_LIMITEXCEEDED = "InvalidFilterValue.LimitExceeded"

	// The specified operation on this CDH instance is not support .
	INVALIDHOST_NOTSUPPORTED = "InvalidHost.NotSupported"

	// Invalid [CDH](https://intl.cloud.tencent.com/document/product/416?from_cn_redirect=1) `ID`. The specified [CDH](https://intl.cloud.tencent.com/document/product/416?from_cn_redirect=1) `ID` has an invalid format. For example, `host-1122` has an invalid `ID` length.
	INVALIDHOSTID_MALFORMED = "InvalidHostId.Malformed"

	// The specified HostId does not exist, or does not belong to your account.
	INVALIDHOSTID_NOTFOUND = "InvalidHostId.NotFound"

	// The image is being shared.
	INVALIDIMAGEID_INSHARED = "InvalidImageId.InShared"

	// Invalid image status.
	INVALIDIMAGEID_INCORRECTSTATE = "InvalidImageId.IncorrectState"

	// Invalid image ID format.
	INVALIDIMAGEID_MALFORMED = "InvalidImageId.Malformed"

	// The image cannot be found.
	INVALIDIMAGEID_NOTFOUND = "InvalidImageId.NotFound"

	// The image size exceeds the limit.
	INVALIDIMAGEID_TOOLARGE = "InvalidImageId.TooLarge"

	// The specified image name already exists.
	INVALIDIMAGENAME_DUPLICATE = "InvalidImageName.Duplicate"

	// The operating system type is not supported.
	INVALIDIMAGEOSTYPE_UNSUPPORTED = "InvalidImageOsType.Unsupported"

	// The operating system version is not supported.
	INVALIDIMAGEOSVERSION_UNSUPPORTED = "InvalidImageOsVersion.Unsupported"

	// This instance is not supported.
	INVALIDINSTANCE_NOTSUPPORTED = "InvalidInstance.NotSupported"

	// Invalid instance `ID`. The specified instance `ID` has an invalid format. For example, `ins-1122` has an invalid `ID` length.
	INVALIDINSTANCEID_MALFORMED = "InvalidInstanceId.Malformed"

	// No instance found.
	INVALIDINSTANCEID_NOTFOUND = "InvalidInstanceId.NotFound"

	// The specified InstanceName exceeds the maximum length of 128 bytes.
	INVALIDINSTANCENAME_TOOLONG = "InvalidInstanceName.TooLong"

	// This instance does not meet the [Return Policy](https://intl.cloud.tencent.com/document/product/213/9711?from_cn_redirect=1) for prepaid instances.
	INVALIDINSTANCENOTSUPPORTEDPREPAIDINSTANCE = "InvalidInstanceNotSupportedPrepaidInstance"

	// This operation cannot be performed due to the current instance status.
	INVALIDINSTANCESTATE = "InvalidInstanceState"

	// The specified `InstanceType` parameter has an invalid format.
	INVALIDINSTANCETYPE_MALFORMED = "InvalidInstanceType.Malformed"

	// The number of key pairs exceeds the limit.
	INVALIDKEYPAIR_LIMITEXCEEDED = "InvalidKeyPair.LimitExceeded"

	// Invalid key pair ID. The specified key pair ID has an invalid format. For example, `skey-1122` has an invalid `ID` length.
	INVALIDKEYPAIRID_MALFORMED = "InvalidKeyPairId.Malformed"

	// Invalid key pair ID. The specified key pair ID does not exist.
	INVALIDKEYPAIRID_NOTFOUND = "InvalidKeyPairId.NotFound"

	// Key pair name already exists.
	INVALIDKEYPAIRNAME_DUPLICATE = "InvalidKeyPairName.Duplicate"

	// The key name cannot be empty.
	INVALIDKEYPAIRNAMEEMPTY = "InvalidKeyPairNameEmpty"

	// The key name contains invalid characters. Key names can only contain letters, numbers and underscores.
	INVALIDKEYPAIRNAMEINCLUDEILLEGALCHAR = "InvalidKeyPairNameIncludeIllegalChar"

	// The key name cannot exceed 25 characters.
	INVALIDKEYPAIRNAMETOOLONG = "InvalidKeyPairNameTooLong"

	// A parameter error occurred.
	INVALIDPARAMETER = "InvalidParameter"

	// Up to one parameter can be specified.
	INVALIDPARAMETER_ATMOSTONE = "InvalidParameter.AtMostOne"

	// Automatic snapshot creation is not supported.
	INVALIDPARAMETER_AUTOSNAPSHOTNOTSUPPORTED = "InvalidParameter.AutoSnapshotNotSupported"

	// The parameter CdcId is not supported.
	INVALIDPARAMETER_CDCNOTSUPPORTED = "InvalidParameter.CdcNotSupported"

	// RootDisk ID should not be passed to DataDiskIds.
	INVALIDPARAMETER_DATADISKIDCONTAINSROOTDISK = "InvalidParameter.DataDiskIdContainsRootDisk"

	// The specified data disk does not belong to the specified instance.
	INVALIDPARAMETER_DATADISKNOTBELONGSPECIFIEDINSTANCE = "InvalidParameter.DataDiskNotBelongSpecifiedInstance"

	// Only one system disk snapshot can be included.
	INVALIDPARAMETER_DUPLICATESYSTEMSNAPSHOTS = "InvalidParameter.DuplicateSystemSnapshots"

	// When specifying the CTCC/CUCC/CMCC public IP address parameter for edge zones, you need to first specify the public IP address parameter for the primary IP address.
	INVALIDPARAMETER_EDGEZONEMISSINTERNETACCESSIBLE = "InvalidParameter.EdgeZoneMissInternetAccessible"

	// The specified CDH host does not support custom instance specifications.
	INVALIDPARAMETER_HOSTIDCUSTOMIZEDINSTANCETYPENOTSUPPORT = "InvalidParameter.HostIdCustomizedInstanceTypeNotSupport"

	// The specified CDH host does not support the instance model specifications.
	INVALIDPARAMETER_HOSTIDINSTANCETYPENOTSUPPORT = "InvalidParameter.HostIdInstanceTypeNotSupport"

	// The specified CDH host does not support standard instance specifications.
	INVALIDPARAMETER_HOSTIDSTANDARDINSTANCETYPENOTSUPPORT = "InvalidParameter.HostIdStandardInstanceTypeNotSupport"

	// This operation is not supported under the current status of the CVM.
	INVALIDPARAMETER_HOSTIDSTATUSNOTSUPPORT = "InvalidParameter.HostIdStatusNotSupport"

	// The specified HostName is invalid.
	INVALIDPARAMETER_HOSTNAMEILLEGAL = "InvalidParameter.HostNameIllegal"

	// Either `ImageIds` or `SnapshotIds` must be specified.
	INVALIDPARAMETER_IMAGEIDSSNAPSHOTIDSMUSTONE = "InvalidParameter.ImageIdsSnapshotIdsMustOne"

	// This API does not support instance images.
	INVALIDPARAMETER_INSTANCEIMAGENOTSUPPORT = "InvalidParameter.InstanceImageNotSupport"

	// No CDH host supports the specified instance specifications.
	INVALIDPARAMETER_INSTANCETYPESUPPORTEDHOSTNOTFOUND = "InvalidParameter.InstanceTypeSupportedHostNotFound"

	// Unable to set the public network bandwidth. 
	INVALIDPARAMETER_INTERNETACCESSIBLENOTSUPPORTED = "InvalidParameter.InternetAccessibleNotSupported"

	// Invalid parameter dependency.
	INVALIDPARAMETER_INVALIDDEPENDENCE = "InvalidParameter.InvalidDependence"

	// Invalid VPC IP address format.
	INVALIDPARAMETER_INVALIDIPFORMAT = "InvalidParameter.InvalidIpFormat"

	// The specified KMS key ID is invalid.
	INVALIDPARAMETER_INVALIDKMSKEYID = "InvalidParameter.InvalidKmsKeyId"

	// `ImageIds` and `Filters` cannot be specified at the same time.
	INVALIDPARAMETER_INVALIDPARAMETERCOEXISTIMAGEIDSFILTERS = "InvalidParameter.InvalidParameterCoexistImageIdsFilters"

	// Invalid URL.
	INVALIDPARAMETER_INVALIDPARAMETERURLERROR = "InvalidParameter.InvalidParameterUrlError"

	// The entered TargetOSType is invalid.
	INVALIDPARAMETER_INVALIDTARGETOSTYPE = "InvalidParameter.InvalidTargetOSType"

	// `CoreCount` and `ThreadPerCore` must be specified at the same time.
	INVALIDPARAMETER_LACKCORECOUNTORTHREADPERCORE = "InvalidParameter.LackCoreCountOrThreadPerCore"

	// Local data disks cannot be used to create instance images.
	INVALIDPARAMETER_LOCALDATADISKNOTSUPPORT = "InvalidParameter.LocalDataDiskNotSupport"

	// Only edge zones support this parameter.
	INVALIDPARAMETER_ONLYSUPPORTFOREDGEZONE = "InvalidParameter.OnlySupportForEdgeZone"

	// Specifying an SSH key will override the original one of the image.
	INVALIDPARAMETER_PARAMETERCONFLICT = "InvalidParameter.ParameterConflict"

	// Setting login password is not supported.
	INVALIDPARAMETER_PASSWORDNOTSUPPORTED = "InvalidParameter.PasswordNotSupported"

	// The specified snapshot does not exist.
	INVALIDPARAMETER_SNAPSHOTNOTFOUND = "InvalidParameter.SnapshotNotFound"

	// This parameter can only be used when the allowlist feature is enabled.
	INVALIDPARAMETER_SPECIALPARAMETERFORSPECIALACCOUNT = "InvalidParameter.SpecialParameterForSpecialAccount"

	// At least one of the multiple parameters must be passed in.
	INVALIDPARAMETER_SPECIFYONEPARAMETER = "InvalidParameter.SpecifyOneParameter"

	// Swap disks are not supported.
	INVALIDPARAMETER_SWAPDISKNOTSUPPORT = "InvalidParameter.SwapDiskNotSupport"

	// The parameter does not contain system disk snapshot.
	INVALIDPARAMETER_SYSTEMSNAPSHOTNOTFOUND = "InvalidParameter.SystemSnapshotNotFound"

	// The length of parameter exceeds the limit.
	INVALIDPARAMETER_VALUETOOLARGE = "InvalidParameter.ValueTooLarge"

	// The parameter combination is invalid.
	INVALIDPARAMETERCOMBINATION = "InvalidParameterCombination"

	// The two specified parameters conflict. An EIP can only be bound to the instance or the specified private IP of the specified ENI.
	INVALIDPARAMETERCONFLICT = "InvalidParameterConflict"

	// Incorrect parameter value.
	INVALIDPARAMETERVALUE = "InvalidParameterValue"

	// The number of request parameters are not equal.
	INVALIDPARAMETERVALUE_AMOUNTNOTEQUAL = "InvalidParameterValue.AmountNotEqual"

	// The shared bandwidth package ID is invalid. Please provide a standard shared bandwidth package ID in the format similar to bwp-xxxxxxxx. In this format, the letter x stands for a lowercase character or a number.
	INVALIDPARAMETERVALUE_BANDWIDTHPACKAGEIDMALFORMED = "InvalidParameterValue.BandwidthPackageIdMalformed"

	// The specified bandwidth package does not exist.
	INVALIDPARAMETERVALUE_BANDWIDTHPACKAGEIDNOTFOUND = "InvalidParameterValue.BandwidthPackageIdNotFound"

	// The ISP of the bandwidth package does not match the ISP parameter.
	INVALIDPARAMETERVALUE_BANDWIDTHPACKAGEISPNOTMATCH = "InvalidParameterValue.BandwidthPackageIspNotMatch"

	// The availability zone of the bandwidth package does not match the specified availability zone.
	INVALIDPARAMETERVALUE_BANDWIDTHPACKAGEZONENOTMATCH = "InvalidParameterValue.BandwidthPackageZoneNotMatch"

	// Only VPC is supported. The network type of the instance is classic network, which cannot be changed.
	INVALIDPARAMETERVALUE_BASICNETWORKINSTANCEFAMILY = "InvalidParameterValue.BasicNetworkInstanceFamily"

	// The bucket does not exist.
	INVALIDPARAMETERVALUE_BUCKETNOTFOUND = "InvalidParameterValue.BucketNotFound"

	// Invalid `CamRoleName`. This parameter must contain only letters, numbers and symbols (`+`, `=`, `,`, `.`, `@`, `_`, `-`).
	INVALIDPARAMETERVALUE_CAMROLENAMEMALFORMED = "InvalidParameterValue.CamRoleNameMalformed"

	// CDH disk expansion only supports LOCAL_BASIC and LOCAL_SSD.
	INVALIDPARAMETERVALUE_CDHONLYLOCALDATADISKRESIZE = "InvalidParameterValue.CdhOnlyLocalDataDiskResize"

	// Corresponding CHC hosts not found.
	INVALIDPARAMETERVALUE_CHCHOSTSNOTFOUND = "InvalidParameterValue.ChcHostsNotFound"

	// No network is configured for this CHC.
	INVALIDPARAMETERVALUE_CHCNETWORKEMPTY = "InvalidParameterValue.ChcNetworkEmpty"

	// The minimum capacity of a SSD data disk is 100 GB.
	INVALIDPARAMETERVALUE_CLOUDSSDDATADISKSIZETOOSMALL = "InvalidParameterValue.CloudSsdDataDiskSizeTooSmall"

	// Invalid number of cores.
	INVALIDPARAMETERVALUE_CORECOUNTVALUE = "InvalidParameterValue.CoreCountValue"

	// CDC does not support the specified billing mode.
	INVALIDPARAMETERVALUE_DEDICATEDCLUSTERNOTSUPPORTEDCHARGETYPE = "InvalidParameterValue.DedicatedClusterNotSupportedChargeType"

	// A deployment VPC already exists.
	INVALIDPARAMETERVALUE_DEPLOYVPCALREADYEXISTS = "InvalidParameterValue.DeployVpcAlreadyExists"

	// Incorrect placement group ID format.
	INVALIDPARAMETERVALUE_DISASTERRECOVERGROUPIDMALFORMED = "InvalidParameterValue.DisasterRecoverGroupIdMalformed"

	// Duplicate parameter value.
	INVALIDPARAMETERVALUE_DUPLICATE = "InvalidParameterValue.Duplicate"

	// Duplicate tags.
	INVALIDPARAMETERVALUE_DUPLICATETAGS = "InvalidParameterValue.DuplicateTags"

	// ENI data does not exist.
	INVALIDPARAMETERVALUE_ELASTICNETWORKNOTEXIST = "InvalidParameterValue.ElasticNetworkNotExist"

	// The eni data vpc subnet is mismatched. it must be in the same vpc but different subnets.
	INVALIDPARAMETERVALUE_ELASTICNETWORKVPCSUBNETMISMATCH = "InvalidParameterValue.ElasticNetworkVpcSubnetMismatch"

	// The number of requested public IP addresses exceeds the quota of this instance type.
	INVALIDPARAMETERVALUE_EXTERNALIPQUOTALIMITED = "InvalidParameterValue.ExternalIpQuotaLimited"

	// Non-GPU instances cannot be changed to the GPU instance.
	INVALIDPARAMETERVALUE_GPUINSTANCEFAMILY = "InvalidParameterValue.GPUInstanceFamily"

	// Your High-Performance Computing (HPC) cluster is already bound to another Availability Zone, so you cannot purchase machines in the current Availability Zone.
	INVALIDPARAMETERVALUE_HPCCLUSTERIDZONEIDNOTMATCH = "InvalidParameterValue.HpcClusterIdZoneIdNotMatch"

	// Invalid IP format
	INVALIDPARAMETERVALUE_IPADDRESSMALFORMED = "InvalidParameterValue.IPAddressMalformed"

	// Invalid IPv6 address
	INVALIDPARAMETERVALUE_IPV6ADDRESSMALFORMED = "InvalidParameterValue.IPv6AddressMalformed"

	// ISO files must be imported by force.
	INVALIDPARAMETERVALUE_ISOMUSTIMPORTBYFORCE = "InvalidParameterValue.ISOMustImportByForce"

	// The value of HostName is invalid.
	INVALIDPARAMETERVALUE_ILLEGALHOSTNAME = "InvalidParameterValue.IllegalHostName"

	// Incorrect request parameter format.
	INVALIDPARAMETERVALUE_INCORRECTFORMAT = "InvalidParameterValue.IncorrectFormat"

	// Invalid instance ID. Please enter a valid ID, such as ins-xxxxxxxx (“x” represents a lower-case letter or a number).
	INVALIDPARAMETERVALUE_INSTANCEIDMALFORMED = "InvalidParameterValue.InstanceIdMalformed"

	// Operation not supported for instances with different billing modes.
	INVALIDPARAMETERVALUE_INSTANCENOTSUPPORTEDMIXPRICINGMODEL = "InvalidParameterValue.InstanceNotSupportedMixPricingModel"

	// The specified instance type does not exist.
	INVALIDPARAMETERVALUE_INSTANCETYPENOTFOUND = "InvalidParameterValue.InstanceTypeNotFound"

	// The instance type does not support eni data.
	INVALIDPARAMETERVALUE_INSTANCETYPENOTSUPPORTELASTICNETWORKS = "InvalidParameterValue.InstanceTypeNotSupportElasticNetworks"

	// This type of instances cannot be added to the HPC cluster.
	INVALIDPARAMETERVALUE_INSTANCETYPENOTSUPPORTHPCCLUSTER = "InvalidParameterValue.InstanceTypeNotSupportHpcCluster"

	// The HPC cluster needs to be specified for the high-performance computing instance.
	INVALIDPARAMETERVALUE_INSTANCETYPEREQUIREDHPCCLUSTER = "InvalidParameterValue.InstanceTypeRequiredHpcCluster"

	// The spot instances are out of stock.
	INVALIDPARAMETERVALUE_INSUFFICIENTOFFERING = "InvalidParameterValue.InsufficientOffering"

	// The bid is lower than the market price.
	INVALIDPARAMETERVALUE_INSUFFICIENTPRICE = "InvalidParameterValue.InsufficientPrice"

	// Invalid AppID
	INVALIDPARAMETERVALUE_INVALIDAPPIDFORMAT = "InvalidParameterValue.InvalidAppIdFormat"

	// Unsupported boot mode.
	INVALIDPARAMETERVALUE_INVALIDBOOTMODE = "InvalidParameterValue.InvalidBootMode"

	// You don’t have the write permission to the bucket.
	INVALIDPARAMETERVALUE_INVALIDBUCKETPERMISSIONFOREXPORT = "InvalidParameterValue.InvalidBucketPermissionForExport"

	// The length of `FileNamePrefixList` does not match `ImageIds` or `SnapshotIds`.
	INVALIDPARAMETERVALUE_INVALIDFILENAMEPREFIXLIST = "InvalidParameterValue.InvalidFileNamePrefixList"

	// Converting to a non-GPU or other type of GPU instance is not supported.
	INVALIDPARAMETERVALUE_INVALIDGPUFAMILYCHANGE = "InvalidParameterValue.InvalidGPUFamilyChange"

	// Invalid format of image family name
	INVALIDPARAMETERVALUE_INVALIDIMAGEFAMILY = "InvalidParameterValue.InvalidImageFamily"

	// The specified image does not support the specified instance type.
	INVALIDPARAMETERVALUE_INVALIDIMAGEFORGIVENINSTANCETYPE = "InvalidParameterValue.InvalidImageForGivenInstanceType"

	// A RAW image cannot be used to create a CVM. Choose another image.
	INVALIDPARAMETERVALUE_INVALIDIMAGEFORMAT = "InvalidParameterValue.InvalidImageFormat"

	// The image does not support this operation.
	INVALIDPARAMETERVALUE_INVALIDIMAGEID = "InvalidParameterValue.InvalidImageId"

	// The image cannot be used to reinstall the current instance.
	INVALIDPARAMETERVALUE_INVALIDIMAGEIDFORRETSETINSTANCE = "InvalidParameterValue.InvalidImageIdForRetsetInstance"

	// The specified image ID is a shared image.
	INVALIDPARAMETERVALUE_INVALIDIMAGEIDISSHARED = "InvalidParameterValue.InvalidImageIdIsShared"

	// The operating system of the specified image is not available in the current region.
	INVALIDPARAMETERVALUE_INVALIDIMAGEOSNAME = "InvalidParameterValue.InvalidImageOsName"

	// The image has another ongoing task. Please check and try again later.
	INVALIDPARAMETERVALUE_INVALIDIMAGESTATE = "InvalidParameterValue.InvalidImageState"

	// The instance configuration is upgraded for free and cannot be downgraded within 3 months.
	INVALIDPARAMETERVALUE_INVALIDINSTANCESOURCE = "InvalidParameterValue.InvalidInstanceSource"

	// The model does not support the support cycle contract.
	INVALIDPARAMETERVALUE_INVALIDINSTANCETYPEPERIODICCONTRACT = "InvalidParameterValue.InvalidInstanceTypePeriodicContract"

	// The specified instance type does not support exclusive sales payment mode.
	INVALIDPARAMETERVALUE_INVALIDINSTANCETYPEUNDERWRITE = "InvalidParameterValue.InvalidInstanceTypeUnderwrite"

	// Invalid IP address.
	INVALIDPARAMETERVALUE_INVALIDIPFORMAT = "InvalidParameterValue.InvalidIpFormat"

	// Instance boot template description format is incorrect.
	INVALIDPARAMETERVALUE_INVALIDLAUNCHTEMPLATEDESCRIPTION = "InvalidParameterValue.InvalidLaunchTemplateDescription"

	// Incorrect format of instance launch template name.
	INVALIDPARAMETERVALUE_INVALIDLAUNCHTEMPLATENAME = "InvalidParameterValue.InvalidLaunchTemplateName"

	// Incorrect format of instance launch template version description.
	INVALIDPARAMETERVALUE_INVALIDLAUNCHTEMPLATEVERSIONDESCRIPTION = "InvalidParameterValue.InvalidLaunchTemplateVersionDescription"

	// Invalid license type.
	INVALIDPARAMETERVALUE_INVALIDLICENSETYPE = "InvalidParameterValue.InvalidLicenseType"

	// The specified eni id is not a vrdma network interface card.
	INVALIDPARAMETERVALUE_INVALIDNETWORKINTERFACEID = "InvalidParameterValue.InvalidNetworkInterfaceId"

	// The value of parameter MinCount must be less than InstanceCount.
	INVALIDPARAMETERVALUE_INVALIDPARAMETERMINCOUNT = "InvalidParameterValue.InvalidParameterMinCount"

	// Invalid parameter value.
	INVALIDPARAMETERVALUE_INVALIDPARAMETERVALUELIMIT = "InvalidParameterValue.InvalidParameterValueLimit"

	// Invalid password. The specified password does not meet the complexity requirements (e.g., too long or too short)
	INVALIDPARAMETERVALUE_INVALIDPASSWORD = "InvalidParameterValue.InvalidPassword"

	// The Region ID is unavailable.
	INVALIDPARAMETERVALUE_INVALIDREGION = "InvalidParameterValue.InvalidRegion"

	// Incorrect time format.
	INVALIDPARAMETERVALUE_INVALIDTIMEFORMAT = "InvalidParameterValue.InvalidTimeFormat"

	// Incorrect UserData format. Use the Base64-encoded format.
	INVALIDPARAMETERVALUE_INVALIDUSERDATAFORMAT = "InvalidParameterValue.InvalidUserDataFormat"

	// Invalid fuzzy query string
	INVALIDPARAMETERVALUE_INVALIDVAGUENAME = "InvalidParameterValue.InvalidVagueName"

	// This special VpcId or SubnetId is not found in the elastic network data structure.
	INVALIDPARAMETERVALUE_INVALIDVPCIDSUBNETIDNOTFOUND = "InvalidParameterValue.InvalidVpcIdSubnetIdNotFound"

	// Edge zones do not support this ISP.
	INVALIDPARAMETERVALUE_ISPNOTSUPPORTFOREDGEZONE = "InvalidParameterValue.IspNotSupportForEdgeZone"

	// Duplicate ISP parameter value specified.
	INVALIDPARAMETERVALUE_ISPVALUEREPEATED = "InvalidParameterValue.IspValueRepeated"

	// The key does not exist.
	INVALIDPARAMETERVALUE_KEYPAIRNOTFOUND = "InvalidParameterValue.KeyPairNotFound"

	// The specified key does not support the operation.
	INVALIDPARAMETERVALUE_KEYPAIRNOTSUPPORTED = "InvalidParameterValue.KeyPairNotSupported"

	// The default launch template version cannot be operated.
	INVALIDPARAMETERVALUE_LAUNCHTEMPLATEDEFAULTVERSION = "InvalidParameterValue.LaunchTemplateDefaultVersion"

	// Incorrect format of instance launch template ID. Please provide a valid instance launch template ID, similar to lt-xxxxxxxx, where x represents a lowercase character or digit.
	INVALIDPARAMETERVALUE_LAUNCHTEMPLATEIDMALFORMED = "InvalidParameterValue.LaunchTemplateIdMalformed"

	// The instance launch template ID does not exist.
	INVALIDPARAMETERVALUE_LAUNCHTEMPLATEIDNOTEXISTED = "InvalidParameterValue.LaunchTemplateIdNotExisted"

	// The combination of the template and the version ID does not exist.
	INVALIDPARAMETERVALUE_LAUNCHTEMPLATEIDVERNOTEXISTED = "InvalidParameterValue.LaunchTemplateIdVerNotExisted"

	// The specified instance launch template ID doesn't exist.
	INVALIDPARAMETERVALUE_LAUNCHTEMPLATEIDVERSETALREADY = "InvalidParameterValue.LaunchTemplateIdVerSetAlready"

	// The instance launch template is not found.
	INVALIDPARAMETERVALUE_LAUNCHTEMPLATENOTFOUND = "InvalidParameterValue.LaunchTemplateNotFound"

	// Invalid instance launch template version number.
	INVALIDPARAMETERVALUE_LAUNCHTEMPLATEVERSION = "InvalidParameterValue.LaunchTemplateVersion"

	// The number of parameter values exceeds the limit.
	INVALIDPARAMETERVALUE_LIMITEXCEEDED = "InvalidParameterValue.LimitExceeded"

	// The parameter value must be a DHCP-enabled VPC.
	INVALIDPARAMETERVALUE_MUSTDHCPENABLEDVPC = "InvalidParameterValue.MustDhcpEnabledVpc"

	// The parameter value must enable the elastic rdma api.
	INVALIDPARAMETERVALUE_MUSTENABLEDISRDMA = "InvalidParameterValue.MustEnabledIsRdma"

	// The subnet is not in the CDC cluster.
	INVALIDPARAMETERVALUE_NOTCDCSUBNET = "InvalidParameterValue.NotCdcSubnet"

	// The parameter value is required.
	INVALIDPARAMETERVALUE_NOTEMPTY = "InvalidParameterValue.NotEmpty"

	// Unsupported operation.
	INVALIDPARAMETERVALUE_NOTSUPPORTED = "InvalidParameterValue.NotSupported"

	// Preheating is not supported on this model.
	INVALIDPARAMETERVALUE_PREHEATNOTSUPPORTEDINSTANCETYPE = "InvalidParameterValue.PreheatNotSupportedInstanceType"

	// Preheating is not supported in this availability zone.
	INVALIDPARAMETERVALUE_PREHEATNOTSUPPORTEDZONE = "InvalidParameterValue.PreheatNotSupportedZone"

	// The pre-warming region is unavailable. Please check if the pre-warming region is correct.
	INVALIDPARAMETERVALUE_PREHEATUNAVAILABLEZONES = "InvalidParameterValue.PreheatUnavailableZones"

	//  Invalid parameter value: invalid parameter value range.
	INVALIDPARAMETERVALUE_RANGE = "InvalidParameterValue.Range"

	// The request requires a regional image.
	INVALIDPARAMETERVALUE_REQUIREDLOCATIONIMAGE = "InvalidParameterValue.RequiredLocationImage"

	// Invalid snapshot ID. Provide a snapshot ID in the format of snap-xxxxxxxx, where the letter x refers to lowercase letter or number.
	INVALIDPARAMETERVALUE_SNAPSHOTIDMALFORMED = "InvalidParameterValue.SnapshotIdMalformed"

	// Invalid subnet ID. Please provide a subnet ID in the format of subnet-xxxxxxxx, where “x” can be a lowercase letter or number.
	INVALIDPARAMETERVALUE_SUBNETIDMALFORMED = "InvalidParameterValue.SubnetIdMalformed"

	// The subnet ID availability zone does not match the instance location.
	INVALIDPARAMETERVALUE_SUBNETIDZONEIDNOTMATCH = "InvalidParameterValue.SubnetIdZoneIdNotMatch"

	// Creation failed: the subnet does not exist. Please specify another subnet.
	INVALIDPARAMETERVALUE_SUBNETNOTEXIST = "InvalidParameterValue.SubnetNotExist"

	// The specified tag does not exist
	INVALIDPARAMETERVALUE_TAGKEYNOTFOUND = "InvalidParameterValue.TagKeyNotFound"

	// Tag quota limit exceeded.
	INVALIDPARAMETERVALUE_TAGQUOTALIMITEXCEEDED = "InvalidParameterValue.TagQuotaLimitExceeded"

	// Invalid thread count per core.
	INVALIDPARAMETERVALUE_THREADPERCOREVALUE = "InvalidParameterValue.ThreadPerCoreValue"

	// The parameter value exceeds the maximum limit.
	INVALIDPARAMETERVALUE_TOOLARGE = "InvalidParameterValue.TooLarge"

	// Invalid parameter value: parameter value is too long.
	INVALIDPARAMETERVALUE_TOOLONG = "InvalidParameterValue.TooLong"

	// Invalid UUID.
	INVALIDPARAMETERVALUE_UUIDMALFORMED = "InvalidParameterValue.UuidMalformed"

	// The VPC ID `xxx` is invalid. Please provide a VPC ID in the format of vpc-xxxxxxxx, where “x” can be a lowercase letter or number.
	INVALIDPARAMETERVALUE_VPCIDMALFORMED = "InvalidParameterValue.VpcIdMalformed"

	// The specified VpcId doesn't exist.
	INVALIDPARAMETERVALUE_VPCIDNOTEXIST = "InvalidParameterValue.VpcIdNotExist"

	// The specified VpcId and SubnetId do not match.
	INVALIDPARAMETERVALUE_VPCIDSUBNETIDNOTMATCH = "InvalidParameterValue.VpcIdSubnetIdNotMatch"

	// The VPC and instance must be in the same availability zone.
	INVALIDPARAMETERVALUE_VPCIDZONEIDNOTMATCH = "InvalidParameterValue.VpcIdZoneIdNotMatch"

	// This VPC does not support the IPv6 addresses.
	INVALIDPARAMETERVALUE_VPCNOTSUPPORTIPV6ADDRESS = "InvalidParameterValue.VpcNotSupportIpv6Address"

	// The availability zone does not support this operation.
	INVALIDPARAMETERVALUE_ZONENOTSUPPORTED = "InvalidParameterValue.ZoneNotSupported"

	// The number of parameter values exceeds the limit.
	INVALIDPARAMETERVALUELIMIT = "InvalidParameterValueLimit"

	// Invalid parameter value: invalid `Offset`.
	INVALIDPARAMETERVALUEOFFSET = "InvalidParameterValueOffset"

	// Invalid password. The specified password does not meet the password requirements. For example, the password length does not meet the requirements.
	INVALIDPASSWORD = "InvalidPassword"

	// Invalid period. Valid values: [1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 24, 36]; unit: month.
	INVALIDPERIOD = "InvalidPeriod"

	// This operation is not supported for the account.
	INVALIDPERMISSION = "InvalidPermission"

	// Invalid project ID: the specified project ID does not exist.
	INVALIDPROJECTID_NOTFOUND = "InvalidProjectId.NotFound"

	// Invalid public key: the specified key already exists.
	INVALIDPUBLICKEY_DUPLICATE = "InvalidPublicKey.Duplicate"

	// Invalid public key: the specified public key does not meet the `OpenSSH RSA` format requirements.
	INVALIDPUBLICKEY_MALFORMED = "InvalidPublicKey.Malformed"

	// The region cannot be found.
	INVALIDREGION_NOTFOUND = "InvalidRegion.NotFound"

	// Currently this region does not support image synchronization.
	INVALIDREGION_UNAVAILABLE = "InvalidRegion.Unavailable"

	// The specified `security group ID` does not exist.
	INVALIDSECURITYGROUPID_NOTFOUND = "InvalidSecurityGroupId.NotFound"

	// The specified `security group ID` is in the wrong format. For example, `sg-ide32` has an invalid `instance ID` length.
	INVALIDSGID_MALFORMED = "InvalidSgId.Malformed"

	// The specified `zone` does not exist.
	INVALIDZONE_MISMATCHREGION = "InvalidZone.MismatchRegion"

	// An instance can be bound with up to 5 security groups.
	LIMITEXCEEDED_ASSOCIATEUSGLIMITEXCEEDED = "LimitExceeded.AssociateUSGLimitExceeded"

	// The quota limit for purchasing instances has been reached.
	LIMITEXCEEDED_CVMINSTANCEQUOTA = "LimitExceeded.CvmInstanceQuota"

	// The CVM ENIs associated with the security group has exceeded the limit.
	LIMITEXCEEDED_CVMSVIFSPERSECGROUPLIMITEXCEEDED = "LimitExceeded.CvmsVifsPerSecGroupLimitExceeded"

	// The quota of the specified placement group is insufficient.
	LIMITEXCEEDED_DISASTERRECOVERGROUP = "LimitExceeded.DisasterRecoverGroup"

	// The number of EIPs of an ENI contained in a specific instance has exceeded the maximum allowed EIPs of the target instance type. Please delete some EIPs and try again.
	LIMITEXCEEDED_EIPNUMLIMIT = "LimitExceeded.EipNumLimit"

	// The number of network interfaces exceeds the maximum limit for the instance.
	LIMITEXCEEDED_ENILIMITINSTANCETYPE = "LimitExceeded.EniLimitInstanceType"

	// The number of ENIs on a specified instance exceeds the maximum ENIs allowed for the target instance type. Delete some ENIs and try again.
	LIMITEXCEEDED_ENINUMLIMIT = "LimitExceeded.EniNumLimit"

	// The number of image export tasks in progress reached the upper limit. Please try again after the running tasks are completed.
	LIMITEXCEEDED_EXPORTIMAGETASKLIMITEXCEEDED = "LimitExceeded.ExportImageTaskLimitExceeded"

	// Number of IPs on this ENI reached the upper limit.
	LIMITEXCEEDED_IPV6ADDRESSNUM = "LimitExceeded.IPv6AddressNum"

	// Reached the upper limit of the ENIs for the instance.
	LIMITEXCEEDED_INSTANCEENINUMLIMIT = "LimitExceeded.InstanceEniNumLimit"

	// You are short of the instance quota.
	LIMITEXCEEDED_INSTANCEQUOTA = "LimitExceeded.InstanceQuota"

	// Unable to adjust: the target instance type does not support the configured public network bandwidth cap. See [Public Network Bandwidth Cap](https://intl.cloud.tencent.com/document/product/213/12523).
	LIMITEXCEEDED_INSTANCETYPEBANDWIDTH = "LimitExceeded.InstanceTypeBandwidth"

	// The number of instance launch templates exceeds the limit.
	LIMITEXCEEDED_LAUNCHTEMPLATEQUOTA = "LimitExceeded.LaunchTemplateQuota"

	// The number of instance launch template versions exceeds the limit.
	LIMITEXCEEDED_LAUNCHTEMPLATEVERSIONQUOTA = "LimitExceeded.LaunchTemplateVersionQuota"

	// Reached the upper limit for preheating in this availability zone. Remove some snapshots first.
	LIMITEXCEEDED_PREHEATIMAGESNAPSHOTOUTOFQUOTA = "LimitExceeded.PreheatImageSnapshotOutOfQuota"

	// Your quota for monthly-subscribed instances is used up. Increase your quota and try again.
	LIMITEXCEEDED_PREPAYQUOTA = "LimitExceeded.PrepayQuota"

	// The purchased quantity of committed instances has reached the maximum quota.
	LIMITEXCEEDED_PREPAYUNDERWRITEQUOTA = "LimitExceeded.PrepayUnderwriteQuota"

	// The number of security groups exceeds the quota limit.
	LIMITEXCEEDED_SINGLEUSGQUOTA = "LimitExceeded.SingleUSGQuota"

	// The spot instance offerings are out of stock.
	LIMITEXCEEDED_SPOTQUOTA = "LimitExceeded.SpotQuota"

	// Exceeded the upper limit of resources bound to the tag.
	LIMITEXCEEDED_TAGRESOURCEQUOTA = "LimitExceeded.TagResourceQuota"

	// Failed to return instances because of the quota limit.
	LIMITEXCEEDED_USERRETURNQUOTA = "LimitExceeded.UserReturnQuota"

	// You are short of the spot instance quota.
	LIMITEXCEEDED_USERSPOTQUOTA = "LimitExceeded.UserSpotQuota"

	// Insufficient subnet IPs.
	LIMITEXCEEDED_VPCSUBNETNUM = "LimitExceeded.VpcSubnetNum"

	// Missing parameter.
	MISSINGPARAMETER = "MissingParameter"

	// Parameter missing. Provide at least one parameter.
	MISSINGPARAMETER_ATLEASTONE = "MissingParameter.AtLeastOne"

	// The DPDK instance requires a VPC.
	MISSINGPARAMETER_DPDKINSTANCETYPEREQUIREDVPC = "MissingParameter.DPDKInstanceTypeRequiredVPC"

	// The instance type must have Cloud Monitor enabled.
	MISSINGPARAMETER_MONITORSERVICE = "MissingParameter.MonitorService"

	// An identical job is running.
	MUTEXOPERATION_TASKRUNNING = "MutexOperation.TaskRunning"

	// Operation not supported for this account.
	OPERATIONDENIED_ACCOUNTNOTSUPPORTED = "OperationDenied.AccountNotSupported"

	// A CHC instance without network configured is not allowed for the installation of a cloud image
	OPERATIONDENIED_CHCINSTALLCLOUDIMAGEWITHOUTDEPLOYNETWORK = "OperationDenied.ChcInstallCloudImageWithoutDeployNetwork"

	// Operation denied: This is a restricted account.
	OPERATIONDENIED_INNERUSERPROHIBITACTION = "OperationDenied.InnerUserProhibitAction"

	// The instance has an operation in progress. Please try again later.
	OPERATIONDENIED_INSTANCEOPERATIONINPROGRESS = "OperationDenied.InstanceOperationInProgress"

	// Bill-by-CVM users are not allowed to apply for edge zone public IP addresses.
	OPERATIONDENIED_NOTBANDWIDTHSHIFTUPUSERAPPLYEDGEZONEEIP = "OperationDenied.NotBandwidthShiftUpUserApplyEdgeZoneEip"

	// The number of shared images exceeds the quota.
	OVERQUOTA = "OverQuota"

	// This region does not support importing images.
	REGIONABILITYLIMIT_UNSUPPORTEDTOIMPORTIMAGE = "RegionAbilityLimit.UnsupportedToImportImage"

	// The resource is in use.
	RESOURCEINUSE = "ResourceInUse"

	// The disk rollback is in progress. Please try again later.
	RESOURCEINUSE_DISKROLLBACKING = "ResourceInUse.DiskRollbacking"

	// The availability zone has been sold out.
	RESOURCEINSUFFICIENT_AVAILABILITYZONESOLDOUT = "ResourceInsufficient.AvailabilityZoneSoldOut"

	// Insufficient subnet resources.
	RESOURCEINSUFFICIENT_CIDRBLOCK = "ResourceInsufficient.CidrBlock"

	// The specified cloud disk has been sold out.
	RESOURCEINSUFFICIENT_CLOUDDISKSOLDOUT = "ResourceInsufficient.CloudDiskSoldOut"

	// The parameters of cloud disk do not meet the specification.
	RESOURCEINSUFFICIENT_CLOUDDISKUNAVAILABLE = "ResourceInsufficient.CloudDiskUnavailable"

	// The number of instances exceeded the quota limit of spread placement groups.
	RESOURCEINSUFFICIENT_DISASTERRECOVERGROUPCVMQUOTA = "ResourceInsufficient.DisasterRecoverGroupCvmQuota"

	// Insufficient security group quota.
	RESOURCEINSUFFICIENT_INSUFFICIENTGROUPQUOTA = "ResourceInsufficient.InsufficientGroupQuota"

	// Inventory fails to satisfy the minimum purchasable quantity.
	RESOURCEINSUFFICIENT_INSUFFICIENTOFFERINGMINIMUM = "ResourceInsufficient.InsufficientOfferingMinimum"

	// The specified instance type is insufficient.
	RESOURCEINSUFFICIENT_SPECIFIEDINSTANCETYPE = "ResourceInsufficient.SpecifiedInstanceType"

	// The instances of the specified types were sold out in the selected AZs.
	RESOURCEINSUFFICIENT_ZONESOLDOUTFORSPECIFIEDINSTANCE = "ResourceInsufficient.ZoneSoldOutForSpecifiedInstance"

	// The HPC cluster does not exist.
	RESOURCENOTFOUND_HPCCLUSTER = "ResourceNotFound.HpcCluster"

	// The specified placement group does not exist.
	RESOURCENOTFOUND_INVALIDPLACEMENTSET = "ResourceNotFound.InvalidPlacementSet"

	// This instance type is not supported in the AZ.
	RESOURCENOTFOUND_INVALIDZONEINSTANCETYPE = "ResourceNotFound.InvalidZoneInstanceType"

	// The specified key pair does not exist.
	RESOURCENOTFOUND_KEYPAIRNOTFOUND = "ResourceNotFound.KeyPairNotFound"

	// No default CBS resources are available.
	RESOURCENOTFOUND_NODEFAULTCBS = "ResourceNotFound.NoDefaultCbs"

	// No default CBS resources are available.
	RESOURCENOTFOUND_NODEFAULTCBSWITHREASON = "ResourceNotFound.NoDefaultCbsWithReason"

	// Resources are unavailable.
	RESOURCEUNAVAILABLE = "ResourceUnavailable"

	// This instance type is unavailable in the availability zone.
	RESOURCEUNAVAILABLE_INSTANCETYPE = "ResourceUnavailable.InstanceType"

	// The snapshot is being created.
	RESOURCEUNAVAILABLE_SNAPSHOTCREATING = "ResourceUnavailable.SnapshotCreating"

	// Resources in this availability zone has been sold out.
	RESOURCESSOLDOUT_AVAILABLEZONE = "ResourcesSoldOut.AvailableZone"

	// The public IP has been sold out.
	RESOURCESSOLDOUT_EIPINSUFFICIENT = "ResourcesSoldOut.EipInsufficient"

	// The specified instance type is sold out.
	RESOURCESSOLDOUT_SPECIFIEDINSTANCETYPE = "ResourcesSoldOut.SpecifiedInstanceType"

	// A general error occurred during the security group service API call.
	SECGROUPACTIONFAILURE = "SecGroupActionFailure"

	// Unauthorized operation.
	UNAUTHORIZEDOPERATION = "UnauthorizedOperation"

	// The specified image does not belong to the user.
	UNAUTHORIZEDOPERATION_IMAGENOTBELONGTOACCOUNT = "UnauthorizedOperation.ImageNotBelongToAccount"

	// Check if the token is valid.
	UNAUTHORIZEDOPERATION_INVALIDTOKEN = "UnauthorizedOperation.InvalidToken"

	// Unauthorized operation. Make sure Multi-Factor Authentication (MFA) is valid.
	UNAUTHORIZEDOPERATION_MFAEXPIRED = "UnauthorizedOperation.MFAExpired"

	// Unauthorized operation. Make sure Multi-Factor Authentication (MFA) exists.
	UNAUTHORIZEDOPERATION_MFANOTFOUND = "UnauthorizedOperation.MFANotFound"

	// You’re not authorized for the operation. Please check the CAM policy.
	UNAUTHORIZEDOPERATION_PERMISSIONDENIED = "UnauthorizedOperation.PermissionDenied"

	// Unknown parameter error.
	UNKNOWNPARAMETER = "UnknownParameter"

	// Unsupported operation.
	UNSUPPORTEDOPERATION = "UnsupportedOperation"

	// This operation is currently not supported for ARM machines.
	UNSUPPORTEDOPERATION_ARMARCHITECTURE = "UnsupportedOperation.ArmArchitecture"

	// The specified instance or network cannot use the bandwidth package.
	UNSUPPORTEDOPERATION_BANDWIDTHPACKAGEIDNOTSUPPORTED = "UnsupportedOperation.BandwidthPackageIdNotSupported"

	// The far end ssd disk does not support this operation.
	UNSUPPORTEDOPERATION_CBSREMOTESSDNOTSUPPORT = "UnsupportedOperation.CbsRemoteSsdNotSupport"

	// Commercial image instance use cannot be adjusted for payment mode.
	UNSUPPORTEDOPERATION_COMMERCIALIMAGECHANGECHARGETYPE = "UnsupportedOperation.CommercialImageChangeChargeType"

	// Only one snapshot can be created in 24 hours. 
	UNSUPPORTEDOPERATION_DISKSNAPCREATETIMETOOOLD = "UnsupportedOperation.DiskSnapCreateTimeTooOld"

	// Edge Zone instances do not support this operation.
	UNSUPPORTEDOPERATION_EDGEZONEINSTANCE = "UnsupportedOperation.EdgeZoneInstance"

	// The selected edge zone does not support cloud disk operations.
	UNSUPPORTEDOPERATION_EDGEZONENOTSUPPORTCLOUDDISK = "UnsupportedOperation.EdgeZoneNotSupportCloudDisk"

	// An ENI is bound to the CVM. Please unbind the ENI from the CVM before switching to VPC.
	UNSUPPORTEDOPERATION_ELASTICNETWORKINTERFACE = "UnsupportedOperation.ElasticNetworkInterface"

	// Encrypted images are not supported.
	UNSUPPORTEDOPERATION_ENCRYPTEDIMAGESNOTSUPPORTED = "UnsupportedOperation.EncryptedImagesNotSupported"

	// You cannot change the model of a heterogeneous instance.
	UNSUPPORTEDOPERATION_HETEROGENEOUSCHANGEINSTANCEFAMILY = "UnsupportedOperation.HeterogeneousChangeInstanceFamily"

	// Instances with hibernation disabled are not supported.
	UNSUPPORTEDOPERATION_HIBERNATIONFORNORMALINSTANCE = "UnsupportedOperation.HibernationForNormalInstance"

	// The current image does not support hibernation.
	UNSUPPORTEDOPERATION_HIBERNATIONOSVERSION = "UnsupportedOperation.HibernationOsVersion"

	// IPv6 instances cannot be migrated from Classiclink to VPC.
	UNSUPPORTEDOPERATION_IPV6NOTSUPPORTVPCMIGRATE = "UnsupportedOperation.IPv6NotSupportVpcMigrate"

	// Failed to export the image: The image is too large.
	UNSUPPORTEDOPERATION_IMAGETOOLARGEEXPORTUNSUPPORTED = "UnsupportedOperation.ImageTooLargeExportUnsupported"

	// This instance billing mode does not support the operation.
	UNSUPPORTEDOPERATION_INSTANCECHARGETYPE = "UnsupportedOperation.InstanceChargeType"

	// A mixed payment mode is not supported.
	UNSUPPORTEDOPERATION_INSTANCEMIXEDPRICINGMODEL = "UnsupportedOperation.InstanceMixedPricingModel"

	// 
	UNSUPPORTEDOPERATION_INSTANCEMIXEDRESETINSTANCETYPE = "UnsupportedOperation.InstanceMixedResetInstanceType"

	// Central AZ and edge zone instances cannot be mixed in batch operation.
	UNSUPPORTEDOPERATION_INSTANCEMIXEDZONETYPE = "UnsupportedOperation.InstanceMixedZoneType"

	// The specified instance does not support operating system switching.
	UNSUPPORTEDOPERATION_INSTANCEOSCONVERTOSNOTSUPPORT = "UnsupportedOperation.InstanceOsConvertOsNotSupport"

	// The instance `ins-xxxxxx` with the `Xserver windows2012cndatacenterx86_64` operating system does not support this operation.
	UNSUPPORTEDOPERATION_INSTANCEOSWINDOWS = "UnsupportedOperation.InstanceOsWindows"

	// The current instance is in a failed OS reinstallation state and does not support this operation. We recommend reinstalling the operating system again. Alternatively, you may terminate/return the instance or submit a support ticket.
	UNSUPPORTEDOPERATION_INSTANCEREINSTALLFAILED = "UnsupportedOperation.InstanceReinstallFailed"

	// This CVM is blocked. Please submit a ticket.
	UNSUPPORTEDOPERATION_INSTANCESTATEBANNING = "UnsupportedOperation.InstanceStateBanning"

	// The instances are permanently corrupted, and this operation is not supported.
	UNSUPPORTEDOPERATION_INSTANCESTATECORRUPTED = "UnsupportedOperation.InstanceStateCorrupted"

	// Instances are entering the rescue mode, and this operation is not supported.
	UNSUPPORTEDOPERATION_INSTANCESTATEENTERRESCUEMODE = "UnsupportedOperation.InstanceStateEnterRescueMode"

	// The instance `ins-xxxxxx` in the `ENTER_SERVICE_LIVE_MIGRATE` status is not supported.
	UNSUPPORTEDOPERATION_INSTANCESTATEENTERSERVICELIVEMIGRATE = "UnsupportedOperation.InstanceStateEnterServiceLiveMigrate"

	// Instances are exiting from the rescue mode, and this operation is not supported.
	UNSUPPORTEDOPERATION_INSTANCESTATEEXITRESCUEMODE = "UnsupportedOperation.InstanceStateExitRescueMode"

	// The instance `ins-xxxxxx` in the `EXIT_SERVICE_LIVE_MIGRATE` status is not supported.
	UNSUPPORTEDOPERATION_INSTANCESTATEEXITSERVICELIVEMIGRATE = "UnsupportedOperation.InstanceStateExitServiceLiveMigrate"

	// The operation is not supported for the frozen instances.
	UNSUPPORTEDOPERATION_INSTANCESTATEFREEZING = "UnsupportedOperation.InstanceStateFreezing"

	// Unable to isolate: the instance is isolated
	UNSUPPORTEDOPERATION_INSTANCESTATEISOLATING = "UnsupportedOperation.InstanceStateIsolating"

	// The instance is failed to create, so the operation is not supported.
	UNSUPPORTEDOPERATION_INSTANCESTATELAUNCHFAILED = "UnsupportedOperation.InstanceStateLaunchFailed"

	// The specified operation is not supported for instances that are not in the running state.
	UNSUPPORTEDOPERATION_INSTANCESTATENOTRUNNING = "UnsupportedOperation.InstanceStateNotRunning"

	// The instances are being created, and this operation is not supported.
	UNSUPPORTEDOPERATION_INSTANCESTATEPENDING = "UnsupportedOperation.InstanceStatePending"

	// The instances are being restarted, and this operation is not supported.
	UNSUPPORTEDOPERATION_INSTANCESTATEREBOOTING = "UnsupportedOperation.InstanceStateRebooting"

	// Instances in the rescue mode are not available for this operation.
	UNSUPPORTEDOPERATION_INSTANCESTATERESCUEMODE = "UnsupportedOperation.InstanceStateRescueMode"

	// Running instances do not support this operation.
	UNSUPPORTEDOPERATION_INSTANCESTATERUNNING = "UnsupportedOperation.InstanceStateRunning"

	// The instances are being migrated, and this operation is not supported.
	UNSUPPORTEDOPERATION_INSTANCESTATESERVICELIVEMIGRATE = "UnsupportedOperation.InstanceStateServiceLiveMigrate"

	// Isolated instances do not support this operation.
	UNSUPPORTEDOPERATION_INSTANCESTATESHUTDOWN = "UnsupportedOperation.InstanceStateShutdown"

	// The instance is starting up, and this operation is not supported.
	UNSUPPORTEDOPERATION_INSTANCESTATESTARTING = "UnsupportedOperation.InstanceStateStarting"

	// The instance has been shut down, and this operation is not supported.
	UNSUPPORTEDOPERATION_INSTANCESTATESTOPPED = "UnsupportedOperation.InstanceStateStopped"

	// The instance is being shut down, and this operation is not supported.
	UNSUPPORTEDOPERATION_INSTANCESTATESTOPPING = "UnsupportedOperation.InstanceStateStopping"

	// Terminated instances are not supported.
	UNSUPPORTEDOPERATION_INSTANCESTATETERMINATED = "UnsupportedOperation.InstanceStateTerminated"

	// The instance is being terminated, and the operation is not supported.
	UNSUPPORTEDOPERATION_INSTANCESTATETERMINATING = "UnsupportedOperation.InstanceStateTerminating"

	// The instance type does not support setting the `Confidentiality` status.
	UNSUPPORTEDOPERATION_INSTANCETYPENOTSUPPORTCONFIDENTIALITY = "UnsupportedOperation.InstanceTypeNotSupportConfidentiality"

	// The instance type does not support setting the `GridDriverService` status.
	UNSUPPORTEDOPERATION_INSTANCETYPENOTSUPPORTGRIDLICENCE = "UnsupportedOperation.InstanceTypeNotSupportGridLicence"

	// The instance type does not support setting the HighDensityMode status.
	UNSUPPORTEDOPERATION_INSTANCETYPENOTSUPPORTHIGHDENSITYMODESETTING = "UnsupportedOperation.InstanceTypeNotSupportHighDensityModeSetting"

	// The instance type does not support setting the `EnableJumboFrame` status.
	UNSUPPORTEDOPERATION_INSTANCETYPENOTSUPPORTJUMBOFRAME = "UnsupportedOperation.InstanceTypeNotSupportJumboFrame"

	// Modifying Jumbo Frame status without a restart is not supported.
	UNSUPPORTEDOPERATION_INSTANCESENABLEJUMBOWITHOUTREBOOT = "UnsupportedOperation.InstancesEnableJumboWithoutReboot"

	// The instance is under termination protection and cannot be terminated. Disable the termination protection and try again.
	UNSUPPORTEDOPERATION_INSTANCESPROTECTED = "UnsupportedOperation.InstancesProtected"

	// Adjusting the data disk is not supported.
	UNSUPPORTEDOPERATION_INVALIDDATADISK = "UnsupportedOperation.InvalidDataDisk"

	// The specified disk is not supported.
	UNSUPPORTEDOPERATION_INVALIDDISK = "UnsupportedOperation.InvalidDisk"

	// Cloud block storage does not support backup points.
	UNSUPPORTEDOPERATION_INVALIDDISKBACKUPQUOTA = "UnsupportedOperation.InvalidDiskBackupQuota"

	// Fast rollback is not supported.
	UNSUPPORTEDOPERATION_INVALIDDISKFASTROLLBACK = "UnsupportedOperation.InvalidDiskFastRollback"

	// The image license type does not match the instance. Select another image.
	UNSUPPORTEDOPERATION_INVALIDIMAGELICENSETYPEFORRESET = "UnsupportedOperation.InvalidImageLicenseTypeForReset"

	// This operation is not supported for the instance with a termination schedule. Please cancel the scheduled termination time in the instance details page and try again.
	UNSUPPORTEDOPERATION_INVALIDINSTANCENOTSUPPORTEDPROTECTEDINSTANCE = "UnsupportedOperation.InvalidInstanceNotSupportedProtectedInstance"

	// Instances with swap disks are not supported.
	UNSUPPORTEDOPERATION_INVALIDINSTANCEWITHSWAPDISK = "UnsupportedOperation.InvalidInstanceWithSwapDisk"

	// The user does not have permissions to operate the current instance.
	UNSUPPORTEDOPERATION_INVALIDINSTANCESOWNER = "UnsupportedOperation.InvalidInstancesOwner"

	// The current operation is supported only for Tencent Cloud users.
	UNSUPPORTEDOPERATION_INVALIDPERMISSIONNONINTERNATIONALACCOUNT = "UnsupportedOperation.InvalidPermissionNonInternationalAccount"

	// Encrypted disks are not available in the selected regions.
	UNSUPPORTEDOPERATION_INVALIDREGIONDISKENCRYPT = "UnsupportedOperation.InvalidRegionDiskEncrypt"

	// Key-pair login is not available to Windows instances.
	UNSUPPORTEDOPERATION_KEYPAIRUNSUPPORTEDWINDOWS = "UnsupportedOperation.KeyPairUnsupportedWindows"

	// A model whose data disks are all local disks does not support cross-model configuration adjustment.
	UNSUPPORTEDOPERATION_LOCALDATADISKCHANGEINSTANCEFAMILY = "UnsupportedOperation.LocalDataDiskChangeInstanceFamily"

	// The specified disk is converting to a cloud disk. Try again later.
	UNSUPPORTEDOPERATION_LOCALDISKMIGRATINGTOCLOUDDISK = "UnsupportedOperation.LocalDiskMigratingToCloudDisk"

	// This request does not support images in this region. Please change to another image.
	UNSUPPORTEDOPERATION_LOCATIONIMAGENOTSUPPORTED = "UnsupportedOperation.LocationImageNotSupported"

	// Marketplace image instances do not support operating system switching.
	UNSUPPORTEDOPERATION_MARKETIMAGECONVERTOSUNSUPPORTED = "UnsupportedOperation.MarketImageConvertOSUnsupported"

	// The custom images created with the market images cannot be exported.
	UNSUPPORTEDOPERATION_MARKETIMAGEEXPORTUNSUPPORTED = "UnsupportedOperation.MarketImageExportUnsupported"

	// This billing mode does not support the MinCount parameter.
	UNSUPPORTEDOPERATION_MINCOUNTUNSUPPORTEDCHARGETYPE = "UnsupportedOperation.MinCountUnsupportedChargeType"

	// This region does not currently support the MinCount parameter.
	UNSUPPORTEDOPERATION_MINCOUNTUNSUPPORTEDREGION = "UnsupportedOperation.MinCountUnsupportedRegion"

	// Encryption attributes of the system disk cannot be modified. 
	UNSUPPORTEDOPERATION_MODIFYENCRYPTIONNOTSUPPORTED = "UnsupportedOperation.ModifyEncryptionNotSupported"

	// An instance bound with CLB does not support modifying its VPC attributes.
	UNSUPPORTEDOPERATION_MODIFYVPCWITHCLB = "UnsupportedOperation.ModifyVPCWithCLB"

	// This instance is configured with ClassLink. Please cancel the association and continue. 
	UNSUPPORTEDOPERATION_MODIFYVPCWITHCLASSLINK = "UnsupportedOperation.ModifyVPCWithClassLink"

	// This instance type does not support spot instances.
	UNSUPPORTEDOPERATION_NOINSTANCETYPESUPPORTSPOT = "UnsupportedOperation.NoInstanceTypeSupportSpot"

	// A physical network is not supported by this instance.
	UNSUPPORTEDOPERATION_NOVPCNETWORK = "UnsupportedOperation.NoVpcNetwork"

	// Failed to configure the scheduled action for the current instance. 
	UNSUPPORTEDOPERATION_NOTSUPPORTIMPORTINSTANCESACTIONTIMER = "UnsupportedOperation.NotSupportImportInstancesActionTimer"

	// The instance does not support this operation.
	UNSUPPORTEDOPERATION_NOTSUPPORTINSTANCEIMAGE = "UnsupportedOperation.NotSupportInstanceImage"

	// There are unpaid orders for the instance.
	UNSUPPORTEDOPERATION_NOTSUPPORTUNPAIDORDER = "UnsupportedOperation.NotSupportUnpaidOrder"

	// Only a prepaid account supports this operation.
	UNSUPPORTEDOPERATION_ONLYFORPREPAIDACCOUNT = "UnsupportedOperation.OnlyForPrepaidAccount"

	// The original instance type is invalid.
	UNSUPPORTEDOPERATION_ORIGINALINSTANCETYPEINVALID = "UnsupportedOperation.OriginalInstanceTypeInvalid"

	// This model is a periodic contract model and does not support manual renewal mode.
	UNSUPPORTEDOPERATION_PERIODICCONTRACTNOTSUPPORTMANUALRENEW = "UnsupportedOperation.PeriodicContractNotSupportManualRenew"

	// Image preheating is not supported under your account.
	UNSUPPORTEDOPERATION_PREHEATIMAGE = "UnsupportedOperation.PreheatImage"

	// Public images and market images cannot be exported.
	UNSUPPORTEDOPERATION_PUBLICIMAGEEXPORTUNSUPPORTED = "UnsupportedOperation.PublicImageExportUnsupported"

	// This image does not support instance reinstallation.
	UNSUPPORTEDOPERATION_RAWLOCALDISKINSREINSTALLTOQCOW2 = "UnsupportedOperation.RawLocalDiskInsReinstalltoQcow2"

	// The RedHat image cannot be exported.
	UNSUPPORTEDOPERATION_REDHATIMAGEEXPORTUNSUPPORTED = "UnsupportedOperation.RedHatImageExportUnsupported"

	// An instance with an enterprise operating system installed cannot be returned.
	UNSUPPORTEDOPERATION_REDHATINSTANCETERMINATEUNSUPPORTED = "UnsupportedOperation.RedHatInstanceTerminateUnsupported"

	// The operating system of the instance is RedHat, so this operation is not supported.
	UNSUPPORTEDOPERATION_REDHATINSTANCEUNSUPPORTED = "UnsupportedOperation.RedHatInstanceUnsupported"

	// The region is unsupported.
	UNSUPPORTEDOPERATION_REGION = "UnsupportedOperation.Region"

	// Purchasing reserved instances is not supported for the current user.
	UNSUPPORTEDOPERATION_RESERVEDINSTANCEINVISIBLEFORUSER = "UnsupportedOperation.ReservedInstanceInvisibleForUser"

	// You’ve used up your quota for Reserved Instances.
	UNSUPPORTEDOPERATION_RESERVEDINSTANCEOUTOFQUATA = "UnsupportedOperation.ReservedInstanceOutofQuata"

	// Shared images cannot be exported.
	UNSUPPORTEDOPERATION_SHAREDIMAGEEXPORTUNSUPPORTED = "UnsupportedOperation.SharedImageExportUnsupported"

	// 
	UNSUPPORTEDOPERATION_SHAREDIMAGEMODIFYUNSUPPORTED = "UnsupportedOperation.SharedImageModifyUnsupported"

	// This special instance type does not support the operation.
	UNSUPPORTEDOPERATION_SPECIALINSTANCETYPE = "UnsupportedOperation.SpecialInstanceType"

	// Spot instance is not supported in this region.
	UNSUPPORTEDOPERATION_SPOTUNSUPPORTEDREGION = "UnsupportedOperation.SpotUnsupportedRegion"

	// The instance does not support the **no charges when shut down** feature.
	UNSUPPORTEDOPERATION_STOPPEDMODESTOPCHARGING = "UnsupportedOperation.StoppedModeStopCharging"

	// Configuration adjustment of the same type is not supported for instances with no charges when shut down.
	UNSUPPORTEDOPERATION_STOPPEDMODESTOPCHARGINGSAMEFAMILY = "UnsupportedOperation.StoppedModeStopChargingSameFamily"

	// The specified image does not support synchronization to an encrypted custom image.
	UNSUPPORTEDOPERATION_SYNCENCRYPTIMAGENOTSUPPORT = "UnsupportedOperation.SyncEncryptImageNotSupport"

	// The request does not support this type of system disk.
	UNSUPPORTEDOPERATION_SYSTEMDISKTYPE = "UnsupportedOperation.SystemDiskType"

	// The operation is not supported when TencentCloud Automation Tools are offline.
	UNSUPPORTEDOPERATION_TATAGENTNOTONLINE = "UnsupportedOperation.TatAgentNotOnline"

	// Monthly subscription to subcontracted exclusive sale does not support exclusive sales discount higher than existing annual/monthly subscription discount.
	UNSUPPORTEDOPERATION_UNDERWRITEDISCOUNTGREATERTHANPREPAIDDISCOUNT = "UnsupportedOperation.UnderwriteDiscountGreaterThanPrepaidDiscount"

	// For an underwriting instance, `RenewFlag` can only be set to `NOTIFY_AND_AUTO_RENEW`.
	UNSUPPORTEDOPERATION_UNDERWRITINGINSTANCETYPEONLYSUPPORTAUTORENEW = "UnsupportedOperation.UnderwritingInstanceTypeOnlySupportAutoRenew"

	// The current instance does not allow resizing to non-ARM instance types.
	UNSUPPORTEDOPERATION_UNSUPPORTEDARMCHANGEINSTANCEFAMILY = "UnsupportedOperation.UnsupportedARMChangeInstanceFamily"

	// The specified model does not support cross-model configuration adjustment.
	UNSUPPORTEDOPERATION_UNSUPPORTEDCHANGEINSTANCEFAMILY = "UnsupportedOperation.UnsupportedChangeInstanceFamily"

	// Non-ARM model instances cannot be changed to the ARM model.
	UNSUPPORTEDOPERATION_UNSUPPORTEDCHANGEINSTANCEFAMILYTOARM = "UnsupportedOperation.UnsupportedChangeInstanceFamilyToARM"

	// Changing to this model type for this instance is not allowed.
	UNSUPPORTEDOPERATION_UNSUPPORTEDCHANGEINSTANCETOTHISINSTANCEFAMILY = "UnsupportedOperation.UnsupportedChangeInstanceToThisInstanceFamily"

	// This operation is not available for Tencent Cloud International users.
	UNSUPPORTEDOPERATION_UNSUPPORTEDINTERNATIONALUSER = "UnsupportedOperation.UnsupportedInternationalUser"

	// The specified pool is illegal.
	UNSUPPORTEDOPERATION_UNSUPPORTEDPOOL = "UnsupportedOperation.UnsupportedPool"

	// The specified user does not support performing operating system switching.
	UNSUPPORTEDOPERATION_USERCONVERTOSNOTSUPPORT = "UnsupportedOperation.UserConvertOsNotSupport"

	// The quota of user limit operations is insufficient.
	UNSUPPORTEDOPERATION_USERLIMITOPERATIONEXCEEDQUOTA = "UnsupportedOperation.UserLimitOperationExceedQuota"

	// Windows images cannot be exported.
	UNSUPPORTEDOPERATION_WINDOWSIMAGEEXPORTUNSUPPORTED = "UnsupportedOperation.WindowsImageExportUnsupported"

	// The VPC IP address is not in the subnet.
	VPCADDRNOTINSUBNET = "VpcAddrNotInSubNet"

	// The VPC IP address is already occupied.
	VPCIPISUSED = "VpcIpIsUsed"
)
