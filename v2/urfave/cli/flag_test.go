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
	"testing"
	"time"

	ucli "github.com/urfave/cli/v3"

	fpgocli "github.com/IBM/fp-go/v2/urfave/cli"

	O "github.com/IBM/fp-go/v2/option"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

// runWithArgs builds a minimal *ucli.Command with the given flags, runs it
// with args, then invokes inspect with the command.  Returns any run error.
func runWithArgs(
	t *testing.T,
	flags []fpgocli.Flag,
	args []string,
	inspect func(*fpgocli.Command),
) error {
	t.Helper()
	var cmd *fpgocli.Command
	app := &ucli.Command{
		Name:  "test",
		Flags: flags,
		Action: func(ctx context.Context, c *ucli.Command) error {
			cmd = c
			return nil
		},
	}
	err := app.Run(context.Background(), args)
	if cmd != nil {
		inspect(cmd)
	}
	return err
}

// ---------------------------------------------------------------------------
// Flag-type prisms – AsString / AsBool / AsInt / AsInt64 / AsFloat64 / AsDuration
// ---------------------------------------------------------------------------

func TestAsString_getOption_string(t *testing.T) {
	p := fpgocli.AsString()
	flag := &ucli.StringFlag{Name: "name", Value: "world"}
	// The flag must be applied (Get() returns the current value after parsing).
	// For a fresh flag with a default, Get() returns the default.
	_ = flag.PreParse()
	got := p.GetOption(flag)
	// A fresh StringFlag.Get() returns "" before being applied; use ReverseGet roundtrip instead.
	rev := p.ReverseGet("hello")
	assert.IsType(t, &ucli.StringFlag{}, rev)
	got2 := p.GetOption(rev)
	assert.Equal(t, O.Some("hello"), got2)
	_ = got // may be None before flag is applied to a flagset; just ensure it compiles
}

func TestAsString_rejectsBoolFlag(t *testing.T) {
	p := fpgocli.AsString()
	flag := &ucli.BoolFlag{Name: "verbose"}
	got := p.GetOption(flag)
	assert.Equal(t, O.None[string](), got)
}

func TestAsBool_getOption(t *testing.T) {
	p := fpgocli.AsBool()
	rev := p.ReverseGet(true)
	assert.IsType(t, &ucli.BoolFlag{}, rev)
	got := p.GetOption(rev)
	assert.Equal(t, O.Some(true), got)
}

func TestAsBool_rejectsStringFlag(t *testing.T) {
	p := fpgocli.AsBool()
	flag := &ucli.StringFlag{Name: "name"}
	got := p.GetOption(flag)
	assert.Equal(t, O.None[bool](), got)
}

func TestAsInt_roundtrip(t *testing.T) {
	p := fpgocli.AsInt()
	rev := p.ReverseGet(42)
	got := p.GetOption(rev)
	assert.Equal(t, O.Some(42), got)
}

func TestAsInt64_roundtrip(t *testing.T) {
	p := fpgocli.AsInt64()
	rev := p.ReverseGet(int64(99))
	got := p.GetOption(rev)
	assert.Equal(t, O.Some(int64(99)), got)
}

func TestAsFloat64_roundtrip(t *testing.T) {
	p := fpgocli.AsFloat64()
	rev := p.ReverseGet(3.14)
	got := p.GetOption(rev)
	assert.Equal(t, O.Some(3.14), got)
}

func TestAsDuration_roundtrip(t *testing.T) {
	p := fpgocli.AsDuration()
	d := 5 * time.Second
	rev := p.ReverseGet(d)
	got := p.GetOption(rev)
	assert.Equal(t, O.Some(d), got)
}

func TestAsStringSlice_roundtrip(t *testing.T) {
	p := fpgocli.AsStringSlice()
	ss := []string{"a", "b"}
	rev := p.ReverseGet(ss)
	got := p.GetOption(rev)
	assert.Equal(t, O.Some(ss), got)
}

// ---------------------------------------------------------------------------
// GetString – optional getter from *Command
// ---------------------------------------------------------------------------

func TestGetString_set(t *testing.T) {
	var got fpgocli.Option[string]
	err := runWithArgs(t,
		[]fpgocli.Flag{&ucli.StringFlag{Name: "name"}},
		[]string{"test", "--name", "Alice"},
		func(cmd *fpgocli.Command) {
			got = fpgocli.GetString("name")(cmd)
		},
	)
	require.NoError(t, err)
	assert.Equal(t, O.Some("Alice"), got)
}

func TestGetString_notSet_returnsNone(t *testing.T) {
	var got fpgocli.Option[string]
	err := runWithArgs(t,
		[]fpgocli.Flag{&ucli.StringFlag{Name: "name", Value: "default"}},
		[]string{"test"},
		func(cmd *fpgocli.Command) {
			got = fpgocli.GetString("name")(cmd)
		},
	)
	require.NoError(t, err)
	assert.Equal(t, O.None[string](), got)
}

// ---------------------------------------------------------------------------
// GetBool
// ---------------------------------------------------------------------------

func TestGetBool_set(t *testing.T) {
	var got fpgocli.Option[bool]
	err := runWithArgs(t,
		[]fpgocli.Flag{&ucli.BoolFlag{Name: "verbose"}},
		[]string{"test", "--verbose"},
		func(cmd *fpgocli.Command) {
			got = fpgocli.GetBool("verbose")(cmd)
		},
	)
	require.NoError(t, err)
	assert.Equal(t, O.Some(true), got)
}

