package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"golang.org/x/net/publicsuffix"

	. "hproxy"
)

var (
	FConfig = flag.String("config", "$HOME/.config/hproxy.json", "config file")
	FSuffix = flag.String("suffix", "", "print public suffix for the given domain")
	FReload = flag.Bool("reload", false, "send signal to reload config file")
)

func serve() {
	L.Printf("Starting...\n")
	L.Printf("PID: %d\n", os.Getpid())

	c, err := NewConfig(*FConfig)
	if err != nil {
		L.Fatalln(err)
	}

	L.Printf("Connecting remote SSH server: %s\n", c.File.RemoteServer)

	go func() {
		normal, err := NewServer(NormalSrv, c)
		if err != nil {
			L.Fatalln(err)
		}
		L.Printf("Local normal HTTP proxy: %s\n", c.File.LocalNormalServer)
		L.Fatalln(http.ListenAndServe(c.File.LocalNormalServer, normal))
	}()

	go func() {
		smart, err := NewServer(SmartSrv, c)
		if err != nil {
			L.Fatalln(err)
		}
		L.Printf("Local smart HTTP proxy: %s\n", c.File.LocalSmartServer)
		L.Fatalln(http.ListenAndServe(c.File.LocalSmartServer, smart))
	}()

	select {}
}

func printSuffix() {
	host := *FSuffix
	tld, _ := publicsuffix.EffectiveTLDPlusOne(host)
	fmt.Printf("EffectiveTLDPlusOne: %s\n", tld)
	suffix, _ := publicsuffix.PublicSuffix(host)
	fmt.Printf("PublicSuffix: %s\n", suffix)
}

func reload() {
	file, err := NewConfigFile(os.ExpandEnv(*FConfig))
	if err != nil {
		L.Fatal(err)
	}
	res, err := http.Get(fmt.Sprintf("http://%s/reload", file.LocalNormalServer))
	if err != nil {
		L.Fatal(err)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		L.Fatal(err)
	}
	fmt.Printf("%s\n", body)
}

func main() {
	flag.Parse()

	if *FSuffix != "" {
		printSuffix()
	} else if *FReload {
		reload()
	} else {
		serve()
	}
}
