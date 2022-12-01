// Copyright (c) 2022 Husarnet sp. z o.o.
// Authors: listed in project_root/README.md
// License: specified in project_root/LICENSE.txt
package main

import (
	"net/url"
	"fmt"
	"strconv"
	"github.com/pterm/pterm"
	"github.com/urfave/cli/v2"
)

var daemonStartCommand = &cli.Command{
	Name:      "start",
	Aliases:   []string{"up"},
	Usage:     "start husarnet daemon",
	ArgsUsage: " ", // No arguments needed
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:        "wait",
			Aliases:     []string{"w"},
			Usage:       "wait until the daemon is running",
			Destination: &wait,
		},
	},
	Action: func(ctx *cli.Context) error {
		if onWindows() {
			runSubcommand(false, "nssm", "start", "husarnet")
		} else {
			runSubcommand(false, "sudo", "systemctl", "start", "husarnet")
		}
		printSuccess("Started husarnet-daemon")

		if wait {
			waitDaemon()
		}

		return nil
	},
}

var daemonRestartCommand = &cli.Command{
	Name:      "restart",
	Aliases:   []string{"down"},
	Usage:     "restart husarnet daemon",
	ArgsUsage: " ", // No arguments needed
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:        "wait",
			Aliases:     []string{"w"},
			Usage:       "wait until the daemon is running",
			Destination: &wait,
		},
	},
	Action: func(ctx *cli.Context) error {
		if onWindows() {
			runSubcommand(false, "nssm", "restart", "husarnet")
		} else {
			runSubcommand(false, "sudo", "systemctl", "restart", "husarnet")
		}
		printSuccess("Restarted husarnet-daemon")

		if wait {
			waitDaemon()
		}

		return nil
	},
}

var daemonStopCommand = &cli.Command{
	Name:      "stop",
	Aliases:   []string{"down"},
	Usage:     "stop husarnet daemon",
	ArgsUsage: " ", // No arguments needed
	Action: func(ctx *cli.Context) error {
		if onWindows() {
			runSubcommand(false, "nssm", "stop", "husarnet")
		} else {
			runSubcommand(false, "sudo", "systemctl", "stop", "husarnet")
		}
		printSuccess("Stopped husarnet-daemon")
		return nil
	},
}

var daemonStatusCommand = &cli.Command{
	Name:  "status",
	Usage: "Display current connectivity status",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:    "verbose",
			Aliases: []string{"v"},
			Usage:   "show more information",
		},
	},
	ArgsUsage: " ", // No arguments needed
	Action: func(ctx *cli.Context) error {
		printStatus(ctx)
		return nil

	},
}

var daemonPingCommand = &cli.Command{
	Name:  "ping",
	Usage: "Ping base server",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:        "count",
			Aliases:     []string{"c"},
			Usage:       "count fo pings to send",
			Destination: &pingCount,
		},
	},
	ArgsUsage: "[count]", // No arguments needed
	Action: func(ctx *cli.Context) error {
		count:=10
		if pingCount {
			tmp, err := strconv.Atoi(ctx.Args().Get(0))
			if err != nil {
				printWarning(fmt.Sprintf("Wrong type of argument was given as count: %s",err.Error()))
				return nil
			}
			count = tmp
		}
		for i:=0; i<count; i++ {
			result:= callDaemonGet[PingResult]("/api/ping").Result.Message
			printInfo(result)
		}
		return nil

	},
}

var daemonSetupServerCommand = &cli.Command{
	Name:      "setup-server",
	Aliases:   []string{"change-dashboard"},
	Usage:     "Connect your Husarnet device to different Husarnet infrastructure",
	ArgsUsage: "[dashboard fqdn]",
	Action: func(ctx *cli.Context) error {
		requiredArgumentsNumber(ctx, 1)

		domain := ctx.Args().Get(0)

		callDaemonPost[EmptyResult]("/api/change-server", url.Values{
			"domain": {domain},
		})

		printSuccess("Successfully requested a change to %s server", pterm.Bold.Sprint(domain))
		printWarning("This action requires you to restart the daemon in order to use the new value")

		if onWindows() {
			runSubcommand(true, "nssm", "restart", "husarnet")
		} else {
			runSubcommand(true, "sudo", "systemctl", "restart", "husarnet")
		}

		waitDaemon()

		return nil
	},
}

