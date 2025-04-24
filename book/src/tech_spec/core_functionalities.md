# Core Functionalities

This section details Khedra's primary technical functionalities, explaining how each core feature is implemented and the technical approaches used.

## Blockchain Indexing

### The Unchained Index

The Unchained Index is the foundational data structure of Khedra, providing a reverse-lookup capability from addresses to their appearances in blockchain data.

#### Technical Implementation

The index is implemented as a specialized data structure with these key characteristics:

1. **Bloom Filter Front-End**: A probabilistic data structure that quickly determines if an address might appear in a block
2. **Address-to-Appearance Mapping**: Maps each address to a list of its appearances
3. **Chunked Storage**: Divides the index into manageable chunks (typically 1,000,000 blocks per chunk)
4. **Versioned Format**: Includes version metadata to handle format evolution

```go
// Simplified representation of the index structure
type UnchainedIndex struct {
    Version string
    Chunks  map[uint64]*IndexChunk  // Key is chunk ID
}

type IndexChunk struct {
    BloomFilter   *BloomFilter
    Appearances   map[string][]Appearance  // Key is hex address
    StartBlock    uint64
    EndBlock      uint64
    LastUpdated   time.Time
}

type Appearance struct {
    BlockNumber    uint64
    TransactionIndex uint16
    AppearanceType  uint8
    LogIndex        uint16
}
```

#### Indexing Process

1. **Block Retrieval**: Fetch blocks from the RPC endpoint in configurable batches
2. **Appearance Extraction**: Process each block to extract address appearances from:
   - Transaction senders and recipients
   - Log topics and indexed parameters
   - Trace calls and results
   - State changes
3. **Deduplication**: Remove duplicate appearances within the same transaction
4. **Storage**: Update the appropriate index chunk with the new appearances
5. **Bloom Filter Update**: Update the bloom filter for quick future lookups

#### Performance Optimizations

- **Parallel Processing**: Multiple blocks processed concurrently
- **Bloom Filters**: Fast negative lookups to avoid unnecessary disk access
- **Binary Encoding**: Compact storage format for index data
- **Caching**: Frequently accessed index portions kept in memory

## Address Monitoring

### Monitor Implementation

The monitoring system tracks specific addresses for on-chain activity and provides notifications when activity is detected.

#### Technical Implementation

Monitors are implemented using these components:

1. **Monitor Registry**: Central store of all monitored addresses
2. **Address Index**: Fast lookup structure for monitored addresses
3. **Activity Tracker**: Records and timestamps address activity
4. **Notification Manager**: Handles alert distribution based on configuration

```go
// Simplified monitor implementation
type Monitor struct {
    Address       string
    Description   string
    CreatedAt     time.Time
    LastActivity  time.Time
    Config        MonitorConfig
    ActivityLog   []Activity
}

type MonitorConfig struct {
    NotificationChannels []string
    Filters              *ActivityFilter
    Thresholds           map[string]interface{}
}

type Activity struct {
    BlockNumber      uint64
    TransactionHash  string
    Timestamp        time.Time
    ActivityType     string
    Details          map[string]interface{}
}
```

#### Monitoring Process

1. **Registration**: Add addresses to the monitor registry
2. **Block Processing**: As new blocks are processed, check for monitored addresses
3. **Activity Detection**: When a monitored address appears, record the activity
4. **Notification**: Based on configuration, send notifications via configured channels
5. **State Update**: Update the monitor's state with the new activity

#### Optimization Approaches

- **Focused Index**: Maintain a separate index for just monitored addresses
- **Early Detection**: Check monitored addresses early in the processing pipeline
- **Configurable Sensitivity**: Allow users to set thresholds for notifications
- **Batched Notifications**: Group notifications to prevent excessive alerts

## API Service

### RESTful Interface

The API service provides HTTP endpoints for querying indexed data and managing Khedra's operations.

#### Technical Implementation

The API is implemented using these components:

1. **HTTP Server**: Handles incoming requests and routing
2. **Route Handlers**: Process specific endpoint requests
3. **Authentication Middleware**: Optional API key verification
4. **Response Formatter**: Structures data in requested format (JSON, CSV, etc.)
5. **Documentation**: Auto-generated Swagger documentation

