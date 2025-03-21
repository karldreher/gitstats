commit=$1
repo=$2

if [ -z "$commit" ] || [ -z "$repo" ]; then
    echo "Usage: $0 <commit> <repo>"
    exit 1
fi
# TODO: ENV that url
curl -s -X POST -H "Content-Type: application/json" -d "{\"commit\": \"$commit\", \"repo\": \"$repo\"}" http://localhost:8000/api/v1/commit

