import { useEffect, useState } from 'react';
import { useSearchParams, useNavigate } from 'react-router-dom';
import { Check, Mail, AlertCircle } from 'lucide-react';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { useAuth } from '@/contexts/auth-context';
import { useToast } from '@/hooks/use-toast';

export function EmailVerificationPage() {
  const [searchParams] = useSearchParams();
  const [status, setStatus] = useState<'verifying' | 'success' | 'error'>('verifying');
  const [error, setError] = useState<string>('');
  
  const { verifyEmail } = useAuth();
  const { toast } = useToast();
  const navigate = useNavigate();

  useEffect(() => {
    const token = searchParams.get('token');
    
    if (!token) {
      setStatus('error');
      setError('Invalid verification link');
      return;
    }

    verifyEmailToken(token);
  }, [searchParams]);

  const verifyEmailToken = async (token: string) => {
    try {
      await verifyEmail(token);
      setStatus('success');
      
      // Redirect to dashboard after a short delay
      setTimeout(() => {
        navigate('/dashboard');
      }, 3000);
      
    } catch (error) {
      setStatus('error');
      setError(error instanceof Error ? error.message : 'Verification failed');
      
      toast({
        title: "Email verification failed",
        description: "The verification link may be expired or invalid",
        variant: "destructive"
      });
    }
  };

  const handleReturnToDashboard = () => {
    navigate('/dashboard');
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-background to-muted/20 p-4">
      <Card className="w-full max-w-md bg-gradient-card shadow-elegant border-0">
        <CardHeader className="text-center space-y-4">
          <div className="w-16 h-16 rounded-2xl flex items-center justify-center mx-auto">
            {status === 'verifying' && (
              <div className="w-16 h-16 bg-muted rounded-2xl flex items-center justify-center animate-pulse">
                <Mail className="h-8 w-8 text-muted-foreground" />
              </div>
            )}
            {status === 'success' && (
              <div className="w-16 h-16 bg-success/10 rounded-2xl flex items-center justify-center">
                <Check className="h-8 w-8 text-success" />
              </div>
            )}
            {status === 'error' && (
              <div className="w-16 h-16 bg-destructive/10 rounded-2xl flex items-center justify-center">
                <AlertCircle className="h-8 w-8 text-destructive" />
              </div>
            )}
          </div>
          
          <div>
            <CardTitle className="text-2xl font-bold">
              {status === 'verifying' && 'Verifying Email'}
              {status === 'success' && 'Email Verified!'}
              {status === 'error' && 'Verification Failed'}
            </CardTitle>
            <CardDescription className="text-muted-foreground">
              {status === 'verifying' && 'Please wait while we verify your email address...'}
              {status === 'success' && 'Your email has been successfully verified. You will be redirected shortly.'}
              {status === 'error' && 'There was a problem verifying your email address.'}
            </CardDescription>
          </div>
        </CardHeader>

        <CardContent className="text-center space-y-6">
          {status === 'error' && (
            <div className="p-4 bg-destructive/10 rounded-lg">
              <p className="text-sm text-destructive font-medium">
                {error}
              </p>
            </div>
          )}

          {status === 'success' && (
            <div className="space-y-4">
              <p className="text-sm text-muted-foreground">
                You can now receive email notifications for your portfolio alerts and updates.
              </p>
              <Button 
                onClick={handleReturnToDashboard}
                className="w-full bg-gradient-primary hover:opacity-90"
              >
                Continue to Dashboard
              </Button>
            </div>
          )}

          {status === 'error' && (
            <div className="space-y-4">
              <p className="text-sm text-muted-foreground">
                You can try requesting a new verification email from your profile settings.
              </p>
              <Button 
                onClick={handleReturnToDashboard}
                variant="outline"
                className="w-full"
              >
                Return to Dashboard
              </Button>
            </div>
          )}
        </CardContent>
      </Card>
    </div>
  );
}