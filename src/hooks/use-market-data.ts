import { useQuery, UseQueryOptions } from '@tanstack/react-query';
import { coinGeckoClient, defiLlamaClient, TOKEN_IDS, CHAIN_MAPPINGS } from '@/lib/api-clients';
import { TokenBalance } from '@/stores/portfolio';

interface PriceHistoryData {
  prices?: [number, number][];
}

// Price data hooks
export function useTokenPrices(tokens: string[]) {
  return useQuery({
    queryKey: ['token-prices', tokens],
    queryFn: async () => {
      try {
        const tokenIds = tokens
          .map(token => TOKEN_IDS[token.toLowerCase() as keyof typeof TOKEN_IDS])
          .filter(Boolean);
        
        if (tokenIds.length === 0) {
          // Silently return fallback data when no valid token IDs (common during loading)
          return getFallbackPrices(tokens);
        }

        const prices = await coinGeckoClient.getTokenPrices(tokenIds);
        
        // Transform to our expected format
        const result: Record<string, { price: number; change24h: number }> = {};
        
        Object.entries(TOKEN_IDS).forEach(([symbol, id]) => {
          if (prices[id]) {
            result[symbol.toUpperCase()] = {
              price: prices[id].usd,
              change24h: prices[id].usd_24h_change || 0,
            };
          }
        });
        
        return result;
      } catch (error) {
        console.error('Failed to fetch token prices:', error);
        // Return fallback mock data
        return getFallbackPrices(tokens);
      }
    },
    enabled: tokens.length > 0, // Only run query when we have tokens
    staleTime: 2 * 60 * 1000, // 2 minutes
    refetchInterval: 5 * 60 * 1000, // 5 minutes
    retry: 3,
    retryDelay: (attemptIndex) => Math.min(1000 * 2 ** attemptIndex, 10000),
  });
}

export function usePriceHistory(tokenSymbol: string, days: number = 7) {
  return useQuery({
    queryKey: ['price-history', tokenSymbol, days],
    queryFn: async () => {
      try {
        const tokenId = TOKEN_IDS[tokenSymbol.toLowerCase() as keyof typeof TOKEN_IDS];
        
        if (!tokenId) {
          throw new Error(`Token ID not found for ${tokenSymbol}`);
        }

        const data = await coinGeckoClient.getPriceHistory(tokenId, days) as any;
        
        // Transform CoinGecko format to our chart format
        if (data && data.prices && Array.isArray(data.prices)) {
          return data.prices.map(([timestamp, price]: [number, number]) => ({
            timestamp,
            value: Math.round(price),
            date: new Date(timestamp).toISOString().split('T')[0],
          }));
        }
        
        throw new Error('Invalid price history data format');
      } catch (error) {
        console.error('Failed to fetch price history:', error);
        // Return fallback mock data
        return getFallbackPriceHistory(days);
      }
    },
    staleTime: 5 * 60 * 1000, // 5 minutes
    refetchInterval: 10 * 60 * 1000, // 10 minutes
    enabled: !!tokenSymbol,
  });
}

