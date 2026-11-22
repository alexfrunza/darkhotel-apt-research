#!/usr/bin/env bash
PASSWORD="password"
OUTPUT_ZIP="top_secret.zip"

if [[ -z "$PASSWORD" ]]; then
  echo "Usage: $0 <password> [output_zip]" >&2
  exit 1
fi

TMP_DIR=$(mktemp -d)
trap 'rm -rf "$TMP_DIR"' EXIT

# Generate 3 files (~50 MB each) containing "file1", "file2", "file3" repeated on new lines
for i in 1 2 3; do
  TARGET="$TMP_DIR/file${i}.txt"
  yes "file${i}" | head -n 8738133 > "$TARGET"
  echo "Created $TARGET ($(du -h "$TARGET" | cut -f1))"
done

# Create password-protected ZIP using zip's built-in encryption (-P)
ABS_OUT="$(cd "$(dirname "$OUTPUT_ZIP")" && pwd)/$(basename "$OUTPUT_ZIP")"
(cd "$TMP_DIR" && zip -0 -q -P "$PASSWORD" "$ABS_OUT" file1.txt file2.txt file3.txt)

echo "Password-protected ZIP written to: $OUTPUT_ZIP"
