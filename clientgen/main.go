package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/luno/jettison/errors"
	"github.com/luno/luno-go/clientgen/clientgen"
	"github.com/luno/luno-go/clientgen/golang"
	"github.com/luno/luno-go/clientgen/php"
	"github.com/luno/luno-go/clientgen/python"
)

var languages = map[string]clientgen.Generator{
	"golang": golang.Generate,
	"php":    php.Generate,
	"python": python.Generate,
}

type options struct {
	spec   *string
	lang   *string
	outdir *string
	dryrun *bool
}

func main() {
	defaultSpec := "https://api.luno.com/api/exchange/1/schema"

	opts := options{
		spec:       flag.String("spec", defaultSpec, "URL of swagger spec"),
		lang:       flag.String("lang", "", "Language to generate"),
		outdir:     flag.String("outdir", "", "Directory to write the client to"),
		dryrun:     flag.Bool("dryrun", false, "Just output the generated files"),
	}

	flag.Parse()

	if err := run(opts); err != nil {
		log.Fatal(err)
	}
}

func run(opts options) error {
	spec, err := clientgen.LoadSpec(*opts.spec)
	if err != nil {
		return err
	}

	api := clientgen.ConvertSpec(spec)

	fl, err := generate(api, *opts.lang)
	if err != nil {
		return err
	}

	if *opts.dryrun {
		return dryrun(fl)
	}

	return output(fl, *opts.outdir)
}

func generate(api clientgen.API, lang string) ([]clientgen.File, error) {
	if lang == "" {
		return nil, errors.New("lang not specified")
	}

	generator := languages[lang]
	if generator == nil {
		return nil, errors.New("lang not implemented")
	}

	return generator(api)
}

func dryrun(fl []clientgen.File) error {
	for _, f := range fl {
		fmt.Println(string(f.Bytes()))
	}
	return nil
}

func output(fl []clientgen.File, outdir string) error {
	if outdir == "" {
		return errors.New("empty outdir")
	}

	for _, f := range fl {
		abspath := path.Join(outdir, path.Dir(f.RelPath()))
		if fi, err := os.Stat(abspath); os.IsNotExist(err) {
			if err := os.Mkdir(abspath, 0755); err != nil {
				return err
			}
		} else if err != nil {
			return errors.New("error checking if file path exists")
		} else if !fi.IsDir() {
			log.Println(abspath)
			return errors.New("file path isn't a directory")
		}

		absfile := path.Join(outdir, f.RelPath())
		if err := ioutil.WriteFile(absfile, f.Bytes(), 0644); err != nil {
			return err
		}
	}

	return nil
}
