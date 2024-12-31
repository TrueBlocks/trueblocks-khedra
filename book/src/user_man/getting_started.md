# Getting Started

This guide will help you set up and run Khedra on your system.

## Requirements

Your system must meet the following requirements

| Requirement | Description                                |
| ----------- | ------------------------------------------ |
| OS          | Linux or macOS                             |
| Golang      | Version 1.23 or later                      |
| Internet    | Active internet connection                 |
| RPC Node    | (Preferably local) RPC node for each chain |

## Installation Guide

1. Clone the repository

   ```bash
   git clone https://github.com/TrueBlocks/trueblocks-khedra.git
   cd trueblocks-khedra
   ```

2. Build **khedra**

   ```bash
   go build -o khedra .
   ```

3. Configure the environment

   ```bash
   cp -p env.example .env # Edit the file for your system
   ```

See below for more information on configuring your `.env` environment file.

## Initial Configuration (Build from Source)

Populate your `.env` file with the necessary parameters:

```env
# required
TB_NODE_DATADIR="<path/to/data/directory>"     # Path to the data directory
TB_NODE_CHAINS="mainnet,sepolia"               # Comma-separated list of chains
TB_NODE_MAINNETRPC="<your-mainnet-rpc>"        # Mainnet RPC endpoint (required)

# optional
TB_NODE_<CHAIN>>RPC="<your-CHAIN-rpc>"         # RPC for other chains in TB_NODE_CHAINS
TB_LOGLEVEL="Info"                             # One of Debug, Info, Warn, Error
TB_NODE_BLOCKCNT=2000                          # Number of blocks to scrape in each pass (default: 2000)
```

## Initial Configuration (Go Install)

If you wish (we don't recommend it), you may avoid cloning the repo and install the application directly with:

```bash
go install github.com/TrueBlocks/trueblocks-khedra
```

## Starting the Application

Run the application with:

```bash
./khedra --version
```

## More Information

See the Technical Specification for more information on the [application's command line options](../tech_spec/command_line_interface.md).
