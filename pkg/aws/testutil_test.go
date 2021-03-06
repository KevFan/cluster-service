package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/service/elasticache"

	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/aws/aws-sdk-go/service/resourcegroupstaggingapi"
	"github.com/integr8ly/cluster-service/pkg/clusterservice"
	"github.com/sirupsen/logrus"

	"github.com/aws/aws-sdk-go/aws"
)

const (
	//generic variables
	fakeARN                = "arn:fake:testIdentifier"
	fakeResourceIdentifier = "testIdentifier"
	fakeClusterId          = "clusterId"

	//rds-specific
	fakeRDSClientTagKey                     = tagKeyClusterId
	fakeRDSClientTagVal                     = "fakeVal"
	fakeRDSClientInstanceIdentifier         = fakeResourceIdentifier
	fakeRDSClientInstanceARN                = fakeARN
	fakeRDSClientInstanceDeletionProtection = true

	//ELasticache-specific
	fakeElasticacheClientName               = "elasticache Replication group"
	fakeElasticacheClientReplicationGroupId = "testRepGroupID"
	fakeElasticacheClientDescription        = "TestDescription"
	fakeElasticacheClientEngine             = "redis"
	fakeElasticacheClientCacheNodeType      = "cache.t2.micro"
	fakeElasticacheClientStatusAvailable    = "available"
	fakeResourceTaggingClientArn            = "arn:fake:testIdentifier"
	fakeResourceTaggingClientTagKey         = "testTag"
	fakeResourceTaggingClientTagValue       = "testValue"
	fakeClusterID                           = "testClusterID"
	fakeCacheClusterStatus                  = "available"
	fakeElasticacheSnapshotName             = "elasticache snapshot"
	fakeElasticacheSnapshotStatus           = "available"

	//resource tagging-specific
	fakeResourceTagMappingARN = fakeARN

	//resource manager-specific
	fakeResourceManagerName = "Fake Action Engine"

	// db snapshot
	fakeSnapshotType    = "manual"
	fakeSnapshotStatus  = "available"
	fakeRDSSnapshotName = "rds-snapshot"
)

func fakeReportItemDeleting() *clusterservice.ReportItem {
	return &clusterservice.ReportItem{
		ID:           fakeRDSClientInstanceARN,
		Name:         fakeRDSClientInstanceIdentifier,
		Action:       clusterservice.ActionDelete,
		ActionStatus: clusterservice.ActionStatusInProgress,
	}
}

func fakeReportItemDryRun() *clusterservice.ReportItem {
	return &clusterservice.ReportItem{
		ID:           fakeRDSClientInstanceARN,
		Name:         fakeRDSClientInstanceIdentifier,
		Action:       clusterservice.ActionDelete,
		ActionStatus: clusterservice.ActionStatusDryRun,
	}
}

func fakeRDSClientTag() *rds.Tag {
	return &rds.Tag{
		Key:   aws.String(fakeRDSClientTagKey),
		Value: aws.String(fakeRDSClientTagVal),
	}
}

func fakeRDSClientDBInstance() *rds.DBInstance {
	return &rds.DBInstance{
		DBInstanceIdentifier: aws.String(fakeRDSClientInstanceIdentifier),
		DBInstanceArn:        aws.String(fakeRDSClientInstanceARN),
		DeletionProtection:   aws.Bool(fakeRDSClientInstanceDeletionProtection),
	}
}

func fakeResourceTagMappingTag() *resourcegroupstaggingapi.Tag {
	return &resourcegroupstaggingapi.Tag{
		Key:   aws.String(tagKeyClusterId),
		Value: aws.String(fakeClusterId),
	}
}

func fakeResourceTagMapping() *resourcegroupstaggingapi.ResourceTagMapping {
	return &resourcegroupstaggingapi.ResourceTagMapping{
		ComplianceDetails: nil,
		ResourceARN:       aws.String(fakeResourceTagMappingARN),
		Tags: []*resourcegroupstaggingapi.Tag{
			fakeResourceTagMappingTag(),
		},
	}
}

