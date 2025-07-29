import { useState, useEffect } from 'react';
import { Navbar } from '@/components/navigation/navbar';
import { PortfolioOverview } from '@/components/dashboard/portfolio-overview';
import { TokenTable } from '@/components/dashboard/token-table';
import { PortfolioChart } from '@/components/dashboard/portfolio-chart';
import { RecentTransactions } from '@/components/dashboard/recent-transactions';
import { useWalletStore } from '@/stores/wallet';
import { usePortfolioStore } from '@/stores/portfolio';
import { useAPIClient } from '@/lib/api/client';
import { Navigate } from 'react-router-dom';
import { useToast } from '@/hooks/use-toast';
import { useAuth } from '@/contexts/auth-context';

export default function Dashboard() {
  console.log('Dashboard: Component starting to render...');
  
  const { isConnected, address } = useWalletStore();
  const { isAuthenticated, isLoading: authLoading } = useAuth();
  console.log('Dashboard: Wallet connected:', isConnected, 'Authenticated:', isAuthenticated);
  
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
  
  const apiClient = useAPIClient();
  const { toast } = useToast();
  
  console.log('Dashboard: Portfolio store loaded, isLoading:', isLoading);

  useEffect(() => {
    if (isConnected && address && isAuthenticated) {
      const fetchPortfolioData = async () => {
        setLoading(true);
        
        try {
          // Debug: Check if auth token exists
          const authToken = localStorage.getItem('auth_token');
          console.log('Dashboard: Auth token exists:', !!authToken);
          console.log('Dashboard: Auth token length:', authToken?.length);
          
          // Fetch balances and transactions in parallel
          const [balancesData, transactionsData] = await Promise.all([
            apiClient.getBalances(address),
            apiClient.getTransactions(address, { limit: 10 })
          ]);
          
          // Transform balance data to match frontend format
          const transformedTokens = balancesData.balances.map(balance => {
            // Convert from wei to proper decimal format
            const decimals = balance.token.decimals || 18;
            const rawBalance = balance.balance || '0';
            const balanceNumber = parseFloat(rawBalance) / Math.pow(10, decimals);
            const balanceFormatted = balanceNumber.toLocaleString('en-US', {
              minimumFractionDigits: 0,
              maximumFractionDigits: 6
            });

            return {
              address: balance.token.address,
              symbol: balance.token.symbol,
              name: balance.token.name,
              decimals: balance.token.decimals,
              balance: balance.balance,
              balanceFormatted: balanceFormatted,
              usdValue: balance.balanceUsd,
              priceUsd: balance.token.price || 0,
              change24h: balance.token.priceChange24h || 0,
              logoUrl: balance.token.logoUri
            };
          });
          
          setTokens(transformedTokens);
          setTransactions(transactionsData.transactions);
          setTotalValue(balancesData.totalValue, 0); // TODO: Calculate 24h change
          
        } catch (error) {
          console.error('Failed to fetch portfolio data:', error);
          toast({
            title: 'Error',
            description: 'Failed to load portfolio data. Please try again.',
            variant: 'destructive'
          });
        } finally {
          setLoading(false);
        }
      };
      
      fetchPortfolioData();
    }
  }, [isConnected, address, isAuthenticated, apiClient, setTokens, setTotalValue, setTransactions, setLoading, toast]);

  // Redirect to home if not connected or not authenticated
  if (!isConnected || !isAuthenticated) {
    console.log('Dashboard: Redirecting - Not connected or authenticated');
    return <Navigate to="/" replace />;
  }

  // Show loading state while auth is being checked
  if (authLoading) {
    console.log('Dashboard: Auth loading...');
    return (
      <div className="min-h-screen bg-gradient-to-br from-background to-muted/20 flex items-center justify-center">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary mx-auto mb-4"></div>
          <p className="text-muted-foreground">Authenticating...</p>
        </div>
      </div>
    );
  }

  console.log('Dashboard: About to render JSX...');
  
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
            <RecentTransactions isLoading={isLoading} />
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