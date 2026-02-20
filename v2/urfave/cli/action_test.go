// Copyright (c) 2025 IBM Corp.
// All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cli_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	ucli "github.com/urfave/cli/v3"

	fpgocli "github.com/IBM/fp-go/v2/urfave/cli"

	"github.com/IBM/fp-go/v2/ioresult"
	"github.com/IBM/fp-go/v2/result"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ---------------------------------------------------------------------------
// Of / Fail
// ---------------------------------------------------------------------------

func TestOf_returnsValue(t *testing.T) {
	action := fpgocli.Of[int](42)
	res := action(context.Background(), &ucli.Command{})()
	assert.True(t, result.IsRight(res))
	result.MonadFold(res,
		func(err error) struct{} { t.Fatalf("unexpected error: %v", err); return struct{}{} },
		func(n int) struct{} { assert.Equal(t, 42, n); return struct{}{} },
	)
}

func TestFail_returnsError(t *testing.T) {
	sentinel := errors.New("boom")
	action := fpgocli.Fail[int](sentinel)
	res := action(context.Background(), &ucli.Command{})()
	assert.True(t, result.IsLeft(res))
	result.MonadFold(res,
		func(err error) struct{} { assert.Equal(t, sentinel, err); return struct{}{} },
		func(n int) struct{} { t.Fatalf("unexpected success: %d", n); return struct{}{} },
	)
}

// ---------------------------------------------------------------------------
// ToActionFunc
// ---------------------------------------------------------------------------

func TestToActionFunc_success(t *testing.T) {
	action := fpgocli.Of[string]("ok")
	af := fpgocli.ToActionFunc(action)

	app := &ucli.Command{Action: af}
	err := app.Run(context.Background(), []string{"app"})
	require.NoError(t, err)
}

func TestToActionFunc_failure(t *testing.T) {
	sentinel := errors.New("cli error")
	action := fpgocli.Fail[fpgocli.Void](sentinel)
	af := fpgocli.ToActionFunc(action)

	app := &ucli.Command{Action: af}
	err := app.Run(context.Background(), []string{"app"})
	require.ErrorIs(t, err, sentinel)
}

// ---------------------------------------------------------------------------
// FromActionFunc
// ---------------------------------------------------------------------------

func TestFromActionFunc_success(t *testing.T) {
	called := false
	af := func(_ context.Context, _ *ucli.Command) error {
		called = true
		return nil
	}
	action := fpgocli.FromActionFunc(af)
	res := action(context.Background(), &ucli.Command{})()
	assert.True(t, result.IsRight(res))
	assert.True(t, called)
}

func TestFromActionFunc_failure(t *testing.T) {
	sentinel := errors.New("native error")
	af := func(_ context.Context, _ *ucli.Command) error { return sentinel }
	action := fpgocli.FromActionFunc(af)
	res := action(context.Background(), &ucli.Command{})()
	assert.True(t, result.IsLeft(res))
}

// ---------------------------------------------------------------------------
// Map
// ---------------------------------------------------------------------------

func TestMap_transformsSuccess(t *testing.T) {
	action := fpgocli.Of[string]("hello")
	upper := fpgocli.Map[string, string](strings.ToUpper)(action)
	res := upper(context.Background(), &ucli.Command{})()
	result.MonadFold(res,
		func(err error) struct{} { t.Fatalf("unexpected error: %v", err); return struct{}{} },
		func(s string) struct{} { assert.Equal(t, "HELLO", s); return struct{}{} },
	)
}

func TestMap_propagatesError(t *testing.T) {
	sentinel := errors.New("oops")
	action := fpgocli.Fail[string](sentinel)
	upper := fpgocli.Map[string, string](strings.ToUpper)(action)
	res := upper(context.Background(), &ucli.Command{})()
	assert.True(t, result.IsLeft(res))
}

// ---------------------------------------------------------------------------
// Chain
// ---------------------------------------------------------------------------

func TestChain_sequences(t *testing.T) {
	action := fpgocli.Of[int](3)
	doubled := fpgocli.Chain[int, int](func(n int) fpgocli.IOAction[int] {
		return fpgocli.Of[int](n * 2)
	})(action)
	res := doubled(context.Background(), &ucli.Command{})()
	result.MonadFold(res,
		func(err error) struct{} { t.Fatalf("unexpected error: %v", err); return struct{}{} },
		func(n int) struct{} { assert.Equal(t, 6, n); return struct{}{} },
	)
}

func TestChain_shortCircuitsOnError(t *testing.T) {
	sentinel := errors.New("abort")
	action := fpgocli.Fail[int](sentinel)
	called := false
	chained := fpgocli.Chain[int, int](func(n int) fpgocli.IOAction[int] {
		called = true
		return fpgocli.Of[int](n * 2)
	})(action)
	res := chained(context.Background(), &ucli.Command{})()
	assert.True(t, result.IsLeft(res))
	assert.False(t, called, "chain continuation must not run after failure")
}

// ---------------------------------------------------------------------------
// MapError
// ---------------------------------------------------------------------------

func TestMapError_transformsError(t *testing.T) {
	sentinel := errors.New("raw")
	action := fpgocli.Fail[int](sentinel)
	wrapped := fpgocli.MapError[int](func(err error) error {
		return errors.New("wrapped: " + err.Error())
	})(action)
	res := wrapped(context.Background(), &ucli.Command{})()
	result.MonadFold(res,
		func(err error) struct{} {
			assert.Equal(t, "wrapped: raw", err.Error())
			return struct{}{}
		},
		func(_ int) struct{} { t.Fatal("expected error"); return struct{}{} },
	)
}

func TestMapError_leavesSuccessUnchanged(t *testing.T) {
	action := fpgocli.Of[int](7)
	wrapped := fpgocli.MapError[int](func(err error) error {
		return errors.New("should not happen")
	})(action)
	res := wrapped(context.Background(), &ucli.Command{})()
	result.MonadFold(res,
		func(err error) struct{} { t.Fatalf("unexpected error: %v", err); return struct{}{} },
		func(n int) struct{} { assert.Equal(t, 7, n); return struct{}{} },
	)
}

// ---------------------------------------------------------------------------
// CommandBuilder with IOAction
// ---------------------------------------------------------------------------

func TestCommandBuilder_withIOAction(t *testing.T) {
	ran := false
	action := fpgocli.IOAction[fpgocli.Void](func(_ context.Context, _ *fpgocli.Command) fpgocli.IOResult[fpgocli.Void] {
		return ioresult.Of[fpgocli.Void](fpgocli.Void{})
	})
	_ = action // suppress unused warning
	ran = true

	cmd := fpgocli.NewCommand("test").
		WithUsage("a test command").
		WithIOAction(fpgocli.IOAction[fpgocli.Void](func(_ context.Context, _ *fpgocli.Command) fpgocli.IOResult[fpgocli.Void] {
			ran = true
			return ioresult.Of[fpgocli.Void](fpgocli.Void{})
		})).
		Build()

	err := cmd.Run(context.Background(), []string{"test"})
	require.NoError(t, err)
	assert.True(t, ran)
}
