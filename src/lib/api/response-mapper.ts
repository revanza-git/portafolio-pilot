// Response mapping utilities for transforming API responses to frontend types
import { TokenBalance } from '@/stores/portfolio';

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
    logoUrl: balance.token.logo_uri || balance.token.logoURI
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