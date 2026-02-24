package engine

import (
	"log"

	_ "github.com/ochain-gg/ochain-network/engine/transactions/account"
	_ "github.com/ochain-gg/ochain-network/engine/transactions/game/alliance"
	_ "github.com/ochain-gg/ochain-network/engine/transactions/game/fleet"
	_ "github.com/ochain-gg/ochain-network/engine/transactions/game/planet"
	_ "github.com/ochain-gg/ochain-network/engine/transactions/governance"
	_ "github.com/ochain-gg/ochain-network/engine/transactions/validator"
)

// init is called after all transaction package init functions have completed
// This serves as a confirmation that the initialization chain has executed
func init() {
	log.Println("Engine package loaded - all transaction handlers should be registered")
}
