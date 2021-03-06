package main

import (
	"fmt"
	"log"
	"os"
	"os/user"

	"github.com/jessevdk/go-flags"
)

const (
	iniFile = ".gof3r.ini"
)

// CommonOpts are Options common to all commands
type CommonOpts struct {
	NoSSL       bool   `long:"no-ssl" description:"Do not use SSL for endpoint connection." ini-name:"no-ssl"`
	NoMd5       bool   `long:"no-md5" description:"Do not use md5 hash checking to ensure data integrity. By default, the md5 hash of is calculated concurrently during puts, stored at <bucket>.md5/<key>.md5, and verified on gets." ini-name:"no-md5"`
	Concurrency int    `long:"concurrency" short:"c" default:"10" description:"Concurrency of transfers" ini-name:"concurrency"`
	PartSize    int64  `long:"partsize" short:"s" description:"Initial size of concurrent parts, in bytes" default:"20971520" ini-name:"partsize"`
	EndPoint    string `long:"endpoint" description:"Amazon S3 endpoint" default:"s3.amazonaws.com" ini-name:"endpoint"`
	Debug       bool   `long:"debug" description:"Enable debug logging." ini-name:"debug"`
}

var appOpts struct {
	Version  func() `long:"version" short:"v" description:"Print version"`
	Man      func() `long:"manpage" short:"m" description:"Create gof3r.man man page in current directory"`
	WriteIni bool   `long:"writeini" short:"i" description:"Write .gof3r.ini in current user's home directory" no-ini:"true"`
}
var parser = flags.NewParser(&appOpts, (flags.HelpFlag | flags.PassDoubleDash))

func init() {

	// set parser fields
	parser.ShortDescription = "streaming, concurrent s3 client"

	appOpts.Version = func() {
		fmt.Fprintf(os.Stderr, "%s version %s\n", name, version)
		os.Exit(0)
	}

	appOpts.Man = func() {
		f, err := os.Create(name + ".man")
		if err != nil {
			log.Fatal(err)
		}
		parser.WriteManPage(f)
		fmt.Fprintf(os.Stderr, "man page written to %s\n", f.Name())
		os.Exit(0)
	}
}

func iniPath() (path string, exist bool, err error) {
	usr, err := user.Current()
	if err != nil {
		return
	}
	path = fmt.Sprintf("%s/%s", usr.HomeDir, iniFile)
	if _, staterr := os.Stat(path); !os.IsNotExist(staterr) {
		exist = true
	}
	return
}

func parseIni() (err error) {
	p, exist, err := iniPath()
	if err != nil || !exist {
		return
	}
	return flags.NewIniParser(parser).ParseFile(p)
}

func writeIni() {
	p, exist, err := iniPath()
	if err != nil {
		log.Fatal(err)
	}
	if exist {
		fmt.Fprintf(os.Stderr, "%s exists, refusing to overwrite.\n", p)
	} else {
		if err := flags.NewIniParser(parser).WriteFile(p,
			(flags.IniIncludeComments | flags.IniIncludeDefaults | flags.IniCommentDefaults)); err != nil {
			log.Fatal(err)
		}
		fmt.Fprintf(os.Stderr, "ini file written to %s\n", p)
	}
	os.Exit(0)
}
