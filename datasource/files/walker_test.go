package files

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"
)

func TestExtract(t *testing.T) {
	// t.Skip("only local testing")

	c := make(chan File)

	wg, ctx := errgroup.WithContext(context.TODO())

	wg.Go(func() error {
		return Extract(ctx, "/Users/nvrusin/black", c)
	})

	wg.Go(func() error {
		for f := range c {
			t.Log("file", f.Package, f.Name, f.Lines, f.Symbols)
		}

		return nil
	})

	require.NoError(t, wg.Wait())
}

func TestExtractImports(t *testing.T) {
	t.Run("go", func(t *testing.T) {
		content := `
package dashboard

import (
	"bytes"
	_ "embed"
	"html/template"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/templates"
	"gorm.io/gorm"

	"devex_dashboard/project"
)

import "errors"
import "devex_dashboard/test"

//go:embed form.gohtml
var form string

type Params struct {
	ProjectIDs []project.ID 'form:"project_ids"'
	PerFiles   bool         'form:"per_files"'
	Filter     string       'form:"filter"'
}
`

		expect := []string{
			"bytes",
			"embed",
			"html/template",
			"net/http",
			"strings",
			"github.com/gin-gonic/gin",
			"github.com/go-echarts/go-echarts/v2/components",
			"github.com/go-echarts/go-echarts/v2/templates",
			"gorm.io/gorm",
			"devex_dashboard/project",
			"errors",
			"devex_dashboard/test",
		}

		assert.Equal(t, expect, extractImports(strings.Split(content, "\n")))
	})

	t.Run("python", func(t *testing.T) {
		content := `
"""
Formatting many files at once via multiprocessing. Contains entrypoint and utilities.

NOTE: this module is only imported if we need to format several files at once.
"""

import asyncio
import logging
import os
import signal
import sys
from concurrent.futures import Executor, ProcessPoolExecutor, ThreadPoolExecutor
from multiprocessing import Manager
from pathlib import Path
from typing import Any, Iterable, Optional, Set

from mypy_extensions import mypyc_attr

from black import WriteBack, format_file_in_place
from black.cache import Cache, filter_cached, read_cache, write_cache
from black.mode import Mode
from black.output import err
from black.report import Changed, Report


def maybe_install_uvloop() -> None:
    """If our environment has uvloop installed we use it.

    This is called only from command-line entry points to avoid
    interfering with the parent process if Black is used as a library.
    """
    try:
        import uvloop

        uvloop.install()
    except ImportError:
        pass
`

		expect := []string{
			"asyncio",
			"logging",
			"os",
			"signal",
			"sys",
			"concurrent.futures",
			"multiprocessing",
			"pathlib",
			"typing",
			"mypy_extensions",
			"black",
			"black.cache",
			"black.mode",
			"black.output",
			"black.report",
			"uvloop",
		}

		assert.Equal(t, expect, extractImports(strings.Split(content, "\n")))
	})
}
