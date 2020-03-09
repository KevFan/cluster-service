package clusterservice

//Action Descriptor of an action
type Action string

//ActionStatus Descriptor of the current status of an action
type ActionStatus string

const (
	//ActionDelete Deletion would be performed
	ActionDelete Action = "delete"
	//ActionStatusDryRun Action will not be performed
	ActionStatusDryRun ActionStatus = "dry run"
	//ActionStatusInProgress Action is being performed currently
	ActionStatusInProgress ActionStatus = "in progress"
	//ActionStatusComplete Action has completed
	ActionStatusComplete ActionStatus = "complete"
	//ActionStatusEmpty Blank status of action
	ActionStatusEmpty ActionStatus = ""
)

//Report Information about what resources are found in the AWS account related to the cluster
type Report struct {
	Items []*ReportItem
}

//MergeForward Merge provided report into this report, assuming the provided report was created after this one
func (r *Report) MergeForward(mergeTarget *Report) {
	mergeTargetItemsLookupMap := make(map[string]*ReportItem)
	for _, item := range mergeTarget.Items {
		mergeTargetItemsLookupMap[item.ID] = item
	}
	for _, item := range r.Items {
		item.MergeForward(mergeTargetItemsLookupMap[item.ID])
	}
}

//ReportItem Information about a specific AWS resource
type ReportItem struct {
	ID           string
	Name         string
	Action       Action
	ActionStatus ActionStatus
}

//MergeForward Merge provided item into this item, assuming the provided item was created after this one
func (r *ReportItem) MergeForward(mergeTarget *ReportItem) {
	//merge target no longer exists, so it must have been deleted
	if mergeTarget == nil {
		r.ActionStatus = ActionStatusComplete
		return
	}
	r.Name = mergeTarget.Name
	r.Action = mergeTarget.Action
	r.ActionStatus = mergeTarget.ActionStatus
}
