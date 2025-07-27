import { ArrowLeftRight, ExternalLink, Clock } from 'lucide-react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';

interface BridgeTransaction {
  id: string;
  fromChain: string;
  toChain: string;
  token: string;
  amount: string;
  status: 'pending' | 'completed' | 'failed';
  timestamp: number;
  txHash: string;
  fee: string;
}

export function BridgeHistory() {
  // Mock data - TODO: Replace with real API
  const transactions: BridgeTransaction[] = [
    {
      id: '1',
      fromChain: 'Ethereum',
      toChain: 'Polygon',
      token: 'USDC',
      amount: '500.00',
      status: 'completed',
      timestamp: Date.now() - 3600000,
      txHash: '0x1234...5678',
      fee: '$5.20',
    },
    {
      id: '2',
      fromChain: 'Polygon',
      toChain: 'Arbitrum',
      token: 'ETH',
      amount: '0.25',
      status: 'pending',
      timestamp: Date.now() - 1800000,
      txHash: '0xabcd...efgh',
      fee: '$3.80',
    },
  ];

  const getStatusVariant = (status: string) => {
    switch (status) {
      case 'completed':
        return 'secondary';
      case 'pending':
        return 'outline';
      case 'failed':
        return 'destructive';
      default:
        return 'secondary';
    }
  };

  const getStatusIcon = (status: string) => {
    switch (status) {
      case 'pending':
        return <Clock className="h-3 w-3" />;
      default:
        return null;
    }
  };

  const formatTime = (timestamp: number) => {
    const now = Date.now();
    const diff = now - timestamp;
    const hours = Math.floor(diff / (1000 * 60 * 60));
    const minutes = Math.floor(diff / (1000 * 60));
    
    if (hours >= 1) {
      return `${hours}h ago`;
    } else {
      return `${minutes}m ago`;
    }
  };

  const openExplorer = (txHash: string) => {
    window.open(`https://etherscan.io/tx/${txHash}`, '_blank');
  };

  return (
    <Card>
      <CardHeader>
        <CardTitle>Bridge History</CardTitle>
      </CardHeader>
      <CardContent>
        <div className="space-y-4">
          {transactions.map((tx) => (
            <div key={tx.id} className="space-y-2 p-3 border rounded-lg">
              <div className="flex items-center justify-between">
                <div className="flex items-center space-x-2">
                  <ArrowLeftRight className="h-4 w-4 text-muted-foreground" />
                  <span className="text-sm font-medium">
                    {tx.fromChain} → {tx.toChain}
                  </span>
                  <Badge variant={getStatusVariant(tx.status)} className="h-5">
                    <div className="flex items-center space-x-1">
                      {getStatusIcon(tx.status)}
                      <span className="text-xs capitalize">{tx.status}</span>
                    </div>
                  </Badge>
                </div>
                <span className="text-xs text-muted-foreground">
                  {formatTime(tx.timestamp)}
                </span>
              </div>
              
              <div className="flex items-center justify-between text-sm">
                <div>
                  <div className="font-medium">{tx.amount} {tx.token}</div>
                  <div className="text-xs text-muted-foreground">Fee: {tx.fee}</div>
                </div>
                <Button
                  variant="ghost"
                  size="sm"
                  onClick={() => openExplorer(tx.txHash)}
                  className="h-auto p-1"
                >
                  <ExternalLink className="h-3 w-3" />
                </Button>
              </div>
              
              <div className="text-xs text-muted-foreground">
                {tx.txHash}
              </div>
            </div>
          ))}
          
          {transactions.length === 0 && (
            <div className="text-center py-6">
              <p className="text-sm text-muted-foreground">
                No bridge transactions found
              </p>
              <p className="text-xs text-muted-foreground mt-1">
                Your bridge history will appear here
              </p>
            </div>
          )}
        </div>
      </CardContent>
    </Card>
  );
}