#!/bin/bash

# Script to recursively replace string values beginning with 'src="' and ending with '"' in all .html files.

# Find all .html files recursively
find /C/Users/user-claude/Documents/wellington/go-class-server/test -name "*.html" -print0 | while IFS= read -r -d $'\0' file; do
  # Check if the file exists
  if [ -f "$file" ]; then
    # Use sed to replace the string
    sed -i "s/li\.open>button>span\.closed/{display:none}/s/li\.open>button>span\.closed/{display:block}/" "$file"
    # sed -i "s/src=\"https\:\/\//src=\"\/websites\//g" "$file"
    # sed -i "s/href=\"http\:\/\/localhost\:22022\/Public\/Downloads\/2025\/mdn2025\//href=\"\/websites\//g" "$file"
    # sed -i "s/src=\"http\:\/\/localhost\:22022\/Public\/Downloads\/2025\/mdn2025\//src=\"\/websites\//g" "$file"
    # sed -i "s/src=\"http\:\/\/localhost\:22022\/Public\/Downloads\/2025\/mdn2025\//src=\"\/websites\//g" "$file"
    # sed -i "s/main.b68f5a8c.css/main.27d5b2a8.css/g" "$file"
    echo "Replaced in: $file"
    
  fi
done

echo "Finished processing .html files."
