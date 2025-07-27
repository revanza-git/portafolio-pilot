import { useState } from 'react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog';
import { useToast } from '@/hooks/use-toast';

interface CreateAlertModalProps {
  isOpen: boolean;
  onClose: () => void;
}

export function CreateAlertModal({ isOpen, onClose }: CreateAlertModalProps) {
  const [type, setType] = useState<'price' | 'apr' | 'allowance'>('price');
  const [token, setToken] = useState('');
  const [condition, setCondition] = useState<'above' | 'below'>('above');
  const [threshold, setThreshold] = useState('');
  const [channel, setChannel] = useState<'email' | 'telegram'>('email');
  const { toast } = useToast();

  const tokens = ['ETH', 'BTC', 'USDC', 'USDT', 'UNI', 'AAVE', 'COMP'];

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    
    if (!token || !threshold) {
      toast({
        title: "Missing Information",
        description: "Please fill in all required fields",
        variant: "destructive",
      });
      return;
    }

    // TODO: Submit to API
    console.log('Creating alert:', {
      type,
      token,
      condition,
      threshold: parseFloat(threshold),
      channel,
    });

    toast({
      title: "Alert Created",
      description: `You'll be notified when ${token} goes ${condition} ${threshold}`,
    });

    onClose();
    resetForm();
  };

  const resetForm = () => {
    setType('price');
    setToken('');
    setCondition('above');
    setThreshold('');
    setChannel('email');
  };

  const handleClose = () => {
    onClose();
    resetForm();
  };

  return (
    <Dialog open={isOpen} onOpenChange={handleClose}>
      <DialogContent className="sm:max-w-md">
        <form onSubmit={handleSubmit}>
          <DialogHeader>
            <DialogTitle>Create Alert</DialogTitle>
            <DialogDescription>
              Get notified when your selected conditions are met
            </DialogDescription>
          </DialogHeader>
          
          <div className="space-y-4 py-4">
            <div className="space-y-2">
              <Label htmlFor="type">Alert Type</Label>
              <Select value={type} onValueChange={(value: any) => setType(value)}>
                <SelectTrigger>
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="price">Price Alert</SelectItem>
                  <SelectItem value="apr">APR Alert</SelectItem>
                  <SelectItem value="allowance">Allowance Alert</SelectItem>
                </SelectContent>
              </Select>
            </div>

            <div className="space-y-2">
              <Label htmlFor="token">Token/Asset</Label>
              <Select value={token} onValueChange={setToken}>
                <SelectTrigger>
                  <SelectValue placeholder="Select token" />
                </SelectTrigger>
                <SelectContent>
                  {tokens.map((t) => (
                    <SelectItem key={t} value={t}>
                      {t}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>

            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-2">
                <Label htmlFor="condition">Condition</Label>
                <Select value={condition} onValueChange={(value: any) => setCondition(value)}>
                  <SelectTrigger>
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="above">Above</SelectItem>
                    <SelectItem value="below">Below</SelectItem>
                  </SelectContent>
                </Select>
              </div>

              <div className="space-y-2">
                <Label htmlFor="threshold">
                  Threshold {type === 'price' ? '($)' : type === 'apr' ? '(%)' : ''}
                </Label>
                <Input
                  id="threshold"
                  type="number"
                  step="0.01"
                  value={threshold}
                  onChange={(e) => setThreshold(e.target.value)}
                  placeholder="0.00"
                />
              </div>
            </div>

            <div className="space-y-2">
              <Label htmlFor="channel">Notification Channel</Label>
              <Select value={channel} onValueChange={(value: any) => setChannel(value)}>
                <SelectTrigger>
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="email">Email</SelectItem>
                  <SelectItem value="telegram">Telegram</SelectItem>
                </SelectContent>
              </Select>
            </div>
          </div>

          <DialogFooter>
            <Button type="button" variant="outline" onClick={handleClose}>
              Cancel
            </Button>
            <Button type="submit" className="bg-gradient-primary hover:opacity-90">
              Create Alert
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  );
}