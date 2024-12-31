# Maintenance and Troubleshooting

## Updating Khedra

To update the application, pull the latest changes from the repository and rebuild the binary:

```bash
git pull
go build -o khedra .
```

## Common Issues and Solutions

- **Missing RPC Provider**: Ensure your `.env` file contains valid RPC URLs.
- **Configuration Errors**: Use `--help` to validate command-line arguments.

## Log Files and Debugging

Logs are written to the standard output by default. Set the log level in the `.env` file:

```env
TB_LOGLEVEL="Debug"
```

## Contacting Support

If you encounter issues not covered in this guide, contact support at:
[TrueBlocks Support](mailto:support@trueblocks.io)