```go
// Simplified API route implementation
type APIRoute struct {
    Path        string
    Method      string
    Handler     http.HandlerFunc
    Description string
    Params      []Parameter
    Responses   map[int]Response
}

// API server initialization
func NewAPIServer(config Config) *APIServer {
    server := &APIServer{
        router: mux.NewRouter(),
        port:   config.Port,
        auth:   config.Auth,
    }
    server.registerRoutes()
    return server
}
```

#### API Endpoints

The API provides endpoints in several categories:

1. **Status Endpoints**: System and service status information
2. **Index Endpoints**: Query the Unchained Index for address appearances
3. **Monitor Endpoints**: Manage and query address monitors
4. **Chain Endpoints**: Blockchain information and operations
5. **Admin Endpoints**: Configuration and management operations

#### Performance Considerations

- **Connection Pooling**: Reuse connections for efficiency
- **Response Caching**: Cache frequent queries with appropriate invalidation
- **Pagination**: Limit response sizes for large result sets
- **Query Optimization**: Efficient translation of API queries to index lookups
- **Rate Limiting**: Prevent resource exhaustion from excessive requests

## IPFS Integration

### Distributed Index Sharing

The IPFS integration enables sharing and retrieving index chunks through the distributed IPFS network.

#### Technical Implementation

The IPFS functionality is implemented with these components:

1. **IPFS Node**: Either embedded or external IPFS node connection
2. **Chunk Manager**: Handles breaking the index into shareable chunks
3. **Publishing Logic**: Manages uploading chunks to IPFS
4. **Discovery Service**: Finds and retrieves chunks from the network
5. **Validation**: Verifies the integrity of downloaded chunks

```go
// Simplified IPFS service implementation
type IPFSService struct {
    node        *ipfs.CoreAPI
    chunkManager *ChunkManager
    config      IPFSConfig
}

type ChunkManager struct {
    chunkSize      uint64
    validationFunc func([]byte) bool
    storage        *Storage
}
```

#### Distribution Process

1. **Chunking**: Divide the index into manageable chunks with metadata
2. **Publishing**: Add chunks to IPFS and record their content identifiers (CIDs)
3. **Announcement**: Share availability information through the network
4. **Discovery**: Find chunks needed by querying the IPFS network
5. **Retrieval**: Download needed chunks from peers
6. **Validation**: Verify chunk integrity before integration

#### Optimization Strategies

- **Incremental Updates**: Share only changed or new chunks
- **Prioritized Retrieval**: Download most useful chunks first
- **Peer Selection**: Connect to reliable peers for better performance
- **Background Syncing**: Retrieve chunks in the background without blocking
- **Compressed Storage**: Minimize bandwidth and storage requirements

## Configuration Management

### Flexible Configuration System

Khedra's configuration system provides multiple ways to configure the application, with clear precedence rules.

#### Technical Implementation

The configuration system is implemented with these components:

1. **YAML Parser**: Reads the configuration file format
2. **Environment Variable Processor**: Overrides from environment variables
3. **Validation Engine**: Ensures configuration values are valid
4. **Defaults Manager**: Provides sensible defaults where needed
5. **Runtime Updater**: Handles configuration changes during operation

```go
// Simplified configuration structure
type Config struct {
    General  GeneralConfig
    Chains   map[string]ChainConfig
    Services map[string]ServiceConfig
    Logging  LoggingConfig
}

// Configuration loading process
func LoadConfig(path string) (*Config, error) {
    config := DefaultConfig()
    
    // Load from file if exists
    if fileExists(path) {
        if err := loadFromFile(path, config); err != nil {
            return nil, err
        }
    }
    
    // Override with environment variables
    applyEnvironmentOverrides(config)
    
    // Validate the final configuration
    if err := validateConfig(config); err != nil {
        return nil, err
    }
    
    return config, nil
}
```

#### Configuration Sources

The system processes configuration from these sources, in order of precedence:

1. **Environment Variables**: Highest precedence, override all other sources
2. **Configuration File**: User-provided settings in YAML format
3. **Default Values**: Built-in defaults for unspecified settings

#### Validation Rules

The configuration system enforces these kinds of validation:

1. **Type Validation**: Ensures values have the correct data type
2. **Range Validation**: Numeric values within acceptable ranges
3. **Format Validation**: Strings matching required patterns (e.g., URLs)
4. **Dependency Validation**: Related settings are consistent
5. **Resource Validation**: Settings are compatible with available resources

These core functionalities form the technical foundation of Khedra, enabling its primary capabilities while providing the flexibility and performance required for blockchain data processing.
