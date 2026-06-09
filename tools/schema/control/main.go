package main

import (
	"fmt"
	"io"
	"os"

	"ariga.io/atlas-provider-gorm/gormschema"
	"github.com/glueops/waggle/internal/models/control"
)

func main() {
	stmts, err := gormschema.New("postgres").Load(
		&control.Account{},
		&control.AccountEmail{},
		&control.AuthAuditEvent{},
		&control.Organization{},
		&control.OrgAPIKey{},
		&control.TokenSession{},
		&control.User{},
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load gorm schema: %v\n", err)
		os.Exit(1)
	}
	io.WriteString(os.Stdout, stmts)
}
