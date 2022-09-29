abigen:
	abigen --out ./contracts/seaport/seaport.go --pkg seaport --type Seaport --abi ./contracts/seaport/abi.json
	abigen --out ./contracts/erc721/erc721.go --pkg erc721 --type ERC721 --abi ./contracts/erc721/abi.json

build:
	go build -o export-nft-data main.go