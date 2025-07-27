import React, { useState } from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { Separator } from '@/components/ui/separator';
import { Clock, Zap, DollarSign, ArrowRight } from 'lucide-react';
import { BridgeRoute } from '@/lib/api/bridge-client';

interface RouteSelectorProps {
  routes: BridgeRoute[];
  selectedRoute: BridgeRoute | null;
  onSelectRoute: (route: BridgeRoute) => void;
  onExecute: (route: BridgeRoute) => void;
  isExecuting?: boolean;
}

export function RouteSelector({ 
  routes, 
  selectedRoute, 
  onSelectRoute, 
  onExecute,
  isExecuting = false 
}: RouteSelectorProps) {
  const [sortBy, setSortBy] = useState<'time' | 'fee' | 'amount'>('amount');

  const sortedRoutes = [...routes].sort((a, b) => {
    switch (sortBy) {
      case 'time':
        return a.estimatedTime - b.estimatedTime;
      case 'fee':
        return parseFloat(a.fees.total) - parseFloat(b.fees.total);
      case 'amount':
        return parseFloat(b.toAmount) - parseFloat(a.toAmount);
      default:
        return 0;
    }
  });

  const formatTime = (seconds: number) => {
    const minutes = Math.floor(seconds / 60);
    return `~${minutes}m`;
  };

  const formatAmount = (amount: string, decimals: number = 6) => {
    const value = parseFloat(amount) / Math.pow(10, decimals);
    return value.toFixed(4);
  };

  return (
    <div className="space-y-4">
      {/* Sort Controls */}
      <div className="flex gap-2">
        <Button
          variant={sortBy === 'amount' ? 'default' : 'outline'}
          size="sm"
          onClick={() => setSortBy('amount')}
        >
          Best Rate
        </Button>
        <Button
          variant={sortBy === 'time' ? 'default' : 'outline'}
          size="sm"
          onClick={() => setSortBy('time')}
        >
          Fastest
        </Button>
        <Button
          variant={sortBy === 'fee' ? 'default' : 'outline'}
          size="sm"
          onClick={() => setSortBy('fee')}
        >
          Lowest Fee
        </Button>
      </div>

      {/* Route Cards */}
      <div className="space-y-3">
        {sortedRoutes.map((route) => (
          <Card
            key={route.id}
            className={`cursor-pointer transition-all hover:shadow-md ${
              selectedRoute?.id === route.id ? 'ring-2 ring-primary' : ''
            }`}
            onClick={() => onSelectRoute(route)}
          >
            <CardContent className="p-4">
              <div className="flex items-center justify-between mb-3">
                <div className="flex items-center gap-2">
                  <Badge variant="outline" className="capitalize">
                    {route.provider}
                  </Badge>
                  <Badge variant="secondary">
                    {route.steps[0].protocol}
                  </Badge>
                </div>
                <div className="flex items-center gap-4 text-sm text-muted-foreground">
                  <div className="flex items-center gap-1">
                    <Clock className="h-4 w-4" />
                    {formatTime(route.estimatedTime)}
                  </div>
                  <div className="flex items-center gap-1">
                    <Zap className="h-4 w-4" />
                    {route.estimatedGas} gas
                  </div>
                  <div className="flex items-center gap-1">
                    <DollarSign className="h-4 w-4" />
                    ${route.fees.total}
                  </div>
                </div>
              </div>

              <div className="flex items-center justify-between">
                <div className="text-lg font-semibold">
                  {formatAmount(route.toAmount)} 
                  <span className="text-sm font-normal text-muted-foreground ml-1">
                    received
                  </span>
                </div>
                <ArrowRight className="h-5 w-5 text-muted-foreground" />
              </div>

              {/* Fee Breakdown */}
              <div className="mt-3 pt-3 border-t border-border/50">
                <div className="flex justify-between text-sm">
                  <span>Bridge Fee:</span>
                  <span>${route.fees.bridgeFee}</span>
                </div>
                <div className="flex justify-between text-sm">
                  <span>Gas Fee:</span>
                  <span>${route.fees.gasFee}</span>
                </div>
                <Separator className="my-2" />
                <div className="flex justify-between text-sm font-medium">
                  <span>Total Fees:</span>
                  <span>${route.fees.total}</span>
                </div>
              </div>
            </CardContent>
          </Card>
        ))}
      </div>

      {/* Execute Button */}
      {selectedRoute && (
        <Button
          className="w-full"
          size="lg"
          onClick={() => onExecute(selectedRoute)}
          disabled={isExecuting}
        >
          {isExecuting ? 'Executing...' : 'Execute Bridge'}
        </Button>
      )}
    </div>
  );
}