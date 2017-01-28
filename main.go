// transfer2go - Go implementation of loosely coupled, distributed agents for data transfer
//
// Copyright (c) 2017 - Valentin Kuznetsov <vkuznet@gmail.com>
//
package main

import (
	"flag"
	"fmt"
	"github.com/vkuznet/transfer2go/client"
	"github.com/vkuznet/transfer2go/server"
	"os"
	"os/user"
)

func main() {
	var agent string
	flag.StringVar(&agent, "agent", "", "Remote agent end-point")
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Errorf("Unable to get current directory, error=%v\n", err)
		os.Exit(1)
	}
	var catalog string
	flag.StringVar(&catalog, "catalog", pwd, "Agent catalog, e.g. dir name or DB uri")
	var status bool
	flag.BoolVar(&status, "status", false, "Return status info about the agent")
	var src string
	flag.StringVar(&src, "src", "", "Source logical file name")
	var dst string
	flag.StringVar(&dst, "dst", "", "Destination end-point")
	var register string
	flag.StringVar(&register, "register", "", "Registration end-point")
	var url string
	flag.StringVar(&url, "url", "", "Server end-point url, e.g. https://a.b.com/transfer2go")
	var port string
	flag.StringVar(&port, "port", "", "Server port number, default 8989")
	var alias string
	flag.StringVar(&alias, "alias", makeSiteName(), "Server alias name, e.g. T3_US_Name")
	var interval int64
	flag.Int64Var(&interval, "interval", 600, "Server metrics interval, default 600 seconds")
	var verbose int
	flag.IntVar(&verbose, "verbose", 0, "Verbosity level, default 0")
	flag.Parse()
	checkX509()
	client.VERBOSE = verbose
	if url != "" {
		server.Server(port, url, alias, register, catalog, interval)
	} else {
		if status {
			client.Status(agent)
		} else {
			client.Process(agent, src, dst)
		}
	}
}

// helper function to construct site name
func makeSiteName() string {
	host, err := os.Hostname()
	if err != nil {
		panic(fmt.Sprintf("Unable to get hostname, error=%v", err))
	}
	return fmt.Sprintf("T4_%s_%v", host, os.Getuid())
}

// helper function to check X509 settings
func checkX509() {
	uproxy := os.Getenv("X509_USER_PROXY")
	uckey := os.Getenv("X509_USER_KEY")
	ucert := os.Getenv("X509_USER_CERT")
	var check int
	if uproxy == "" {
		// check if /tmp/x509up_u$UID exists
		u, err := user.Current()
		if err == nil {
			fname := fmt.Sprintf("/tmp/x509up_u%s", u.Uid)
			if _, err := os.Stat(fname); err != nil {
				check += 1
			}
		}
	}
	if uckey == "" && ucert == "" {
		check += 1
	}
	if check > 1 {
		fmt.Println("Neither X509_USER_PROXY or X509_USER_KEY/X509_USER_CERT are set")
		fmt.Println("In order to run please obtain valid proxy via \"voms-proxy-init -voms cms -rfc\"")
		fmt.Println("and setup X509_USER_PROXY or setup X509_USER_KEY/X509_USER_CERT in your environment")
		os.Exit(-1)
	}
}