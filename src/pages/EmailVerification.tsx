import { useEffect, useState } from 'react';
import { useSearchParams, useNavigate } from 'react-router-dom';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { useAuth } from '@/contexts/auth-context';
import { CheckCircle, XCircle, Loader2, Mail } from 'lucide-react';

export default function EmailVerification() {
  const [searchParams] = useSearchParams();
  const [status, setStatus] = useState<'loading' | 'success' | 'error'>('loading');
  const [message, setMessage] = useState('');
  const { verifyEmail } = useAuth();
  const navigate = useNavigate();

  const token = searchParams.get('token');

  useEffect(() => {
    if (!token) {
      setStatus('error');
      setMessage('Invalid verification link. Token is missing.');
      return;
    }

    verifyEmailToken();
  }, [token]);

  const verifyEmailToken = async () => {
    if (!token) return;

    try {
      await verifyEmail(token);
      setStatus('success');
      setMessage('Your email has been successfully verified!');
    } catch (error) {
      setStatus('error');
      setMessage('Email verification failed. The link may be expired or invalid.');
    }
  };

  const handleContinue = () => {
    navigate('/dashboard');
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-background to-muted/20 flex items-center justify-center p-8">
      <Card className="w-full max-w-md">
        <CardHeader className="text-center">
          <div className="mx-auto w-12 h-12 rounded-full flex items-center justify-center mb-4">
            {status === 'loading' && (
              <div className="bg-primary/10">
                <Loader2 className="h-6 w-6 text-primary animate-spin" />
              </div>
            )}
            {status === 'success' && (
              <div className="bg-green-100 dark:bg-green-900">
                <CheckCircle className="h-6 w-6 text-green-600 dark:text-green-400" />
              </div>
            )}
            {status === 'error' && (
              <div className="bg-red-100 dark:bg-red-900">
                <XCircle className="h-6 w-6 text-red-600 dark:text-red-400" />
              </div>
            )}
          </div>
          <CardTitle className="text-2xl">
            {status === 'loading' && 'Verifying Email'}
            {status === 'success' && 'Email Verified'}
            {status === 'error' && 'Verification Failed'}
          </CardTitle>
        </CardHeader>
        <CardContent className="text-center space-y-4">
          <div className="flex items-center justify-center mb-4">
            <Mail className="h-8 w-8 text-muted-foreground" />
          </div>
          
          <p className="text-muted-foreground">
            {message}
          </p>

          {status === 'success' && (
            <Button onClick={handleContinue} className="w-full">
              Continue to Dashboard
            </Button>
          )}

          {status === 'error' && (
            <div className="space-y-2">
              <Button onClick={() => navigate('/auth/signin')} className="w-full">
                Back to Sign In
              </Button>
              <Button 
                variant="outline" 
                onClick={() => navigate('/')} 
                className="w-full"
              >
                Go to Home
              </Button>
            </div>
          )}

          {status === 'loading' && (
            <div className="flex items-center justify-center">
              <Loader2 className="h-5 w-5 animate-spin mr-2" />
              <span className="text-sm text-muted-foreground">
                Please wait while we verify your email...
              </span>
            </div>
          )}
        </CardContent>
      </Card>
    </div>
  );
}