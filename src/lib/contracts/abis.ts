// Standard ERC20 Token ABI (for allowance and transfer operations)
export const ERC20_ABI = [
  {
    "inputs": [
      { "name": "spender", "type": "address" },
      { "name": "amount", "type": "uint256" }
    ],
    "name": "approve",
    "outputs": [{ "name": "", "type": "bool" }],
    "stateMutability": "nonpayable",
    "type": "function"
  },
  {
    "inputs": [
      { "name": "owner", "type": "address" },
      { "name": "spender", "type": "address" }
    ],
    "name": "allowance",
    "outputs": [{ "name": "", "type": "uint256" }],
    "stateMutability": "view",
    "type": "function"
  },
  {
    "inputs": [{ "name": "account", "type": "address" }],
    "name": "balanceOf",
    "outputs": [{ "name": "", "type": "uint256" }],
    "stateMutability": "view",
    "type": "function"
  },
  {
    "inputs": [],
    "name": "decimals",
    "outputs": [{ "name": "", "type": "uint8" }],
    "stateMutability": "view",
    "type": "function"
  },
  {
    "inputs": [],
    "name": "symbol",
    "outputs": [{ "name": "", "type": "string" }],
    "stateMutability": "view",
    "type": "function"
  },
  {
    "inputs": [],
    "name": "name",
    "outputs": [{ "name": "", "type": "string" }],
    "stateMutability": "view",
    "type": "function"
  }
] as const;

// Aave V3 Pool ABI (for claiming rewards)
export const AAVE_V3_POOL_ABI = [
  {
    "inputs": [
      { "name": "assets", "type": "address[]" },
      { "name": "to", "type": "address" }
    ],
    "name": "claimRewards",
    "outputs": [{ "name": "", "type": "uint256[]" }],
    "stateMutability": "nonpayable",
    "type": "function"
  },
  {
    "inputs": [
      { "name": "user", "type": "address" },
      { "name": "assets", "type": "address[]" }
    ],
    "name": "getUserRewards",
    "outputs": [{ "name": "", "type": "uint256[]" }],
    "stateMutability": "view",
    "type": "function"
  }
] as const;

// Compound V3 Comet ABI (for claiming COMP rewards)
export const COMPOUND_V3_COMET_ABI = [
  {
    "inputs": [
      { "name": "src", "type": "address" },
      { "name": "shouldAccrue", "type": "bool" }
    ],
    "name": "claim",
    "outputs": [],
    "stateMutability": "nonpayable",
    "type": "function"
  },
  {
    "inputs": [{ "name": "account", "type": "address" }],
    "name": "baseTrackingAccrued",
    "outputs": [{ "name": "", "type": "uint64" }],
    "stateMutability": "view",
    "type": "function"
  }
] as const;

// Uniswap V3 Position Manager ABI (for collecting fees)
export const UNISWAP_V3_POSITION_MANAGER_ABI = [
  {
    "inputs": [
      {
        "components": [
          { "name": "tokenId", "type": "uint256" },
          { "name": "recipient", "type": "address" },
          { "name": "amount0Max", "type": "uint128" },
          { "name": "amount1Max", "type": "uint128" }
        ],
        "name": "params",
        "type": "tuple"
      }
    ],
    "name": "collect",
    "outputs": [
      { "name": "amount0", "type": "uint256" },
      { "name": "amount1", "type": "uint256" }
    ],
    "stateMutability": "payable",
    "type": "function"
  }
] as const;