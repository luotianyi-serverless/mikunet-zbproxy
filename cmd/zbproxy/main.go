package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/layou233/zbproxy/v3"
	"github.com/layou233/zbproxy/v3/common"
	"github.com/layou233/zbproxy/v3/common/console"
	"github.com/layou233/zbproxy/v3/common/console/color"
	"github.com/layou233/zbproxy/v3/version"
)

func main() {
	console.SetTitle(fmt.Sprintf("ZBProxy %v | Running...", version.Version))
	fmt.Println(color.Apply(color.FgHiRed, ` ______  _____   _____   _____    _____  __    __ __    __
|___  / |  _  \ |  _  \ |  _  \  /  _  \ \ \  / / \ \  / /
   / /  | |_| | | |_| | | |_| |  | | | |  \ \/ /   \ \/ /`), color.Apply(color.FgHiWhite, `
  / /   |  _  { |  ___/ |  _  /  | | | |   }  {     \  /
 / /__  | |_| | | |     | | \ \  | |_| |  / /\ \    / /
/_____| |_____/ |_|     |_|  \_\ \_____/ /_/  \_\  /_/`))
	fmt.Printf(color.Apply(color.FgHiGreen, "Welcome to ZBProxy %s (%s)!\n"), version.Version, version.CommitHash)
	fmt.Printf(color.Apply(color.FgHiBlack, "Build Information: %s, %s/%s, CGO %s\n"),
		runtime.Version(), runtime.GOOS, runtime.GOARCH, common.CGOHint)
	// go version.CheckUpdate()

	ctx, cancel := context.WithCancel(context.Background())
	instance, err := zbproxy.NewInstance(ctx, zbproxy.Options{
		ConfigFilePath: "ZBProxy.json",
		DisableReload:  false,
	})
	if err != nil {
		panic(err)
	}

	err = instance.Start()
	if err != nil {
		panic(err)
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGHUP, os.Interrupt)
	for {
		select {
		case s := <-signalChan:
			switch s {
			case syscall.SIGHUP:
				instance.Reload()
			case os.Interrupt:
				cancel()
				return
			}
		}
	}
}
