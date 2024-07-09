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



response=$( curl -s -fail -X POST "$ENDPOINT" -H "Content-Type: text/plain;" -d "$ORIGINAL_URL" )
short_url=$(echo "$response" | tail -n 1)


printf "\nresponse>>>>>>>>>>>>\n$response\nresponse<<<<<<<<<<<<\n\n"
printf "\nshort_url>>>>>>>>>>>>\n$short_url\nshort_url<<<<<<<<<<<<\n\n"

response=$( curl -s -fail -X GET $short_url )

printf "\nresponse>>>>>>>>>>>>\n$response\nresponse<<<<<<<<<<<<\n\n"