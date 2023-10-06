# YouTube Video Clipper

This project allows you to download a YouTube video and clip it into multiple segments based on provided time ranges. It's built in Go and makes use of multi-threading for improved performance.

## Pre-requisites

- Go (1.16 or higher)
- FFmpeg

## Usage

1. Clone the repository:

    ```bash
    git clone https://github.com/zucchivan/yt-video-clipper.git
    ```

2. Navigate to the project's cmd directory:

    ```bash
    cd yt-video-clipper/cmd
    ```

3. Run the program with the `-url` and `-timePairs` flags:

    ```bash
    go run cmd/main.go -url="https://www.youtube.com/watch?v=dQw4w9WgXcQ" -timePairs="00:00:10-00:00:20,00:01:00-00:01:10"
    ``````

    The time pairs should be in the `HH:mm:SS` format separated by a `-`.