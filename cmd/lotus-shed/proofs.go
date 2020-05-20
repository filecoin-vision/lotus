package main

import (
	"encoding/hex"
	"fmt"

	"gopkg.in/urfave/cli.v2"

	ffi "github.com/filecoin-project/filecoin-ffi"
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/specs-actors/actors/abi"
	"github.com/ipfs/go-cid"
)

var proofsCmd = &cli.Command{
	Name: "proofs",
	Subcommands: []*cli.Command{
		verifySealProofCmd,
	},
}

var verifySealProofCmd = &cli.Command{
	Name:        "verify-seal",
	Description: "Verify a seal proof with manual inputs",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name: "ticket",
		},
		&cli.StringFlag{
			Name: "proof-rand",
		},
		&cli.StringFlag{
			Name: "miner",
		},
		&cli.Uint64Flag{
			Name: "sector-id",
		},
		&cli.Int64Flag{
			Name: "proof-type",
		},
	},
	Action: func(cctx *cli.Context) error {
		if cctx.Args().Len() != 3 {
			return fmt.Errorf("must specify commR, commD, and proof to verify")
		}

		commr, err := cid.Decode(cctx.Args().Get(0))
		if err != nil {
			return err
		}

		commd, err := cid.Decode(cctx.Args().Get(1))
		if err != nil {
			return err
		}

		proof, err := hex.DecodeString(cctx.Args().Get(2))
		if err != nil {
			return fmt.Errorf("failed to decode hex proof input: %w", err)
		}

		maddr, err := address.NewFromString(cctx.String("miner"))
		if err != nil {
			return err
		}

		mid, err := address.IDFromAddress(maddr)
		if err != nil {
			return err
		}

		snum := abi.SectorNumber(cctx.Uint64("sector-id"))

		ok, err := ffi.VerifySeal(abi.SealVerifyInfo{
			SectorID: abi.SectorID{
				Miner:  abi.ActorID(mid),
				Number: snum,
			},
			OnChain: abi.OnChainSealVerifyInfo{
				SealedCID:        commr,
				InteractiveEpoch: 0,
				RegisteredProof:  abi.RegisteredProof(cctx.Int64("proof-type")),
				Proof:            proof,
				DealIDs:          nil,
				SectorNumber:     snum,
				SealRandEpoch:    0,
			},
			Randomness:            abi.SealRandomness(nil),
			InteractiveRandomness: abi.InteractiveSealRandomness(nil),
			UnsealedCID:           commd,
		})
		if err != nil {
			return err
		}
		if !ok {
			return fmt.Errorf("invalid proof")
		}

		fmt.Println("proof valid!")
		return nil
	},
}
