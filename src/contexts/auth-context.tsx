import { createContext, useContext, useEffect, useState, ReactNode } from 'react';
import { useAccount, useSignMessage } from 'wagmi';
// import { SiweMessage } from 'siwe'; // Temporarily disabled
import { useToast } from '@/hooks/use-toast';

interface User {
  id: string;
  address: string;
  email?: string;
  emailVerified?: boolean;
  lastLoginAt: string;
  createdAt: string;
}

interface AuthState {
  user: User | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  error: string | null;
}

interface AuthContextType extends AuthState {
  signIn: () => Promise<void>;
  signOut: () => void;
  linkEmail: (email: string) => Promise<void>;
  verifyEmail: (token: string) => Promise<void>;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

interface AuthProviderProps {
  children: ReactNode;
}

export function AuthProvider({ children }: AuthProviderProps) {
  const [authState, setAuthState] = useState<AuthState>({
    user: null,
    isAuthenticated: false,
    isLoading: true,
    error: null
  });

  const { address, isConnected } = useAccount();
  const { signMessageAsync } = useSignMessage();
  const { toast } = useToast();

  // Check for existing session on mount
  useEffect(() => {
    checkAuthSession();
  }, []);

  // Clear auth when wallet disconnects
  useEffect(() => {
    if (!isConnected) {
      signOut();
    }
  }, [isConnected]);

  const checkAuthSession = async () => {
    console.log('AuthContext: Checking auth session...');
    try {
      const token = localStorage.getItem('auth_token');
      if (!token) {
        console.log('AuthContext: No token found, setting loading false');
        setAuthState(prev => ({ ...prev, isLoading: false }));
        return;
      }

      console.log('AuthContext: Token found, verifying with backend...');
      // Verify session with backend
      const response = await fetch('/api/auth/verify', {
        headers: {
          'Authorization': `Bearer ${token}`
        }
      });

      if (response.ok) {
        const user = await response.json();
        console.log('AuthContext: Session verified, user:', user);
        setAuthState({
          user,
          isAuthenticated: true,
          isLoading: false,
          error: null
        });
      } else {
        console.log('AuthContext: Session invalid, clearing token');
        // Invalid token, clear it
        localStorage.removeItem('auth_token');
        setAuthState(prev => ({ ...prev, isLoading: false }));
      }
    } catch (error) {
      console.error('AuthContext: Auth session check failed:', error);
      // Backend not available - continue without auth for now
      console.log('AuthContext: Backend not available, continuing without auth');
      setAuthState(prev => ({ ...prev, isLoading: false }));
    }
  };

  const signIn = async () => {
    if (!address || !isConnected) {
      toast({
        title: "Wallet not connected",
        description: "Please connect your wallet first",
        variant: "destructive"
      });
      return;
    }

    setAuthState(prev => ({ ...prev, isLoading: true, error: null }));

    try {
      // Step 1: Get nonce from backend
      const nonceResponse = await fetch('/api/auth/nonce', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ address })
      });

      if (!nonceResponse.ok) {
        throw new Error('Failed to get nonce');
      }

      const { nonce } = await nonceResponse.json();

      // Step 2: Create SIWE message
      const domain = window.location.host;
      const origin = window.location.origin;
      const statement = 'Sign in with Ethereum to DeFi Portfolio';

      // Temporarily disabled SIWE message creation
      // const message = new SiweMessage({
      //   domain,
      //   address,
      //   statement,
      //   uri: origin,
      //   version: '1',
      //   chainId: 1, // Ethereum mainnet
      //   nonce,
      //   issuedAt: new Date().toISOString(),
      // });

      const messageToSign = 'Temporarily disabled for testing';

      // Step 3: Sign the message
      const signature = await signMessageAsync({
        message: messageToSign,
        account: address as `0x${string}`
      });

      // Step 4: Verify signature with backend
      const verifyResponse = await fetch('/api/auth/verify', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          message: messageToSign,
          signature,
          address
        })
      });

      if (!verifyResponse.ok) {
        throw new Error('Signature verification failed');
      }

      const { user, token } = await verifyResponse.json();

      // Step 5: Store token and update state
      localStorage.setItem('auth_token', token);
      setAuthState({
        user,
        isAuthenticated: true,
        isLoading: false,
        error: null
      });

      toast({
        title: "Successfully signed in",
        description: `Welcome back, ${user.address.slice(0, 6)}...${user.address.slice(-4)}`
      });

    } catch (error) {
      console.error('SIWE sign in failed:', error);
      const errorMessage = error instanceof Error ? error.message : 'Authentication failed';
      
      setAuthState(prev => ({
        ...prev,
        isLoading: false,
        error: errorMessage
      }));

      toast({
        title: "Sign in failed",
        description: errorMessage,
        variant: "destructive"
      });
    }
  };

  const signOut = () => {
    localStorage.removeItem('auth_token');
    setAuthState({
      user: null,
      isAuthenticated: false,
      isLoading: false,
      error: null
    });

    toast({
      title: "Signed out",
      description: "You have been signed out successfully"
    });
  };

  const linkEmail = async (email: string) => {
    if (!authState.isAuthenticated) {
      throw new Error('Must be authenticated to link email');
    }

    try {
      const token = localStorage.getItem('auth_token');
      const response = await fetch('/api/auth/link-email', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`
        },
        body: JSON.stringify({ email })
      });

      if (!response.ok) {
        throw new Error('Failed to link email');
      }

      toast({
        title: "Verification email sent",
        description: `Check your email at ${email} to complete verification`
      });

    } catch (error) {
      console.error('Email linking failed:', error);
      throw error;
    }
  };

  const verifyEmail = async (token: string) => {
    try {
      const response = await fetch('/api/auth/verify-email', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ token })
      });

      if (!response.ok) {
        throw new Error('Email verification failed');
      }

      const { user } = await response.json();
      setAuthState(prev => ({ ...prev, user }));

      toast({
        title: "Email verified",
        description: "Your email has been successfully verified"
      });

    } catch (error) {
      console.error('Email verification failed:', error);
      throw error;
    }
  };

  const contextValue: AuthContextType = {
    ...authState,
    signIn,
    signOut,
    linkEmail,
    verifyEmail
  };

  return (
    <AuthContext.Provider value={contextValue}>
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
}