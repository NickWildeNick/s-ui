package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime/debug"
	"strings"
	"syscall"

	"github.com/admin8800/s-ui/app"
	"github.com/admin8800/s-ui/cmd"
	"github.com/admin8800/s-ui/config"
)

func runApp(localOnly bool) {
	config.SetLocalOnly(localOnly)

	app := app.NewApp()

	err := app.Init()
	if err != nil {
		log.Fatal(err)
	}

	err = app.Start()
	if err != nil {
		log.Fatal(err)
	}

	sigCh := make(chan os.Signal, 1)
	// Trap shutdown signals
	signal.Notify(sigCh, syscall.SIGHUP, syscall.SIGTERM)
	for {
		sig := <-sigCh

		switch sig {
		case syscall.SIGHUP:
			app.RestartApp()
		default:
			app.Stop()
			return
		}
	}
}

func main() {
	if len(os.Args) < 2 {
		runApp(false)
		return
	}

	if strings.HasPrefix(os.Args[1], "-") {
		appCmd := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
		var localOnly bool
		var showVersion bool
		appCmd.BoolVar(&localOnly, "local", false, "listen on 127.0.0.1 only")
		appCmd.BoolVar(&localOnly, "localhost", false, "listen on 127.0.0.1 only")
		appCmd.BoolVar(&localOnly, "local-only", false, "listen on 127.0.0.1 only")
		appCmd.BoolVar(&showVersion, "v", false, "show version")
		appCmd.Usage = func() {
			out := appCmd.Output()
			fmt.Fprintf(out, "Usage of %s:\n", os.Args[0])
			appCmd.PrintDefaults()
			fmt.Fprintln(out)
			fmt.Fprintln(out, "Commands:")
			fmt.Fprintln(out, "    admin          set/reset/show first admin credentials")
			fmt.Fprintln(out, "    uri            Show panel URI")
			fmt.Fprintln(out, "    migrate        migrate form older version")
			fmt.Fprintln(out, "    setting        set/reset/show settings")
			fmt.Fprintln(out)
			fmt.Fprintln(out, "Admin options:")
			fmt.Fprintln(out, "    -username string")
			fmt.Fprintln(out, "        set login username")
			fmt.Fprintln(out, "    -password string")
			fmt.Fprintln(out, "        set login password")
			fmt.Fprintln(out, "    -reset")
			fmt.Fprintln(out, "        reset first admin credentials")
			fmt.Fprintln(out, "    -show")
			fmt.Fprintln(out, "        show first admin credentials")
			fmt.Fprintln(out)
			fmt.Fprintln(out, "Setting options:")
			fmt.Fprintln(out, "    -port int")
			fmt.Fprintln(out, "        set panel port")
			fmt.Fprintln(out, "    -path string")
			fmt.Fprintln(out, "        set panel path")
			fmt.Fprintln(out, "    -subPort int")
			fmt.Fprintln(out, "        set sub port")
			fmt.Fprintln(out, "    -subPath string")
			fmt.Fprintln(out, "        set sub path")
			fmt.Fprintln(out, "    -reset")
			fmt.Fprintln(out, "        reset all settings")
			fmt.Fprintln(out, "    -show")
			fmt.Fprintln(out, "        show current settings")
		}
		err := appCmd.Parse(os.Args[1:])
		if err != nil {
			log.Fatal(err)
		}
		if showVersion {
			fmt.Println("S-UI Panel\t", config.GetVersion())
			info, ok := debug.ReadBuildInfo()
			if ok {
				for _, dep := range info.Deps {
					if dep.Path == "github.com/sagernet/sing-box" {
						fmt.Println("Sing-Box\t", dep.Version)
						break
					}
				}
			}
			return
		}
		if appCmd.NArg() > 0 {
			log.Fatal("invalid app arguments")
		}
		runApp(localOnly)
		return
	}

	cmd.ParseCmd()
}
