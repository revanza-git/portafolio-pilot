// Contract addresses by chain ID
export const CONTRACT_ADDRESSES = {
  // Ethereum Mainnet (Chain ID: 1)
  1: {
    // Aave V3
    AAVE_V3_POOL: '0x87870Bca3F3fD6335C3F4ce8392D69350B4fA4E2',
    AAVE_V3_REWARDS_CONTROLLER: '0x8164Cc65827dcFe994AB23944CBC90e0aa80bFcb',
    
    // Compound V3
    COMPOUND_V3_USDC: '0xc3d688B66703497DAA19211EEdff47f25384cdc3',
    COMPOUND_V3_ETH: '0xA17581A9E3356d9A858b789D68B4d866e593aE94',
    
    // Uniswap V3
    UNISWAP_V3_POSITION_MANAGER: '0xC36442b4a4522E871399CD717aBDD847Ab11FE88',
    UNISWAP_V3_FACTORY: '0x1F98431c8aD98523631AE4a59f267346ea31F984',
    
    // Common tokens
    USDC: '0xA0b86991c431E4dFe7bb8E5f2D5E8b8A8A8b3c8B',
    USDT: '0xdAC17F958D2ee523a2206206994597C13D831ec7',
    WETH: '0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2',
    WBTC: '0x2260FAC5E5542a773Aa44fBCfeDf7C193bc2C599',
    UNI: '0x1f9840a85d5aF5bf1D1762F925BDADdC4201F984',
    AAVE: '0x7Fc66500c84A76Ad7e9c93437bFc5Ac33E2DDaE9',
  },
  
  // Polygon (Chain ID: 137)
  137: {
    // Aave V3
    AAVE_V3_POOL: '0x794a61358D6845594F94dc1DB02A252b5b4814aD',
    AAVE_V3_REWARDS_CONTROLLER: '0x929EC64c34a17401F460460D4B9390518E5B473e',
    
    // Common tokens
    USDC: '0x2791Bca1f2de4661ED88A30C99A7a9449Aa84174',
    USDT: '0xc2132D05D31c914a87C6611C10748AEb04B58e8F',
    WETH: '0x7ceB23fD6bC0adD59E62ac25578270cFf1b9f619',
    WMATIC: '0x0d500B1d8E8eF31E21C99d1Db9A6444d3ADf1270',
  },
  
  // Arbitrum (Chain ID: 42161)
  42161: {
    // Aave V3
    AAVE_V3_POOL: '0x794a61358D6845594F94dc1DB02A252b5b4814aD',
    AAVE_V3_REWARDS_CONTROLLER: '0x929EC64c34a17401F460460D4B9390518E5B473e',
    
    // Common tokens
    USDC: '0xaf88d065e77c8cC2239327C5EDb3A432268e5831',
    USDT: '0xFd086bC7CD5C481DCC9C85ebE478A1C0b69FCbb9',
    WETH: '0x82aF49447D8a07e3bd95BD0d56f35241523fBab1',
    ARB: '0x912CE59144191C1204E64559FE8253a0e49E6548',
  },
  
  // Optimism (Chain ID: 10)
  10: {
    // Aave V3
    AAVE_V3_POOL: '0x794a61358D6845594F94dc1DB02A252b5b4814aD',
    AAVE_V3_REWARDS_CONTROLLER: '0x929EC64c34a17401F460460D4B9390518E5B473e',
    
    // Common tokens
    USDC: '0x0b2C639c533813f4Aa9D7837CAf62653d097Ff85',
    USDT: '0x94b008aA00579c1307B0EF2c499aD98a8ce58e58',
    WETH: '0x4200000000000000000000000000000000000006',
    OP: '0x4200000000000000000000000000000000000042',
  },
} as const;

// Helper function to get contract address by chain
export function getContractAddress(
  chainId: keyof typeof CONTRACT_ADDRESSES,
  contractName: keyof typeof CONTRACT_ADDRESSES[1]
): string {
  const addresses = CONTRACT_ADDRESSES[chainId];
  if (!addresses) {
    throw new Error(`Unsupported chain ID: ${chainId}`);
  }
  
  const address = addresses[contractName];
  if (!address) {
    throw new Error(`Contract ${contractName} not found on chain ${chainId}`);
  }
  
  return address;
}

// Common token addresses helper
export function getTokenAddress(chainId: number, symbol: string): string {
  const symbolUpper = symbol.toUpperCase();
  return getContractAddress(chainId as keyof typeof CONTRACT_ADDRESSES, symbolUpper as any);
}