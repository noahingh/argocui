package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/hanjunlee/argocui/cmd/mock/serializer"
	"github.com/hanjunlee/argocui/internal/config"
	svc "github.com/hanjunlee/argocui/internal/runtime"
	"github.com/hanjunlee/argocui/internal/runtime/mock"
	"github.com/hanjunlee/argocui/internal/runtime/namespace"
	workflow "github.com/hanjunlee/argocui/internal/runtime/workflow/fake"
	"github.com/hanjunlee/argocui/internal/ui"
	colorutil "github.com/hanjunlee/argocui/pkg/util/color"

	"github.com/argoproj/argo/pkg/client/clientset/versioned"
	af "github.com/argoproj/argo/pkg/client/clientset/versioned/fake"
	ai "github.com/argoproj/argo/pkg/client/informers/externalversions"
	"github.com/jroimartin/gocui"
	log "github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
	"k8s.io/apimachinery/pkg/runtime"
	ki "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	kf "k8s.io/client-go/kubernetes/fake"
)

const (
	noResyncPeriod = 0
)

var (
	version = flag.Bool("version", false, "Check the version.")
	debug   = flag.Bool("debug", false, "Debug mode.")
	trace   = flag.Bool("trace", false, "Debug as trace level.")
)

var (
	namespaces []runtime.Object
	workflows  []runtime.Object
)

func init() {
	namespaces = getMockingNamespaces()
	workflows = getMockingWorkflows()
}

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

	kc := kf.NewSimpleClientset(namespaces...)
	kfactory := ki.NewSharedInformerFactory(kc, noResyncPeriod)
	ac := af.NewSimpleClientset(workflows...)
	afactory := ai.NewSharedInformerFactory(ac, noResyncPeriod)

	svcs := getRuntimeServices(kc, kfactory, ac, afactory)

	neverStop := make(<-chan struct{})
	kfactory.WaitForCacheSync(neverStop)
	afactory.WaitForCacheSync(neverStop)

	m := ui.NewManager(svcs["wf"], svcs)
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

func getMockingNamespaces() []runtime.Object {
	const (
		file = "testdata/namespaces.yaml"
	)
	ret := make([]runtime.Object, 0)

	yamls, _ := serializer.ReadYamlAndSplit(file)
	for _, y := range yamls {
		o, err := serializer.ConvertToNamespace([]byte(y))
		if err != nil {
			panic(err)
		}

		ret = append(ret, o)
	}
	return ret
}

func getMockingWorkflows() []runtime.Object {
	const (
		file = "testdata/workflows.yaml"
	)
	ret := make([]runtime.Object, 0)

	yamls, _ := serializer.ReadYamlAndSplit(file)
	for _, y := range yamls {
		o, err := serializer.ConvertToWorkflow([]byte(y))
		if err != nil {
			panic(err)
		}

		ret = append(ret, o)
	}
	return ret
}

func getRuntimeServices(kc kubernetes.Interface, kfactory ki.SharedInformerFactory, ac versioned.Interface, afactory ai.SharedInformerFactory) map[string]svc.UseCase {
	// mock service
	mr := &mock.Repo{}
	ms := svc.NewService(mr)

	// namespace service
	ni := kfactory.Core().V1().Namespaces()
	for _, n := range namespaces {
		ni.Informer().GetIndexer().Add(n)
	}

	nr := namespace.NewRepo(ni.Lister())
	ns := svc.NewService(nr)

	// workflow service
	wi := afactory.Argoproj().V1alpha1().Workflows()
	for _, w := range workflows {
		wi.Informer().GetIndexer().Add(w)
	}

	wr := workflow.NewRepo(ac, wi.Informer(), wi.Lister())
	ws := svc.NewService(wr)

	return map[string]svc.UseCase{
		"mock": ms,
		"ns":   ns,
		"wf":   ws,
	}
}
