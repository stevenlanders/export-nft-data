package main

import (
	"errors"
	"export-nft-data/cmd/discover"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	"export-nft-data/cmd/edges"
	"export-nft-data/cmd/owners"
	"export-nft-data/cmd/supernodes"
)

func main() {
	app := &cli.App{
		Name:  "export-nft-data",
		Usage: "a cli for processing the nft poc",
		Commands: []*cli.Command{
			{
				Name:  "supernodes",
				Usage: "(step 1) find supernodes from names in a file",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "centerdev-key",
						Value: os.Getenv("CENTERDEV_KEY"),
					}, &cli.StringFlag{
						Name:  "mainnet-json-rpc",
						Value: os.Getenv("MAINNET_JSON_RPC"),
					}, &cli.StringFlag{
						Name:  "etherscan-key",
						Value: os.Getenv("ETHERSCAN_KEY"),
					}},
				ArgsUsage: "<collections-csv>",
				Action: func(c *cli.Context) error {
					file := c.Args().Get(0)
					if file == "" {
						return errors.New("<collections-csv> is required")
					}
					return supernodes.Run(c.Context, supernodes.Config{
						File:         file,
						CenterDevKey: c.String("centerdev-key"),
						EtherscanKey: c.String("etherscan-key"),
						JsonRpcUrl:   c.String("mainnet-json-rpc"),
					})
				},
			},
			{
				Name:  "owners",
				Usage: "(step 2) find owners from names in a collection json",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "mainnet-json-rpc",
						Value: os.Getenv("MAINNET_JSON_RPC"),
					},
					&cli.IntFlag{
						Name:  "days",
						Value: 180,
					},
				},
				ArgsUsage: "<collections-json>",
				Action: func(c *cli.Context) error {
					file := c.Args().Get(0)
					if file == "" {
						return errors.New("<collections-json> is required")
					}
					return owners.Run(c.Context, owners.Config{
						File:       file,
						JsonRpcUrl: c.String("mainnet-json-rpc"),
						Days:       c.Int("days"),
					})
				},
			},
			{
				Name:  "edges",
				Usage: "(step 3) find new edges given existing owners & collections",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "mainnet-json-rpc",
						Value: os.Getenv("MAINNET_JSON_RPC"),
					},
					&cli.Uint64Flag{
						Name:     "start-block",
						Required: true,
					},
					&cli.IntFlag{
						Name:  "num-blocks",
						Value: 30,
					},
				},
				ArgsUsage: "<buyer-and-collections-json>",
				Action: func(c *cli.Context) error {
					file := c.Args().Get(0)
					if file == "" {
						return errors.New("<buyer-and-collections-json> is required")
					}
					return edges.Run(c.Context, edges.Config{
						File:           file,
						JsonRpcUrl:     c.String("mainnet-json-rpc"),
						BlockStart:     c.Uint64("start-block"),
						NumberOfBlocks: c.Int("num-blocks"),
					})
				},
			},
			{
				Name:  "discover",
				Usage: "start with a small set of nodes and discover a graph",
				Flags: []cli.Flag{
					&cli.Uint64Flag{
						Name:  "start-block",
						Value: 14946474,
					},
					&cli.IntFlag{
						Name:  "num-blocks",
						Value: 1000000,
					},
					&cli.IntFlag{
						Name:  "owner-days",
						Value: 30,
					},
					&cli.IntFlag{
						Name:  "iterations",
						Value: 1,
					},
					&cli.StringFlag{
						Name:  "mainnet-json-rpc",
						Value: os.Getenv("MAINNET_JSON_RPC"),
					},
					&cli.StringFlag{
						Name:  "centerdev-key",
						Value: os.Getenv("CENTERDEV_KEY"),
					},
					&cli.StringFlag{
						Name:  "etherscan-key",
						Value: os.Getenv("ETHERSCAN_KEY"),
					},
				},
				ArgsUsage: "<input-json>",
				Action: func(c *cli.Context) error {
					file := c.Args().Get(0)
					if file == "" {
						return errors.New("<input-json> is required")
					}
					return discover.Run(c.Context, discover.Config{
						File:           file,
						JsonRpcUrl:     c.String("mainnet-json-rpc"),
						CenterDevKey:   c.String("centerdev-key"),
						EtherscanKey:   c.String("etherscan-key"),
						BlockStart:     c.Uint64("start-block"),
						NumberOfBlocks: c.Uint64("num-blocks"),
						OwnerDays:      c.Int("owner-days"),
						Iterations:     c.Int("iterations"),
					})
				},
			},
		},
	}

	log.SetLevel(log.DebugLevel)
	err := app.Run(os.Args)
	if err != nil {
		log.WithError(err).Fatal(err)
	}
}
