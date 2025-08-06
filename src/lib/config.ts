// DeFi Portfolio Configuration
export const config = {
  // Chains
  defaultChain: 'ethereum',
  supportedChains: ['ethereum', 'polygon', 'arbitrum', 'optimism', 'polygonAmoy'],
  
  // API endpoints
  api: {
    baseUrl: process.env.NODE_ENV === 'production' 
      ? 'https://api.defiportfolio.com' 
      : 'http://localhost:3001',
    timeout: 10000,
  },
  
  // RPC URLs (TODO: Replace with your own)
  rpc: {
    ethereum: 'https://eth-mainnet.alchemyapi.io/v2/YOUR_API_KEY',
    polygon: 'https://polygon-mainnet.alchemyapi.io/v2/YOUR_API_KEY',
    arbitrum: 'https://arb-mainnet.g.alchemy.com/v2/YOUR_API_KEY',
    optimism: 'https://opt-mainnet.g.alchemy.com/v2/YOUR_API_KEY',
    polygonAmoy: 'https://rpc-amoy.polygon.technology',
  },
  
  // External APIs
  external: {
    coingecko: 'https://api.coingecko.com/api/v3',
    defillama: 'https://api.llama.fi',
  },
  
  // Feature flags
  features: {
    swapEnabled: false, // TODO: Enable when swap integration is ready
    notificationsEnabled: false, // TODO: Implement notifications
    bridgeEnabled: false, // TODO: Implement bridge functionality
  },
} as const;

export type SupportedChain = typeof config.supportedChains[number];