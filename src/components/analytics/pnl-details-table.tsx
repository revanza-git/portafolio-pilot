import { useState } from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table';
import { Badge } from '@/components/ui/badge';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { PnLCalculation, formatCurrency, formatPercent } from '@/lib/pnl-calculator';
import { format } from 'date-fns';

interface PnLDetailsTableProps {
  calculation: PnLCalculation;
}

export function PnLDetailsTable({ calculation }: PnLDetailsTableProps) {
  const [sortBy, setSortBy] = useState<'timestamp' | 'pnl' | 'symbol'>('timestamp');
  const [sortOrder, setSortOrder] = useState<'asc' | 'desc'>('desc');

  const sortedTrades = [...calculation.trades].sort((a, b) => {
    let comparison = 0;
    
    switch (sortBy) {
      case 'timestamp':
        comparison = a.timestamp - b.timestamp;
        break;
      case 'pnl':
        comparison = (a.realizedPnL || 0) - (b.realizedPnL || 0);
        break;
      case 'symbol':
        comparison = a.symbol.localeCompare(b.symbol);
        break;
    }
    
    return sortOrder === 'asc' ? comparison : -comparison;
  });

  const sortedLots = [...calculation.lots].sort((a, b) => {
    if (sortBy === 'symbol') {
      return sortOrder === 'asc' 
        ? a.symbol.localeCompare(b.symbol)
        : b.symbol.localeCompare(a.symbol);
    }
    return sortOrder === 'asc'
      ? a.timestamp - b.timestamp
      : b.timestamp - a.timestamp;
  });

  const handleSort = (column: typeof sortBy) => {
    if (sortBy === column) {
      setSortOrder(sortOrder === 'asc' ? 'desc' : 'asc');
    } else {
      setSortBy(column);
      setSortOrder('desc');
    }
  };

  return (
    <Card className="bg-gradient-card shadow-card border-0">
      <CardHeader>
        <CardTitle>P&L Details</CardTitle>
      </CardHeader>
      <CardContent>
        <Tabs defaultValue="trades" className="w-full">
          <TabsList className="grid w-full grid-cols-2">
            <TabsTrigger value="trades">Trade History</TabsTrigger>
            <TabsTrigger value="lots">Current Lots</TabsTrigger>
          </TabsList>
          
          <TabsContent value="trades" className="space-y-4">
            <div className="rounded-md border border-border">
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead 
                      className="cursor-pointer hover:bg-muted/50"
                      onClick={() => handleSort('timestamp')}
                    >
                      Date {sortBy === 'timestamp' && (sortOrder === 'asc' ? '↑' : '↓')}
                    </TableHead>
                    <TableHead 
                      className="cursor-pointer hover:bg-muted/50"
                      onClick={() => handleSort('symbol')}
                    >
                      Asset {sortBy === 'symbol' && (sortOrder === 'asc' ? '↑' : '↓')}
                    </TableHead>
                    <TableHead>Type</TableHead>
                    <TableHead className="text-right">Quantity</TableHead>
                    <TableHead className="text-right">Price</TableHead>
                    <TableHead 
                      className="text-right cursor-pointer hover:bg-muted/50"
                      onClick={() => handleSort('pnl')}
                    >
                      Realized P&L {sortBy === 'pnl' && (sortOrder === 'asc' ? '↑' : '↓')}
                    </TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {sortedTrades.map((trade, index) => (
                    <TableRow key={`${trade.hash}-${index}`}>
                      <TableCell>
                        {format(new Date(trade.timestamp), 'MMM dd, yyyy')}
                      </TableCell>
                      <TableCell className="font-medium">{trade.symbol}</TableCell>
                      <TableCell>
                        <Badge variant={trade.type === 'buy' ? 'secondary' : 'outline'}>
                          {trade.type.toUpperCase()}
                        </Badge>
                      </TableCell>
                      <TableCell className="text-right">
                        {trade.quantity.toFixed(6)}
                      </TableCell>
                      <TableCell className="text-right">
                        {formatCurrency(trade.price)}
                      </TableCell>
                      <TableCell className="text-right">
                        {trade.realizedPnL !== undefined ? (
                          <span className={trade.realizedPnL >= 0 ? 'text-profit' : 'text-loss'}>
                            {formatCurrency(trade.realizedPnL)}
                          </span>
                        ) : (
                          <span className="text-muted-foreground">-</span>
                        )}
                      </TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
            </div>
          </TabsContent>
          
          <TabsContent value="lots" className="space-y-4">
            <div className="rounded-md border border-border">
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead 
                      className="cursor-pointer hover:bg-muted/50"
                      onClick={() => handleSort('symbol')}
                    >
                      Asset {sortBy === 'symbol' && (sortOrder === 'asc' ? '↑' : '↓')}
                    </TableHead>
                    <TableHead className="text-right">Quantity</TableHead>
                    <TableHead className="text-right">Cost Basis</TableHead>
                    <TableHead className="text-right">Total Cost</TableHead>
                    <TableHead 
                      className="cursor-pointer hover:bg-muted/50"
                      onClick={() => handleSort('timestamp')}
                    >
                      Acquired {sortBy === 'timestamp' && (sortOrder === 'asc' ? '↑' : '↓')}
                    </TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {sortedLots.map((lot) => (
                    <TableRow key={lot.id}>
                      <TableCell className="font-medium">{lot.symbol}</TableCell>
                      <TableCell className="text-right">
                        {lot.quantity.toFixed(6)}
                      </TableCell>
                      <TableCell className="text-right">
                        {formatCurrency(lot.costBasis)}
                      </TableCell>
                      <TableCell className="text-right">
                        {formatCurrency(lot.quantity * lot.costBasis)}
                      </TableCell>
                      <TableCell>
                        {format(new Date(lot.timestamp), 'MMM dd, yyyy')}
                      </TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
            </div>
          </TabsContent>
        </Tabs>
      </CardContent>
    </Card>
  );
}