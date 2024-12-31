# Using Khedra

## Indexing Blockchains

To index a blockchain, ensure the required environment variables are set for your RPC endpoints, then run:

```bash
./khedra --init all --scrape on
```

This will initialize the blockchain index and start the scraping process.

## Accessing the REST API

Enable the REST API by running the application with:

```bash
./khedra --api on
```

Access the API through the default endpoint:

```bash
curl http://localhost:8080
```

Refer to the API documentation for available endpoints and usage.

## Monitoring Addresses

You can monitor specific blockchain addresses for transactions. Configure the monitored addresses in your `.env` file or through the API, and enable monitoring:

```bash
./trueblocks-node --monitor on
```

## Managing Configurations

Khedra configurations can be managed using the `.env` file. Changes to the `.env` file require a restart of the application to take effect.
