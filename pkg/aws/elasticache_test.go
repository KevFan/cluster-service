package aws

import (
	"github.com/aws/aws-sdk-go/service/elasticache"
	"github.com/aws/aws-sdk-go/service/resourcegroupstaggingapi"
	"github.com/integr8ly/cluster-service/pkg/clusterservice"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	//"reflect"
	"testing"
)

func TestElasticacheEngine_DeleteResourcesForCluster(t *testing.T) {
	fakeClusterId := "testClusterId"
	fakeLogger, err := fakeLogger(func(l *logrus.Entry) error {
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}

	type fields struct {
		elasticacheClient func() *elasticacheClientMock
		taggingClient     func() *resourcetaggingClientMock
		logger            *logrus.Entry
	}
	type args struct {
		clusterId string
		tags      map[string]string
		dryRun    bool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*clusterservice.ReportItem
		wantFn  func(mock *elasticacheClientMock) error
		wantErr string
	}{
		{
			name: "error when describing clusters fail",
			fields: fields{
				elasticacheClient: func() *elasticacheClientMock {
					fakeClient, err := fakeElasticacheClient(func(c *elasticacheClientMock) error {
						return nil
					})
					if err != nil {
						t.Fatal(err)
					}
					return fakeClient
				},
				taggingClient: func() *resourcetaggingClientMock {
					fakeTaggingClient, err := fakeResourcetaggingClient(func(c *resourcetaggingClientMock) error {
						c.GetResourcesFunc = func(in1 *resourcegroupstaggingapi.GetResourcesInput) (output *resourcegroupstaggingapi.GetResourcesOutput, e error) {
							return nil, errors.New("")
						}
						return nil
					})
					if err != nil {
						t.Fatal(err)
					}

					return fakeTaggingClient
				},
				logger: fakeLogger,
			},
			args: args{
				clusterId: fakeClusterId,
				dryRun:    true,
			},
			wantErr: "failed to describe cache clusters: ",
		}, {
			name: "error when getting cacheCluster output",
			fields: fields{
				elasticacheClient: func() *elasticacheClientMock {
					fakeClient, err := fakeElasticacheClient(func(c *elasticacheClientMock) error {
						c.DescribeCacheClustersFunc = func(in1 *elasticache.DescribeCacheClustersInput) (*elasticache.DescribeCacheClustersOutput, error) {
							return nil, errors.New("")
						}
						return nil
					})
					if err != nil {
						t.Fatal(err)
					}
					return fakeClient
				},
				taggingClient: func() *resourcetaggingClientMock {
					fakeTaggingClient, err := fakeResourcetaggingClient(func(c *resourcetaggingClientMock) error {
						return nil
					})
					if err != nil {
						t.Fatal(err)
					}

					return fakeTaggingClient
				},
				logger: fakeLogger,
			},
			args: args{
				clusterId: fakeClusterId,
				dryRun:    true,
			},
			wantErr: "cannot get cacheCluster output: ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeClient := tt.fields.elasticacheClient()
			r := &ElasticacheEngine{
				elasticacheClient: fakeClient,
				taggingClient:     tt.fields.taggingClient(),
				logger:            tt.fields.logger,
			}
			got, err := r.DeleteResourcesForCluster(tt.args.clusterId, nil, tt.args.dryRun)
			if tt.wantErr != "" && err.Error() != tt.wantErr {
				t.Errorf("DeleteResourcesForCluster() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !equalReportItems(got, tt.want) {
				t.Errorf("DeleteResourcesForCluster() got = %v, want %v", got, tt.want)
			}
			if tt.wantFn != nil {
				if err := tt.wantFn(fakeClient); err != nil {
					t.Errorf("DeleteResourcesForCluster() err = %v", err)
				}
			}
		})
	}
}

//func equalReportItems(a, b []*clusterservice.ReportItem) bool {
//	if len(a) != len(b) {
//		return false
//	}
//	for i, _ := range a {
//		if !reflect.DeepEqual(*a[i], *b[i]) {
//			return false
//		}
//	}
//	return true
//}
