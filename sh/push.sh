#!/bin/bash

echo "Введите комментарий к коммиту:"
read msg
echo "Комментарий = $msg."
echo 


git add -A .
git commit -m "$msg."

git push origin --all 

git push origin --tags 