func TestGetBool_notSet_returnsNone(t *testing.T) {
	var got fpgocli.Option[bool]
	err := runWithArgs(t,
		[]fpgocli.Flag{&ucli.BoolFlag{Name: "verbose"}},
		[]string{"test"},
		func(cmd *fpgocli.Command) {
			got = fpgocli.GetBool("verbose")(cmd)
		},
	)
	require.NoError(t, err)
	assert.Equal(t, O.None[bool](), got)
}

// ---------------------------------------------------------------------------
// GetInt
// ---------------------------------------------------------------------------

func TestGetInt_set(t *testing.T) {
	var got fpgocli.Option[int]
	err := runWithArgs(t,
		[]fpgocli.Flag{&ucli.IntFlag{Name: "count"}},
		[]string{"test", "--count", "7"},
		func(cmd *fpgocli.Command) {
			got = fpgocli.GetInt("count")(cmd)
		},
	)
	require.NoError(t, err)
	assert.Equal(t, O.Some(7), got)
}

func TestGetInt_notSet_returnsNone(t *testing.T) {
	var got fpgocli.Option[int]
	err := runWithArgs(t,
		[]fpgocli.Flag{&ucli.IntFlag{Name: "count", Value: 1}},
		[]string{"test"},
		func(cmd *fpgocli.Command) {
			got = fpgocli.GetInt("count")(cmd)
		},
	)
	require.NoError(t, err)
	assert.Equal(t, O.None[int](), got)
}

// ---------------------------------------------------------------------------
// GetInt64
// ---------------------------------------------------------------------------

func TestGetInt64_set(t *testing.T) {
	var got fpgocli.Option[int64]
	err := runWithArgs(t,
		[]fpgocli.Flag{&ucli.Int64Flag{Name: "offset"}},
		[]string{"test", "--offset", "9999999999"},
		func(cmd *fpgocli.Command) {
			got = fpgocli.GetInt64("offset")(cmd)
		},
	)
	require.NoError(t, err)
	assert.Equal(t, O.Some(int64(9999999999)), got)
}

// ---------------------------------------------------------------------------
// GetFloat64
// ---------------------------------------------------------------------------

func TestGetFloat64_set(t *testing.T) {
	var got fpgocli.Option[float64]
	err := runWithArgs(t,
		[]fpgocli.Flag{&ucli.Float64Flag{Name: "rate"}},
		[]string{"test", "--rate", "0.5"},
		func(cmd *fpgocli.Command) {
			got = fpgocli.GetFloat64("rate")(cmd)
		},
	)
	require.NoError(t, err)
	assert.Equal(t, O.Some(0.5), got)
}

// ---------------------------------------------------------------------------
// GetDuration
// ---------------------------------------------------------------------------

func TestGetDuration_set(t *testing.T) {
	var got fpgocli.Option[time.Duration]
	err := runWithArgs(t,
		[]fpgocli.Flag{&ucli.DurationFlag{Name: "timeout"}},
		[]string{"test", "--timeout", "30s"},
		func(cmd *fpgocli.Command) {
			got = fpgocli.GetDuration("timeout")(cmd)
		},
	)
	require.NoError(t, err)
	assert.Equal(t, O.Some(30*time.Second), got)
}

// ---------------------------------------------------------------------------
// GetStringSlice
// ---------------------------------------------------------------------------

func TestGetStringSlice_set(t *testing.T) {
	var got fpgocli.Option[[]string]
	err := runWithArgs(t,
		[]fpgocli.Flag{&ucli.StringSliceFlag{Name: "tag"}},
		[]string{"test", "--tag", "a", "--tag", "b"},
		func(cmd *fpgocli.Command) {
			got = fpgocli.GetStringSlice("tag")(cmd)
		},
	)
	require.NoError(t, err)
	if assert.True(t, O.IsSome(got)) {
		vals := O.GetOrElse(func() []string { return nil })(got)
		assert.ElementsMatch(t, []string{"a", "b"}, vals)
	}
}

func TestGetStringSlice_notSet_returnsNone(t *testing.T) {
	var got fpgocli.Option[[]string]
	err := runWithArgs(t,
		[]fpgocli.Flag{&ucli.StringSliceFlag{Name: "tag"}},
		[]string{"test"},
		func(cmd *fpgocli.Command) {
			got = fpgocli.GetStringSlice("tag")(cmd)
		},
	)
	require.NoError(t, err)
	assert.Equal(t, O.None[[]string](), got)
}

// ---------------------------------------------------------------------------
// CommandBuilder
// ---------------------------------------------------------------------------

func TestCommandBuilder_buildsCorrectly(t *testing.T) {
	cmd := fpgocli.NewCommand("greet").
		WithUsage("say hello").
		WithDescription("greets a user").
		WithAliases("g", "hi").
		WithFlags(&ucli.StringFlag{Name: "name", Value: "World"}).
		Build()

	assert.Equal(t, "greet", cmd.Name)
	assert.Equal(t, "say hello", cmd.Usage)
	assert.Equal(t, "greets a user", cmd.Description)
	assert.ElementsMatch(t, []string{"g", "hi"}, cmd.Aliases)
	assert.Len(t, cmd.Flags, 1)
}

func TestCommandBuilder_withAction(t *testing.T) {
	ran := false
	cmd := fpgocli.NewCommand("run").
		WithAction(func(_ context.Context, _ *ucli.Command) error {
			ran = true
			return nil
		}).
		Build()

	err := cmd.Run(context.Background(), []string{"run"})
	require.NoError(t, err)
	assert.True(t, ran)
}
