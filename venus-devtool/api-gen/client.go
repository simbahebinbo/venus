package main

import (
	"bytes"
	"fmt"
	"log"
	"text/template"

	"github.com/filecoin-project/venus/venus-devtool/util"
	"github.com/urfave/cli/v2"
)

var clientCmd = &cli.Command{
	Name:  "client",
	Flags: []cli.Flag{},
	Action: func(cctx *cli.Context) error {
		for _, target := range apiTargets {
			err := genClientForAPI(target)
			if err != nil {
				log.Fatalf("got error while generating client codes for %s: %s", target.Type, err)
			}
		}
		return nil
	},
}

const clientGenTemplate = `
// Code generated by github.com/filecoin-project/venus/venus-devtool/api-gen. DO NOT EDIT.
package {{ .PkgName }}

import (
	"context"
	"fmt"
	"net/http"

	"github.com/filecoin-project/go-jsonrpc"

	"github.com/filecoin-project/venus/venus-shared/api"
)

const MajorVersion = {{ .MajorVersion }}
const APINamespace = "{{ .APINs }}"
const MethodNamespace = "{{ .MethNs }}"

// New{{ .APIName }}RPC creates a new httpparse jsonrpc remotecli.
func New{{ .APIName }}RPC(ctx context.Context, addr string, requestHeader http.Header, opts ...jsonrpc.Option) ({{ .APIName }}, jsonrpc.ClientCloser, error) {
	endpoint, err := api.Endpoint(addr, MajorVersion)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid addr %s: %w", addr, err)
	}

	if requestHeader == nil {
		requestHeader = http.Header{}
	}
	requestHeader.Set(api.VenusAPINamespaceHeader, APINamespace) 

	var res {{ .APIStruct }}
	closer, err := jsonrpc.NewMergeClient(ctx, endpoint, MethodNamespace, api.GetInternalStructs(&res), requestHeader, opts...)

	return &res, closer, err
}
`

func genClientForAPI(t util.APIMeta) error {
	ifaceMetas, astMeta, err := util.ParseInterfaceMetas(t.ParseOpt)
	if err != nil {
		return err
	}

	apiName := t.Type.Name()

	var apiIface *util.InterfaceMeta
	for i := range ifaceMetas {
		if ifaceMetas[i].Name == apiName {
			apiIface = ifaceMetas[i]
			break
		}
	}

	if apiIface == nil {
		return fmt.Errorf("api %s not found", apiName)
	}

	tmpl, err := template.New("client").Parse(clientGenTemplate)
	if err != nil {
		return fmt.Errorf("parse template: %w", err)
	}

	ns := t.RPCMeta.Namespace
	if ns == "" {
		ns = fmt.Sprintf("%s.%s", apiIface.Pkg.Name, apiIface.Name)
	}

	methNs := t.RPCMeta.MethodNamespace
	if methNs == "" {
		methNs = "Filecoin"
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, map[string]interface{}{
		"PkgName":      apiIface.Pkg.Name,
		"APIName":      apiName,
		"APIStruct":    structName(apiName),
		"APINs":        ns,
		"MethNs":       methNs,
		"MajorVersion": t.RPCMeta.Version,
	})
	if err != nil {
		return fmt.Errorf("exec template: %w", err)
	}

	return outputSourceFile(astMeta.Location, "client_gen.go", &buf)
}