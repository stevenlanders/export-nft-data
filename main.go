package main

import (
	"errors"
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
						Name:  "alchemy-key",
						Value: os.Getenv("ALCHEMY_MAINNET_API_KEY"),
					},
					&cli.IntFlag{
						Name:  "days",
						Value: 30,
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
						AlchemyKey: c.String("alchemy-key"),
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
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
