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

	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/ioresult"
	"github.com/IBM/fp-go/v2/result"
)

// Of creates an IOAction that always succeeds with value a, ignoring the
// context and command.
//
//go:inline
func Of[A any](a A) IOAction[A] {
	return func(_ context.Context, _ *Command) IOResult[A] {
		return ioresult.Of[A](a)
	}
}

// Fail creates an IOAction that always fails with err.
//
//go:inline
func Fail[A any](err error) IOAction[A] {
	return func(_ context.Context, _ *Command) IOResult[A] {
		return ioresult.Left[A](err)
	}
}

// ToActionFunc converts an fp-go [IOAction][A] into a urfave/cli [ActionFunc].
//
// The IOAction is executed eagerly when urfave/cli invokes the action. Any
// error captured inside the [IOResult] is returned directly; a successful
// result is discarded (only the side effect matters in the context of a CLI).
func ToActionFunc[A any](f IOAction[A]) ActionFunc {
	return func(ctx context.Context, cmd *ucli.Command) error {
		res := f(ctx, cmd)()
		return result.MonadFold(res,
			F.Identity[error],
			func(_ A) error { return nil },
		)
	}
}

// FromActionFunc promotes a plain urfave/cli [ActionFunc] to an [IOAction][Void].
//
// The resulting IOAction wraps the ActionFunc call: on success it returns
// [Void], on failure it captures the error inside the [IOResult].
func FromActionFunc(f ActionFunc) IOAction[Void] {
	return func(ctx context.Context, cmd *Command) IOResult[Void] {
		return func() Result[Void] {
			if err := f(ctx, cmd); err != nil {
				return result.Left[Void](err)
			}
			return result.Right[Void](F.VOID)
		}
	}
}

// Map transforms the success value of an IOAction without changing its
// error channel, producing a new IOAction[B].
func Map[A, B any](f func(A) B) func(IOAction[A]) IOAction[B] {
	return func(fa IOAction[A]) IOAction[B] {
		return func(ctx context.Context, cmd *Command) IOResult[B] {
			return ioresult.Map[A, B](f)(fa(ctx, cmd))
		}
	}
}

// Chain sequences two IOActions: the second receives the success value of
// the first. If the first fails the chain is short-circuited.
func Chain[A, B any](f func(A) IOAction[B]) func(IOAction[A]) IOAction[B] {
	return func(fa IOAction[A]) IOAction[B] {
		return func(ctx context.Context, cmd *Command) IOResult[B] {
			return ioresult.Chain[A, B](func(a A) IOResult[B] {
				return f(a)(ctx, cmd)
			})(fa(ctx, cmd))
		}
	}
}

// MapError transforms the error of a failing IOAction, leaving success
// values untouched.
func MapError[A any](f func(error) error) func(IOAction[A]) IOAction[A] {
	return func(fa IOAction[A]) IOAction[A] {
		return func(ctx context.Context, cmd *Command) IOResult[A] {
			return func() Result[A] {
				return result.MonadFold(fa(ctx, cmd)(),
					func(err error) Result[A] { return result.Left[A](f(err)) },
					result.Right[A],
				)
			}
		}
	}
}
