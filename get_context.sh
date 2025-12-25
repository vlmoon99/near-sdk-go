#!/usr/bin/env bash

OUTPUT="go_context.txt"

> "$OUTPUT"

echo "### GO PROJECT CONTEXT" >> "$OUTPUT"
echo "### Root: $(pwd)" >> "$OUTPUT"
echo "### Generated at: $(date)" >> "$OUTPUT"
echo "" >> "$OUTPUT"

find . \
  -type d -name examples -prune -o \
  -type f \( \
    -name "*.go" -o \
    -name "go.mod" -o \
    -name "go.sum" -o \
    -name "README.md" \
  \) -print | sort | while read -r file; do

  echo "==================================================" >> "$OUTPUT"
  echo "FILE: ${file#./}" >> "$OUTPUT"
  echo "==================================================" >> "$OUTPUT"
  echo "" >> "$OUTPUT"

  cat "$file" >> "$OUTPUT"
  echo -e "\n\n" >> "$OUTPUT"
done

echo "Context written to $OUTPUT"
