package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/hanjunlee/argocui/internal/config"
	am "github.com/hanjunlee/argocui/internal/managers/argo"
	"github.com/hanjunlee/argocui/internal/managers/etc"
	"github.com/hanjunlee/argocui/pkg/argo"
	"github.com/hanjunlee/argocui/pkg/argo/repo"
	"github.com/hanjunlee/argocui/pkg/kube"
	argoutil "github.com/hanjunlee/argocui/pkg/util/argo"
	colorutil "github.com/hanjunlee/argocui/pkg/util/color"

	informers "github.com/argoproj/argo/pkg/client/informers/externalversions"
	"github.com/asaskevich/EventBus"
	"github.com/jroimartin/gocui"
	log "github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	version  = flag.Bool("version", false, "Check the version.")
	debug    = flag.Bool("debug", false, "Debug mode.")
	trace    = flag.Bool("trace", false, "Debug as trace level.")
	readOnly = flag.Bool("ro", false, "Read only mode.")
)

func main() {
	// flag command
	flag.Parse()
	if *version {
		currentVersion()
		return 
	}
	setConfig()
	setLog()

	// create a new repo and syncronize.
	var (
		argoService *argo.Service
		kubeService *kube.Service
	)
	argoClientset := argoutil.GetClientset()
	kubeClientset := argoutil.GetKubeClientset()

	factory := informers.NewSharedInformerFactory(argoClientset, 1*time.Minute)
	argoInformer := factory.Argoproj().V1alpha1().Workflows()

	repo := repo.NewArgoRepository(argoClientset, argoInformer, kubeClientset)

	neverStop := make(chan struct{}, 1)
	factory.Start(neverStop)
	repo.WaitForSync(neverStop)

	// create new services
	argoService = argo.NewService(repo)
	kubeService = kube.NewService(kubeClientset)

	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panic(err)
	}
	defer g.Close()

	g.Highlight = true
	g.SelFgColor = gocui.ColorYellow
	g.InputEsc = true

	em := etc.NewManager()

	bus := EventBus.New()
	m := am.NewManager(argoService, kubeService, bus)
	g.SetManager(em, m)

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panic(err)
	}
}

func currentVersion() {
	fmt.Println(colorutil.ChangeColor(config.Logo, gocui.ColorYellow))
	fmt.Printf("Version: %s\n", config.Version)
	fmt.Printf("Argo Version: %s\n", config.ArgoVersion)
}

func setConfig() {
	if *readOnly {
		config.ReadOnly = true
	}
}

func setLog() {
	log.SetLevel(log.InfoLevel)
	if *debug {
		log.SetLevel(log.DebugLevel)
	}
	if *trace {
		log.SetLevel(log.TraceLevel)
	}

	path := filepath.Join(os.Getenv("HOME"), "/.argocui/log")
	log.SetOutput(&lumberjack.Logger{
		Filename:   path,
		MaxSize:    500,
		MaxBackups: 1,
		MaxAge:     7,
		Compress:   true,
	})
}
