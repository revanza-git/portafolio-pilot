import { ArrowUpRight, ArrowDownLeft, RefreshCw, CheckCircle, ExternalLink, Clock, XCircle } from 'lucide-react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Skeleton } from '@/components/ui/skeleton';
import { Transaction } from '@/stores/portfolio';

interface TransactionListProps {
  transactions: Transaction[];
  isLoading: boolean;
  filters: {
    type: string;
    status: string;
    timeRange: string;
  };
}

export function TransactionList({ transactions, isLoading, filters }: TransactionListProps) {
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

  const getStatusIcon = (status: string) => {
    switch (status) {
      case 'success':
        return <CheckCircle className="h-3 w-3 text-profit" />;
      case 'pending':
        return <Clock className="h-3 w-3 text-warning" />;
      case 'failed':
        return <XCircle className="h-3 w-3 text-loss" />;
      default:
        return null;
    }
  };

  const getStatusVariant = (status: string) => {
    switch (status) {
      case 'success':
        return 'secondary';
      case 'pending':
        return 'outline';
      case 'failed':
        return 'destructive';
      default:
        return 'secondary';
    }
  };

  const formatDate = (timestamp: number) => {
    return new Date(timestamp).toLocaleDateString('en-US', {
      month: 'short',
      day: 'numeric',
      year: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    });
  };

  const openEtherscan = (hash: string) => {
    window.open(`https://etherscan.io/tx/${hash}`, '_blank');
  };

  // Filter transactions based on filters
  const filteredTransactions = transactions.filter((tx) => {
    if (filters.type !== 'all' && tx.type !== filters.type) return false;
    if (filters.status !== 'all' && tx.status !== filters.status) return false;
    
    if (filters.timeRange !== 'all') {
      const now = Date.now();
      const timeRanges: Record<string, number> = {
        '1d': 24 * 60 * 60 * 1000,
        '7d': 7 * 24 * 60 * 60 * 1000,
        '30d': 30 * 24 * 60 * 60 * 1000,
        '90d': 90 * 24 * 60 * 60 * 1000,
      };
      
      const range = timeRanges[filters.timeRange];
      if (range && (now - tx.timestamp) > range) return false;
    }
    
    return true;
  });

  if (isLoading) {
    return (
      <Card>
        <CardHeader>
          <CardTitle>Transactions</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="space-y-4">
            {[...Array(5)].map((_, i) => (
              <div key={i} className="flex items-center space-x-4 p-4 border rounded-lg">
                <Skeleton className="h-8 w-8" />
                <div className="flex-1 space-y-2">
                  <Skeleton className="h-4 w-32" />
                  <Skeleton className="h-3 w-48" />
                </div>
                <div className="space-y-2">
                  <Skeleton className="h-4 w-20" />
                  <Skeleton className="h-3 w-16" />
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
        <CardTitle>
          Transactions 
          {filteredTransactions.length !== transactions.length && (
            <span className="text-muted-foreground font-normal">
              ({filteredTransactions.length} of {transactions.length})
            </span>
          )}
        </CardTitle>
      </CardHeader>
      <CardContent>
        <div className="space-y-4">
          {filteredTransactions.map((tx) => (
            <div key={tx.hash} className="flex items-center space-x-4 p-4 border rounded-lg hover:bg-muted/50 transition-colors">
              <div className="flex-shrink-0">
                {getTransactionIcon(tx.type)}
              </div>
              
              <div className="flex-1 min-w-0">
                <div className="flex items-center space-x-2 mb-1">
                  <span className="text-sm font-medium capitalize">{tx.type}</span>
                  <Badge variant={getStatusVariant(tx.status)} className="h-5">
                    <div className="flex items-center space-x-1">
                      {getStatusIcon(tx.status)}
                      <span className="text-xs capitalize">{tx.status}</span>
                    </div>
                  </Badge>
                </div>
                
                <div className="flex items-center justify-between text-xs text-muted-foreground">
                  <span>{formatDate(tx.timestamp)}</span>
                  <div className="flex items-center space-x-2">
                    <span>Gas: {tx.gasFee}</span>
                    <Button
                      variant="ghost"
                      size="sm"
                      onClick={() => openEtherscan(tx.hash)}
                      className="h-auto p-1"
                    >
                      <ExternalLink className="h-3 w-3" />
                    </Button>
                  </div>
                </div>
                
                <div className="text-xs text-muted-foreground mt-1">
                  {tx.hash.slice(0, 10)}...{tx.hash.slice(-8)}
                </div>
              </div>
              
              <div className="text-right">
                {tx.tokenOut && tx.type !== 'approve' && (
                  <div className="text-sm text-loss">
                    -{tx.tokenOut.amount} {tx.tokenOut.symbol}
                  </div>
                )}
                {tx.tokenIn && (
                  <div className="text-sm text-profit">
                    +{tx.tokenIn.amount} {tx.tokenIn.symbol}
                  </div>
                )}
                {tx.type === 'approve' && tx.tokenOut && (
                  <div className="text-sm text-muted-foreground">
                    Approved {tx.tokenOut.symbol}
                  </div>
                )}
              </div>
            </div>
          ))}
          
          {filteredTransactions.length === 0 && (
            <div className="text-center py-8">
              <p className="text-muted-foreground">
                No transactions found matching the selected filters.
              </p>
            </div>
          )}
        </div>
      </CardContent>
    </Card>
  );
}