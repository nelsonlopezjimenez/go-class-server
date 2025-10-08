#!/bin/bash

# Script to recursively replace string values beginning with 'src="' and ending with '"' in all .html files.

# Find all .html files recursively
find . -name "*.html" -print0 | while IFS= read -r -d $'\0' file; do
  # Check if the file exists
  if [ -f "$file" ]; then
    # Use sed to replace the string
    sed -i "s/src=\"[^\"]*\"/src=\"new_value\"/g" "$file"
    echo "Replaced in: $file"
  fi
done

echo "Finished processing .html files."