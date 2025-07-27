import { RefreshCw, ExternalLink } from 'lucide-react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { usePortfolioStore } from '@/stores/portfolio';

export function SwapHistory() {
  const { transactions } = usePortfolioStore();
  
  // Filter for swap transactions only
  const swapTransactions = transactions.filter(tx => tx.type === 'swap').slice(0, 5);

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

  const openEtherscan = (hash: string) => {
    window.open(`https://etherscan.io/tx/${hash}`, '_blank');
  };

  return (
    <Card>
      <CardHeader>
        <div className="flex items-center justify-between">
          <CardTitle>Recent Swaps</CardTitle>
          <RefreshCw className="h-4 w-4 text-muted-foreground" />
        </div>
      </CardHeader>
      <CardContent>
        <div className="space-y-4">
          {swapTransactions.map((tx) => (
            <div key={tx.hash} className="space-y-2 p-3 border rounded-lg">
              <div className="flex items-center justify-between">
                <div className="flex items-center space-x-2">
                  <span className="text-sm font-medium">
                    {tx.tokenIn?.symbol} → {tx.tokenOut?.symbol}
                  </span>
                  <Badge variant="secondary" className="h-5">
                    {tx.status}
                  </Badge>
                </div>
                <span className="text-xs text-muted-foreground">
                  {formatTime(tx.timestamp)}
                </span>
              </div>
              
              <div className="flex items-center justify-between text-sm">
                <div>
                  <div className="text-loss">-{tx.tokenIn?.amount} {tx.tokenIn?.symbol}</div>
                  <div className="text-profit">+{tx.tokenOut?.amount} {tx.tokenOut?.symbol}</div>
                </div>
                <div className="text-right">
                  <div className="text-muted-foreground">${tx.tokenIn?.usdValue.toLocaleString()}</div>
                  <div className="text-xs text-muted-foreground">Gas: {tx.gasFee}</div>
                </div>
              </div>
              
              <div className="flex items-center justify-between text-xs">
                <span className="text-muted-foreground">
                  {tx.hash.slice(0, 10)}...{tx.hash.slice(-8)}
                </span>
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
          ))}
          
          {swapTransactions.length === 0 && (
            <div className="text-center py-6">
              <p className="text-sm text-muted-foreground">
                No swap history found
              </p>
              <p className="text-xs text-muted-foreground mt-1">
                Your completed swaps will appear here
              </p>
            </div>
          )}
        </div>
      </CardContent>
    </Card>
  );
}