export function usePortfolioValueHistory(tokens: TokenBalance[], days: number = 7) {
  return useQuery({
    queryKey: ['portfolio-value-history', tokens, days],
    queryFn: async () => {
      try {
        if (!tokens || tokens.length === 0) {
          return getFallbackPriceHistory(days);
        }

        // Get unique token symbols from portfolio
        const tokenSymbols = tokens
          .map(t => t?.symbol)
          .filter(Boolean)
          .filter((symbol, index, arr) => arr.indexOf(symbol) === index); // Remove duplicates

        // Fetch price history for all tokens in parallel
        const priceHistoryPromises = tokenSymbols.map(async (symbol) => {
          const tokenId = TOKEN_IDS[symbol.toLowerCase() as keyof typeof TOKEN_IDS];
          if (!tokenId) return null;

          try {
            const data = await coinGeckoClient.getPriceHistory(tokenId, days);
            return { symbol, data };
          } catch (error) {
            console.error(`Failed to fetch history for ${symbol}:`, error);
            return null;
          }
        });

        const priceHistories = await Promise.all(priceHistoryPromises);
        const validHistories = priceHistories.filter(Boolean) as Array<{ symbol: string; data: PriceHistoryData }>;

        if (validHistories.length === 0) {
          return getFallbackPortfolioHistory(tokens, days);
        }

        // Create timeline from the first token's data
        const firstHistory = validHistories[0];
        if (!firstHistory.data?.prices || !Array.isArray(firstHistory.data.prices)) {
          return getFallbackPortfolioHistory(tokens, days);
        }

        // Calculate portfolio value at each timestamp
        const portfolioHistory = firstHistory.data.prices.map(([timestamp]: [number, number]) => {
          let portfolioValue = 0;

          // Calculate value for each token at this timestamp
          tokens.forEach(token => {
            if (!token?.symbol || !token?.balanceFormatted) return;

            const balance = parseFloat((token.balanceFormatted || '0').replace(/,/g, ''));
            if (balance > 1000000000) return; // Safety check

            // Find price for this token at this timestamp
            const tokenHistory = validHistories.find(h => h.symbol === token.symbol);
            if (tokenHistory?.data?.prices) {
              // Find closest timestamp in price data
              const priceEntry = tokenHistory.data.prices.find(([ts]: [number, number]) => 
                Math.abs(ts - timestamp) < 24 * 60 * 60 * 1000 // Within 24 hours
              );
              
              if (priceEntry) {
                const [, price] = priceEntry;
                portfolioValue += balance * price;
              }
            }
          });

          return {
            timestamp,
            value: Math.round(portfolioValue),
            date: new Date(timestamp).toISOString().split('T')[0],
          };
        });

        return portfolioHistory;
      } catch (error) {
        console.error('Failed to fetch portfolio value history:', error);
        return getFallbackPortfolioHistory(tokens, days);
      }
    },
    staleTime: 5 * 60 * 1000, // 5 minutes
    refetchInterval: 10 * 60 * 1000, // 10 minutes
    enabled: tokens.length > 0,
  });
}

// Yield/Pool data hooks
export function useYieldPools(chain?: string) {
  return useQuery({
    queryKey: ['yield-pools', chain],
    queryFn: async () => {
      try {
        let pools;
        
        if (chain) {
          const chainName = CHAIN_MAPPINGS[chain.toLowerCase() as keyof typeof CHAIN_MAPPINGS];
          pools = await defiLlamaClient.getPoolsByChain(chainName || chain);
        } else {
          const allPools = await defiLlamaClient.getYieldPools();
          // Ensure pools is an array before slicing
          if (Array.isArray(allPools)) {
            pools = allPools.slice(0, 50); // Limit to top 50 pools
          } else {
            throw new Error('Invalid pools data format');
          }
        }
        
        // Transform to our expected format
        if (!Array.isArray(pools)) {
          throw new Error('Pools data is not an array');
        }
        return pools.map((pool: any) => ({
          id: pool.pool || `${pool.project}-${pool.symbol}`,
          protocol: pool.project || 'Unknown',
          pair: pool.symbol || 'Unknown',
          chain: pool.chain || 'Unknown',
          apr: pool.apy || pool.apyBase || 0,
          tvl: pool.tvlUsd || 0,
          userStaked: 0, // This would come from user's wallet data
          rewards: 0, // This would come from user's wallet data
        }));
      } catch (error) {
        console.error('Failed to fetch yield pools:', error);
        // Return fallback mock data
        return getFallbackYieldPools();
      }
    },
    staleTime: 10 * 60 * 1000, // 10 minutes
    refetchInterval: 15 * 60 * 1000, // 15 minutes
  });
}

export function useProtocolTVL(protocol: string) {
  return useQuery({
    queryKey: ['protocol-tvl', protocol],
    queryFn: async () => {
      try {
        return await defiLlamaClient.getProtocolTVL(protocol);
      } catch (error) {
        console.error('Failed to fetch protocol TVL:', error);
        return { tvl: 0, change_1d: 0 };
      }
    },
    staleTime: 15 * 60 * 1000, // 15 minutes
    enabled: !!protocol,
  });
}

