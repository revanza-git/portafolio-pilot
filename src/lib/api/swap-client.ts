import { Address } from 'viem';

export interface SwapRoute {
  id: string;
  fromToken: Address;
  toToken: Address;
  fromAmount: string;
  toAmount: string;
  estimatedGas: string;
  gasPrice: string;
  priceImpact: number; // percentage
  fees: {
    protocolFee: string;
    gasFee: string;
    total: string;
  };
  path: Address[];
  provider: '0x' | '1inch' | 'uniswap';
  dex: string;
  calldata: string;
  value: string;
}

export interface SwapQuoteRequest {
  chainId: number;
  fromToken: Address;
  toToken: Address;
  fromAmount: string;
  userAddress: Address;
  slippage?: number; // percentage, default 0.5
  gasPrice?: string;
}

export class SwapClient {
  private baseUrl: string;

  constructor(baseUrl: string = '/api') {
    this.baseUrl = baseUrl;
  }

  async getQuote(request: SwapQuoteRequest): Promise<SwapRoute[]> {
    try {
      const response = await fetch(`${this.baseUrl}/swap/quote`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(request),
      });

      if (!response.ok) {
        throw new Error(`Swap API error: ${response.statusText}`);
      }

      return await response.json();
    } catch (error) {
      console.error('Failed to fetch swap quotes:', error);
      
      // Mock data for development
      return this.getMockQuotes(request);
    }
  }

  async executeSwap(routeId: string, userAddress: Address): Promise<{ txHash: string }> {
    try {
      const response = await fetch(`${this.baseUrl}/swap/execute`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ routeId, userAddress }),
      });

      if (!response.ok) {
        throw new Error(`Swap execution error: ${response.statusText}`);
      }

      return await response.json();
    } catch (error) {
      console.error('Failed to execute swap:', error);
      throw error;
    }
  }

  private getMockQuotes(request: SwapQuoteRequest): SwapRoute[] {
    const baseToAmount = BigInt(request.fromAmount) * BigInt(95) / BigInt(100);
    
    return [
      {
        id: '0x-route-1',
        fromToken: request.fromToken,
        toToken: request.toToken,
        fromAmount: request.fromAmount,
        toAmount: baseToAmount.toString(),
        estimatedGas: '120000',
        gasPrice: request.gasPrice || '20000000000',
        priceImpact: 0.1,
        fees: {
          protocolFee: '0.003',
          gasFee: '0.002',
          total: '0.005',
        },
        path: [request.fromToken, request.toToken],
        provider: '0x',
        dex: 'Uniswap V3',
        calldata: '0x',
        value: '0',
      },
      {
        id: '1inch-route-1',
        fromToken: request.fromToken,
        toToken: request.toToken,
        fromAmount: request.fromAmount,
        toAmount: (baseToAmount * BigInt(101) / BigInt(100)).toString(),
        estimatedGas: '140000',
        gasPrice: request.gasPrice || '20000000000',
        priceImpact: 0.08,
        fees: {
          protocolFee: '0.002',
          gasFee: '0.0025',
          total: '0.0045',
        },
        path: [request.fromToken, request.toToken],
        provider: '1inch',
        dex: 'Multiple DEXs',
        calldata: '0x',
        value: '0',
      },
      {
        id: 'uniswap-route-1',
        fromToken: request.fromToken,
        toToken: request.toToken,
        fromAmount: request.fromAmount,
        toAmount: (baseToAmount * BigInt(99) / BigInt(100)).toString(),
        estimatedGas: '110000',
        gasPrice: request.gasPrice || '20000000000',
        priceImpact: 0.12,
        fees: {
          protocolFee: '0.003',
          gasFee: '0.0018',
          total: '0.0048',
        },
        path: [request.fromToken, request.toToken],
        provider: 'uniswap',
        dex: 'Uniswap V3',
        calldata: '0x',
        value: '0',
      },
    ];
  }
}