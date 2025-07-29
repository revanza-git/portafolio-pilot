"use client";

import { useState, useEffect } from 'react';
import { Wallet, ChevronDown, Copy, ExternalLink, LogOut } from 'lucide-react';
import { useAccount, useConnect, useDisconnect } from 'wagmi';
import { Button } from '@/components/ui/button';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import { useToast } from '@/hooks/use-toast';
import { useWalletStore } from '@/stores/wallet';
import { useAuth } from '@/contexts/auth-context';

export function WalletConnectButton() {
  const { address, isConnected, chainId } = useAccount();
  const { connect, connectors } = useConnect();
  const { disconnect } = useDisconnect();
  const { toast } = useToast();
  const { setWallet, disconnect: disconnectStore } = useWalletStore();
  const { signIn, signOut, isAuthenticated, isLoading: authLoading } = useAuth();

  // Sync wagmi state with zustand store and trigger authentication
  useEffect(() => {
    console.log('Wallet state change:', { address, isConnected, chainId, isAuthenticated });
    if (isConnected && address && chainId) {
      setWallet(address, chainId);
      
      // Auto-trigger SIWE authentication if wallet connected but not authenticated
      if (!isAuthenticated && !authLoading) {
        console.log('Wallet connected but not authenticated, triggering SIWE...');
        signIn().catch(console.error);
      }
    } else if (!isConnected) {
      disconnectStore();
    }
  }, [isConnected, address, chainId, isAuthenticated, authLoading, setWallet, disconnectStore, signIn]);

  const handleConnect = async () => {
    try {
      const connector = connectors.find(c => c.name === 'MetaMask') || connectors[0];
      if (connector) {
        connect({ connector });
      }
    } catch (error) {
      console.error('Wallet connection error:', error);
      toast({
        title: "Connection Failed",
        description: "Failed to connect wallet. Please try again.",
        variant: "destructive",
      });
    }
  };

  const handleDisconnect = () => {
    disconnect();
    disconnectStore();
    signOut(); // Clear authentication state
    toast({
      title: "Wallet Disconnected",
      description: "Your wallet has been disconnected",
    });
  };

  const copyAddress = () => {
    if (address) {
      navigator.clipboard.writeText(address);
      toast({
        title: "Address Copied",
        description: "Wallet address copied to clipboard",
      });
    }
  };

  const formatAddress = (addr: string) => {
    return `${addr.slice(0, 6)}...${addr.slice(-4)}`;
  };

  if (!isConnected || !address) {
    return (
      <Button 
        onClick={handleConnect}
        className="bg-gradient-primary hover:opacity-90 transition-opacity"
        size="lg"
      >
        <Wallet className="mr-2 h-4 w-4" />
        Connect Wallet
      </Button>
    );
  }

  // Show authentication loading state
  if (isConnected && !isAuthenticated && authLoading) {
    return (
      <Button variant="outline" className="gap-2" disabled>
        <div className="w-2 h-2 bg-yellow-500 rounded-full animate-pulse" />
        Authenticating...
      </Button>
    );
  }

  // Show unauthenticated state (shouldn't happen with auto-auth but just in case)
  if (isConnected && !isAuthenticated) {
    return (
      <Button 
        onClick={() => signIn().catch(console.error)}
        variant="outline" 
        className="gap-2"
      >
        <div className="w-2 h-2 bg-yellow-500 rounded-full" />
        Sign Message
      </Button>
    );
  }

  // Fully authenticated state
  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button variant="outline" className="gap-2">
          <div className="w-2 h-2 bg-green-500 rounded-full" />
          {formatAddress(address)}
          <ChevronDown className="h-4 w-4 opacity-50" />
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="end" className="w-56">
        <DropdownMenuItem onClick={copyAddress}>
          <Copy className="mr-2 h-4 w-4" />
          Copy Address
        </DropdownMenuItem>
        <DropdownMenuItem>
          <ExternalLink className="mr-2 h-4 w-4" />
          View on Etherscan
        </DropdownMenuItem>
        <DropdownMenuSeparator />
        <DropdownMenuItem onClick={handleDisconnect} className="text-destructive">
          <LogOut className="mr-2 h-4 w-4" />
          Disconnect
        </DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>
  );
}