// Portfolio data with real prices
export function usePortfolioWithRealPrices(tokens: TokenBalance[]) {
  const tokenSymbols = (tokens || []).map(t => t?.symbol).filter(Boolean);
  const { data: prices, isLoading: pricesLoading } = useTokenPrices(tokenSymbols);
  
  return useQuery({
    queryKey: ['portfolio-real-prices', tokens, prices],
    queryFn: async () => {
      // Handle empty tokens case
      if (!tokens || tokens.length === 0) {
        return { tokens: [], totalValue: 0, change24h: 0 };
      }
      
      if (!prices) return { tokens, totalValue: 0, change24h: 0 };
      
      const updatedTokens = tokens.map(token => {
        // Safe property access
        if (!token?.symbol || !token?.balanceFormatted) return token;
        
        const priceData = prices[token.symbol];
        if (priceData) {
          const balanceNum = parseFloat((token.balanceFormatted || '0').replace(/,/g, ''));
          
          // Safety check to prevent astronomical values
          if (balanceNum > 1000000000) { // More than 1 billion tokens seems unrealistic
            console.warn(`Suspicious balance for ${token.symbol}: ${balanceNum}. Using 0 instead.`);
            return {
              ...token,
              priceUsd: priceData.price,
              change24h: priceData.change24h,
              usdValue: 0,
            };
          }
          
          const newUsdValue = balanceNum * priceData.price;
          return {
            ...token,
            priceUsd: priceData.price,
            change24h: priceData.change24h,
            usdValue: newUsdValue,
          };
        }
        return token;
      });
      
      const totalValue = updatedTokens.reduce((sum, token) => sum + (token?.usdValue || 0), 0);
      const totalValueYesterday = updatedTokens.reduce((sum, token) => {
        if (!token?.priceUsd || !token?.change24h || !token?.balanceFormatted) return sum;
        
        const yesterdayPrice = token.priceUsd / (1 + token.change24h / 100);
        const balanceNum = parseFloat(token.balanceFormatted.replace(/,/g, ''));
        
        // Safety check for yesterday's calculation too
        if (balanceNum > 1000000000) return sum;
        
        const yesterdayValue = balanceNum * yesterdayPrice;
        return sum + yesterdayValue;
      }, 0);
      
      const change24h = totalValue - totalValueYesterday;
      
      return {
        tokens: updatedTokens,
        totalValue,
        change24h,
      };
    },
    enabled: tokens.length > 0 && (!pricesLoading || !!prices), // Improved enable condition
    staleTime: 2 * 60 * 1000, // 2 minutes
  });
}

// Fallback data functions (in case APIs fail)
function getFallbackPrices(tokens: string[]) {
  // Handle empty tokens array
  if (!tokens || tokens.length === 0) {
    return {};
  }
  
  const fallbackPrices: Record<string, { price: number; change24h: number }> = {
    'ETH': { price: 2458.30, change24h: 2.5 },
    'USDC': { price: 1.00, change24h: 0.01 },
    'WBTC': { price: 43000, change24h: -1.2 },
    'UNI': { price: 7.00, change24h: 5.8 },
    'AAVE': { price: 95.50, change24h: 3.2 },
  };
  
  return Object.fromEntries(
    tokens
      .filter(token => token && typeof token === 'string')
      .map(token => [
        token.toUpperCase(),
        fallbackPrices[token.toUpperCase()] || { price: 1, change24h: 0 }
      ])
  );
}

function getFallbackPriceHistory(days: number) {
  const data = [];
  const now = new Date();
  const baseValue = 14400;
  
  for (let i = days; i >= 0; i--) {
    const date = new Date(now);
    date.setDate(date.getDate() - i);
    const variation = (Math.random() - 0.5) * 0.1;
    const value = baseValue * (1 + variation);
    
    data.push({
      timestamp: date.getTime(),
      value: Math.round(value),
      date: date.toISOString().split('T')[0],
    });
  }
  
  return data;
}

function getFallbackYieldPools() {
  return [
    {
      id: '1',
      protocol: 'Aave',
      pair: 'USDC',
      chain: 'Ethereum',
      apr: 4.2,
      tvl: 2100000000,
      userStaked: 0,
      rewards: 0,
    },
    {
      id: '2',
      protocol: 'Compound',
      pair: 'ETH',
      chain: 'Ethereum',
      apr: 3.8,
      tvl: 1500000000,
      userStaked: 0,
      rewards: 0,
    },
  ];
}

function getFallbackPortfolioHistory(tokens: TokenBalance[], days: number) {
  const data = [];
  const now = new Date();
  
  // Calculate base portfolio value from current tokens
  const baseValue = tokens.reduce((sum, token) => {
    if (!token?.balanceFormatted || !token?.usdValue) return sum;
    return sum + (token.usdValue || 0);
  }, 0) || 14400; // Default to 14400 if no tokens
  
  for (let i = days; i >= 0; i--) {
    const date = new Date(now);
    date.setDate(date.getDate() - i);
    const variation = (Math.random() - 0.5) * 0.1;
    const value = baseValue * (1 + variation);
    
    data.push({
      timestamp: date.getTime(),
      value: Math.round(value),
      date: date.toISOString().split('T')[0],
    });
  }
  
  return data;
}