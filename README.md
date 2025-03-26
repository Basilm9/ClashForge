# ClashForge 🏰⚔️

A high-performance Clash Royale style game server implementation using WebSockets in Go.

## Overview

ClashForge is a lightweight, scalable WebSocket server that powers real-time card battling games inspired by Clash Royale. This server handles player matchmaking, battle synchronization, card deployment, and game state management with minimal latency.

## Features

- **Real-time WebSocket Communication**: Fast, bidirectional communication between clients and server
- **Match Management**: Dynamic matchmaking and match state synchronization
- **Card System**: Comprehensive card management using authentic data
- **Scalable Architecture**: Designed to handle thousands of concurrent connections
- **Cross-platform Compatibility**: Works on all major operating systems
- **Lightweight**: Minimal resource footprint
- **Battle Replay System**: Record and replay past battles

## Data Source

The game data (cards, arenas, etc.) was sourced from the [smlbiobot/cr-csv](https://github.com/smlbiobot/cr-csv) repository, which provides an extensive collection of Clash Royale data in CSV format.

## Prerequisites

- Go 1.19 or higher
- Working network connection (for websocket communications)

## Quick Start

```bash
# Clone this repository
git clone https://github.com/yourusername/clashforge.git
cd clashforge

# Build the server
go build -o out && ./out
```

By default, the server will start on port 8080. You can access the WebSocket endpoint at `ws://localhost:8080/game`.

## Configuration

Server configuration is handled via environment variables or a config file:

```bash
# Set custom port
PORT=9000 ./out

# Use configuration file
CONFIG_PATH=./config.json ./out
```

## Client Connection

Connect to the WebSocket server from client applications:

```javascript
// JavaScript example
const socket = new WebSocket("ws://localhost:8080/game");

socket.onopen = () => {
  console.log("Connected to ClashForge server");

  // Join matchmaking queue
  socket.send(
    JSON.stringify({
      action: "queue_join",
      deck: [
        "knight",
        "archers",
        "fireball",
        "goblin_barrel",
        "skeleton_army",
        "inferno_tower",
        "log",
        "princess",
      ],
    })
  );
};

socket.onmessage = (event) => {
  const message = JSON.parse(event.data);
  console.log("Received:", message);

  // Handle different message types
  switch (message.type) {
    case "match_found":
      // Match found logic
      break;
    case "card_deployed":
      // Card deployment logic
      break;
    // Other message types...
  }
};
```

## API Documentation

The WebSocket API documentation can be found in the [docs/api.md](docs/api.md) file.

## Development

### Project Structure

```
clashforge/
├── cmd/
│   └── server/
│       └── main.go      # Entry point
├── internal/
│   ├── game/            # Game logic
│   ├── matchmaking/     # Matchmaking system
│   ├── models/          # Data models
│   └── websocket/       # WebSocket handling
├── data/                # Game data from cr-csv
└── config/              # Configuration files
```

### Adding New Cards

1. Update the CSV data files in the `data/` directory
2. Restart the server to load the new card data

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments

- [smlbiobot/cr-csv](https://github.com/smlbiobot/cr-csv) for providing the comprehensive Clash Royale data
- The Go WebSocket community for excellent libraries and examples
