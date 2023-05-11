docker run -d --name rediseach \
-p 6379:6379 redis/redis-stack-server:latest

echo "Use RediSearch CLI with"
echo "docker exec -it rediseach redis-cli -p 6379"