func fakeRDSClientDBSnapshots() []*rds.DBSnapshot {
	return []*rds.DBSnapshot{
		fakeRDSSnapshot(),
	}
}

func fakeRDSSnapshot() *rds.DBSnapshot {
	return &rds.DBSnapshot{
		Engine:               aws.String(fakeElasticacheClientEngine),
		DBInstanceIdentifier: aws.String(fakeResourceIdentifier),
		DBSnapshotIdentifier: aws.String(fakeRDSSnapshotName),
		Status:               aws.String(fakeSnapshotStatus),
		SnapshotType:         aws.String(fakeSnapshotType),
	}
}

func fakeRDSClient(modifyFn func(c *rdsClientMock) error) (*rdsClientMock, error) {
	if modifyFn == nil {
		return nil, errorMustBeDefined("modifyFn")
	}
	client := &rdsClientMock{
		DescribeDBInstancesFunc: func(in1 *rds.DescribeDBInstancesInput) (output *rds.DescribeDBInstancesOutput, e error) {
			return &rds.DescribeDBInstancesOutput{
				DBInstances: []*rds.DBInstance{
					fakeRDSClientDBInstance(),
				},
			}, nil
		},
		ListTagsForResourceFunc: func(in1 *rds.ListTagsForResourceInput) (output *rds.ListTagsForResourceOutput, e error) {
			return &rds.ListTagsForResourceOutput{
				TagList: []*rds.Tag{
					fakeRDSClientTag(),
				},
			}, nil
		},
		ModifyDBInstanceFunc: func(in1 *rds.ModifyDBInstanceInput) (output *rds.ModifyDBInstanceOutput, e error) {
			return &rds.ModifyDBInstanceOutput{
				DBInstance: fakeRDSClientDBInstance(),
			}, nil
		},
		DeleteDBInstanceFunc: func(in1 *rds.DeleteDBInstanceInput) (output *rds.DeleteDBInstanceOutput, e error) {
			return &rds.DeleteDBInstanceOutput{
				DBInstance: fakeRDSClientDBInstance(),
			}, nil
		},
		DescribeDBSnapshotsFunc: func(in1 *rds.DescribeDBSnapshotsInput) (*rds.DescribeDBSnapshotsOutput, error) {
			return &rds.DescribeDBSnapshotsOutput{
				DBSnapshots: fakeRDSClientDBSnapshots(),
			}, nil
		},
		DeleteDBSnapshotFunc: func(in1 *rds.DeleteDBSnapshotInput) (*rds.DeleteDBSnapshotOutput, error) {
			return &rds.DeleteDBSnapshotOutput{}, nil
		},
	}
	if err := modifyFn(client); err != nil {
		return nil, errorModifyFailed(err)
	}
	return client, nil
}

func fakeS3Client(modifyFn func(c *s3ClientMock) error) (*s3ClientMock, error) {
	if modifyFn == nil {
		return nil, errorMustBeDefined("modifyFn")
	}
	client := &s3ClientMock{
		DeleteBucketFunc: func(in1 *s3.DeleteBucketInput) (output *s3.DeleteBucketOutput, e error) {
			return &s3.DeleteBucketOutput{}, nil
		},
	}
	if err := modifyFn(client); err != nil {
		return nil, errorModifyFailed(err)
	}
	return client, nil
}

func fakeTaggingClient(modifyFn func(c *taggingClientMock) error) (*taggingClientMock, error) {
	if modifyFn == nil {
		return nil, errorMustBeDefined("modifyFn")
	}
	client := &taggingClientMock{
		GetResourcesFunc: func(in1 *resourcegroupstaggingapi.GetResourcesInput) (output *resourcegroupstaggingapi.GetResourcesOutput, e error) {
			return &resourcegroupstaggingapi.GetResourcesOutput{
				ResourceTagMappingList: []*resourcegroupstaggingapi.ResourceTagMapping{
					fakeResourceTagMapping(),
				},
			}, nil
		},
	}
	if err := modifyFn(client); err != nil {
		return nil, fmt.Errorf("error occurred in modify function: %w", err)
	}
	return client, nil
}

