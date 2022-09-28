abigen:
	abigen --out ./contracts/seaport/seaport.go --pkg seaport --type Seaport --abi ./contracts/seaport/abi.json

build:
	go build -o export-nft-data main.go