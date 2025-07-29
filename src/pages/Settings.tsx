import { useState, useEffect } from 'react';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Separator } from '@/components/ui/separator';
import { Badge } from '@/components/ui/badge';
import { Eye, EyeOff, Save, Trash2, AlertTriangle, CheckCircle } from 'lucide-react';
import { toast } from 'sonner';

interface ApiKeys {
  alchemyApiKey: string;
  coingeckoApiKey: string;
  etherscanApiKey: string;
  infuraApiKey: string;
}

const Settings = () => {
  const [apiKeys, setApiKeys] = useState<ApiKeys>({
    alchemyApiKey: '',
    coingeckoApiKey: '',
    etherscanApiKey: '',
    infuraApiKey: ''
  });

  const [showKeys, setShowKeys] = useState({
    alchemyApiKey: false,
    coingeckoApiKey: false,
    etherscanApiKey: false,
    infuraApiKey: false
  });

  const [hasChanges, setHasChanges] = useState(false);

  // Load saved API keys on component mount
  useEffect(() => {
    const savedKeys = localStorage.getItem('defi_api_keys');
    if (savedKeys) {
      try {
        const parsedKeys = JSON.parse(savedKeys);
        setApiKeys(parsedKeys);
      } catch (error) {
        console.error('Failed to load API keys from localStorage:', error);
      }
    }
  }, []);

  const handleInputChange = (key: keyof ApiKeys, value: string) => {
    setApiKeys(prev => ({ ...prev, [key]: value }));
    setHasChanges(true);
  };

  const toggleShowKey = (key: keyof typeof showKeys) => {
    setShowKeys(prev => ({ ...prev, [key]: !prev[key] }));
  };

  const saveApiKeys = () => {
    try {
      localStorage.setItem('defi_api_keys', JSON.stringify(apiKeys));
      setHasChanges(false);
      toast.success('API keys saved successfully!');
    } catch (error) {
      console.error('Failed to save API keys:', error);
      toast.error('Failed to save API keys. Please try again.');
    }
  };

  const clearApiKeys = () => {
    setApiKeys({
      alchemyApiKey: '',
      coingeckoApiKey: '',
      etherscanApiKey: '',
      infuraApiKey: ''
    });
    localStorage.removeItem('defi_api_keys');
    setHasChanges(false);
    toast.success('All API keys cleared!');
  };

  const apiKeyConfigs = [
    {
      key: 'alchemyApiKey' as keyof ApiKeys,
      label: 'Alchemy API Key',
      description: 'Required for blockchain data and RPC calls',
      required: true,
      placeholder: 'Enter your Alchemy API key...'
    },
    {
      key: 'coingeckoApiKey' as keyof ApiKeys,
      label: 'CoinGecko API Key',
      description: 'Required for cryptocurrency prices and market data',
      required: true,
      placeholder: 'Enter your CoinGecko API key...'
    },
    {
      key: 'etherscanApiKey' as keyof ApiKeys,
      label: 'Etherscan API Key',
      description: 'Optional: Enhanced transaction data and gas tracking',
      required: false,
      placeholder: 'Enter your Etherscan API key...'
    },
    {
      key: 'infuraApiKey' as keyof ApiKeys,
      label: 'Infura API Key',
      description: 'Optional: Additional blockchain infrastructure',
      required: false,
      placeholder: 'Enter your Infura API key...'
    }
  ];

  const getConnectionStatus = () => {
    const requiredKeys = ['alchemyApiKey', 'coingeckoApiKey'];
    const hasRequired = requiredKeys.every(key => apiKeys[key as keyof ApiKeys].trim() !== '');
    const optionalKeys = ['etherscanApiKey', 'infuraApiKey'];
    const hasOptional = optionalKeys.some(key => apiKeys[key as keyof ApiKeys].trim() !== '');
    
    return { hasRequired, hasOptional };
  };

  const { hasRequired, hasOptional } = getConnectionStatus();

  return (
    <div className="container mx-auto py-8 px-4 max-w-4xl">
      <div className="space-y-6">
        {/* Header */}
        <div>
          <h1 className="text-3xl font-bold">Settings</h1>
          <p className="text-muted-foreground mt-2">
            Configure your API keys for real-time DeFi data integration
          </p>
        </div>

        {/* Connection Status */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              {hasRequired ? (
                <CheckCircle className="h-5 w-5 text-green-500" />
              ) : (
                <AlertTriangle className="h-5 w-5 text-yellow-500" />
              )}
              Connection Status
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="flex flex-wrap gap-2">
              <Badge variant={hasRequired ? "default" : "secondary"}>
                {hasRequired ? "Essential APIs Connected" : "Missing Essential APIs"}
              </Badge>
              {hasOptional && (
                <Badge variant="outline">
                  Optional Enhancers Connected
                </Badge>
              )}
            </div>
            <p className="text-sm text-muted-foreground mt-2">
              {hasRequired 
                ? "Your app is ready for real-time data!" 
                : "Configure essential API keys to enable real-time features."}
            </p>
          </CardContent>
        </Card>

        {/* API Keys Configuration */}
        <Card>
          <CardHeader>
            <CardTitle>API Keys</CardTitle>
            <CardDescription>
              Your API keys are stored locally on your device and never sent to our servers.
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-6">
            {apiKeyConfigs.map((config) => (
              <div key={config.key} className="space-y-2">
                <div className="flex items-center gap-2">
                  <Label htmlFor={config.key}>{config.label}</Label>
                  {config.required && (
                    <Badge variant="destructive" className="text-xs">
                      Essential
                    </Badge>
                  )}
                </div>
                
                <div className="relative">
                  <Input
                    id={config.key}
                    type={showKeys[config.key] ? "text" : "password"}
                    value={apiKeys[config.key]}
                    onChange={(e) => handleInputChange(config.key, e.target.value)}
                    placeholder={config.placeholder}
                    className="pr-10"
                  />
                  <Button
                    type="button"
                    variant="ghost"
                    size="sm"
                    className="absolute right-0 top-0 h-full px-3 py-2 hover:bg-transparent"
                    onClick={() => toggleShowKey(config.key)}
                  >
                    {showKeys[config.key] ? (
                      <EyeOff className="h-4 w-4" />
                    ) : (
                      <Eye className="h-4 w-4" />
                    )}
                  </Button>
                </div>
                
                <p className="text-sm text-muted-foreground">
                  {config.description}
                </p>
              </div>
            ))}

            <Separator />

            {/* Actions */}
            <div className="flex flex-col sm:flex-row gap-3">
              <Button 
                onClick={saveApiKeys}
                disabled={!hasChanges}
                className="flex items-center gap-2"
              >
                <Save className="h-4 w-4" />
                Save Changes
              </Button>
              
              <Button 
                variant="outline"
                onClick={clearApiKeys}
                className="flex items-center gap-2"
              >
                <Trash2 className="h-4 w-4" />
                Clear All Keys
              </Button>
            </div>

            {hasChanges && (
              <p className="text-sm text-yellow-600 dark:text-yellow-400">
                You have unsaved changes. Don't forget to save!
              </p>
            )}
          </CardContent>
        </Card>

        {/* Security Notice */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <AlertTriangle className="h-5 w-5 text-yellow-500" />
              Security Notice
            </CardTitle>
          </CardHeader>
          <CardContent className="space-y-3">
            <div className="text-sm text-muted-foreground space-y-2">
              <p>
                • Your API keys are stored locally in your browser and are never transmitted to our servers
              </p>
              <p>
                • Always keep your API keys secure and never share them publicly
              </p>
              <p>
                • If you suspect your keys have been compromised, regenerate them immediately
              </p>
              <p>
                • Clear your browser data will remove all stored API keys
              </p>
            </div>
          </CardContent>
        </Card>

        {/* Getting API Keys Help */}
        <Card>
          <CardHeader>
            <CardTitle>How to Get API Keys</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="grid gap-4 md:grid-cols-2">
              <div className="space-y-2">
                <h4 className="font-medium">Essential APIs</h4>
                <div className="text-sm text-muted-foreground space-y-1">
                  <p>• <strong>Alchemy:</strong> Sign up at alchemy.com</p>
                  <p>• <strong>CoinGecko:</strong> Get API key at coingecko.com/api</p>
                </div>
              </div>
              <div className="space-y-2">
                <h4 className="font-medium">Optional Enhancers</h4>
                <div className="text-sm text-muted-foreground space-y-1">
                  <p>• <strong>Etherscan:</strong> Register at etherscan.io/apis</p>
                  <p>• <strong>Infura:</strong> Create project at infura.io</p>
                </div>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  );
};

export default Settings;