var daemonWhitelistCommand = &cli.Command{
	Name:  "whitelist",
	Usage: "Manage whitelist on the device.",
	Subcommands: []*cli.Command{
		{
			Name:      "enable",
			Aliases:   []string{"on"},
			Usage:     "enable whitelist",
			ArgsUsage: " ", // No arguments needed
			Action: func(ctx *cli.Context) error {
				callDaemonPost[EmptyResult]("/api/whitelist/enable", url.Values{})
				printSuccess("Enabled the whitelist")

				return nil
			},
		},
		{
			Name:      "disable",
			Aliases:   []string{"off"},
			Usage:     "disable whitelist",
			ArgsUsage: " ", // No arguments needed
			Action: func(ctx *cli.Context) error {
				callDaemonPost[EmptyResult]("/api/whitelist/whitelist", url.Values{})
				printSuccess("Disabled the whitelist")

				return nil
			},
		},
		{
			Name:      "ls",
			Aliases:   []string{"show", "dir", "get"},
			Usage:     "list entries on the whitelist",
			ArgsUsage: " ", // No arguments needed
			Action: func(ctx *cli.Context) error {
				status := getDaemonStatus()
				printWhitelist(status, false)

				return nil
			},
		},
		{
			Name:      "add",
			Usage:     "Add a device to your whitelist by Husarnet address",
			ArgsUsage: "[device's ip address]",
			Action: func(ctx *cli.Context) error {
				requiredArgumentsNumber(ctx, 1)

				addr := makeCannonicalAddr(ctx.Args().Get(0))

				callDaemonPost[EmptyResult]("/api/whitelist/add", url.Values{
					"address": {addr},
				})
				printSuccess("Added %s to whitelist", addr)

				return nil
			},
		},
		{
			Name:      "rm",
			Usage:     "Remove device from the whitelist",
			ArgsUsage: "[device's ip address]",
			Action: func(ctx *cli.Context) error {
				requiredArgumentsNumber(ctx, 1)

				addr := makeCannonicalAddr(ctx.Args().Get(0))

				callDaemonPost[EmptyResult]("/api/whitelist/rm", url.Values{
					"address": {addr},
				})
				printSuccess("Removed %s from whitelist", addr)

				return nil
			},
		},
	},
}

var daemonWaitCommand = &cli.Command{
	Name:  "wait",
	Usage: "Wait until certain events occur. If no events provided will wait for as many elements as it can (the best case scenario). Husarnet will continue working even if some of those elements are unreachable, so consider narrowing your search down a bit.",
	Action: func(ctx *cli.Context) error {
		ignoreExtraArguments(ctx)
		err := waitBaseANY()
		if err != nil {
			return err
		}

		err = waitBaseUDP()
		if err != nil {
			return err
		}

		return waitWebsetup()
	},
	Subcommands: []*cli.Command{
		{
			Name:      "daemon",
			Aliases:   []string{"ready"},
			Usage:     "Wait until the deamon is able to return it's own status",
			ArgsUsage: " ", // No arguments needed
			Action: func(ctx *cli.Context) error {
				ignoreExtraArguments(ctx)
				return waitDaemon()
			},
		},
		{
			Name:      "base",
			Usage:     "Wait until there is a base-server connection established (via any protocol)",
			ArgsUsage: " ", // No arguments needed
			Action: func(ctx *cli.Context) error {
				ignoreExtraArguments(ctx)
				return waitBaseANY()
			},
			Subcommands: []*cli.Command{
				{
					Name:      "udp",
					Usage:     "Wait until there is a base-server connection established via UDP. This is the best case scenario. Husarnet will work even without it.",
					ArgsUsage: " ", // No arguments needed
					Action: func(ctx *cli.Context) error {
						ignoreExtraArguments(ctx)
						return waitBaseUDP()
					},
				},
			},
		},
		{
			Name:      "joinable",
			Usage:     "Wait until there is enough connectivity to join to a network",
			ArgsUsage: " ", // No arguments needed
			Action: func(ctx *cli.Context) error {
				ignoreExtraArguments(ctx)

				err := waitBaseANY()
				if err != nil {
					return err
				}

				return waitWebsetup()
			},
		},
		{
			Name:      "joined",
			Usage:     "Wait until there the daemon has joined the network",
			ArgsUsage: " ", // No arguments needed
			Action: func(ctx *cli.Context) error {
				ignoreExtraArguments(ctx)
				return waitJoined()
			},
		},
		{
			Name:      "host",
			Usage:     "Wait until there is an established connection to a given host",
			ArgsUsage: "[device name or ip]",
			Action: func(ctx *cli.Context) error {
				requiredArgumentsNumber(ctx, 1)
				return waitHost(ctx.Args().First())
			},
		},
		{
			Name:      "hostnames",
			Usage:     "Wait until given hosts are known to the system",
			ArgsUsage: "[device name] […]",
			Action: func(ctx *cli.Context) error {
				minimumArguments(ctx, 1)
				return waitHostnames(ctx.Args().Slice())
			},
		},
	},
}

var daemonCommand = &cli.Command{
	Name:  "daemon",
	Usage: "Control the local daemon",
	Subcommands: []*cli.Command{
		daemonStatusCommand,
		joinCommand,
		daemonPingCommand,
		daemonSetupServerCommand,
		daemonStartCommand,
		daemonRestartCommand,
		daemonStopCommand,
		daemonWhitelistCommand,
		daemonWaitCommand,
	},
}
