import { useState } from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Separator } from '@/components/ui/separator';
import { useAuth } from '@/contexts/auth-context';
import { useWalletStore } from '@/stores/wallet';
import { WalletConnectButton } from '@/components/wallet/wallet-connect-button';
import { useNavigate, useLocation } from 'react-router-dom';
import { useEffect } from 'react';
import { Shield, Mail, Loader2 } from 'lucide-react';

export default function SignIn() {
  const [email, setEmail] = useState('');
  const [isLinkingEmail, setIsLinkingEmail] = useState(false);
  const { signIn, linkEmail, isAuthenticated, isLoading } = useAuth();
  const { isConnected } = useWalletStore();
  const navigate = useNavigate();
  const location = useLocation();

  const from = location.state?.from?.pathname || '/dashboard';

  useEffect(() => {
    if (isAuthenticated) {
      navigate(from, { replace: true });
    }
  }, [isAuthenticated, navigate, from]);

  const handleSignIn = async () => {
    try {
      await signIn();
    } catch (error) {
      console.error('Sign in failed:', error);
    }
  };

  const handleLinkEmail = async () => {
    if (!email.trim()) return;
    
    setIsLinkingEmail(true);
    try {
      await linkEmail(email);
      setEmail('');
    } catch (error) {
      console.error('Email linking failed:', error);
    } finally {
      setIsLinkingEmail(false);
    }
  };

  if (isLoading) {
    return (
      <div className="min-h-screen bg-gradient-to-br from-background to-muted/20 flex items-center justify-center p-8">
        <Card className="w-full max-w-md">
          <CardContent className="p-6 text-center">
            <Loader2 className="h-8 w-8 animate-spin mx-auto mb-4 text-primary" />
            <p className="text-muted-foreground">Checking authentication...</p>
          </CardContent>
        </Card>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-background to-muted/20 flex items-center justify-center p-8">
      <Card className="w-full max-w-md">
        <CardHeader className="text-center">
          <div className="mx-auto w-12 h-12 bg-primary/10 rounded-full flex items-center justify-center mb-4">
            <Shield className="h-6 w-6 text-primary" />
          </div>
          <CardTitle className="text-2xl">Sign In</CardTitle>
          <p className="text-muted-foreground">
            Connect your wallet and sign with Ethereum
          </p>
        </CardHeader>
        <CardContent className="space-y-6">
          <div className="space-y-4">
            <WalletConnectButton />
            
            {isConnected && (
              <Button 
                onClick={handleSignIn}
                className="w-full"
                disabled={isLoading}
              >
                {isLoading ? (
                  <>
                    <Loader2 className="h-4 w-4 mr-2 animate-spin" />
                    Signing...
                  </>
                ) : (
                  'Sign with Ethereum'
                )}
              </Button>
            )}
          </div>

          {isAuthenticated && (
            <>
              <Separator />
              <div className="space-y-4">
                <div className="text-center">
                  <h3 className="font-medium">Link Email (Optional)</h3>
                  <p className="text-sm text-muted-foreground">
                    Link your email to recover profile settings
                  </p>
                </div>
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
                <Button 
                  onClick={handleLinkEmail}
                  variant="outline"
                  className="w-full"
                  disabled={isLinkingEmail || !email.trim()}
                >
                  {isLinkingEmail ? (
                    <>
                      <Loader2 className="h-4 w-4 mr-2 animate-spin" />
                      Sending...
                    </>
                  ) : (
                    <>
                      <Mail className="h-4 w-4 mr-2" />
                      Link Email
                    </>
                  )}
                </Button>
              </div>
            </>
          )}
        </CardContent>
      </Card>
    </div>
  );
}