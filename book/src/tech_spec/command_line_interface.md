# Command-Line Interface

## Available Commands and Options

### Initialization

```bash
./khedra --init all
```

- Options: `all`, `blooms`, `none`

### Scraper

```bash
./khedra --scrape on
```

- Enables or disables the blockchain scraper.

### REST API

```bash
./khedra --api on
```

- Starts the API server.

### Sleep Duration

```bash
./khedra --sleep 60
```

- Sets the duration (in seconds) between updates.

## Detailed Behavior for Each Command

1. **`--init`**: Controls how the blockchain index is initialized.
2. **`--scrape`**: Toggles the blockchain scraper.
3. **`--api`**: Starts or stops the API server.
