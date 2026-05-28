commit=$1
repo=$2
author=$3

if [ -z "$commit" ] || [ -z "$repo" ] || [ -z "$author" ]; then
    echo "Usage: $0 <commit> <repo> <author>"
    exit 1
fi

curl -s -X POST \
    -H "Content-Type: application/json" \
    -H "x-api-key: $GITSTATS_API_KEY" \
    -d "{\"commit\": \"$commit\", \"repo\": \"$repo\", \"author\": \"$author\"}" \
    http://$GITSTATS_HOST/api/v1/commit
