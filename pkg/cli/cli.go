package cli

import (
	"flag"
	"github.com/ryotarai/github-api-authz-proxy/pkg/handler"
	"github.com/ryotarai/github-api-authz-proxy/pkg/opa"
	"log"
	"net/http"
	"net/url"
)

type CLI struct {
}

func New() *CLI {
	return &CLI{}
}

func (c *CLI) Start(args []string) error {
	fs := flag.NewFlagSet(args[0], flag.ContinueOnError)
	listen := fs.String("listen", ":8080", "")
	originURLStr := fs.String("origin-url", "", "")
	opaServerURLStr := fs.String("opa-server-url", "", "")
	accessToken := fs.String("access-token", "", "")

	err := fs.Parse(args[1:])
	if err != nil {
		return err
	}

	originURL, err := url.Parse(*originURLStr)
	if err != nil {
		return err
	}

	opaServerURL, err := url.Parse(*opaServerURLStr)
	if err != nil {
		return err
	}

	authz := opa.NewClient(opaServerURL)

	h, err := handler.New(originURL, *accessToken, authz)
	if err != nil {
		return err
	}

	log.Printf("INFO: Listening on %s", *listen)
	err = http.ListenAndServe(*listen, h)
	if err != nil {
		return err
	}

	return nil
}