func fakeS3BatchClient(modifyFn func(c *s3BatchDeleteClientMock) error) (*s3BatchDeleteClientMock, error) {
	if modifyFn == nil {
		return nil, errorMustBeDefined("modifyFn")
	}
	client := &s3BatchDeleteClientMock{
		DeleteFunc: func(in1 context.Context, in2 s3manager.BatchDeleteIterator) error {
			return nil
		},
	}
	if err := modifyFn(client); err != nil {
		return nil, fmt.Errorf("error occurred in modify function: %w", err)
	}
	return client, nil
}

//ELASTICACHE
func fakeElasticacheSnapshot() *elasticache.Snapshot {
	return &elasticache.Snapshot{
		CacheClusterId: aws.String(fakeClusterID),
		CacheNodeType:  aws.String(fakeElasticacheClientCacheNodeType),
		Engine:         aws.String(fakeElasticacheClientEngine),
		SnapshotName:   aws.String(fakeElasticacheSnapshotName),
		SnapshotStatus: aws.String(fakeElasticacheSnapshotStatus),
	}
}
func fakeReportItemElasticacheSnapshotDeleting() *clusterservice.ReportItem {
	return &clusterservice.ReportItem{
		ID:           fakeARN,
		Name:         fakeResourceIdentifier,
		Action:       clusterservice.ActionDelete,
		ActionStatus: clusterservice.ActionStatusInProgress,
	}
}
func fakeReportItemElasticacheSnapshotDryRun() *clusterservice.ReportItem {
	return &clusterservice.ReportItem{
		ID:           fakeARN,
		Name:         fakeResourceIdentifier,
		Action:       clusterservice.ActionDelete,
		ActionStatus: clusterservice.ActionStatusDryRun,
	}
}
func fakeReportItemRDSSnapshotDeleting() *clusterservice.ReportItem {
	return &clusterservice.ReportItem{
		ID:           fakeARN,
		Name:         fakeResourceIdentifier,
		Action:       clusterservice.ActionDelete,
		ActionStatus: clusterservice.ActionStatusInProgress,
	}
}
func fakeReportItemRDSSnapshotDryRun() *clusterservice.ReportItem {
	return &clusterservice.ReportItem{
		ID:           fakeARN,
		Name:         fakeResourceIdentifier,
		Action:       clusterservice.ActionDelete,
		ActionStatus: clusterservice.ActionStatusDryRun,
	}
}
func fakeReportItemReplicationGroupDeleting() *clusterservice.ReportItem {
	return &clusterservice.ReportItem{
		ID:           fakeElasticacheClientReplicationGroupId,
		Name:         fakeElasticacheClientName,
		Action:       clusterservice.ActionDelete,
		ActionStatus: clusterservice.ActionStatusInProgress,
	}
}

func fakeReportItemReplicationGroupDryRun() *clusterservice.ReportItem {
	return &clusterservice.ReportItem{
		ID:           fakeElasticacheClientReplicationGroupId,
		Name:         fakeElasticacheClientName,
		Action:       clusterservice.ActionDelete,
		ActionStatus: clusterservice.ActionStatusDryRun,
	}
}

func fakeElasticacheReplicationGroup() *elasticache.ReplicationGroup {
	return &elasticache.ReplicationGroup{
		CacheNodeType:      aws.String(fakeElasticacheClientCacheNodeType),
		Description:        aws.String(fakeElasticacheClientDescription),
		ReplicationGroupId: aws.String(fakeElasticacheClientReplicationGroupId),
		Status:             aws.String(fakeElasticacheClientStatusAvailable),
	}
}
func fakeElasticacheCacheCluster() *elasticache.CacheCluster {
	return &elasticache.CacheCluster{
		CacheClusterId:     aws.String(fakeClusterID),
		CacheClusterStatus: aws.String(fakeCacheClusterStatus),
		CacheNodeType:      aws.String(fakeElasticacheClientCacheNodeType),
		Engine:             aws.String(fakeElasticacheClientEngine),
		ReplicationGroupId: aws.String(fakeElasticacheClientReplicationGroupId)}
}

