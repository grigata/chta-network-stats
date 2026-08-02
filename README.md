# CHTA Network Stats

A lightweight console and local web analyzer for the Cheetahcoin blockchain.

## Features

- Local CHTA Core RPC support
- Public API support
- Automatic source selection
- Network statistics
- Pool analysis
- LOW-DIFF detection
- Configurable scan depth
- Fast block scanner
- Progress indicator
- Local dark web dashboard
- Block table with direct Cheetahcoin Explorer links
- Pool, block-type and logarithmic difficulty charts

## Usage

```bash
CHTA-NetworkStats.exe
CHTA-NetworkStats.exe 100
CHTA-NetworkStats.exe 1000
CHTA-NetworkStats.exe web
CHTA-NetworkStats.exe web 100
```

Web mode scans the requested blocks, starts a private server at
`http://127.0.0.1:8080`, and opens the dashboard in the default browser.
Stop it with `Ctrl+C`. No blockchain data is sent anywhere by the dashboard.

## Example

```
============================================================
CHTA Network Stats v0.2.0
============================================================

Mode        : AUTO
Source      : Local CHTA Core RPC
Height      : 4971412
Blocks      : 100
```

## Build

```bash
go build -o bin/CHTA-NetworkStats.exe ./cmd/networkstats
```

## Roadmap

### v0.2.0
- [x] Local CHTA Core RPC
- [x] Public API
- [x] AUTO mode
- [x] Pool Analysis
- [x] Scan progress
- [x] Scan timing

### v0.3.0
- [x] Local web dashboard
- [x] Difficulty timeline
- [x] Pool and block-type analysis
- [x] Searchable block explorer table
- [ ] Streak Analysis
- [ ] Gap Analysis
- [ ] CSV export
- [ ] JSON export
- [ ] Watch mode

## License

MIT
