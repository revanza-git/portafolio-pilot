import { useState } from 'react';
import { Star, StarOff, TrendingUp, DollarSign } from 'lucide-react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';

interface WatchlistItem {
  id: string;
  type: 'token' | 'pool' | 'protocol';
  name: string;
  symbol?: string;
  price?: number;
  change24h?: number;
  apr?: number;
  tvl?: number;
  chain?: string;
}

export function WatchlistGrid() {
  const [watchlist, setWatchlist] = useState<WatchlistItem[]>([
    {
      id: '1',
      type: 'token',
      name: 'Ethereum',
      symbol: 'ETH',
      price: 2458.30,
      change24h: 2.5,
    },
    {
      id: '2',
      type: 'token',
      name: 'Chainlink',
      symbol: 'LINK',
      price: 14.82,
      change24h: -1.2,
    },
    {
      id: '3',
      type: 'pool',
      name: 'USDC/ETH Pool',
      apr: 18.5,
      tvl: 450000000,
      chain: 'Ethereum',
    },
    {
      id: '4',
      type: 'protocol',
      name: 'Aave',
      tvl: 12000000000,
      change24h: 0.8,
    },
  ]);

  const removeFromWatchlist = (id: string) => {
    setWatchlist(watchlist.filter(item => item.id !== id));
  };

  const formatCurrency = (amount: number) => {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'USD',
      minimumFractionDigits: 2,
    }).format(amount);
  };

  const formatLargeCurrency = (amount: number) => {
    if (amount >= 1e9) return `$${(amount / 1e9).toFixed(1)}B`;
    if (amount >= 1e6) return `$${(amount / 1e6).toFixed(1)}M`;
    if (amount >= 1e3) return `$${(amount / 1e3).toFixed(1)}K`;
    return formatCurrency(amount);
  };

  const tokens = watchlist.filter(item => item.type === 'token');
  const pools = watchlist.filter(item => item.type === 'pool');
  const protocols = watchlist.filter(item => item.type === 'protocol');

  const renderTokenCard = (item: WatchlistItem) => (
    <Card key={item.id} className="hover:shadow-md transition-shadow">
      <CardContent className="p-4">
        <div className="flex justify-between items-start mb-3">
          <div className="flex items-center space-x-3">
            <div className="w-10 h-10 rounded-full bg-muted flex items-center justify-center">
              <span className="text-sm font-medium">
                {item.symbol?.slice(0, 2)}
              </span>
            </div>
            <div>
              <div className="font-medium">{item.symbol}</div>
              <div className="text-sm text-muted-foreground">{item.name}</div>
            </div>
          </div>
          <Button
            variant="ghost"
            size="icon"
            onClick={() => removeFromWatchlist(item.id)}
            className="text-warning hover:text-warning"
          >
            <Star className="h-4 w-4 fill-current" />
          </Button>
        </div>
        
        <div className="space-y-2">
          <div className="text-lg font-bold">
            {formatCurrency(item.price || 0)}
          </div>
          <div className={`text-sm flex items-center ${
            (item.change24h || 0) >= 0 ? 'text-profit' : 'text-loss'
          }`}>
            <TrendingUp className={`h-3 w-3 mr-1 ${
              (item.change24h || 0) < 0 ? 'rotate-180' : ''
            }`} />
            {(item.change24h || 0) >= 0 ? '+' : ''}{item.change24h?.toFixed(2)}%
          </div>
        </div>
      </CardContent>
    </Card>
  );

  const renderPoolCard = (item: WatchlistItem) => (
    <Card key={item.id} className="hover:shadow-md transition-shadow">
      <CardContent className="p-4">
        <div className="flex justify-between items-start mb-3">
          <div>
            <div className="font-medium">{item.name}</div>
            <div className="text-sm text-muted-foreground">
              <Badge variant="outline">{item.chain}</Badge>
            </div>
          </div>
          <Button
            variant="ghost"
            size="icon"
            onClick={() => removeFromWatchlist(item.id)}
            className="text-warning hover:text-warning"
          >
            <Star className="h-4 w-4 fill-current" />
          </Button>
        </div>
        
        <div className="space-y-2">
          <div className="text-lg font-bold text-profit">
            {item.apr?.toFixed(1)}% APR
          </div>
          <div className="text-sm text-muted-foreground">
            TVL: {formatLargeCurrency(item.tvl || 0)}
          </div>
        </div>
      </CardContent>
    </Card>
  );

  const renderProtocolCard = (item: WatchlistItem) => (
    <Card key={item.id} className="hover:shadow-md transition-shadow">
      <CardContent className="p-4">
        <div className="flex justify-between items-start mb-3">
          <div className="flex items-center space-x-3">
            <div className="w-10 h-10 rounded-full bg-muted flex items-center justify-center">
              <span className="text-sm font-medium">
                {item.name.slice(0, 2)}
              </span>
            </div>
            <div>
              <div className="font-medium">{item.name}</div>
              <div className="text-sm text-muted-foreground">Protocol</div>
            </div>
          </div>
          <Button
            variant="ghost"
            size="icon"
            onClick={() => removeFromWatchlist(item.id)}
            className="text-warning hover:text-warning"
          >
            <Star className="h-4 w-4 fill-current" />
          </Button>
        </div>
        
        <div className="space-y-2">
          <div className="text-lg font-bold">
            {formatLargeCurrency(item.tvl || 0)} TVL
          </div>
          <div className={`text-sm flex items-center ${
            (item.change24h || 0) >= 0 ? 'text-profit' : 'text-loss'
          }`}>
            <TrendingUp className={`h-3 w-3 mr-1 ${
              (item.change24h || 0) < 0 ? 'rotate-180' : ''
            }`} />
            {(item.change24h || 0) >= 0 ? '+' : ''}{item.change24h?.toFixed(2)}% 24h
          </div>
        </div>
      </CardContent>
    </Card>
  );

  return (
    <Tabs defaultValue="all" className="w-full">
      <TabsList className="grid w-full grid-cols-4">
        <TabsTrigger value="all">All ({watchlist.length})</TabsTrigger>
        <TabsTrigger value="tokens">Tokens ({tokens.length})</TabsTrigger>
        <TabsTrigger value="pools">Pools ({pools.length})</TabsTrigger>
        <TabsTrigger value="protocols">Protocols ({protocols.length})</TabsTrigger>
      </TabsList>
      
      <TabsContent value="all" className="space-y-4">
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {watchlist.map((item) => {
            if (item.type === 'token') return renderTokenCard(item);
            if (item.type === 'pool') return renderPoolCard(item);
            if (item.type === 'protocol') return renderProtocolCard(item);
            return null;
          })}
        </div>
      </TabsContent>
      
      <TabsContent value="tokens" className="space-y-4">
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {tokens.map(renderTokenCard)}
        </div>
      </TabsContent>
      
      <TabsContent value="pools" className="space-y-4">
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {pools.map(renderPoolCard)}
        </div>
      </TabsContent>
      
      <TabsContent value="protocols" className="space-y-4">
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {protocols.map(renderProtocolCard)}
        </div>
      </TabsContent>

      {watchlist.length === 0 && (
        <div className="text-center py-12">
          <StarOff className="h-12 w-12 text-muted-foreground mx-auto mb-4" />
          <h3 className="text-lg font-medium mb-2">Your watchlist is empty</h3>
          <p className="text-muted-foreground">
            Add tokens, pools, and protocols to track their performance
          </p>
        </div>
      )}
    </Tabs>
  );
}