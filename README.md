# TrueBlocks Khedra

## Intro

`trueblocks-khedra` is an extension (or plugin) to the [TrueBlocks](https://github.com/TrueBlocks/trueblocks-core) system that focuses on providing specialized data extraction, analysis, or other functionality related to Ethereum blockchain indexing. Khedra aims to simplify the process of gathering on-chain data and building advanced, queryable indexes for Ethereum addresses.

Key features:

- **Custom Indexing**: Provides specialized indexing capabilities tailored to specific use-cases beyond the core TrueBlocks functionality.  
- **Plugin-Based Architecture**: Easily integrates with TrueBlocks while maintaining modular, extensible design.  
- **Efficient Data Retrieval**: Optimized for quick querying and data lookups, especially when dealing with large Ethereum datasets.

## Installation

### Prerequisites

- Make sure you have [TrueBlocks Core](https://github.com/TrueBlocks/trueblocks-core) installed.  
- A C++ build environment (such as `g++` or `clang++`) if you plan to compile from source.  
- [CMake](https://cmake.org/) (version 3.16 or higher recommended).  
- (Optional) [Docker](https://docs.docker.com/get-docker/) if you plan to run via container.

### Clone this Repository

  ```[bash]
      git clone https://github.com/TrueBlocks/trueblocks-khedra.git  
      cd trueblocks-khedra  
  ```

### Build from Source

  ```[bash]
      mkdir build && cd build  
      cmake ..  
      make  
  ```

   After a successful build, you’ll find the `khedra` executable (or library, depending on how the project is organized) in the build output.

### Install

```[bash]
sudo make install  
```

## Configuration

Before using `khedra`, you may need to configure it to point at the TrueBlocks indexing data or specify custom indexing rules:

- **Config File**: By default, `khedra` may look for a configuration file at `~/.trueblocks/trueblocks-khedra.conf`.  
- **Environment Variables**:  
  - `KHEDRA_DATA_DIR`: Path to where you want `khedra` to store or read data.  
  - `KHEDRA_LOG_LEVEL`: Adjusts the verbosity of logs (`DEBUG`, `INFO`, `WARN`, `ERROR`).

Refer to the sample configuration file (`.conf.example`) in this repo for a template of possible settings.

---

## Docker Version - Building & Running

Build the Docker image:

```bash
docker build -t trueblocks-khedra .
```

Run the Docker container (showing the help message by default):

```bash
docker run --rm -it trueblocks-khedra
```

Use a custom command, for example to specify a subcommand or different flags:

```bash
docker run --rm -it trueblocks-khedra some-subcommand --flag
```

Adjust paths, environment variables, or your config file strategy as needed. You can also mount external volumes (e.g., a local ~/.trueblocks directory) if you prefer to maintain data outside the container.

---

## Documentation

<!--
  BEGIN SECTION: (Exact text from trueblocks-core README)
  Copy/Paste the "Documentation" section here verbatim.
-->

**(Paste the *exact* Documentation text from the trueblocks-core README here.)**

---

## Linting

<!--
  BEGIN SECTION: (Exact text from trueblocks-core README)
  Copy/Paste the "Linting" section here verbatim.
-->

**(Paste the *exact* Linting text from the trueblocks-core README here.)**

---

## Contributing

<!--
  BEGIN SECTION: (Exact text from trueblocks-core README)
  Copy/Paste the "Contributing" section here verbatim.
-->

**(Paste the *exact* Contributing text from the trueblocks-core README here.)**

---

## Contact

<!--
  BEGIN SECTION: (Exact text from trueblocks-core README)
  Copy/Paste the "Contact" section here verbatim.
-->

**(Paste the *exact* Contact text from the trueblocks-core README here.)**

---

## Contributors

<!--
  BEGIN SECTION: (Exact text from trueblocks-core README)
  Copy/Paste the "Contributors" section here verbatim.
-->

**(Paste the *exact* Contributors text from the trueblocks-core README here.)**

---

This project is part of the [TrueBlocks](https://github.com/TrueBlocks) ecosystem.
