#!/bin/bash

while true; do
    new_clipboard=$(xclip -o)
    if [[ "$new_clipboard" != "$old_clipboard" ]]; then
        # Ваша команда для выполнения при изменении буфера обмена
        echo "Clipboard contents changed: $new_clipboard"
        old_clipboard="$new_clipboard"
    fi
    sleep 1
done
