# CHTA Network Stats

A lightweight command-line analyzer for the Cheetahcoin blockchain.

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

## Usage

```bash
CHTA-NetworkStats.exe
CHTA-NetworkStats.exe 100
CHTA-NetworkStats.exe 1000
```

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
- [ ] Streak Analysis
- [ ] Difficulty Analysis
- [ ] Gap Analysis
- [ ] CSV export
- [ ] JSON export
- [ ] Watch mode

## License

MIT