// Response mapping utilities for transforming API responses to frontend types
import { TokenBalance } from '@/stores/portfolio';

/**
 * Chain ID to name mapping
 */
export const CHAIN_NAMES: Record<number, string> = {
  1: 'Ethereum',
  137: 'Polygon',
  42161: 'Arbitrum',
  10: 'Optimism',
  56: 'BSC',
  43114: 'Avalanche',
  250: 'Fantom',
  25: 'Cronos',
  100: 'Gnosis'
};

/**
 * Get chain name from chain ID
 */
export function getChainName(chainId: number): string {
  return CHAIN_NAMES[chainId] || `Chain ${chainId}`;
}

/**
 * Get unique chains from token list
 */
export function getUniqueChains(tokens: TokenBalance[]): string[] {
  const chainIds = Array.from(new Set(tokens.map(token => token.chainId || 1)));
  return chainIds.map(id => getChainName(id));
}

/**
 * Find best performing token from token list
 */
export function getBestPerformer(tokens: TokenBalance[]): { symbol: string; change24h: number; priceUsd: number } | null {
  if (!tokens || tokens.length === 0) return null;
  
  const validTokens = tokens.filter(token => 
    token.change24h !== undefined && 
    token.change24h !== null && 
    !isNaN(token.change24h) &&
    token.usdValue > 1 // Only consider tokens with meaningful value
  );
  
  if (validTokens.length === 0) return null;
  
  const bestToken = validTokens.reduce((best, current) => 
    current.change24h > best.change24h ? current : best
  );
  
  return {
    symbol: bestToken.symbol,
    change24h: bestToken.change24h,
    priceUsd: bestToken.priceUsd
  };
}

/**
 * Safely converts a balance string (in wei) to a decimal number
 */
export function parseTokenBalance(balance: string, decimals: number): number {
  try {
    if (!balance || balance === '0') return 0;
    
    const balanceBigInt = BigInt(balance);
    const divisor = BigInt(10) ** BigInt(decimals);
    
    // Convert to number safely - this might lose precision for very large numbers
    // but is sufficient for display purposes
    return Number(balanceBigInt) / Number(divisor);
  } catch (error) {
    console.warn('Error parsing token balance:', { balance, decimals, error });
    return 0;
  }
}

/**
 * Formats a numeric balance for display
 */
export function formatTokenBalance(balance: number, decimals: number): string {
  return balance.toLocaleString('en-US', {
    minimumFractionDigits: 0,
    maximumFractionDigits: Math.min(decimals, 6) // Cap at 6 decimal places
  });
}

/**
 * Maps API balance response to frontend TokenBalance format
 */
export function mapBalanceToTokenBalance(balance: any): TokenBalance | null {
  if (!balance?.token) {
    console.warn('Balance missing token data:', balance);
    return null;
  }

  const decimals = balance.token.decimals || 18;
  const rawBalance = balance.balance || '0';
  const balanceNumber = parseTokenBalance(rawBalance, decimals);
  const balanceFormatted = formatTokenBalance(balanceNumber, decimals);

  return {
    address: balance.token.address,
    symbol: balance.token.symbol,
    name: balance.token.name,
    decimals: balance.token.decimals,
    balance: balance.balance,
    balanceFormatted,
    usdValue: balance.balance_usd || balance.balanceUSD || 0,
    priceUsd: balance.token.price_usd || balance.token.priceUSD || 0,
    change24h: balance.token.price_change_24h || balance.token.priceChange24h || 0,
    logoUrl: balance.token.logo_uri || balance.token.logoURI,
    chainId: balance.token.chain_id || balance.token.chainID || balance.chainId || 1 // Default to Ethereum mainnet
  };
}

/**
 * Maps API transaction response to frontend Transaction format
 */
export function mapTransactionToFrontend(transaction: any) {
  // This would map the API transaction format to the frontend format
  // Implementation depends on the final API response structure
  return {
    hash: transaction.hash,
    type: transaction.type,
    timestamp: new Date(transaction.timestamp).getTime(),
    status: transaction.status,
    // Add more mappings as needed based on API structure
  };
}