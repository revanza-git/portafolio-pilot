import { useState, useCallback } from 'react';
import { useAccount } from 'wagmi';
import { Address } from 'viem';
import { useToast } from '@/hooks/use-toast';
import { BridgeClient, BridgeRoute, BridgeQuoteRequest } from '@/lib/api/bridge-client';

export function useBridgeRoutes() {
  const [routes, setRoutes] = useState<BridgeRoute[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const { address } = useAccount();
  const { toast } = useToast();
  
  const bridgeClient = new BridgeClient();

  const fetchRoutes = useCallback(async (request: Omit<BridgeQuoteRequest, 'userAddress'>) => {
    if (!address) {
      toast({
        title: "Wallet Not Connected",
        description: "Please connect your wallet to get bridge routes",
        variant: "destructive",
      });
      return;
    }

    setIsLoading(true);
    setError(null);

    try {
      const bridgeRoutes = await bridgeClient.getRoutes({
        ...request,
        userAddress: address,
      });
      
      setRoutes(bridgeRoutes);
      
      if (bridgeRoutes.length === 0) {
        toast({
          title: "No Routes Found",
          description: "No bridge routes available for this token pair",
          variant: "destructive",
        });
      }
    } catch (err: any) {
      const errorMessage = err.message || 'Failed to fetch bridge routes';
      setError(errorMessage);
      toast({
        title: "Failed to Fetch Routes",
        description: errorMessage,
        variant: "destructive",
      });
    } finally {
      setIsLoading(false);
    }
  }, [address, bridgeClient, toast]);

  const executeRoute = useCallback(async (route: BridgeRoute) => {
    if (!address) {
      toast({
        title: "Wallet Not Connected",
        description: "Please connect your wallet to execute bridge",
        variant: "destructive",
      });
      return;
    }

    try {
      const result = await bridgeClient.executeRoute(route.id, address);
      
      toast({
        title: "Bridge Executed",
        description: `Transaction submitted: ${result.txHash}`,
      });

      return result;
    } catch (err: any) {
      const errorMessage = err.message || 'Failed to execute bridge';
      toast({
        title: "Bridge Failed",
        description: errorMessage,
        variant: "destructive",
      });
      throw err;
    }
  }, [address, bridgeClient, toast]);

  const clearRoutes = useCallback(() => {
    setRoutes([]);
    setError(null);
  }, []);

  return {
    routes,
    isLoading,
    error,
    fetchRoutes,
    executeRoute,
    clearRoutes,
  };
}