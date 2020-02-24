package argo

import (
	"github.com/asaskevich/EventBus"
	"github.com/hanjunlee/argocui/pkg/argo"
	"github.com/hanjunlee/argocui/pkg/kube"
	"github.com/jroimartin/gocui"
)

// Manager is the manager of the Argo cui.
type Manager struct {
	s  *searchManager
	cm *collectionManager
	f  *followerManager
	t  *treeManager
	nm *namespaceManager
}

// NewManager create a new manager of the Argo cui.
func NewManager(au argo.UseCase, ku kube.UseCase, bus EventBus.Bus) *Manager {
	return &Manager{
		s:  newSearchManager(au, bus),
		cm: newCollectionManager(au, bus),
		f:  newFollowerManager(au, bus),
		t:  newTreeManager(au, bus),
		nm: newNamespaceManager(ku, bus),
	}
}

// Layout lay out the Argo cui.
func (m *Manager) Layout(g *gocui.Gui) error {
	x, y := g.Size()
	if err := m.s.layout(g, 0, y/4, x-1, y/4+2); err != nil {
		return err
	}

	if err := m.f.layout(g, 0, y/4+3, x-1, y-1); err != nil {
		return err
	}

	if err := m.t.layout(g, 0, y/4+3, x-1, y-1); err != nil {
		return err
	}

	if err := m.nm.layout(g, 0, y/4+3, x-1, y-1); err != nil {
		return err
	}

	// collection view should be on the top when it start.
	if err := m.cm.layout(g, 0, y/4+3, x-1, y-1); err != nil {
		return err
	}

	return nil
}
