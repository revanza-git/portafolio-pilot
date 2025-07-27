import { Transaction } from '@/stores/portfolio';

export interface PnLLot {
  id: string;
  symbol: string;
  quantity: number;
  costBasis: number; // Price per unit when acquired
  timestamp: number;
  transactionHash: string;
}

export interface PnLCalculation {
  realizedPnL: number;
  unrealizedPnL: number;
  totalReturn: number;
  totalReturnPercent: number;
  lots: PnLLot[];
  trades: Array<{
    hash: string;
    symbol: string;
    type: 'buy' | 'sell';
    quantity: number;
    price: number;
    realizedPnL?: number;
    timestamp: number;
  }>;
}

export type AccountingMethod = 'fifo' | 'lifo';

export class PnLCalculator {
  private lots: Map<string, PnLLot[]> = new Map();
  private trades: PnLCalculation['trades'] = [];
  private totalRealizedPnL = 0;

  constructor(private method: AccountingMethod = 'fifo') {}

  // Process a transaction and update lots
  processTransaction(tx: Transaction, currentPrices: Record<string, number>) {
    if (tx.type === 'swap') {
      // Handle swap as sell tokenIn + buy tokenOut
      if (tx.tokenIn) {
        this.processSell(
          tx.tokenIn.symbol,
          parseFloat(tx.tokenIn.amount),
          tx.tokenIn.usdValue / parseFloat(tx.tokenIn.amount),
          tx.timestamp,
          tx.hash
        );
      }
      if (tx.tokenOut) {
        this.processBuy(
          tx.tokenOut.symbol,
          parseFloat(tx.tokenOut.amount),
          tx.tokenOut.usdValue / parseFloat(tx.tokenOut.amount),
          tx.timestamp,
          tx.hash
        );
      }
    } else if (tx.type === 'receive' && tx.tokenOut) {
      // Handle receive as buy
      this.processBuy(
        tx.tokenOut.symbol,
        parseFloat(tx.tokenOut.amount),
        tx.tokenOut.usdValue / parseFloat(tx.tokenOut.amount),
        tx.timestamp,
        tx.hash
      );
    } else if (tx.type === 'send' && tx.tokenIn) {
      // Handle send as sell
      this.processSell(
        tx.tokenIn.symbol,
        parseFloat(tx.tokenIn.amount),
        tx.tokenIn.usdValue / parseFloat(tx.tokenIn.amount),
        tx.timestamp,
        tx.hash
      );
    }
  }

  private processBuy(symbol: string, quantity: number, price: number, timestamp: number, hash: string) {
    if (!this.lots.has(symbol)) {
      this.lots.set(symbol, []);
    }

    const lot: PnLLot = {
      id: `${hash}-${symbol}-${timestamp}`,
      symbol,
      quantity,
      costBasis: price,
      timestamp,
      transactionHash: hash
    };

    this.lots.get(symbol)!.push(lot);
    
    this.trades.push({
      hash,
      symbol,
      type: 'buy',
      quantity,
      price,
      timestamp
    });
  }

  private processSell(symbol: string, quantity: number, price: number, timestamp: number, hash: string) {
    const lots = this.lots.get(symbol) || [];
    if (lots.length === 0) {
      // No lots to sell from - this shouldn't happen in practice
      return;
    }

    let remainingToSell = quantity;
    let realizedPnL = 0;
    const lotsToUpdate = [...lots];

    // Sort lots based on accounting method
    if (this.method === 'fifo') {
      lotsToUpdate.sort((a, b) => a.timestamp - b.timestamp);
    } else {
      lotsToUpdate.sort((a, b) => b.timestamp - a.timestamp);
    }

    for (let i = 0; i < lotsToUpdate.length && remainingToSell > 0; i++) {
      const lot = lotsToUpdate[i];
      const sellQuantity = Math.min(lot.quantity, remainingToSell);
      
      // Calculate realized P&L for this portion
      const lotPnL = sellQuantity * (price - lot.costBasis);
      realizedPnL += lotPnL;
      
      // Update lot quantity
      lot.quantity -= sellQuantity;
      remainingToSell -= sellQuantity;
      
      // Remove lot if fully sold
      if (lot.quantity <= 0) {
        const lotIndex = lots.findIndex(l => l.id === lot.id);
        if (lotIndex !== -1) {
          lots.splice(lotIndex, 1);
        }
      }
    }

    this.totalRealizedPnL += realizedPnL;
    
    this.trades.push({
      hash,
      symbol,
      type: 'sell',
      quantity,
      price,
      realizedPnL,
      timestamp
    });
  }

  // Calculate unrealized P&L based on current prices
  calculateUnrealizedPnL(currentPrices: Record<string, number>): number {
    let unrealizedPnL = 0;

    for (const [symbol, lots] of this.lots.entries()) {
      const currentPrice = currentPrices[symbol.toLowerCase()] || 0;
      
      for (const lot of lots) {
        const unrealizedGain = lot.quantity * (currentPrice - lot.costBasis);
        unrealizedPnL += unrealizedGain;
      }
    }

    return unrealizedPnL;
  }

  // Get complete P&L calculation
  getCalculation(currentPrices: Record<string, number>): PnLCalculation {
    const unrealizedPnL = this.calculateUnrealizedPnL(currentPrices);
    const totalReturn = this.totalRealizedPnL + unrealizedPnL;
    
    // Calculate total cost basis for percentage calculation
    let totalCostBasis = 0;
    for (const [, lots] of this.lots.entries()) {
      for (const lot of lots) {
        totalCostBasis += lot.quantity * lot.costBasis;
      }
    }
    
    const totalReturnPercent = totalCostBasis > 0 ? (totalReturn / totalCostBasis) * 100 : 0;

    const allLots: PnLLot[] = [];
    for (const lots of this.lots.values()) {
      allLots.push(...lots);
    }

    return {
      realizedPnL: this.totalRealizedPnL,
      unrealizedPnL,
      totalReturn,
      totalReturnPercent,
      lots: allLots,
      trades: [...this.trades]
    };
  }

  // Reset calculator
  reset() {
    this.lots.clear();
    this.trades = [];
    this.totalRealizedPnL = 0;
  }
}

// Utility functions
export function formatCurrency(value: number): string {
  return new Intl.NumberFormat('en-US', {
    style: 'currency',
    currency: 'USD',
    minimumFractionDigits: 2,
    maximumFractionDigits: 2
  }).format(value);
}

export function formatPercent(value: number): string {
  return new Intl.NumberFormat('en-US', {
    style: 'percent',
    minimumFractionDigits: 2,
    maximumFractionDigits: 2
  }).format(value / 100);
}