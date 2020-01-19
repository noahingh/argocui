package main

import (
	"os"
	"time"

	informers "github.com/argoproj/argo/pkg/client/informers/externalversions"
	"github.com/hanjunlee/argocui/internal/app"
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
	argoClientset := argoutil.GetClientset()
	kubeClientset := argoutil.GetKubeClientset()

	factory := informers.NewSharedInformerFactory(argoClientset, 1*time.Minute)
	argoInformer := factory.Argoproj().V1alpha1().Workflows()

	// create a new repo and syncronize.
	repo := repo.NewArgoRepository(argoClientset, argoInformer, kubeClientset)

	neverStop := make(chan struct{}, 1)
	factory.Start(neverStop)
	repo.WaitForSync(neverStop)

	service := argo.NewService(repo)

	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panic(err)
	}
	defer g.Close()

	app.ConfigureGui(g)
	g.SetManagerFunc(app.ManagerFunc(service, g))

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panic(err)
	}
}
