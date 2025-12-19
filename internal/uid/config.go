package uid

import (
	"os"
	"sync"

	"github.com/doffy007/tlab-test/internal/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Snowflake struct {
	mu        sync.Mutex
	timestamp int64
	nodeID    int64
	sequence  int64
}

const (
	nodeBits     = 10
	sequenceBits = 12

	sequenceMax = -1 ^ (-1 << sequenceBits)
	nodeMax     = -1 ^ (-1 << nodeBits)

	timeShift = nodeBits + sequenceBits
	nodeShift = sequenceBits
	epoch     = int64(1700000000000)
)

var (
	sf   *Snowflake
	once sync.Once
)

func init() {
	once.Do(func() {
		config.ConfigApps()

		zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

		nodeID := config.AppConfig.Snowflake.NodeID
		if nodeID < 0 {
			idFromIP, err := getNodeIDFromIP(config.AppConfig.Snowflake.IP)
			if err != nil {
				log.Fatal().Err(err).Msg("failed to derive node ID from IP")
			}
			nodeID = idFromIP
		}

		if nodeID < 0 || nodeID > nodeMax {
			log.Fatal().Int64("node_id", nodeID).Msg("invalid snowflake node id")
		}

		sf = &Snowflake{nodeID: nodeID}

		log.Info().
			Int64("node_id", nodeID).
			Msg("initialize snowflake uid generator")
	})
}
