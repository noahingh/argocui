package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/hanjunlee/argocui/internal/config"
	svc "github.com/hanjunlee/argocui/internal/runtime"
	"github.com/hanjunlee/argocui/internal/runtime/namespace"
	"github.com/hanjunlee/argocui/internal/runtime/workflow"
	"github.com/hanjunlee/argocui/internal/ui"
	argoutil "github.com/hanjunlee/argocui/pkg/util/argo"
	colorutil "github.com/hanjunlee/argocui/pkg/util/color"

	"github.com/argoproj/argo/pkg/client/clientset/versioned"
	ai "github.com/argoproj/argo/pkg/client/informers/externalversions"
	"github.com/jroimartin/gocui"
	log "github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
	ki "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
)

const (
	noResyncPeriod = 0
)

var (
	version = flag.Bool("version", false, "Check the version.")
	debug   = flag.Bool("debug", false, "Debug mode.")
	trace   = flag.Bool("trace", false, "Debug as trace level.")
)

func main() {
	// flag command
	flag.Parse()

	if *version {
		currentVersion()
		return
	}
	setLog()

	g := newGui()
	defer g.Close()

	kc := argoutil.GetKubeClientset()
	kfactory := ki.NewSharedInformerFactory(kc, 1*time.Second)
	ac := argoutil.GetClientset()
	afactory := ai.NewSharedInformerFactory(ac, 1*time.Second)

	svcs := getRuntimeServices(kc, kfactory, ac, afactory)

	neverStop := make(<-chan struct{})
	kfactory.WaitForCacheSync(neverStop)
	afactory.WaitForCacheSync(neverStop)

	m := ui.NewManager(svcs["mock"], svcs["ns"], svcs["wf"])
	g.SetManager(m)

	if err := m.Keybinding(g); err != nil {
		log.Panic(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panic(err)
	}
}

func currentVersion() {
	fmt.Println(colorutil.ChangeColor(config.Logo, gocui.ColorYellow))
	fmt.Printf("Version: %s\n", config.Version)
	fmt.Printf("Argo Version: %s\n", config.ArgoVersion)
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

func newGui() *gocui.Gui {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panic(err)
	}

	g.Highlight = false
	g.InputEsc = true
	return g
}

func getRuntimeServices(kc kubernetes.Interface, kfactory ki.SharedInformerFactory, ac versioned.Interface, afactory ai.SharedInformerFactory) map[string]svc.UseCase {
	neverStop := make(<-chan struct{})
	// namespace service
	ni := kfactory.Core().V1().Namespaces()
	go ni.Informer().Run(neverStop)
	nr := namespace.NewRepo(ni.Lister())
	ns := svc.NewService(nr)

	// workflow service
	wi := afactory.Argoproj().V1alpha1().Workflows()
	go wi.Informer().Run(neverStop)
	wr := workflow.NewRepo(ac, wi.Lister(), kc)
	ws := svc.NewService(wr)

	cache.WaitForCacheSync(neverStop, ni.Informer().HasSynced, wi.Informer().HasSynced)

	return map[string]svc.UseCase{
		"ns": ns,
		"wf": ws,
	}
}
