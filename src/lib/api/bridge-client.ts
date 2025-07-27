import { Address } from 'viem';

export interface BridgeRoute {
  id: string;
  fromChain: number;
  toChain: number;
  fromToken: Address;
  toToken: Address;
  fromAmount: string;
  toAmount: string;
  estimatedGas: string;
  estimatedTime: number; // seconds
  fees: {
    bridgeFee: string;
    gasFee: string;
    total: string;
  };
  steps: BridgeStep[];
  provider: 'lifi' | 'socket';
}

export interface BridgeStep {
  type: 'swap' | 'bridge';
  protocol: string;
  fromChain: number;
  toChain: number;
  fromToken: Address;
  toToken: Address;
  fromAmount: string;
  toAmount: string;
  data: string;
  value: string;
  gasLimit: string;
}

export interface BridgeQuoteRequest {
  fromChain: number;
  toChain: number;
  fromToken: Address;
  toToken: Address;
  fromAmount: string;
  userAddress: Address;
  slippage?: number; // percentage, default 0.5
}

export class BridgeClient {
  private baseUrl: string;

  constructor(baseUrl: string = '/api') {
    this.baseUrl = baseUrl;
  }

  async getRoutes(request: BridgeQuoteRequest): Promise<BridgeRoute[]> {
    try {
      const response = await fetch(`${this.baseUrl}/bridge/routes`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(request),
      });

      if (!response.ok) {
        throw new Error(`Bridge API error: ${response.statusText}`);
      }

      return await response.json();
    } catch (error) {
      console.error('Failed to fetch bridge routes:', error);
      
      // Mock data for development
      return this.getMockRoutes(request);
    }
  }

  async executeRoute(routeId: string, userAddress: Address): Promise<{ txHash: string }> {
    try {
      const response = await fetch(`${this.baseUrl}/bridge/execute`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ routeId, userAddress }),
      });

      if (!response.ok) {
        throw new Error(`Bridge execution error: ${response.statusText}`);
      }

      return await response.json();
    } catch (error) {
      console.error('Failed to execute bridge route:', error);
      throw error;
    }
  }

  private getMockRoutes(request: BridgeQuoteRequest): BridgeRoute[] {
    return [
      {
        id: 'lifi-route-1',
        fromChain: request.fromChain,
        toChain: request.toChain,
        fromToken: request.fromToken,
        toToken: request.toToken,
        fromAmount: request.fromAmount,
        toAmount: (BigInt(request.fromAmount) * BigInt(98) / BigInt(100)).toString(),
        estimatedGas: '150000',
        estimatedTime: 300,
        fees: {
          bridgeFee: '0.005',
          gasFee: '0.002',
          total: '0.007',
        },
        steps: [
          {
            type: 'bridge',
            protocol: 'Stargate',
            fromChain: request.fromChain,
            toChain: request.toChain,
            fromToken: request.fromToken,
            toToken: request.toToken,
            fromAmount: request.fromAmount,
            toAmount: (BigInt(request.fromAmount) * BigInt(98) / BigInt(100)).toString(),
            data: '0x',
            value: '0',
            gasLimit: '150000',
          },
        ],
        provider: 'lifi',
      },
      {
        id: 'socket-route-1',
        fromChain: request.fromChain,
        toChain: request.toChain,
        fromToken: request.fromToken,
        toToken: request.toToken,
        fromAmount: request.fromAmount,
        toAmount: (BigInt(request.fromAmount) * BigInt(97) / BigInt(100)).toString(),
        estimatedGas: '180000',
        estimatedTime: 450,
        fees: {
          bridgeFee: '0.008',
          gasFee: '0.003',
          total: '0.011',
        },
        steps: [
          {
            type: 'bridge',
            protocol: 'Hop Protocol',
            fromChain: request.fromChain,
            toChain: request.toChain,
            fromToken: request.fromToken,
            toToken: request.toToken,
            fromAmount: request.fromAmount,
            toAmount: (BigInt(request.fromAmount) * BigInt(97) / BigInt(100)).toString(),
            data: '0x',
            value: '0',
            gasLimit: '180000',
          },
        ],
        provider: 'socket',
      },
    ];
  }
}