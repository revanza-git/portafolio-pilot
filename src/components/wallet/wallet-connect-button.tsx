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

export function WalletConnectButton() {
  const { address, isConnected, chainId } = useAccount();
  const { connect, connectors } = useConnect();
  const { disconnect } = useDisconnect();
  const { toast } = useToast();
  const { setWallet, disconnect: disconnectStore } = useWalletStore();

  // Sync wagmi state with zustand store
  useEffect(() => {
    console.log('Wallet state change:', { address, isConnected, chainId });
    if (isConnected && address && chainId) {
      setWallet(address, chainId);
    } else if (!isConnected) {
      disconnectStore();
    }
  }, [isConnected, address, chainId, setWallet, disconnectStore]);

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

  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button variant="outline" className="gap-2">
          <div className="w-2 h-2 bg-success rounded-full" />
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