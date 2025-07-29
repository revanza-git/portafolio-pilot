// API Configuration
export const API_CONFIG = {
  // Backend base URL - uses environment variable or fallback
  baseUrl: import.meta.env.VITE_API_BASE_URL || 
    (import.meta.env.PROD 
      ? 'https://api.your-domain.com' 
      : 'http://localhost:3000'),
  
  // API version
  version: 'v1',
  
  // Endpoints
  endpoints: {
    // Auth
    authNonce: '/api/v1/auth/nonce',
    authVerify: '/api/v1/auth/verify',
    
    // Portfolio
    portfolioBalances: (address: string) => `/api/v1/portfolio/${address}/balances`,
    portfolioHistory: (address: string) => `/api/v1/portfolio/${address}/history`,
    
    // Transactions
    transactions: (address: string) => `/api/v1/transactions/${address}`,
    approvals: (address: string) => `/api/v1/transactions/${address}/approvals`,
    revokeApproval: (address: string, token: string) => `/api/v1/transactions/${address}/approvals/${token}`,
    
    // Yield
    yieldPools: '/api/v1/yield/pools',
    yieldPositions: (address: string) => `/api/v1/yield/positions/${address}`,
    
    // Bridge & Swap
    bridgeRoutes: '/api/v1/bridge/routes',
    swapQuote: '/api/v1/swap/quote',
    
    // Analytics
    pnlExport: (address: string) => `/api/v1/analytics/pnl/${address}`,
    
    // Alerts
    alerts: '/api/v1/alerts',
    alert: (id: string) => `/api/v1/alerts/${id}`,
    
    // Watchlists
    watchlists: '/api/v1/watchlists',
    watchlist: (id: string) => `/api/v1/watchlists/${id}`,
  }
};

// Helper to build full URL
export function buildApiUrl(endpoint: string): string {
  return `${API_CONFIG.baseUrl}${endpoint}`;
}