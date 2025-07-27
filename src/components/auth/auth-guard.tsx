import { ReactNode } from 'react';
import { Navigate, useLocation } from 'react-router-dom';
import { useAuth } from '@/contexts/auth-context';
import { useWalletStore } from '@/stores/wallet';
import { Skeleton } from '@/components/ui/skeleton';
import { Card, CardContent } from '@/components/ui/card';

interface AuthGuardProps {
  children: ReactNode;
  requireAuth?: boolean;
  fallbackPath?: string;
}

export function AuthGuard({ 
  children, 
  requireAuth = true, 
  fallbackPath = '/auth/signin' 
}: AuthGuardProps) {
  const { isAuthenticated, isLoading } = useAuth();
  const { isConnected } = useWalletStore();
  const location = useLocation();

  // Show loading state while checking authentication
  if (isLoading) {
    return <AuthLoadingState />;
  }

  // If auth is required but user is not authenticated
  if (requireAuth && !isAuthenticated) {
    // If wallet is not connected, redirect to home
    if (!isConnected) {
      return <Navigate to="/" state={{ from: location }} replace />;
    }
    // If wallet is connected but not authenticated, redirect to sign in
    return <Navigate to={fallbackPath} state={{ from: location }} replace />;
  }

  // If auth is not required but user is authenticated, allow access
  return <>{children}</>;
}

function AuthLoadingState() {
  return (
    <div className="min-h-screen bg-gradient-to-br from-background to-muted/20 p-8">
      <div className="max-w-4xl mx-auto space-y-8">
        <div className="text-center space-y-4">
          <Skeleton className="h-8 w-48 mx-auto" />
          <Skeleton className="h-4 w-64 mx-auto" />
        </div>
        
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
          {Array.from({ length: 3 }).map((_, i) => (
            <Card key={i}>
              <CardContent className="p-6">
                <Skeleton className="h-6 w-24 mb-4" />
                <Skeleton className="h-8 w-32 mb-2" />
                <Skeleton className="h-4 w-20" />
              </CardContent>
            </Card>
          ))}
        </div>
        
        <Card>
          <CardContent className="p-6">
            <Skeleton className="h-96 w-full" />
          </CardContent>
        </Card>
      </div>
    </div>
  );
}

// Higher-order component for route protection
export function withAuthGuard<P extends object>(
  Component: React.ComponentType<P>,
  requireAuth: boolean = true
) {
  return function AuthGuardedComponent(props: P) {
    return (
      <AuthGuard requireAuth={requireAuth}>
        <Component {...props} />
      </AuthGuard>
    );
  };
}