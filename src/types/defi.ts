// Shared types for DeFi Portfolio

export interface Chain {
  id: number;
  name: string;
  nativeCurrency: {
    name: string;
    symbol: string;
    decimals: number;
  };
  rpcUrls: string[];
  blockExplorerUrls: string[];
}

export interface Token {
  address: string;
  symbol: string;
  name: string;
  decimals: number;
  logoUrl?: string;
  isNative?: boolean;
}

export interface PriceData {
  price: number;
  change24h: number;
  marketCap?: number;
  volume24h?: number;
}

export interface PortfolioOverview {
  totalValue: number;
  change24h: number;
  changePercent24h: number;
  tokenCount: number;
  chain: string;
}

export interface SwapQuote {
  fromToken: Token;
  toToken: Token;
  fromAmount: string;
  toAmount: string;
  priceImpact: number;
  gasEstimate: string;
  route: string[];
  provider: '0x' | '1inch' | 'uniswap';
}

// API Response types
export interface ApiResponse<T> {
  success: boolean;
  data?: T;
  error?: string;
  timestamp: number;
}

export interface PaginatedResponse<T> {
  data: T[];
  pagination: {
    page: number;
    limit: number;
    total: number;
    hasNext: boolean;
  };
}

// Mock data flags
export const MOCK_DATA_ENABLED = true; // TODO: Set to false when backend is ready