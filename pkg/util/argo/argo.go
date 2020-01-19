package argo

import (
	wf "github.com/argoproj/argo/pkg/apis/workflow/v1alpha1"
	"github.com/argoproj/argo/workflow/util"
)

// WorkflowStatus return the status of workflow.
func WorkflowStatus(w *wf.Workflow) string {
	switch w.Status.Phase {
	case wf.NodeRunning:
		if util.IsWorkflowSuspended(w) {
			return "Running (Suspended)"
		}
		return string(w.Status.Phase)
	case wf.NodeFailed:
		if util.IsWorkflowTerminated(w) {
			return "Failed (Terminated)"
		}
		return string(w.Status.Phase)
	case "", wf.NodePending:
		if !w.ObjectMeta.CreationTimestamp.IsZero() {
			return string(wf.NodePending)
		}
		return "Unknown"
	default:
		return string(w.Status.Phase)
	}
}
