import { useState } from 'react';
import { User, Mail, Shield, Key, ExternalLink } from 'lucide-react';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Badge } from '@/components/ui/badge';
import { Separator } from '@/components/ui/separator';
import { useAuth } from '@/contexts/auth-context';
import { useToast } from '@/hooks/use-toast';
import { formatAddress } from '@/lib/utils';

export function ProfileSettings() {
  const [email, setEmail] = useState('');
  const [isLinkingEmail, setIsLinkingEmail] = useState(false);
  
  const { user, linkEmail, signOut } = useAuth();
  const { toast } = useToast();

  const handleEmailLink = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!email) return;

    setIsLinkingEmail(true);
    try {
      await linkEmail(email);
      setEmail('');
      toast({
        title: "Verification email sent",
        description: "Check your email to complete verification"
      });
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

  const handleSignOut = () => {
    signOut();
  };

  if (!user) {
    return (
      <div className="max-w-2xl mx-auto p-6">
        <Card>
          <CardContent className="p-6 text-center">
            <Shield className="h-12 w-12 text-muted-foreground mx-auto mb-4" />
            <h3 className="text-lg font-medium mb-2">Not Authenticated</h3>
            <p className="text-muted-foreground">Please sign in to view your profile</p>
          </CardContent>
        </Card>
      </div>
    );
  }

  return (
    <div className="max-w-2xl mx-auto space-y-6">
      {/* Profile Overview */}
      <Card className="bg-gradient-card shadow-elegant border-0">
        <CardHeader>
          <div className="flex items-center justify-between">
            <div className="flex items-center space-x-3">
              <div className="w-12 h-12 bg-gradient-primary rounded-full flex items-center justify-center">
                <User className="h-6 w-6 text-primary-foreground" />
              </div>
              <div>
                <CardTitle>Profile Settings</CardTitle>
                <CardDescription>Manage your account and security settings</CardDescription>
              </div>
            </div>
            <Badge variant="outline" className="text-success border-success">
              Authenticated
            </Badge>
          </div>
        </CardHeader>
      </Card>

      {/* Wallet Information */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center space-x-2">
            <Shield className="h-5 w-5" />
            <span>Wallet Information</span>
          </CardTitle>
          <CardDescription>
            Your connected Ethereum wallet details
          </CardDescription>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="grid grid-cols-2 gap-4">
            <div>
              <Label className="text-sm font-medium">Wallet Address</Label>
              <div className="flex items-center space-x-2 mt-1">
                <code className="text-sm bg-muted px-2 py-1 rounded">
                  {formatAddress(user.address)}
                </code>
                <Button
                  variant="ghost"
                  size="sm"
                  onClick={() => {
                    navigator.clipboard.writeText(user.address);
                    toast({ title: "Address copied!" });
                  }}
                >
                  <ExternalLink className="h-3 w-3" />
                </Button>
              </div>
            </div>
            <div>
              <Label className="text-sm font-medium">Last Login</Label>
              <p className="text-sm text-muted-foreground mt-1">
                {new Date(user.lastLoginAt).toLocaleDateString()}
              </p>
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Email Settings */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center space-x-2">
            <Mail className="h-5 w-5" />
            <span>Email Settings</span>
          </CardTitle>
          <CardDescription>
            Link an email for account recovery and notifications
          </CardDescription>
        </CardHeader>
        <CardContent className="space-y-4">
          {user.email ? (
            <div className="flex items-center justify-between p-4 bg-muted/50 rounded-lg">
              <div className="flex items-center space-x-3">
                <Mail className="h-5 w-5 text-muted-foreground" />
                <div>
                  <p className="font-medium">{user.email}</p>
                  <p className="text-sm text-muted-foreground">
                    {user.emailVerified ? (
                      <span className="text-success">✓ Verified</span>
                    ) : (
                      <span className="text-warning">⚠ Pending verification</span>
                    )}
                  </p>
                </div>
              </div>
            </div>
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
                <p className="text-xs text-muted-foreground mt-1">
                  Used for account recovery and portfolio notifications
                </p>
              </div>
              <Button 
                type="submit" 
                disabled={isLinkingEmail || !email}
                className="w-full"
              >
                {isLinkingEmail ? 'Sending verification...' : 'Add Email'}
              </Button>
            </form>
          )}
        </CardContent>
      </Card>

      {/* Security Actions */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center space-x-2">
            <Key className="h-5 w-5" />
            <span>Security</span>
          </CardTitle>
          <CardDescription>
            Manage your authentication and security settings
          </CardDescription>
        </CardHeader>
        <CardContent>
          <div className="space-y-4">
            <div className="flex items-center justify-between p-4 border border-border rounded-lg">
              <div>
                <h4 className="font-medium">Sign out of all sessions</h4>
                <p className="text-sm text-muted-foreground">
                  Sign out of your current session and clear stored authentication
                </p>
              </div>
              <Button variant="destructive" onClick={handleSignOut}>
                Sign Out
              </Button>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}