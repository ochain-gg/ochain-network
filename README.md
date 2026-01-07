# OChain Network Validator

A blockchain validator node for the OChain Network, a space strategy game built on CometBFT with Ethereum bridge integration.

## Overview

OChain Network is a blockchain-based space strategy game that combines DeFi mechanics with gaming. Players manage planets, build fleets, form alliances, and engage in space combat while earning and spending OCT tokens. The validator node handles all game logic, economic transactions, and cross-chain operations.

## Features

### Core Blockchain
- **Consensus Engine**: CometBFT v1.0.1 (formerly Tendermint)
- **Database**: BadgerDB for persistent storage
- **Bridge**: Ethereum integration via smart contracts (Sepolia testnet)
- **Scheduler**: Automated background tasks and epoch management

### Game Mechanics

#### Planet Management
- Building construction and upgrades
- Technology research and development
- Resource production and management
- Resource trading marketplace

#### Fleet Operations
- Fleet creation, merging, and splitting
- Cargo management and space missions
- Fleet combat system
- Recycling and lending mechanics

#### Alliance System
- Alliance creation and management
- Membership request system
- Role assignment within alliances
- Collaborative gameplay features

### Economic System
- **OCT Token**: Native token with 8 decimals
- **Cross-chain**: Ethereum ↔ OChain bridge
- **Limits**: 10,000 OCT weekly deposit limit
- **Market**: Integrated resource exchange
- **Governance**: Proposal and voting system with staking

## Architecture

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│  Ethereum       │    │  OChain Network  │    │  Game Client    │
│  Smart Contracts│◄──►│  Validator       │◄──►│  Interface      │
│  (Bridge)       │    │  (CometBFT)      │    │                 │
└─────────────────┘    └──────────────────┘    └─────────────────┘
```

### Transaction Types (95 types)
- **System Transactions**: Validator management, epochs, bridge operations
- **Account Transactions**: Registration, authentication, faucet claims
- **Governance Transactions**: Proposals, voting, staking mechanisms
- **Game Transactions**: All space game actions (planets, fleets, combat)

## Installation

### Prerequisites
- Go 1.23.5+
- Docker & Docker Compose (optional)

### Build from Source
```bash
git clone https://github.com/ochain-gg/ochain-network
cd ochain-network-validator
go build -o ochain-validator main.go
```

### Docker Setup
```bash
docker-compose up -d
```

## Configuration

### Environment Variables
- `CMT_HOME`: CometBFT configuration directory (default: `$HOME/.cometbft`)
- `CHAIN_ID`: Ethereum chain ID (default: `11155111` for Sepolia)
- `EVM_RPC`: Ethereum RPC endpoint
- `PORTAL_ADDRESS`: OChain Portal contract address

### Example Start
```bash
./ochain-validator \
  --cmt-home ~/.cometbft \
  --chainId 11155111 \
  --evmRpc https://ethereum-sepolia.core.chainstack.com/... \
  --portalAddress 0x4Dd9d772C67fbC858918f364E5CB9e0B6E53Fd44
```

## API Endpoints

### Queries
- `GET /accounts` - Account information
- `GET /planet/{id}` - Planet details
- `GET /market` - Resource market data
- `GET /universe/{id}` - Universe state

### Game Actions
Submit transactions via CometBFT RPC or use the game client interface.

## Smart Contracts

### Deployed Contracts (Sepolia)
- **OChain Portal**: `0x4Dd9d772C67fbC858918f364E5CB9e0B6E53Fd44`
- **OChain Token**: ERC-20 token contract
- **Bridge Logic**: Cross-chain deposit/withdrawal handling

## Development

### Project Structure
```
├── application/          # ABCI application logic
├── cmd/                 # CLI commands
├── config/              # Configuration management
├── contracts/           # Ethereum contract bindings
├── engine/              # Transaction processing engine
│   ├── database/        # Data models and storage
│   └── transactions/    # Transaction handlers
├── queries/             # Query handlers
├── scheduler/           # Background job scheduler
├── types/               # Type definitions
└── tests/               # Test utilities
```

### Running Tests
```bash
go test ./...
```

### Adding New Transaction Types
1. Define transaction type in `engine/transactions/transaction.go`
2. Create handler in appropriate subdirectory
3. Register in transaction engine
4. Add validation and finalization logic

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## Security

- All transactions use EIP-712 typed data signing
- Cross-chain operations are secured by smart contract validation
- Validator consensus ensures network security and integrity

## License

This project is licensed under the terms specified in the LICENSE file.

## Links

- **Game Website**: https://ochain.gg
- **Documentation**: [Coming Soon]
- **Discord Community**: [Link TBA]
- **GitHub Issues**: Report bugs and feature requests

---

**Note**: This is alpha software. Use at your own risk on testnets only.