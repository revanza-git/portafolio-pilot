import { useState, useEffect } from 'react';
import { Navbar } from '@/components/navigation/navbar';
import { TransactionFilters } from '@/components/transactions/transaction-filters';
import { TransactionList } from '@/components/transactions/transaction-list';
import { useWalletStore } from '@/stores/wallet';
import { usePortfolioStore } from '@/stores/portfolio';
import { generateMockTransactions } from '@/lib/mock-data';
import { Navigate } from 'react-router-dom';

export default function Transactions() {
  const { isConnected } = useWalletStore();
  const { 
    transactions, 
    transactionsLoading, 
    setTransactions, 
    setTransactionsLoading 
  } = usePortfolioStore();
  
  const [filters, setFilters] = useState({
    type: 'all',
    status: 'all',
    timeRange: '7d',
  });

  useEffect(() => {
    if (isConnected) {
      // TODO: Replace with real API calls
      setTransactionsLoading(true);
      
      setTimeout(() => {
        const mockTransactions = generateMockTransactions();
        setTransactions(mockTransactions);
        setTransactionsLoading(false);
      }, 800);
    }
  }, [isConnected, setTransactions, setTransactionsLoading]);

  if (!isConnected) {
    return <Navigate to="/" replace />;
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-background to-muted/20">
      <Navbar />
      
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="mb-8">
          <h1 className="text-3xl font-bold">Transaction History</h1>
          <p className="text-muted-foreground mt-2">
            View and filter your recent DeFi activity
          </p>
        </div>

        <div className="space-y-6">
          <TransactionFilters 
            filters={filters} 
            onFiltersChange={setFilters} 
          />
          
          <TransactionList 
            transactions={transactions}
            isLoading={transactionsLoading}
            filters={filters}
          />
        </div>
      </div>
    </div>
  );
}