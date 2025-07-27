import { create } from 'zustand';
import { persist } from 'zustand/middleware';

interface WalletState {
  address: string | null;
  chainId: number | null;
  isConnected: boolean;
  balance: string | null;
  
  // Actions
  setWallet: (address: string, chainId: number) => void;
  setBalance: (balance: string) => void;
  disconnect: () => void;
}

export const useWalletStore = create<WalletState>()(
  persist(
    (set) => ({
      address: null,
      chainId: null,
      isConnected: false,
      balance: null,
      
      setWallet: (address, chainId) =>
        set({ address, chainId, isConnected: true }),
      
      setBalance: (balance) => set({ balance }),
      
      disconnect: () =>
        set({ 
          address: null, 
          chainId: null, 
          isConnected: false, 
          balance: null 
        }),
    }),
    {
      name: 'wallet-storage',
      partialize: (state) => ({ 
        address: state.address,
        chainId: state.chainId,
        isConnected: state.isConnected,
      }),
    }
  )
);