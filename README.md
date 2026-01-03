# Cat Image Viewer

A desktop GUI application built with Go and Gio UI that fetches and displays random cat images from [cataas.com](https://cataas.com/).

## Prerequisites

- Go 1.25 or higher

## Getting Started

### Clone the Repository

```bash
git clone <repository-url>
cd go-gui
```

### Build the Application

```bash
go build -o build/catfetch ./cmd/catfetch
```

### Run the Application

From the project directory:

```bash
./build/catfetch
```

Or run directly without building:

```bash
go run ./cmd/catfetch/main.go
```

### Optional: Install System-wide

To run the application from anywhere without `./`:

**Linux/macOS (user-local install):**
```bash
cp build/catfetch ~/bin/catfetch
# Ensure ~/bin is in your PATH
```

**Linux (system-wide install):**
```bash
sudo cp build/catfetch /usr/local/bin/catfetch
```

Then you can run it from anywhere:
```bash
catfetch
```

## Usage

Launch the application and click the "Fetch Image" button to load a random cat picture. The image will automatically scale to fit the window while maintaining its aspect ratio.

## Roadmap

- **Cat History**: Browse previously fetched cat images
- **Text Overlays**: Add custom text overlays to cat images
- **Tag Search**: Search for cats by specific tags
- **Image Filters**: Add sliders and options to apply filters (sepia, blur, brightness, etc.) using cataas API parameters

## License

MIT License - see [LICENSE](LICENSE) file for details.

Copyright (c) 2026 NovelGit LLC