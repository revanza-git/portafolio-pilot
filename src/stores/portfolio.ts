import { create } from 'zustand';

export interface TokenBalance {
  address: string;
  symbol: string;
  name: string;
  decimals: number;
  balance: string;
  balanceFormatted: string;
  usdValue: number;
  priceUsd: number;
  change24h: number;
  logoUrl?: string;
}

export interface Transaction {
  hash: string;
  type: 'send' | 'receive' | 'swap' | 'approve';
  timestamp: number;
  tokenIn?: {
    symbol: string;
    amount: string;
    usdValue: number;
  };
  tokenOut?: {
    symbol: string;
    amount: string;
    usdValue: number;
  };
  status: 'success' | 'pending' | 'failed';
  gasUsed?: string;
  gasFee?: string;
}

export interface Allowance {
  id: string;
  token: {
    address: string;
    symbol: string;
    name: string;
    logoUrl?: string;
  };
  spender: {
    address: string;
    name: string;
    logoUrl?: string;
  };
  amount: string;
  amountFormatted: string;
  isUnlimited: boolean;
  lastUpdated: number;
}

interface PortfolioState {
  // Portfolio data
  tokens: TokenBalance[];
  totalValue: number;
  change24h: number;
  isLoading: boolean;
  
  // Transactions
  transactions: Transaction[];
  transactionsLoading: boolean;
  
  // Allowances
  allowances: Allowance[];
  allowancesLoading: boolean;
  
  // Actions
  setTokens: (tokens: TokenBalance[]) => void;
  setTotalValue: (value: number, change24h: number) => void;
  setLoading: (loading: boolean) => void;
  setTransactions: (transactions: Transaction[]) => void;
  setTransactionsLoading: (loading: boolean) => void;
  setAllowances: (allowances: Allowance[]) => void;
  setAllowancesLoading: (loading: boolean) => void;
  updateTokenBalance: (address: string, balance: string, usdValue: number) => void;
}

export const usePortfolioStore = create<PortfolioState>((set, get) => ({
  // Initial state
  tokens: [],
  totalValue: 0,
  change24h: 0,
  isLoading: false,
  transactions: [],
  transactionsLoading: false,
  allowances: [],
  allowancesLoading: false,
  
  // Actions
  setTokens: (tokens) => set({ tokens }),
  
  setTotalValue: (totalValue, change24h) => {
    // Safety check to prevent astronomical portfolio values
    const safeTotalValue = totalValue > 1000000000000 ? 0 : totalValue; // Cap at $1 trillion
    if (totalValue > 1000000000000) {
      console.warn(`Unrealistic portfolio value detected: $${totalValue}. Using $0 instead.`);
    }
    set({ totalValue: safeTotalValue, change24h });
  },
  
  setLoading: (isLoading) => set({ isLoading }),
  
  setTransactions: (transactions) => set({ transactions }),
  
  setTransactionsLoading: (transactionsLoading) => set({ transactionsLoading }),
  
  setAllowances: (allowances) => set({ allowances }),
  
  setAllowancesLoading: (allowancesLoading) => set({ allowancesLoading }),
  
  updateTokenBalance: (address, balance, usdValue) => {
    // Safety check to prevent astronomical token values
    const safeUsdValue = usdValue > 1000000000000 ? 0 : usdValue; // Cap at $1 trillion per token
    if (usdValue > 1000000000000) {
      console.warn(`Unrealistic token value detected for ${address}: $${usdValue}. Using $0 instead.`);
    }
    
    const { tokens } = get();
    const updatedTokens = tokens.map(token =>
      token.address === address 
        ? { ...token, balance, usdValue: safeUsdValue }
        : token
    );
    set({ tokens: updatedTokens });
  },
}));