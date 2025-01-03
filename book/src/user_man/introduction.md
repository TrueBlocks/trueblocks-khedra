# Introduction

Blockchains are long running-processes that continually create new data (in the form of blocks). For this reason, any process that wishes to monitor, index, or access data from a blockchain must also be long running.

**Khedra** is such a long-running process.

In order to remain decentralized and permissionless, blockchains must be "freed" from the stranglehold of large data providers. One way to do that is to help people run blockchain nodes locally. However, as soon as one does that, one learns that blockchains are not very good databases. This is for a simple reason, they lack an index.

TrueBlocks Core (of which **chifra** and **khedra** are a part) is a set of command-line tools, SDKs, and packages that help users who are running their own blockchain nodes make better use of the data. **Khedra** indexes and monitors the data. **Chifra** helps access the data providing various useful commands for exporting, filtering, and processing on-chain activity.

Of primary importance in the design of both systems are:

- **speed** - we cache nearly everything
- **permisionless access** - no servers, no API keys, you run your own infrastructure
- **accuracy** - the goal is 100% off-chain reconciliation of account balances and state history
- **depth of detail** - required to enable 100% accurate reconciliations
- **ease of use** - so shoot us - this one is hard

*Enjoy!*

Please help us improve this software by providing any feedback or suggestions. Contact information and links to our socials are available [at our website](https://trueblocks.io).
