commit=$1
repo=$2

if [ -z "$commit" ] || [ -z "$repo" ]; then
    echo "Usage: $0 <commit> <repo>"
    exit 1
fi

curl -s -X POST -H "Content-Type: application/json" -H "x-api-key: $GITSTATS_API_KEY" -d "{\"commit\": \"$commit\", \"repo\": \"$repo\"}" http://$GITSTATS_HOST/api/v1/commit

