import { useState } from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import { Shield, Wallet, Mail } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Separator } from '@/components/ui/separator';
import { WalletConnectButton } from '@/components/wallet/wallet-connect-button';
import { useAuth } from '@/contexts/auth-context';
import { useWalletStore } from '@/stores/wallet';
import { useToast } from '@/hooks/use-toast';

export function SignInPage() {
  const [email, setEmail] = useState('');
  const [isLinkingEmail, setIsLinkingEmail] = useState(false);
  const [showEmailForm, setShowEmailForm] = useState(false);
  
  const { signIn, linkEmail, isLoading, isAuthenticated } = useAuth();
  const { isConnected } = useWalletStore();
  const { toast } = useToast();
  const navigate = useNavigate();
  const location = useLocation();

  // Redirect if already authenticated
  if (isAuthenticated) {
    const from = location.state?.from?.pathname || '/dashboard';
    navigate(from, { replace: true });
    return null;
  }

  const handleSiweSignIn = async () => {
    if (!isConnected) {
      toast({
        title: "Connect wallet first",
        description: "Please connect your wallet to sign in with Ethereum",
        variant: "destructive"
      });
      return;
    }

    await signIn();
  };

  const handleEmailLink = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!email) return;

    setIsLinkingEmail(true);
    try {
      await linkEmail(email);
      setShowEmailForm(false);
      setEmail('');
    } catch (error) {
      toast({
        title: "Failed to link email",
        description: error instanceof Error ? error.message : "Please try again",
        variant: "destructive"
      });
    } finally {
      setIsLinkingEmail(false);
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-background to-muted/20 p-4">
      <Card className="w-full max-w-md bg-gradient-card shadow-elegant border-0">
        <CardHeader className="text-center space-y-4">
          <div className="w-16 h-16 bg-gradient-primary rounded-2xl flex items-center justify-center mx-auto shadow-glow">
            <Shield className="h-8 w-8 text-primary-foreground" />
          </div>
          <div>
            <CardTitle className="text-2xl font-bold">Sign In</CardTitle>
            <CardDescription className="text-muted-foreground">
              Authenticate with your Ethereum wallet to access your portfolio
            </CardDescription>
          </div>
        </CardHeader>

        <CardContent className="space-y-6">
          {/* Wallet Connection */}
          <div className="space-y-4">
            <div className="text-center">
              <h3 className="text-lg font-semibold mb-2">Connect Your Wallet</h3>
              <p className="text-sm text-muted-foreground mb-4">
                {isConnected ? 'Wallet connected! Now sign the message to authenticate.' : 'Connect your wallet to get started'}
              </p>
            </div>

            {!isConnected ? (
              <div className="flex justify-center">
                <WalletConnectButton />
              </div>
            ) : (
              <Button 
                onClick={handleSiweSignIn}
                disabled={isLoading}
                className="w-full bg-gradient-primary hover:opacity-90"
                size="lg"
              >
                <Wallet className="mr-2 h-4 w-4" />
                {isLoading ? 'Signing...' : 'Sign Message to Authenticate'}
              </Button>
            )}
          </div>

          {/* Email Recovery Option */}
          {isConnected && (
            <>
              <Separator />
              
              <div className="space-y-4">
                <div className="text-center">
                  <h3 className="text-lg font-semibold mb-2">Optional: Link Email</h3>
                  <p className="text-sm text-muted-foreground">
                    Link an email to recover your profile settings and receive notifications
                  </p>
                </div>

                {!showEmailForm ? (
                  <Button 
                    variant="outline" 
                    onClick={() => setShowEmailForm(true)}
                    className="w-full"
                  >
                    <Mail className="mr-2 h-4 w-4" />
                    Add Email Recovery
                  </Button>
                ) : (
                  <form onSubmit={handleEmailLink} className="space-y-4">
                    <div>
                      <Label htmlFor="email">Email Address</Label>
                      <Input
                        id="email"
                        type="email"
                        value={email}
                        onChange={(e) => setEmail(e.target.value)}
                        placeholder="your@email.com"
                        required
                      />
                    </div>
                    <div className="flex space-x-2">
                      <Button 
                        type="submit" 
                        disabled={isLinkingEmail || !email}
                        className="flex-1"
                      >
                        {isLinkingEmail ? 'Sending...' : 'Send Verification'}
                      </Button>
                      <Button 
                        type="button" 
                        variant="outline"
                        onClick={() => {
                          setShowEmailForm(false);
                          setEmail('');
                        }}
                      >
                        Cancel
                      </Button>
                    </div>
                  </form>
                )}
              </div>
            </>
          )}
        </CardContent>
      </Card>
    </div>
  );
}