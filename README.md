# GO-SCRIPTING

This project provides a set of scripts designed for bulk updates using CSV files, leveraging Go routines. It supports multiple database options including DynamoDB, Firestore, and MySQL. Each database repository is encapsulated within its respective folder under `repositories`, while configurations for each database can be found in `configs`.

## Project Structure

```
📁GO-SCRIPTING
├── .env-example               # Example environment variables
├── README.md                  # Project overview and setup instructions
├── 📁configs
│   ├── 📁dynamo_db
│   │   └── setup.go           # DynamoDB setup configuration
│   ├── 📁firestore
│   │   └── setup.go           # Firestore setup configuration
│   └── 📁mySql
│       └── setup.go           # MySQL setup configuration
├── 📁entities
│   └── early_booking_access.go  # Entity models
│   └── user.go                 # User entity
├── go.mod                     # Go modules definition
├── go.sum                     # Go modules lock file
├── 📁pkg
│   ├── 📁constant
│   │   └── pt_package.go      # Constants related to package types
│   │   └── time_format.go     # Time formatting constants
│   │   └── transaction.go     # Transaction related constants
│   ├── 📁csv_processor
│   │   └── reader.go          # CSV file reader package
│   └── 📁logger
│       └── logger.go          # Logging package
├── 📁repositories
│   ├── 📁dynamo_db
│   │   ├── 📁pt_packages
│   │   │   └── pt_package.go  # DynamoDB package repository
│   │   └── 📁transactions
│   │       └── cs_transactions_v6.go  # DynamoDB transaction repository
│   │       └── transactions_v1.go    # Previous version transactions
│   ├── 📁firestore
│   │   └── user.go            # Firestore user repository
│   └── 📁mySql
│       └── early_booking_access.go  # MySQL early booking access repository
└── 📁scripts
    ├── 📁early_book_access
    │   ├── main.go            # Main script for early booking access
    │   └── 📁service
    │       └── execute.go     # Service for executing the main script
    ├── 📁extend_expired_pt_package
    │   ├── example.csv        # Example CSV for extending packages
    │   ├── main.go            # Main script for extending expired packages
    │   └── 📁service
    │       └── execute_svc.go # Service for executing the main script
    └── 📁inject_free_days
        ├── example.csv        # Example CSV for injecting free days
        ├── main.go            # Main script for injecting free days
        └── 📁service
            └── execute_svc.go # Service for executing the main script
```

## Usage

### Setup

1. Clone the repository.
2. Review `.env-example` and configure environment variables as required.
3. Install dependencies using `go mod tidy`.
4. Update database configurations in `configs` folder.

### Running Scripts

- Navigate to the desired script folder under `scripts`.
- Modify CSV files as necessary for data input.
- Execute `go run main.go` to run the script.

## Features

- **Concurrency**: Utilizes Go routines for concurrent processing of data.
- **Modular Structure**: Organized into packages and modules for clarity and maintainability.
- **Database Support**: Includes support for DynamoDB, Firestore, and MySQL databases.
- **Logging**: Integrated logging functionality for debugging and monitoring.

## Contributing

Contributions are welcome! Please fork the repository and submit a pull request with your changes.

## Acknowledgments

- Inspired by real-world scenarios requiring bulk data processing and database operations.

---