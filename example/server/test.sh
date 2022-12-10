hostname=localhost
declare -A keyvaluemap

put_cache() {
    curl --location --request PUT "http://${hostname}:3000/api/v1/cache/$1/$2" -w "" -s
}

get_cache() {
    curl --location --request GET "http://${hostname}:3000/api/v1/cache/$1" -w "" -s
}

requestNum=10*10*10

for ((n=0;n<$requestNum;n++))
do
    key=$(cat /dev/urandom | tr -dc '[:alpha:]' | fold -w ${1:-32} | head -n 1)
    value=$(cat /dev/urandom | tr -dc '[:alpha:]' | fold -w ${1:-128} | head -n 1)

    keyvaluemap[$key]=$value
    echo $key
done


echo "running PUT cache"
for i in "${!keyvaluemap[@]}"
do
  key=$i
  value=${keyvaluemap[$i]}
  
  put_cache $key $value &
done
echo "end of running PUT cache"

sleep 2

echo "running GET cache"
declare -i match=0
declare -i mismatch=0

for i in "${!keyvaluemap[@]}"
do
  key=$i
  value=${keyvaluemap[$i]}

  resultKey=$(get_cache $key)

  if [ "$value" = "$resultKey" ]; then
    match=$match+1
    else
        mismatch=$mismatch+1
    fi
done
echo "end of running GET cache"

echo "Match $match"
echo "Mismatch $mismatch"

# DDOS
# ddosify -m PUT -p HTTP -t mnvs:3000/api/v1/cache/{{_randomUUID}}/{{_randomFullName}} -l linear -d 10 -n 30000
