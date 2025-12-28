#!/usr/bin/env bash

OUTPUT="go_context.txt"

> "$OUTPUT"

echo "### GO PROJECT CONTEXT" >> "$OUTPUT"
echo "### Root: ." >> "$OUTPUT"
echo "" >> "$OUTPUT"

find . \
  -type d -name examples -prune -o \
  -type f \( \
    -name "*.go" -o \
    -name "go.mod" -o \
    -name "go.sum" -o \
    -name "README.md" \
  \) -print | sort | while read -r file; do

  rel="${file#./}"

  echo "==================================================" >> "$OUTPUT"
  echo "FILE: $rel" >> "$OUTPUT"
  echo "==================================================" >> "$OUTPUT"
  echo "" >> "$OUTPUT"

  cat "$file" >> "$OUTPUT"
  printf "\n\n" >> "$OUTPUT"
done

echo "Context written to $OUTPUT"
