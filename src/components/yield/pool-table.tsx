import { useState } from 'react';
import { ExternalLink, Gift } from 'lucide-react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table';
import { Badge } from '@/components/ui/badge';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Skeleton } from '@/components/ui/skeleton';

interface Pool {
  id: string;
  protocol: string;
  pair: string;
  chain: string;
  apr: number;
  tvl: number;
  userStaked: number;
  rewards: number;
  logoUrl?: string;
}

export function PoolTable() {
  const [filter, setFilter] = useState('all');
  const [isLoading] = useState(false);

  // Mock data - TODO: Replace with real API
  const pools: Pool[] = [
    {
      id: '1',
      protocol: 'Aave',
      pair: 'USDC',
      chain: 'Ethereum',
      apr: 4.2,
      tvl: 2100000000,
      userStaked: 5000,
      rewards: 42.50,
    },
    {
      id: '2',
      protocol: 'Compound',
      pair: 'ETH',
      chain: 'Ethereum',
      apr: 3.8,
      tvl: 1500000000,
      userStaked: 2.5,
      rewards: 15.20,
    },
    {
      id: '3',
      protocol: 'Uniswap V3',
      pair: 'USDC/ETH',
      chain: 'Ethereum',
      apr: 18.5,
      tvl: 450000000,
      userStaked: 1000,
      rewards: 67.80,
    },
  ];

  const filteredPools = pools.filter(pool => 
    filter === 'all' || pool.chain.toLowerCase() === filter
  );

  const handleClaim = (poolId: string) => {
    console.log('Claiming rewards for pool:', poolId);
    // TODO: Implement claim transaction
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

  if (isLoading) {
    return (
      <Card>
        <CardHeader>
          <CardTitle>Yield Pools</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="space-y-4">
            {[...Array(3)].map((_, i) => (
              <div key={i} className="flex items-center space-x-4">
                <Skeleton className="h-10 w-10 rounded-full" />
                <div className="space-y-2 flex-1">
                  <Skeleton className="h-4 w-32" />
                  <Skeleton className="h-3 w-48" />
                </div>
                <Skeleton className="h-8 w-20" />
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
        <div className="flex justify-between items-center">
          <CardTitle>Available Pools</CardTitle>
          <div className="flex items-center space-x-4">
            <Select value={filter} onValueChange={setFilter}>
              <SelectTrigger className="w-32">
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">All Chains</SelectItem>
                <SelectItem value="ethereum">Ethereum</SelectItem>
                <SelectItem value="polygon">Polygon</SelectItem>
                <SelectItem value="arbitrum">Arbitrum</SelectItem>
              </SelectContent>
            </Select>
            <Button size="sm" variant="outline">
              <Gift className="h-4 w-4 mr-2" />
              Claim All
            </Button>
          </div>
        </div>
      </CardHeader>
      <CardContent>
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Pool</TableHead>
              <TableHead>Chain</TableHead>
              <TableHead>APR</TableHead>
              <TableHead>TVL</TableHead>
              <TableHead>Your Stake</TableHead>
              <TableHead>Rewards</TableHead>
              <TableHead className="text-right">Action</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {filteredPools.map((pool) => (
              <TableRow key={pool.id}>
                <TableCell>
                  <div className="flex items-center space-x-3">
                    <div className="w-8 h-8 rounded-full bg-muted flex items-center justify-center">
                      <span className="text-xs font-medium">
                        {pool.protocol.slice(0, 2)}
                      </span>
                    </div>
                    <div>
                      <div className="font-medium">{pool.pair}</div>
                      <div className="text-sm text-muted-foreground">{pool.protocol}</div>
                    </div>
                  </div>
                </TableCell>
                <TableCell>
                  <Badge variant="outline">{pool.chain}</Badge>
                </TableCell>
                <TableCell>
                  <span className="font-medium text-profit">
                    {pool.apr.toFixed(1)}%
                  </span>
                </TableCell>
                <TableCell>{formatLargeCurrency(pool.tvl)}</TableCell>
                <TableCell>
                  {pool.userStaked > 0 ? formatCurrency(pool.userStaked) : '-'}
                </TableCell>
                <TableCell>
                  {pool.rewards > 0 ? (
                    <span className="font-medium text-profit">
                      {formatCurrency(pool.rewards)}
                    </span>
                  ) : (
                    '-'
                  )}
                </TableCell>
                <TableCell className="text-right">
                  <div className="flex items-center justify-end space-x-2">
                    {pool.rewards > 0 && (
                      <Button
                        size="sm"
                        onClick={() => handleClaim(pool.id)}
                        className="bg-gradient-primary hover:opacity-90"
                      >
                        Claim
                      </Button>
                    )}
                    <Button variant="ghost" size="sm">
                      <ExternalLink className="h-4 w-4" />
                    </Button>
                  </div>
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
        
        {filteredPools.length === 0 && (
          <div className="text-center py-8">
            <p className="text-muted-foreground">
              No pools found for the selected filter.
            </p>
          </div>
        )}
      </CardContent>
    </Card>
  );
}