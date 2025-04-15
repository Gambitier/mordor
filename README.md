# Mordor

A collection of cryptographic and security utilities.

## Available Tools

### 1. RSA Key Generator

A utility to generate RSA key pairs with timestamp-based naming.

#### Installation

```bash
go install github.com/Gambitier/mordor/cmd/RSAKeyGenerator@latest
```

#### Usage

```bash
# Generate keys with default settings (2048 bits)
RSAKeyGenerator -output ./keys
```

The RSA key generator will:
- Create the output directory if it doesn't exist
- Generate a private key in PEM format
- Extract the corresponding public key
- Name the files using UTC timestamp (format: YYYYMMDDHHMMSS)

#### Features
- 2048-bit RSA keys
- PEM format output
- Uses OpenSSL for secure key generation
- Automatic output directory creation
- Timestamp-based file naming

## Project Structure

```
cmd/
  ├── RSAKeyGenerator/    # RSA key pair generation utility
  └── ... (more tools to come)
```