#!/bin/bash

# Define paths
input_folder="./gif"
output_folder="./gif/compressed"
mkdir -p "$output_folder"

compress_gif() {
    input_path="$1"
    output_path="$2"

    # Optimize the GIF using gifsicle
    gifsicle --optimize=3 --threads=4 --colors 64 --lossy "$input_path" -o "$output_path"
}

# Process each GIF in the input folder
for gif_file in "$input_folder"/*.gif; do
    gif_name=$(basename "$gif_file")
    output_path="$output_folder/$gif_name"
    compress_gif "$gif_file" "$output_path"
done

echo "Compression completed."c