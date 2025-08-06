import { useState, useEffect } from 'react';
import { Navbar } from '@/components/navigation/navbar';
import { TransactionFilters } from '@/components/transactions/transaction-filters';
import { TransactionList } from '@/components/transactions/transaction-list';
import { useWalletStore } from '@/stores/wallet';
import { usePortfolioStore } from '@/stores/portfolio';
import { useAPIClient } from '@/lib/api/client';
import { mapTransactionToFrontend } from '@/lib/api/response-mapper';
import { Navigate } from 'react-router-dom';
import { useToast } from '@/hooks/use-toast';
import { useAuth } from '@/contexts/auth-context';

export default function Transactions() {
  const { isConnected, address, chainId } = useWalletStore();
  const { isAuthenticated } = useAuth();
  const { 
    transactions, 
    transactionsLoading, 
    setTransactions, 
    setTransactionsLoading 
  } = usePortfolioStore();
  
  const apiClient = useAPIClient();
  const { toast } = useToast();
  
  const [filters, setFilters] = useState({
    type: 'all',
    status: 'all',
    timeRange: '7d',
  });

  useEffect(() => {
    if (isConnected && address && isAuthenticated) {
      const fetchTransactions = async () => {
        setTransactionsLoading(true);
        
        try {
          const params = {
            limit: 50,
            chainId,
            ...(filters.type !== 'all' && { type: filters.type })
          };
          
          const transactionsData = await apiClient.getTransactions(address, params);
          
          // Transform API response to frontend format
          const transformedTransactions = (transactionsData.data || [])
            .map(mapTransactionToFrontend)
            .filter(tx => tx !== null);
          
          setTransactions(transformedTransactions);
        } catch (error) {
          console.error('Failed to fetch transactions:', error);
          toast({
            title: 'Error',
            description: 'Failed to load transaction history. Please try again.',
            variant: 'destructive'
          });
          // Set empty array on error
          setTransactions([]);
        } finally {
          setTransactionsLoading(false);
        }
      };
      
      fetchTransactions();
    }
  }, [isConnected, address, chainId, isAuthenticated, filters.type, apiClient, setTransactions, setTransactionsLoading, toast]);

  if (!isConnected || !isAuthenticated) {
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