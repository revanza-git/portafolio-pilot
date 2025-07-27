import { useState, useEffect, useMemo } from 'react';
import { usePortfolioStore } from '@/stores/portfolio';
import { useTokenPrices } from '@/hooks/use-market-data';
import { PnLCalculator, AccountingMethod, PnLCalculation } from '@/lib/pnl-calculator';

export interface UsePnLCalculatorOptions {
  method: AccountingMethod;
  dateRange: {
    start: Date;
    end: Date;
  };
}

export function usePnLCalculator(options: UsePnLCalculatorOptions) {
  const { transactions, tokens } = usePortfolioStore();
  const [isCalculating, setIsCalculating] = useState(false);
  const [calculation, setCalculation] = useState<PnLCalculation | null>(null);

  // Get current prices for all tokens
  const tokenSymbols = useMemo(() => 
    tokens.map(t => t.symbol.toLowerCase())
  , [tokens]);
  
  const { data: priceData } = useTokenPrices(tokenSymbols);

  useEffect(() => {
    const calculatePnL = async () => {
      if (!priceData || transactions.length === 0) return;

      setIsCalculating(true);

      try {
        const calculator = new PnLCalculator(options.method);
        
        // Filter transactions by date range
        const filteredTransactions = transactions.filter(tx => {
          const txDate = new Date(tx.timestamp);
          return txDate >= options.dateRange.start && txDate <= options.dateRange.end;
        });

        // Sort transactions by timestamp
        filteredTransactions.sort((a, b) => a.timestamp - b.timestamp);

        // Build current price map
        const currentPrices: Record<string, number> = {};
        for (const [symbol, data] of Object.entries(priceData)) {
          currentPrices[symbol] = data.price;
        }

        // Process each transaction
        for (const tx of filteredTransactions) {
          calculator.processTransaction(tx, currentPrices);
        }

        // Get final calculation
        const result = calculator.getCalculation(currentPrices);
        setCalculation(result);
      } catch (error) {
        console.error('Error calculating P&L:', error);
        setCalculation(null);
      } finally {
        setIsCalculating(false);
      }
    };

    calculatePnL();
  }, [transactions, priceData, options.method, options.dateRange]);

  return {
    calculation,
    isCalculating,
    refetch: () => {
      // Trigger recalculation by updating state
      setCalculation(null);
    }
  };
}