func fakeElasticacheClient(modifyFn func(c *elasticacheClientMock) error) (*elasticacheClientMock, error) {
	if modifyFn == nil {
		return nil, fmt.Errorf("modifyFn must be defined")
	}
	client := &elasticacheClientMock{
		DescribeReplicationGroupsFunc: func(in1 *elasticache.DescribeReplicationGroupsInput) (output *elasticache.DescribeReplicationGroupsOutput, e error) {
			return &elasticache.DescribeReplicationGroupsOutput{
				ReplicationGroups: []*elasticache.ReplicationGroup{
					fakeElasticacheReplicationGroup(),
				}}, nil
		},
		DescribeSnapshotsFunc: func(in1 *elasticache.DescribeSnapshotsInput) (output *elasticache.DescribeSnapshotsOutput, e error) {
			return &elasticache.DescribeSnapshotsOutput{
				Snapshots: []*elasticache.Snapshot{
					fakeElasticacheSnapshot(),
				}}, nil
		},
		DescribeCacheClustersFunc: func(in1 *elasticache.DescribeCacheClustersInput) (output *elasticache.DescribeCacheClustersOutput, e error) {
			return &elasticache.DescribeCacheClustersOutput{
				CacheClusters: []*elasticache.CacheCluster{
					fakeElasticacheCacheCluster(),
				}}, nil
		},
		DeleteReplicationGroupFunc: func(in1 *elasticache.DeleteReplicationGroupInput) (output *elasticache.DeleteReplicationGroupOutput, e error) {
			return &elasticache.DeleteReplicationGroupOutput{
				ReplicationGroup: fakeElasticacheReplicationGroup(),
			}, nil
		},
		DeleteSnapshotFunc: func(in1 *elasticache.DeleteSnapshotInput) (output *elasticache.DeleteSnapshotOutput, e error) {
			return &elasticache.DeleteSnapshotOutput{
				Snapshot: fakeElasticacheSnapshot(),
			}, nil
		},
	}
	if err := modifyFn(client); err != nil {
		return nil, fmt.Errorf("error occurred in modify function: %w", err)
	}
	return client, nil
}

func fakeLogger(modifyFn func(l *logrus.Entry) error) (*logrus.Entry, error) {
	if modifyFn == nil {
		return nil, errorMustBeDefined("modifyFn")
	}
	logger := logrus.NewEntry(logrus.StandardLogger())
	if err := modifyFn(logger); err != nil {
		return nil, errorModifyFailed(err)
	}
	return logger, nil
}

func fakeClusterManager(modifyFn func(e *ClusterResourceManagerMock) error) (*ClusterResourceManagerMock, error) {
	if modifyFn == nil {
		return nil, errorMustBeDefined("modifyFn")
	}
	clusterManager := &ClusterResourceManagerMock{
		DeleteResourcesForClusterFunc: func(clusterId string, tags map[string]string, dryRun bool) (items []*clusterservice.ReportItem, e error) {
			return []*clusterservice.ReportItem{
				fakeReportItemDeleting(),
			}, nil
		},
		GetNameFunc: func() string {
			return fakeResourceManagerName
		},
	}
	if err := modifyFn(clusterManager); err != nil {
		return nil, errorModifyFailed(err)
	}
	return clusterManager, nil
}

func errorMustBeDefined(varName string) error {
	return fmt.Errorf("%s must be defined", varName)
}

func errorModifyFailed(err error) error {
	return fmt.Errorf("error occurred while modifying resource: %w", err)
}
