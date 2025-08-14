# Core Functionalities

This section details Khedra's primary technical functionalities, explaining how each core feature is implemented and the technical approaches used.

## Control Service

### Service Management Interface

The Control Service exposes a minimal HTTP interface for pausing and unpausing supported services (`scraper`, `monitor`) and for reporting their pause status.

#### Technical Implementation

Implemented functions:

1. Pause a pausable service
2. Unpause a pausable service
3. Report paused / running status

Scope: does not provide start/stop/restart, runtime config mutation, or metrics collection.

```go
// Simplified Control Service interface
type ControlService struct {
    serviceManager *ServiceManager
    httpServer     *http.Server
    logger         *slog.Logger
}

type ServiceStatus struct {
    Name        string
    State       ServiceState
    LastStarted time.Time
    Uptime      time.Duration
    Metrics     map[string]interface{}
}

type ServiceState int
const (
    StateStopped ServiceState = iota
    StateStarting
    StateRunning
    StatePausing
    StatePaused
    StateStopping
)
```

#### Management Endpoints (Implemented)

1. `GET /isPaused` — status for all services
2. `GET /isPaused?name={service}` — status for one service
3. `GET /pause?name={service|all}` — pause service(s)
4. `GET /unpause?name={service|all}` — unpause service(s)

Mutating operations currently use GET.

#### Pausable Services

Only services implementing the `Pauser` interface can be paused:

- **Scraper**: Blockchain indexing service (pausable)
- **Monitor**: Address monitoring service (pausable)

Non‑pausable services: `control`, `api`, `ipfs` (if enabled). The monitor service is disabled by default but is pausable when enabled.

#### Service Coordination

Coordination is limited to toggling internal paused state.

## Blockchain Indexing

### The Unchained Index (High-Level Overview)

The Unchained Index implementation resides in upstream TrueBlocks libraries. This repository configures and invokes indexing.

#### Technical Implementation

The index is implemented as a specialized data structure with these key characteristics:

1. **Bloom Filter Front-End**: A probabilistic data structure that quickly determines if an address might appear in a block
2. **Address-to-Appearance Mapping**: Maps each address to a list of its appearances
3. **Chunked Storage**: Divides the index into manageable chunks (typically 1,000,000 blocks per chunk)
4. **Versioned Format**: Includes version metadata to handle format evolution

Internal storage specifics are handled upstream and not duplicated here.

#### Indexing Process (Conceptual)

High level only: batches of blocks are processed, appearances extracted, and persisted through the underlying TrueBlocks indexing subsystem; batch size and sleep are configured in `config.yaml`.

#### Performance Optimizations

- **Parallel Processing**: Multiple blocks processed concurrently
- **Bloom Filters**: Fast negative lookups to avoid unnecessary disk access
- **Binary Encoding**: Compact storage format for index data
- **Caching**: Frequently accessed index portions kept in memory

## Address Monitoring (Experimental / Limited)

### Monitor Implementation

The monitoring system currently provides service enablement/disablement and pause control; advanced notification features are outside this repository.

#### Technical Implementation

Monitors are implemented using these components:

1. **Monitor Registry**: Central store of all monitored addresses
2. **Address Index**: Fast lookup structure for monitored addresses
3. **Activity Tracker**: Records and timestamps address activity
4. **Notification Manager**: Handles alert distribution based on configuration

Implementation structs are managed upstream.

#### Monitoring Process

1. **Registration**: Add addresses to the monitor registry
2. **Block Processing**: As new blocks are processed, check for monitored addresses
3. **Activity Detection**: When a monitored address appears, record the activity
4. **Notification**: Based on configuration, send notifications via configured channels
5. **State Update**: Update the monitor's state with the new activity

#### Optimization Approaches

Optimizations will be added over time as needed.

## API Service (When Enabled)

### RESTful Interface

The API service provides HTTP endpoints for querying indexed data and managing Khedra's operations.

#### Technical Implementation

The API is implemented using these components:

1. **HTTP Server**: Handles incoming requests and routing
2. **Route Handlers**: Process specific endpoint requests
3. **Authentication Middleware**: Optional API key verification
4. **Response Formatter**: Structures data in requested format (JSON, CSV, etc.)
5. **Documentation**: Auto-generated Swagger documentation

Server implementation is provided by upstream services packages.

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

## IPFS Integration (Optional)

### Distributed Index Sharing

The IPFS integration enables sharing and retrieving index chunks through the distributed IPFS network.

#### Technical Implementation

The IPFS functionality is implemented with these components:

1. **IPFS Node**: Either embedded or external IPFS node connection
2. **Chunk Manager**: Handles breaking the index into shareable chunks
3. **Publishing Logic**: Manages uploading chunks to IPFS
4. **Discovery Service**: Finds and retrieves chunks from the network
5. **Validation**: Verifies the integrity of downloaded chunks

Implementation details are abstracted via the services layer.

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

## Configuration Management (YAML)

### Flexible Configuration System

Khedra's configuration system provides multiple ways to configure the application, with clear precedence rules.

#### Technical Implementation

The configuration system is implemented with these components:

1. **YAML Parser**: Reads the configuration file format
2. **Environment Variable Processor**: Overrides from environment variables
3. **Validation Engine**: Ensures configuration values are valid
4. **Defaults Manager**: Provides sensible defaults where needed
5. **Runtime Updater**: Handles configuration changes during operation

Authoritative structure lives in `pkg/types/config.go`.

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

The descriptions above match the repository's current functionality.
