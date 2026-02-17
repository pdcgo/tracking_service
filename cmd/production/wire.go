//go:build wireinject
// +build wireinject

package main

import (
	"github.com/urfave/cli/v3"
)

func InitializeApp() (App, error) {

	return &cli.Command{}, nil
}
