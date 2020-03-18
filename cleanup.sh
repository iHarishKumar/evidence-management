docker stop $(docker ps -aq)
docker rm -f $(docker ps -aq)
docker volume prune
docker ps -aq
rm -rf fabric-client-kv-*
rm -rf crypto-config/*
rm -rf channel-artifacts/*