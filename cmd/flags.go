package main

import (
	"flag"
	"os"
	"path/filepath"
	"time"

	"github.com/jeffersongirao/source-ranges-controller/controller"
	"k8s.io/client-go/util/homedir"
)

type Flags struct {
	flagSet *flag.FlagSet

	Development bool
	ResyncSec   int
	KubeConfig  string
}

func (f *Flags) ControllerConfig() controller.Config {
	return controller.Config{
		ResyncPeriod: time.Duration(f.ResyncSec) * time.Second,
	}
}

func NewFlags() *Flags {
	f := &Flags{
		flagSet: flag.NewFlagSet(os.Args[0], flag.ExitOnError),
	}

	kubehome := filepath.Join(homedir.HomeDir(), ".kube", "config")

	f.flagSet.IntVar(&f.ResyncSec, "resync-seconds", 30, "The number of seconds the controller will resync the resources")
	f.flagSet.StringVar(&f.KubeConfig, "kubeconfig", kubehome, "kubernetes configuration path, only used when development mode enabled")
	f.flagSet.BoolVar(&f.Development, "development", false, "development flag will allow to run the operator outside a kubernetes cluster")

	f.flagSet.Parse(os.Args[1:])

	return f
}
