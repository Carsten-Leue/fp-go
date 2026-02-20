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
	"context"

	ucli "github.com/urfave/cli/v3"

	"github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/ioresult"
	O "github.com/IBM/fp-go/v2/option"
	P "github.com/IBM/fp-go/v2/optics/prism"
	"github.com/IBM/fp-go/v2/result"
)

type (
	// Command is the urfave/cli/v3 Command type.
	Command = ucli.Command

	// Flag is the urfave/cli/v3 Flag interface.
	Flag = ucli.Flag

	// ActionFunc is the native urfave/cli/v3 action signature.
	ActionFunc = ucli.ActionFunc

	// Result is a computation that may fail with an error.
	Result[A any] = result.Result[A]

	// IOResult is a lazy IO computation that may fail with an error.
	IOResult[A any] = ioresult.IOResult[A]

	// Option is an optional value.
	Option[A any] = O.Option[A]

	// Prism is an optic for focusing on a variant of a sum type.
	Prism[S, A any] = P.Prism[S, A]

	// Void is the unit type – a type with exactly one value.
	Void = function.Void

	// IOAction is an fp-go style action that returns a lazy IOResult[A].
	// It is the fp-go equivalent of [ActionFunc] but richer: the error is
	// captured inside the returned IOResult rather than returned directly,
	// which enables composition through the standard fp-go combinators.
	//
	// Convert to a plain [ActionFunc] with [ToActionFunc], or promote an
	// existing [ActionFunc] with [FromActionFunc].
	IOAction[A any] func(context.Context, *Command) IOResult[A]
)
