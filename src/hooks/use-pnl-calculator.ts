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
    (tokens || []).map(t => t.symbol.toLowerCase())
  , [tokens]);
  
  const { data: priceData } = useTokenPrices(tokenSymbols);

  useEffect(() => {
    const calculatePnL = async () => {
      if (!priceData || !transactions || transactions.length === 0) return;

      // Validate date range
      if (!options.dateRange.start || !options.dateRange.end || 
          options.dateRange.start >= options.dateRange.end) {
        console.warn('Invalid date range for P&L calculation');
        return;
      }

      setIsCalculating(true);

      try {
        const calculator = new PnLCalculator(options.method);
        
        // Filter transactions by date range with additional validation
        const filteredTransactions = transactions.filter(tx => {
          if (!tx || !tx.timestamp) return false;
          
          try {
            const txDate = new Date(tx.timestamp);
            return !isNaN(txDate.getTime()) && 
                   txDate >= options.dateRange.start && 
                   txDate <= options.dateRange.end;
          } catch (error) {
            console.warn('Invalid transaction timestamp:', tx.timestamp);
            return false;
          }
        });

        // Sort transactions by timestamp with error handling
        filteredTransactions.sort((a, b) => {
          try {
            return a.timestamp - b.timestamp;
          } catch (error) {
            console.warn('Error sorting transactions:', error);
            return 0;
          }
        });

        // Build current price map with validation
        const currentPrices: Record<string, number> = {};
        for (const [symbol, data] of Object.entries(priceData || {})) {
          if (data && typeof data.price === 'number' && !isNaN(data.price)) {
            currentPrices[symbol] = data.price;
          } else {
            console.warn(`Invalid price data for symbol ${symbol}:`, data);
          }
        }

        // Process each transaction with error handling
        for (const tx of filteredTransactions) {
          try {
            calculator.processTransaction(tx, currentPrices);
          } catch (error) {
            console.warn('Error processing transaction:', tx, error);
            // Continue processing other transactions
          }
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