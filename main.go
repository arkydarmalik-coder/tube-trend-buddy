// Package main is the entrypoint for tube-trend-buddy, an 8-in-1 YouTube
// creator toolkit. See cmd/ for subcommands and internal/ for shared logic.
package main

import "github.com/arkydarmalik-coder/tube-trend-buddy/cmd"

func main() {
	cmd.Execute()
}
