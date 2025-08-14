# Tall Affiliate Common Library

This library contains shared code for the Tall Affiliate microservices.

## Packages

- `pkg/models`: Shared data models
- `pkg/events`: Event definitions
- `pkg/interfaces`: Common interfaces
- `pkg/constants`: Shared constants

## Usage

Add this library as a dependency to your service:

```bash
go get github.com/maltedev/tall-affiliate/tall-affiliate-common
```

Import the packages you need:

```go
import (
    "github.com/maltedev/tall-affiliate/tall-affiliate-common/pkg/models"
    "github.com/maltedev/tall-affiliate/tall-affiliate-common/pkg/events"
    "github.com/maltedev/tall-affiliate/tall-affiliate-common/pkg/interfaces"
    "github.com/maltedev/tall-affiliate/tall-affiliate-common/pkg/constants"
)
```
