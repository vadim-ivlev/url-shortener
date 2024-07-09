#!/usr/bin/env bash




ENDPOINT=http://localhost:8080
ORIGINAL_URL=https://practicum.yandex.ru/

echo "Введите URL для сокращения."
echo "По умолчанию: $ORIGINAL_URL, нажмите Enter."
echo "Для выхода  Сtrl+C"
read -r url
if [ -n "$url" ]; then
    ORIGINAL_URL=$url
fi

echo "--------------------------------------------------"
echo "Doing POST $ENDPOINT -d '$ORIGINAL_URL'"
post_response=$( curl -s -fail -X POST "$ENDPOINT" -H "Content-Type: text/plain;" -d "$ORIGINAL_URL" )
printf "\npost_response>>>>>>>>>>>>\n $post_response \npost_response<<<<<<<<<<<<\n\n"

short_url=$(echo "$post_response" | tail -n 1)
printf "\nshort_url>>>>>>>>>>>>\n $short_url \nshort_url<<<<<<<<<<<<\n\n"

echo "--------------------------------------------------"
echo "Doing GET $short_url"
get_response=$( curl -s -fail -X GET $short_url )
printf "\nget_response>>>>>>>>>>>>\n $get_response \nget_response<<<<<<<<<<<<\n\n"