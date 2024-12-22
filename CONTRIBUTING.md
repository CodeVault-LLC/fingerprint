# Contribution Guidelines

## General Guidelines

- **Use our Logger for Logging**
  Always use our logger for logging. This ensures that logs are consistent and can be easily filtered.

  ```go
  import "github.com/codevault-llc/fingerprint/pkg/logger"

  logger.Log.Info("This is an info log")
  logger.Log.Error("This is an error log")
  ```
