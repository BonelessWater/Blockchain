# Blockchain Algorithm

## TO-DO
:x: RSA encryption algorithm such that transactions can be securely signed
:x: Blockchain algorithm where hashes are created with their respective nonce value
:x: Establish P2P network
:x: Securely add crypo currency to the miner
:x: Broadcast new blocks to connected nodes, nodes recursively broadcast a block until all nodes are updated. Note: Each node must verify transaction integrity and hash integrity


## Notes
Blockchain technology is a type of technology that is used to create and maintain a decentralized ledger.

Each block will be represented as a struct called 'Block' and is composed of the following: an index, a timestamp, relevant data (BPM), a hash and the hash of the last block in the chain. Hashes are useful because when a block is composed of millions and millions of data, a hash provides a convenient way to secure the integrity of blocks without having to take into account the millions of data points of all the previous blocks. These hashes are SHA256 identifiers which are overwhelmingly secure for the purposes of this project.
