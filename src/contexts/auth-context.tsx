import React, { createContext, useContext, useReducer, useEffect } from 'react';
import { toast } from 'sonner';
import { getAPIClient } from '@/lib/api/client';
import { SiweMessage } from 'siwe';

export interface User {
  id: string;
  address?: string;
  email?: string;
  isEmailVerified?: boolean;
  emailVerified?: boolean; // Add alias for compatibility
  isAdmin?: boolean;
  lastLoginAt?: string; // Add missing property
}

export interface AuthState {
  user: User | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  error: string | null;
}

export interface AuthContextType extends AuthState {
  signIn: (address?: string, signMessage?: (message: string) => Promise<string>, chainId?: number) => Promise<void>;
  signOut: () => void;
  linkEmail: (email: string) => Promise<void>;
  verifyEmail: (code: string) => Promise<void>;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

type AuthAction =
  | { type: 'AUTH_START' }
  | { type: 'AUTH_SUCCESS'; payload: User }
  | { type: 'AUTH_ERROR'; payload: string }
  | { type: 'AUTH_LOGOUT' }
  | { type: 'CLEAR_ERROR' };

const initialState: AuthState = {
  user: null,
  isAuthenticated: false,
  isLoading: false,
  error: null,
};

function authReducer(state: AuthState, action: AuthAction): AuthState {
  switch (action.type) {
    case 'AUTH_START':
      return { ...state, isLoading: true, error: null };
    case 'AUTH_SUCCESS':
      return { ...state, user: action.payload, isAuthenticated: true, isLoading: false, error: null };
    case 'AUTH_ERROR':
      return { ...state, isLoading: false, error: action.payload };
    case 'AUTH_LOGOUT':
      return { ...initialState };
    case 'CLEAR_ERROR':
      return { ...state, error: null };
    default:
      return state;
  }
}

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [authState, dispatch] = useReducer(authReducer, initialState);

  // Check for existing session on mount
  useEffect(() => {
    const checkExistingSession = async () => {
      console.log('AuthProvider: Checking existing session...');
      const token = localStorage.getItem('auth_token');
      if (token) {
        try {
          dispatch({ type: 'AUTH_START' });
          const apiClient = getAPIClient();
          const userData = await apiClient.getCurrentUser();
          dispatch({ type: 'AUTH_SUCCESS', payload: userData });
          console.log('AuthProvider: Session restored for user:', userData.id);
        } catch (error) {
          console.error('AuthProvider: Session restoration failed:', error);
          localStorage.removeItem('auth_token');
          dispatch({ type: 'AUTH_ERROR', payload: 'Session expired' });
        }
      } else {
        console.log('AuthProvider: No existing session found');
      }
    };

    checkExistingSession();
  }, []);

  const signIn = async (address?: string, signMessage?: (message: string) => Promise<string>, chainId?: number) => {
    try {
      console.log('AuthProvider: Starting SIWE sign in process...');
      dispatch({ type: 'AUTH_START' });

      // If called without parameters, show error asking user to connect wallet first
      if (!address || !signMessage) {
        throw new Error('Please connect your wallet first and try again');
      }

      const apiClient = getAPIClient();
      
      // Step 1: Get nonce from backend
      console.log('AuthProvider: Getting nonce for address:', address);
      const nonceResponse = await apiClient.getNonce(address);
      const { nonce, message: backendMessage } = nonceResponse;
      
      // Step 2: Create SIWE message
      const actualChainId = chainId || 1; // Default to Ethereum mainnet if not provided
      console.log('AuthProvider: Creating SIWE message with nonce:', nonce, 'chainId:', actualChainId);
      const siweMessage = new SiweMessage({
        domain: window.location.hostname,
        address: address,
        statement: 'Sign in to DeFi Portfolio Dashboard',
        uri: window.location.origin,
        version: '1',
        chainId: actualChainId,
        nonce: nonce,
        issuedAt: new Date().toISOString(),
      });

      const messageToSign = siweMessage.prepareMessage();
      console.log('AuthProvider: SIWE message to sign:', messageToSign);

      // Step 3: Sign the message
      console.log('AuthProvider: Requesting signature from wallet...');
      const signature = await signMessage(messageToSign);
      console.log('AuthProvider: Signature received');

      // Step 4: Verify signature with backend
      console.log('AuthProvider: Verifying signature with backend...');
      const authResponse = await apiClient.verifySiwe(messageToSign, signature);
      
      localStorage.setItem('auth_token', authResponse.token);
      dispatch({ type: 'AUTH_SUCCESS', payload: authResponse.user });

      toast.success('Successfully signed in with wallet!');
      console.log('AuthProvider: SIWE authentication successful');
    } catch (error) {
      console.error('AuthProvider: SIWE sign in failed:', error);
      let errorMessage = 'Failed to sign in';
      
      if (error instanceof Error) {
        if (error.message.includes('User rejected')) {
          errorMessage = 'Signature was rejected. Please try again.';
        } else if (error.message.includes('Invalid SIWE message format')) {
          errorMessage = 'Authentication format error. Please try again.';
        } else if (error.message.includes('rate limit')) {
          errorMessage = 'Too many attempts. Please wait a moment and try again.';
        }
      }
      
      dispatch({ type: 'AUTH_ERROR', payload: errorMessage });
      toast.error(errorMessage);
    }
  };

  const signOut = () => {
    console.log('AuthProvider: Signing out...');
    localStorage.removeItem('auth_token');
    dispatch({ type: 'AUTH_LOGOUT' });
    toast.success('Signed out successfully');
  };

  const linkEmail = async (email: string) => {
    try {
      console.log('AuthProvider: Linking email:', email);
      dispatch({ type: 'AUTH_START' });
      
      // Mock email linking
      toast.success('Verification email sent to ' + email);
      
      dispatch({ type: 'CLEAR_ERROR' });
    } catch (error) {
      console.error('AuthProvider: Email linking failed:', error);
      dispatch({ type: 'AUTH_ERROR', payload: 'Failed to link email' });
      toast.error('Failed to link email. Please try again.');
    }
  };

  const verifyEmail = async (code: string) => {
    try {
      console.log('AuthProvider: Verifying email with code:', code);
      dispatch({ type: 'AUTH_START' });

      // Mock email verification
      if (authState.user) {
        const updatedUser = { ...authState.user, isEmailVerified: true };
        dispatch({ type: 'AUTH_SUCCESS', payload: updatedUser });
      }

      toast.success('Email verified successfully!');
    } catch (error) {
      console.error('AuthProvider: Email verification failed:', error);
      dispatch({ type: 'AUTH_ERROR', payload: 'Failed to verify email' });
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