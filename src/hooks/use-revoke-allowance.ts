import React, { useCallback, useState, useEffect } from 'react';
import { useWriteContract, useWaitForTransactionReceipt, useAccount } from 'wagmi';
import { Address } from 'viem';
import { useToast } from '@/hooks/use-toast';
import { ERC20_ABI } from '@/lib/contracts/abis';

interface RevokeAllowanceParams {
  tokenAddress: Address;
  spenderAddress: Address;
  tokenSymbol: string;
  spenderName: string;
}

export function useRevokeAllowance() {
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

  const revokeAllowance = useCallback(async (params: RevokeAllowanceParams) => {
    if (!userAddress || !chainId) {
      toast({
        title: "Wallet Not Connected",
        description: "Please connect your wallet to revoke allowances",
        variant: "destructive",
      });
      return;
    }

    try {
      setIsSubmitting(true);

      console.log('Revoking allowance:', {
        token: params.tokenAddress,
        spender: params.spenderAddress,
        user: userAddress,
        chainId
      });

      writeContract({
        address: params.tokenAddress,
        abi: ERC20_ABI,
        functionName: 'approve',
        args: [params.spenderAddress, 0n],
      });

      toast({
        title: "Transaction Submitted",
        description: `Revoking ${params.tokenSymbol} allowance for ${params.spenderName}`,
      });

    } catch (err: any) {
      console.error('Revoke allowance error:', err);
      
      let errorMessage = 'Failed to revoke allowance';
      
      if (err.message?.includes('User rejected')) {
        errorMessage = 'Transaction was rejected by user';
      } else if (err.message?.includes('insufficient funds')) {
        errorMessage = 'Insufficient funds for gas fees';
      } else if (err.shortMessage) {
        errorMessage = err.shortMessage;
      }

      toast({
        title: "Transaction Failed",
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
        title: "Allowance Revoked",
        description: "The token allowance has been successfully revoked",
      });
    }
  }, [isConfirmed, toast]);

  // Handle transaction error
  useEffect(() => {
    if (error) {
      setIsSubmitting(false);
      toast({
        title: "Transaction Failed",
        description: error.message || "Failed to revoke allowance",
        variant: "destructive",
      });
    }
  }, [error, toast]);

  return {
    revokeAllowance,
    isLoading: isSubmitting || isConfirming,
    isSuccess: isConfirmed,
    error: error?.message || null,
    txHash: hash,
    reset: () => {
      setIsSubmitting(false);
    },
  };
}