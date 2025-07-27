import { useState, useEffect } from 'react';
import { Navbar } from '@/components/navigation/navbar';
import { PortfolioOverview } from '@/components/dashboard/portfolio-overview';
import { TokenTable } from '@/components/dashboard/token-table';
import { PortfolioChart } from '@/components/dashboard/portfolio-chart';
import { RecentTransactions } from '@/components/dashboard/recent-transactions';
import { useWalletStore } from '@/stores/wallet';
import { usePortfolioStore } from '@/stores/portfolio';
import { generateMockTokens, generateMockTransactions } from '@/lib/mock-data';
import { Navigate } from 'react-router-dom';

export default function Dashboard() {
  const { isConnected } = useWalletStore();
  const { 
    tokens, 
    totalValue, 
    change24h, 
    setTokens, 
    setTotalValue, 
    setTransactions,
    isLoading,
    setLoading 
  } = usePortfolioStore();

  useEffect(() => {
    if (isConnected) {
      // TODO: Replace with real API calls
      setLoading(true);
      
      setTimeout(() => {
        const mockTokens = generateMockTokens();
        const mockTransactions = generateMockTransactions();
        
        setTokens(mockTokens);
        setTransactions(mockTransactions);
        
        const total = mockTokens.reduce((sum, token) => sum + token.usdValue, 0);
        setTotalValue(total, 5.2);
        
        setLoading(false);
      }, 1000);
    }
  }, [isConnected, setTokens, setTotalValue, setTransactions, setLoading]);

  if (!isConnected) {
    return <Navigate to="/" replace />;
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-background to-muted/20">
      <Navbar />
      
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {/* Portfolio Overview */}
        <PortfolioOverview 
          totalValue={totalValue}
          change24h={change24h}
          isLoading={isLoading}
        />

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-8 mt-8">
          {/* Portfolio Chart */}
          <div className="lg:col-span-2">
            <PortfolioChart />
          </div>

          {/* Recent Transactions */}
          <div>
            <RecentTransactions />
          </div>
        </div>

        {/* Token Holdings */}
        <div className="mt-8">
          <TokenTable tokens={tokens} isLoading={isLoading} />
        </div>
      </div>
    </div>
  );
}