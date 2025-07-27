import { useState } from 'react';
import { ArrowLeftRight, Info } from 'lucide-react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Separator } from '@/components/ui/separator';

export function BridgeForm() {
  const [fromChain, setFromChain] = useState('ethereum');
  const [toChain, setToChain] = useState('polygon');
  const [token, setToken] = useState('USDC');
  const [amount, setAmount] = useState('');

  const chains = [
    { id: 'ethereum', name: 'Ethereum', symbol: 'ETH' },
    { id: 'polygon', name: 'Polygon', symbol: 'MATIC' },
    { id: 'arbitrum', name: 'Arbitrum', symbol: 'ARB' },
    { id: 'optimism', name: 'Optimism', symbol: 'OP' },
  ];

  const tokens = ['USDC', 'USDT', 'ETH', 'WBTC'];

  const handleSwapChains = () => {
    const temp = fromChain;
    setFromChain(toChain);
    setToChain(temp);
  };

  return (
    <Card>
      <CardHeader>
        <CardTitle>Bridge Assets</CardTitle>
      </CardHeader>
      <CardContent className="space-y-6">
        {/* From Chain */}
        <div className="space-y-2">
          <label className="text-sm font-medium">From</label>
          <Select value={fromChain} onValueChange={setFromChain}>
            <SelectTrigger>
              <SelectValue />
            </SelectTrigger>
            <SelectContent>
              {chains.map((chain) => (
                <SelectItem key={chain.id} value={chain.id}>
                  {chain.name}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        </div>

        {/* Swap Button */}
        <div className="flex justify-center">
          <Button
            variant="ghost"
            size="icon"
            onClick={handleSwapChains}
            className="rounded-full border"
          >
            <ArrowLeftRight className="h-4 w-4" />
          </Button>
        </div>

        {/* To Chain */}
        <div className="space-y-2">
          <label className="text-sm font-medium">To</label>
          <Select value={toChain} onValueChange={setToChain}>
            <SelectTrigger>
              <SelectValue />
            </SelectTrigger>
            <SelectContent>
              {chains.map((chain) => (
                <SelectItem key={chain.id} value={chain.id}>
                  {chain.name}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        </div>

        {/* Token Selection */}
        <div className="space-y-2">
          <label className="text-sm font-medium">Token</label>
          <Select value={token} onValueChange={setToken}>
            <SelectTrigger>
              <SelectValue />
            </SelectTrigger>
            <SelectContent>
              {tokens.map((token) => (
                <SelectItem key={token} value={token}>
                  {token}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        </div>

        {/* Amount */}
        <div className="space-y-2">
          <div className="flex justify-between items-center">
            <label className="text-sm font-medium">Amount</label>
            <span className="text-xs text-muted-foreground">
              Balance: 1,000.00 {token}
            </span>
          </div>
          <Input
            placeholder="0.0"
            value={amount}
            onChange={(e) => setAmount(e.target.value)}
            className="text-lg h-12"
          />
        </div>

        {/* Route Information */}
        {amount && (
          <>
            <Separator />
            <div className="space-y-2 text-sm">
              <div className="flex justify-between">
                <span className="text-muted-foreground">Bridge Fee</span>
                <span>~$2.50</span>
              </div>
              <div className="flex justify-between">
                <span className="text-muted-foreground">Gas Cost</span>
                <span>~$8.20</span>
              </div>
              <div className="flex justify-between">
                <span className="text-muted-foreground">Time</span>
                <span>~5 minutes</span>
              </div>
              <div className="flex justify-between font-medium">
                <span>You'll receive</span>
                <span>{amount} {token}</span>
              </div>
            </div>
          </>
        )}

        {/* Bridge Button */}
        <Button 
          size="lg" 
          className="w-full bg-gradient-primary hover:opacity-90"
          disabled={!amount}
        >
          {!amount ? 'Enter Amount' : 'Continue in Wallet'}
        </Button>

        {/* Development Notice */}
        <div className="flex items-start space-x-2 p-3 bg-primary/5 border border-primary/20 rounded-lg">
          <Info className="h-4 w-4 text-primary mt-0.5 flex-shrink-0" />
          <div className="text-xs">
            <p className="font-medium text-primary mb-1">Development Mode</p>
            <p className="text-muted-foreground">
              Bridge functionality will be available when LI.FI or Socket integration is completed.
            </p>
          </div>
        </div>
      </CardContent>
    </Card>
  );
}