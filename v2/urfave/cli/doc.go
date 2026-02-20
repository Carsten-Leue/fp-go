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

// Package cli provides fp-go functional wrappers for the [github.com/urfave/cli/v3] package.
//
// # Actions as Effects
//
// The standard urfave/cli [ActionFunc] has the signature:
//
//	func(context.Context, *Command) error
//
// This package introduces [IOAction], which lifts that signature into an [IOResult]:
//
//	IOAction[A] = func(context.Context, *Command) IOResult[A]
//
// An [IOAction] is a lazy, referentially-transparent description of a computation
// that may fail. Use [ToActionFunc] to wire it back into urfave's Command.Action field,
// and [FromActionFunc] to promote an existing ActionFunc.
//
// # Prisms for Flag Types
//
// [AsString], [AsBool], [AsInt], [AsInt64], [AsFloat64] and [AsDuration] each return a
// [Prism] over the [Flag] interface (the sum type). They let you safely inspect
// an arbitrary flag value, returning [Option][T] on the GetOption path and
// constructing a concrete flag on the ReverseGet path.
//
// # Optional Getters from *Command
//
// [GetString], [GetBool], [GetInt], [GetInt64], [GetFloat64], [GetDuration] and
// [GetStringSlice] each return a curried getter
//
//	func(*Command) Option[T]
//
// that returns [option.Some] only when the flag has been explicitly set by the
// caller, and [option.None] otherwise. This preserves the distinction between
// "flag not provided" and "flag set to its default value".
//
// # Functional Command Building
//
// [NewCommand] starts a [CommandBuilder] that lets you assemble a *[Command] in a
// pipeline style, attaching flags, sub-commands, and an [IOAction] in one
// expression.
package cli
