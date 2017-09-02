package main

import (
	"os"

	logs "github.com/appscode/log/golog"
	_ "github.com/appscode/stash/apis/install"
	_ "github.com/appscode/stash/client/fake"
	"github.com/appscode/stash/pkg/cmds"
	_ "k8s.io/client-go/kubernetes/fake"
)

func main() {
	logs.InitLogs()
	defer logs.FlushLogs()

	if err := cmds.NewCmdStash(Version).Execute(); err != nil {
		os.Exit(1)
	}
	os.Exit(0)
}
