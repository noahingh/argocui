package argo

import (
	"context"

)

type followerManager struct {
	ctx context.Context
	cancel context.CancelFunc
	
}

func (f *followerManager) layout()
