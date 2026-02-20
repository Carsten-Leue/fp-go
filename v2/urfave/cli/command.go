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

import ucli "github.com/urfave/cli/v3"

// CommandBuilder is a fluent builder for a *[Command].
// Start with [NewCommand] and finish with [CommandBuilder.Build].
type CommandBuilder struct {
	cmd *ucli.Command
}

// NewCommand starts a [CommandBuilder] with the given name.
//
// Example:
//
//	cmd := cli.NewCommand("greet").
//	    WithUsage("say hello").
//	    WithFlags(&cli.StringFlag{Name: "name", Value: "World"}).
//	    WithIOAction(greetAction).
//	    Build()
func NewCommand(name string) *CommandBuilder {
	return &CommandBuilder{cmd: &ucli.Command{Name: name}}
}

// WithUsage sets the one-line usage description.
func (b *CommandBuilder) WithUsage(usage string) *CommandBuilder {
	b.cmd.Usage = usage
	return b
}

// WithDescription sets the long description shown in help text.
func (b *CommandBuilder) WithDescription(desc string) *CommandBuilder {
	b.cmd.Description = desc
	return b
}

// WithAliases adds alternative names for the command.
func (b *CommandBuilder) WithAliases(aliases ...string) *CommandBuilder {
	b.cmd.Aliases = append(b.cmd.Aliases, aliases...)
	return b
}

// WithFlags appends flags to the command.
func (b *CommandBuilder) WithFlags(flags ...Flag) *CommandBuilder {
	b.cmd.Flags = append(b.cmd.Flags, flags...)
	return b
}

// WithCommands appends sub-commands.
func (b *CommandBuilder) WithCommands(cmds ...*Command) *CommandBuilder {
	b.cmd.Commands = append(b.cmd.Commands, cmds...)
	return b
}

// WithAction attaches a plain urfave/cli [ActionFunc] to the command.
func (b *CommandBuilder) WithAction(action ActionFunc) *CommandBuilder {
	b.cmd.Action = action
	return b
}

// WithIOAction attaches an fp-go [IOAction][Void] to the command, wrapping it
// with [ToActionFunc] so it integrates seamlessly with urfave/cli.
func (b *CommandBuilder) WithIOAction(action IOAction[Void]) *CommandBuilder {
	b.cmd.Action = ToActionFunc(action)
	return b
}

// WithIOActionOf attaches an [IOAction][A] to the command.
// The success value is discarded; only errors are propagated.
func WithIOActionOf[A any](action IOAction[A]) func(*CommandBuilder) *CommandBuilder {
	return func(b *CommandBuilder) *CommandBuilder {
		b.cmd.Action = ToActionFunc(action)
		return b
	}
}

// Build returns the fully assembled *[Command].
func (b *CommandBuilder) Build() *Command {
	return b.cmd
}
