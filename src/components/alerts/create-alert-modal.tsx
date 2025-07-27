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
import { useAlertsStore } from '@/stores/alerts';
import { AccountingMethod } from '@/lib/pnl-calculator';

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
  const [cooldown, setCooldown] = useState('60');
  const [retryAttempts, setRetryAttempts] = useState('3');
  const [webhookUrl, setWebhookUrl] = useState('');
  const [email, setEmail] = useState('user@example.com');
  
  const { toast } = useToast();
  const { addAlert } = useAlertsStore();

  const tokens = ['ETH', 'BTC', 'USDC', 'USDT', 'UNI', 'AAVE', 'COMP'];
  const aprPools = ['AAVE Pool', 'Compound USDC', 'Uniswap V3'];

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

    if (channel === 'telegram' && !webhookUrl) {
      toast({
        title: "Missing Webhook URL",
        description: "Please provide a Telegram webhook URL",
        variant: "destructive",
      });
      return;
    }

    const alertData = {
      type,
      token,
      condition,
      threshold: parseFloat(threshold),
      isActive: true,
      channel,
      notificationSettings: {
        cooldown: parseInt(cooldown),
        retryAttempts: parseInt(retryAttempts),
        ...(channel === 'telegram' ? { webhookUrl } : { email })
      }
    };

    addAlert(alertData);

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
    setCooldown('60');
    setRetryAttempts('3');
    setWebhookUrl('');
  };

  const handleClose = () => {
    onClose();
    resetForm();
  };

  const getTokenOptions = () => {
    if (type === 'apr') {
      return aprPools;
    }
    return tokens;
  };

  return (
    <Dialog open={isOpen} onOpenChange={handleClose}>
      <DialogContent className="sm:max-w-md max-h-[90vh] overflow-y-auto">
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
              <Label htmlFor="token">{type === 'apr' ? 'Pool' : 'Token/Asset'}</Label>
              <Select value={token} onValueChange={setToken}>
                <SelectTrigger>
                  <SelectValue placeholder={`Select ${type === 'apr' ? 'pool' : 'token'}`} />
                </SelectTrigger>
                <SelectContent>
                  {getTokenOptions().map((t) => (
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

            {channel === 'telegram' && (
              <div className="space-y-2">
                <Label htmlFor="webhook">Telegram Webhook URL</Label>
                <Input
                  id="webhook"
                  value={webhookUrl}
                  onChange={(e) => setWebhookUrl(e.target.value)}
                  placeholder="https://api.telegram.org/bot..."
                />
              </div>
            )}

            {channel === 'email' && (
              <div className="space-y-2">
                <Label htmlFor="email">Email Address</Label>
                <Input
                  id="email"
                  type="email"
                  value={email}
                  onChange={(e) => setEmail(e.target.value)}
                  placeholder="your@email.com"
                />
              </div>
            )}

            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-2">
                <Label htmlFor="cooldown">Cooldown (minutes)</Label>
                <Select value={cooldown} onValueChange={setCooldown}>
                  <SelectTrigger>
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="15">15 minutes</SelectItem>
                    <SelectItem value="30">30 minutes</SelectItem>
                    <SelectItem value="60">1 hour</SelectItem>
                    <SelectItem value="120">2 hours</SelectItem>
                    <SelectItem value="240">4 hours</SelectItem>
                  </SelectContent>
                </Select>
              </div>

              <div className="space-y-2">
                <Label htmlFor="retries">Retry Attempts</Label>
                <Select value={retryAttempts} onValueChange={setRetryAttempts}>
                  <SelectTrigger>
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="1">1 attempt</SelectItem>
                    <SelectItem value="2">2 attempts</SelectItem>
                    <SelectItem value="3">3 attempts</SelectItem>
                    <SelectItem value="5">5 attempts</SelectItem>
                  </SelectContent>
                </Select>
              </div>
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