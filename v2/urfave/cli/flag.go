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

package cli

import (
	"time"

	ucli "github.com/urfave/cli/v3"

	O "github.com/IBM/fp-go/v2/option"
	P "github.com/IBM/fp-go/v2/optics/prism"
)

// ---------------------------------------------------------------------------
// Prisms on the Flag interface (sum-type discrimination)
//
// cli.Flag is an interface implemented by StringFlag, BoolFlag, IntFlag, etc.
// Each prism below focuses on one variant:
//   GetOption: extracts the typed default value if the underlying flag
//              holds that type, returns None otherwise.
//   ReverseGet: constructs a minimal concrete Flag from a typed value.
// ---------------------------------------------------------------------------

// AsString returns a Prism that focuses on the string default value of a Flag.
// GetOption returns Some(s) when the flag holds a string value, None otherwise.
// ReverseGet wraps a string into a *[ucli.StringFlag] with that default.
func AsString() Prism[Flag, string] {
	return P.MakePrismWithName(
		func(f Flag) O.Option[string] {
			if s, ok := f.Get().(string); ok {
				return O.Some(s)
			}
			return O.None[string]()
		},
		func(s string) Flag {
			return &ucli.StringFlag{Value: s}
		},
		"AsStringFlag",
	)
}

// AsBool returns a Prism that focuses on the bool default value of a Flag.
// GetOption returns Some(b) when the flag holds a bool value, None otherwise.
// ReverseGet wraps a bool into a *[ucli.BoolFlag] with that default.
func AsBool() Prism[Flag, bool] {
	return P.MakePrismWithName(
		func(f Flag) O.Option[bool] {
			if b, ok := f.Get().(bool); ok {
				return O.Some(b)
			}
			return O.None[bool]()
		},
		func(b bool) Flag {
			return &ucli.BoolFlag{Value: b}
		},
		"AsBoolFlag",
	)
}

// AsInt returns a Prism that focuses on the int default value of a Flag.
func AsInt() Prism[Flag, int] {
	return P.MakePrismWithName(
		func(f Flag) O.Option[int] {
			if n, ok := f.Get().(int); ok {
				return O.Some(n)
			}
			return O.None[int]()
		},
		func(n int) Flag {
			return &ucli.IntFlag{Value: n}
		},
		"AsIntFlag",
	)
}

// AsInt64 returns a Prism that focuses on the int64 default value of a Flag.
func AsInt64() Prism[Flag, int64] {
	return P.MakePrismWithName(
		func(f Flag) O.Option[int64] {
			if n, ok := f.Get().(int64); ok {
				return O.Some(n)
			}
			return O.None[int64]()
		},
		func(n int64) Flag {
			return &ucli.Int64Flag{Value: n}
		},
		"AsInt64Flag",
	)
}

// AsFloat64 returns a Prism that focuses on the float64 default value of a Flag.
func AsFloat64() Prism[Flag, float64] {
	return P.MakePrismWithName(
		func(f Flag) O.Option[float64] {
			if v, ok := f.Get().(float64); ok {
				return O.Some(v)
			}
			return O.None[float64]()
		},
		func(v float64) Flag {
			return &ucli.Float64Flag{Value: v}
		},
		"AsFloat64Flag",
	)
}

// AsDuration returns a Prism that focuses on the time.Duration default value of a Flag.
func AsDuration() Prism[Flag, time.Duration] {
	return P.MakePrismWithName(
		func(f Flag) O.Option[time.Duration] {
			if d, ok := f.Get().(time.Duration); ok {
				return O.Some(d)
			}
			return O.None[time.Duration]()
		},
		func(d time.Duration) Flag {
			return &ucli.DurationFlag{Value: d}
		},
		"AsDurationFlag",
	)
}

// AsStringSlice returns a Prism that focuses on the []string default value of a Flag.
func AsStringSlice() Prism[Flag, []string] {
	return P.MakePrismWithName(
		func(f Flag) O.Option[[]string] {
			if ss, ok := f.Get().([]string); ok {
				return O.Some(ss)
			}
			return O.None[[]string]()
		},
		func(ss []string) Flag {
			return &ucli.StringSliceFlag{Value: ss}
		},
		"AsStringSliceFlag",
	)
}

// ---------------------------------------------------------------------------
// Optional getters from *Command
//
// Each getter is a curried function  name → (*Command → Option[T]).
// It returns Some(value) only when the flag was explicitly set by the user
// (cmd.IsSet returns true), and None otherwise.  This lets callers
// distinguish "flag absent" from "flag set to the zero value".
// ---------------------------------------------------------------------------

// GetString returns a function that extracts the named string flag value from
// a *Command, returning Some(s) when the flag was explicitly set, None otherwise.
func GetString(name string) func(*Command) Option[string] {
	return func(cmd *Command) Option[string] {
		if cmd.IsSet(name) {
			return O.Some(cmd.String(name))
		}
		return O.None[string]()
	}
}

// GetBool returns a function that extracts the named bool flag value from a
// *Command, returning Some(b) when the flag was explicitly set.
func GetBool(name string) func(*Command) Option[bool] {
	return func(cmd *Command) Option[bool] {
		if cmd.IsSet(name) {
			return O.Some(cmd.Bool(name))
		}
		return O.None[bool]()
	}
}

// GetInt returns a function that extracts the named int flag value from a
// *Command, returning Some(n) when the flag was explicitly set.
func GetInt(name string) func(*Command) Option[int] {
	return func(cmd *Command) Option[int] {
		if cmd.IsSet(name) {
			return O.Some(cmd.Int(name))
		}
		return O.None[int]()
	}
}

// GetInt64 returns a function that extracts the named int64 flag value from a
// *Command, returning Some(n) when the flag was explicitly set.
func GetInt64(name string) func(*Command) Option[int64] {
	return func(cmd *Command) Option[int64] {
		if cmd.IsSet(name) {
			return O.Some(cmd.Int64(name))
		}
		return O.None[int64]()
	}
}

// GetFloat64 returns a function that extracts the named float64 flag value from
// a *Command, returning Some(v) when the flag was explicitly set.
func GetFloat64(name string) func(*Command) Option[float64] {
	return func(cmd *Command) Option[float64] {
		if cmd.IsSet(name) {
			return O.Some(cmd.Float64(name))
		}
		return O.None[float64]()
	}
}

// GetDuration returns a function that extracts the named Duration flag value
// from a *Command, returning Some(d) when the flag was explicitly set.
func GetDuration(name string) func(*Command) Option[time.Duration] {
	return func(cmd *Command) Option[time.Duration] {
		if cmd.IsSet(name) {
			return O.Some(cmd.Duration(name))
		}
		return O.None[time.Duration]()
	}
}

// GetStringSlice returns a function that extracts the named string-slice flag
// value from a *Command, returning Some(ss) when the flag was explicitly set.
func GetStringSlice(name string) func(*Command) Option[[]string] {
	return func(cmd *Command) Option[[]string] {
		if cmd.IsSet(name) {
			return O.Some(cmd.StringSlice(name))
		}
		return O.None[[]string]()
	}
}
