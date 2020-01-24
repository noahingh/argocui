package main

import (
	"os"
	"time"

	"github.com/hanjunlee/argocui/internal/app/views/etc"
	"github.com/hanjunlee/argocui/internal/app/views/list"
	"github.com/hanjunlee/argocui/internal/app/views/search"
	informers "github.com/argoproj/argo/pkg/client/informers/externalversions"
	"github.com/asaskevich/EventBus"
	"github.com/hanjunlee/argocui/pkg/argo"
	"github.com/hanjunlee/argocui/pkg/argo/repo"
	argoutil "github.com/hanjunlee/argocui/pkg/util/argo"
	"github.com/jroimartin/gocui"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetLevel(log.DebugLevel)
	file, err := os.OpenFile(".argocui.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		log.SetOutput(file)
	} else {
		log.Info("Failed to log to file, using default stderr")
	}
}

func main() {
	var (
		service *argo.Service
	)
	argoClientset := argoutil.GetClientset()
	kubeClientset := argoutil.GetKubeClientset()

	factory := informers.NewSharedInformerFactory(argoClientset, 1*time.Minute)
	argoInformer := factory.Argoproj().V1alpha1().Workflows()

	// create a new repo and syncronize.
	repo := repo.NewArgoRepository(argoClientset, argoInformer, kubeClientset)

	neverStop := make(chan struct{}, 1)
	factory.Start(neverStop)
	repo.WaitForSync(neverStop)

	// create a new service
	service = argo.NewService(repo)

	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panic(err)
	}
	defer g.Close()
	g.Highlight = true
	g.SelFgColor = gocui.ColorYellow
	g.InputEsc = true

	// lay out the gui
	var (
		bus = EventBus.New() 
	)
	g.SetManagerFunc(func(g *gocui.Gui) error {
		maxX, maxY := g.Size()

		ic := etc.NewInfoConfig()
		if err := ic.Layout(g, 1, 0, maxX/5-1, maxY/4-1); err != nil {
			return err
		}

		bc := etc.NewBrandConfig()
		if err := bc.Layout(g, maxX/5, 0, maxX-1, maxY/4-1); err != nil {
			return err
		}

		sc := search.NewConfig()
		if err := sc.Layout(g, service, bus, 0, maxY/4-2, maxX-1, maxY/4); err != nil {
			return err
		}

		lc := list.NewConfig()
		if err := lc.Layout(g, service, bus, 0, maxY/4+1, maxX-1, maxY-1); err != nil {
			return err
		}
		return nil
	})

	etc.GlobalKeybinding(g)

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panic(err)
	}
}
