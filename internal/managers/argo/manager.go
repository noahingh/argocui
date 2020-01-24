package argo

import (
	"github.com/jroimartin/gocui"
	"github.com/asaskevich/EventBus"
	"github.com/hanjunlee/argocui/pkg/argo"
)

// Manager is the manager of the Argo cui.
type Manager struct {
	cm *collectionManager
}

// NewManager create a new manager of the Argo cui.
func NewManager(uc argo.UseCase, bus EventBus.Bus) *Manager {
	return &Manager{
		cm: newCollectionManager(uc, bus),
	}
}

// Layout lay out the Argo cui.
func (m *Manager) Layout(g *gocui.Gui) error {
	x, y := g.Size()
	if err := m.cm.layout(g, 0, y/4, x-1, y-3); err != nil {
		return err
	}
	return nil
}
