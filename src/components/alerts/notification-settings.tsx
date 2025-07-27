import { useState } from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Button } from '@/components/ui/button';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Separator } from '@/components/ui/separator';
import { Mail, MessageSquare, Settings } from 'lucide-react';
import { useToast } from '@/hooks/use-toast';

interface NotificationSettings {
  email: {
    address: string;
    enabled: boolean;
  };
  telegram: {
    webhookUrl: string;
    chatId: string;
    enabled: boolean;
  };
  defaults: {
    cooldown: number;
    retryAttempts: number;
  };
}

export function NotificationSettings() {
  const [settings, setSettings] = useState<NotificationSettings>({
    email: {
      address: 'user@example.com',
      enabled: true
    },
    telegram: {
      webhookUrl: '',
      chatId: '',
      enabled: false
    },
    defaults: {
      cooldown: 60,
      retryAttempts: 3
    }
  });
  const [isLoading, setIsLoading] = useState(false);
  const { toast } = useToast();

  const handleSave = async () => {
    setIsLoading(true);
    
    try {
      // TODO: Save to backend/Supabase
      console.log('Saving notification settings:', settings);
      
      await new Promise(resolve => setTimeout(resolve, 1000)); // Simulate API call
      
      toast({
        title: "Settings Saved",
        description: "Your notification preferences have been updated.",
      });
    } catch (error) {
      toast({
        title: "Error",
        description: "Failed to save settings. Please try again.",
        variant: "destructive",
      });
    } finally {
      setIsLoading(false);
    }
  };

  const testNotification = async (channel: 'email' | 'telegram') => {
    setIsLoading(true);
    
    try {
      if (channel === 'telegram' && !settings.telegram.webhookUrl) {
        toast({
          title: "Missing Configuration",
          description: "Please configure your Telegram webhook URL first.",
          variant: "destructive",
        });
        return;
      }

      // TODO: Send test notification
      console.log(`Testing ${channel} notification`);
      
      await new Promise(resolve => setTimeout(resolve, 2000)); // Simulate API call
      
      toast({
        title: "Test Notification Sent",
        description: `Check your ${channel} for the test message.`,
      });
    } catch (error) {
      toast({
        title: "Test Failed",
        description: `Failed to send ${channel} test notification.`,
        variant: "destructive",
      });
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <Card className="bg-gradient-card shadow-card border-0">
      <CardHeader>
        <div className="flex items-center space-x-2">
          <Settings className="h-5 w-5" />
          <CardTitle>Notification Settings</CardTitle>
        </div>
      </CardHeader>
      <CardContent className="space-y-6">
        {/* Email Settings */}
        <div className="space-y-4">
          <div className="flex items-center space-x-2">
            <Mail className="h-4 w-4" />
            <h3 className="text-lg font-medium">Email Notifications</h3>
          </div>
          
          <div className="space-y-2">
            <Label htmlFor="email">Email Address</Label>
            <Input
              id="email"
              type="email"
              value={settings.email.address}
              onChange={(e) => setSettings(prev => ({
                ...prev,
                email: { ...prev.email, address: e.target.value }
              }))}
              placeholder="your@email.com"
            />
          </div>

          <Button 
            variant="outline" 
            onClick={() => testNotification('email')}
            disabled={isLoading || !settings.email.address}
          >
            Test Email
          </Button>
        </div>

        <Separator />

        {/* Telegram Settings */}
        <div className="space-y-4">
          <div className="flex items-center space-x-2">
            <MessageSquare className="h-4 w-4" />
            <h3 className="text-lg font-medium">Telegram Notifications</h3>
          </div>
          
          <div className="space-y-2">
            <Label htmlFor="telegram-webhook">Webhook URL</Label>
            <Input
              id="telegram-webhook"
              value={settings.telegram.webhookUrl}
              onChange={(e) => setSettings(prev => ({
                ...prev,
                telegram: { ...prev.telegram, webhookUrl: e.target.value }
              }))}
              placeholder="https://api.telegram.org/bot..."
            />
            <p className="text-xs text-muted-foreground">
              Create a Telegram bot and get your webhook URL from @BotFather
            </p>
          </div>

          <div className="space-y-2">
            <Label htmlFor="telegram-chat">Chat ID (Optional)</Label>
            <Input
              id="telegram-chat"
              value={settings.telegram.chatId}
              onChange={(e) => setSettings(prev => ({
                ...prev,
                telegram: { ...prev.telegram, chatId: e.target.value }
              }))}
              placeholder="123456789"
            />
          </div>

          <Button 
            variant="outline" 
            onClick={() => testNotification('telegram')}
            disabled={isLoading || !settings.telegram.webhookUrl}
          >
            Test Telegram
          </Button>
        </div>

        <Separator />

        {/* Default Settings */}
        <div className="space-y-4">
          <h3 className="text-lg font-medium">Default Alert Settings</h3>
          
          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-2">
              <Label htmlFor="cooldown">Cooldown Period</Label>
              <Select 
                value={settings.defaults.cooldown.toString()} 
                onValueChange={(value) => setSettings(prev => ({
                  ...prev,
                  defaults: { ...prev.defaults, cooldown: parseInt(value) }
                }))}
              >
                <SelectTrigger>
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="15">15 minutes</SelectItem>
                  <SelectItem value="30">30 minutes</SelectItem>
                  <SelectItem value="60">1 hour</SelectItem>
                  <SelectItem value="120">2 hours</SelectItem>
                  <SelectItem value="240">4 hours</SelectItem>
                  <SelectItem value="480">8 hours</SelectItem>
                  <SelectItem value="1440">24 hours</SelectItem>
                </SelectContent>
              </Select>
            </div>

            <div className="space-y-2">
              <Label htmlFor="retries">Retry Attempts</Label>
              <Select 
                value={settings.defaults.retryAttempts.toString()} 
                onValueChange={(value) => setSettings(prev => ({
                  ...prev,
                  defaults: { ...prev.defaults, retryAttempts: parseInt(value) }
                }))}
              >
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

        <div className="flex justify-end">
          <Button 
            onClick={handleSave}
            disabled={isLoading}
            className="bg-gradient-primary hover:opacity-90"
          >
            {isLoading ? 'Saving...' : 'Save Settings'}
          </Button>
        </div>
      </CardContent>
    </Card>
  );
}