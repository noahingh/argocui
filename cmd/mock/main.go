package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/hanjunlee/argocui/internal/config"
	"github.com/hanjunlee/argocui/internal/ui"
	"github.com/hanjunlee/argocui/pkg/runtime"
	"github.com/hanjunlee/argocui/pkg/runtime/mock"
	"github.com/hanjunlee/argocui/pkg/runtime/namespace"
	argoutil "github.com/hanjunlee/argocui/pkg/util/argo"
	colorutil "github.com/hanjunlee/argocui/pkg/util/color"

	"github.com/jroimartin/gocui"
	log "github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
	kubeInformers "k8s.io/client-go/informers"
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

	// TODO: mocking the clientset with test data.
	kubeFactory := kubeInformers.NewSharedInformerFactory(argoutil.GetKubeClientset(), 1*time.Minute)
	svcs := getRuntimeServices(kubeFactory)

	neverStop := make(<-chan struct{})
	kubeFactory.WaitForCacheSync(neverStop)

	m := ui.NewManager(svcs["mock"], svcs)
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

func getRuntimeServices(kubeFactory kubeInformers.SharedInformerFactory) map[string]runtime.UseCase {
	neverStop := make(<-chan struct{})

	// mock service
	mockRepo := &mock.Repo{}
	mockSvc := runtime.NewService(mockRepo)

	// namespace service
	nsInformer := kubeFactory.Core().V1().Namespaces()
	go nsInformer.Informer().Run(neverStop)

	nsRepo := namespace.NewRepo(nsInformer.Lister())
	nsSvc := runtime.NewService(nsRepo)

	return map[string]runtime.UseCase{
		"mock": mockSvc,
		"ns":   nsSvc,
	}
}
