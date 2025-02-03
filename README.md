<div align="center">
  <img src="WikiStream.svg" width="200" height="200" alt="WikiStream Logo">
  
  # WikiStream
  
  *Discover recent Wikipedia changes in real-time*
</div>

## Installation

Copy repo:
```bash
git clone https://github.com/imNauryzbay/wikistream
```

## Run

Open cloned repo:
```bash
cd wikistream/cmd
```

Take token from Discord:
```bash
go run main.go -token TOKEN
```

## Usage

After running the bot, two commands are available:

- `!recent`: Shows up to 5 last changes (waits 5 seconds)
- `!setLang`: Change language (default: en)
  - Example: `!setLang ru`