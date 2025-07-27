import React, { useState } from 'react';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { ArrowUpDown, Settings } from 'lucide-react';
import { RouteSelector } from './route-selector';
import { useBridgeRoutes } from '@/hooks/use-bridge-routes';
import { BridgeRoute } from '@/lib/api/bridge-client';

const CHAINS = [
  { id: 1, name: 'Ethereum', symbol: 'ETH' },
  { id: 137, name: 'Polygon', symbol: 'MATIC' },
  { id: 42161, name: 'Arbitrum', symbol: 'ARB' },
  { id: 10, name: 'Optimism', symbol: 'OP' },
];

const TOKENS = [
  { address: '0xA0b86991c431E4dFe7bb8E5f2D5E8b8A8A8b3c8B', symbol: 'USDC', name: 'USD Coin' },
  { address: '0xdAC17F958D2ee523a2206206994597C13D831ec7', symbol: 'USDT', name: 'Tether USD' },
  { address: '0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2', symbol: 'WETH', name: 'Wrapped Ether' },
];

export function BridgeForm() {
  const [fromChain, setFromChain] = useState<number>(1);
  const [toChain, setToChain] = useState<number>(137);
  const [fromToken, setFromToken] = useState(TOKENS[0].address);
  const [toToken, setToToken] = useState(TOKENS[0].address);
  const [amount, setAmount] = useState('');
  const [slippage, setSlippage] = useState(0.5);
  const [selectedRoute, setSelectedRoute] = useState<BridgeRoute | null>(null);
  const [isExecuting, setIsExecuting] = useState(false);

  const { routes, isLoading, fetchRoutes, executeRoute } = useBridgeRoutes();

  const handleSearch = async () => {
    if (!amount || !fromToken || !toToken) return;
    
    await fetchRoutes({
      fromChain,
      toChain,
      fromToken: fromToken as `0x${string}`,
      toToken: toToken as `0x${string}`,
      fromAmount: (parseFloat(amount) * Math.pow(10, 6)).toString(), // Assuming 6 decimals
      slippage,
    });
  };

  const handleExecute = async (route: BridgeRoute) => {
    setIsExecuting(true);
    try {
      await executeRoute(route);
    } finally {
      setIsExecuting(false);
    }
  };

  const swapChains = () => {
    const tempChain = fromChain;
    setFromChain(toChain);
    setToChain(tempChain);
  };

  return (
    <div className="space-y-6">
      <Card>
        <CardHeader>
          <CardTitle>Bridge Assets</CardTitle>
          <CardDescription>
            Transfer your tokens across different blockchains
          </CardDescription>
        </CardHeader>
        <CardContent className="space-y-4">
          {/* From Section */}
          <div className="space-y-2">
            <Label>From</Label>
            <div className="grid grid-cols-2 gap-4">
              <Select value={fromChain.toString()} onValueChange={(value) => setFromChain(Number(value))}>
                <SelectTrigger>
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  {CHAINS.map((chain) => (
                    <SelectItem key={chain.id} value={chain.id.toString()}>
                      {chain.name}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
              <Select value={fromToken} onValueChange={setFromToken}>
                <SelectTrigger>
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  {TOKENS.map((token) => (
                    <SelectItem key={token.address} value={token.address}>
                      {token.symbol}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
            <Input
              type="number"
              placeholder="0.0"
              value={amount}
              onChange={(e) => setAmount(e.target.value)}
            />
          </div>

          {/* Swap Button */}
          <div className="flex justify-center">
            <Button variant="outline" size="icon" onClick={swapChains}>
              <ArrowUpDown className="h-4 w-4" />
            </Button>
          </div>

          {/* To Section */}
          <div className="space-y-2">
            <Label>To</Label>
            <div className="grid grid-cols-2 gap-4">
              <Select value={toChain.toString()} onValueChange={(value) => setToChain(Number(value))}>
                <SelectTrigger>
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  {CHAINS.map((chain) => (
                    <SelectItem key={chain.id} value={chain.id.toString()}>
                      {chain.name}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
              <Select value={toToken} onValueChange={setToToken}>
                <SelectTrigger>
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  {TOKENS.map((token) => (
                    <SelectItem key={token.address} value={token.address}>
                      {token.symbol}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
          </div>

          {/* Settings */}
          <div className="flex items-center gap-2">
            <Settings className="h-4 w-4" />
            <Label>Slippage Tolerance: {slippage}%</Label>
            <Input
              type="number"
              className="w-20"
              value={slippage}
              onChange={(e) => setSlippage(parseFloat(e.target.value))}
              min="0.1"
              max="5"
              step="0.1"
            />
          </div>

          <Button 
            onClick={handleSearch} 
            disabled={!amount || isLoading}
            className="w-full"
          >
            {isLoading ? 'Finding Routes...' : 'Find Routes'}
          </Button>
        </CardContent>
      </Card>

      {/* Route Selection */}
      {routes.length > 0 && (
        <Card>
          <CardHeader>
            <CardTitle>Available Routes</CardTitle>
            <CardDescription>
              Choose the best route for your bridge transaction
            </CardDescription>
          </CardHeader>
          <CardContent>
            <RouteSelector
              routes={routes}
              selectedRoute={selectedRoute}
              onSelectRoute={setSelectedRoute}
              onExecute={handleExecute}
              isExecuting={isExecuting}
            />
          </CardContent>
        </Card>
      )}
    </div>
  );
}