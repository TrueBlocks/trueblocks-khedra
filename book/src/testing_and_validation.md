# Testing and Validation

## Unit Testing

Unit tests cover:

- Blockchain indexing logic.
- Configuration parsing and validation.
- REST API endpoint functionality.

Run tests with:

```bash
go test ./...
```

## Integration Testing

Integration tests ensure all components work together as expected. Tests include:

- RPC connectivity validation.
- Multi-chain indexing workflows.

## Testing Guidelines for Developers

1. Use mock RPC endpoints for testing without consuming live resources.
2. Validate `.env` configuration in test environments before deployment.
3. Automate tests with CI/CD pipelines to ensure reliability.
