package workflow

import (
	wf "github.com/argoproj/argo/pkg/apis/workflow/v1alpha1"
)

type workflows []*wf.Workflow

func (w workflows) Len() int {
	return len(w)
}

func (w workflows) Less(i, j int) bool {
	wi, wj := w[i], w[j]
	return less(wi, wj)
}

func (w workflows) Swap(i, j int) {
	tmp := w[i]
	w[i] = w[j]
	w[j] = tmp
}

func less(item, comp *wf.Workflow) bool {
	iStart := item.ObjectMeta.CreationTimestamp
	iFinish := item.Status.FinishedAt
	cStart := comp.ObjectMeta.CreationTimestamp
	cFinish := comp.Status.FinishedAt

	if iFinish.IsZero() && cFinish.IsZero() {
		return cStart.Before(&iStart)
	}
	if iFinish.IsZero() && !cFinish.IsZero() {
		return true
	}
	if !iFinish.IsZero() && cFinish.IsZero() {
		return false
	}
	// comp finished eariler
	return cFinish.Before(&iFinish)
}

