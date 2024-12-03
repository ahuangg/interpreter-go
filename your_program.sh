set -e 

(
  cd "$(dirname "$0")" 
  go build -o /tmp/interpreter-target ./cmd/interpreter
)

exec /tmp/interpreter-target "$@"
