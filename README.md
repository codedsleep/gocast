# gocast

A fast, lightweight command-line weather application written in Go that provides real-time weather information using the Open-Meteo API.

## Features

- **Real-time Weather Data**: Get current weather conditions for any location worldwide
- **Historical Weather**: View weather data from the past 24 hours or 7 days
- **Location Support**: Search by city name with optional country code specification
- **Cross-platform**: Available as binary executable and packaged in Deb and RPM formats
- **Fast Performance**: Lightweight Go implementation with minimal dependencies
- **Free API**: Uses the Open-Meteo API which requires no API key

## Installation

### Binary Release
Download the latest binary from the releases page and run:
```bash
./gocast
```

### Package Installation
- **Debian/Ubuntu**: Install the `.deb` package
- **RPM-based systems**: Install the `.rpm` package

## Usage

### Basic Usage
```bash
gocast <location>
```

### Command Options
- `gocast <city>` - Get current weather for a city
- `gocast <city> <country_code>` - Get weather with country specification
- `gocast -24h <city>` - Show weather from the past 24 hours
- `gocast -7d <city>` - Show weather from the past 7 days

## Examples

### Current Weather for London
```bash
gocast london
```

![swappy-20250708_043811](https://github.com/user-attachments/assets/dce4e63f-8ddc-4001-bd11-60892d92b7bc)

### Weather with Country Code
```bash
gocast manchester gb
```

![swappy-20250708_043930](https://github.com/user-attachments/assets/cb57ad53-8cec-4b15-a386-7f3bb972e9fe)

### 24-Hour Weather History
```bash
gocast -24h london
```

![swappy-20250708_044020](https://github.com/user-attachments/assets/bd43eda0-e8b4-4bb3-8ad8-cca78ac5edb3)

### 7-Day Weather History
```bash
gocast -7d newyork
```

![swappy-20250708_044204](https://github.com/user-attachments/assets/75fef58c-dd29-4d79-9f95-ee69f1864b0a)

## Building from Source

Ensure you have Go 1.24.4 or later installed:

```bash
go build -o gocast main.go
```

## API

This application uses the [Open-Meteo API](https://open-meteo.com/), which provides free weather data without requiring an API key.

## License

This project is open source. See the repository for license details.
