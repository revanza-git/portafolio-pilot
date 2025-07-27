import React, { useCallback, useState, useEffect } from 'react';
import { useWriteContract, useWaitForTransactionReceipt, useAccount } from 'wagmi';
import { Address } from 'viem';
import { useToast } from '@/hooks/use-toast';
import { 
  AAVE_V3_POOL_ABI, 
  COMPOUND_V3_COMET_ABI, 
  UNISWAP_V3_POSITION_MANAGER_ABI 
} from '@/lib/contracts/abis';
import { getContractAddress } from '@/lib/contracts/addresses';

interface ClaimRewardsParams {
  protocol: 'aave' | 'compound' | 'uniswap';
  poolId: string;
  tokenAddress?: Address;
  positionId?: string;
}

interface BatchClaimState {
  isLoading: boolean;
  progress: number;
  currentClaim: number;
  totalClaims: number;
  completedClaims: string[];
  failedClaims: { poolId: string; error: string }[];
}

export function useClaimRewards() {
  const [isSubmitting, setIsSubmitting] = useState(false);
  const { toast } = useToast();
  const { address: userAddress, chainId } = useAccount();
  
  const { 
    data: hash,
    error,
    writeContract 
  } = useWriteContract();
  
  const { isLoading: isConfirming, isSuccess: isConfirmed } = useWaitForTransactionReceipt({
    hash,
  });

  const claimRewards = useCallback(async (params: ClaimRewardsParams) => {
    if (!userAddress || !chainId) {
      toast({
        title: "Wallet Not Connected",
        description: "Please connect your wallet to claim rewards",
        variant: "destructive",
      });
      return;
    }

    try {
      setIsSubmitting(true);

      console.log('Claiming rewards:', {
        protocol: params.protocol,
        poolId: params.poolId,
        user: userAddress,
        chainId
      });

      switch (params.protocol) {
        case 'aave':
          try {
            const aavePoolAddress = getContractAddress(chainId as any, 'AAVE_V3_REWARDS_CONTROLLER');
            writeContract({
              address: aavePoolAddress as Address,
              abi: AAVE_V3_POOL_ABI,
              functionName: 'claimRewards',
              args: [
                params.tokenAddress ? [params.tokenAddress] : [],
                userAddress
              ],
            });
          } catch (contractError) {
            console.warn('Aave contract not found, using mock transaction');
            // For demo purposes - show success after delay
            setTimeout(() => {
              setIsSubmitting(false);
              toast({
                title: "Mock Claim Successful",
                description: "This is a demo transaction (Aave contract not available)",
              });
            }, 2000);
            return;
          }
          break;
          
        case 'compound':
          try {
            const compoundAddress = getContractAddress(chainId as any, 'COMPOUND_V3_USDC');
            writeContract({
              address: compoundAddress as Address,
              abi: COMPOUND_V3_COMET_ABI,
              functionName: 'claim',
              args: [userAddress, true],
            });
          } catch (contractError) {
            console.warn('Compound contract not found, using mock transaction');
            setTimeout(() => {
              setIsSubmitting(false);
              toast({
                title: "Mock Claim Successful",
                description: "This is a demo transaction (Compound contract not available)",
              });
            }, 2000);
            return;
          }
          break;
          
        case 'uniswap':
          try {
            const uniswapAddress = getContractAddress(chainId as any, 'UNISWAP_V3_POSITION_MANAGER');
            writeContract({
              address: uniswapAddress as Address,
              abi: UNISWAP_V3_POSITION_MANAGER_ABI,
              functionName: 'collect',
              args: [{
                tokenId: BigInt(params.positionId || '0'),
                recipient: userAddress,
                amount0Max: BigInt('340282366920938463463374607431768211455'),
                amount1Max: BigInt('340282366920938463463374607431768211455'),
              }],
            });
          } catch (contractError) {
            console.warn('Uniswap contract not found, using mock transaction');
            setTimeout(() => {
              setIsSubmitting(false);
              toast({
                title: "Mock Claim Successful",
                description: "This is a demo transaction (Uniswap contract not available)",
              });
            }, 2000);
            return;
          }
          break;
          
        default:
          throw new Error(`Unsupported protocol: ${params.protocol}`);
      }

      toast({
        title: "Transaction Submitted",
        description: `Claiming rewards from ${params.protocol} pool`,
      });

    } catch (err: any) {
      console.error('Claim rewards error:', err);
      
      let errorMessage = 'Failed to claim rewards';
      
      if (err.message?.includes('User rejected')) {
        errorMessage = 'Transaction was rejected by user';
      } else if (err.message?.includes('insufficient funds')) {
        errorMessage = 'Insufficient funds for gas fees';
      } else if (err.message?.includes('No rewards to claim')) {
        errorMessage = 'No rewards available to claim';
      } else if (err.shortMessage) {
        errorMessage = err.shortMessage;
      }

      toast({
        title: "Claim Failed",
        description: errorMessage,
        variant: "destructive",
      });
      
      setIsSubmitting(false);
    }
  }, [userAddress, chainId, writeContract, toast]);

  // Handle transaction confirmation
  useEffect(() => {
    if (isConfirmed) {
      setIsSubmitting(false);
      toast({
        title: "Rewards Claimed",
        description: "Your rewards have been successfully claimed",
      });
    }
  }, [isConfirmed, toast]);

  // Handle transaction error
  useEffect(() => {
    if (error) {
      setIsSubmitting(false);
      toast({
        title: "Claim Failed",
        description: error.message || "Failed to claim rewards",
        variant: "destructive",
      });
    }
  }, [error, toast]);

  return {
    claimRewards,
    isLoading: isSubmitting || isConfirming,
    isSuccess: isConfirmed,
    error: error?.message || null,
    txHash: hash,
    reset: () => {
      setIsSubmitting(false);
    },
  };
}

