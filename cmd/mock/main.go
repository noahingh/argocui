package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/hanjunlee/argocui/cmd/mock/serializer"
	"github.com/hanjunlee/argocui/internal/config"
	"github.com/hanjunlee/argocui/internal/ui"
	svc "github.com/hanjunlee/argocui/pkg/runtime"
	"github.com/hanjunlee/argocui/pkg/runtime/mock"
	"github.com/hanjunlee/argocui/pkg/runtime/namespace"
	colorutil "github.com/hanjunlee/argocui/pkg/util/color"

	"github.com/jroimartin/gocui"
	log "github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
	"k8s.io/apimachinery/pkg/runtime"
	kubeInformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes/fake"
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
)

func init() {
	namespaces = getMockingNamespaces()
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

	client := fake.NewSimpleClientset(namespaces...)
	kubeFactory := kubeInformers.NewSharedInformerFactory(client, noResyncPeriod)
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

func getMockingNamespaces() []runtime.Object {
	const (
		namespaces = "testdata/namespaces.yaml"
	)
	ret := make([]runtime.Object, 0)

	yamls, _ := serializer.ReadYamlAndSplit(namespaces)
	for _, y := range yamls {
		o, err := serializer.ConvertToNamespace([]byte(y))
		if err != nil {
			panic(err)
		}

		ret = append(ret, o)
	}
	return ret
}

func getRuntimeServices(kubeFactory kubeInformers.SharedInformerFactory) map[string]svc.UseCase {
	// mock service
	mockRepo := &mock.Repo{}
	mockSvc := svc.NewService(mockRepo)

	// namespace service
	i := kubeFactory.Core().V1().Namespaces()
	for _, n := range namespaces {
		i.Informer().GetIndexer().Add(n)
	}

	nsRepo := namespace.NewRepo(i.Lister())
	nsSvc := svc.NewService(nsRepo)

	return map[string]svc.UseCase{
		"mock": mockSvc,
		"ns":   nsSvc,
	}
}
