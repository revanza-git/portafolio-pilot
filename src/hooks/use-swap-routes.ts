import { useState, useCallback } from 'react';
import { useAccount } from 'wagmi';
import { Address } from 'viem';
import { useToast } from '@/hooks/use-toast';
import { SwapClient, SwapRoute, SwapQuoteRequest } from '@/lib/api/swap-client';

export function useSwapRoutes() {
  const [routes, setRoutes] = useState<SwapRoute[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const { address, chainId } = useAccount();
  const { toast } = useToast();
  
  const swapClient = new SwapClient();

  const fetchQuotes = useCallback(async (request: Omit<SwapQuoteRequest, 'userAddress' | 'chainId'>) => {
    if (!address || !chainId) {
      toast({
        title: "Wallet Not Connected",
        description: "Please connect your wallet to get swap quotes",
        variant: "destructive",
      });
      return;
    }

    setIsLoading(true);
    setError(null);

    try {
      const swapRoutes = await swapClient.getQuote({
        ...request,
        chainId,
        userAddress: address,
      });
      
      setRoutes(swapRoutes);
      
      if (swapRoutes.length === 0) {
        toast({
          title: "No Routes Found",
          description: "No swap routes available for this token pair",
          variant: "destructive",
        });
      }
    } catch (err: any) {
      const errorMessage = err.message || 'Failed to fetch swap quotes';
      setError(errorMessage);
      toast({
        title: "Failed to Fetch Quotes",
        description: errorMessage,
        variant: "destructive",
      });
    } finally {
      setIsLoading(false);
    }
  }, [address, chainId, swapClient, toast]);

  const executeSwap = useCallback(async (route: SwapRoute) => {
    if (!address) {
      toast({
        title: "Wallet Not Connected",
        description: "Please connect your wallet to execute swap",
        variant: "destructive",
      });
      return;
    }

    try {
      const result = await swapClient.executeSwap(route.id, address);
      
      toast({
        title: "Swap Executed",
        description: `Transaction submitted: ${result.txHash}`,
      });

      return result;
    } catch (err: any) {
      const errorMessage = err.message || 'Failed to execute swap';
      toast({
        title: "Swap Failed",
        description: errorMessage,
        variant: "destructive",
      });
      throw err;
    }
  }, [address, swapClient, toast]);

  const clearRoutes = useCallback(() => {
    setRoutes([]);
    setError(null);
  }, []);

  return {
    routes,
    isLoading,
    error,
    fetchQuotes,
    executeSwap,
    clearRoutes,
  };
}