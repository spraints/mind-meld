package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/spf13/cobra"

	"github.com/spraints/mind-meld/appcmd"
	"github.com/spraints/mind-meld/appcmd/fetch"
	"github.com/spraints/mind-meld/appcmd/watch"
	"github.com/spraints/mind-meld/apps/mindstormsapp"
	"github.com/spraints/mind-meld/githooks"
	"github.com/spraints/mind-meld/lmsdump"
	"github.com/spraints/mind-meld/lmsp"
	"github.com/spraints/mind-meld/ui"
)

func main() {
	finish(mkRootCmd().Execute())
}

func mkRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:   "mind-meld",
		Short: "Manage your LEGO MINDSTORMS",
	}

	root.AddCommand(mkBrowseCmd())
	root.AddCommand(mkDumpCmd())
	root.AddCommand(mkPreCommitCmd())

	root.AddCommand(mkAppSubcommandCmd("mindstorms", mindstormsapp.New()))

	return root
}

func mkBrowseCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "browse",
		Short: "Browse for a file and see inside of it!",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			var opts ui.Opts
			if len(args) == 1 {
				opts.Workdir = args[0]
			}
			return ui.Run(opts)
		},
	}
}

func mkDumpCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "dump",
		Short: "Print a plain text version of a mindstorms program.",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			return dump(args[0])
		},
	}
}

func mkPreCommitCmd() *cobra.Command {
	var cached bool
	cmd := &cobra.Command{
		Use:    "pre-commit",
		Args:   cobra.NoArgs,
		Hidden: true,
		RunE: func(*cobra.Command, []string) error {
			mode := githooks.UpdateWorkingCopy
			if cached {
				mode = githooks.UpdateCache
			}
			return githooks.RunPreCommit(mode)
		},
	}
	cmd.PersistentFlags().BoolVar(&cached, "cached", false, "update index instead of working copy")
	return cmd
}

func mkAppSubcommandCmd(name string, a appcmd.App) *cobra.Command {
	subCmd := &cobra.Command{
		Use:   name,
		Short: "Manage " + a.FullName() + " programs.",
	}

	subCmd.AddCommand(mkAppFetchCommand(a))
	subCmd.AddCommand(mkAppWatchCommand(a))

	return subCmd
}

func mkAppFetchCommand(a appcmd.App) *cobra.Command {
	return &cobra.Command{
		Use:   "fetch",
		Short: "Synchronize python programs between " + a.FullName() + " and a Git repository in the current directory.",
		Args:  cobra.NoArgs,
		RunE: func(*cobra.Command, []string) error {
			return fetch.Run(a)
		},
	}
}

func mkAppWatchCommand(a appcmd.App) *cobra.Command {
	return &cobra.Command{
		Use:   "watch",
		Short: "Continuously synchronize python programs between " + a.FullName() + " and a Git repository in the current directory.",
		Args:  cobra.NoArgs,
		RunE: func(*cobra.Command, []string) error {
			ctx := context.Background()
			ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
			defer cancel()
			return watch.Run(ctx, a)
		},
	}
}

func finish(err error) {
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
}

func dump(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	l, err := lmsp.ReadFile(f)
	if err != nil {
		return err
	}
	/*
		man, err := l.Manifest()
		if err != nil {
			return err
		}
		spew.Dump(man)
	*/

	proj, err := l.Project()
	if err != nil {
		return err
	}
	lmsdump.Dump(os.Stdout, proj)

	if os.Getenv("WRITE_PROJECT_JSON") != "" {
		log.Print("writing JSON back out to 'testing.json'...")
		f, err = os.Create("testing.json")
		if err != nil {
			return err
		}
		defer f.Close()
		if err := json.NewEncoder(f).Encode(proj); err != nil {
			return err
		}
	}

	// todo later - print out programs in pybricks

	return nil
}