// Batch claim hook for claiming multiple rewards
export function useBatchClaimRewards() {
  const [batchState, setBatchState] = useState<BatchClaimState>({
    isLoading: false,
    progress: 0,
    currentClaim: 0,
    totalClaims: 0,
    completedClaims: [],
    failedClaims: [],
  });

  const { claimRewards } = useClaimRewards();
  const { toast } = useToast();

  const batchClaimRewards = useCallback(async (claimParams: ClaimRewardsParams[]) => {
    if (claimParams.length === 0) return;

    setBatchState({
      isLoading: true,
      progress: 0,
      currentClaim: 0,
      totalClaims: claimParams.length,
      completedClaims: [],
      failedClaims: [],
    });

    for (let i = 0; i < claimParams.length; i++) {
      const params = claimParams[i];
      
      setBatchState(prev => ({
        ...prev,
        currentClaim: i + 1,
        progress: ((i + 1) / claimParams.length) * 100,
      }));

      try {
        await claimRewards(params);
        
        setBatchState(prev => ({
          ...prev,
          completedClaims: [...prev.completedClaims, params.poolId],
        }));

        // Add delay between transactions to avoid nonce issues
        if (i < claimParams.length - 1) {
          await new Promise(resolve => setTimeout(resolve, 2000));
        }
        
      } catch (error: any) {
        console.error(`Failed to claim from pool ${params.poolId}:`, error);
        
        setBatchState(prev => ({
          ...prev,
          failedClaims: [...prev.failedClaims, {
            poolId: params.poolId,
            error: error.message || 'Unknown error'
          }],
        }));
      }
    }

    setBatchState(prev => ({ ...prev, isLoading: false }));

    // Show summary toast
    const completedCount = batchState.completedClaims.length;
    const failedCount = batchState.failedClaims.length;
    
    if (completedCount > 0) {
      toast({
        title: "Batch Claim Complete",
        description: `Successfully claimed from ${completedCount} pools${failedCount > 0 ? `, ${failedCount} failed` : ''}`,
      });
    }
  }, [claimRewards, toast, batchState.completedClaims.length, batchState.failedClaims.length]);

  return {
    batchClaimRewards,
    ...batchState,
    reset: () => setBatchState({
      isLoading: false,
      progress: 0,
      currentClaim: 0,
      totalClaims: 0,
      completedClaims: [],
      failedClaims: [],
    }),
  };
}