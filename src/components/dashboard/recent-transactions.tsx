import { ArrowUpRight, ArrowDownLeft, RefreshCw, CheckCircle } from 'lucide-react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Skeleton } from '@/components/ui/skeleton';
import { usePortfolioStore } from '@/stores/portfolio';
import { Link } from 'react-router-dom';

interface RecentTransactionsProps {
  isLoading?: boolean;
}

export function RecentTransactions({ isLoading }: RecentTransactionsProps = {}) {
  const { transactions } = usePortfolioStore();
  const recentTransactions = (transactions || []).slice(0, 5);

  const getTransactionIcon = (type: string) => {
    switch (type) {
      case 'send':
        return <ArrowUpRight className="h-4 w-4 text-loss" />;
      case 'receive':
        return <ArrowDownLeft className="h-4 w-4 text-profit" />;
      case 'swap':
        return <RefreshCw className="h-4 w-4 text-primary" />;
      default:
        return <CheckCircle className="h-4 w-4 text-muted-foreground" />;
    }
  };

  const formatTime = (timestamp: number) => {
    const now = Date.now();
    const diff = now - timestamp;
    const hours = Math.floor(diff / (1000 * 60 * 60));
    const minutes = Math.floor(diff / (1000 * 60));
    
    if (hours >= 24) {
      const days = Math.floor(hours / 24);
      return `${days}d ago`;
    } else if (hours >= 1) {
      return `${hours}h ago`;
    } else {
      return `${minutes}m ago`;
    }
  };

  if (isLoading) {
    return (
      <Card>
        <CardHeader>
          <div className="flex items-center justify-between">
            <CardTitle>Recent Activity</CardTitle>
            <Button variant="outline" size="sm" disabled>
              View All
            </Button>
          </div>
        </CardHeader>
        <CardContent>
          <div className="space-y-4">
            {[...Array(3)].map((_, i) => (
              <div key={i} className="flex items-center space-x-3">
                <Skeleton className="h-8 w-8 rounded-full" />
                <div className="flex-1 space-y-2">
                  <Skeleton className="h-4 w-[250px]" />
                  <Skeleton className="h-3 w-[200px]" />
                </div>
              </div>
            ))}
          </div>
        </CardContent>
      </Card>
    );
  }

  return (
    <Card>
      <CardHeader>
        <div className="flex items-center justify-between">
          <CardTitle>Recent Activity</CardTitle>
          <Button variant="outline" size="sm" asChild>
            <Link to="/transactions">View All</Link>
          </Button>
        </div>
      </CardHeader>
      <CardContent>
        <div className="space-y-4">
          {recentTransactions.map((tx) => (
            <div key={tx.hash} className="flex items-center space-x-3">
              <div className="flex-shrink-0">
                {getTransactionIcon(tx.type)}
              </div>
              <div className="flex-1 min-w-0">
                <div className="flex items-center justify-between">
                  <p className="text-sm font-medium capitalize">
                    {tx.type}
                    {tx.tokenIn && tx.tokenOut && (
                      <span className="text-muted-foreground">
                        {' '}{tx.tokenIn.symbol} → {tx.tokenOut.symbol}
                      </span>
                    )}
                    {tx.tokenIn && !tx.tokenOut && (
                      <span className="text-muted-foreground">
                        {' '}{tx.tokenIn.symbol}
                      </span>
                    )}
                    {!tx.tokenIn && tx.tokenOut && (
                      <span className="text-muted-foreground">
                        {' '}{tx.tokenOut.symbol}
                      </span>
                    )}
                  </p>
                  <p className="text-sm text-muted-foreground">
                    {formatTime(tx.timestamp)}
                  </p>
                </div>
                <div className="flex items-center justify-between">
                  <p className="text-xs text-muted-foreground">
                    {tx.hash.slice(0, 10)}...{tx.hash.slice(-8)}
                  </p>
                  <p className="text-xs font-medium">
                    {tx.tokenOut && tx.type !== 'approve' && (
                      <span className="text-loss">-{tx.tokenOut.amount}</span>
                    )}
                    {tx.tokenIn && (
                      <span className="text-profit">+{tx.tokenIn.amount}</span>
                    )}
                  </p>
                </div>
              </div>
            </div>
          ))}
          
          {recentTransactions.length === 0 && (
            <div className="text-center py-6">
              <p className="text-sm text-muted-foreground">
                No recent transactions
              </p>
            </div>
          )}
        </div>
      </CardContent>
    </Card>
  );
}