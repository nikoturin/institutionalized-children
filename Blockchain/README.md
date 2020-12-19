
Steps to add library IPFS

1.- Please check URL:

https://github.com/ipfs/go-ipfs-api

2.- Test it installing:

https://github.com/ipfs/go-ipfs-api

Run test code, "go run <name file golang>"

3.- Installing Chaincode:

peer chaincode install -n ipfs -v 1 -l golang -p cars/

4.- Instance Chaincode:

peer chaincode instantiate -o orderer.example.com:7050 -C general -n ccipfs -l golang -v 0 -c '{"Args":[]}' --tls true --cafile $ORDERER_CA 

Note: Try to check any issue with crypto/ library
