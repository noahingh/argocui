package view

import (
	"time"

	"github.com/jroimartin/gocui"
	log "github.com/sirupsen/logrus"
)

const (
	refreshSecond = 3
)

// RefreshViewPeriodic refresh table periodically.
func RefreshViewPeriodic(g *gocui.Gui, v *gocui.View, duration time.Duration, f func() error) {
	t := time.NewTicker(duration)
	go g.Update(func(g *gocui.Gui) error {
		return f()
	})

	for {
		select {
		case <-t.C:
			_, err := g.View(v.Name())
			if err != nil {
				log.Warnf("'%s' view does not exist anymore", v.Name())
				return
			}

			log.Tracef("refresh '%s' view", v.Name())
			go g.Update(func(g *gocui.Gui) error {
				return f()
			})
		}
	}
}
