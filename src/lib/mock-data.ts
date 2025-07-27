import { TokenBalance, Transaction, Allowance } from '@/stores/portfolio';

// Mock token data
export function generateMockTokens(): TokenBalance[] {
  return [
    {
      address: '0xA0b86991c431E4dFe7bbb8E5f2D5E8b8A8A8b3c8B',
      symbol: 'USDC',
      name: 'USD Coin',
      decimals: 6,
      balance: '5000000000',
      balanceFormatted: '5,000.00',
      usdValue: 5000,
      priceUsd: 1.00,
      change24h: 0.01,
      logoUrl: 'https://tokens.1inch.io/0xa0b86991c431e4dfe7bb2d8e8b8a8b3c8b8a8b3c8b.png',
    },
    {
      address: '0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2',
      symbol: 'WETH',
      name: 'Wrapped Ether',
      decimals: 18,
      balance: '1500000000000000000',
      balanceFormatted: '1.50',
      usdValue: 3750,
      priceUsd: 2500,
      change24h: 2.5,
      logoUrl: 'https://tokens.1inch.io/0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2.png',
    },
    {
      address: '0x2260FAC5E5542a773Aa44fBCfeDf7C193bc2C599',
      symbol: 'WBTC',
      name: 'Wrapped Bitcoin',
      decimals: 8,
      balance: '5000000',
      balanceFormatted: '0.05',
      usdValue: 2150,
      priceUsd: 43000,
      change24h: -1.2,
      logoUrl: 'https://tokens.1inch.io/0x2260fac5e5542a773aa44fbcfedf7c193bc2c599.png',
    },
    {
      address: '0x1f9840a85d5aF5bf1D1762F925BDADdC4201F984',
      symbol: 'UNI',
      name: 'Uniswap',
      decimals: 18,
      balance: '500000000000000000000',
      balanceFormatted: '500.00',
      usdValue: 3500,
      priceUsd: 7.00,
      change24h: 5.8,
      logoUrl: 'https://tokens.1inch.io/0x1f9840a85d5af5bf1d1762f925bdaddc4201f984.png',
    },
  ];
}

// Mock transaction data
export function generateMockTransactions(): Transaction[] {
  const now = Date.now();
  
  return [
    {
      hash: '0x1234567890abcdef1234567890abcdef12345678',
      type: 'swap',
      timestamp: now - 3600000, // 1 hour ago
      tokenIn: {
        symbol: 'USDC',
        amount: '1,000.00',
        usdValue: 1000,
      },
      tokenOut: {
        symbol: 'WETH',
        amount: '0.40',
        usdValue: 1000,
      },
      status: 'success',
      gasUsed: '150,000',
      gasFee: '$12.50',
    },
    {
      hash: '0xabcdef1234567890abcdef1234567890abcdef12',
      type: 'receive',
      timestamp: now - 7200000, // 2 hours ago
      tokenIn: {
        symbol: 'UNI',
        amount: '50.00',
        usdValue: 350,
      },
      status: 'success',
      gasUsed: '21,000',
      gasFee: '$2.10',
    },
    {
      hash: '0x567890abcdef1234567890abcdef1234567890ab',
      type: 'approve',
      timestamp: now - 14400000, // 4 hours ago
      tokenOut: {
        symbol: 'WBTC',
        amount: 'Unlimited',
        usdValue: 0,
      },
      status: 'success',
      gasUsed: '46,000',
      gasFee: '$4.60',
    },
    {
      hash: '0x90abcdef1234567890abcdef1234567890abcdef',
      type: 'send',
      timestamp: now - 86400000, // 1 day ago
      tokenOut: {
        symbol: 'USDC',
        amount: '500.00',
        usdValue: 500,
      },
      status: 'success',
      gasUsed: '21,000',
      gasFee: '$8.20',
    },
  ];
}

// Mock allowance data
export function generateMockAllowances(): Allowance[] {
  return [
    {
      id: '1',
      token: {
        address: '0xA0b86991c431E4dFe7bbb8E5f2D5E8b8A8A8b3c8B',
        symbol: 'USDC',
        name: 'USD Coin',
        logoUrl: 'https://tokens.1inch.io/0xa0b86991c431e4dfe7bb2d8e8b8a8b3c8b8a8b3c8b.png',
      },
      spender: {
        address: '0x1111111254EEB25477B68fb85Ed929f73A960582',
        name: '1inch Router',
        logoUrl: 'https://1inch.io/img/favicon/favicon-32x32.png',
      },
      amount: '115792089237316195423570985008687907853269984665640564039457584007913129639935',
      amountFormatted: 'Unlimited',
      isUnlimited: true,
      lastUpdated: Date.now() - 86400000,
    },
    {
      id: '2',
      token: {
        address: '0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2',
        symbol: 'WETH',
        name: 'Wrapped Ether',
        logoUrl: 'https://tokens.1inch.io/0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2.png',
      },
      spender: {
        address: '0x7a250d5630B4cF539739dF2C5dAcb4c659F2488D',
        name: 'Uniswap V2 Router',
        logoUrl: 'https://uniswap.org/favicon.ico',
      },
      amount: '5000000000000000000',
      amountFormatted: '5.00',
      isUnlimited: false,
      lastUpdated: Date.now() - 172800000,
    },
    {
      id: '3',
      token: {
        address: '0x1f9840a85d5aF5bf1D1762F925BDADdC4201F984',
        symbol: 'UNI',
        name: 'Uniswap',
        logoUrl: 'https://tokens.1inch.io/0x1f9840a85d5af5bf1d1762f925bdaddc4201f984.png',
      },
      spender: {
        address: '0xE592427A0AEce92De3Edee1F18E0157C05861564',
        name: 'Uniswap V3 Router',
        logoUrl: 'https://uniswap.org/favicon.ico',
      },
      amount: '115792089237316195423570985008687907853269984665640564039457584007913129639935',
      amountFormatted: 'Unlimited',
      isUnlimited: true,
      lastUpdated: Date.now() - 259200000,
    },
  ];
}

// Mock price history data for charts
export function generateMockPriceHistory(days: number = 7) {
  const data = [];
  const now = new Date();
  
  for (let i = days; i >= 0; i--) {
    const date = new Date(now);
    date.setDate(date.getDate() - i);
    
    // Generate realistic price movement
    const baseValue = 14400; // $14.4k base portfolio value
    const variation = (Math.random() - 0.5) * 0.1; // ±10% variation
    const value = baseValue * (1 + variation);
    
    data.push({
      timestamp: date.getTime(),
      value: Math.round(value),
      date: date.toISOString().split('T')[0],
    });
  }
  
  return data;
}