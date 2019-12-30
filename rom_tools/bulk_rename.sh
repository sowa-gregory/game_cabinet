for file in *.zip.done; do
  mv "$file" "$(basename "$file" .done).zip"
done
