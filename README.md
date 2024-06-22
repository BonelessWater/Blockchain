# Blockchain Algorithm
- This source code is my attempt at making my own crypto currency. Such cyrpto currency will be designed to be traded for items on x market. Whether this currency will appreciate at all is unknown to me. Regardless, the currency must also be able to be mined. Once mined, the miner will obtain a certain amount of currency to grow the money supply. Since crypto currencies such as these do not have a central financial institution, money can only be added to the money supply through mining. Mining is the process of discovering some number called a nonce (number only used once) that when appended to the block, will result in a sha256 hash that satisfies a certain condition. In this case, and in nearly all cases, the condition will be that the sha256 hash must have a certain amount of leading zeros. The precise amount of leading zeros that a hash should have will be changed throughout the currency's lifetime to make sure money is added to the money supply at a near-constant rate (regardless of the amount of miners in the network).
- Additionally, the currency must be cryptographically secure, meaning that ill-intending actors cannot simply give themselves money without being verified. I will utilize RSA (Rivest–Shamir–Adleman) encryption to ensure that transactions are securely signed. Another method to stop bad actors is to append the value of the hash of the latest block each new block. In this way, bad actors cannot change the order of blocks because that would change the hollistic hash of the chain. Nodes will never trust blockchains by themselves and always take an additional step to verify them before mining.

## TO-DO
- :white_check_mark: Blockchain
  - :white_check_mark: Genesis block
  - :white_check_mark: Relevant data: ID, Timestamp, Transactions, Merkleroot, Prevhash, Target difficulty, Nonce, Hash
  - :white_check_mark: Assign public key pairs to node
  - :white_check_mark: Create successful Block hashes based on target
  - :white_check_mark: RSA encryption algorithm
  - :white_check_mark: Make transaction signitures
  - :white_check_mark: Verify transaction signitures
  - :white_check_mark: Merkle root algorithm
- :white_check_mark: Establish P2P network
  - :white_check_mark: Establish local host connection
  - :white_check_mark: Send messages via JSON
  - :white_check_mark: Receive messages via JSON
  - :white_check_mark: Create and utilize trust threshold
  
## Notes
Blockchain technology is a type of technology that is used to create and maintain a decentralized ledger.

Each block will be represented as a struct called 'Block' and is composed of the following: an index, a timestamp, relevant data (transactions), a hash and the hash of the last block in the chain. Hashes are useful because when a block is composed of millions and millions of data, a hash provides a convenient way to secure the integrity of blocks without having to take into account the millions of data points of all the previous blocks. These hashes are SHA256 identifiers which are overwhelmingly secure for the purposes of this project.
