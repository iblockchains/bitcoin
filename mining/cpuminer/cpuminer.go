package cpuminer

import (
	"fmt"

	"github.com/iblockchains/bitcoin/mining"
)

type Config struct {
	// BlockTemplateGenerator 区块生成模板
	// identifies the instance to use in order to
	// generate block templates that the miner will attempt to solve.
	BlockTemplateGenerator *mining.BlkTmplGenerator
}
type CPUMiner struct{}

func New(cfg *Config) *CPUMiner {
	fmt.Println("Unfinished:cpuminer.New")
	return nil
}
