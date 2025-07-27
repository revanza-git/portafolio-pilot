// API clients for external data sources

const COINGECKO_BASE_URL = 'https://api.coingecko.com/api/v3';
const DEFILLAMA_BASE_URL = 'https://api.llama.fi';

// Rate limiting helper
class RateLimiter {
  private requests: number[] = [];
  private maxRequests: number;
  private timeWindow: number;

  constructor(maxRequests: number, timeWindowMs: number) {
    this.maxRequests = maxRequests;
    this.timeWindow = timeWindowMs;
  }

  async throttle(): Promise<void> {
    const now = Date.now();
    this.requests = this.requests.filter(time => now - time < this.timeWindow);
    
    if (this.requests.length >= this.maxRequests) {
      const oldestRequest = Math.min(...this.requests);
      const waitTime = this.timeWindow - (now - oldestRequest);
      await new Promise(resolve => setTimeout(resolve, waitTime));
    }
    
    this.requests.push(now);
  }
}

// CoinGecko API client
export class CoinGeckoClient {
  private rateLimiter = new RateLimiter(50, 60000); // 50 requests per minute
  private cache = new Map<string, { data: any; timestamp: number }>();
  private cacheTimeout = 5 * 60 * 1000; // 5 minutes

  private async request<T>(endpoint: string): Promise<T> {
    await this.rateLimiter.throttle();
    
    const cacheKey = endpoint;
    const cached = this.cache.get(cacheKey);
    
    if (cached && Date.now() - cached.timestamp < this.cacheTimeout) {
      return cached.data;
    }

    try {
      const response = await fetch(`${COINGECKO_BASE_URL}${endpoint}`);
      
      if (!response.ok) {
        throw new Error(`CoinGecko API error: ${response.status}`);
      }
      
      const data = await response.json();
      this.cache.set(cacheKey, { data, timestamp: Date.now() });
      
      return data;
    } catch (error) {
      console.error('CoinGecko API error:', error);
      
      // Return cached data if available
      if (cached) {
        console.warn('Using stale CoinGecko data due to API error');
        return cached.data;
      }
      
      throw error;
    }
  }

  async getTokenPrices(tokenIds: string[]): Promise<Record<string, { usd: number; usd_24h_change: number }>> {
    const ids = tokenIds.join(',');
    return this.request(`/simple/price?ids=${ids}&vs_currencies=usd&include_24hr_change=true`);
  }

  async getTokenInfo(tokenId: string) {
    return this.request(`/coins/${tokenId}?localization=false&tickers=false&market_data=true&community_data=false&developer_data=false`);
  }

  async getPriceHistory(tokenId: string, days: number = 7) {
    return this.request(`/coins/${tokenId}/market_chart?vs_currency=usd&days=${days}`);
  }
}

// DefiLlama API client
export class DefiLlamaClient {
  private rateLimiter = new RateLimiter(300, 60000); // 300 requests per minute
  private cache = new Map<string, { data: any; timestamp: number }>();
  private cacheTimeout = 10 * 60 * 1000; // 10 minutes

  private async request<T>(endpoint: string): Promise<T> {
    await this.rateLimiter.throttle();
    
    const cacheKey = endpoint;
    const cached = this.cache.get(cacheKey);
    
    if (cached && Date.now() - cached.timestamp < this.cacheTimeout) {
      return cached.data;
    }

    try {
      const response = await fetch(`${DEFILLAMA_BASE_URL}${endpoint}`);
      
      if (!response.ok) {
        throw new Error(`DefiLlama API error: ${response.status}`);
      }
      
      const data = await response.json();
      this.cache.set(cacheKey, { data, timestamp: Date.now() });
      
      return data;
    } catch (error) {
      console.error('DefiLlama API error:', error);
      
      // Return cached data if available
      if (cached) {
        console.warn('Using stale DefiLlama data due to API error');
        return cached.data;
      }
      
      throw error;
    }
  }

  async getYieldPools() {
    return this.request('/pools');
  }

  async getProtocolTVL(protocol: string) {
    return this.request(`/tvl/${protocol}`);
  }

  async getChainTVL() {
    return this.request('/v2/chains');
  }

  async getPoolsByChain(chain: string) {
    const pools = await this.request('/pools');
    return (pools as any[]).filter(pool => 
      pool.chain?.toLowerCase() === chain.toLowerCase()
    );
  }
}

// Token ID mappings for CoinGecko
export const TOKEN_IDS = {
  'ethereum': 'ethereum',
  'bitcoin': 'bitcoin',
  'usdc': 'usd-coin',
  'usdt': 'tether',
  'uni': 'uniswap',
  'aave': 'aave',
  'comp': 'compound-governance-token',
  'link': 'chainlink',
  'wbtc': 'wrapped-bitcoin',
  'dai': 'dai',
  'matic': 'matic-network',
  'arb': 'arbitrum',
  'op': 'optimism',
} as const;

// Chain mappings for DefiLlama
export const CHAIN_MAPPINGS = {
  'ethereum': 'Ethereum',
  'polygon': 'Polygon',
  'arbitrum': 'Arbitrum',
  'optimism': 'Optimism',
  'bsc': 'BSC',
  'avalanche': 'Avax',
} as const;

// Singleton instances
export const coinGeckoClient = new CoinGeckoClient();
export const defiLlamaClient = new DefiLlamaClient();